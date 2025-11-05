package explorer

import (
	"github.com/charmbracelet/bubbles/v2/key"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/ctrl-alt-boop/dribbler/internal/core/util"
	"github.com/ctrl-alt-boop/dribbler/internal/page"
)

const ExplorerPageID page.ID = "explorer"

type focusBackwardsMsg struct{}

type focusChangedMsg page.ID

func focusBackwards() tea.Msg {
	return focusBackwardsMsg{}
}

const (
	defaultSidePanelWidthRatio = 0.2
	defaultSidePanelMinWidth   = 50

	defaultPromptNumLines = 1
	defaultHelpNumLines   = 1
)

type panelDefinition struct {
	x, y          int
	width, height int

	style        lipgloss.Style
	focusedStyle lipgloss.Style
}

type ExplorerPage struct {
	sidePanel *sidePanel
	workspace *workspacePanel
	prompt    *promptPanel
	help      *helpPanel

	width, height int

	focusRing *util.Ring[string]

	panelDefs map[string]panelDefinition
}

func NewExplorerPage(width, height int) *ExplorerPage {
	return &ExplorerPage{
		width:  width,
		height: height,

		sidePanel: newSidePanel(),
		workspace: newWorkspace(),
		prompt:    newPrompt(),
		help:      newHelpBar(),

		focusRing: util.NewRing(sidePanelID, workspacePanelID),
	}
}

func (e *ExplorerPage) Init() tea.Cmd {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	e.panelDefs = compose(e.width, e.height)

	sidePanelDef := e.panelDefs[sidePanelID]
	e.sidePanel.SetInnerSize(
		sidePanelDef.width-sidePanelDef.style.GetHorizontalFrameSize(),
		sidePanelDef.height-sidePanelDef.style.GetVerticalFrameSize(),
	)
	cmd = e.sidePanel.Init()
	cmds = append(cmds, cmd)

	workspacePanelDef := e.panelDefs[workspacePanelID]
	e.workspace.SetInnerSize(
		workspacePanelDef.width-workspacePanelDef.style.GetHorizontalFrameSize(),
		workspacePanelDef.height-workspacePanelDef.style.GetVerticalFrameSize(),
	)
	cmd = e.workspace.Init()
	cmds = append(cmds, cmd)

	promptPanelDef := e.panelDefs[promptPanelID]
	e.prompt.SetInnerWidth(promptPanelDef.width - promptPanelDef.style.GetHorizontalFrameSize())
	cmd = e.prompt.Init()
	cmds = append(cmds, cmd)

	helpPanelDef := e.panelDefs[helpPanelID]
	e.help.SetInnerWidth(helpPanelDef.width - helpPanelDef.style.GetHorizontalFrameSize())
	cmd = e.help.Init()
	cmds = append(cmds, cmd)

	return tea.Batch(cmds...)
}

func (e *ExplorerPage) Update(msg tea.Msg) (page.Page, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	// Handle message that can be received at any time, regardless of focus or popups.
	switch msg := msg.(type) {
	case tea.WindowSizeMsg: // Handle application window resize
		e.width, e.height = msg.Width, msg.Height
		e.setSize()
		return e, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keyMap.FocusCycleNext):
			e.focusRing.Forward()
			cmd = e.notifyFocusChange()
			cmds = append(cmds, cmd)
		case key.Matches(msg, keyMap.FocusFocusPrev):
			e.focusRing.Backward()
			cmd = e.notifyFocusChange()
			cmds = append(cmds, cmd)
		}
	}

	cmd = e.updatePanels(msg)
	cmds = append(cmds, cmd)

	return e, tea.Batch(cmds...)
}

func (e *ExplorerPage) Render() *lipgloss.Canvas {
	sidePanelDef := e.panelDefs[sidePanelID]
	sidePanelRender := e.sidePanel.Render().X(sidePanelDef.x).Y(sidePanelDef.y).Width(sidePanelDef.width).Height(sidePanelDef.height)

	workspacePanelDef := e.panelDefs[workspacePanelID]
	workspacePanelRender := e.workspace.Render().X(workspacePanelDef.x).Y(workspacePanelDef.y).Width(workspacePanelDef.width).Height(workspacePanelDef.height)

	promptPanelDef := e.panelDefs[promptPanelID]
	promptPanelRender := e.prompt.Render().X(promptPanelDef.x).Y(promptPanelDef.y).Width(promptPanelDef.width).Height(promptPanelDef.height)

	helpPanelDef := e.panelDefs[helpPanelID]
	helpPanelRender := e.help.Render().X(helpPanelDef.x).Y(helpPanelDef.y).Width(helpPanelDef.width).Height(helpPanelDef.height)

	return lipgloss.NewCanvas(sidePanelRender, workspacePanelRender, promptPanelRender, helpPanelRender)
}

func (e *ExplorerPage) notifyFocusChange() tea.Cmd {
	var cmds []tea.Cmd

	newFocusMsg := focusChangedMsg(e.focusRing.Value())
	_, cmd := e.sidePanel.Update(newFocusMsg)
	cmds = append(cmds, cmd)

	_, cmd = e.workspace.Update(newFocusMsg)
	cmds = append(cmds, cmd)

	_, cmd = e.prompt.Update(newFocusMsg)
	cmds = append(cmds, cmd)

	_, cmd = e.help.Update(newFocusMsg)
	cmds = append(cmds, cmd)

	return tea.Batch(cmds...)
}

func (e *ExplorerPage) updatePanels(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	e.sidePanel, cmd = e.sidePanel.Update(msg)
	cmds = append(cmds, cmd)

	e.workspace, cmd = e.workspace.Update(msg)
	cmds = append(cmds, cmd)

	e.prompt, cmd = e.prompt.Update(msg)
	cmds = append(cmds, cmd)

	e.help, cmd = e.help.Update(msg)
	cmds = append(cmds, cmd)

	return tea.Batch(cmds...)
}

func (e *ExplorerPage) setSize() {
	e.panelDefs = compose(e.width, e.height)

	sidePanelDef := e.panelDefs[sidePanelID]
	e.sidePanel.SetInnerSize(
		sidePanelDef.width-sidePanelDef.style.GetHorizontalFrameSize(),
		sidePanelDef.height-sidePanelDef.style.GetVerticalFrameSize(),
	)

	workspacePanelDef := e.panelDefs[workspacePanelID]
	e.workspace.SetInnerSize(
		workspacePanelDef.width-workspacePanelDef.style.GetHorizontalFrameSize(),
		workspacePanelDef.height-workspacePanelDef.style.GetVerticalFrameSize(),
	)

	promptPanelDef := e.panelDefs[promptPanelID]
	e.prompt.SetInnerWidth(promptPanelDef.width - promptPanelDef.style.GetHorizontalFrameSize())

	helpPanelDef := e.panelDefs[helpPanelID]
	e.help.SetInnerWidth(helpPanelDef.width - helpPanelDef.style.GetHorizontalFrameSize())
}

func compose(width, height int) map[string]panelDefinition {
	sidePanelWidth := max(int(float64(width)*defaultSidePanelWidthRatio), defaultSidePanelMinWidth)

	return map[string]panelDefinition{
		sidePanelID: {
			x:      0,
			y:      0,
			width:  int(float64(width) * defaultSidePanelWidthRatio),
			height: height - 3 - 3,
		},
		workspacePanelID: {
			x:      sidePanelWidth,
			y:      0,
			width:  width - sidePanelWidth,
			height: height - 3 - 3,
		},
		promptPanelID: {
			x:      0,
			y:      height - 3 - 3,
			width:  width,
			height: 3,
		},
		helpPanelID: {
			x:      0,
			y:      height - 3,
			width:  width,
			height: 3,
		},
	}
}
