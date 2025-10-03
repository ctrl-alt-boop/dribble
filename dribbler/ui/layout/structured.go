package layout

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var _ Manager = (*StructuredLayout)(nil)

type StructuredLayout struct {
	Width, Height int
	X, Y          int

	renderDefinition RenderDefinition
}

// We build it by horizontally joining left, middle, right from top, center, bottom then joining each row
//
// One future idea is to think of it like a 3*3 grid?
func NewStructuredLayout() *StructuredLayout {
	return &StructuredLayout{}
}

func (s *StructuredLayout) SetSize(width, height int) {
	s.Width = width
	s.Height = height
}

// Built by horizontally joining left, middle, right from top, center, bottom then joining each row
// supplied models will be calculated in same order as Definitions list
// in case of duplicated positions, ...
func (s *StructuredLayout) Layout(models []tea.Model) []tea.Model {
	if len(models) == 0 {
		return models

	}

	models = models[:len(s.renderDefinition.Layouts)]

	updatedModels := make([]tea.Model, len(models))
	for i, model := range models { // Add calculations based on positions
		layout := s.renderDefinition.Layouts[i]

		// Apply min/max constraints
		width := layout.MaxWidth
		height := layout.MaxHeight

		msg := tea.WindowSizeMsg{Width: width, Height: height}
		updatedModel, _ := model.Update(msg)
		updatedModels[i] = updatedModel

		updatedDef := layout
		updatedDef.actualWidth = width
		updatedDef.actualHeight = height
		s.renderDefinition.Layouts[i] = updatedDef
	}

	return updatedModels
}

// Rendered by horizontally joining left, middle, right from top, center, bottom then joining each row
// supplied models will be rendered in same order as Definitions list
func (s *StructuredLayout) View(models []tea.Model) string {
	if len(models) == 0 {
		return ""

	}

	models = models[:len(s.renderDefinition.Layouts)]

	rows := []string{}

	topRow := []string{}
	middleRow := []string{}
	bottomRow := []string{}

	rows = append(rows, lipgloss.JoinHorizontal(lipgloss.Top, topRow...))
	rows = append(rows, lipgloss.JoinHorizontal(lipgloss.Top, middleRow...))
	rows = append(rows, lipgloss.JoinHorizontal(lipgloss.Bottom, bottomRow...))

	return lipgloss.JoinVertical(lipgloss.Left, rows...)
}

func (s *StructuredLayout) AddLayout(definition LayoutDefinition) {
	s.renderDefinition.Layouts = append(s.renderDefinition.Layouts, definition)
}

func (s *StructuredLayout) GetDefinition() RenderDefinition {
	return s.renderDefinition
}

func (s *StructuredLayout) GetLayout(index int) LayoutDefinition {
	return s.renderDefinition.Layouts[index]
}

// If position is not set, returns empty LayoutDefinition
func (s *StructuredLayout) GetLayoutForPosition(position Position) LayoutDefinition {
	if index, ok := s.renderDefinition.indexForPosition[position]; ok {
		return s.renderDefinition.Layouts[index]
	}
	return LayoutDefinition{}
}

func (s *StructuredLayout) SetDefinition(definition RenderDefinition) {
	s.renderDefinition = definition
}

func (s *StructuredLayout) SetLayout(index int, definition LayoutDefinition) {
	s.renderDefinition.Layouts[index] = definition
}

func (s *StructuredLayout) UpdateLayout(index int, opts ...LayoutOption) {
	for _, opt := range opts {
		opt(&s.renderDefinition.Layouts[index])
	}
}
