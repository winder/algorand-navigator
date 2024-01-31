package util

import "github.com/charmbracelet/bubbles/key"

// AppKeyMap contains references to all the key bindings.
type AppKeyMap struct {
	Generic      key.Binding
	Quit         key.Binding
	Catchup      key.Binding
	AbortCatchup key.Binding
	Section      key.Binding
	Forward      key.Binding
	Back         key.Binding
	Help         key.Binding
}

// ShortHelp implements the AppKeyMap interface.
func (k *AppKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Section, k.Forward, k.Back, k.Generic, k.Quit, k.Help}
}

// FullHelp implements the AppKeyMap interface.
func (k *AppKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{k.ShortHelp()}
}

// AppKeys is a global for accessing the AppKeyMap.
var AppKeys = &AppKeyMap{
	// Not sure how to group help together.
	Generic: key.NewBinding(
		key.WithHelp("↑/↓", "navigate")),
	Help: key.NewBinding(
		key.WithHelp("?", "help")),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit")),
	Catchup: key.NewBinding(
		key.WithKeys("f"),
		key.WithHelp("f", "start fast catchup")),
	AbortCatchup: key.NewBinding(
		key.WithKeys("a"),
		key.WithHelp("a", "abort catchup")),
	Section: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "section")),
	Forward: key.NewBinding(
		key.WithKeys("enter", "→"),
		key.WithHelp("enter", "forwards")),
	Back: key.NewBinding(
		key.WithKeys("esc", "←"),
		key.WithHelp("esc", "backwards")),
}
