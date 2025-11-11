package layout

import (
	lipgloss "charm.land/lipgloss/v2"
	"github.com/ctrl-alt-boop/dribbler/internal/panel"
)

func Panel(opts ...panelOption) panel.Definition {
	definition := panel.Definition{
		Position:   panel.None,
		AlignmentX: 0,
		AlignmentY: 0,
		Focusable:  false,
	}

	for _, option := range opts {
		option(&definition)
	}

	return definition
}

func Panels(panels ...panel.Definition) panel.DefinitionList {
	definition := panel.DefinitionList{}
	definition = append(definition, panels...)

	return definition
}

type panelOption func(*panel.Definition)

func WithPosition(position panel.Position) panelOption {
	return func(def *panel.Definition) {
		def.Position = position
	}
}

func WithInnerSize(width, height int) panelOption {
	return func(def *panel.Definition) {
		def.Width = width + frameThickness
		def.Height = height + frameThickness
		def.WidthRatio = 0.0
		def.HeightRatio = 0.0
	}
}

func WithSize(width, height int) panelOption {
	return func(def *panel.Definition) {
		def.Width = width
		def.Height = height
		def.WidthRatio = 0.0
		def.HeightRatio = 0.0
	}
}

func WithInnerWidth(width int) panelOption {
	return func(def *panel.Definition) {
		def.Width = width + frameThickness
		def.WidthRatio = 0.0
	}
}

func WithInnerHeight(height int) panelOption {
	return func(def *panel.Definition) {
		def.Height = height + frameThickness
		def.HeightRatio = 0.0
	}
}

func WithWidth(width int) panelOption {
	return func(def *panel.Definition) {
		def.Width = width
		def.WidthRatio = 0.0
	}
}

func WithHeight(height int) panelOption {
	return func(def *panel.Definition) {
		def.Height = height
		def.HeightRatio = 0.0
	}
}

func WithHorizontalAlignment(alignment lipgloss.Position) panelOption {
	return func(def *panel.Definition) {
		def.AlignmentX = alignment
	}
}

func WithVerticalAlignment(alignment lipgloss.Position) panelOption {
	return func(def *panel.Definition) {
		def.AlignmentY = alignment
	}
}

func WithWidthRatio(ratio float64) panelOption {
	if ratio < 0.0 {
		panic("ratio < 0.0")
	}
	if ratio > 1.0 {
		panic("ratio > 1.0")
	}
	return func(def *panel.Definition) {
		def.Width = 0
		def.WidthRatio = ratio
	}
}

func WithHeightRatio(ratio float64) panelOption {
	if ratio < 0.0 {
		panic("ratio < 0.0")
	}
	if ratio > 1.0 {
		panic("ratio > 1.0")
	}
	return func(def *panel.Definition) {
		def.Height = 0
		def.HeightRatio = ratio
	}
}

func Focusable() panelOption {
	return func(def *panel.Definition) {
		def.Focusable = true
	}
}
