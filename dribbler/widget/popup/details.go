package popup

import (
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ctrl-alt-boop/dribbler/widget"
)

var _ PopupModel = (*Details)(nil)

type Details struct {
	viewport.Model

	Value     string
	CancelCmd tea.Cmd
}

func newDetails(values ...string) Details {
	v := strings.Join(values, "\n")
	return Details{
		Value:     v,
		CancelCmd: func() tea.Msg { return widget.PopupCancelMsg{} },
	}
}

// GetContentHeight implements PopupModel.
func (d Details) GetContentHeight() int {
	panic("unimplemented")
}

// GetContentSize implements PopupModel.
func (d Details) GetContentSize() (int, int) {
	panic("unimplemented")
}

// GetContentWidth implements PopupModel.
func (d Details) GetContentWidth() int {
	panic("unimplemented")
}

// Init implements PopupModel.
// Subtle: this method shadows the method (Model).Init of Details.Model.
func (d Details) Init() tea.Cmd {
	panic("unimplemented")
}

// SetMaxSize implements PopupModel.
func (d Details) SetMaxSize(width int, height int) {
	panic("unimplemented")
}

// Update implements PopupModel.
// Subtle: this method shadows the method (Model).Update of Details.Model.
func (d Details) Update(tea.Msg) (tea.Model, tea.Cmd) {
	panic("unimplemented")
}

// View implements PopupModel.
// Subtle: this method shadows the method (Model).View of Details.Model.
func (d Details) View() string {
	panic("unimplemented")
}
