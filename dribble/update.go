package dribble

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ctrl-alt-boop/gooldb/dribble/config"
	"github.com/ctrl-alt-boop/gooldb/dribble/io"
	"github.com/ctrl-alt-boop/gooldb/dribble/widget"
	"github.com/ctrl-alt-boop/gooldb/dribble/widget/popup"
	"github.com/ctrl-alt-boop/gooldb/internal/app/gooldb"
	"github.com/ctrl-alt-boop/gooldb/pkg/connection"
)

func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	connectMsg, ok := msg.(io.ConnectMsg)
	if ok {
		m.popupHandler.Close()
		return m, tea.Batch(
			widget.ChangeFocus(widget.KindPanel),
			m.Connect(connectMsg),
		)
	}

	if m.popupHandler.IsOpen() {
		switch msg := msg.(type) {
		case widget.PopupConfirmMsg:
			m.popupHandler.Close()
			return m, m.popupConfirm(msg)
		case widget.PopupCancelMsg:
			m.popupHandler.Close()
			return m, widget.ChangeFocus(widget.KindPanel)
		default:
			_, cmd = m.popupHandler.Update(msg)
			cmds = append(cmds, cmd)
			return m, tea.Batch(cmds...)
		}
	}

	// AppModel messages
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
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
		return m, tea.Batch(cmds...)

	case widget.RequestFocus:
		m.inFocus = widget.Kind(msg)
		return m, nil

	case io.GoolDbEventMsg:
		logger.Infof("GoolDbEvent received: %+v", msg)
		switch msg.Type {
		case gooldb.DriverLoadError:

		case gooldb.ConnectError:
		case gooldb.Connected:
			cmd = func() tea.Msg {
				m.gooldb.FetchDatabaseList()
				return nil
			}
			return m, cmd

		case gooldb.DatabaseListFetchError:
		case gooldb.DatabaseListFetched:

		case gooldb.DisconnectError:
		}

	case widget.SelectServerMsg:
		return m, m.SelectServer(msg)

	case widget.SelectDatabaseMsg:
		return m, m.SelectDatabase(msg)

	case widget.SelectTableMsg:
		return m, m.SelectTable(msg)

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
				m.inFocus = widget.KindWorkspace
			case widget.KindWorkspace:
				m.inFocus = widget.KindPanel
			}
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
	case io.GoolDbEventMsg:
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

func (m AppModel) onPanelSelect(msg widget.PanelSelectMsg) []tea.Cmd {
	var cmds []tea.Cmd

	switch msg.CurrentMode {
	case widget.ServerList:
		if config, ok := config.SavedConfigs[msg.Selected]; ok {
			cmds = append(cmds, func() tea.Msg { return io.ConnectMsg{Settings: config} })
			return cmds
		}
		cmds = append(cmds, m.popupHandler.Popup(popup.KindConnect, msg.Selected))
		cmds = append(cmds, widget.ChangeFocus(widget.KindPopupHandler))
	case widget.DatabaseList:
		cmds = append(cmds, func() tea.Msg { return widget.SelectDatabaseMsg(msg.Selected) })
	case widget.TableList:
		cmds = append(cmds, func() tea.Msg { return widget.SelectTableMsg(msg.Selected) })
	}

	return cmds
}

func (m AppModel) popupConfirm(msg widget.PopupConfirmMsg) tea.Cmd {
	settings := connection.NewSettings(
		connection.WithDriver(msg.DriverName),
		connection.WithHost(msg.Ip, msg.Port),
		connection.WithUser(msg.Username),
		connection.WithPassword(msg.Password),
	)
	return func() tea.Msg {
		m.gooldb.Connect(settings)
		return nil
	}
}

func (m AppModel) Connect(msg io.ConnectMsg) tea.Cmd {
	return func() tea.Msg {
		m.gooldb.Connect(msg.Settings)
		return nil
	}
}

func (m AppModel) SelectServer(msg widget.SelectServerMsg) tea.Cmd {
	logger.Infof("Server selected: %s", string(msg))
	saved, ok := config.SavedConfigs[string(msg)]
	if !ok {
		logger.Infof("Config not saved: %s", string(msg))
		var cmds []tea.Cmd
		cmds = append(cmds, m.popupHandler.Popup(popup.KindConnect, string(msg)))
		cmds = append(cmds, widget.ChangeFocus(widget.KindPopupHandler))
		return tea.Batch(cmds...)
	}
	logger.Infof("Config found: %+v", saved)
	return func() tea.Msg {
		m.gooldb.Connect(saved)
		return nil
	}
}

func (m AppModel) SelectDatabase(msg widget.SelectDatabaseMsg) tea.Cmd {
	return func() tea.Msg {
		m.gooldb.FetchTableList(string(msg))
		return nil
	}
}

func (m AppModel) SelectTable(msg widget.SelectTableMsg) tea.Cmd {
	return func() tea.Msg {
		m.gooldb.FetchTable(string(msg))
		return nil
	}
}
