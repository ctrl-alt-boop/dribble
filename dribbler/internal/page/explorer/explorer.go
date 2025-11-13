package explorer

import (
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"
	"github.com/ctrl-alt-boop/dribbler/internal/dribbleapi"
	"github.com/ctrl-alt-boop/dribbler/internal/page"
	"github.com/ctrl-alt-boop/dribbler/logging"
)

const ExplorerPageID page.ID = "explorer"

type focusBackwardsMsg struct{}

type focusChangedMsg string

func focusBackwards() tea.Msg {
	return focusBackwardsMsg{}
}

type ExplorerPage struct {
	sidebar     *sidebar
	workspace   *workspace
	commandline *commandline // should it be hidable?
	help        *cheatsheet

	focusCycle []string
	focusIndex int
	focused    string

	width, height int

	normalStyle, focusedStyle lipgloss.Style

	panels map[string]*rect

	showHelp bool // May change to full help instead
	keybinds *ExplorerPageKeys
}

func NewExplorerPage() *ExplorerPage {
	return &ExplorerPage{
		sidebar:     newSidebar(),
		workspace:   newWorkspace(),
		commandline: newCommandline(),
		help:        newCheatsheet(),

		normalStyle: baseStyle.
			Border(lipgloss.NormalBorder()).
			Faint(true),

		focusedStyle: baseStyle.
			Border(lipgloss.ThickBorder()),

		focusCycle: []string{sidebarID, workspaceID},
		focusIndex: 0,
		focused:    sidebarID,

		showHelp: false,
		keybinds: DefaultExplorerKeyBindings(),
	}
}

// SetSize implements page.Page.
func (e *ExplorerPage) SetSize(width, height int) {
	logging.GlobalLogger().Infof("ExplorerPage.SetSize: %d, %d", width, height)
	e.width, e.height = width, height
	e.recompose(width, height)
}

// Init implements page.Page.
func (e *ExplorerPage) Init() tea.Cmd {
	logging.GlobalLogger().Infof("ExplorerPage.Init")

	var cmd tea.Cmd
	var cmds []tea.Cmd

	cmd = e.sidebar.Init()
	cmds = append(cmds, cmd)

	cmd = e.workspace.Init()
	cmds = append(cmds, cmd)

	cmd = e.commandline.Init()
	cmds = append(cmds, cmd)

	cmd = e.help.Init()
	cmds = append(cmds, cmd)

	e.help.registerKeymap(string(ExplorerPageID), e.keybinds)
	e.help.registerKeymap(sidebarID, e.sidebar.keybinds)
	e.help.registerKeymap(workspaceID, e.workspace.keybinds)
	e.help.registerKeymap(commandlineID, e.commandline.keybinds)

	return tea.Batch(cmds...)
}

// Update implements page.Page.
func (e *ExplorerPage) Update(msg tea.Msg) (page.Page, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	// Handle message that can be received at any time, regardless of focus or popups.
	switch msg := msg.(type) {
	case tea.WindowSizeMsg: // Handle application window resize
		e.width, e.height = msg.Width, msg.Height
		e.recompose(msg.Width, msg.Height)
		return e, nil
	}

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, e.keybinds.FocusCycleNext):
			if e.focused == commandlineID {
				break
			}
			e.focusIndex++
			if e.focusIndex >= len(e.focusCycle) {
				e.focusIndex = 0
			}
			e.focused = e.focusCycle[e.focusIndex]
			cmd = e.notifyFocusChange()
			return e, cmd
		case key.Matches(msg, e.keybinds.FocusCyclePrev):
			if e.focused == commandlineID {
				break
			}
			e.focusIndex--
			if e.focusIndex < 0 {
				e.focusIndex = len(e.focusCycle) - 1
			}
			e.focused = e.focusCycle[e.focusIndex]
			cmd = e.notifyFocusChange()
			return e, cmd
		case key.Matches(msg, e.keybinds.CommandlineMode):
			if e.focused != commandlineID {
				e.focused = commandlineID
				cmd = e.notifyFocusChange()
				return e, cmd
			}
		case key.Matches(msg, e.keybinds.ToggleHelp):
			if e.focused != commandlineID {
				e.toggleHelp()
			}
		}
	case focusBackwardsMsg:
		if e.focused == commandlineID {
			e.focused = e.focusCycle[e.focusIndex]
			cmd = e.notifyFocusChange()
			return e, cmd
		}

	case dribbleapi.DribbleResponseMsg:
		cmd := e.onDribbleResponse(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
		return e, cmd
	}

	cmd = e.updatePanels(msg)
	cmds = append(cmds, cmd)

	return e, tea.Batch(cmds...)
}

// Render implements page.Page.
func (e *ExplorerPage) Render() *lipgloss.Canvas {
	canvas := lipgloss.NewCanvas()

	for id, rect := range e.panels {
		var layer *lipgloss.Layer
		switch id {
		case sidebarID:
			if e.focused == sidebarID {
				e.sidebar.SetStyle(e.focusedStyle)
			} else {
				e.sidebar.SetStyle(e.normalStyle)
			}
			layer = e.sidebar.Render()
		case workspaceID:
			if e.focused == workspaceID {
				e.workspace.SetStyle(unstyled.Inherit(e.focusedStyle))
			} else {
				e.workspace.SetStyle(unstyled.Inherit(e.normalStyle))
			}
			layer = e.workspace.Render()
		case commandlineID:
			if e.focused == commandlineID {
				e.commandline.SetStyle(e.focusedStyle)
			} else {
				e.commandline.SetStyle(e.normalStyle)
			}
			layer = e.commandline.Render()
		case helpPanelID:
			if !e.showHelp {
				continue
			}
			layer = e.help.Render()
		default:
			logging.GlobalLogger().Warnf("Unknown panel id: %s", id)
			continue
		}
		canvas.AddLayers(layer.
			ID(id).
			X(rect.x).Y(rect.y).
			Width(rect.width).Height(rect.height))
	}
	return canvas
}

func (e *ExplorerPage) toggleHelp() {
	e.showHelp = !e.showHelp
	e.recompose(e.width, e.height)
}

func (e *ExplorerPage) notifyFocusChange() tea.Cmd {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	newFocusMsg := focusChangedMsg(e.focused)
	logging.GlobalLogger().Infof("notifyFocusChange: %s", newFocusMsg)
	e.sidebar, cmd = e.sidebar.Update(newFocusMsg)
	if cmd != nil {
		cmds = append(cmds, cmd)
	}
	if cmd != nil {
		cmds = append(cmds, cmd)
	}

	e.workspace, cmd = e.workspace.Update(newFocusMsg)
	if cmd != nil {
		cmds = append(cmds, cmd)
	}

	e.commandline, cmd = e.commandline.Update(newFocusMsg)
	if cmd != nil {
		cmds = append(cmds, cmd)
	}

	e.help, cmd = e.help.Update(newFocusMsg)
	if cmd != nil {
		cmds = append(cmds, cmd)
	}
	logging.GlobalLogger().Infof("notifyFocusChange: cmds: %s", cmds)
	return tea.Batch(cmds...)
}

func (e *ExplorerPage) updatePanels(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch e.focused {
	case sidebarID:
		e.sidebar, cmd = e.sidebar.Update(msg)
		cmds = append(cmds, cmd)

	case workspaceID:
		e.workspace, cmd = e.workspace.Update(msg)
		cmds = append(cmds, cmd)

	case commandlineID:
		e.commandline, cmd = e.commandline.Update(msg)
		cmds = append(cmds, cmd)
	}
	e.help, cmd = e.help.Update(msg)
	cmds = append(cmds, cmd)

	return tea.Batch(cmds...)
}

func (e *ExplorerPage) onDribbleResponse(msg dribbleapi.DribbleResponseMsg) tea.Cmd {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	e.sidebar, cmd = e.sidebar.Update(msg)
	if cmd != nil {
		cmds = append(cmds, cmd)
	}
	e.workspace, cmd = e.workspace.Update(msg)
	if cmd != nil {
		cmds = append(cmds, cmd)
	}

	return tea.Batch(cmds...)
}

func (e *ExplorerPage) recompose(width, height int) {
	helpHeight := defaultHelpHeight
	if !e.showHelp {
		helpHeight = 0
	}
	helpRect := newRect(0, height-helpHeight, width, helpHeight)
	e.help.SetWidth(helpRect.width)

	commandlineHeight := defaultCommandlineHeight + e.normalStyle.GetVerticalBorderSize()
	commandlineRect := newRect(0, height-commandlineHeight-helpHeight, width, commandlineHeight)
	e.commandline.SetWidth(commandlineRect.width)

	sidebarWidth := max(int(float64(width)*defaultSidebarWidthRatio), defaultSidebarMinWidth)
	sidebarRect := newRect(0, 0, sidebarWidth, height-commandlineHeight-helpHeight)

	e.sidebar.SetSize(
		sidebarRect.width,
		sidebarRect.height,
	)

	workspaceRect := newRect(sidebarWidth, 0, width-sidebarWidth, height-commandlineHeight-helpHeight)
	e.workspace.SetSize(
		workspaceRect.width,
		workspaceRect.height,
	)

	e.panels = map[string]*rect{
		sidebarID:     sidebarRect,
		workspaceID:   workspaceRect,
		commandlineID: commandlineRect,
		helpPanelID:   helpRect,
	}
	for id, rect := range e.panels {
		logging.GlobalLogger().Infof("%s: %+v", id, rect)
	}
}
