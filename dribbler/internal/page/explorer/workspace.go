package explorer

import (
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"
	"github.com/ctrl-alt-boop/dribbler/internal/dribbleapi"
	"github.com/ctrl-alt-boop/dribbler/logging"
)

const workspaceID string = string(ExplorerPageID) + ".workspace"

type workspace struct {
	width, height int

	tabs *Tabs

	style lipgloss.Style

	keybinds *WorkspaceKeys
}

func (w *workspace) SetSize(width, height int) {
	w.width, w.height = width, height
}

func newWorkspace() *workspace {
	return &workspace{
		keybinds: DefaultWorkspaceKeyBindings(),
		tabs:     NewTabs(),
	}
}

func (w *workspace) Init() tea.Cmd {
	logging.GlobalLogger().Infof("workspace.Init")

	w.tabs.Add(&StringTab{
		name:    "Tabu oneo",
		content: "This is one mighty fine tab\nKABOOM!",
	})
	w.tabs.Add(&StringTab{
		name:    "Tabu twoo",
		content: "This is the second mighty fine tab\nKABOOM CHakalaka!",
	})
	w.tabs.Add(&StringTab{
		name:    "Tabu threeo",
		content: "This is third mighty fine tab\nKABOOMISHOOKABOOM!",
	})
	w.tabs.Add(&StringTab{
		name:    "Tabu fouro",
		content: "This is fourth mighty fine tab\nKABOOMBABOOM!",
	})

	return nil
}

func (w *workspace) Update(msg tea.Msg) (*workspace, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		logging.GlobalLogger().Infof("workspace.KeyPressMsg: %s", msg)
		switch {
		case key.Matches(msg, w.keybinds.NextTab):
			w.tabs.Move(1)

		case key.Matches(msg, w.keybinds.PrevTab):
			w.tabs.Move(-1)
		}

	case dribbleapi.DribbleResponseMsg:
		logging.GlobalLogger().Infof("workspace.DribbleResponseMsg: %s", msg)
		// w.AddTab(typeResString(msg))
	}
	return w, nil
}

func (w *workspace) AddTab(tab Tab) {
	w.tabs.Add(tab)
}

func (w *workspace) SetStyle(style lipgloss.Style) {
	w.style = style
}

func (w *workspace) Render() *lipgloss.Layer {
	current, bar := w.tabs.Render()

	view := lipgloss.JoinVertical(lipgloss.Left, bar, current)
	style := w.style.
		Width(w.width).Height(w.height)

	return lipgloss.NewLayer(style.Render(view))
}
