// Package widget provides functionality for Ui
package component

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ctrl-alt-boop/dribble/database"
	"github.com/ctrl-alt-boop/dribbler/datastore"
)

type (
	// IsIdentifiable is used when a widget or component supports identification
	IsIdentifiable interface {
		ID() int
		SetID(int)
	}
	Identifiable struct {
		id int
	}
)

func (i Identifiable) ID() int {
	return i.id
}

func (i *Identifiable) SetID(id int) {
	i.id = id
}

type (
	// ShouldAlwaysUpdate is implemented by widgets or components that shall always be updated
	ShouldAlwaysUpdate interface {
		ShouldAlwaysUpdate()
	}
	AlwaysUpdate struct{}
)

func (a AlwaysUpdate) ShouldAlwaysUpdate() {}

type (
	// IsNamed is implemented by widgets or components that wants their name to be visible to other widgets or components
	IsNamed interface {
		SetName(string)
		Name() string
	}
	Named struct {
		name string
	}
)

func (n Named) Name() string {
	return n.name
}

func (n *Named) SetName(name string) {
	n.name = name
}

type targetBinding struct {
	name string
}

// Base for widgets, implements tea.Model except Update()
type Base struct {
	Identifiable
	Named

	model tea.Model

	waitingForResponse bool // Maybe request type?
	binding            targetBinding
}

// Init implements tea.Model
func (b Base) Init() tea.Cmd {
	return b.model.Init()
}

// View implements tea.Model
func (b Base) View() string {
	return b.model.View()
}

func (b Base) createTargetedRequest(r database.Request) tea.Cmd {
	if b.binding.name == "" {
		return func() tea.Msg {
			return fmt.Errorf("Widget: %d has no target binding", b.ID)
		}
	}
	return func() tea.Msg {
		return datastore.DribbleRequestMsg{
			TargetName: b.binding.name,
			Request:    r,
		}
	}
}

func (b Base) createAllRequest(r database.Request) tea.Cmd {
	return func() tea.Msg {
		return datastore.DribbleRequestMsg{
			TargetName: "*",
			Request:    r,
		}
	}
}
