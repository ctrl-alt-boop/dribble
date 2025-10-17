package layout

import "github.com/charmbracelet/lipgloss"

type panelDefinition struct { // FIXME: unexport
	position Position

	width, height           int
	widthRatio, heightRatio float64

	fillRemaining bool
	focusable     bool

	actualWidth, actualHeight int
	actualX, actualY          int

	topLeftChar, topRightChar, bottomLeftChar, bottomRightChar string
	topBorder, rightBorder, bottomBorder, leftBorder           bool

	alignmentX, alignmentY lipgloss.Position
}

func (l panelDefinition) GetBoundingBox() BoundingBox {
	return BoundingBox{
		TopLeft:     Coord{X: l.actualX, Y: l.actualY},
		BottomRight: Coord{X: l.actualX + l.actualWidth - 1, Y: l.actualY + l.actualHeight - 1},
	}
}

func Panel(position Position, opts ...panelOption) panelDefinition {
	definition := panelDefinition{
		position:   position,
		alignmentX: 0,
		alignmentY: 0,
		focusable:  true,
	}

	for _, option := range opts {
		option(&definition)
	}

	return definition
}

type panelsDefinition []panelDefinition

func Panels(panels ...panelDefinition) panelsDefinition {
	definition := panelsDefinition{}
	definition = append(definition, panels...)

	return definition
}

type panelOption func(*panelDefinition)

func WithSize(width, height int) panelOption {
	return func(def *panelDefinition) {
		def.width = width
		def.height = height
	}
}

func WithWidth(width int) panelOption {
	return func(def *panelDefinition) {
		def.width = width
		def.widthRatio = 0.0
	}
}

func WithHeight(height int) panelOption {
	return func(def *panelDefinition) {
		def.height = height
		def.heightRatio = 0.0
	}
}

func WithHorizontalAlignment(alignment lipgloss.Position) panelOption {
	return func(def *panelDefinition) {
		def.alignmentX = alignment
	}
}

func WithVerticalAlignment(alignment lipgloss.Position) panelOption {
	return func(def *panelDefinition) {
		def.alignmentY = alignment
	}
}

func WithWidthRatio(ratio float64) panelOption {
	if ratio < 0.0 {
		panic("ratio < 0.0")
	}
	if ratio > 1.0 {
		panic("ratio > 1.0")
	}
	return func(def *panelDefinition) {
		def.width = 0
		def.widthRatio = ratio
	}
}

func WithHeightRatio(ratio float64) panelOption {
	if ratio < 0.0 {
		panic("ratio < 0.0")
	}
	if ratio > 1.0 {
		panic("ratio > 1.0")
	}
	return func(def *panelDefinition) {
		def.height = 0
		def.heightRatio = ratio
	}
}

func Unfocusable() panelOption {
	return func(def *panelDefinition) {
		def.focusable = false
	}
}
