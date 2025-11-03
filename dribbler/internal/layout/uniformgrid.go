package layout

import (
	"github.com/ctrl-alt-boop/dribbler/internal/panel"
)

var _ panel.Composer = (*UniformGridLayoutComposer)(nil)

// UniformGridLayoutComposer composes panels in a uniform grid layout.
type UniformGridLayoutComposer struct {
	panel.ComposerBase

	numRows    int
	numColumns int
}

// NewUniformGridLayout creates a new UniformGridLayoutComposer.
func NewUniformGridLayout(numRows, numColumns int, opts ...panel.ComposerOption) *UniformGridLayoutComposer {
	composer := &UniformGridLayoutComposer{
		numRows:    numRows,
		numColumns: numColumns,
	}
	for _, opt := range opts {
		opt(composer)
	}

	return composer
}

// Compose implements panel.Composer.
func (g *UniformGridLayoutComposer) Compose(width int, height int) *panel.Composition {
	var baseCellWidth, baseCellHeight int
	var widthRemainder, heightRemainder int

	if g.numColumns > 0 {
		baseCellWidth = width / g.numColumns
		widthRemainder = width % g.numColumns
	}
	if g.numRows > 0 {
		baseCellHeight = height / g.numRows
		heightRemainder = height % g.numRows
	}

	panelStates := make([]*panel.State, g.numRows*g.numColumns)
	widths := make([]int, g.numColumns)
	heights := make([]int, g.numRows)
	for i := range widths {
		if i < widthRemainder {
			widths[i] = baseCellWidth + 1
		} else {
			widths[i] = baseCellWidth
		}
	}
	for i := range heights {
		if i < heightRemainder {
			heights[i] = baseCellHeight + 1
		} else {
			heights[i] = baseCellHeight
		}
	}
	currentY := 0
	for i := range g.numRows {
		currentX := 0
		for j := range g.numColumns {
			index := i*g.numColumns + j

			panelStates[index].Width = widths[j]
			panelStates[index].Height = heights[i]
			panelStates[index].X = currentX
			panelStates[index].Y = currentY

			currentX += widths[j]
		}
		currentY += heights[i]
	}
	return g.BuildComposition(panelStates...)
}
