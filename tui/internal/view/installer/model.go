package installer

import (
	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/algorand/node-ui/tui/internal/view"
)

type Model struct {
	heightMargin int
	help         help.Model
}

var istyle = lipgloss.NewStyle().Height(20)

func New(height, width, heightMargin int) Model {
	//constants.Keys.SetRunnable(m.runnable)
	return Model{
		heightMargin: heightMargin,
		help:         help.New(),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	//var cmd tea.Cmd
	//var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// -1 for the help
		istyle.Height(msg.Height - m.heightMargin - 1)
		istyle.Width(msg.Width)
	case tea.KeyMsg:

	}

	return m, nil
}

func (m Model) View() string {
	return lipgloss.JoinVertical(0,
		istyle.Render("Installer"),
		m.help.View(view.Keys))
}
