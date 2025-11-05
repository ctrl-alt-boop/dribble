package components

import (
	"github.com/charmbracelet/bubbles/v2/help"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/ctrl-alt-boop/dribble/database"
	"github.com/ctrl-alt-boop/dribble/request"
	"github.com/ctrl-alt-boop/dribbler/keys"
)

type (
	OpenCellDataMsg struct {
		Value string
	}

	OpenIntentBuilderMsg struct {
		Method database.RequestType
		Table  string
	}

	IntentBuilderInitMsg struct {
		Database string
		Method   database.RequestType
		Table    string
	}

	IntentBuilderDataMsg struct{}

	IntentBuilderConfirmMsg struct {
		Intent *request.Intent
	}

	// GotFocusMsg is sent to a widget or component when it is focused
	GotFocusMsg struct{}
	// LostFocusMsg is sent to a widget or component when it loses focused
	LostFocusMsg struct{}

	// FocusMsg is to be used when the main model wants to grant focus to a widget or component
	FocusMsg struct {
		Index []int
	}

	// LoseFocusMsg is to be used when a widget or component wants to lose focus
	LoseFocusMsg struct {
		ID int
	}
	// RedirectFocusMsg is to be used when a widget or component wants to send the focus to someone else
	RedirectFocusMsg struct {
		ID int
	}

	// UpdateHelpMsg is used by the View method of the Help widget
	UpdateHelpMsg struct {
		KeyMap help.KeyMap
	}
)

// LoseFocusCmd is a helper for LoseFocusMsg
func LoseFocusCmd(id int) tea.Cmd {
	return func() tea.Msg {
		return LoseFocusMsg{
			ID: id,
		}
	}
}

// RedirectFocusCmd is a helper for RedirectFocusMsg
func RedirectFocusCmd(id int) tea.Cmd {
	return func() tea.Msg {
		return RedirectFocusMsg{
			ID: id,
		}
	}
}

// UpdateHelpCmd is a small helper method for creating a cmd to update the help bar
func (f GotFocusMsg) UpdateHelpCmd(keys keys.KeyMap) tea.Cmd {
	return func() tea.Msg {
		return UpdateHelpMsg{KeyMap: keys}
	}
}
