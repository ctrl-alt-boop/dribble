package dribbler

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ctrl-alt-boop/dribble/dribble/config"
	"github.com/ctrl-alt-boop/dribble/dribble/io"
	"github.com/ctrl-alt-boop/dribble/dribble/widget"
	"github.com/ctrl-alt-boop/dribble/dribble/widget/popup"
	"github.com/ctrl-alt-boop/dribble/internal/app/dribble"
	"github.com/ctrl-alt-boop/dribble/playbook/connection"
)

func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	logger.Infof("msg type received: %T", msg)
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		logger.Infof("key msg received: %v", msg.String())
	case tea.WindowSizeMsg:
		m.Width, m.Height = msg.Width, msg.Height
		m.updateDimensions(msg)
	case io.ConnectMsg:
		m.popupHandler.Close()
		m.ChangeFocus(widget.KindPanel)
		return m, tea.Batch(
			m.Connect(msg),
		)
	}

	if m.popupHandler.IsOpen() {
		switch msg := msg.(type) {
		case widget.ConnectPopupConfirmMsg:
			m.popupHandler.Close()
			m.ChangeFocus(m.prevFocus)
			return m, tea.Batch(m.connectPopupConfirm(msg))
		case widget.PopupCancelMsg:
			m.popupHandler.Close()
			m.ChangeFocus(m.prevFocus)
			_, cmd = m.popupHandler.Update(msg)
			return m, cmd
		default:
			_, cmd = m.popupHandler.Update(msg)
			cmds = append(cmds, cmd)
			return m, tea.Batch(cmds...)
		}
	}

	// AppModel messages
	switch msg := msg.(type) {
	case widget.RequestFocus:
		m.ChangeFocus(widget.Kind(msg))
		return m, nil

	case io.DribbleEventMsg:
		logger.Infof("DribbleEvent received: %+v", msg)
		switch msg.Type {
		case dribble.DriverLoadError:

		case dribble.ConnectError:
		case dribble.Connected:
			cmd = func() tea.Msg {
				m.dribble.FetchDatabaseList()
				return nil
			}
			return m, cmd

		case dribble.DatabaseListFetchError:
		case dribble.DatabaseListFetched:

		case dribble.DisconnectError:
		}

	case widget.SelectServerMsg:
		return m, m.SelectServer(msg)

	case widget.SelectDatabaseMsg:
		return m, m.SelectDatabase(msg)

	case widget.SelectTableMsg:
		return m, m.SelectTable(msg)

	case widget.SelectTableColumnsMsg:
		return m, m.SelectTableColumns(msg)

	case widget.OpenCellDataMsg:
		m.popupHandler.Popup(popup.KindTableCell, msg.Value)
		m.ChangeFocus(widget.KindPopupHandler)
		return m, nil

	case widget.OpenQueryBuilderMsg:
		m.popupHandler.Popup(popup.KindQueryBuilder, msg.Method, msg.Table)
		m.ChangeFocus(widget.KindPopupHandler)
		return m, nil

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, config.Keys.Quit):
			return m, tea.Quit
		}
		// case message.CommandExecMsg:
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, config.Keys.CycleView):
			switch m.inFocus {
			case widget.KindPanel:
				m.ChangeFocus(widget.KindWorkspace)
			case widget.KindWorkspace:
				m.ChangeFocus(widget.KindPanel)
			}
			m.help.FocusChanged(m.inFocus)
			return m, tea.Batch(cmds...)
		case m.inFocus == widget.KindPanel:
			_, cmd = m.panel.Update(msg)
			cmds = append(cmds, cmd)
		case m.inFocus == widget.KindWorkspace:
			_, cmd = m.workspace.Update(msg)
			cmds = append(cmds, cmd)
		case m.inFocus == widget.KindPrompt:
			_, cmd = m.prompt.Update(msg)
			cmds = append(cmds, cmd)
		case m.inFocus == widget.KindPopupHandler:
			_, cmd = m.popupHandler.Update(msg)
			cmds = append(cmds, cmd)
		}
		_, cmd = m.help.Update(msg)
		cmds = append(cmds, cmd)
	case io.DribbleEventMsg:
		_, cmd = m.panel.Update(msg)
		cmds = append(cmds, cmd)
		_, cmd = m.workspace.Update(msg)
		cmds = append(cmds, cmd)
		_, cmd = m.prompt.Update(msg)
		cmds = append(cmds, cmd)
		_, cmd = m.help.Update(msg)
		cmds = append(cmds, cmd)
		_, cmd = m.popupHandler.Update(msg)
		cmds = append(cmds, cmd)
	}

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
			cmds = append(cmds, func() tea.Msg { return io.ConnectMsg{Settings: config} })
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
	settings := connection.NewSettings(
		connection.WithDriver(msg.DriverName),
		connection.WithHost(msg.Ip, msg.Port),
		connection.WithUser(msg.Username),
		connection.WithPassword(msg.Password),
	)

	return func() tea.Msg { return io.ConnectMsg{Settings: settings} }
}

func (m AppModel) Connect(msg io.ConnectMsg) tea.Cmd {
	return func() tea.Msg {
		m.dribble.Connect(msg.Settings)
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
		m.dribble.Connect(saved)
		return nil
	}
}

func (m AppModel) SelectDatabase(msg widget.SelectDatabaseMsg) tea.Cmd {
	return func() tea.Msg {
		m.dribble.FetchTableList(string(msg))
		return nil
	}
}

func (m AppModel) SelectTable(msg widget.SelectTableMsg) tea.Cmd {
	return func() tea.Msg {
		m.dribble.FetchTable(string(msg))
		return nil
	}
}

func (m AppModel) SelectTableColumns(msg widget.SelectTableColumnsMsg) tea.Cmd {
	return func() tea.Msg {
		m.dribble.FetchTableColumns(string(msg))
		return nil
	}
}
