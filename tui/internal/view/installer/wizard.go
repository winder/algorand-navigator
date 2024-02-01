package installer

import (
	_ "embed"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/algorand/node-ui/tui/internal/util"
)

//go:embed update.sh
var updateScript string

const (
	networkQuestion = iota
	installQuestion
	installing
)

// DataDirReady is sent when once the node has been installed and started.
type DataDirReady struct {
	DataDir string
	BinDir  string
}

type WizardModel struct {
	heightMargin int
	width        int
	height       int

	question int
	list     networkPicker

	installYesNoContent string
	installingContent   string
	outputBuffer        string

	network   int
	configDir string
	binDir    string
	dataDir   string
	progress  string
}

func NewWizardModel(h, w, heightMargin int) WizardModel {
	networks := []networkItem{
		{title: "mainnet", desc: "Top banana."},
		{title: "testnet", desc: "Assessment arena."},
		{title: "betanet", desc: "Where bugs vacation."},
	}

	l := NewNetworkPicker(w, h, heightMargin, networks...)

	return WizardModel{
		heightMargin: heightMargin,
		width:        w,
		height:       h,
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
		m.installYesNoContent = renderYesNoContent(m.width, m.height, m.configDir, m.list.Selected())
		for i := range m.list.networks {
			n := &m.list.networks[i]
			_, err := os.Stat(path.Join(m.configDir, n.title))
			n.present = err == nil
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height - m.heightMargin
		m.installYesNoContent = renderYesNoContent(m.width, m.height, m.configDir, m.list.Selected())
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, util.InstallerKeys.Forward):
			m.question++
		case key.Matches(msg, util.InstallerKeys.Back):
			if m.question > 0 {
				m.question--
			}
		case key.Matches(msg, util.InstallerKeys.Yes):
			m.question++
			m.binDir, m.dataDir, cmd = installAndStartNodeReturnDirs(m.configDir, m.list.Selected())
			cmds = append(cmds, cmd)
		case key.Matches(msg, util.InstallerKeys.No):
			return m, tea.Quit
		}
	case installProgress:
		if msg.err != nil {
			fmt.Fprintf(os.Stderr, "A problem occurred during installation: %s\n", msg.err)
			return m, tea.Quit
		}

		m.progress = msg.msg
		m.outputBuffer = msg.output.String()
		if msg.done {
			return m, tea.Tick(2*time.Second, func(_ time.Time) tea.Msg {
				return DataDirReady{DataDir: msg.datadir, BinDir: msg.bindir}
			})
		}
		m.installingContent = renderInstallingContent(m.width, m.height, m.binDir, m.dataDir, m.progress, m.outputBuffer)

		return m, tea.Tick(50*time.Millisecond, func(t time.Time) tea.Msg {
			// come back in a moment to append more output.
			return installProgress{
				msg:    msg.msg,
				err:    msg.err,
				output: msg.output,
			}
		})
	}

	// Need to update the submodules first to ensure the right file is "picked"
	m.list, cmd = m.list.Update(msg)
	cmds = append(cmds, cmd)

	util.InstallerKeys.Back.SetEnabled(m.question > 0 && m.question < installing)
	util.InstallerKeys.Forward.SetEnabled(m.question < installQuestion)
	util.InstallerKeys.Yes.SetEnabled(m.question == installQuestion)
	util.InstallerKeys.No.SetEnabled(m.question == installQuestion)
	util.InstallerKeys.CursorUp.SetEnabled(m.question == networkQuestion)
	util.InstallerKeys.CursorDown.SetEnabled(m.question == networkQuestion)

	return m, tea.Batch(cmds...)
}

func (m WizardModel) View() string {
	switch m.question {
	case networkQuestion:
		return m.list.View()
	case installQuestion:
		return m.installYesNoContent
	case installing:
		return m.installingContent
	default:
		return "unknown state"
	}
}
