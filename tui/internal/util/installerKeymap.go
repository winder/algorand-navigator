package util

import "github.com/charmbracelet/bubbles/key"

// InstallerKeyMap contains references to all the key bindings.
type InstallerKeyMap struct {
	Generic    key.Binding
	Yes        key.Binding
	No         key.Binding
	CursorUp   key.Binding
	CursorDown key.Binding
	Quit       key.Binding
	Install    key.Binding
	Forward    key.Binding
	Back       key.Binding
	Help       key.Binding
}

// ShortHelp implements the InstallerKeyMap interface.
func (k *InstallerKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Install, k.Yes, k.No, k.CursorUp, k.CursorDown, k.Forward, k.Generic, k.Quit, k.Help}
}

// FullHelp implements the InstallerKeyMap interface.
func (k *InstallerKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{k.ShortHelp()}
}

// InstallerKeys is a global for accessing the InstallerKeyMap.
var InstallerKeys = &InstallerKeyMap{
	// Not sure how to group help together.
	Install: key.NewBinding(
		key.WithKeys("i"),
		key.WithHelp("i", "install")),
	CursorUp: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "up"),
	),
	CursorDown: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "down"),
	),
	Yes: key.NewBinding(
		key.WithKeys("y"),
		key.WithHelp("y", "yes"),
		key.WithDisabled()),
	No: key.NewBinding(
		key.WithKeys("n"),
		key.WithHelp("n", "no"),
		key.WithDisabled()),
	Generic: key.NewBinding(
		key.WithHelp("↑/↓", "navigate")),
	Help: key.NewBinding(
		key.WithHelp("?", "help")),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit")),
	Forward: key.NewBinding(
		key.WithKeys("enter", "→"),
		key.WithHelp("enter", "select")),
	Back: key.NewBinding(
		key.WithKeys("esc", "←"),
		key.WithHelp("esc", "backwards")),
}
