package widget

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ctrl-alt-boop/dribble/database"
	"github.com/ctrl-alt-boop/dribble/request"
)

type (
	ConnectPopupConfirmMsg struct {
		DriverName string

		DefaultServer bool
		Ip            string
		Port          int

		Username string
		Password string
	}

	PopupCancelMsg struct{}

	SelectServerMsg       string
	SelectDatabaseMsg     string
	SelectTableMsg        string
	SelectTableColumnsMsg string

	WorkspaceSetMsg struct{}

	OpenCellDataMsg struct {
		Value string
	}

	OpenIntentBuilderMsg struct {
		Method database.RequestType
		Table  string
	}

	IntentBuilderInitMsg struct {
		Database string
		Method   database.RequestType
		Table    string
	}

	IntentBuilderDataMsg struct {
	}

	IntentBuilderConfirmMsg struct {
		Intent *request.Intent
	}
)

func WorkspaceSet() tea.Msg {
	return WorkspaceSetMsg{}
}
