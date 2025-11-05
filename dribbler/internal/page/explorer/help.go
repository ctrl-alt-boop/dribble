package explorer

import (
	"github.com/charmbracelet/bubbles/v2/help"
	"github.com/charmbracelet/bubbles/v2/key"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
)

const helpPanelID string = string(ExplorerPageID) + ".help"

type helpKeyMap struct {
	keys []key.Binding
}

func (k helpKeyMap) ShortHelp() []key.Binding {
	return k.keys
}

func (k helpKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{k.keys}
}

type helpPanel struct {
	help help.Model

	width int

	keyMaps      []helpKeyMap
	focusedPanel int
}

func newHelpBar() *helpPanel {
	return &helpPanel{
		keyMaps: []helpKeyMap{
			{},
		},
	}
}

func (h *helpPanel) SetInnerWidth(width int) {
	h.width = width
}

func (h *helpPanel) Init() tea.Cmd {
	h.help = help.New()

	return nil
}

func (h *helpPanel) Update(msg tea.Msg) (*helpPanel, tea.Cmd) {
	return h, nil
}

func (h *helpPanel) Render() *lipgloss.Layer {
	h.help.Width = h.width

	box := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), true).
		Width(h.width).Height(3)

	return lipgloss.NewLayer(box.Render(h.help.View(h.keyMaps[h.focusedPanel])))
}
