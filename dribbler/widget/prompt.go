package widget

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/ctrl-alt-boop/dribble/dribble/ui"
	"github.com/ctrl-alt-boop/dribble/internal/app/dribble"
)

type Prompt struct {
	width, height int
	dribbleClient *dribble.Client

	input *huh.Input
}

func NewPromptBar(dribbleClient *dribble.Client) *Prompt {
	return &Prompt{
		dribbleClient: dribbleClient,
		input:         huh.NewInput().Prompt(">"),
	}
}

func (c *Prompt) Init() tea.Cmd {
	return nil
}

func (c *Prompt) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var input tea.Model
	input, cmd = c.input.Update(msg)
	c.input = input.(*huh.Input)

	return c, cmd
}

func (c *Prompt) UpdateSize(width, height int) {
	c.width, c.height = width, height
}

func (c *Prompt) View() string {
	contentWidth := c.width - ui.PromptStyle.GetHorizontalFrameSize()
	contentHeight := c.height - ui.PromptStyle.GetVerticalFrameSize()

	return ui.PromptStyle.Width(contentWidth).Height(contentHeight).Render(c.input.View())
}
