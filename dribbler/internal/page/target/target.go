package target

import (
	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"
	"github.com/ctrl-alt-boop/dribbler/internal/page"
)

type TargetingPage struct{}

func NewTargetingPage() *TargetingPage {
	return &TargetingPage{}
}

func (t *TargetingPage) SetSize(width, height int) {}

func (t *TargetingPage) Init() tea.Cmd {
	return nil
}

func (t *TargetingPage) Update(msg tea.Msg) (page.Page, tea.Cmd) {
	return t, nil
}

func (t *TargetingPage) Render() *lipgloss.Canvas {
	return lipgloss.NewCanvas()
}
