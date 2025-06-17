package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ctrl-alt-boop/gooldb/pkg/data"
	"github.com/ctrl-alt-boop/gooldb/pkg/logging"
)

var logger = logging.NewLogger("ui.log")

type ListTable struct {
	table []*List
}

func (t *ListTable) SetHeight(height int) {
	if t.table == nil {
		return
	}
	for _, table := range t.table {
		table.Model.SetHeight(height)
	}
}

func NewTable() *ListTable {
	return &ListTable{}
}

func (t *ListTable) SetTable(table *data.Table) {
	lists := make([]*List, table.NumColumns())
	for i := range lists {
		lists[i] = NewList()
	}
	for i := range table.Columns() {
		column, width := table.GetColumnRows(i)
		name := table.Columns()[i].Name
		if width < len(name) {
			width = len(name)
		}
		lists[i].Model.SetWidth(width + 4)
		lists[i].Model.SetHeight(len(column))
		lists[i].Model.Styles.Title = lipgloss.NewStyle().
			PaddingRight(width - len(name)).
			BorderStyle(lipgloss.NormalBorder()).
			BorderRight(true)
		lists[i].SetStringItems(column)
		lists[i].Model.SetShowTitle(true)
		lists[i].Title = name
	}
	t.table = lists
}

func (t *ListTable) View() string {
	if t.table == nil {
		return ""
	}
	views := make([]string, len(t.table))
	for i, table := range t.table[:2] {
		logger.Info(table.View())
		views[i] = table.View()
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, views...)
}

func (t *ListTable) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	return t, nil
}

func (t *ListTable) Init() tea.Cmd {
	return nil
}
