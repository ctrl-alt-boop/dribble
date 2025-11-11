package components

import (
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
)

type Prompt struct {
	input textinput.Model
}

func NewPromptBar() Prompt {
	return Prompt{
		input: textinput.New(),
	}
}

func (c Prompt) Init() tea.Cmd {
	// c.input.Cursor.SetChar(">")
	c.input.CharLimit = 128
	return nil
}

func (c Prompt) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	updated := c
	var cmd tea.Cmd
	updated.input, cmd = c.input.Update(msg)

	return updated, cmd
}

func (c Prompt) Render() string {
	return c.input.View()
}

func (c Prompt) View() tea.View {
	return tea.NewView(c.Render())
}
