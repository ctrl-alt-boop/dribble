package layout

import (
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/ctrl-alt-boop/dribbler/internal/panel"
)

var _ panel.Composer = (*PrioritySplitComposer)(nil)

// PrioritySplitComposer composes panels in a priority split layout.
// One panel takes a primary position (e.g., top, bottom, left, right)
// and the remaining panels are stacked in the remaining space.
type PrioritySplitComposer struct {
	panel.ComposerBase

	PrimarySizeRatio float64
	Position         panel.Position
	stackComposer    *StackLayoutComposer
}

// NewPrioritySplitLayout creates a new PrioritySplitComposer.
func NewPrioritySplitLayout(primaryPosition panel.Position, opts ...panel.ComposerOption) *PrioritySplitComposer {
	var stackDirection panel.Direction
	switch primaryPosition {
	case panel.Left, panel.Right:
		stackDirection = panel.South
	case panel.Top, panel.Bottom:
		stackDirection = panel.East
	}

	composer := &PrioritySplitComposer{
		PrimarySizeRatio: 0.5,
		Position:         primaryPosition,
		stackComposer:    NewStackLayoutComposer(stackDirection, opts...),
	}

	for _, opt := range opts {
		opt(composer)
	}
	composer.stackComposer.SetNumPanels(composer.GetNumPanels() - 1)
	return composer
}

// Compose implements panel.Composer.
func (p *PrioritySplitComposer) Compose(width, height int) *panel.Composition {
	if p.PrimarySizeRatio <= 0 || p.PrimarySizeRatio >= 1 {
		p.PrimarySizeRatio = 0.5 // Default to 50/50 split
	}

	// Calculate sizes for primary and secondary models

	primaryWidth, primaryHeight := width, height
	primaryX, primaryY := 0, 0

	secondaryWidth, secondaryHeight := width, height
	secondaryX, secondaryY := 0, 0

	switch p.Position {
	case panel.Left, panel.Right:
		totalUsableWidth := width

		primaryWidth = int(float64(totalUsableWidth) * p.PrimarySizeRatio)
		secondaryWidth = totalUsableWidth - primaryWidth
		if p.Position == panel.Right {
			primaryX = width - primaryWidth
		} else {
			secondaryX = width - primaryWidth
			secondaryY = 0
		}
	case panel.Top, panel.Bottom:
		totalUsableHeight := height

		primaryHeight = int(float64(totalUsableHeight) * p.PrimarySizeRatio)
		secondaryHeight = totalUsableHeight - primaryHeight
		if p.Position == panel.Bottom {
			primaryY = height - primaryHeight
		} else {
			secondaryX = 0
			secondaryY = height - primaryHeight
		}
	}
	primaryState := &panel.State{
		Width:  primaryWidth,
		Height: primaryHeight,
		X:      primaryX,
		Y:      primaryY,
	}

	// Compose secondary panels using the stack layout
	secondaryComposition := p.stackComposer.Compose(secondaryWidth, secondaryHeight)

	// Adjust secondary panel positions based on the primary panel's position
	secondaryComposition.X += secondaryX
	secondaryComposition.Y += secondaryY

	// Combine primary and secondary compositions
	allLayers := make([]*lipgloss.Layer, 0, len(secondaryComposition.Layers)+1)
	allLayers = append(allLayers, lipgloss.NewLayer("").
		Width(primaryState.Width).Height(primaryState.Height).
		X(primaryState.X).Y(primaryState.Y))
	allLayers = append(allLayers, secondaryComposition.Layers...)

	return &panel.Composition{
		Layers:      allLayers,
		PreRendered: secondaryComposition.PreRendered,
	}
}
