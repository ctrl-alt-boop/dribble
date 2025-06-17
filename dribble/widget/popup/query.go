package popup

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ctrl-alt-boop/gooldb/pkg/query"
)

type QueryOptions struct {
	Query *query.Statement
}

func newQueryOptions(query *query.Statement) *QueryOptions {
	return &QueryOptions{
		Query: query,
	}
}

// GetContentHeight implements PopupModel.
func (q *QueryOptions) GetContentHeight() int {
	panic("unimplemented")
}

// GetContentWidth implements PopupModel.
func (q *QueryOptions) GetContentWidth() int {
	panic("unimplemented")
}

// Exec implements PopupModel.
func (q *QueryOptions) Exec() tea.Cmd {
	panic("unimplemented")
}

// Cancel implements PopupModel.
func (q *QueryOptions) Cancel() tea.Cmd {
	panic("unimplemented")
}

// Init implements tea.Model.
func (q *QueryOptions) Init() tea.Cmd {
	panic("unimplemented")
}

// Update implements tea.Model.
func (q *QueryOptions) Update(tea.Msg) (tea.Model, tea.Cmd) {
	return q, nil
}

// GetContentSize implements PopupModel.
func (q *QueryOptions) GetContentSize() (int, int) {
	panic("unimplemented")
}

// SetMaxSize implements PopupModel.
func (q *QueryOptions) SetMaxSize(width int, height int) {
	panic("unimplemented")
}

// View implements tea.Model.
func (q *QueryOptions) View() string {
	panic("unimplemented")
}
