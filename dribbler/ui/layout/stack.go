package layout

import (
	"slices"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ctrl-alt-boop/dribbler/logging"
)

var _ Manager = (*StackLayout)(nil)

type StackLayout struct {
	managerBase
	StackDirection Direction
}

func NewStackLayout(direction Direction, opts ...layoutOption) *StackLayout {
	return &StackLayout{
		managerBase: managerBase{
			layoutDefinition: New(
				[]panelDefinition{},
				opts...,
			),
			focusPassThrough: false,
		},
		StackDirection: direction,
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

	fullSize := 0
	switch s.StackDirection {
	case West, East:
		fullSize = s.Width
	case North, South:
		fullSize = s.Height
	}

	baseSplitSize := fullSize / numModels
	sizeRemainder := fullSize % numModels

	sizes := make([]int, numModels)
	for i := range numModels {
		if i < sizeRemainder {
			sizes[i] = baseSplitSize + 1
		} else {
			sizes[i] = baseSplitSize
		}
	}

	updatedDefinitions := make([]panelDefinition, numModels)
	current := 0
	for i := range numModels {
		switch s.StackDirection {
		case West, East:
			updatedDefinitions[i].actualWidth = sizes[i]
			updatedDefinitions[i].actualHeight = s.Height
			updatedDefinitions[i].actualX = current
			updatedDefinitions[i].actualY = 0
		case North, South:
			updatedDefinitions[i].actualWidth = s.Width
			updatedDefinitions[i].actualHeight = sizes[i]
			updatedDefinitions[i].actualX = 0
			updatedDefinitions[i].actualY = current
		default:
			updatedDefinitions[i].actualWidth = s.Width
			updatedDefinitions[i].actualHeight = s.Height
		}
		current += sizes[i]
	}

	s.layoutDefinition.panels = updatedDefinitions

	return s.layout(models)
}

// View implements Manager.
func (s *StackLayout) View(models []tea.Model) string {
	if len(models) == 0 || s.Height == 0 || s.Width == 0 {
		return lipgloss.NewStyle().Width(s.Width).Height(s.Height).Render("")
	}

	var views []string
	for i, model := range models {
		logging.GlobalLogger().Infof("model %d: %+v", i, model)
		style := s.getDefinitionStyle(i)
		modelRender := style.Render(model.View())

		views = append(views, modelRender)
	}

	switch s.StackDirection {
	case North:
		slices.Reverse(views)
		return lipgloss.JoinVertical(0, views...)
	case East:
		return lipgloss.JoinHorizontal(0, views...)
	case South:
		return lipgloss.JoinVertical(0, views...)
	case West:
		slices.Reverse(views)
		return lipgloss.JoinHorizontal(0, views...)
	default:
		return views[0]
	}
}

// // 1 argument means that size
// // 2 arguments chooses based on direction, (width, height)
// func (s *StackLayout) CreateStackSeparator(direction Direction, size ...int) string {
// 	if len(size) == 0 {
// 		return ""
// 	}
// 	if len(size) == 1 {
// 		size = append(size, size[0])
// 	}

// 	switch direction {
// 	case West, East:
// 		return lipgloss.PlaceVertical(
// 			size[1],
// 			0,
// 			s.VerticalGutter,
// 			lipgloss.WithWhitespaceChars(s.VerticalGutter),
// 		)
// 	case North, South:
// 		return lipgloss.PlaceHorizontal(
// 			size[0],
// 			0,
// 			s.HorizontalGutter,
// 			lipgloss.WithWhitespaceChars(s.HorizontalGutter),
// 		)
// 	default:
// 		return ""
// 	}
// }
