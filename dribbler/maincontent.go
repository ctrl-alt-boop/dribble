package dribbler

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ctrl-alt-boop/dribbler/ui/content"
	"github.com/ctrl-alt-boop/dribbler/ui/layout"
	"github.com/ctrl-alt-boop/dribbler/widget"
)

func CreateMainContent() widget.ContentArea {
	models := []tea.Model{
		// widget.New(
		// 	"bot",
		// 	layout.NewDockLayout(
		// 		layout.Panels(
		// 			layout.Panel(layout.Center),
		// 			layout.Panel(layout.Bottom, layout.WithHeight(3)),
		// 		),
		// 		layout.WithPanelBorder(lipgloss.NormalBorder()),
		// 	),
		// 	content.Text("PromptBar will be here."),
		// 	//huh.NewInput().Prompt(">"), // prompt bar
		// 	content.Text("HelpText will be here, i promise!"),
		// ),
		widget.New("SidePanel", layout.NewSimpleLayout(),
			content.NewList([]content.Item{
				{Value: "foo1"},
				{Value: "bar2"},
				{Value: "baz3"},
			}),
		), // Side Panel

		widget.New("SidePanel", layout.NewSimpleLayout(),
			content.NewList([]content.Item{
				{Value: "foo1"},
				{Value: "bar2"},
				{Value: "baz3"},
			}),
		), // Side Panel
		widget.New("Workspace", layout.NewSimpleLayout(),
			content.NewList([]content.Item{
				{Value: "foo1"},
				{Value: "bar2"},
				{Value: "baz3"},
			}),
		),
	}
	dockLayout := layout.NewDockLayout(
		layout.Panels(
			layout.Panel(layout.Bottom, layout.WithHeight(6)),
			layout.Panel(layout.Left, layout.WithWidthRatio(0.25)),
			layout.Panel(layout.Center),
		),
		layout.WithPanelBorder(lipgloss.DoubleBorder()),
	)

	layout.DebugBackgrounds = true
	contentArea := widget.New("dribbler", dockLayout, models...)
	contentArea.SetStyle(lipgloss.NewStyle().Border(lipgloss.NormalBorder()))
	return contentArea
}
