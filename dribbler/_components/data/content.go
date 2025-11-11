// Package data contain data structs for UI widgets and components
package data

import (
	"charm.land/bubbles/v2/help"
	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"
)

// Base for components, implements Model except Render and Update
type Base struct{}

// Init implements Model
func (b Base) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model
func (b Base) Update(_ tea.Msg) (tea.Model, tea.Cmd) {
	return b, nil
}

// View implements tea.Model
func (b Base) View() tea.View { // FIXME:
	return tea.NewView("")
}

type (
	Container struct {
		Content any
		Extras  []any
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

func (c Container) Update(msg tea.Msg) (Container, tea.Cmd) {
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

func (c Container) Render() string {
	switch model := c.Content.(type) {
	case interface{ Render() string }:
		return model.Render()
	case help.Model:
		helpKeys := c.Extras[0].(help.KeyMap)
		return model.View(helpKeys)
	}
	return ""
}
