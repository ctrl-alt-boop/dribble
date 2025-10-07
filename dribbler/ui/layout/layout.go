package layout

import (
	"github.com/charmbracelet/lipgloss"
)

type Direction int

const (
	Horizontal Direction = iota
	Vertical
)

func CreateConnectedCorners(style lipgloss.Style) lipgloss.Border {
	border := style.GetBorderStyle()
	border.TopLeft = border.MiddleTop
	border.TopRight = border.MiddleTop
	border.BottomLeft = border.MiddleBottom
	border.BottomRight = border.MiddleBottom
	return border
}

type LayoutDefinition struct {
	Position Position

	MinWidth, MinHeight int
	MaxWidth, MaxHeight int
	FillRemaining       bool

	actualWidth, actualHeight int
	actualX, actualY          int
}

type RenderDefinition struct {
	Definitions []LayoutDefinition

	PanelBorder                  lipgloss.Border
	FocusedStyle, UnfocusedStyle lipgloss.Style

	indexForPosition      map[Position]int
	hasFocusUnfocusStyles bool
	customPanelBorder     bool
	connectCorners        bool
}

func (r *RenderDefinition) Update(opts ...Option) {
	for _, option := range opts {
		option(r)
	}
}

func (r RenderDefinition) PositionsInUse() []Position {
	positions := make([]Position, 0, len(r.indexForPosition))
	for position := range r.indexForPosition {
		positions = append(positions, position)
	}
	return positions
}

func NewDefinition(definitions []LayoutDefinition, opts ...Option) RenderDefinition {
	definition := RenderDefinition{
		Definitions:       make([]LayoutDefinition, len(definitions)),
		PanelBorder:       lipgloss.NormalBorder(),
		FocusedStyle:      lipgloss.NewStyle().Border(lipgloss.NormalBorder()),
		customPanelBorder: false,
		connectCorners:    true,
		indexForPosition:  make(map[Position]int),
	}

	copy(definition.Definitions, definitions)

	for _, option := range opts {
		option(&definition)
	}

	if definition.customPanelBorder || definition.connectCorners {
		if definition.PanelBorder.MiddleTop == "" || definition.PanelBorder.MiddleBottom == "" {
			panic("custom panel border requires MiddleTop and MiddleBottom")
		}
	}

	for i, def := range definition.Definitions {
		definition.indexForPosition[def.Position] = i
	}

	return definition
}

func NewLayoutDefinition(position Position, opts ...LayoutOption) LayoutDefinition {
	definition := LayoutDefinition{
		Position: position,
	}

	for _, option := range opts {
		option(&definition)
	}

	return definition
}

type Option func(*RenderDefinition)
type LayoutOption func(*LayoutDefinition)

func WithMaxSize(width, height int) LayoutOption {
	return func(def *LayoutDefinition) {
		def.MaxWidth = width
		def.MaxHeight = height
	}
}

func WithMinSize(width, height int) LayoutOption {
	return func(def *LayoutDefinition) {
		def.MinWidth = width
		def.MinHeight = height
	}
}

func WithMaxWidth(width int) LayoutOption {
	return func(def *LayoutDefinition) {
		def.MaxWidth = width
	}
}

func WithMaxHeight(height int) LayoutOption {
	return func(def *LayoutDefinition) {
		def.MaxHeight = height
	}
}

func WithMinWidth(width int) LayoutOption {
	return func(def *LayoutDefinition) {
		def.MinWidth = width
	}
}

func WithMinHeight(height int) LayoutOption {
	return func(def *LayoutDefinition) {
		def.MinHeight = height
	}
}

func WithFillRemaining() LayoutOption {
	return func(def *LayoutDefinition) {
		def.FillRemaining = true
	}
}

func FillRemainingAt(position Position) Option {
	return func(renderModel *RenderDefinition) {
		if position == None {
			return
		}
		if _, ok := renderModel.indexForPosition[position]; !ok {
			return
		}
		for i := range renderModel.Definitions {
			renderModel.Definitions[i].FillRemaining = false
		}
		renderModel.Definitions[renderModel.indexForPosition[position]].FillRemaining = true
	}
}

// The custom border uses MiddleTop and MiddleBottom if connectCorners is true
func WithPanelBorder(border lipgloss.Border) Option {
	return func(renderModel *RenderDefinition) {
		renderModel.PanelBorder = border
		renderModel.customPanelBorder = true
	}
}

func WithoutConnectCorners() Option {
	return func(renderModel *RenderDefinition) {
		renderModel.connectCorners = false
	}
}

func WithStyle(style lipgloss.Style) Option {
	return func(renderModel *RenderDefinition) {
		renderModel.hasFocusUnfocusStyles = false
		renderModel.FocusedStyle = style
		renderModel.UnfocusedStyle = lipgloss.NewStyle()
	}
}

func WithFocusUnfocusStyles(focused, unfocused lipgloss.Style) Option {
	return func(renderModel *RenderDefinition) {
		renderModel.hasFocusUnfocusStyles = true
		renderModel.FocusedStyle = focused
		renderModel.UnfocusedStyle = unfocused
	}
}
