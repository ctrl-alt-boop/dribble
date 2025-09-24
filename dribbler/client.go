package dribbler

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (c *AppModel) RequestTableNames() tea.Cmd {
	return func() tea.Msg {
		c.dribbleClient.Request(context.w())
		return nil
	}
}
