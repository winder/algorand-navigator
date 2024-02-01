package installer

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/indent"
)

// getInstallationContent is the static markdown content when
// Node UI is first run with no configuration.
func getInstallationContent() string {
	return `
# Algorand Node UI Installer :smiley_cat:

You are about to install an Algorand Node. Congratulations!
If this is a mistake, and you already have a node you'd like
to connect to please exit and run the UI again using the -d
or -u/-t options.

When you are ready press **i** to start the installation wizard.
This will guide you through the installation process.

# What to expect
The wizard will ask you to select a network to install. You can
install multiple networks, but you will need to run the wizard
for each network you want to install.

If a network already exists, the wizard will check for an update
and then start or restart the node.

Once the node is running the UI will automatically connect to it.

# Node management
The node will not be shutdown when you exit the UI. You can
manage the node using the goal command line tool. It will be
located in your user config directory. More information about
this is provided at the end of the installation.
`
}

func renderInstallingContent(width, height int, binDir, dataDir, progress, outputBuffer string) string {
	//pad := 7
	padder := lipgloss.NewStyle().Height(height).MaxWidth(width)

	header := lipgloss.NewStyle().Bold(true).Align(lipgloss.Center).Width(width).
		Foreground(lipgloss.Color("228")).
		Background(lipgloss.Color("63"))

	paragraph := lipgloss.NewStyle().Width(width).PaddingLeft(5)
	code := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, false, false, true).
		Width(width).
		Foreground(lipgloss.Color("244")).
		PaddingLeft(10)

	pretruncate := lipgloss.NewStyle().Width(width - (7 * 2))
	bufferStyle := lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder()).
		MaxWidth(width).
		Margin(0, 5, 5, 5).Padding(1)

	var buf strings.Builder
	buf.WriteString(header.Render("\nInstallation in Progress!\n"))
	buf.WriteString(paragraph.Render(`

The node is not stopped after Node UI is closed
To stop it yourself use the following command:`))
	buf.WriteString("\n\n")
	buf.WriteString(code.Render(binDir + "/goal node stop -d " + dataDir))
	buf.WriteString("\n\n")
	buf.WriteString(paragraph.Render("Progress: " + progress))
	buf.WriteString("\n\n")
	buf.WriteString(bufferStyle.Render(pretruncate.Render(outputBuffer)))

	return padder.Render(buf.String())
}

// renderYesNoContent is the static markdown content when the user is prompted.
func renderYesNoContent(width, height int, configDir string, network string) string {
	pad := 7
	r, _ := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(width-pad*2),
		glamour.WithEmoji(),
	)
	padder := lipgloss.NewStyle().Height(height).Padding(1, pad, 0, pad).MaxWidth(width)
	header, _ := r.Render(`
# Do you want to install a ` + network + ` node?
Press [y]es or [n]o to start the node installation.

If a node has been previously installed a software update will be attempted before it is started or restarted.`)

	configB := lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#F793FF", Dark: "#AD58B4"}).
		Bold(true)

	joined := lipgloss.JoinVertical(0,
		header,
		indent.String(lipgloss.JoinVertical(0,
			"NodeUI directory: "+configB.Render(configDir),
			"         Network: "+configB.Render(network)),
			7))

	return padder.Render(joined)
}

type installProgress struct {
	output *bytes.Buffer

	msg     string
	err     error
	done    bool
	datadir string
	bindir  string
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

	var b bytes.Buffer
	// TODO: when using '-i' the node is re-installed and restarted every time
	// update.sh is called. This is much slower, but always results in a running
	// node after the command completes. It would be better to start the node
	// manually after the script is called.
	c := exec.Command(fmt.Sprintf("%s/update.sh", bin), "-c", "stable", "-p", bin, "-d", data, "-g", network, "-i")
	c.Stdout = &b
	c.Stderr = &b

	return bin, data, tea.Batch(
		func() tea.Msg {
			return installProgress{
				output: &b,
				msg:    "Running update.sh",
				err:    nil,
			}
		},
		func() tea.Msg {
			err := c.Run()
			return installProgress{
				msg:     "Finished running update.sh!",
				err:     err,
				done:    true,
				datadir: data,
				bindir:  bin,
			}
		})
}
