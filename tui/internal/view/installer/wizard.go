package installer

import (
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/algorand/node-ui/tui/internal/util"
)

type DataDirReady struct {
	DataDir string
}

//go:embed update.sh
var updateScript string

var enter = key.NewBinding(
	key.WithKeys("enter"),
	key.WithHelp("enter", "select"))

type WizardModel struct {
	heightMargin int

	question int
	list     list.Model

	network   int
	configDir string
	progress  string
}

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

func NewWizardModel(h, w, heightMargin int) WizardModel {
	os.UserHomeDir()
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

	return WizardModel{
		heightMargin: heightMargin,
		list:         l,
	}
}

func (m WizardModel) Init() tea.Cmd {
	return nil
}

func (m WizardModel) Update(msg tea.Msg) (WizardModel, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case util.NodeUIConfigDir:
		if msg.Err != nil {
			fmt.Fprintf(os.Stderr, "Unable to get config: %s\n", msg.Err)
			return m, tea.Quit
		}
		m.configDir = msg.Dir

	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v-m.heightMargin)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, enter):
			m.question++
		case key.Matches(msg, util.InstallerKeys.Back):
			m.question--
		case key.Matches(msg, util.InstallerKeys.Yes):
			m.question++
			cmds = append(cmds, installAndStartNode(m.configDir, m.list.SelectedItem().FilterValue()))

		case key.Matches(msg, util.InstallerKeys.No):
			return m, tea.Quit
		}
	case installProgress:
		if msg.err != nil {
			fmt.Println(msg.err)
			return m, tea.Quit
		}
		m.progress = msg.msg
		if msg.done {
			return m, func() tea.Msg {
				return DataDirReady{DataDir: msg.datadir}
			}
		}
	}

	// Need to update the submodules first to ensure the right file is "picked"
	m.list, cmd = m.list.Update(msg)
	cmds = append(cmds, cmd)

	util.InstallerKeys.Yes.SetEnabled(m.question == 2)
	util.InstallerKeys.No.SetEnabled(m.question == 2)

	return m, tea.Batch(cmds...)
}

func (m WizardModel) View() string {
	switch m.question {
	case 0:
		return m.list.View()
	case 1:
		return istyle.Render(fmt.Sprintf("Do you want to install?\n\ndata directory: %s/algod_data\nbin directory: %s/algod_bin\n\nNetwork: %s\n\n Press [y]es or [n]o.",
			m.configDir, m.configDir, m.list.SelectedItem().FilterValue()))
	default:
		return istyle.Render("installing!", "\n\n", "Progress: ", m.progress)
	}
}

type installProgress struct {
	msg     string
	err     error
	done    bool
	datadir string
}

func installAndStartNode(rootDir, network string) tea.Cmd {
	data := path.Join(rootDir, "algod_data")
	bin := path.Join(rootDir, "algod_bin")
	err := os.Mkdir(data, 0755)
	if err != nil {
		return func() tea.Msg {
			return installProgress{msg: "failed to create data dir", err: err}
		}
	}
	err = os.Mkdir(bin, 0755)
	if err != nil {
		return func() tea.Msg {
			return installProgress{msg: "failed to create bin dir", err: err}
		}
	}
	err = os.WriteFile(path.Join(bin, "update.sh"), []byte(updateScript), 0755)
	if err != nil {
		return func() tea.Msg {
			return installProgress{msg: "failed to write update.sh", err: err}
		}
	}

	return tea.Batch(
		func() tea.Msg {
			return installProgress{msg: "Running update.sh", err: nil}
		},
		func() tea.Msg {
			c := exec.Command(fmt.Sprintf("%s/update.sh", bin), "-i", "-c", "stable", "-p", bin, "-d", data, "-i", "-g", network)
			out, err := c.CombinedOutput()
			return installProgress{msg: fmt.Sprintf("Finished running update.sh! \n%s\n" + string(out)), err: err, done: true, datadir: data}
		})
}
