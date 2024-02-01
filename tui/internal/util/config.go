package util

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/kirsle/configdir"
)

type NavigatorUIConfigDir struct {
	Dir string
	Err error
}

// MakeConfigCmd creates the config directory.
func MakeConfigCmd() tea.Msg {
	configPath := configdir.LocalConfig("algorand-navigator")
	err := configdir.MakePath(configPath) // Ensure it exists.
	return NavigatorUIConfigDir{
		Dir: configPath,
		Err: err,
	}
}
