package installer

import (
	"github.com/algorand/node-ui/tui/internal/style"
	"github.com/algorand/node-ui/tui/internal/util"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"strings"
)

var defaultstyle = style.DefaultStyles()

type networkPicker struct {
	networks             []networkItem
	selected             int
	w, h, verticalMargin int

	header              string
	printer             lipgloss.Style
	listPad             lipgloss.Style
	selectedLine        lipgloss.Style
	selectedLine2       lipgloss.Style
	selectedPresentLine lipgloss.Style
	nonSelectedLine     lipgloss.Style
}

type networkItem struct {
	title, desc string
	present     bool
}

func NewNetworkPicker(width, height, verticalMargin int, items ...networkItem) networkPicker {
	n := networkPicker{
		networks:       items,
		w:              width,
		h:              height,
		verticalMargin: verticalMargin,
	}

	n.printer = lipgloss.NewStyle().Height(height - verticalMargin)
	n.selectedLine = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, false, false, true).
		Bold(true).
		BorderForeground(lipgloss.AdaptiveColor{Light: "#F793FF", Dark: "#AD58B4"}).
		Foreground(lipgloss.AdaptiveColor{Light: "#EE6FF8", Dark: "#EE6FF8"}).
		Padding(0, 0, 0, 1)
	n.selectedLine2 = n.selectedLine.Copy().UnsetBold()
	n.selectedPresentLine = n.selectedLine2.Copy().
		Foreground(lipgloss.AdaptiveColor{Light: "226", Dark: "228"})
	n.nonSelectedLine = lipgloss.NewStyle().Foreground(lipgloss.Color("240")).PaddingLeft(2)
	n.listPad = lipgloss.NewStyle().PaddingLeft(9)

	r, _ := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(width),
		glamour.WithEmoji(),
	)
	padder := lipgloss.NewStyle().PaddingLeft(7).PaddingTop(1)
	n.header, _ = r.Render(`
# Network Selector
Choose one of the long running networks to install.
`)
	n.header = padder.Render(n.header)

	return n
}

func (n networkPicker) Selected() string {
	return n.networks[n.selected].title
}

func (n networkPicker) Init() tea.Cmd {
	return nil
}

func (n networkPicker) Update(msg tea.Msg) (networkPicker, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		n.w = msg.Width
		n.h = msg.Height
		n.printer.Height(n.h - n.verticalMargin)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, util.InstallerKeys.CursorUp):
			if n.selected-1 >= 0 {
				n.selected--
			}
		case key.Matches(msg, util.InstallerKeys.CursorDown):
			if n.selected+1 < len(n.networks) {
				n.selected++
			}
		}
	}
	return n, nil
}

func (n networkPicker) View() string {
	var bldr strings.Builder

	for i, network := range n.networks {
		var s1, s2, s3 lipgloss.Style
		if i == n.selected {
			s1 = n.selectedLine
			s2 = n.selectedLine2
			s3 = n.selectedPresentLine
		} else {
			s1 = n.nonSelectedLine
			s2 = s1
			s3 = s1
		}
		bldr.WriteString(s1.Render(network.title))
		bldr.WriteString("\n")

		bldr.WriteString(s2.Render(network.desc))
		bldr.WriteString("\n")

		if network.present {
			bldr.WriteString(s3.Render("* pre-existing config directory *"))
			bldr.WriteString("\n")
		}
		bldr.WriteString("\n")
	}

	return n.printer.Render(lipgloss.JoinVertical(0, n.header, n.listPad.Render(bldr.String())))
}
