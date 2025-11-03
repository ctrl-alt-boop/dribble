package dribbler

import (
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/ctrl-alt-boop/dribbler/internal/layout"
	"github.com/ctrl-alt-boop/dribbler/internal/panel"
)

// MainContentModel contain the main ui for dribbler
type MainContentModel struct {
	PanelManager *panel.Manager

	popupContent any
}

// NewMainContentModel creates a new MainContentModel
func NewMainContentModel() MainContentModel {
	return MainContentModel{}
}

// Init implements tea.Model.
func (m *MainContentModel) Init() tea.Cmd {
	m.PanelManager = panel.NewPanelManager(
		layout.NewDockComposer(
			layout.Panels(
				layout.Panel(layout.WithPosition(panel.Bottom), layout.WithInnerHeight(1)),
				layout.Panel(layout.WithPosition(panel.Bottom), layout.WithInnerHeight(1)),
				layout.Panel(layout.WithPosition(panel.Left), layout.WithWidthRatio(0.2), layout.Focusable()),
				layout.Panel(layout.WithPosition(panel.Center), layout.Focusable()),
			),
			panel.WithStyle(lipgloss.NewStyle().Background(lipgloss.Color("155")).Foreground(lipgloss.Color("111"))),
			panel.WithFocusedStyle(lipgloss.NewStyle().Background(lipgloss.Color("179")).Foreground(lipgloss.Color("225"))),
			panel.WithPanelBorder(lipgloss.RoundedBorder()),
		),
		panel.NewPanel("bottom1"),
		panel.NewPanel("bottom2"),
		panel.NewPanel("side"),
		panel.NewPanel("center"),
	)

	return nil
}

// Update implements tea.Model.
func (m *MainContentModel) Update(msg tea.Msg) (*MainContentModel, tea.Cmd) {
	var cmd tea.Cmd
	m.PanelManager, cmd = m.PanelManager.Update(msg)

	return m, cmd
}

// View implements tea.Model.
func (m MainContentModel) View() tea.View {
	// logging.Log.Infof("MainContentModel.Render: %s", m.PanelManager.Render())
	return tea.NewView(m.PanelManager.Render())
}
