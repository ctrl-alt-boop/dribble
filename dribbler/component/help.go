package component

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ctrl-alt-boop/dribbler/config"
	"github.com/ctrl-alt-boop/dribbler/ui"
)

// Help is a wrapper around help.Model
type Help struct {
	AlwaysUpdate

	help help.Model
	Keys config.KeyMap
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
		Keys: config.Keys,
	}
}

// Update implements tea.Model.
func (h Help) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	updated := h
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, config.Keys.Help):
			updated.help.ShowAll = !h.help.ShowAll
		}
	}
	return updated, nil
}

// View implements tea.Model.
func (h Help) View() string {
	return ui.HelpStyle.Width(h.help.Width).Height(1).Render(h.help.View(h.Keys))
}
