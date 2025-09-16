package dribbler

import (
	"context"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ctrl-alt-boop/dribble"
	"github.com/ctrl-alt-boop/dribble/database"
	"github.com/ctrl-alt-boop/dribbler/config"
	"github.com/ctrl-alt-boop/dribbler/io"
	"github.com/ctrl-alt-boop/dribbler/widget"
	"github.com/ctrl-alt-boop/dribbler/widget/popup"
)

func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	logger.Infof("msg type received: %T", msg)
	var cmd tea.Cmd
	var cmds []tea.Cmd

	// Handle messages that can be received at any time, regardless of focus or popups.
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width, m.Height = msg.Width, msg.Height
		m.updateDimensions(msg)
	case tea.KeyMsg:
		if key.Matches(msg, config.Keys.Quit) {
			return m, tea.Quit
		}
	}

	// If a popup is open, it's the primary message handler.
	if m.popupHandler.IsOpen() {
		return m.updatePopupOpened(msg)
	}

	// If no popup is open, handle main application logic.
	m, cmd = m.updatePopupClosed(msg)
	cmds = append(cmds, cmd)

	// Update components that are always active (like help).
	// Note: tea.KeyMsg is handled by the focused widget, but help also needs it to toggle.
	_, cmd = m.help.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

// updatePopupOpened handles all message processing when a popup is active.
func (m AppModel) updatePopupOpened(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	// Messages that close the popup
	case widget.ConnectPopupConfirmMsg:
		m.popupHandler.Close()
		m.ChangeFocus(m.prevFocus)
		return m, m.connectPopupConfirm(msg)
	case widget.PopupCancelMsg:
		m.popupHandler.Close()
		m.ChangeFocus(m.prevFocus)
		return m, nil

	// Delegate all other messages to the popup handler.
	default:
		var popupModel tea.Model
		var cmd tea.Cmd
		popupModel, cmd = m.popupHandler.Update(msg)
		m.popupHandler = popupModel.(*popup.PopupHandler)
		return m, cmd
	}
}

// updatePopupClosed handles all message processing when no popup is active.
func (m AppModel) updatePopupClosed(msg tea.Msg) (AppModel, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	// App-level messages & actions
	case widget.RequestFocus:
		m.ChangeFocus(widget.Kind(msg))
	case io.ConnectMsg:
		m.ChangeFocus(widget.KindPanel)
		cmd = m.Connect(msg)
	case widget.SelectServerMsg:
		cmd = m.SelectServer(msg)
	case widget.SelectDatabaseMsg:
		cmd = m.SelectDatabase(msg)
	case widget.SelectTableMsg:
		cmd = m.SelectTable(msg)
	case widget.SelectTableColumnsMsg:
		cmd = m.SelectTableColumns(msg)

	// Messages that open popups
	case widget.OpenCellDataMsg:
		m.popupHandler.Popup(popup.KindTableCell, msg.Value)
		m.ChangeFocus(widget.KindPopupHandler)
	case widget.OpenQueryBuilderMsg:
		m.popupHandler.Popup(popup.KindQueryBuilder, msg.Method, msg.Table)
		m.ChangeFocus(widget.KindPopupHandler)

	// Dribble client events (broadcast)
	case io.DribbleEventMsg:
		logger.Infof("DribbleEvent received: %+v", msg)
		// Propagate to interested children
		var panelCmd, workspaceCmd tea.Cmd
		_, panelCmd = m.panel.Update(msg)
		_, workspaceCmd = m.workspace.Update(msg)
		cmds = append(cmds, panelCmd, workspaceCmd)

		// App-level handling of the event
		switch msg.Type {
		case dribble.SuccessConnect:
			// cmd = func() tea.Msg { return nil }
			// cmds = append(cmds, cmd)
		case dribble.SuccessExecute:

		}

	// Key presses
	case tea.KeyMsg:
		if key.Matches(msg, config.Keys.CycleView) {
			switch m.inFocus {
			case widget.KindPanel:
				m.ChangeFocus(widget.KindWorkspace)
			case widget.KindWorkspace:
				m.ChangeFocus(widget.KindPanel)
			}
			m.help.FocusChanged(m.inFocus)
		} else {
			// Dispatch to the focused widget.
			switch m.inFocus {
			case widget.KindPanel:
				_, cmd = m.panel.Update(msg)
			case widget.KindWorkspace:
				_, cmd = m.workspace.Update(msg)
			case widget.KindPrompt:
				_, cmd = m.prompt.Update(msg)
			}
		}
	}

	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m AppModel) updateDimensions(msg tea.WindowSizeMsg) {
	widgetDimensions := widget.GetWidgetDimensions(msg.Width, msg.Height)

	panelDimensions := widgetDimensions[widget.KindPanel]
	m.panel.UpdateSize(panelDimensions.Width, panelDimensions.Height)

	workspaceDimensions := widgetDimensions[widget.KindWorkspace]
	m.workspace.UpdateSize(workspaceDimensions.Width, workspaceDimensions.Height)

	helpDimensions := widgetDimensions[widget.KindHelp]
	m.help.UpdateSize(helpDimensions.Width, helpDimensions.Height)

	promptDimensions := widgetDimensions[widget.KindPrompt]
	m.prompt.UpdateSize(promptDimensions.Width, promptDimensions.Height)

	// popupDimensions := widgetDimensions[widget.KindPopupHandler]
	popupDimensions := widgetDimensions[widget.KindWorkspace] // Temporarily use workspace dimensions
	m.popupHandler.UpdateSize(popupDimensions.Width, popupDimensions.Height)
}

func (m AppModel) onPanelSelect(msg widget.PanelSelectMsg) []tea.Cmd {
	var cmds []tea.Cmd

	switch msg.CurrentMode {
	case widget.ServerList:
		if config, ok := config.SavedConfigs[msg.Selected]; ok {
			cmds = append(cmds, func() tea.Msg { return io.ConnectMsg{Target: config} })
			return cmds
		}
		cmds = append(cmds, m.popupHandler.Popup(popup.KindConnect, msg.Selected))
		m.ChangeFocus(widget.KindPopupHandler)
	case widget.DatabaseList:
		cmds = append(cmds, func() tea.Msg { return widget.SelectDatabaseMsg(msg.Selected) })
	case widget.TableList:
		cmds = append(cmds, func() tea.Msg { return widget.SelectTableMsg(msg.Selected) })
	}

	return cmds
}

func (m *AppModel) ChangeFocus(widget widget.Kind) {
	m.prevFocus = m.inFocus
	m.inFocus = widget
}

func (m AppModel) connectPopupConfirm(msg widget.ConnectPopupConfirmMsg) tea.Cmd {
	connectMsg := io.ConnectMsg{Target: database.NewTarget("", database.DBDriver,
		database.WithDriver(msg.DriverName),
		database.WithHost(msg.Ip, msg.Port),
		database.WithUser(msg.Username),
		database.WithPassword(msg.Password),
	)}
	return func() tea.Msg { return connectMsg }
}

func (m AppModel) Connect(msg io.ConnectMsg) tea.Cmd {
	return func() tea.Msg {
		err := m.dribbleClient.OpenTarget(context.TODO(), msg.Target)
		if err != nil {
			logger.Error(err)
			return io.NewDribbleError(err)
		}
		return nil
	}
}

func (m AppModel) SelectServer(msg widget.SelectServerMsg) tea.Cmd {
	saved, ok := config.SavedConfigs[string(msg)]
	if !ok {
		var cmds []tea.Cmd
		cmds = append(cmds, m.popupHandler.Popup(popup.KindConnect, string(msg)))
		m.ChangeFocus(widget.KindPopupHandler)
		return tea.Batch(cmds...)
	}
	logger.Infof("Config found: %+v", saved)
	return func() tea.Msg {
		m.dribbleClient.OpenTarget(context.TODO(), saved)
		return nil
	}
}

func (m AppModel) SelectDatabase(msg widget.SelectDatabaseMsg) tea.Cmd {
	return func() tea.Msg {
		// m.dribble.FetchTableList(string(msg)) //FIXME
		return nil
	}
}

func (m AppModel) SelectTable(msg widget.SelectTableMsg) tea.Cmd {
	return func() tea.Msg {
		// m.dribble.FetchTable(string(msg)) //FIXME
		return nil
	}
}

func (m AppModel) SelectTableColumns(msg widget.SelectTableColumnsMsg) tea.Cmd {
	return func() tea.Msg {
		// m.dribble.FetchTableColumns(string(msg)) //FIXME
		return nil
	}
}
