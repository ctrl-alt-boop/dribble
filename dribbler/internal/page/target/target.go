package target

import (
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/ctrl-alt-boop/dribbler/internal/page"
)

type TargetingPage struct{}

func NewTargetingPage() *TargetingPage {
	return &TargetingPage{}
}

func (t *TargetingPage) Init() tea.Cmd {
	return nil
}

func (t *TargetingPage) Update(msg tea.Msg) (page.Page, tea.Cmd) {
	return t, nil
}

func (t *TargetingPage) Render() *lipgloss.Canvas {
	return lipgloss.NewCanvas()
}
