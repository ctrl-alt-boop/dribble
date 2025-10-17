package layout

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var _ Manager = (*PrioritySplitLayout)(nil)

type PrioritySplitLayout struct { // FIXME: Border junctions
	managerBase
	PrimarySizeRatio float64
	Position         Position

	stackLayout *StackLayout
}

func NewPrioritySplitLayout(primaryPosition Position, opts ...layoutOption) *PrioritySplitLayout {
	var direction Direction = 0
	switch primaryPosition {
	case Left, Right:
		direction = South
	case Top, Bottom:
		direction = East
	}
	layout := &PrioritySplitLayout{
		managerBase: managerBase{
			layoutDefinition: New(
				[]panelDefinition{},
				opts...,
			),
			focusPassThrough: true,
			focusedIndex:     -1,
		},
		PrimarySizeRatio: 0.5,
		Position:         primaryPosition,
		stackLayout:      NewStackLayout(direction, opts...),
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

	if p.PrimarySizeRatio <= 0 || p.PrimarySizeRatio >= 1 {
		p.PrimarySizeRatio = 0.5 // Default to 50/50 split
	}

	// Calculate sizes for primary and secondary models

	primaryWidth, primaryHeight := p.Width, p.Height
	primaryX, primaryY := 0, 0

	secondaryWidth, secondaryHeight := p.Width, p.Height
	secondaryX, secondaryY := 0, 0

	switch p.Position {
	case Left, Right:
		totalUsableWidth := p.Width

		primaryWidth = int(float64(totalUsableWidth) * p.PrimarySizeRatio)
		secondaryWidth = totalUsableWidth - primaryWidth
		if p.Position == Right {
			primaryX = p.Width - primaryWidth
		} else {
			secondaryX = p.Width - primaryWidth
			secondaryY = 0
		}
	case Top, Bottom:
		totalUsableHeight := p.Height

		primaryHeight = int(float64(totalUsableHeight) * p.PrimarySizeRatio)
		secondaryHeight = totalUsableHeight - primaryHeight
		if p.Position == Bottom {
			primaryY = p.Height - primaryHeight
		} else {
			secondaryX = 0
			secondaryY = p.Height - primaryHeight
		}
	}

	p.layoutDefinition.panels = []panelDefinition{
		{
			actualWidth:  primaryWidth,
			actualHeight: primaryHeight,
			actualX:      primaryX,
			actualY:      primaryY,
		},
		{
			actualWidth:  secondaryWidth,
			actualHeight: secondaryHeight,
			actualX:      secondaryX,
			actualY:      secondaryY,
		},
	}

	updatedModels := p.layout(models[:1])

	p.stackLayout.SetSize(secondaryWidth, secondaryHeight)
	updatedSecondaryModels := p.stackLayout.Layout(models[1:])

	updatedModels = append(updatedModels, updatedSecondaryModels...)

	return updatedModels
}

func (p *PrioritySplitLayout) View(models []tea.Model) string {
	if len(models) == 0 || p.Height == 0 || p.Width == 0 {
		return lipgloss.NewStyle().Width(p.Width).Height(p.Height).Render("")
	}
	primaryRender := p.getDefinitionStyle(0).Render(models[0].View())

	if p.focusedIndex > 0 {
		p.stackLayout.focusedIndex = p.focusedIndex - 1
	}
	stackView := p.stackLayout.View(models[1:])

	switch p.Position {
	case Top:
		return lipgloss.JoinVertical(0, primaryRender, stackView)
	case Right:
		return lipgloss.JoinHorizontal(0, stackView, primaryRender)
	case Bottom:
		return lipgloss.JoinVertical(0, stackView, primaryRender)
	case Left:
		return lipgloss.JoinHorizontal(0, primaryRender, stackView)
	default:
		return ""
	}
}
