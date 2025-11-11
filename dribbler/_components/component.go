// Package component provides functionality for Ui
package components

import (
	tea "charm.land/bubbletea/v2"
	"github.com/ctrl-alt-boop/dribble/database"
	"github.com/ctrl-alt-boop/dribbler/internal/dribbleapi"
)

type Model interface {
	tea.Model

	// Render is the old View method
	Render() string
}

// Base for components, implements Model except Render and Update
type Base struct {
	Identification
	Name

	model Model

	binding binding
}

// Init implements Model
func (b Base) Init() tea.Cmd {
	return b.model.Init()
}

// View implements tea.Model
func (b Base) View() tea.View {
	return tea.NewView(b.model.Render())
}

type status int

const (
	statusNone status = iota
	statusWaiting
	statusDone
	statusError
)

type binding struct {
	target Identification
	status status
}

func (b binding) createTargetedRequest(r database.Request) tea.Cmd {
	return func() tea.Msg {
		id := b.target.ID()
		return dribbleapi.DribbleApiMsg{
			Request: dribbleapi.DataSourceRequestMsg{
				TargetID: &id,
				Request:  r,
			},
		}
	}
}

func (b binding) createAllRequest(r database.Request) tea.Cmd {
	return func() tea.Msg {
		return dribbleapi.DribbleApiMsg{
			Request: dribbleapi.DataSourceRequestMsg{
				TargetID: nil,
				Request:  r,
			},
		}
	}
}
