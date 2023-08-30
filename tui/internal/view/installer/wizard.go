package installer

import (
	"os"

	"github.com/charmbracelet/bubbles/filepicker"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/algorand/node-ui/tui/internal/view"
)

var enter = key.NewBinding(
	key.WithKeys("enter"),
	key.WithHelp("enter", "select"))

type WizardModel struct {
	margines int

	question   int
	list       list.Model
	filepicker filepicker.Model

	// answers
	network    int
	installDir string // bubbles/filepicker ?

}

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

func NewWizardModel(h, w, heightMargine int) WizardModel {
	networks := []list.Item{
		item{title: "mainnet", desc: "Top banana."},
		item{title: "testnet", desc: "Assessment arena."},
		item{title: "betanet", desc: "Where bugs vacation."},
	}

	l := list.New(networks, list.NewDefaultDelegate(), w, h)
	l.Title = "Network"
	l.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{enter}
	}
	l.SetShowHelp(false)
	l.DisableQuitKeybindings()

	fp := filepicker.New()
	fp.DirAllowed = true
	fp.FileAllowed = false
	fp.CurrentDirectory, _ = os.UserHomeDir()
	fp.Height = h - heightMargine
	return WizardModel{
		margines:   heightMargine,
		list:       l,
		filepicker: fp,
	}
}

func (m WizardModel) Init() tea.Cmd {
	return m.filepicker.Init()
}

func (m WizardModel) Update(msg tea.Msg) (WizardModel, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	startQuestion := m.question
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v-m.margines)
		m.filepicker.Height = msg.Height - v - m.margines
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, enter):
			m.question++
		case key.Matches(msg, view.InstallerKeys.Back):
			m.question--
			if m.question == 1 {
				m.installDir = ""
			}
		}
	}

	// Need to update the submodules first to ensure the right file is "picked"
	m.list, cmd = m.list.Update(msg)
	cmds = append(cmds, cmd)

	// don't update the filepicker on transitions
	if startQuestion == m.question || startQuestion == 1 {
		m.filepicker, cmd = m.filepicker.Update(msg)
		cmds = append(cmds, cmd)
	}

	if didSelect, path := m.filepicker.DidSelectFile(msg); didSelect {
		m.installDir = path
	}

	return m, tea.Batch(cmds...)
}

func (m WizardModel) View() string {
	switch m.question {
	case 0:
		return m.list.View()
	case 1:
		return m.filepicker.View()
	default:
		return m.installDir + " " + m.list.SelectedItem().FilterValue()
	}
}
