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

var _ tea.Model = (*TestModel)(nil)

type TestModel struct {
	thing                   tea.Model
	Width, Height           int
	InnerWidth, InnerHeight int
	ViewIndex               int
}

// Init implements tea.Model.
func (t *TestModel) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model.
func (t *TestModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
	CreateTestDockContentAreaNotFull,
	CreateTestDockContentAreaFull,
	CreateTestPriorityContentAreaHorizontal,
	CreateTestPriorityContentAreaVertical,
	CreateTestStackContentAreaHorizontal,
	CreateTestStackContentAreaVertical,
	CreateTestTabbedContentArea,
	CreateTest2x2GridContentArea,
	CreateTest2x3GridContentArea,
	CreateTest3x3GridContentArea,
	CreateTestSimpleContentArea,
	CreateTestEmptyContentArea,
}

func (t *TestModel) OnCycleView() tea.Cmd {
	t.ViewIndex = (t.ViewIndex + 1) % len(ViewList)
	t.thing = ViewList[t.ViewIndex]()

	return func() tea.Msg {
		return tea.WindowSizeMsg{Width: t.Width, Height: t.Height}
	}
}

func (t *TestModel) StyledView() string {
	style := lipgloss.NewStyle().Align(lipgloss.Center, lipgloss.Center).Border(lipgloss.NormalBorder())
	return style.
		Width(t.InnerWidth).
		Height(t.InnerHeight).
		Render(t.thing.View())
}

func (t *TestModel) View() string {
	return t.thing.View()
}

func CreateTestModel(dribbleClient *dribble.Client) tea.Model {
	// config.LoadConfig()
	return &TestModel{
		thing: ViewList[0](),
	}
}

func CreateTestDockContentAreaFull() widget.ContentArea {
	top := content.Text{Item: content.Item{Value: "Top"}}
	bottom := content.Text{Item: content.Item{Value: "Bottom"}}
	left := content.Text{Item: content.Item{Value: "Left"}}
	right := content.Text{Item: content.Item{Value: "Right"}}
	center := content.Text{Item: content.Item{Value: "Center"}}

	contentArea := widget.NewContentArea(0, "DockLayoutFull", top, bottom, left, right, center)

	contentArea.SetStyle(lipgloss.NewStyle().Border(lipgloss.NormalBorder()))
	dockLayout := layout.NewDockLayout()
	dockLayout.SetDefinition(
		layout.NewDefinition(
			[]layout.LayoutDefinition{
				layout.NewLayoutDefinition(layout.Top, layout.WithMaxHeight(5)),
				layout.NewLayoutDefinition(layout.Bottom, layout.WithMaxHeight(15)),
				layout.NewLayoutDefinition(layout.Left, layout.WithMaxWidth(45)),
				layout.NewLayoutDefinition(layout.Right, layout.WithMaxWidth(25)),
				layout.NewLayoutDefinition(layout.Center),
			},
			layout.WithStyle(lipgloss.NewStyle().Border(lipgloss.RoundedBorder())),
		))
	contentArea.SetLayoutManager(dockLayout)

	return contentArea
}

func CreateTestDockContentAreaNotFull() widget.ContentArea {
	bottom := content.Text{Item: content.Item{Value: "Bottom"}}
	left := content.Text{Item: content.Item{Value: "Left"}}
	center := content.Text{Item: content.Item{Value: "Center"}}
	// contentList := CreateTestList(3)

	contentArea := widget.NewContentArea(0, "DockLayoutNotFull", bottom, left, center)

	contentArea.SetStyle(lipgloss.NewStyle().Border(lipgloss.NormalBorder()))
	dockLayout := layout.NewDockLayout()
	dockLayout.SetDefinition(
		layout.NewDefinition(
			[]layout.LayoutDefinition{
				layout.NewLayoutDefinition(layout.Bottom, layout.WithMaxHeight(25)),
				layout.NewLayoutDefinition(layout.Left, layout.WithMaxWidth(45)),
				layout.NewLayoutDefinition(layout.Center),
			},
			layout.WithStyle(lipgloss.NewStyle().Border(lipgloss.RoundedBorder())),
		))
	contentArea.SetLayoutManager(dockLayout)

	return contentArea
}

func CreateTestPriorityContentAreaHorizontal() widget.ContentArea {
	contentList := CreateTestList(4)

	contentArea := widget.NewContentArea(0, "PrioritySplitHorizontal", contentList...)

	contentArea.SetStyle(lipgloss.NewStyle().Border(lipgloss.RoundedBorder()))
	contentArea.SetLayoutManager(layout.NewPrioritySplitLayout(layout.Horizontal))
	return contentArea
}

func CreateTestPriorityContentAreaVertical() widget.ContentArea {
	contentList := CreateTestList(4)

	contentArea := widget.NewContentArea(0, "PrioritySplitVertical", contentList...)
	contentArea.SetStyle(lipgloss.NewStyle().Border(lipgloss.NormalBorder()))
	contentArea.SetLayoutManager(layout.NewPrioritySplitLayout(layout.Vertical))
	return contentArea
}

func CreateTestStackContentAreaHorizontal() widget.ContentArea {
	contentList := CreateTestList(3)

	contentArea := widget.NewContentArea(0, "StackHorizontal", contentList...)
	contentArea.SetStyle(lipgloss.NewStyle().Border(lipgloss.NormalBorder()))
	contentArea.SetLayoutManager(layout.NewStackLayout(layout.Horizontal))
	return contentArea
}

func CreateTestStackContentAreaVertical() widget.ContentArea {
	contentList := CreateTestList(3)

	contentArea := widget.NewContentArea(0, "StackVertical", contentList...)
	contentArea.SetStyle(lipgloss.NewStyle().Border(lipgloss.NormalBorder()))
	contentArea.SetLayoutManager(layout.NewStackLayout(layout.Vertical))
	return contentArea
}

func CreateTestTabbedContentArea() widget.ContentArea {
	contentList := CreateTestList(4)

	contentArea := widget.NewContentArea(0, "Tabbed", contentList...)
	contentArea.SetStyle(lipgloss.NewStyle().Border(lipgloss.NormalBorder()))
	contentArea.SetLayoutManager(layout.NewTabbedLayout(layout.Top))
	return contentArea
}

func CreateTest2x2GridContentArea() widget.ContentArea {
	contentList := CreateTestList(4)

	contentArea := widget.NewContentArea(0, "2x2Grid", contentList...)
	contentArea.SetStyle(lipgloss.NewStyle().Border(lipgloss.NormalBorder()))
	contentArea.SetLayoutManager(layout.NewUniformGridLayout(2))

	return contentArea
}

func CreateTest2x3GridContentArea() widget.ContentArea {
	contentList := CreateTestList(6)

	contentArea := widget.NewContentArea(0, "2x3Grid", contentList...)
	contentArea.SetStyle(lipgloss.NewStyle().Border(lipgloss.NormalBorder()))
	contentArea.SetLayoutManager(layout.NewUniformGridLayout(3))

	return contentArea
}

func CreateTest3x3GridContentArea() widget.ContentArea {
	contentList := CreateTestList(9)

	contentArea := widget.NewContentArea(0, "3x3Grid", contentList...)
	contentArea.SetStyle(lipgloss.NewStyle().Border(lipgloss.NormalBorder()))
	contentArea.SetLayoutManager(layout.NewUniformGridLayout(3))

	return contentArea
}

func CreateTestSimpleContentArea() widget.ContentArea {

	contentArea := widget.NewContentArea(0, "Simple", content.Text{Item: content.Item{Value: "Simple Layout"}})
	contentArea.SetStyle(lipgloss.NewStyle().Border(lipgloss.NormalBorder()))
	contentArea.SetLayoutManager(layout.NewSimpleLayout())

	return contentArea
}

func CreateTestEmptyContentArea() widget.ContentArea {

	contentArea := widget.NewContentArea(0, "empty")
	contentArea.SetStyle(lipgloss.NewStyle().Border(lipgloss.NormalBorder()))
	return contentArea
}

func CreateTestList(numItems int) []tea.Model {
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
func (t *TestModel) View() string {

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
