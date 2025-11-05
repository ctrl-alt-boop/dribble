package components

import (
	"github.com/charmbracelet/bubbles/v2/help"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/ctrl-alt-boop/dribbler/keys"

	"github.com/charmbracelet/lipgloss/v2"
)

var helpStyle = lipgloss.NewStyle().
	Align(lipgloss.Left, lipgloss.Center).
	PaddingLeft(1)

// Help is a wrapper around help.Model
type Help struct {
	AlwaysUpdate

	help help.Model
	Keys keys.KeyMap
}

// Init implements tea.Model.
func (h Help) Init() tea.Cmd {
	h.help.FullSeparator = " \u2502 "
	h.help.ShortSeparator = " \u2502 "
	return nil
}

// NewHelp creates a new Help widget
func NewHelp() Help {
	return Help{
		help: help.New(),
		Keys: keys.Map,
	}
}

// Update implements tea.Model.
func (h Help) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	updated := h
	updated.help, cmd = h.help.Update(msg)

	return updated, cmd
}

func (h Help) Render() string {
	return helpStyle.Width(h.help.Width).Height(1).Render(h.help.View(h.Keys))
}

// View implements tea.Model.
func (h Help) View() tea.View {
	return tea.NewView(h.Render())
}
