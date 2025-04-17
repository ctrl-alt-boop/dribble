package ui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/ctrl-alt-boop/gooldb/pkg/connection"
)

type (
	ConnectionItem struct {
		*connection.Settings
		Name string
	}
)

func (item ConnectionItem) FilterValue() string { return "" }
func (item ConnectionItem) Title() string       { return item.Name }
func (item ConnectionItem) Description() string { return "" }
func (item ConnectionItem) Inspect() string {
	inspectStyle := lipgloss.NewStyle().
		Align(lipgloss.Center, lipgloss.Center).
		Border(lipgloss.RoundedBorder(), true)

	return inspectStyle.Render(item.AsString())
}
