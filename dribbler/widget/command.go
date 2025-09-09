package widget

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ctrl-alt-boop/dribble/database"
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

	OpenQueryBuilderMsg struct {
		Method database.QueryType
		Table  string
	}

	QueryBuilderInitMsg struct {
		Database string
		Method   string
		Table    string
	}

	QueryBuilderDataMsg struct {
	}

	QueryBuilderConfirmMsg struct {
		Query *database.QueryIntent
	}
)

func WorkspaceSet() tea.Msg {
	return WorkspaceSetMsg{}
}
