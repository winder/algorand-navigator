package setup

import (
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/algorand/node-ui/messages"
	"github.com/algorand/node-ui/tui/args"
	"github.com/algorand/node-ui/tui/internal/bubbles/footer"
	"github.com/algorand/node-ui/tui/internal/style"
	"github.com/algorand/node-ui/tui/internal/util"
	"github.com/algorand/node-ui/tui/internal/view/app"
	"github.com/algorand/node-ui/tui/internal/view/installer"
)

type setupState int

const (
	installerState setupState = iota + 1
	appState
	shutdownState
)

type Model struct {
	state     setupState
	configDir string

	args args.Arguments

	installer installer.Model
	app       app.Model
	shutdown  string

	Footer tea.Model

	sizeMsg tea.WindowSizeMsg
}

func New(args args.Arguments) (m Model) {
	requestor, err := getRequestor(args.AlgodDataDir, args.AlgodBinDir, args.AlgodURL, args.AlgodToken, args.AlgodAdminToken)
	if err == nil {
		addresses := getAddressesOrExit(args.AddressWatchList)
		m.app = app.New(util.InitialWidth, util.InitialHeight, requestor, addresses)
		m.state = appState

	} else {
		m.installer = installer.New(util.InitialHeight, util.InitialWidth, style.FooterHeight)
		m.state = installerState
	}
	m.Footer = footer.New(style.DefaultStyles())
	return m
}

func (m Model) Init() tea.Cmd {
	cmds := []tea.Cmd{
		tea.EnterAltScreen,
		util.MakeConfigCmd,
	}

	switch m.state {
	case appState:
		cmds = append(cmds, m.app.Init())
	case installerState:
		cmds = append(cmds, m.installer.Init())
	}

	return tea.Batch(cmds...)
}

type nodeShutdownComplete int

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.sizeMsg = msg
	case tea.KeyMsg:
		// only handle input with the appropriate view
		switch m.state {
		case appState:
			m.app, cmd = m.app.Update(msg)
			return m, cmd
		case installerState:
			m.installer, cmd = m.installer.Update(msg)
			return m, cmd
		}
	case util.NodeUIConfigDir:
		if msg.Err != nil {
			fmt.Fprintf(os.Stderr, "Problem fetching config dir: %v\n", msg.Err)
			return m, tea.Quit
		}
		m.configDir = msg.Dir
	// message from installer
	case installer.DataDirReady:
		requestor, err := getRequestor(msg.DataDir, msg.BinDir, "", "", "")
		if err == nil {
			addresses := getAddressesOrExit(m.args.AddressWatchList)
			m.app = app.New(m.sizeMsg.Width, m.sizeMsg.Height, requestor, addresses)
			m.state = appState
			return m, m.app.Init()
		}
	case nodeShutdownComplete:
		// TODO: go back to the installer?
		//m.installer = installer.New(m.sizeMsg.Height, m.sizeMsg.Width, style.FooterHeight)
		//m.state = installerState
		return m, tea.Quit
	case messages.StopNodeResult:
		if msg.Err != nil {
			m.shutdown = fmt.Sprintf("Error stopping node: %v", msg.Err)
		} else {
			m.shutdown = "Node stopped."
		}

		// show result for 5 seconds before transitioning
		return m, tea.Tick(5*time.Second, func(t time.Time) tea.Msg {
			return nodeShutdownComplete(0)
		})

	case messages.StopNodeMsg:
		m.state = shutdownState
		m.shutdown = "Stopping node..."
		return m, msg.Stop()
	}

	if m.state == appState {
		m.app, cmd = m.app.Update(msg)
		cmds = append(cmds, cmd)
	}

	m.installer, cmd = m.installer.Update(msg)
	cmds = append(cmds, cmd)

	m.Footer, cmd = m.Footer.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	var primary string
	switch m.state {
	case appState:
		primary = m.app.View()
	case installerState:
		primary = m.installer.View()
	case shutdownState:
		primary = m.shutdown
	}
	return lipgloss.JoinVertical(0,
		primary,
		m.Footer.View())
}
