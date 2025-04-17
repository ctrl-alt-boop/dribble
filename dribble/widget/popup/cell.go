package popup

import tea "github.com/charmbracelet/bubbletea"

type TableCell struct {
	Value string
}

// SetSize implements PopupModel.
func (t *TableCell) SetSize(width int, height int) {
	panic("unimplemented")
}

// Exec implements PopupModel.
func (t *TableCell) Exec() tea.Cmd {
	panic("unimplemented")
}

// Cancel implements PopupModel.
func (t *TableCell) Cancel() tea.Cmd {
	panic("unimplemented")
}

// Init implements tea.Model.
func (t *TableCell) Init() tea.Cmd {
	panic("unimplemented")
}

// Update implements tea.Model.
func (t *TableCell) Update(tea.Msg) (tea.Model, tea.Cmd) {
	return t, nil
}

// View implements tea.Model.
func (t *TableCell) View() string {
	panic("unimplemented")
}
