package ui

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

type (
	ListItem string

	List struct {
		list.Model
		delegate list.ItemDelegate
	}
)

func (item ListItem) FilterValue() string { return "" }
func (item ListItem) Title() string       { return string(item) }
func (item ListItem) Description() string { return "" }

func NewList() *List {
	delegate := list.NewDefaultDelegate()
	delegate.Styles.NormalTitle = lipgloss.NewStyle().
		Foreground(lipgloss.NoColor{}).
		Faint(true).
		PaddingLeft(1)

	delegate.Styles.SelectedTitle = lipgloss.NewStyle().
		BorderForeground(lipgloss.NoColor{}).
		Foreground(lipgloss.NoColor{}).
		MarginLeft(1).
		Reverse(true)

	delegate.SetSpacing(0)
	delegate.ShowDescription = false

	list := list.New(make([]list.Item, 0), delegate, 0, 0)
	list.SetShowHelp(false)
	list.SetShowStatusBar(false)
	list.SetShowFilter(false)
	list.SetFilteringEnabled(false)
	list.SetShowPagination(false)
	list.SetShowTitle(false)

	return &List{
		Model:    list,
		delegate: delegate,
	}
}

func (l *List) SetStringItems(items []string) {
	listItems := make([]list.Item, len(items))
	for i, item := range items {
		listItems[i] = ListItem(item)
	}
	l.Model.ResetSelected()
	l.Model.SetItems(listItems)
}

func (l *List) SetConnectionItems(items []ConnectionItem) {
	listItems := make([]list.Item, len(items))
	for i, item := range items {
		listItems[i] = item
	}
	l.Model.ResetSelected()
	l.Model.SetItems(listItems)
}
