package panel

import (
	"github.com/charmbracelet/bubbles/v2/key"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/ctrl-alt-boop/dribbler/config"

	"github.com/ctrl-alt-boop/dribbler/util"
)

type Manager struct {
	composer    Composer
	composition *Composition

	Panels []Panel

	Width, Height int

	focusRing *util.Ring[int]

	normalStyle  lipgloss.Style
	focusedStyle lipgloss.Style
	panelBorder  lipgloss.Border
}

func NewPanelManager(composer Composer, panels ...Panel) *Manager {
	composer.SetNumPanels(len(panels))
	return &Manager{
		composer:     composer,
		Panels:       panels,
		focusRing:    util.NewRing(0, 1, 2), // FIXME: Fix focus ring
		normalStyle:  composer.GetStyle().Border(composer.GetPanelBorder()),
		focusedStyle: composer.GetFocusedStyle().Border(composer.GetPanelBorder()),
		panelBorder:  composer.GetPanelBorder(),
	}
}

func (p Manager) GetFocused() int {
	return p.focusRing.Value()
}

func (p Manager) Init() tea.Cmd {
	return nil
}

func (p *Manager) Update(msg tea.Msg) (*Manager, tea.Cmd) {
	var cmd tea.Cmd
	updated := p
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		updated.Layout(msg.Width, msg.Height)
		return updated, nil

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, config.Keys.CycleViewPrev):
			updated.focusRing.Backward()
			return updated, nil
		case key.Matches(msg, config.Keys.CycleViewNext):
			updated.focusRing.Forward()
			return updated, nil
		}
	}
	for i, panel := range p.Panels {
		updated.Panels[i], cmd = panel.Update(msg)
		if cmd != nil {
			return updated, cmd
		}
	}

	return updated, cmd
}

func (p Manager) Render() string {
	if p.composition == nil || len(p.composition.Layers) == 0 {
		return ""
	}
	render := lipgloss.NewCanvas()

	for i, panl := range p.Panels {
		style := p.normalStyle
		z := 0
		if i == p.focusRing.Value() {
			style = p.focusedStyle
			z = 1
		}

		panelRender := style.
			Width(p.composition.Layers[i].GetWidth()).
			Height(p.composition.Layers[i].GetHeight()).
			Render(panl.Render())

		layer := p.composition.Layers[i].
			SetContent(panelRender).
			Z(z)

		render.AddLayers(layer)
	}

	render.AddLayers(p.composition.PreRendered...)

	return render.Render()
}

func (p Manager) View() tea.View {
	return tea.NewView(p.Render())
}

func (p *Manager) Layout(width, height int) {
	p.Width, p.Height = width, height
	p.composition = p.composer.Compose(width, height)
}

func (p *Manager) String() string {
	return p.Render()
}

func (p *Manager) adjustPosition(x, y int) (int, int) {
	x--
	y--
	if x < 0 {
		x = 0
	}
	if y < 0 {
		y = 0
	}
	return x, y
}
