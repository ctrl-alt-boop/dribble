package dribble

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss/list"
	"github.com/charmbracelet/lipgloss/tree"
	"github.com/ctrl-alt-boop/gooldb/dribble/config"
	"github.com/ctrl-alt-boop/gooldb/dribble/io"
	"github.com/ctrl-alt-boop/gooldb/dribble/ui"
	"github.com/ctrl-alt-boop/gooldb/dribble/widget"
	"github.com/ctrl-alt-boop/gooldb/dribble/widget/popup"
	"github.com/ctrl-alt-boop/gooldb/internal/app/gooldb"
	"github.com/ctrl-alt-boop/gooldb/pkg/connection"
	"github.com/ctrl-alt-boop/gooldb/pkg/logging"
)

var logger = logging.NewLogger("tea.log")

type testData struct {
	x int
	y int
}

type AppModel struct {
	gooldb        *gooldb.GoolDb
	Width, Height int

	panel        *widget.Panel
	prompt       *widget.Prompt
	workspace    *widget.Workspace
	help         *widget.Help
	popupHandler *popup.PopupHandler

	inFocus   widget.Kind
	prevFocus widget.Kind

	programSend func(msg tea.Msg)

	testData *testData
}

func NewModel(gool *gooldb.GoolDb) AppModel {
	config.LoadConfig()

	// testTree()

	return AppModel{
		gooldb:       gool,
		panel:        widget.NewPanel(gool),
		prompt:       widget.NewPromptBar(gool),
		workspace:    widget.NewWorkspace(gool),
		help:         widget.NewHelp(),
		popupHandler: popup.NewHandler(gool),

		testData: &testData{},
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

	cmds = append(cmds, widget.RequestFocusChange(widget.KindPanel))

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
