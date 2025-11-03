package layout

import (
	"github.com/ctrl-alt-boop/dribbler/internal/panel"
)

var _ panel.Composer = (*StackLayoutComposer)(nil)

// StackLayoutComposer composes panels in a stack layout.
type StackLayoutComposer struct {
	panel.ComposerBase

	StackDirection panel.Direction
}

// NewStackLayoutComposer creates a new StackLayoutComposer.
func NewStackLayoutComposer(direction panel.Direction, opts ...panel.ComposerOption) *StackLayoutComposer {
	composer := &StackLayoutComposer{
		StackDirection: direction,
	}
	for _, opt := range opts {
		opt(composer)
	}

	return composer
}

// Compose implements panel.Composer.
func (s *StackLayoutComposer) Compose(width, height int) *panel.Composition {
	fullSize := 0
	switch s.StackDirection {
	case panel.West, panel.East:
		fullSize = width
	case panel.North, panel.South:
		fullSize = height
	}

	baseSplitSize := fullSize / s.GetNumPanels()
	sizeRemainder := fullSize % s.GetNumPanels()

	sizes := make([]int, s.GetNumPanels())
	for i := range s.GetNumPanels() {
		if i < sizeRemainder {
			sizes[i] = baseSplitSize + 1
		} else {
			sizes[i] = baseSplitSize
		}
	}

	panelStates := make([]*panel.State, s.GetNumPanels())
	for i := range s.GetNumPanels() {
		panelStates[i] = &panel.State{}
	}

	current := 0
	for i := range s.GetNumPanels() {
		switch s.StackDirection {
		case panel.West, panel.East:
			panelStates[i].Width = sizes[i]
			panelStates[i].Height = height
			panelStates[i].X = current
			panelStates[i].Y = 0
		case panel.North, panel.South:
			panelStates[i].Width = width
			panelStates[i].Height = sizes[i]
			panelStates[i].X = 0
			panelStates[i].Y = current
		default:
			panelStates[i].Width = width
			panelStates[i].Height = height
		}
		current += sizes[i]
	}

	return s.BuildComposition(panelStates...)
}
