// Package dribbler is the main model package
package dribbler

import (
	"github.com/charmbracelet/bubbles/v2/key"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/ctrl-alt-boop/dribble"
	"github.com/ctrl-alt-boop/dribbler/components"
	"github.com/ctrl-alt-boop/dribbler/datastore"
	"github.com/ctrl-alt-boop/dribbler/keys"
	"github.com/ctrl-alt-boop/dribbler/logging"
)

var logger = logging.GlobalLogger()

// Dribbler is the main model for the Dribbler bubbletea application
type Dribbler struct {
	Width, Height int

	MainContent *MainContentModel

	dribbleClient *dribble.Client
}

// NewDribblerModel creates a default Dribbler bubbletea model
func NewDribblerModel() Dribbler {
	m := Dribbler{
		dribbleClient: dribble.NewClient(),
	}

	m.MainContent = &MainContentModel{}

	return m
}

// Init implements tea.Model.
func (m Dribbler) Init() tea.Cmd {
	// config.LoadConfig()

	return m.MainContent.Init()
}

// Update implements tea.Model.
func (m Dribbler) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	updated := m
	var cmd tea.Cmd
	var cmds []tea.Cmd
	logger.Infof("Main Model got %T: %v", msg, msg)

	// Handle message that can be received at any time, regardless of focus or popups.
	switch msg := msg.(type) {
	case error:
		logger.Error(msg)

	case tea.KeyMsg: // Handle application wide key presses
		if key.Matches(msg, keys.Map.Quit) {
			return updated, tea.Quit
		} else if key.Matches(msg, keys.Map.Help) {
		}

	case tea.WindowSizeMsg: // Handle application window resize
		updated.Width, updated.Height = msg.Width, msg.Height
		updated.MainContent, cmd = updated.MainContent.Update(tea.WindowSizeMsg{Width: updated.Width, Height: updated.Height})
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
		return updated, tea.Batch(cmds...)
	}

	// Handle general dribbler message
	switch msg := msg.(type) {
	case components.UpdateHelpMsg:

		return updated, nil
	case datastore.DribbleRequestMsg:
		cmd := m.handleDribbleRequestMsg(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	updated.MainContent, cmd = updated.MainContent.Update(msg)

	return updated, tea.Batch(cmds...)
}

// View implements tea.Model.
func (m Dribbler) View() tea.View {
	contentView := m.MainContent.View()
	contentView.WindowTitle = "Dribbler"
	contentView.AltScreen = true

	return contentView
}
