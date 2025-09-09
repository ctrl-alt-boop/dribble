package widget

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/ctrl-alt-boop/dribbler/config"
	"github.com/ctrl-alt-boop/dribbler/logging"
)

var logger = logging.GlobalLogger()

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

func RequestFocusChange(id Kind) tea.Cmd {
	return func() tea.Msg {
		return RequestFocus(id)
	}
}

func fullHelpPanelFunc() [][]key.Binding {
	return fullHelpPanel
}

func fullHelpWorkspaceFunc() [][]key.Binding {
	return fullHelpWorkspace
}

type WidgetDimensions map[Kind]struct {
	Width  int
	Height int
}

func GetWidgetDimensions(termWidth, termHeight int) WidgetDimensions {
	helpWidth := termWidth
	helpHeight := 1

	promptWidth := termWidth
	promptHeight := 1 + 1 // +1 for bottom border

	separatorHeight := 1

	footerHeight := helpHeight + promptHeight + separatorHeight

	panelWidth := 35
	panelHeight := termHeight - footerHeight

	workspaceWidth := termWidth - panelWidth
	workspaceHeight := termHeight - footerHeight

	popupWidth, popupHeight := workspaceWidth-10, workspaceHeight-10

	return WidgetDimensions{
		KindPanel:        {Width: panelWidth, Height: panelHeight},
		KindWorkspace:    {Width: workspaceWidth, Height: workspaceHeight},
		KindHelp:         {Width: helpWidth, Height: helpHeight},
		KindPrompt:       {Width: promptWidth, Height: promptHeight},
		KindPopupHandler: {Width: popupWidth, Height: popupHeight},
	}
}

var (
	fullHelpPanel [][]key.Binding = [][]key.Binding{
		{config.Keys.Help},
		{config.Keys.Quit},
		{config.Keys.CycleView},
		{config.Keys.Details},
		{config.Keys.Nav},
		{config.Keys.Select},
		{config.Keys.Back},
		{config.Keys.New},
	}
	fullHelpWorkspace [][]key.Binding = [][]key.Binding{
		{config.Keys.Help},
		{config.Keys.Quit},
		{config.Keys.CycleView},
		{config.Keys.Nav},
		{config.Keys.Zoom},
		{config.Keys.Select},
		{config.Keys.New},
	}
)

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
		Connect      string
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
		Connect:      "Popup_Connect",
		QueryOptions: "Popup_QueryOptions",
		TableCell:    "Popup_TableCell",
	},
}
