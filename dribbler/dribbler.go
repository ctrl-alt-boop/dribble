// Package dribbler is the main model package
package dribbler

import (
	"context"
	"fmt"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/term"
	"github.com/ctrl-alt-boop/dribble"
	"github.com/ctrl-alt-boop/dribble/datasource"
	"github.com/ctrl-alt-boop/dribble/request"
	"github.com/ctrl-alt-boop/dribble/target"

	"github.com/ctrl-alt-boop/dribbler/internal/dribbleapi"
	"github.com/ctrl-alt-boop/dribbler/internal/keys"
	"github.com/ctrl-alt-boop/dribbler/internal/page"
	"github.com/ctrl-alt-boop/dribbler/internal/page/explorer"
	"github.com/ctrl-alt-boop/dribbler/logging"
)

// Model is the main model for the Dribbler bubbletea application
type Model struct {
	currentPage page.Page

	dribble *dribble.Client
}

// NewModel creates a new Model
func NewModel() Model {
	logging.GlobalLogger().Infof("NewModel")
	return Model{
		currentPage: explorer.NewExplorerPage(),
		dribble:     dribble.NewClient(),
	}
}

// Init implements tea.Model.
func (m Model) Init() tea.Cmd {
	// load config then check if there is a config.startup.target or similar

	//
	logging.GlobalLogger().Infof("Model.Init")
	if width, height, err := term.GetSize(0); err != nil {
		logging.GlobalLogger().Infof("error getting size: %v", err)
	} else {
		m.currentPage.SetSize(width, height)
	}
	return m.currentPage.Init()
}

// Update implements tea.Model.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// logging.GlobalLogger().Infof("Model.Update: %T: %v", msg, msg)

	switch msg := msg.(type) {
	case tea.KeyMsg: // Handle application wide key presses
		if key.Matches(msg, keys.Map.Quit) {
			return m, tea.Quit
		}

	case dribbleapi.DribbleApiMsg:
		return m, m.handleDribbleAPIMsg(msg)

	case dribbleapi.DribbleResponseMsg:
		//
	}

	updatedPage, cmd := m.currentPage.Update(msg)
	m.currentPage = updatedPage

	return m, cmd
}

// View implements tea.Model.
func (m Model) View() tea.View {
	canvas := m.currentPage.Render()

	view := tea.NewView(canvas.Render())
	view.WindowTitle = "Dribbler"
	view.AltScreen = true

	return view
}

func (m Model) handleDribbleAPIMsg(msg dribbleapi.DribbleApiMsg) tea.Cmd {
	logging.GlobalLogger().Infof("Model.handle.DribbleApiMsg: %T, %v", msg.Request, msg)
	switch msg := msg.Request.(type) {
	case dribbleapi.DSNOpenTargetMsg:
		return m.handleDSNOpenMsg(msg)
	case dribbleapi.AdapterOpenTargetMsg:
		return m.handleAdapterOpenMsg(msg)
	case dribbleapi.DataSourceRequestMsg:
		return m.handleDribbleRequestMsg(msg)
	}
	return nil
}

func (m Model) handleDSNOpenMsg(msg dribbleapi.DSNOpenTargetMsg) tea.Cmd {
	logging.GlobalLogger().Infof("Model.handleDSNOpenMsg: %v", msg)
	return func() tea.Msg {
		logging.GlobalLogger().Infof("handleDSNOpenMsg.cmd")
		t, err := target.New(msg.Name, msg.DSN)
		logging.GlobalLogger().Infof("target.New: %s, %s", msg.Name, msg.DSN)
		if err != nil {
			logging.GlobalLogger().Infof("error creating target: %v", err)
			return dribbleapi.NewAPIError(err)
		}
		if err := m.dribble.OpenTarget(context.TODO(), t); err != nil {
			logging.GlobalLogger().Infof("error opening target: %v", err)
			return dribbleapi.NewAPIError(err)
		}
		logging.GlobalLogger().Infof("target opened: %v", t)
		return dribbleapi.TargetOpened(t)
	}
}

func (m Model) handleAdapterOpenMsg(msg dribbleapi.AdapterOpenTargetMsg) tea.Cmd {
	panic("unimplemented")
}

func (m Model) handleDribbleRequestMsg(msg dribbleapi.DataSourceRequestMsg) tea.Cmd {
	var responseChan chan *request.Response
	var err error

	if msg.TargetID == nil { // Target all
		responseChan, err = m.dribble.RequestForAll(msg.Context, msg.Request)
	} else {
		responseChan, err = m.dribble.Request(msg.Context, fmt.Sprint(msg.TargetID), msg.Request) // FIXME: not fmt.Sprint
	}

	if err != nil {
		return dribbleapi.NewAPIError(err)
	}

	return func() tea.Msg {
		return dribbleapi.ResponseMsg{
			Channel: responseChan,
		}
	}
}

func (m Model) createOnSuccessCmds(onSuccess ...datasource.Request) []tea.Cmd {
	var cmds []tea.Cmd
	for _, req := range onSuccess {
		logging.GlobalLogger().Infof("on success, request: %v", req)
		cmd := func() tea.Msg {
			res, err := m.dribble.Request(context.TODO(), "msg.Name", req)
			logging.GlobalLogger().Infof("request: %v", req)
			if err != nil {
				logging.GlobalLogger().Infof("error requesting: %v", err)
				return dribbleapi.NewAPIError(err)
			}
			return dribbleapi.ResponseMsg{
				Channel: res,
			}
		}
		cmds = append(cmds, cmd)
	}
	return cmds
}

// Target tries to get target by name from the dribble client
func (m Model) Target(targetName string) (*dribbleapi.Requester, error) {
	target, ok := m.dribble.Target(targetName)
	if !ok {
		return nil, dribbleapi.NewTargetingError(targetName)
	}
	return &dribbleapi.Requester{
		Target: target,
	}, nil
}

// func (m Model) supportedModelsCmd() tea.Msg {
// 	return datastore.SupportedModelsMsg{
// 		Models: dribble.ListSupportedTargetTypes(),
// 	}
// }
