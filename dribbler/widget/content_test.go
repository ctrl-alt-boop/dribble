package widget_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ctrl-alt-boop/dribbler/config"
	"github.com/ctrl-alt-boop/dribbler/widget"
	"github.com/ctrl-alt-boop/dribbler/widget/layout"
	"github.com/stretchr/testify/assert"
)

// mockModel is a simple tea.Model for testing purposes.
type mockModel struct {
	id            int
	width, height int
	viewContent   string
	updateCount   int
	lastMsg       tea.Msg
}

func (m *mockModel) Init() tea.Cmd { return nil }

func (m *mockModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.updateCount++
	m.lastMsg = msg
	if size, ok := msg.(tea.WindowSizeMsg); ok {
		m.width = size.Width
		m.height = size.Height
	}
	return m, nil
}

func (m *mockModel) Name() string {
	return fmt.Sprintf("mock-name-%d", m.id)
}

func (m *mockModel) View() string {
	if m.viewContent != "" {
		return m.viewContent
	}
	return fmt.Sprintf("mock-%d", m.id)
}

// mockLayout is a mock implementation of layout.Manager for testing ContentArea.
type mockLayout struct {
	layoutCalled           bool
	layoutCalledWithWidth  int
	layoutCalledWithHeight int
	layoutCmdToReturn      tea.Cmd

	viewCalled         bool
	viewStringToReturn string
}

func (m *mockLayout) Layout(width, height int, models []tea.Model) tea.Cmd {
	m.layoutCalled = true
	m.layoutCalledWithWidth = width
	m.layoutCalledWithHeight = height
	return m.layoutCmdToReturn
}

func (m *mockLayout) View(models []tea.Model) string {
	m.viewCalled = true
	// In a real scenario, you might check the models passed in.
	// For this test, we just return a predefined string.
	return m.viewStringToReturn
}

func TestNewContentArea(t *testing.T) {
	ca := widget.NewContentArea(1, "test-area")
	assert.NotNil(t, ca, "NewContentArea should not return nil")
	assert.NotNil(t, ca.Layout, "Layout should be initialized")
	assert.IsType(t, &layout.SimpleLayout{}, ca.Layout, "Default layout should be SimpleLayout")
}

func TestContentArea_AddChild(t *testing.T) {
	ca := widget.NewContentArea(1, "test")
	child1 := widget.NewContentArea(2, "child1")
	child2 := widget.NewContentArea(3, "child2")

	ca.AddChild(child1)
	assert.Len(t, ca.Children, 1, "Should add the first child")
	assert.Equal(t, child1, ca.Children[0])

	ca.AddChild(child2)
	assert.Len(t, ca.Children, 2, "Should add the second child")
	assert.Equal(t, child2, ca.Children[1])
}

func TestContentArea_Update(t *testing.T) {
	ca := widget.NewContentArea(1, "test")
	mockL := &mockLayout{}
	ca.Layout = mockL
	child1 := &mockModel{id: 1}
	child2 := &mockModel{id: 2}
	ca.AddChild(widget.NewContentArea(2, "child1"))
	ca.AddChild(widget.NewContentArea(3, "child2"))
	// Replace with mocks for testing update delegation
	ca.Children[0] = child1
	ca.Children[1] = child2

	t.Run("WindowSizeMsg", func(t *testing.T) {
		mockL.layoutCalled = false
		mockL.layoutCmdToReturn = func() tea.Msg { return "layout-cmd" }

		_, cmd := ca.Update(tea.WindowSizeMsg{Width: 100, Height: 50})

		assert.True(t, mockL.layoutCalled, "Layout.Layout should be called on WindowSizeMsg")
		assert.Equal(t, 100, mockL.layoutCalledWithWidth, "Layout should be called with correct width")
		assert.Equal(t, 50, mockL.layoutCalledWithHeight, "Layout should be called with correct height")
		assert.Equal(t, "layout-cmd", cmd(), "Update should return the command from the layout manager")
	})

	t.Run("CycleView KeyMsg", func(t *testing.T) {
		// Reset
		ca.FocusedChild = 0

		// Cycle once
		_, _ = ca.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(config.Keys.CycleView.Keys()[0])})
		assert.Equal(t, 1, ca.FocusedChild, "FocusedChild should cycle to the next child")

		// Cycle again to wrap around
		_, _ = ca.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(config.Keys.CycleView.Keys()[0])})
		assert.Equal(t, 0, ca.FocusedChild, "FocusedChild should wrap around to the first child")
	})

	t.Run("Back KeyMsg", func(t *testing.T) {
		ca.FocusedChild = 1
		_, cmd := ca.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(config.Keys.Back.Keys()[0])})

		msg := cmd()
		unfocusMsg, ok := msg.(widget.UnfocusChildMsg)
		assert.True(t, ok, "Should return an UnfocusChildMsg")
		assert.Equal(t, 1, unfocusMsg.ID, "Unfocus message should have the ID of the focused child")
	})

	t.Run("Delegate to focused child", func(t *testing.T) {
		// Reset
		child1.updateCount = 0
		child2.updateCount = 0
		ca.FocusedChild = 0

		// A key that is not 'back' or 'cycle'
		upKey := key.NewBinding(key.WithKeys("up"))
		_, _ = ca.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(upKey.Keys()[0])})

		assert.Equal(t, 1, child1.updateCount, "Focused child should have its Update method called")
		assert.Equal(t, 0, child2.updateCount, "Non-focused child should not have its Update method called for keys")
	})

	t.Run("Delegate non-key messages to all children", func(t *testing.T) {
		// Reset
		child1.updateCount = 0
		child2.updateCount = 0

		type customMsg struct{}
		_, _ = ca.Update(customMsg{})

		assert.Equal(t, 1, child1.updateCount, "Child 1 should receive the broadcast message")
		assert.Equal(t, 1, child2.updateCount, "Child 2 should receive the broadcast message")
	})
}

func TestContentArea_View(t *testing.T) {
	ca := widget.NewContentArea(1, "test")
	mockL := &mockLayout{viewStringToReturn: "inner-layout-view"}
	ca.Layout = mockL
	ca.SetSyle(lipgloss.NewStyle().Border(lipgloss.NormalBorder()))
	child := &mockModel{id: 1}
	ca.Children = []tea.Model{child}
	ca.UpdateSize(20, 5)

	view := ca.View()

	assert.True(t, mockL.viewCalled, "Layout.View should have been called")
	assert.True(t, strings.Contains(view, "inner-layout-view"), "Rendered view should contain the layout's view")
	assert.True(t, strings.Contains(view, "┌"), "Rendered view should contain the style's border")
}
