package installer

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	heightMargin int
}

var istyle = lipgloss.NewStyle().Height(20)

func New(height, width, heightMargin int) Model {
	//constants.Keys.SetRunnable(m.runnable)
	return Model{
		heightMargin: heightMargin,
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
		istyle.Height(msg.Height - m.heightMargin)
		istyle.Width(msg.Width)
	}

	return m, nil
}

func (m Model) View() string {
	return istyle.Render("Installer")
}
