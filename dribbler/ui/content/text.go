package content

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ctrl-alt-boop/dribbler/ui"
)

var _ tea.Model = (*Text)(nil)

type Text struct {
	Item          any
	Width, Height int
}

func (t Text) Init() tea.Cmd {
	return nil
}

func (t Text) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return t, nil
}

func (t *Text) UpdateSize(width int, height int) {
	t.Width, t.Height = width, height
}

func (t *Text) String() string {
	return fmt.Sprint(t.Item)
}

func (t Text) View() string {
	return ui.Inline(t.Width, t.Height, t.String())
}
