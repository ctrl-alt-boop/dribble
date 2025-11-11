package explorer

import (
	"charm.land/bubbles/v2/help"
	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"
)

const helpPanelID string = string(ExplorerPageID) + ".help"

type cheatsheet struct {
	help help.Model

	width int

	keymaps map[string]help.KeyMap
	current string

	style lipgloss.Style
}

func (h *cheatsheet) registerKeymap(id string, keymap help.KeyMap) {
	h.keymaps[id] = keymap
}

func newCheatsheet() *cheatsheet {
	return &cheatsheet{
		keymaps: map[string]help.KeyMap{},
	}
}

func (h *cheatsheet) SetWidth(width int) {
	h.width = width
}

func (h *cheatsheet) Init() tea.Cmd {
	h.help = help.New()

	return nil
}

func (h *cheatsheet) Update(msg tea.Msg) (*cheatsheet, tea.Cmd) {
	switch msg := msg.(type) {
	case focusChangedMsg:
		h.current = string(msg)
	}
	return h, nil
}

func (h *cheatsheet) Render() *lipgloss.Layer {
	h.help.Width = h.width

	style := lipgloss.NewStyle().
		Padding(0, 3).
		Faint(true).
		Width(h.width).Height(1)

	keymap, ok := h.keymaps[h.current]
	if ok {
		return lipgloss.NewLayer(style.Render(h.help.View(keymap)))
	}

	return lipgloss.NewLayer(style.Render(h.help.View(h.keymaps[string(ExplorerPageID)])))
}
