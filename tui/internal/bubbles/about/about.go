package about

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

// Model represents the about bubble.
type Model struct {
	heightMargin int
	viewport     viewport.Model
	content      string
}

// New creates the about Model.
func New(heightMargin int, content string) Model {
	m := Model{
		heightMargin: heightMargin,
		viewport:     viewport.New(0, 0),
		content:      content,
	}
	m.setSize(80, 20)
	m.viewport.SetContent(render(76, content, 1, 7, 0, 7))
	return m
}

func render(wrap int, content string, padding ...int) string {
	r, _ := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(wrap-14),
		glamour.WithEmoji(),
	)
	c, _ := r.Render(content)
	return lipgloss.NewStyle().Padding(padding...).Render(c)
}

// Init is part of the tea.Model interface.
func (m Model) Init() tea.Cmd {
	return nil
}

// Update is part of the tea.Model interface.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.setSize(msg.Width, msg.Height)
		m.viewport.SetContent(render(msg.Width, m.content, 1, 7, 0, 7))
	}
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

// View is part of the tea.Model interface.
func (m Model) View() string {
	return m.viewport.View()
}

func (m *Model) setSize(width, height int) {
	m.viewport.Width = width
	m.viewport.Height = height - m.heightMargin
}
