package app

import (
	"github.com/algorand/go-algorand-sdk/v2/types"
	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/algorand/node-ui/messages"
	"github.com/algorand/node-ui/tui/internal/bubbles/about"
	"github.com/algorand/node-ui/tui/internal/bubbles/accounts"
	"github.com/algorand/node-ui/tui/internal/bubbles/configs"
	"github.com/algorand/node-ui/tui/internal/bubbles/explorer"
	"github.com/algorand/node-ui/tui/internal/bubbles/status"
	"github.com/algorand/node-ui/tui/internal/bubbles/tabs"
	"github.com/algorand/node-ui/tui/internal/style"
	"github.com/algorand/node-ui/tui/internal/util"
)

type activeComponent int

const (
	explorerTab activeComponent = iota
	utilitiesTab
	accountTab
	configTab
	helpTab
)

// Model represents the top level of the TUI.
type Model struct {
	Status        tea.Model
	Accounts      tea.Model
	Tabs          tabs.Model
	BlockExplorer tea.Model
	Configs       tea.Model
	Utilities     tea.Model
	About         tea.Model
	help          help.Model

	network messages.NetworkMsg

	styles *style.Styles

	requestor *messages.Requestor

	active activeComponent
	// remember the last resize so we can re-send it when selecting a different bottom component.
	lastResize tea.WindowSizeMsg
}

// New initializes the TUI.
func New(initialWidth, initialHeight int, requestor *messages.Requestor, addresses []types.Address) Model {
	util.AppKeys.Shutdown.SetEnabled(requestor.CanShutdown())

	styles := style.DefaultStyles()
	tab := tabs.New([]string{"EXPLORER", "UTILITIES", "ACCOUNTS", "CONFIGURATION", "HELP"})
	// The tab content is the only flexible element.
	// This means the height must grow or shrink to fill the available
	// window height. It has access to the absolute height but needs to
	// be informed about the space used by other elements.
	// +1 for the help
	tabContentMargin := style.TopHeight + tab.Height() + style.FooterHeight + 1
	return Model{
		active:        explorerTab,
		styles:        styles,
		Status:        status.New(styles, requestor),
		Tabs:          tab,
		BlockExplorer: explorer.New(styles, requestor, initialWidth, 0, initialHeight, tabContentMargin),
		Configs:       configs.New(requestor, tabContentMargin),
		Accounts:      accounts.New(styles, requestor, initialHeight, tabContentMargin, addresses),
		About:         about.New(tabContentMargin, about.GetHelpContent()),
		Utilities:     about.New(tabContentMargin, about.GetUtilsContent()),
		help:          help.New(),
		requestor:     requestor,
	}
}
