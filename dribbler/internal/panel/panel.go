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
}

type Model interface {
	ID() string
	Name() string
	Init() tea.Cmd
	Update(msg tea.Msg) (Model, tea.Cmd)
	Render() *lipgloss.Layer
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

func (m model) Render() *lipgloss.Layer {
	switch content := m.content.(type) {
	case string:
		return lipgloss.NewLayer(content)
	case interface{ Render() string }:
		return lipgloss.NewLayer(content.Render())
	case Model:
		return content.Render()
	}
	return lipgloss.NewLayer(fmt.Sprintf("%s", m.content))
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

	borderStyle lipgloss.Border
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

func (p *ManagerPanel) Render() *lipgloss.Layer {
	return lipgloss.NewLayer(p.Manager.Render())
}

func (p *ManagerPanel) View() tea.View {
	return tea.NewView(p.Render())
}

func (p *ManagerPanel) Init() tea.Cmd {
	return nil
}

func (p *ManagerPanel) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	updated := p

	updated.Manager, cmd = updated.Manager.Update(msg)

	return updated, cmd
}

func (p *ManagerPanel) SetBorderStyle(panelBorder lipgloss.Border) {
	p.borderStyle = panelBorder
}

// IsFocused implements Panel.
func (p *ManagerPanel) IsFocused() bool {
	return p.isFocused
}

func (p *ManagerPanel) SetFocused(f bool) {
	p.isFocused = f
}
