package layout

import (
	"github.com/ctrl-alt-boop/dribbler/internal/panel"
)

var _ panel.Composer = (*SimpleLayoutComposer)(nil)

// SimpleLayoutComposer composes panels in a simple layout.
type SimpleLayoutComposer struct {
	panel.ComposerBase
}

// NewSimpleComposer creates a new SimpleLayoutComposer.
func NewSimpleComposer(opts ...panel.ComposerOption) *SimpleLayoutComposer {
	composer := &SimpleLayoutComposer{}
	for _, opt := range opts {
		opt(composer)
	}

	return composer
}

// Compose implements panel.Composer.
func (s *SimpleLayoutComposer) Compose(width int, height int) *panel.Composition {
	state := &panel.State{
		Width:  width,
		Height: height,
		X:      0,
		Y:      0,
	}
	return s.BuildComposition(state)
}
