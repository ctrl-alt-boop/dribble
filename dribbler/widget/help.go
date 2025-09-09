package widget

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ctrl-alt-boop/dribbler/config"
	"github.com/ctrl-alt-boop/dribbler/ui"
)

type Help struct {
	help help.Model
	Keys config.KeyMap
}

// Init implements tea.Model.
func (h *Help) Init() tea.Cmd {
	h.help.FullSeparator = " \u2502 "
	h.help.ShortSeparator = " \u2502 "
	return nil
}

func NewHelp() *Help {
	return &Help{
		help: help.New(),
		Keys: config.Keys,
	}
}

func (h *Help) FocusChanged(focus Kind) {
	switch focus {
	case KindPanel:
		h.Keys.FullHelpFunc = fullHelpPanelFunc
	case KindWorkspace:
		h.Keys.FullHelpFunc = fullHelpWorkspaceFunc
	}
}

func (h *Help) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h.UpdateSize(msg.Width, msg.Height)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, config.Keys.Help):
			h.help.ShowAll = !h.help.ShowAll
		}
	}
	return h, nil
}

func (h *Help) UpdateSize(width int, _ int) {
	h.help.Width = width - ui.HelpStyle.GetHorizontalFrameSize()
}

func (h *Help) View() string {

	return ui.HelpStyle.Width(h.help.Width).Height(1).Render(h.help.View(h.Keys))
}
