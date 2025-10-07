package layout

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ctrl-alt-boop/dribbler/ui"
)

var _ Manager = (*StackLayout)(nil)

type StackLayout struct {
	managerBase
	StackDirection Direction

	perModelWidth, perModelHeight    int
	HorizontalGutter, VerticalGutter string
}

func NewStackLayout(direction Direction) *StackLayout {
	return &StackLayout{
		managerBase: managerBase{
			focusPassThrough: false,
		},
		StackDirection:   direction,
		HorizontalGutter: ui.DefaultHorizontalGutter,
		VerticalGutter:   ui.DefaultVerticalGutter,
	}
}

func (s *StackLayout) SetSize(width, height int) {
	s.Width = width
	s.Height = height
}

func (s *StackLayout) GetSize() (width, height int) {
	return s.Width, s.Height
}

// Layout implements Manager.
func (s *StackLayout) Layout(models []tea.Model) []tea.Model {
	numModels := len(models)
	if numModels == 0 {
		return models
	}

	if s.StackDirection == Horizontal {
		s.perModelWidth = s.Width / numModels
		s.perModelHeight = s.Height
	} else {
		s.perModelWidth = s.Width
		s.perModelHeight = s.Height / numModels
	}

	for i, model := range models {
		msg := tea.WindowSizeMsg{Width: s.perModelWidth, Height: s.perModelHeight}
		updatedModel, _ := model.Update(msg)
		models[i] = updatedModel
	}
	return models
}

// View implements Manager.
func (s *StackLayout) View(models []tea.Model) string {
	if len(models) == 0 {
		return ""
	}

	var views []string
	for i, model := range models {
		style := lipgloss.NewStyle().
			Width(s.perModelWidth).Height(s.perModelHeight)
		modelRender := style.Render(model.View())
		if i > 0 {
			separator := s.CreateStackSeparator(s.StackDirection, s.perModelWidth, s.perModelHeight)
			views = append(views, separator)
		}

		views = append(views, modelRender)
	}

	switch s.StackDirection {
	case Horizontal:
		return lipgloss.JoinHorizontal(0, views...)
	case Vertical:
		return lipgloss.JoinVertical(0, views...)
	default:
		return ""
	}
}

// 1 argument means that size
// 2 arguments chooses based on direction, (width, height)
func (s *StackLayout) CreateStackSeparator(direction Direction, size ...int) string {
	if len(size) == 0 {
		return ""
	}
	if len(size) == 1 {
		size = append(size, size[0])
	}

	switch direction {
	case Horizontal:
		return lipgloss.PlaceVertical(
			size[1],
			lipgloss.Center,
			s.VerticalGutter,
			lipgloss.WithWhitespaceChars(s.VerticalGutter),
		)
	case Vertical:
		return lipgloss.PlaceHorizontal(
			size[0],
			lipgloss.Center,
			s.HorizontalGutter,
			lipgloss.WithWhitespaceChars(s.HorizontalGutter),
		)
	default:
		return ""
	}
}
