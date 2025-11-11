package panel

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"
)

type DefinitionList []Definition

type Definition struct {
	ID                      string
	Name                    string
	Position                Position
	Width, Height           int
	WidthRatio, HeightRatio float64

	ShouldFillRemaining bool
	Focusable           bool

	AlignmentX, AlignmentY lipgloss.Position
}

type State struct {
	Width, Height int
	X, Y          int

	// IsFocused     bool
	FillRemaining bool
}

type Model interface {
	ID() string
	Name() string
	Init() tea.Cmd
	Update(msg tea.Msg) (Model, tea.Cmd)
	Canvas() *lipgloss.Canvas
	IsFocused() bool
	SetFocused(bool)
	SetBorderStyle(panelBorder lipgloss.Border)
}

type model struct {
	content any

	id        string
	name      string
	isFocused bool

	Definition                                                 Definition
	TopLeftChar, TopRightChar, BottomLeftChar, BottomRightChar string
	TopBorder, RightBorder, BottomBorder, LeftBorder           bool

	borderStyle lipgloss.Border
}

func NewPanel(content any, opts ...PanelOption) *model {
	panel := &model{
		content: content,
	}

	for _, opt := range opts {
		opt(panel)
	}

	return panel
}

type PanelOption func(*model)

func WithID(id string) PanelOption {
	return func(p *model) {
		p.id = id
	}
}

func WithName(name string) PanelOption {
	return func(p *model) {
		p.name = name
	}
}

func (m *model) ID() string {
	return m.id
}

func (m *model) Name() string {
	return m.name
}

func (m *model) Init() tea.Cmd { return nil }

func (m *model) Update(msg tea.Msg) (Model, tea.Cmd) {
	updated := m

	return updated, nil
}

func (m *model) IsFocused() bool {
	return m.isFocused
}

func (m *model) SetFocused(f bool) {
	m.isFocused = f
}

func (m *model) SetBorderStyle(panelBorder lipgloss.Border) {
	m.borderStyle = panelBorder
}

func (m model) Canvas() *lipgloss.Canvas {
	switch content := m.content.(type) {
	case string:
		return lipgloss.NewCanvas(lipgloss.NewLayer(content))
	case interface{ Render() string }:
		return lipgloss.NewCanvas(lipgloss.NewLayer(content.Render()))
	case Model:
		return content.Canvas()
	}
	return lipgloss.NewCanvas(lipgloss.NewLayer(fmt.Sprintf("%s", m.content)))
}

func GetBoundingBox(panel *State) BoundingBox {
	return BoundingBox{
		TopLeft:     Coord{X: panel.X, Y: panel.Y},
		BottomRight: Coord{X: panel.X + panel.Width - 1, Y: panel.Y + panel.Height - 1},
	}
}

func GetAllBoundingBoxes(panels ...*State) []BoundingBox {
	var boundingBoxes []BoundingBox
	for _, panel := range panels {
		boundingBoxes = append(boundingBoxes, GetBoundingBox(panel))
	}
	return boundingBoxes
}
