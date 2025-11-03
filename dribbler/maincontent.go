package dribbler

// import (
// 	tea "github.com/charmbracelet/bubbletea/v2"
// 	"github.com/charmbracelet/lipgloss/v2"
// 	"github.com/ctrl-alt-boop/dribbler/components"
// 	"github.com/ctrl-alt-boop/dribbler/components/data"
// 	"github.com/ctrl-alt-boop/dribbler/ui/layout"
// )

// func (m *Dribbler) createMainContent() components.ContentArea {
// 	bottom := components.New("BottomArea",
// 		layout.NewDockLayout(
// 			layout.Panels(
// 				layout.Panel(layout.Top, layout.WithHeight(4), layout.WithVerticalAlignment(lipgloss.Center), layout.Unfocusable()),
// 				layout.Panel(layout.Center, layout.Unfocusable()),
// 			),
// 			layout.WithPanelBorder(lipgloss.DoubleBorder()),
// 			layout.WithFocusedStyle(lipgloss.NewStyle().BorderForeground(lipgloss.Color("179")).Foreground(lipgloss.Color("179"))),
// 		),
// 		components.NewPromptBar(),
// 		m.help,
// 		// content.NewParamContainer(m.help, m.helpKeyMap),
// 	)

// 	models := []tea.Model{
// 		bottom, // Bottom

// 		components.New("SidePanel", layout.NewSimpleLayout(),
// 			data.NewList([]data.Item{
// 				{Value: "foo1"},
// 				{Value: "bar2"},
// 				{Value: "baz3"},
// 			}...),
// 		), // Side Panel
// 		components.New("Workspace", layout.NewSimpleLayout(),
// 			data.NewList([]data.Item{
// 				{Value: "foo1"},
// 				{Value: "bar2"},
// 				{Value: "baz3"},
// 			}...),
// 		),
// 	}

// 	dockLayout := layout.NewDockLayout(
// 		layout.Panels(
// 			layout.Panel(layout.Bottom, layout.WithHeight(6), layout.Unfocusable()),
// 			layout.Panel(layout.Left, layout.WithWidthRatio(0.25)),
// 			layout.Panel(layout.Center),
// 		),
// 		layout.WithPanelBorder(lipgloss.DoubleBorder()),
// 		layout.WithFocusedStyle(lipgloss.NewStyle().Background(lipgloss.Color("179")).Foreground(lipgloss.Color("225"))),
// 		layout.AllowNoFocus(),
// 	)
// 	// contentList := CreateDemoList(3)

// 	// layout.DebugBackgrounds = true
// 	contentArea := components.New("MainContent", dockLayout, models...)
// 	contentArea.SetStyle(lipgloss.NewStyle().Border(lipgloss.NormalBorder()))
// 	return contentArea
// }
