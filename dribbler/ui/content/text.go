package content

import (
	tea "github.com/charmbracelet/bubbletea"
)

var _ tea.Model = (*Text)(nil)

type Text string

func (t Text) Init() tea.Cmd {
	return nil
}

func (t Text) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return t, nil
}

func (t Text) View() string {
	return string(t)
}
