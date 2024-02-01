package app

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/winder/algorand-navigator/messages"
	"github.com/winder/algorand-navigator/tui/internal/util"
)

func networkFromID(genesisID string) string {
	return strings.Split(genesisID, "-")[0]
}

// Update is part of the tea.Model interface.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case messages.NetworkMsg:
		m.network = msg

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, util.AppKeys.Quit):
			return m, tea.Quit
		case key.Matches(msg, util.AppKeys.Catchup):
			return m, m.requestor.StartFastCatchup(networkFromID(m.network.GenesisID))
		case key.Matches(msg, util.AppKeys.AbortCatchup):
			return m, m.requestor.StopFastCatchup(networkFromID(m.network.GenesisID))
		case key.Matches(msg, util.AppKeys.Shutdown):
			// trigger shutdown from a different level.
			return m, func() tea.Msg {
				return messages.MakeStopNodeMsg(m.requestor)
			}
		case key.Matches(msg, util.AppKeys.Section):
			m.active++
			m.active %= 5
			m.Tabs.SetActiveIndex(int(m.active))
			return m, nil
		}
		switch m.active {
		case explorerTab:
			var explorerCommand tea.Cmd
			m.BlockExplorer, explorerCommand = m.BlockExplorer.Update(msg)
			return m, explorerCommand
		case accountTab:
		case configTab:
		case helpTab:
		case utilitiesTab:
		}

	case tea.WindowSizeMsg:
		m.lastResize = msg
	}

	m.Status, cmd = m.Status.Update(msg)
	cmds = append(cmds, cmd)

	m.Accounts, cmd = m.Accounts.Update(msg)
	cmds = append(cmds, cmd)

	m.BlockExplorer, cmd = m.BlockExplorer.Update(msg)
	cmds = append(cmds, cmd)

	m.Configs, cmd = m.Configs.Update(msg)
	cmds = append(cmds, cmd)

	m.Tabs, cmd = m.Tabs.Update(msg)
	cmds = append(cmds, cmd)

	m.About, cmd = m.About.Update(msg)
	cmds = append(cmds, cmd)

	m.Utilities, cmd = m.Utilities.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}
