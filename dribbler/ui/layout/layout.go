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

/*
Position is used for positioning views and components in a layout

- 5 base positions, using only these in a structured layout leads to
a layout built around a center panel with a fullwidth top and bottom row, then a horizontally stacked middle row.
Each skipped position leads to it being filled equally by the other views.

- vertical, horizontal (row, column) based positions, when used in a structured layout
the view is expanded based on the row prefix until it meets another view

- horizontal, vertical (column, row) based positions, when used in a structured layout
the view is expanded based on the column prefix until it meets another view
*/
type Position int

const (
	None   Position = iota
	Center          // 5 based positions
	Top             // 5 based positions
	Bottom          // 5 based positions
	Left            // 5 based positions
	Right           // 5 based positions

	TopLeft      // vertical, horizontal based
	TopCenter    // vertical, horizontal based
	TopRight     // vertical, horizontal based
	MiddleLeft   // vertical, horizontal based
	MiddleCenter // vertical, horizontal based
	MiddleRight  // vertical, horizontal based
	BottomLeft   // vertical, horizontal based
	BottomCenter // vertical, horizontal based
	BottomRight  // vertical, horizontal based

	LeftTop      // horizontal, vertical based
	LeftMiddle   // horizontal, vertical based
	LeftBottom   // horizontal, vertical based
	CenterTop    // horizontal, vertical based
	CenterMiddle // horizontal, vertical based
	CenterBottom // horizontal, vertical based
	RightTop     // horizontal, vertical based
	RightMiddle  // horizontal, vertical based
	RightBottom  // horizontal, vertical based
)

const Middle = Center // 5 based positions

func (p Position) PositionInRow() Position {
	switch p {
	case Center, Left, Right:
		return p
	case TopLeft, MiddleLeft, BottomLeft, LeftTop, LeftMiddle, LeftBottom:
		return Left
	case TopCenter, MiddleCenter, BottomCenter, CenterTop, CenterMiddle, CenterBottom:
		return Center
	case TopRight, MiddleRight, BottomRight, RightTop, RightMiddle, RightBottom:
		return Right
	default:
		return None
	}
}

func (p Position) PositionInColumn() Position {
	switch p {
	case Center, Left, Right:
		return p
	case TopLeft, LeftTop, TopCenter, CenterTop, TopRight, RightTop:
		return Top
	case MiddleLeft, LeftMiddle, MiddleCenter, CenterMiddle, MiddleRight, RightMiddle:
		return Middle
	case BottomLeft, LeftBottom, BottomCenter, CenterBottom, BottomRight, RightBottom:
		return Bottom
	default:
		return None
	}
}

func (p Position) AsLipglossHorizontal() lipgloss.Position {
	horizontalPos := p.PositionInRow()
	switch horizontalPos {
	case Left:
		return lipgloss.Left
	case Right:
		return lipgloss.Right
	default:
		return lipgloss.Center
	}
}

func (p Position) AsLipglossVertical() lipgloss.Position {
	verticalPos := p.PositionInColumn()
	switch verticalPos {
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
}

type RenderDefinition struct {
	Layouts []LayoutDefinition

	Width, Height int

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
		Layouts: make([]LayoutDefinition, len(definitions)),
	}

	copy(renderDefinition.Layouts, definitions)

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
		for i := range renderModel.Layouts {
			renderModel.Layouts[i].FillRemaining = false
		}
		renderModel.Layouts[renderModel.indexForPosition[position]].FillRemaining = true
	}
}

func WithSize(width, height int) Option {
	return func(renderModel *RenderDefinition) {
		renderModel.Width = width
		renderModel.Height = height
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
