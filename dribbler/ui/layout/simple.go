package layout

import tea "github.com/charmbracelet/bubbletea"

var _ Manager = (*SimpleLayout)(nil)

type SimpleLayout struct {
	Width, Height    int
	renderDefinition RenderDefinition
}

func NewSimpleLayout() *SimpleLayout {
	return &SimpleLayout{}
}

func (s *SimpleLayout) SetSize(width, height int) {
	s.Width = width
	s.Height = height
}

func (s *SimpleLayout) GetSize() (width, height int) {
	return s.Width, s.Height
}

func (s SimpleLayout) Layout(models []tea.Model) []tea.Model {
	if len(models) == 0 {
		return models
	}

	// Issue a warning if there are too many children (optional, but helpful)
	if len(models) > 1 {
		// TODO: LOGGING
	}

	child := models[0]

	msg := tea.WindowSizeMsg{Width: s.Width, Height: s.Height}

	updatedChild, _ := child.Update(msg)

	models[0] = updatedChild
	return models
}

func (s *SimpleLayout) View(models []tea.Model) string {
	if len(models) == 0 {
		return ""
	}

	// Only render the first model
	return models[0].View()
}

func (s *SimpleLayout) AddLayout(definition LayoutDefinition) {
	s.renderDefinition.Definitions = append(s.renderDefinition.Definitions, definition)
}

func (s *SimpleLayout) GetDefinition() RenderDefinition {
	return s.renderDefinition
}

func (s *SimpleLayout) GetLayout(index int) LayoutDefinition {
	return s.renderDefinition.Definitions[index]
}

// If position is not set, returns empty LayoutDefinition
func (s *SimpleLayout) GetLayoutForPosition(position Position) LayoutDefinition {
	if index, ok := s.renderDefinition.indexForPosition[position]; ok {
		return s.renderDefinition.Definitions[index]
	}
	return LayoutDefinition{}
}

func (s *SimpleLayout) SetDefinition(definition RenderDefinition) {
	s.renderDefinition = definition
}

func (s *SimpleLayout) SetLayout(index int, definition LayoutDefinition) {
	s.renderDefinition.Definitions[index] = definition
}

func (s *SimpleLayout) UpdateLayout(index int, opts ...LayoutOption) {
	for _, opt := range opts {
		opt(&s.renderDefinition.Definitions[index])
	}
}
