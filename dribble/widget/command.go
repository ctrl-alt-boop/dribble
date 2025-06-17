package widget

import tea "github.com/charmbracelet/bubbletea"

type (
	PopupConfirmMsg struct {
		DriverName string

		DefaultServer bool
		Ip            string
		Port          int

		Username string
		Password string
	}
	PopupCancelMsg struct{}

	SelectServerMsg   string
	SelectDatabaseMsg string
	SelectTableMsg    string

	WorkspaceSetMsg struct{}

	CellDataMsg struct {
		Value string
	}
)

func WorkspaceSet() tea.Msg {
	return WorkspaceSetMsg{}
}
