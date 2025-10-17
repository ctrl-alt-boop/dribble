package dribbler

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ctrl-alt-boop/dribbler/component"
	"github.com/ctrl-alt-boop/dribbler/config"
	"github.com/ctrl-alt-boop/dribbler/ui/content"
	"github.com/ctrl-alt-boop/dribbler/ui/layout"
)

var _ tea.Model = (*DemoModel)(nil)

// DemoModel is used by the demo application to showcase UI
type DemoModel struct {
	thing                   tea.Model
	Width, Height           int
	InnerWidth, InnerHeight int
	ViewIndex               int
}

// Init implements tea.Model.
func (t *DemoModel) Init() tea.Cmd {
	layout.DebugBackgrounds = true
	return nil
}

// Update implements tea.Model.
func (t *DemoModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:

		t.Width, t.Height = msg.Width, msg.Height
		// borderStyle := lipgloss.NewStyle().Border(lipgloss.NormalBorder())
		// t.InnerWidth, t.InnerHeight = msg.Width-borderStyle.GetHorizontalFrameSize(), msg.Height-borderStyle.GetVerticalFrameSize()
		updatedThing, cmd := t.thing.Update(tea.WindowSizeMsg{Width: t.Width, Height: t.Height})
		t.thing = updatedThing
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
		return t, tea.Batch(cmds...)
	case tea.KeyMsg:
		if key.Matches(msg, config.Keys.Quit) {
			return t, tea.Quit
		}
		if key.Matches(msg, config.Keys.CycleViewNext) {
			cmd := t.onCycleView()
			return t, cmd
		}
		wsMsg := func() tea.Msg {
			return tea.WindowSizeMsg{Width: t.Width, Height: t.Height}
		}
		if key.Matches(msg, config.Keys.Increase) {
			t.ViewIndex = (t.ViewIndex + 1) % 8
			return t, wsMsg
		}
		if key.Matches(msg, config.Keys.Decrease) {
			t.ViewIndex = (t.ViewIndex - 1) % 8
			if t.ViewIndex < 0 {
				t.ViewIndex = len(viewList) - 1
			}
			return t, wsMsg
		}
	}

	switch msg := msg.(type) {
	default:
		thing, cmd := t.thing.Update(msg)
		t.thing = thing

		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}
	return t, tea.Batch(cmds...)
}

var viewList = []func() component.ContentArea{
	createDemoTabbedContentArea,

	createDemoDockContentAreaFull,
	createDemoDockContentAreaFullRatios,
	createDemoDockContentAreaUnorderedModels,
	createDemoDockContentAreaNotFull,

	createDemoPriorityContentAreaTop,
	createDemoPriorityContentAreaRight,
	createDemoPriorityContentAreaBottom,
	createDemoPriorityContentAreaLeft,

	createDemoStackContentAreaHorizontal,
	createDemoStackContentAreaVertical,

	createDemo2x2GridContentArea,
	createDemo2x3GridContentArea,
	createDemo3x3GridContentArea,

	createDemoSimpleContentArea,
	createDemoEmptyContentArea,
}

func (t *DemoModel) onCycleView() tea.Cmd {
	t.ViewIndex = (t.ViewIndex + 1) % len(viewList)
	t.thing = viewList[t.ViewIndex]()

	return func() tea.Msg {
		return tea.WindowSizeMsg{Width: t.Width, Height: t.Height}
	}
}

// View implements tea.Model.
func (t *DemoModel) View() string {
	return t.thing.View()
}

// CreateDemoModel is used by the demo application
func CreateDemoModel() tea.Model {
	// config.LoadConfig()
	return &DemoModel{
		thing: viewList[0](),
	}
}

// func createDemoLayout() widget.ContentArea {
// 	contentList := createDemoList(4)

// 	contentArea := widget.New(
// 		"Area",
// 		layout.NewDockLayout(
// 			layout.Panels(
// 				layout.Panel(layout.Top),
// 				layout.Panel(layout.Bottom),
// 			),
// 			layout.WithPanelBorder(lipgloss.DoubleBorder()),
// 			layout.WithStyle(
// 				lipgloss.NewStyle().
// 					BorderForeground(lipgloss.Color("205")).
// 					Foreground(lipgloss.Color("205")),
// 			),
// 		),
// 		contentList...,
// 	)
// 	return contentArea
// }

func createDemoDockContentAreaFull() component.ContentArea {
	top := content.Text("Top")
	bottom := content.Text("Bottom")
	left := content.Text("Left")
	right := content.Text("Right")
	center := content.Text("Center")

	dockLayout := layout.NewDockLayout(
		layout.Panels(
			layout.Panel(layout.Top, layout.WithHeight(7)),
			layout.Panel(layout.Bottom, layout.WithHeight(10)),
			layout.Panel(layout.Left, layout.WithWidth(67)),
			layout.Panel(layout.Right, layout.WithWidth(40)),
			layout.Panel(layout.Center),
		),
		layout.WithStyle(lipgloss.NewStyle().Border(lipgloss.RoundedBorder())),
	)

	contentArea := component.New("DockLayoutFull", dockLayout, top, bottom, left, right, center)

	contentArea.SetStyle(lipgloss.NewStyle().Border(lipgloss.NormalBorder()))

	return contentArea
}

func createDemoDockContentAreaFullRatios() component.ContentArea {
	top := content.Text("Top")
	bottom := content.Text("Bottom")
	left := content.Text("Left")
	right := content.Text("Right")
	center := content.Text("Center")

	dockLayout := layout.NewDockLayout(
		layout.Panels(
			layout.Panel(layout.Top, layout.WithHeightRatio(0.10)),
			layout.Panel(layout.Bottom, layout.WithHeightRatio(0.15)),
			layout.Panel(layout.Left, layout.WithWidthRatio(0.25)),
			layout.Panel(layout.Right, layout.WithWidthRatio(0.15)),
			layout.Panel(layout.Center),
		),
		layout.WithStyle(lipgloss.NewStyle().Border(lipgloss.RoundedBorder())),
	)

	contentArea := component.New("DockLayoutFull", dockLayout, top, bottom, left, right, center)
	contentArea.SetStyle(lipgloss.NewStyle().Border(lipgloss.NormalBorder()))

	return contentArea
}

func createDemoDockContentAreaNotFull() component.ContentArea {
	bottom := content.Text("Bottom")
	left := content.Text("Left")
	center := content.Text("Center")
	// contentList := createDemoList(3)
	dockLayout := layout.NewDockLayout(
		layout.Panels(
			layout.Panel(layout.Bottom, layout.WithHeight(25)),
			layout.Panel(layout.Left, layout.WithWidth(45)),
			layout.Panel(layout.Center),
		),
		layout.WithPanelBorder(lipgloss.DoubleBorder()),
		layout.WithDefaultUnfocusedStyle(),
	)

	contentArea := component.New("DockLayoutNotFull", dockLayout, bottom, left, center)

	contentArea.SetStyle(lipgloss.NewStyle().Border(lipgloss.ThickBorder()))

	return contentArea
}

func createDemoDockContentAreaUnorderedModels() component.ContentArea {
	left := content.Text("Left")
	top := content.Text("Top")
	bottom := content.Text("Bottom")
	center := content.Text("Center")
	right := content.Text("Right")
	// contentList := createDemoList(3)
	dockLayout := layout.NewDockLayout(
		layout.Panels(
			layout.Panel(layout.Left, layout.WithWidth(10)),
			layout.Panel(layout.Top, layout.WithHeight(16)),
			layout.Panel(layout.Center),
			layout.Panel(layout.Bottom, layout.WithHeight(42)),
			layout.Panel(layout.Right, layout.WithWidth(16)),
		),
		layout.WithPanelBorder(lipgloss.DoubleBorder()),
		layout.WithStyle(lipgloss.NewStyle().BorderForeground(lipgloss.Color("205")).Foreground(lipgloss.Color("205"))),
		layout.WithFocusedStyle(lipgloss.NewStyle().BorderForeground(lipgloss.Color("179")).Foreground(lipgloss.Color("179"))),
	)

	contentArea := component.New("DockLayoutNotFull", dockLayout, left, top, center, bottom, right)

	contentArea.SetStyle(lipgloss.NewStyle().Border(lipgloss.BlockBorder()))

	return contentArea
}

func createDemoPriorityContentAreaTop() component.ContentArea {
	contentList := createDemoList(4)

	contentArea := component.New("PrioritySplitTop",
		layout.NewPrioritySplitLayout(
			layout.Top,
			layout.WithPanelBorder(lipgloss.DoubleBorder()),
		),
		contentList...,
	)
	contentArea.SetStyle(lipgloss.NewStyle().Border(lipgloss.NormalBorder()))

	return contentArea
}

func createDemoPriorityContentAreaRight() component.ContentArea {
	contentList := createDemoList(4)

	contentArea := component.New("PrioritySplitRight", layout.NewPrioritySplitLayout(layout.Right), contentList...)

	contentArea.SetStyle(lipgloss.NewStyle().Border(lipgloss.RoundedBorder()))

	return contentArea
}

func createDemoPriorityContentAreaBottom() component.ContentArea {
	contentList := createDemoList(4)

	contentArea := component.New("PrioritySplitBottom", layout.NewPrioritySplitLayout(layout.Bottom), contentList...)
	contentArea.SetStyle(lipgloss.NewStyle().Border(lipgloss.NormalBorder()))

	return contentArea
}

func createDemoPriorityContentAreaLeft() component.ContentArea {
	contentList := createDemoList(4)

	contentArea := component.New("PrioritySplitLeft", layout.NewPrioritySplitLayout(layout.Left), contentList...)

	contentArea.SetStyle(lipgloss.NewStyle().Border(lipgloss.RoundedBorder()))

	return contentArea
}

func createDemoStackContentAreaHorizontal() component.ContentArea {
	contentList := createDemoList(3)

	contentArea := component.New("StackHorizontal", layout.NewStackLayout(layout.East), contentList...)
	contentArea.SetStyle(lipgloss.NewStyle().Border(lipgloss.NormalBorder()))

	return contentArea
}

func createDemoStackContentAreaVertical() component.ContentArea {
	contentList := createDemoList(3)

	contentArea := component.New("StackVertical", layout.NewStackLayout(layout.South), contentList...)
	contentArea.SetStyle(lipgloss.NewStyle().Border(lipgloss.NormalBorder()))

	return contentArea
}

func createDemoTabbedContentArea() component.ContentArea {
	contentList := createDemoList(4)

	contentArea := component.New("Tabbed", layout.NewTabbedLayout(layout.North), contentList...)
	contentArea.SetStyle(lipgloss.NewStyle().Border(lipgloss.NormalBorder()))

	return contentArea
}

func createDemo2x2GridContentArea() component.ContentArea {
	contentList := createDemoList(4)

	contentArea := component.New("2x2Grid",
		layout.NewUniformGridLayout(
			2,
			layout.WithPanelBorder(lipgloss.ASCIIBorder()),
		),
		contentList...,
	)
	contentArea.SetStyle(lipgloss.NewStyle().Border(lipgloss.NormalBorder()))

	return contentArea
}

func createDemo2x3GridContentArea() component.ContentArea {
	contentList := createDemoList(6)

	contentArea := component.New("2x3Grid",
		layout.NewUniformGridLayout(
			3,
			layout.WithPanelBorder(lipgloss.ASCIIBorder()),
			layout.WithStyle(lipgloss.NewStyle().BorderForeground(lipgloss.Color("205"))),
		),
		contentList...,
	)
	contentArea.SetStyle(lipgloss.NewStyle().Border(lipgloss.NormalBorder()))

	return contentArea
}

func createDemo3x3GridContentArea() component.ContentArea {
	contentList := createDemoList(9)

	contentArea := component.New("3x3Grid",
		layout.NewUniformGridLayout(
			3,
			layout.WithPanelBorder(lipgloss.ASCIIBorder()),
			layout.WithStyle(lipgloss.NewStyle().BorderForeground(lipgloss.Color("205"))),
		),
		contentList...,
	)
	contentArea.SetStyle(lipgloss.NewStyle().Border(lipgloss.NormalBorder()))

	return contentArea
}

func createDemoSimpleContentArea() component.ContentArea {
	contentArea := component.New("Simple", layout.NewSimpleLayout(), content.Text("Simple Layout"))
	contentArea.SetStyle(lipgloss.NewStyle().Border(lipgloss.NormalBorder()))

	return contentArea
}

func createDemoEmptyContentArea() component.ContentArea {
	contentArea := component.New("empty", layout.NewSimpleLayout())
	contentArea.SetStyle(lipgloss.NewStyle().Border(lipgloss.NormalBorder()))
	return contentArea
}

func createDemoList(numItems int) []tea.Model {
	list := make([]tea.Model, 0, numItems)

	for i := range numItems {
		list = append(list, content.NewList([]content.Item{
			{Value: fmt.Sprintf("foo%d", i)},
			{Value: fmt.Sprintf("bar%d", i)},
			{Value: fmt.Sprintf("baz%d", i)},
		}))
	}
	return list
}

// Exampling
/*
if key.Matches(msg, config.Keys.Up) {
			t.posY = t.posY + 0.1
			if t.posY < 0 {
				t.posY = 0
			}
			logger.Infof("posX: %f, posY: %f\n", t.posX, t.posY)
		} else if key.Matches(msg, config.Keys.Down) {
			t.posY = t.posY - 0.1
			if t.posY > 1 {
				t.posY = 1
			}
			logger.Infof("posX: %f, posY: %f\n", t.posX, t.posY)
		} else if key.Matches(msg, config.Keys.Left) {
			t.posX = t.posX + 0.1
			if t.posX < 0 {
				t.posX = 0
			}
			logger.Infof("posX: %f, posY: %f\n", t.posX, t.posY)
		} else if key.Matches(msg, config.Keys.Right) {
			t.posX = t.posX - 0.1
			if t.posX > 1 {
				t.posX = 1
			}
			logger.Infof("posX: %f, posY: %f\n", t.posX, t.posY)
		}
func (t *DemoModel) View() string {

	var view string
	viewStyle := lipgloss.NewStyle().Border(lipgloss.NormalBorder())
	switch t.ViewIndex {
	case 0: // center, center
		view = lipgloss.Place(
			25,
			20,
			lipgloss.Center,
			lipgloss.Center,
			viewStyle.Render("center center"))
	case 1:
		view = lipgloss.Place(
			25,
			20,
			lipgloss.Left,
			lipgloss.Center,
			viewStyle.Render("left center"))
	case 2:
		view = lipgloss.Place(
			25,
			20,
			lipgloss.Right,
			lipgloss.Center,
			viewStyle.Render("right center"))
	case 3:
		view = lipgloss.Place(
			25,
			20,
			lipgloss.Center,
			lipgloss.Top,
			viewStyle.Render("center top"))
	case 4:
		view = lipgloss.Place(
			25,
			20,
			lipgloss.Center,
			lipgloss.Bottom,
			viewStyle.Render("center bottom"))
	case 5:
		view = lipgloss.Place(
			25,
			20,
			0.25,
			lipgloss.Center,
			viewStyle.Render("0.25 center"))
	case 6:
		view = lipgloss.Place(
			25,
			20,
			lipgloss.Center,
			0.25,
			viewStyle.Render("center 0.25"))
	case 7:
		view = lipgloss.Place(
			25,
			20,
			0.25,
			0.25,
			viewStyle.Render("0.25 0.25"))

	default:
		view = "default"
	}

	return lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		Render(lipgloss.Place(
			t.InnerWidth,
			t.InnerHeight,
			lipgloss.Position(t.posX),
			lipgloss.Position(t.posY),
			view))
	// return t.thing.View()
}
*/
