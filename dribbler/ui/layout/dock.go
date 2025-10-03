package layout

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var _ Manager = (*DockLayout)(nil)

type DockLayout struct {
	Width, Height int
	X, Y          int

	renderDefinition RenderDefinition

	usableWidth, usableHeight int
	usableX, usableY          int
}

// We build it by horizontally joining left, middle, right from top, center, bottom then joining each row
//
// One future idea is to think of it like a 3*3 grid?
func NewStructuredLayout() *DockLayout {
	return &DockLayout{}
}

func (s *DockLayout) SetSize(width, height int) {
	s.Width = width
	s.Height = height
}

// supplied models will be calculated in same order as Definitions list
func (s *DockLayout) Layout(models []tea.Model) []tea.Model {
	if len(models) == 0 {
		return models

	}

	updatedModels := models
	for i, definition := range s.renderDefinition.Definitions {
		pos := definition.Position
		if pos == None {
			continue
		}
		updatedDefinition := s.Allocate(definition)

		// updatedLayout := layout
		// updatedLayout.actualWidth = minWidth
		// updatedLayout.actualHeight = minHeight
		// updatedLayout.actualX = 0
		// updatedLayout.actualY = 0
		s.renderDefinition.Definitions[i] = updatedDefinition
	}

	for i, layout := range s.renderDefinition.Definitions {
		var model tea.Model
		if i < len(models) {
			model = models[i]
			msg := tea.WindowSizeMsg{Width: layout.actualWidth, Height: layout.actualHeight}
			updatedModel, _ := model.Update(msg)
			updatedModels[i] = updatedModel
		}
	}

	return updatedModels
}

func (s *DockLayout) Allocate(definition LayoutDefinition) LayoutDefinition {
	updated := definition
	switch definition.Position {
	case Bottom:
		updated.actualWidth = s.usableWidth
		updated.actualHeight = s.usableHeight - s.usableY
		updated.actualX = 0
		updated.actualY = s.usableY

		s.usableY += updated.actualHeight
		s.usableHeight -= updated.actualHeight
	case Left:
		updated.actualWidth = s.usableWidth - s.usableX
		updated.actualHeight = s.usableHeight
		updated.actualX = s.usableX
		updated.actualY = 0

		s.usableX += updated.actualWidth
		s.usableWidth -= updated.actualWidth
	case Right:
		updated.actualWidth = s.usableWidth - s.usableX
		updated.actualHeight = s.usableHeight
		updated.actualX = s.usableX
		updated.actualY = 0

		s.usableX += updated.actualWidth
		s.usableWidth -= updated.actualWidth
	case Top:
		updated.actualWidth = s.usableWidth
		updated.actualHeight = s.usableHeight - s.usableY
		updated.actualX = 0
		updated.actualY = s.usableY

		s.usableHeight -= updated.actualHeight
		s.usableY += updated.actualHeight
	default:
		updated.FillRemaining = true
	}
	return updated
}

// supplied models will be rendered in same order as Definitions list
func (s *DockLayout) View(models []tea.Model) string {
	if len(models) == 0 {
		return ""

	}

	var composed string
	for i, model := range models {
		definition := s.renderDefinition.Definitions[i]
		x := HorizontalPosInBox(definition.actualX, s.Width)
		y := VerticalPosInBox(definition.actualY, s.Height)
		composed += lipgloss.Place( // How Merge??
			definition.actualWidth,
			definition.actualHeight,
			lipgloss.Position(x),
			lipgloss.Position(y),
			model.View())
	}

	return lipgloss.Place(
		s.Width,
		s.Height,
		lipgloss.Center, // Which is best?
		lipgloss.Center, // Which is best?
		composed)
}

func (s *DockLayout) AddLayout(definition LayoutDefinition) {
	s.renderDefinition.Definitions = append(s.renderDefinition.Definitions, definition)
}

func (s *DockLayout) GetDefinition() RenderDefinition {
	return s.renderDefinition
}

func (s *DockLayout) GetLayout(index int) LayoutDefinition {
	return s.renderDefinition.Definitions[index]
}

// If position is not set, returns empty LayoutDefinition
func (s *DockLayout) GetLayoutForPosition(position Position) LayoutDefinition {
	if index, ok := s.renderDefinition.indexForPosition[position]; ok {
		return s.renderDefinition.Definitions[index]
	}
	return LayoutDefinition{}
}

func (s *DockLayout) SetDefinition(definition RenderDefinition) {
	s.renderDefinition = definition
}

func (s *DockLayout) SetLayout(index int, definition LayoutDefinition) {
	s.renderDefinition.Definitions[index] = definition
}

func (s *DockLayout) UpdateLayout(index int, opts ...LayoutOption) {
	for _, opt := range opts {
		opt(&s.renderDefinition.Definitions[index])
	}
}
