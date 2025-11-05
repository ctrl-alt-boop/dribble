package explorer

import (
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
)

const workspacePanelID string = string(ExplorerPageID) + ".workspace"

type workspacePanel struct {
	innerWidth, innerHeight int
}

func (w *workspacePanel) SetInnerSize(width, height int) {
	w.innerWidth, w.innerHeight = width, height
}

func newWorkspace() *workspacePanel {
	return &workspacePanel{}
}

func (w *workspacePanel) Init() tea.Cmd {
	return nil
}

func (w *workspacePanel) Update(msg tea.Msg) (*workspacePanel, tea.Cmd) {
	return w, nil
}

func (w *workspacePanel) Render() *lipgloss.Layer {
	box := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), true).
		Width(w.innerWidth).Height(w.innerHeight)

	return lipgloss.NewLayer(box.Render("boox"))
}
