package view

import "github.com/charmbracelet/bubbles/key"

// InstallerKeyMap contains references to all the key bindings.
type InstallerKeyMap struct {
	Generic key.Binding
	Quit    key.Binding
	Install key.Binding
	Forward key.Binding
	Back    key.Binding
	Help    key.Binding
}

// ShortHelp implements the InstallerKeyMap interface.
func (k *InstallerKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Forward, k.Back, k.Generic, k.Quit, k.Help}
}

// FullHelp implements the InstallerKeyMap interface.
func (k *InstallerKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{k.ShortHelp()}
}

// Keys is a global for accessing the InstallerKeyMap.
var Keys = &InstallerKeyMap{
	// Not sure how to group help together.
	Install: key.NewBinding(
		key.WithHelp("i", "install")),
	Generic: key.NewBinding(
		key.WithHelp("↑/↓", "navigate")),
	Help: key.NewBinding(
		key.WithHelp("?", "help")),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit")),
	Forward: key.NewBinding(
		key.WithKeys("enter", "→"),
		key.WithHelp("enter", "forwards")),
	Back: key.NewBinding(
		key.WithKeys("esc", "←"),
		key.WithHelp("esc", "backwards")),
}
