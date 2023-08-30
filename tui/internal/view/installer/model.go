package installer

import (
	"github.com/algorand/node-ui/tui/internal/bubbles/about"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/algorand/node-ui/tui/internal/view"
)

type phase int

const (
	intro phase = iota
	install
)

type Model struct {
	active       phase
	heightMargin int

	installationInfo tea.Model
	wizard           WizardModel
	help             help.Model
}

var istyle = lipgloss.NewStyle().Height(20)

func New(height, width, heightMargin int) Model {
	return Model{
		active:           intro,
		heightMargin:     heightMargin,
		help:             help.New(),
		installationInfo: about.New(heightMargin+1, GetInstallationContent()),
		wizard:           NewWizardModel(height, width, heightMargin),
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.installationInfo.Init(),
		m.wizard.Init())
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// -1 for the help
		istyle.Height(msg.Height - m.heightMargin - 1)
		istyle.Width(msg.Width)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, view.InstallerKeys.Install):
			if m.active == intro {
				m.active = install
				view.InstallerKeys.Install.SetEnabled(false)
			}
		}

	}

	var cmd tea.Cmd
	var cmds []tea.Cmd
	m.installationInfo, cmd = m.installationInfo.Update(msg)
	cmds = append(cmds, cmd)
	m.wizard, cmd = m.wizard.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	content := ""
	switch m.active {
	case intro:
		content = m.installationInfo.View()
	case install:
		content = m.wizard.View()
	default:
		content = istyle.Render("Installing...")
	}
	return lipgloss.JoinVertical(0,
		content,
		m.help.View(view.InstallerKeys))
}
