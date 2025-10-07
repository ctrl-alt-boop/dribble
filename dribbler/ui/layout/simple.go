package layout

import tea "github.com/charmbracelet/bubbletea"

var _ Manager = (*SimpleLayout)(nil)

type SimpleLayout struct {
	managerBase
}

func NewSimpleLayout() *SimpleLayout {
	return &SimpleLayout{
		managerBase: managerBase{
			focusPassThrough: false,
		},
	}
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
