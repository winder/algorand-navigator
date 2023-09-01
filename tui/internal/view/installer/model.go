package installer

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/algorand/node-ui/tui/internal/bubbles/about"
	"github.com/algorand/node-ui/tui/internal/util"
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
		installationInfo: about.New(heightMargin+1, getInstallationContent()),
		wizard:           NewWizardModel(height, width, heightMargin+1), // add 1 for the help
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.installationInfo.Init(),
		m.wizard.Init())
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	keyMsg := false
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		istyle.Height(msg.Height - m.heightMargin)
		istyle.Width(msg.Width)
	case tea.KeyMsg:
		keyMsg = true
		switch {
		case key.Matches(msg, util.InstallerKeys.Install):
			if m.active == intro {
				m.active = install
				util.InstallerKeys.Install.SetEnabled(false)
			}
		}
	}

	if !keyMsg || m.active == intro {
		m.installationInfo, cmd = m.installationInfo.Update(msg)
		cmds = append(cmds, cmd)
	}

	if !keyMsg || m.active == install {
		m.wizard, cmd = m.wizard.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	switch m.active {
	case intro:
		help := m.help.Styles.ShortKey.Inline(true).Render("i") + " " +
			m.help.Styles.ShortDesc.Inline(true).Render("install")

		return lipgloss.JoinVertical(0,
			m.installationInfo.View(),
			help)
	case install:
		return lipgloss.JoinVertical(0,
			m.wizard.View(),
			m.help.View(util.InstallerKeys))
	default:
		return istyle.Render("Installing...")
	}
}
