package setup

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/algorand/node-ui/tui/args"
	"github.com/algorand/node-ui/tui/internal/bubbles/footer"
	"github.com/algorand/node-ui/tui/internal/style"
	"github.com/algorand/node-ui/tui/internal/util"
	"github.com/algorand/node-ui/tui/internal/view/app"
	"github.com/algorand/node-ui/tui/internal/view/installer"
)

type Model struct {
	runnable  bool
	configDir string

	args      args.Arguments
	installer installer.Model
	app       app.Model
	Footer    tea.Model

	sizeMsg tea.WindowSizeMsg
}

func New(args args.Arguments) (m Model) {
	requestor, err := getRequestor(args.AlgodDataDir, args.AlgodURL, args.AlgodToken, args.AlgodAdminToken)
	if err == nil {
		addresses := getAddressesOrExit(args.AddressWatchList)
		m.app = app.New(util.InitialWidth, util.InitialHeight, requestor, addresses)
		m.runnable = true
	} else {
		m.installer = installer.New(util.InitialHeight, util.InitialWidth, style.FooterHeight)
	}
	m.Footer = footer.New(style.DefaultStyles())
	return m
}

func (m Model) Init() tea.Cmd {
	cmds := []tea.Cmd{
		tea.EnterAltScreen,
		util.MakeConfigCmd,
	}

	if m.runnable {
		cmds = append(cmds, m.app.Init())
	} else {
		cmds = append(cmds, m.installer.Init())
	}

	return tea.Batch(
		cmds...,
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	// short circuit when things are runnable.
	if m.runnable {
		m.app, cmd = m.app.Update(msg)
		return m, cmd
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.sizeMsg = msg
	case util.NodeUIConfigDir:
		if msg.Err != nil {
			fmt.Fprintf(os.Stderr, "Problem fetching config dir: %v\n", msg.Err)
			return m, tea.Quit
		}
		m.configDir = msg.Dir
	case installer.DataDirReady:
		requestor, err := getRequestor(msg.DataDir, "", "", "")
		if err == nil {
			addresses := getAddressesOrExit(m.args.AddressWatchList)
			m.app = app.New(m.sizeMsg.Width, m.sizeMsg.Height, requestor, addresses)
			m.runnable = true
			return m, m.app.Init()
		}
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, util.AppKeys.Quit):
			fallthrough
		case key.Matches(msg, util.InstallerKeys.Quit):
			return m, tea.Quit
		}
	}

	m.installer, cmd = m.installer.Update(msg)
	cmds = append(cmds, cmd)

	m.Footer, cmd = m.Footer.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	var primary string
	if m.runnable {
		primary = m.app.View()
	} else {
		primary = m.installer.View()
	}
	return lipgloss.JoinVertical(0,
		primary,
		m.Footer.View())
}
