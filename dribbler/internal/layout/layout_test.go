package layout_test

// import (
// 	"fmt"
// 	"strings"
// 	"testing"

// 	tea "charm.land/bubbletea/v2"
// 	"github.com/ctrl-alt-boop/dribbler/layout"
// 	"github.com/stretchr/testify/assert"
// )

// // mockModel is a simple tea.Model for testing purposes.
// type mockModel struct {
// 	id            int
// 	width, height int
// 	viewContent   string
// 	updateCount   int
// 	lastMsg       tea.Msg
// }

// func (m mockModel) Init() tea.Cmd { return nil }

// func (m mockModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
// 	m.updateCount++
// 	m.lastMsg = msg
// 	if size, ok := msg.(tea.WindowSizeMsg); ok {
// 		m.width = size.Width
// 		m.height = size.Height
// 	}
// 	return m, nil
// }

// func (m mockModel) Name() string {
// 	return fmt.Sprintf("mock-name-%d", m.id)
// }

// func (m mockModel) View() string {
// 	if m.viewContent != "" {
// 		return m.viewContent
// 	}
// 	return fmt.Sprintf("mock-view-%d", m.id)
// }

// func TestSimpleLayout(t *testing.T) {
// 	l := layout.SimpleLayout{}
// 	child1 := mockModel{id: 1, viewContent: "child1"}
// 	child2 := mockModel{id: 2, viewContent: "child2"}
// 	models := []tea.Model{child1, child2}

// 	t.Run("Layout", func(t *testing.T) {
// 		l.SetSize(100, 50)
// 		models = l.Layout(models)
// 		c1 := models[0].(mockModel)
// 		c2 := models[1].(mockModel)

// 		assert.Equal(t, 100, c1.width, "First child should get full width")
// 		assert.Equal(t, 50, c1.height, "First child should get full height")
// 		assert.Equal(t, 0, c2.width, "Second child should not be updated")
// 	})

// 	t.Run("View", func(t *testing.T) {
// 		view := l.View([]tea.Model{child1, child2})
// 		assert.Equal(t, "child1", view, "Should only render the first child's view")
// 	})
// }

// func TestUniformGridLayout(t *testing.T) {
// 	l := layout.NewUniformGridLayout(2)
// 	child1 := mockModel{id: 1, viewContent: "View1"}
// 	child2 := mockModel{id: 2, viewContent: "View2"}
// 	child3 := mockModel{id: 3, viewContent: "View3"}
// 	models := []tea.Model{child1, child2, child3}

// 	t.Run("Layout", func(t *testing.T) {
// 		// 101 width - 1 gutter = 100. 100 / 2 cols = 50 width per cell.
// 		// 3 children / 2 cols = 1.5 -> 2 rows. 100 height / 2 rows = 50 height per cell.
// 		l.SetSize(101, 100)
// 		models = l.Layout(models)
// 		c1 := models[0].(mockModel)
// 		c2 := models[1].(mockModel)
// 		c3 := models[2].(mockModel)

// 		assert.Equal(t, 50, c1.width, "Child 1 width incorrect")
// 		assert.Equal(t, 49, c1.height, "Child 1 height incorrect")
// 		assert.Equal(t, 50, c2.width, "Child 2 width incorrect")
// 		assert.Equal(t, 49, c2.height, "Child 2 height incorrect")
// 		assert.Equal(t, 50, c3.width, "Child 3 width incorrect")
// 		assert.Equal(t, 49, c3.height, "Child 3 height incorrect")
// 	})

// 	t.Run("View", func(t *testing.T) {
// 		view := l.View(models)
// 		// lipgloss joins add spaces, so we check for substrings and order
// 		assert.Contains(t, view, "View1", "Should contain View1")
// 		assert.Contains(t, view, "View2", "Should contain View2")
// 		assert.Contains(t, view, "View3", "Should contain View3")
// 		assert.True(t, strings.Index(view, "View3") > strings.Index(view, "View1"), "View3 should appear after View1")
// 	})
// }

// func TestPrioritySplitLayout(t *testing.T) {
// 	child1 := mockModel{id: 1, viewContent: "Primary"}
// 	child2 := mockModel{id: 2, viewContent: "Secondary1"}
// 	child3 := mockModel{id: 3, viewContent: "Secondary2"}
// 	models := []tea.Model{child1, child2, child3}

// 	t.Run("Layout Horizontal", func(t *testing.T) {
// 		l := layout.NewPrioritySplitLayout(layout.Right)
// 		l.PrimarySizeRatio = 0.7
// 		// 101 width - 1 gutter = 100.
// 		// Primary: 100 * 0.7 = 70 width.
// 		// Secondary: 100 - 70 = 30 width total. 30 / 2 children = 15 width each.
// 		l.SetSize(101, 80)
// 		models = l.Layout(models)
// 		c1 := models[0].(mockModel)
// 		c2 := models[1].(mockModel)
// 		c3 := models[2].(mockModel)

// 		assert.Equal(t, 70, c1.width, "Primary child width incorrect")
// 		assert.Equal(t, 80, c1.height, "Primary child height incorrect")
// 		assert.Equal(t, 15, c2.width, "Secondary child 1 width incorrect")
// 		assert.Equal(t, 80, c2.height, "Secondary child 1 height incorrect")
// 		assert.Equal(t, 15, c3.width, "Secondary child 2 width incorrect")
// 		assert.Equal(t, 80, c3.height, "Secondary child 2 height incorrect")
// 	})

// 	t.Run("Layout Vertical", func(t *testing.T) {
// 		l := layout.NewPrioritySplitLayout(layout.Top)
// 		l.PrimarySizeRatio = 0.6

// 		// 101 height - 1 gutter = 100.
// 		// Primary: 100 * 0.6 = 60 height.
// 		// Secondary: 100 - 60 = 40 height total. 40 / 2 children = 20 height each.
// 		l.SetSize(90, 101)
// 		models = l.Layout(models)
// 		c1 := models[0].(mockModel)
// 		c2 := models[1].(mockModel)
// 		c3 := models[2].(mockModel)

// 		assert.Equal(t, 90, c1.width, "Primary child width incorrect")
// 		assert.Equal(t, 60, c1.height, "Primary child height incorrect")
// 		assert.Equal(t, 90, c2.width, "Secondary child 1 width incorrect")
// 		assert.Equal(t, 20, c2.height, "Secondary child 1 height incorrect")
// 		assert.Equal(t, 90, c3.width, "Secondary child 2 width incorrect")
// 		assert.Equal(t, 20, c3.height, "Secondary child 2 height incorrect")
// 	})

// 	t.Run("View", func(t *testing.T) {
// 		l := layout.NewPrioritySplitLayout(layout.Left)
// 		view := l.View(models)
// 		assert.Contains(t, view, "Primary", "View should contain primary view")
// 		assert.Contains(t, view, "Secondary1", "View should contain secondary view 1")
// 		assert.Contains(t, view, "Secondary2", "View should contain secondary view 2")
// 		assert.True(t, strings.Index(view, "Secondary1") > strings.Index(view, "Primary"), "Secondary should appear after primary")
// 	})
// }

// func TestStackLayout(t *testing.T) {
// 	child1 := mockModel{id: 1, viewContent: "View1"}
// 	child2 := mockModel{id: 2, viewContent: "View2"}
// 	models := []tea.Model{child1, child2}

// 	t.Run("Layout", func(t *testing.T) {
// 		l := layout.NewStackLayout(0) // 0 should not break anything
// 		l.SetSize(100, 50)
// 		models = l.Layout(models)
// 		c1 := models[0].(mockModel)
// 		c2 := models[1].(mockModel)

// 		// None should lead to no changes
// 		assert.Equal(t, 100, c1.width, "Child 1 should get full width")
// 		assert.Equal(t, 50, c1.height, "Child 1 should get full height")
// 		assert.Equal(t, 100, c2.width, "Child 2 should get full width")
// 		assert.Equal(t, 50, c2.height, "Child 2 should get full height")
// 	})

// 	t.Run("View Horizontal", func(t *testing.T) {
// 		l := &layout.StackLayout{StackDirection: layout.East}
// 		view := l.View(models)
// 		assert.True(t, strings.Index(view, "View2") > strings.Index(view, "View1"), "Views should be joined horizontally")
// 	})

// 	t.Run("View Vertical", func(t *testing.T) {
// 		l := &layout.StackLayout{StackDirection: layout.South}
// 		view := l.View(models)
// 		assert.True(t, strings.Index(view, "View2") > strings.Index(view, "View1"), "Views should be joined vertically")
// 		assert.Contains(t, view, "\n", "Vertical join should contain newlines")
// 	})
// }

// func TestTabbedLayout(t *testing.T) {
// 	child1 := mockModel{id: 1, viewContent: "Content1"}
// 	child2 := mockModel{id: 2, viewContent: "Content2"}
// 	models := []tea.Model{child1, child2}

// 	t.Run("Layout", func(t *testing.T) {
// 		l := &layout.TabbedLayout{ActiveIndex: 1}
// 		l.SetSize(100, 50)
// 		models = l.Layout(models)
// 		c1 := models[0].(mockModel)
// 		c2 := models[1].(mockModel)

// 		// Tab height is hardcoded to 1 + vertical frame size. Let's assume 2 for a simple border.
// 		// 50 - (1 + 2) = 47
// 		assert.Equal(t, 0, c1.width, "Inactive child should not be resized")
// 		assert.Equal(t, 100, c2.width, "Active child should get full width")
// 		assert.InDelta(t, 47, c2.height, 2, "Active child height should account for tabs")
// 	})

// 	t.Run("View", func(t *testing.T) {
// 		l := &layout.TabbedLayout{ActiveIndex: 0}
// 		view := l.View(models)

// 		assert.Contains(t, view, "mock-name-1", "Should contain tab for child 1")
// 		assert.Contains(t, view, "mock-name-2", "Should contain tab for child 2")
// 		assert.Contains(t, view, "Content1", "Should contain active child's content")
// 		assert.NotContains(t, view, "Content2", "Should not contain inactive child's content")
// 	})
// }
