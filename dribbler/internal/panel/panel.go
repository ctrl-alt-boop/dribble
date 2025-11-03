package panel

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
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

	TopLeftChar, TopRightChar, BottomLeftChar, BottomRightChar string
	TopBorder, RightBorder, BottomBorder, LeftBorder           bool
}

type Panel interface {
	ID() string
	Name() string
	Init() tea.Cmd
	Update(msg tea.Msg) (Panel, tea.Cmd)
	Render() string
	IsFocused() bool
	SetFocused(bool)
}

type Model struct {
	content any

	id        string
	name      string
	isFocused bool

	Definition Definition
}

func NewPanel(content any, opts ...PanelOption) *Model {
	panel := &Model{
		content: content,
	}

	for _, opt := range opts {
		opt(panel)
	}

	return panel
}

type PanelOption func(*Model)

func WithID(id string) PanelOption {
	return func(p *Model) {
		p.id = id
	}
}

func WithName(name string) PanelOption {
	return func(p *Model) {
		p.name = name
	}
}

func (p *Model) ID() string {
	return p.id
}

func (p *Model) Name() string {
	return p.name
}

func (p *Model) Init() tea.Cmd { return nil }

func (p *Model) Update(msg tea.Msg) (Panel, tea.Cmd) {
	updated := p

	return updated, nil
}

func (p *Model) IsFocused() bool {
	return p.isFocused
}

func (p *Model) SetFocused(f bool) {
	p.isFocused = f
}

func (p Model) Render() string {
	switch content := p.content.(type) {
	case string:
		return content
	case interface{ Render() string }:
		return content.Render()
	}
	return fmt.Sprintf("%s", p.content)
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

type ManagerPanel struct {
	*Manager
	id        string
	name      string
	isFocused bool

	Definition Definition
}

// ID implements Panel.
func (p *ManagerPanel) ID() string {
	return p.id
}

// Name implements Panel.
func (p *ManagerPanel) Name() string {
	return p.name
}

func NewManagerPanel(manager *Manager, opts ...PanelOption) *ManagerPanel {
	return &ManagerPanel{
		Manager: manager,
	}
}

func (p *ManagerPanel) Render() string {
	return p.Manager.Render()
}

func (p *ManagerPanel) View() tea.View {
	return tea.NewView(p.Render())
}

func (p *ManagerPanel) Init() tea.Cmd {
	return nil
}

func (p *ManagerPanel) Update(msg tea.Msg) (Panel, tea.Cmd) {
	var cmd tea.Cmd
	updated := p

	updated.Manager, cmd = updated.Manager.Update(msg)

	return updated, cmd
}

// IsFocused implements Panel.
func (p *ManagerPanel) IsFocused() bool {
	return p.isFocused
}

func (p *ManagerPanel) SetFocused(f bool) {
	p.isFocused = f
}
