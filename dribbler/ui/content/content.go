// Package content contain data structs for UI widgets and components
package content

import (
	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ctrl-alt-boop/dribbler/logging"
)

type (
	Selection interface {
		CursorX() int
		CursorY() int
		Cursor() (int, int)

		SetCursor(x, y int)
		MoveCursor(dX, dY int)

		MoveCursorUp(y ...int)
		MoveCursorDown(y ...int)
		MoveCursorLeft(x ...int)
		MoveCursorRight(x ...int)

		GetSelected() any
	}

	Container struct {
		Content any
		Extras  []any
	}
	Base struct {
		ID   int
		Name string

		model tea.Model
	}
)

var DefaultStyle = lipgloss.NewStyle().Margin(1, 2)

// Totally unnecessary
func ListToString[L ~[]T, T any, F func(T) string](list L, fn F) []string {
	out := make([]string, len(list))
	for i := range list {
		out[i] = fn(list[i])
	}
	return out
}

func NewContainer(content any, extras ...any) Container {
	return Container{
		Content: content,
		Extras:  extras,
	}
}

func (c *Container) SetContent(content any, extras ...any) {
	c.Content = content
	c.Extras = extras
}

func (c Container) Init() tea.Cmd {
	switch model := c.Content.(type) {
	case tea.Model:
		return model.Init()
	}
	return nil
}

func (c Container) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	updated := c

	var cmd tea.Cmd
	var content any

	switch model := c.Content.(type) {
	case tea.Model:
		content, cmd = model.Update(msg)
	}
	updated.Content = content
	return updated, cmd
}

func (c Container) View() string {
	switch model := c.Content.(type) {
	case tea.Model:
		return model.View()
	case help.Model:
		helpKeys := c.Extras[0].(help.KeyMap)
		return model.View(helpKeys)
	}
	return ""
}

type (
	// Initable to check if Content implements tea.Model.Init
	Initable interface {
		Init() tea.Cmd
	}
	// Updatable to check if Content implements tea.Model.Update
	Updatable interface {
		Update(msg tea.Msg) (tea.Model, tea.Cmd)
	}
	// Viewable to check if Content implements tea.Model.View
	Viewable interface {
		View() string
	}
	// FieldUpdatable to check if Content implements SetFields
	// Should return itself with the updates
	FieldUpdatable interface {
		UpdateFields(...any) any
	}
	// ParamContainer is a type parameterized Container
	ParamContainer[T any] struct {
		Content T
		Extras  []any
	}
)

// NewParamContainer creates a new ParamContainer
func NewParamContainer[T any](content T, extras ...any) ParamContainer[T] {
	return ParamContainer[T]{
		Content: content,
		Extras:  extras,
	}
}

// Init implements tea.Model.
func (c ParamContainer[T]) Init() tea.Cmd {
	if initable, ok := any(c.Content).(Initable); ok {
		return initable.Init()
	}

	return nil
}

// Update implements tea.Model.
func (c ParamContainer[T]) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	updated := c
	if fieldSettable, ok := any(c.Content).(FieldUpdatable); ok {
		updatedContent, ok := fieldSettable.UpdateFields(c.Extras...).(T)
		if ok {
			updated.Content = updatedContent
		} else {
			logging.GlobalLogger().Warnf("Something strange happened when trying fieldSettable.UpdateFields(c.Extras...).(T)")
		}
	}
	if updatable, ok := any(updated).(Updatable); ok {
		updatedContent, cmd := updatable.Update(msg)
		updated.Content = updatedContent.(T)

		return updated, cmd
	}
	return updated, nil
}

// View implements tea.Model.
func (c ParamContainer[T]) View() string {
	if viewable, ok := any(c.Content).(Viewable); ok {
		return viewable.View()
	}
	return ""
}
