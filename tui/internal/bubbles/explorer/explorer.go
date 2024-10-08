package explorer

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"time"

	table "github.com/calyptia/go-bubble-table"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/algorand/go-algorand-sdk/v2/encoding/msgpack"

	"github.com/winder/algorand-navigator/messages"
	"github.com/winder/algorand-navigator/tui/internal/style"
	"github.com/winder/algorand-navigator/tui/internal/util"
)

type state int

const (
	blockState = iota
	paysetState
	txnState
)

const initialBlocks = 25

type blocks []BlockItem
type txnItems []transactionItem

// Model for the block explorer bubble.
type Model struct {
	errCnt uint64
	err    error
	state  state

	width        int
	widthMargin  int
	height       int
	heightMargin int
	style        *style.Styles

	// for blocks page
	blocks blocks

	// listener for status command, used for recovery
	statusRound uint64

	// cache for transactions page
	transactions txnItems

	table     table.Model
	txnView   viewport.Model
	requestor *messages.Requestor
}

// New constructs the explorer Model.
func New(styles *style.Styles, requestor *messages.Requestor, width, widthMargin, height, heightMargin int) Model {
	m := Model{
		state:        blockState,
		style:        styles,
		width:        width,
		widthMargin:  widthMargin,
		height:       height,
		heightMargin: heightMargin,
		requestor:    requestor,
	}
	m.initBlocks()
	return m
}

// BlocksMsg contains new block information.
type BlocksMsg struct {
	Blocks []BlockItem
	Err    error
}

// initBlocksCmd is the initializer command.
func (m Model) initBlocksCmd() tea.Msg {
	status, err := m.requestor.Client.Status().Do(context.Background())
	if err != nil {
		return BlocksMsg{
			Err: err,
		}
	}
	return m.getBlocks(status.LastRound-initialBlocks, status.LastRound)()
}

func (m *Model) getBlocks(first, last uint64) tea.Cmd {
	return func() tea.Msg {
		var result BlocksMsg
		for i := last; i >= first; i-- {
			block, err := m.requestor.Client.BlockRaw(i).Do(context.Background())
			if err != nil {
				result.Err = err
				return result
			}
			item := BlockItem{Round: i}
			err = lenientDecode(block, &item.Block)
			if err != nil {
				result.Err = err
				return result
			}
			result.Blocks = append(result.Blocks, item)
		}
		return result
	}
}

// Init is part of the tea.Model interface.
func (m Model) Init() tea.Cmd {
	return m.initBlocksCmd
}

func lenientDecode(data []byte, objptr interface{}) error {
	return msgpack.NewLenientDecoder(bytes.NewReader(data)).Decode(&objptr)
}

func (m Model) nextBlockCmd(round uint64) tea.Cmd {
	return func() tea.Msg {
		_, err := m.requestor.Client.StatusAfterBlock(round).Do(context.Background())
		if err != nil {
			return BlocksMsg{Err: err}
		}
		blk, err := m.requestor.Client.BlockRaw(round).Do(context.Background())
		if err != nil {
			return BlocksMsg{Err: err}
		}
		item := BlockItem{Round: round}
		err = lenientDecode(blk, &item.Block)
		if err != nil {
			return err
		}

		if err != nil {
			return BlocksMsg{
				Err: err,
			}
		}
		return BlocksMsg{
			Blocks: []BlockItem{item},
		}
	}
}

func (m *Model) setSize(width, height int) {
	m.width = width
	m.height = height
	verticalFrameSize := m.style.Bottom.GetVerticalFrameSize()
	m.table.SetSize(width-m.widthMargin, height-m.heightMargin-verticalFrameSize)
	m.txnView.Width = width - m.widthMargin
	m.txnView.Height = height - m.heightMargin - lipgloss.Height(m.headerView()) - lipgloss.Height(m.footerView())
}

// Update is part of the tea.Model interface.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var updateCmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// navigate into explorer views
		switch {
		case key.Matches(msg, util.AppKeys.Forward):
			switch m.state {
			case blockState:
				// Select transactions.
				m.state = paysetState
				switch block := m.table.SelectedRow().(type) {
				case BlockItem:
					m.transactions = make([]transactionItem, 0)
					for _, txn := range block.Block.Block.Payset {
						t := txn
						m.transactions = append(m.transactions, transactionItem{&t})
					}
				}
				m.initTransactions()
			case paysetState:
				m.state = txnState
				switch txn := m.table.SelectedRow().(type) {
				case transactionItem:
					m.initTransaction(txn.SignedTxnInBlock)
				}
			}

		// navigate out of explorer views
		case key.Matches(msg, util.AppKeys.Back):
			switch m.state {
			case paysetState:
				m.state = blockState
				m.initBlocks()
			case txnState:
				m.state = paysetState
			}
		}

	case tea.WindowSizeMsg:
		m.setSize(msg.Width, msg.Height)
		return m, nil

	case messages.StatusMsg:
		m.statusRound = msg.Status.LastRound
		return m, nil

	case BlocksMsg:
		next := uint64(0)

		// report the error and wait 1 second before the next attempt.
		if msg.Err != nil {
			next = m.statusRound
			m.errCnt++
			m.err = msg.Err

			return m, tea.Tick(1*time.Second, func(time.Time) tea.Msg {
				return m.nextBlockCmd(next)()
			})
		}
		m.err = nil

		// append Blocks
		backup := m.blocks
		m.blocks = msg.Blocks
		m.blocks = append(m.blocks, backup...)
		if len(m.blocks) > 0 {
			next = m.blocks[0].Round + 1
		}
		cmds = append(cmds, m.nextBlockCmd(next))
	}

	t, tableCmd := m.table.Update(msg)
	m.table = t
	cmds = append(cmds, tableCmd)

	switch m.state {
	case blockState:
		m, updateCmd = m.updateBlocks(msg)
		return m, tea.Batch(append(cmds, updateCmd)...)
	case paysetState:
		return m, nil
	case txnState:
		m.txnView, updateCmd = m.txnView.Update(msg)
		return m, tea.Batch(append(cmds, updateCmd)...)
	}

	return m, nil
}

// View is part of the tea.Model interface.
func (m Model) View() string {
	prefix := ""
	if m.err != nil {
		var bldr strings.Builder
		bldr.WriteString("\n")
		bldr.WriteString("\n")
		bldr.WriteString(fmt.Sprintf("Error(%d): %s\n", m.errCnt, strings.ReplaceAll(m.err.Error(), "\n", "")))
		for i := 0; i < (m.height - m.heightMargin - 3 - 1); i++ {
			bldr.WriteString("\n")
		}
		return bldr.String()
	}
	switch m.state {
	case blockState, paysetState:
		return prefix + m.style.Bottom.Render(m.table.View())
	case txnState:
		return prefix + m.viewTransaction()
	}
	return ""
}
