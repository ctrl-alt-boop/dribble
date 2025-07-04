package popup

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/ctrl-alt-boop/dribble/dribble/config"
	"github.com/ctrl-alt-boop/dribble/dribble/ui"
	"github.com/ctrl-alt-boop/dribble/dribble/widget"
	"github.com/ctrl-alt-boop/dribble/playbook/database"
)

type QueryBuilder struct {
	Query *database.Statement

	QueryForm *huh.Form
	CancelCmd tea.Cmd

	stableWidth, stableHeight int
}

func newEmptyQueryBuilder() *QueryBuilder {
	return &QueryBuilder{
		CancelCmd: func() tea.Msg { return widget.PopupCancelMsg{} },
	}
}

func newTableQueryBuilder(method database.SqlMethod, tableName string) *QueryBuilder {
	q := newEmptyQueryBuilder()
	q.Query = &database.Statement{
		Method: method,
		Table:  tableName,
	}
	return q
}

func newQueryBuilder(query *database.Query) *QueryBuilder {
	q := newEmptyQueryBuilder()
	// q.Query = query
	return q
}

// Init implements tea.Model.
func (q *QueryBuilder) Init() tea.Cmd {
	if q.Query == nil {
		q.Query = &database.Statement{
			Table: "table",
		}
	}
	formTitle := fmt.Sprintf("New %s query", q.Query.Table)
	q.QueryForm = huh.NewForm(
		huh.NewGroup(
			huh.NewNote().Title(formTitle),
		),
		huh.NewGroup(
			huh.NewSelect[database.SqlMethod]().Title("Method:").
				Options(huh.NewOptions(database.SqlMethods...)...).
				Value(&q.Query.Method),
		),
		huh.NewGroup(
			huh.NewInput().Title("Table:").
				Value(&q.Query.Table),
		),
	)

	q.stableWidth = lipgloss.Width(q.QueryForm.View())
	q.stableHeight = lipgloss.Height(q.QueryForm.View())
	return func() tea.Msg {
		return widget.QueryBuilderInitMsg{}
	}
}

// Update implements tea.Model.
func (q *QueryBuilder) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	form, cmd := q.QueryForm.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		q.QueryForm = f
	}

	switch q.QueryForm.State {
	case huh.StateCompleted:
		return q, q.ConfirmCmd
	case huh.StateAborted:
		return q, q.CancelCmd
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, config.Keys.Back):
			return q, q.CancelCmd
		}
	case widget.QueryBuilderDataMsg:
		// so this can add things like table schema etc.
	}

	// keyMsg, ok := msg.(tea.KeyMsg)
	// if ok && key.Matches(keyMsg, config.Keys.Back) {
	// 	return q, q.CancelCmd
	// }

	return q, cmd
}

func (q *QueryBuilder) ConfirmCmd() tea.Msg {

	return widget.QueryBuilderConfirmMsg{
		Query: q.Query,
	}
}

// GetContentSize implements PopupModel.
func (q *QueryBuilder) GetContentSize() (int, int) {
	return q.GetContentWidth(), q.GetContentHeight()
}

// GetContentWidth implements PopupModel.
func (q *QueryBuilder) GetContentWidth() int {
	return lipgloss.Width(q.QueryForm.View())
}

// GetContentHeight implements PopupModel.
func (q *QueryBuilder) GetContentHeight() int {
	return lipgloss.Height(q.QueryForm.View())
}

// SetMaxSize implements PopupModel.
func (q *QueryBuilder) SetMaxSize(width int, height int) {
	w := min(width-ui.PopupStyle.GetHorizontalFrameSize(), q.GetContentWidth())
	h := min(height-ui.PopupStyle.GetVerticalFrameSize(), q.GetContentHeight())

	q.QueryForm = q.QueryForm.WithWidth(w).WithHeight(h)
}

// View implements tea.Model.
func (q *QueryBuilder) View() string {
	if q.QueryForm == nil || q.QueryForm.State != huh.StateNormal {
		return ""
	}

	errors := q.QueryForm.Errors()
	errorText := ""
	for _, err := range errors {
		errorText += err.Error() + "\n"
	}

	errorRender := lipgloss.NewStyle().
		Foreground(lipgloss.Color("9")).
		Width(q.GetContentWidth()).Height(2).
		Render(errorText)

	render := lipgloss.JoinVertical(lipgloss.Left, q.QueryForm.View(), errorRender)

	return render
}
