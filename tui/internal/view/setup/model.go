package setup

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/algorand/node-ui/tui/args"
	"github.com/algorand/node-ui/tui/internal/constants"
	"github.com/algorand/node-ui/tui/internal/view/app"
)

type Model struct {
	runnable bool
	app      app.Model
}

func New(args args.Arguments) (m Model) {
	requestor := getRequestorOrExit(args.AlgodDataDir, args.AlgodURL, args.AlgodToken, args.AlgodAdminToken)
	if requestor != nil {
		addresses := getAddressesOrExit(args.AddressWatchList)
		m.app = app.New(constants.InitialWidth, constants.InitialHeight, requestor, addresses)
		m.runnable = true
	}
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
	m.app, cmd = m.app.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	return m.app.View()
}
