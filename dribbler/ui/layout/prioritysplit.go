package layout

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ctrl-alt-boop/dribbler/ui"
)

var _ Manager = (*PrioritySplitLayout)(nil)

type PrioritySplitLayout struct {
	managerBase
	Ratio                            float64
	Direction                        Direction
	HorizontalGutter, VerticalGutter string

	// I'm not entierly sure how I want to do this, either letting the content decide or Layout decide...
	// One option is to reset the style of Children and using the one the parent decides

	primaryWidth, primaryHeight                 int
	secondaryWidth, secondaryHeight             int
	secondaryWidthForOne, secondaryHeightForOne int
	stackLayout                                 *StackLayout
}

func NewPrioritySplitLayout(direction Direction) *PrioritySplitLayout {
	layout := &PrioritySplitLayout{
		managerBase: managerBase{
			focusPassThrough: false,
		},
		Ratio:            0.6,
		Direction:        direction,
		HorizontalGutter: ui.DefaultHorizontalGutter,
		VerticalGutter:   ui.DefaultVerticalGutter,
		stackLayout:      NewStackLayout(direction),
	}

	return layout
}

func (p *PrioritySplitLayout) SetSize(width, height int) {
	p.Width = width
	p.Height = height
}

func (p *PrioritySplitLayout) GetSize() (width, height int) {
	return p.Width, p.Height
}

// Layout implements Manager.
func (p *PrioritySplitLayout) Layout(models []tea.Model) []tea.Model {
	if len(models) == 0 {
		return models
	}

	updatedModels := make([]tea.Model, len(models))

	// if p.Ratio <= 0 || p.Ratio >= 1 {
	// 	p.Ratio = 0.5 // Default to 50/50 split
	// }

	// Calculate sizes for primary and secondary models

	secondaryCount := len(models) - 1

	switch p.Direction {
	case Horizontal:
		p.primaryWidth = p.Width
		p.secondaryWidth = p.Width

		totalUsableHeight := p.Height - 1
		p.primaryHeight = int(float64(totalUsableHeight) * p.Ratio)
		p.secondaryHeight = totalUsableHeight - p.primaryHeight
		if secondaryCount > 0 {
			p.secondaryHeightForOne = p.secondaryHeight / secondaryCount
		}
	case Vertical:
		p.primaryHeight = p.Height
		p.secondaryHeight = p.Height

		totalUsableWidth := p.Width - 1
		p.primaryWidth = int(float64(totalUsableWidth) * p.Ratio)
		p.secondaryWidth = totalUsableWidth - p.primaryWidth
		if secondaryCount > 0 {
			p.secondaryWidthForOne = p.secondaryWidth / secondaryCount
		}
	}

	updatedPrimary, _ := models[0].Update(tea.WindowSizeMsg{Width: p.primaryWidth, Height: p.primaryHeight})
	updatedModels[0] = updatedPrimary

	p.stackLayout.SetSize(p.secondaryWidth, p.secondaryHeight)

	updatedSecondaryModels := p.stackLayout.Layout(models[1:])

	for i, updatedModel := range updatedSecondaryModels {
		updatedModels[i+1] = updatedModel
	}

	return updatedModels
}

func (p *PrioritySplitLayout) View(models []tea.Model) string {
	return p.StackView(models)
}

func (p *PrioritySplitLayout) StackView(models []tea.Model) string {
	if len(models) == 0 {
		return ""
	}

	secondaryModels := models[1:]

	stackView := p.stackLayout.View(secondaryModels)

	primaryRender := lipgloss.NewStyle().
		Width(p.primaryWidth).Height(p.primaryHeight).
		Render(models[0].View())

	primarySeparator := p.CreatePrimarySeparator(p.Direction, p.Width, p.Height)
	switch p.Direction {
	case Horizontal:
		return lipgloss.JoinVertical(0, primaryRender, primarySeparator, stackView)

	case Vertical:
		return lipgloss.JoinHorizontal(0, primaryRender, primarySeparator, stackView)

	default:
		return ""
	}
}

// 1 argument means that size
// 2 arguments chooses based on direction, (width, height)
func (p *PrioritySplitLayout) CreatePrimarySeparator(direction Direction, size ...int) string {
	if len(size) == 0 {
		return ""
	}
	if len(size) == 1 {
		size = append(size, size[0])
	}

	switch direction {
	case Horizontal:
		return lipgloss.PlaceHorizontal(
			size[0],
			lipgloss.Center,
			p.HorizontalGutter,
			lipgloss.WithWhitespaceChars(p.HorizontalGutter),
		)
	case Vertical:
		return lipgloss.PlaceVertical(
			size[1],
			lipgloss.Center,
			p.VerticalGutter,
			lipgloss.WithWhitespaceChars(p.VerticalGutter),
		)
	default:
		return ""
	}
}
