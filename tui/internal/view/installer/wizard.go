package installer

import (
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"

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
}

type WizardModel struct {
	heightMargin int

	question int
	list     networkPicker

	installYesNoContent string

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
		list:         l,
	}
}

func renderYesNoContent(width, height int, configDir string, network string) string {
	r, _ := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(width),
		glamour.WithEmoji(),
	)
	padder := lipgloss.NewStyle().Height(height).PaddingLeft(7).PaddingTop(1)
	header, _ := r.Render(`
# Do you want to install a ` + network + ` node?
Press [y]es or [n]o to start the node installation.

If a node has been previously installed a software update will be attempted before it is started or restarted.

### NodeUI directory:` + configDir + `
### Network:` + network + `
`)
	return padder.Render(header)
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
		m.installYesNoContent = renderYesNoContent(m.list.w, m.list.h, m.configDir, m.list.Selected())
		for i := range m.list.networks {
			n := &m.list.networks[i]
			_, err := os.Stat(path.Join(m.configDir, n.title))
			n.present = err == nil
		}
	case tea.WindowSizeMsg:
		m.list.w = msg.Width
		m.list.h = msg.Height - m.heightMargin
		m.installYesNoContent = renderYesNoContent(m.list.w, m.list.h, m.configDir, m.list.Selected())
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
		if msg.done {
			return m, func() tea.Msg {
				return DataDirReady{DataDir: msg.datadir}
			}
		}
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
		return istyle.Render("installing!", "\n\n",
			"The node is not stopped after Node UI is closed\n",
			"To stop it yourself use the following command:\n",
			fmt.Sprintf("    %s/goal node stop -d %s\n\n", m.binDir, m.dataDir),
			"Progress: ", m.progress)
	default:
		return "unknown state"
	}
}

type installProgress struct {
	msg     string
	err     error
	done    bool
	datadir string
}

func installAndStartNodeReturnDirs(rootDir, network string) (string, string, tea.Cmd) {
	data := path.Join(rootDir, network, "algod_data")
	bin := path.Join(rootDir, network, "algod_bin")
	err := os.MkdirAll(data, 0755)
	if err != nil {
		return "", "", func() tea.Msg {
			return installProgress{msg: "failed to create data dir", err: err}
		}
	}
	err = os.MkdirAll(bin, 0755)
	if err != nil {
		return "", "", func() tea.Msg {
			return installProgress{msg: "failed to create bin dir", err: err}
		}
	}
	err = os.WriteFile(path.Join(bin, "update.sh"), []byte(updateScript), 0755)
	if err != nil {
		return "", "", func() tea.Msg {
			return installProgress{msg: "failed to write update.sh", err: err}
		}
	}

	return bin, data, tea.Batch(
		func() tea.Msg {
			return installProgress{msg: "Running update.sh", err: nil}
		},
		func() tea.Msg {
			c := exec.Command(fmt.Sprintf("%s/update.sh", bin), "-i", "-c", "stable", "-p", bin, "-d", data, "-i", "-g", network)
			out, err := c.CombinedOutput()
			return installProgress{msg: fmt.Sprintf("Finished running update.sh! \n%s\n" + string(out)), err: err, done: true, datadir: data}
		})
}
