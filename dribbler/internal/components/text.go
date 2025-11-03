package components

import (
	tea "github.com/charmbracelet/bubbletea/v2"
)

type TextItem struct {
	content string
}

func Text(content string) TextItem {
	return TextItem{
		content: content,
	}
}

func (t TextItem) Init() tea.Cmd {
	return nil
}

func (t TextItem) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return t, nil
}

func (t TextItem) Render() string {
	return t.content
}

func (t TextItem) View() tea.View {
	return tea.NewView(t.Render())
}
