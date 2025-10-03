package layout

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Direction int

const (
	Horizontal Direction = iota
	Vertical
)

type Position int

const (
	None Position = iota
	Center
	Top
	Bottom
	Left
	Right
)

const Middle = Center

func (p Position) AsLipglossHorizontal() lipgloss.Position {
	switch p {
	case Left:
		return lipgloss.Left
	case Right:
		return lipgloss.Right
	default:
		return lipgloss.Center
	}
}

func (p Position) AsLipglossVertical() lipgloss.Position {
	switch p {
	case Top:
		return lipgloss.Top
	case Bottom:
		return lipgloss.Bottom
	default:
		return lipgloss.Center
	}
}

type Manager interface {
	SetDefinition(definition RenderDefinition)
	GetDefinition() RenderDefinition

	AddLayout(definition LayoutDefinition)

	UpdateLayout(index int, opts ...LayoutOption)
	SetLayout(index int, definition LayoutDefinition)
	GetLayout(index int) LayoutDefinition

	// If position is not set, returns empty LayoutDefinition
	GetLayoutForPosition(position Position) LayoutDefinition

	Layout(models []tea.Model) []tea.Model
	View(models []tea.Model) string

	SetSize(width, height int)
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

	PanelBorders                 lipgloss.Border
	FocusedStyle, UnfocusedStyle lipgloss.Style

	indexForPosition      map[Position]int
	hasFocusUnfocusStyles bool
}

func (r *RenderDefinition) Update(opts ...Option) {
	for _, option := range opts {
		option(r)
	}
}

func NewDefinition(definitions []LayoutDefinition, opts ...Option) RenderDefinition {
	renderDefinition := RenderDefinition{
		Definitions: make([]LayoutDefinition, len(definitions)),
	}

	copy(renderDefinition.Definitions, definitions)

	for _, option := range opts {
		option(&renderDefinition)
	}

	return renderDefinition
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

// func WithSize(width, height int) Option {
// 	return func(renderModel *RenderDefinition) {
// 		renderModel.Width = width
// 		renderModel.Height = height
// 	}
// }

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
