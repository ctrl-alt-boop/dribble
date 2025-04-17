package widget

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ctrl-alt-boop/gooldb/pkg/logging"
)

var logger = logging.NewLogger("widgets.log")

type Kind int

const (
	KindPanel Kind = iota
	KindWorkspace
	KindHelp
	KindPrompt

	KindPopupHandler
	KindQueryOptions
	KindTableCell
)

type RequestFocus Kind

func ChangeFocus(id Kind) tea.Cmd {
	return func() tea.Msg {
		return RequestFocus(id)
	}
}

type (
	WidgetNames struct {
		Prompt    string
		Workspace string
		Panel     string
		Help      string

		Popups PopupNames
	}

	PopupNames struct {
		Handler      string
		QueryOptions string
		TableCell    string
	}
)

var Widgets = WidgetNames{
	Prompt:    "Prompt",
	Panel:     "Panel",
	Workspace: "Workspace",

	Help: "Help",

	Popups: PopupNames{
		Handler:      "Popups",
		QueryOptions: "Popup_QueryOptions",
		TableCell:    "Popup_TableCell",
	},
}
