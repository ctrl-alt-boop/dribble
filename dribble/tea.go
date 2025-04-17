package dribble

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ctrl-alt-boop/gooldb/dribble/io"
	"github.com/ctrl-alt-boop/gooldb/dribble/widget"
	"github.com/ctrl-alt-boop/gooldb/dribble/widget/popup"
	"github.com/ctrl-alt-boop/gooldb/internal/app/gooldb"
	"github.com/ctrl-alt-boop/gooldb/pkg/logging"
)

var logger = logging.NewLogger("tea.log")

type AppModel struct {
	gooldb        *gooldb.GoolDb
	Width, Height int

	panel        *widget.Panel
	prompt       *widget.Prompt
	workspace    *widget.Workspace
	help         *widget.Help
	popupHandler *popup.PopupHandler

	inFocus widget.Kind

	programSend func(msg tea.Msg)
}

func NewModel(gool *gooldb.GoolDb) AppModel {
	return AppModel{
		gooldb:       gool,
		panel:        widget.NewPanel(gool),
		prompt:       widget.NewPromptBar(gool),
		workspace:    widget.NewWorkspace(gool),
		help:         widget.NewHelp(),
		popupHandler: popup.NewHandler(gool),
	}
}

func (m AppModel) SetProgramSend(send func(msg tea.Msg)) {
	m.programSend = send

	m.gooldb.OnEvent(func(eventType gooldb.EventType, args any, err error) {
		event := io.GoolDbEventMsg{
			Type: eventType,
			Args: args,
			Err:  err,
		}
		m.programSend(event)
	})
}

func (m AppModel) Init() tea.Cmd {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	cmd = m.panel.Init()
	cmds = append(cmds, cmd)
	cmd = m.prompt.Init()
	cmds = append(cmds, cmd)
	cmd = m.workspace.Init()
	cmds = append(cmds, cmd)
	cmd = m.help.Init()
	cmds = append(cmds, cmd)
	cmd = m.popupHandler.Init()
	cmds = append(cmds, cmd)

	cmds = append(cmds, widget.ChangeFocus(widget.KindPanel))

	return tea.Batch(cmds...)
}

// func (m AppModel) onEventFunc(eventType gooldb.EventType) func(a any, err error) {
// 	return func(a any, err error) {
// 		event := io.GoolDbEventMsg{
// 			Type: eventType,
// 			Args: a,
// 			Err:  err,
// 		}
// 		m.programSend(event)
// 	}
// }
