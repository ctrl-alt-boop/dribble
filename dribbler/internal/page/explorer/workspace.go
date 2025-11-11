package explorer

import (
	"fmt"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"
	"github.com/ctrl-alt-boop/dribbler/internal/dribbleapi"
	"github.com/ctrl-alt-boop/dribbler/logging"
)

const workspaceID string = string(ExplorerPageID) + ".workspace"

type workspace struct {
	width, height           int
	innerWidth, innerHeight int

	tabs       []string // type?
	currentTab int

	style lipgloss.Style

	keybinds *WorkspaceKeys
}

func (w *workspace) SetSize(width, height int) {
	w.width, w.height = width, height
	w.innerWidth, w.innerHeight = width-w.style.GetHorizontalFrameSize(), height-w.style.GetVerticalFrameSize()
}

func newWorkspace() *workspace {
	return &workspace{
		keybinds:   DefaultWorkspaceKeyBindings(),
		tabs:       []string{},
		currentTab: 0,
	}
}

func (w *workspace) Init() tea.Cmd {
	return nil
}

func (w *workspace) Update(msg tea.Msg) (*workspace, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		logging.GlobalLogger().Infof("workspace.KeyPressMsg: %s", msg)
		switch {
		case key.Matches(msg, w.keybinds.NextTab):
			w.currentTab++
			if w.currentTab >= len(w.tabs) {
				w.currentTab = 0
			}
		case key.Matches(msg, w.keybinds.PrevTab):
			w.currentTab--
			if w.currentTab < 0 {
				w.currentTab = len(w.tabs) - 1
			}
		}

	case dribbleapi.DribbleResponseMsg:
		logging.GlobalLogger().Infof("workspace.DribbleResponseMsg: %s", msg)
		w.AddTab(fmt.Sprintf("%T: %v", msg.Response, msg.Response))
	}
	return w, nil
}

func (w *workspace) AddTab(tab string) {
	w.tabs = append(w.tabs, tab)
	w.currentTab = len(w.tabs) - 1
}

func (w *workspace) SetStyle(style lipgloss.Style) {
	w.style = style
}

func (w *workspace) Render() *lipgloss.Layer {
	box := w.style.
		Width(w.width).Height(w.height)

	tab := ""
	if w.currentTab < len(w.tabs) {
		tab = w.tabs[w.currentTab]
	}
	logging.GlobalLogger().Infof("workspace.tab %v, tabs: %v", tab, w.tabs)
	return lipgloss.NewLayer(box.Render(tab))
}
