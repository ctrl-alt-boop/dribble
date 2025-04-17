package widget

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/ctrl-alt-boop/gooldb/dribble/ui"
	"github.com/ctrl-alt-boop/gooldb/internal/app/gooldb"
)

type Prompt struct {
	width, height int
	goolDb        *gooldb.GoolDb
}

func NewPromptBar(gool *gooldb.GoolDb) *Prompt {
	return &Prompt{
		goolDb: gool,
	}
}

func (c *Prompt) Init() tea.Cmd {
	return nil
}

func (c *Prompt) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		c.UpdateSize(msg.Width, 1)
	}
	return c, nil
}

func (c *Prompt) UpdateSize(termWidth, termHeight int) {
	c.width, c.height = termWidth-ui.BorderThicknessDouble, termHeight
}

func (c *Prompt) View() string {
	promptBorder := lipgloss.Border{
		Bottom:      "─",
		Top:         "─",
		Left:        "│",
		Right:       "│",
		TopLeft:     "├",
		TopRight:    "┤",
		BottomLeft:  "└",
		BottomRight: "┘",
	}

	input := huh.NewInput().
		Prompt(">")

	inputStyle := lipgloss.NewStyle().
		Width(c.width).
		Height(c.height).
		Align(lipgloss.Left, lipgloss.Center).
		Border(promptBorder, false, true, true, true)

	return inputStyle.Render(input.View())
}
