package dribbler

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss/list"
	"github.com/charmbracelet/lipgloss/tree"
	"github.com/ctrl-alt-boop/dribble"
	"github.com/ctrl-alt-boop/dribble/dribble/config"
	"github.com/ctrl-alt-boop/dribble/dribble/io"
	"github.com/ctrl-alt-boop/dribble/dribble/ui"
	"github.com/ctrl-alt-boop/dribble/dribble/widget"
	"github.com/ctrl-alt-boop/dribble/dribble/widget/popup"
	"github.com/ctrl-alt-boop/dribble/internal/app/dribble"
	"github.com/ctrl-alt-boop/dribble/playbook/connection"
	"github.com/ctrl-alt-boop/dribble/playbook/logging"
)

var logger = logging.GlobalLogger()

type AppModel struct {
	dribbleClient *dribble.Client
	Width, Height int

	panel        *widget.Panel
	prompt       *widget.Prompt
	workspace    *widget.Workspace
	help         *widget.Help
	popupHandler *popup.PopupHandler

	inFocus   widget.Kind
	prevFocus widget.Kind

	programSend func(msg tea.Msg)
}

func InitialModel(dribbleClient *dribble.DB) AppModel {
	config.LoadConfig()

	// testTree()

	return AppModel{
		dribbleClient: dribbleClient,
		panel:         widget.NewPanel(dribbleClient),
		prompt:        widget.NewPromptBar(dribbleClient),
		workspace:     widget.NewWorkspace(dribbleClient),
		help:          widget.NewHelp(),
		popupHandler:  popup.NewHandler(dribbleClient),
	}
}

func (m AppModel) SetProgramSend(send func(msg tea.Msg)) {
	m.programSend = send

	m.dribbleClient.OnEvent(func(eventType dribble.EventType, args any, err error) {
		event := io.DribbleEventMsg{
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

	cmds = append(cmds, widget.RequestFocusChange(widget.KindPanel))

	return tea.Batch(cmds...)
}

func testTree() {
	uiTree := ui.NewTree()

	var treeItems []*ui.TreeNode
	connItems := ui.GetSavedConfigsSorted()
	for _, item := range connItems {
		treeItem := ui.NewConnectionNode(connection.Server, item)
		treeItems = append(treeItems, treeItem)
	}

	categoryItem := ui.NewCategoryNode("Things", treeItems)

	uiTree.Child(categoryItem)

	logger.Infof("Tree:\n%+v", uiTree)

	logger.Infof("EnumeratorTest:\n%+v", list.New(uiTree.Children()))

	logger.Infof("Nested:\n%+v", ui.CreateNestedList())

	l := list.New()
	l.Item("aa")
	l.Item("ba")
	l.Item("bb")
	logger.Infof("List:\n%+v", l)

	treeTree := tree.New().Child("aa")
	treeTree.Child("ba")
	treeTree.Child("bb")
	logger.Infof("tree.Tree:\n%+v", treeTree)
}
