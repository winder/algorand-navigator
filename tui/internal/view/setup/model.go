package setup

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/algorand/node-ui/tui/args"
	"github.com/algorand/node-ui/tui/internal/bubbles/footer"
	"github.com/algorand/node-ui/tui/internal/constants"
	"github.com/algorand/node-ui/tui/internal/style"
	"github.com/algorand/node-ui/tui/internal/view/app"
)

type Model struct {
	runnable bool
	help     help.Model
	app      app.Model
	Footer   tea.Model
}

func New(args args.Arguments) (m Model) {
	requestor, err := getRequestor(args.AlgodDataDir, args.AlgodURL, args.AlgodToken, args.AlgodAdminToken)
	if err == nil {
		addresses := getAddressesOrExit(args.AddressWatchList)
		m.app = app.New(constants.InitialWidth, constants.InitialHeight, requestor, addresses)
		m.runnable = true
	}
	m.Footer = footer.New(style.DefaultStyles())
	m.help = help.New()
	constants.Keys.SetRunnable(m.runnable)
	return m
}

func (m Model) Init() tea.Cmd {

	if m.runnable {
		return tea.Batch(
			tea.EnterAltScreen,
			m.app.Init(),
		)
	}

	// else.... installer
	return tea.Batch(
		tea.EnterAltScreen,
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, constants.Keys.Quit):
			return m, tea.Quit
		}
	}

	m.help, cmd = m.help.Update(msg)
	cmds = append(cmds, cmd)

	m.Footer, cmd = m.Footer.Update(msg)
	cmds = append(cmds, cmd)

	if m.runnable {
		m.app, cmd = m.app.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	var primary string
	if m.runnable {
		primary = m.app.View()
	} else {
		primary = "Installer"
	}
	return lipgloss.JoinVertical(0,
		primary,
		m.help.View(constants.Keys),
		m.Footer.View())
}
