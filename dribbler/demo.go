package dribbler

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ctrl-alt-boop/dribble"
	"github.com/ctrl-alt-boop/dribbler/config"
	"github.com/ctrl-alt-boop/dribbler/ui/content"
	"github.com/ctrl-alt-boop/dribbler/ui/layout"
	"github.com/ctrl-alt-boop/dribbler/widget"
)

var _ tea.Model = (*DemoModel)(nil)

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
		borderStyle := lipgloss.NewStyle().Border(lipgloss.NormalBorder())
		t.InnerWidth, t.InnerHeight = msg.Width-borderStyle.GetHorizontalFrameSize(), msg.Height-borderStyle.GetVerticalFrameSize()
		updatedThing, cmd := t.thing.Update(tea.WindowSizeMsg{Width: t.InnerWidth, Height: t.InnerHeight})
		t.thing = updatedThing
		if cmd != nil {
			cmds = append(cmds, cmd)

		}
		return t, tea.Batch(cmds...)
	case tea.KeyMsg:
		if key.Matches(msg, config.Keys.Quit) {
			return t, tea.Quit
		}
		if key.Matches(msg, config.Keys.CycleView) {
			cmd := t.OnCycleView()
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
				t.ViewIndex = len(ViewList) - 1
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

var ViewList = []func() widget.ContentArea{
	CreateDemoPriorityContentAreaTop,
	CreateDemoPriorityContentAreaRight,
	CreateDemoPriorityContentAreaBottom,
	CreateDemoPriorityContentAreaLeft,

	CreateDemoStackContentAreaHorizontal,
	CreateDemoStackContentAreaVertical,

	CreateDemo2x2GridContentArea,
	CreateDemo2x3GridContentArea,
	CreateDemo3x3GridContentArea,

	CreateDemoDockContentAreaFull,
	CreateDemoDockContentAreaFullRatios,
	CreateDemoDockContentAreaUnorderedModels,
	CreateDemoDockContentAreaNotFull,

	CreateDemoTabbedContentArea,

	CreateDemoSimpleContentArea,
	CreateDemoEmptyContentArea,
}

func (t *DemoModel) OnCycleView() tea.Cmd {
	t.ViewIndex = (t.ViewIndex + 1) % len(ViewList)
	t.thing = ViewList[t.ViewIndex]()

	return func() tea.Msg {
		return tea.WindowSizeMsg{Width: t.Width, Height: t.Height}
	}
}

// Unused
func (t *DemoModel) StyledView() string {
	style := lipgloss.NewStyle().Align(lipgloss.Center, lipgloss.Center).Border(lipgloss.NormalBorder())
	return style.
		Width(t.InnerWidth).
		Height(t.InnerHeight).
		Render(t.thing.View())
}

func (t *DemoModel) View() string {
	return t.thing.View()
}

func CreateDemoModel(dribbleClient *dribble.Client) tea.Model {
	// config.LoadConfig()
	return &DemoModel{
		thing: ViewList[0](),
	}
}

func CreateDemoLayout() widget.ContentArea {
	contentList := CreateDemoList(4)

	contentArea := widget.New(
		"Area",
		layout.NewDockLayout(
			layout.Panels(
				layout.Panel(layout.Top),
				layout.Panel(layout.Bottom),
			),
			layout.WithPanelBorder(lipgloss.DoubleBorder()),
			layout.WithStyle(
				lipgloss.NewStyle().
					BorderForeground(lipgloss.Color("205")).
					Foreground(lipgloss.Color("205")),
			),
		),
		contentList...,
	)
	return contentArea
}

func CreateDemoDockContentAreaFull() widget.ContentArea {
	top := content.Text{Item: content.Item{Value: "Top"}}
	bottom := content.Text{Item: content.Item{Value: "Bottom"}}
	left := content.Text{Item: content.Item{Value: "Left"}}
	right := content.Text{Item: content.Item{Value: "Right"}}
	center := content.Text{Item: content.Item{Value: "Center"}}

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

	contentArea := widget.New("DockLayoutFull", dockLayout, top, bottom, left, right, center)

	contentArea.SetStyle(lipgloss.NewStyle().Border(lipgloss.NormalBorder()))

	return contentArea
}

func CreateDemoDockContentAreaFullRatios() widget.ContentArea {
	top := content.Text{Item: content.Item{Value: "Top"}}
	bottom := content.Text{Item: content.Item{Value: "Bottom"}}
	left := content.Text{Item: content.Item{Value: "Left"}}
	right := content.Text{Item: content.Item{Value: "Right"}}
	center := content.Text{Item: content.Item{Value: "Center"}}

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

	contentArea := widget.New("DockLayoutFull", dockLayout, top, bottom, left, right, center)
	contentArea.SetStyle(lipgloss.NewStyle().Border(lipgloss.NormalBorder()))

	return contentArea
}

func CreateDemoDockContentAreaNotFull() widget.ContentArea {
	bottom := content.Text{Item: content.Item{Value: "Bottom"}}
	left := content.Text{Item: content.Item{Value: "Left"}}
	center := content.Text{Item: content.Item{Value: "Center"}}
	// contentList := CreateDemoList(3)
	dockLayout := layout.NewDockLayout(
		layout.Panels(
			layout.Panel(layout.Bottom, layout.WithHeight(25)),
			layout.Panel(layout.Left, layout.WithWidth(45)),
			layout.Panel(layout.Center),
		),
		layout.WithPanelBorder(lipgloss.DoubleBorder()),
		layout.WithDefaultUnfocusedStyle(),
	)

	contentArea := widget.New("DockLayoutNotFull", dockLayout, bottom, left, center)

	contentArea.SetStyle(lipgloss.NewStyle().Border(lipgloss.ThickBorder()))

	return contentArea
}

func CreateDemoDockContentAreaUnorderedModels() widget.ContentArea {
	left := content.Text{Item: "Left"}
	top := content.Text{Item: "Top"}
	bottom := content.Text{Item: "Bottom"}
	center := content.Text{Item: "Center"}
	right := content.Text{Item: "Right"}
	// contentList := CreateDemoList(3)
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

	contentArea := widget.New("DockLayoutNotFull", dockLayout, left, top, center, bottom, right)

	contentArea.SetStyle(lipgloss.NewStyle().Border(lipgloss.BlockBorder()))

	return contentArea
}

func CreateDemoPriorityContentAreaTop() widget.ContentArea {
	contentList := CreateDemoList(4)

	contentArea := widget.New("PrioritySplitTop",
		layout.NewPrioritySplitLayout(
			layout.Top,
			layout.WithPanelBorder(lipgloss.DoubleBorder()),
		),
		contentList...,
	)
	contentArea.SetStyle(lipgloss.NewStyle().Border(lipgloss.NormalBorder()))

	return contentArea
}

func CreateDemoPriorityContentAreaRight() widget.ContentArea {
	contentList := CreateDemoList(4)

	contentArea := widget.New("PrioritySplitRight", layout.NewPrioritySplitLayout(layout.Right), contentList...)

	contentArea.SetStyle(lipgloss.NewStyle().Border(lipgloss.RoundedBorder()))

	return contentArea
}

func CreateDemoPriorityContentAreaBottom() widget.ContentArea {
	contentList := CreateDemoList(4)

	contentArea := widget.New("PrioritySplitBottom", layout.NewPrioritySplitLayout(layout.Bottom), contentList...)
	contentArea.SetStyle(lipgloss.NewStyle().Border(lipgloss.NormalBorder()))

	return contentArea
}

func CreateDemoPriorityContentAreaLeft() widget.ContentArea {
	contentList := CreateDemoList(4)

	contentArea := widget.New("PrioritySplitLeft", layout.NewPrioritySplitLayout(layout.Left), contentList...)

	contentArea.SetStyle(lipgloss.NewStyle().Border(lipgloss.RoundedBorder()))

	return contentArea
}

func CreateDemoStackContentAreaHorizontal() widget.ContentArea {
	contentList := CreateDemoList(3)

	contentArea := widget.New("StackHorizontal", layout.NewStackLayout(layout.East), contentList...)
	contentArea.SetStyle(lipgloss.NewStyle().Border(lipgloss.NormalBorder()))

	return contentArea
}

func CreateDemoStackContentAreaVertical() widget.ContentArea {
	contentList := CreateDemoList(3)

	contentArea := widget.New("StackVertical", layout.NewStackLayout(layout.South), contentList...)
	contentArea.SetStyle(lipgloss.NewStyle().Border(lipgloss.NormalBorder()))

	return contentArea
}

func CreateDemoTabbedContentArea() widget.ContentArea {
	contentList := CreateDemoList(4)

	contentArea := widget.New("Tabbed", layout.NewTabbedLayout(layout.North), contentList...)
	contentArea.SetStyle(lipgloss.NewStyle().Border(lipgloss.NormalBorder()))

	return contentArea
}

func CreateDemo2x2GridContentArea() widget.ContentArea {
	contentList := CreateDemoList(4)

	contentArea := widget.New("2x2Grid",
		layout.NewUniformGridLayout(
			2,
			layout.WithPanelBorder(lipgloss.ASCIIBorder()),
		),
		contentList...,
	)
	contentArea.SetStyle(lipgloss.NewStyle().Border(lipgloss.NormalBorder()))

	return contentArea
}

func CreateDemo2x3GridContentArea() widget.ContentArea {
	contentList := CreateDemoList(6)

	contentArea := widget.New("2x3Grid",
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

func CreateDemo3x3GridContentArea() widget.ContentArea {
	contentList := CreateDemoList(9)

	contentArea := widget.New("3x3Grid",
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

func CreateDemoSimpleContentArea() widget.ContentArea {

	contentArea := widget.New("Simple", layout.NewSimpleLayout(), content.Text{Item: "Simple Layout"})
	contentArea.SetStyle(lipgloss.NewStyle().Border(lipgloss.NormalBorder()))

	return contentArea
}

func CreateDemoEmptyContentArea() widget.ContentArea {

	contentArea := widget.New("empty", layout.NewSimpleLayout())
	contentArea.SetStyle(lipgloss.NewStyle().Border(lipgloss.NormalBorder()))
	return contentArea
}

func CreateDemoList(numItems int) []tea.Model {
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
