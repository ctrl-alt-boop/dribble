package layout

import "github.com/ctrl-alt-boop/dribbler/internal/panel"

var _ panel.Composer = (*DockLayoutComposer)(nil)

// DockLayoutComposer composes panels in a dock layout.
// Panels are docked to the edges of the available space, and the remaining
// space is filled by a central panel.
type DockLayoutComposer struct {
	panel.ComposerBase
	definitions panel.DefinitionList

	usableWidth, usableHeight int
	usableX, usableY          int
	usableNegX, usableNegY    int
}

// NewDockComposer creates a new DockLayoutComposer.
func NewDockComposer(panels panel.DefinitionList, opts ...panel.ComposerOption) *DockLayoutComposer {
	composer := &DockLayoutComposer{
		definitions: panels,
	}
	for _, opt := range opts {
		opt(composer)
	}
	composer.checkPositionsAndFix()

	return composer
}

// Compose implements panel.Composer.
func (d *DockLayoutComposer) Compose(width, height int) *panel.Composition {
	d.usableWidth = width
	d.usableHeight = height
	d.usableX = 0
	d.usableY = 0

	panelStates := make([]*panel.State, len(d.definitions))
	fillRemainingIndex := -1

	for i, definition := range d.definitions {
		if definition.Position == panel.None {
			continue
		}
		newPanelState := d.Allocate(width, height, definition)

		if newPanelState.FillRemaining {
			if fillRemainingIndex != -1 {
				panic("multiple fill remaining panels")
			}
			fillRemainingIndex = i
		}
		panelStates[i] = newPanelState

	}
	if fillRemainingIndex != -1 {
		panelStates[fillRemainingIndex].Width = d.usableWidth
		panelStates[fillRemainingIndex].Height = d.usableHeight
		panelStates[fillRemainingIndex].X = d.usableX
		panelStates[fillRemainingIndex].Y = d.usableY
	}

	return d.BuildComposition(panelStates...)
}

// Allocate calculates each panels size and position.
func (d *DockLayoutComposer) Allocate(width, height int, definition panel.Definition) *panel.State {
	state := &panel.State{}

	if definition.WidthRatio > 0.0 {
		state.Width = int(float64(width) * definition.WidthRatio)
	} else {
		state.Width = definition.Width
	}
	if definition.HeightRatio > 0.0 {
		state.Height = int(float64(height) * definition.HeightRatio)
	} else {
		state.Height = definition.Height
	}

	switch definition.Position {
	case panel.Top:
		panelHeight := min(state.Height, d.usableHeight)

		panelWidth := d.usableWidth

		state.Width = panelWidth
		state.Height = panelHeight
		state.X = d.usableX
		state.Y = 0

		d.usableY += state.Height
		d.usableHeight -= state.Height

	case panel.Bottom:
		panelHeight := min(state.Height, d.usableHeight)

		panelWidth := d.usableWidth

		state.Width = panelWidth
		state.Height = panelHeight
		state.X = d.usableX
		state.Y = height - state.Height - d.usableNegY
		d.usableNegY += state.Height
		d.usableHeight -= state.Height

	case panel.Left:
		panelWidth := min(state.Width, d.usableWidth)

		panelHeight := d.usableHeight

		state.Width = panelWidth
		state.Height = panelHeight
		state.X = 0
		state.Y = d.usableY

		d.usableX += state.Width
		d.usableWidth -= state.Width

	case panel.Right:
		panelWidth := min(state.Width, d.usableWidth)

		panelHeight := d.usableHeight

		state.Width = panelWidth
		state.Height = panelHeight
		state.X = width - state.Width - d.usableNegX
		state.Y = d.usableY

		d.usableNegX += state.Width

		d.usableWidth -= state.Width

	default:
		state.FillRemaining = true
	}

	return state
}

// checks the position of all panels, if all are unpositioned, they will be arranged: top, right, bottom, left, center
func (d *DockLayoutComposer) checkPositionsAndFix() {
	numDefs := len(d.definitions)
	unpositioned := 0
	for _, definition := range d.definitions {
		if definition.Position == panel.None {
			unpositioned++
		}
	}
	if unpositioned == numDefs {
		updated := d.definitions
		for i := range d.definitions {
			if i < 6 {
				updated[i].Position = panel.Position(i + 1)
			}
		}
		d.definitions = updated
	}
}
