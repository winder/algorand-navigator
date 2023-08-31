package util

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/kirsle/configdir"
)

type NodeUIConfigDir struct {
	Dir string
	Err error
}

// MakeConfigCmd creates the config directory.
func MakeConfigCmd() tea.Msg {
	configPath := configdir.LocalConfig("nodeui")
	err := configdir.MakePath(configPath) // Ensure it exists.
	return NodeUIConfigDir{
		Dir: configPath,
		Err: err,
	}
}
