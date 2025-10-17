package component

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Prompt struct {
	// input *huh.input
	input textinput.Model
}

func NewPromptBar() Prompt {
	return Prompt{
		input: textinput.New(),
	}
}

func (c Prompt) Init() tea.Cmd {
	c.input.Cursor.SetChar(">")
	c.input.CharLimit = 128
	return nil
}

func (c Prompt) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	updated := c
	var cmd tea.Cmd
	updated.input, cmd = c.input.Update(msg)

	return updated, cmd
}

func (c Prompt) View() string {
	return c.input.View()
}
