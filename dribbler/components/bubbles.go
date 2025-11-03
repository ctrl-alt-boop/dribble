package components

import (
	"github.com/charmbracelet/bubbles/v2/list"
	"github.com/charmbracelet/bubbles/v2/table"
	tea "github.com/charmbracelet/bubbletea/v2"
)

type List struct {
	list list.Model
}

func (b List) Init() tea.Cmd {
	return nil
}

func (b List) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	updated := b
	updated.list, cmd = b.list.Update(msg)

	return updated, cmd
}

func (b List) Render() string {
	return b.list.View()
}

func (b List) View() tea.View {
	return tea.NewView(b.list.View())
}

type Table struct {
	table table.Model
}

func (b Table) Init() tea.Cmd {
	return nil
}

func (b Table) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	updated := b
	updated.table, cmd = b.table.Update(msg)

	return updated, cmd
}

func (b Table) Render() string {
	return b.table.View()
}

func (b Table) View() tea.View {
	return tea.NewView(b.table.View())
}
