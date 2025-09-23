package popup

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/ctrl-alt-boop/dribble/database"
	"github.com/ctrl-alt-boop/dribble/request"
	"github.com/ctrl-alt-boop/dribble/target"
	"github.com/ctrl-alt-boop/dribbler/config"
	"github.com/ctrl-alt-boop/dribbler/ui"
	"github.com/ctrl-alt-boop/dribbler/widget"
)

type IntentBuilder struct {
	Intent *request.Intent
	Target target.Target

	QueryForm *huh.Form
	CancelCmd tea.Cmd

	stableWidth, stableHeight int
}

func newEmptyQueryBuilder() *IntentBuilder {
	return &IntentBuilder{
		CancelCmd: func() tea.Msg { return widget.PopupCancelMsg{} },
	}
}

func newTableQueryBuilder(method database.RequestType, tableName string) *IntentBuilder {
	q := newEmptyQueryBuilder()
	q.Intent = &request.Intent{
		Type: method,
	}
	return q
}

func newQueryBuilder(dbType *database.DBType) *IntentBuilder {
	q := newEmptyQueryBuilder()
	// q.Query = query
	return q
}

// Init implements tea.Model.
func (q *IntentBuilder) Init() tea.Cmd {
	if q.Intent == nil {
		q.Intent = &request.Intent{}
	}
	formTitle := fmt.Sprintf("New %s query", q.Target.Name)
	q.QueryForm = huh.NewForm(
		huh.NewGroup(
			huh.NewNote().Title(formTitle),
		),
		huh.NewGroup(
			huh.NewSelect[database.RequestType]().Title("Method:").
				Options(huh.NewOptions(database.RequestTypes...)...).
				Value(&q.Intent.Type),
		),
		huh.NewGroup(
			huh.NewInput().Title("Table:").
				Value(&q.Target.Name),
		),
	)

	q.stableWidth = lipgloss.Width(q.QueryForm.View())
	q.stableHeight = lipgloss.Height(q.QueryForm.View())
	return func() tea.Msg {
		return widget.IntentBuilderInitMsg{}
	}
}

// Update implements tea.Model.
func (q *IntentBuilder) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
	case widget.IntentBuilderDataMsg:
		// so this can add things like table schema etc.
	}

	// keyMsg, ok := msg.(tea.KeyMsg)
	// if ok && key.Matches(keyMsg, config.Keys.Back) {
	// 	return q, q.CancelCmd
	// }

	return q, cmd
}

func (q *IntentBuilder) ConfirmCmd() tea.Msg {

	return widget.IntentBuilderConfirmMsg{
		Intent: q.Intent,
	}
}

// GetContentSize implements PopupModel.
func (q *IntentBuilder) GetContentSize() (int, int) {
	return q.GetContentWidth(), q.GetContentHeight()
}

// GetContentWidth implements PopupModel.
func (q *IntentBuilder) GetContentWidth() int {
	return lipgloss.Width(q.QueryForm.View())
}

// GetContentHeight implements PopupModel.
func (q *IntentBuilder) GetContentHeight() int {
	return lipgloss.Height(q.QueryForm.View())
}

// SetMaxSize implements PopupModel.
func (q *IntentBuilder) SetMaxSize(width int, height int) {
	w := min(width-ui.PopupStyle.GetHorizontalFrameSize(), q.GetContentWidth())
	h := min(height-ui.PopupStyle.GetVerticalFrameSize(), q.GetContentHeight())

	q.QueryForm = q.QueryForm.WithWidth(w).WithHeight(h)
}

// View implements tea.Model.
func (q *IntentBuilder) View() string {
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
