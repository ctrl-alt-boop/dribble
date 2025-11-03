package panel

import (
	"github.com/charmbracelet/lipgloss/v2"
)

type Composer interface {
	Compose(width, height int) *Composition

	SetStyle(style lipgloss.Style)
	SetFocusedStyle(style lipgloss.Style)
	SetPanelBorder(border lipgloss.Border)

	GetStyle() lipgloss.Style
	GetFocusedStyle() lipgloss.Style
	GetPanelBorder() lipgloss.Border

	SetNumPanels(numPanels int)
	GetNumPanels() int
}

// ComposerBase implements shared logic for Composer
type ComposerBase struct {
	numPanels int

	normalStyle  lipgloss.Style
	focusedStyle lipgloss.Style
	panelBorder  lipgloss.Border
}

func (b *ComposerBase) SetNumPanels(n int) {
	b.numPanels = n
}

func (b *ComposerBase) GetNumPanels() int {
	return b.numPanels
}

func (b *ComposerBase) SetStyle(style lipgloss.Style) {
	b.normalStyle = style
}

func (b *ComposerBase) SetFocusedStyle(style lipgloss.Style) {
	b.focusedStyle = style
}

func (b *ComposerBase) SetPanelBorder(border lipgloss.Border) {
	b.panelBorder = border
}

func (b *ComposerBase) GetStyle() lipgloss.Style {
	return b.normalStyle
}

func (b *ComposerBase) GetFocusedStyle() lipgloss.Style {
	return b.focusedStyle
}

func (b *ComposerBase) GetPanelBorder() lipgloss.Border {
	return b.panelBorder
}

func (b *ComposerBase) BuildComposition(panelStates ...*State) *Composition {
	layers := make([]*lipgloss.Layer, len(panelStates))

	for i, state := range panelStates {
		layer := lipgloss.NewLayer("").
			Width(state.Width).Height(state.Height).
			X(state.X).Y(state.Y)
		layers[i] = layer
	}

	return &Composition{
		Layers:      layers,
		PreRendered: []*lipgloss.Layer{},
	}
}

type Composition struct {
	Layers      []*lipgloss.Layer
	PreRendered []*lipgloss.Layer

	X, Y int
}

type ComposerOption func(Composer)

func WithStyle(style lipgloss.Style) ComposerOption {
	return func(c Composer) {
		c.SetStyle(style)
	}
}

func WithFocusedStyle(style lipgloss.Style) ComposerOption {
	return func(c Composer) {
		c.SetFocusedStyle(style)
	}
}

func WithPanelBorder(border lipgloss.Border) ComposerOption {
	return func(c Composer) {
		c.SetPanelBorder(border)
	}
}
