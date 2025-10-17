// Package dribbler is the main model package
package dribbler

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ctrl-alt-boop/dribble"
	"github.com/ctrl-alt-boop/dribbler/component"
	"github.com/ctrl-alt-boop/dribbler/config"
	"github.com/ctrl-alt-boop/dribbler/datastore"
	"github.com/ctrl-alt-boop/dribbler/logging"
)

var logger = logging.GlobalLogger()

// Model is the main model for the Dribbler bubbletea application
type Model struct {
	Width, Height int

	MainContent tea.Model

	showHelp     bool
	showFullHelp bool
	help         component.Help
	helpKeyMap   help.KeyMap

	popup tea.Model

	borderStyle lipgloss.Style

	dribbleClient *dribble.Client
}

// NewDribblerModel creates a default Dribbler bubbletea model
func NewDribblerModel() Model {
	m := Model{
		dribbleClient: dribble.NewClient(),
		borderStyle:   lipgloss.NewStyle().Border(lipgloss.NormalBorder()),
		popup:         nil,
		help:          component.NewHelp(),
		// help: help.Model{
		// 	ShortSeparator: " • ",
		// 	FullSeparator:  "    ",
		// 	Ellipsis:       "…",
		// },
		helpKeyMap: config.Keys,
	}

	m.MainContent = m.createMainContent()

	return m
}

// Init implements tea.Model.
func (m Model) Init() tea.Cmd {
	// config.LoadConfig()

	return m.MainContent.Init()
}

// Update implements tea.Model.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	updated := m
	var cmd tea.Cmd
	var cmds []tea.Cmd
	logger.Infof("Main Model got %T: %v", msg, msg)

	// Handle message that can be received at any time, regardless of focus or popups.
	switch msg := msg.(type) {
	case error:
		logger.Error(msg)

	case tea.KeyMsg: // Handle application wide key presses
		if key.Matches(msg, config.Keys.Quit) {
			return updated, tea.Quit
		} else if key.Matches(msg, config.Keys.Help) {
			updated.showFullHelp = !updated.showFullHelp
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
	case component.UpdateHelpMsg:
		updated.helpKeyMap = msg.KeyMap
		return updated, nil
	case datastore.DribbleRequestMsg:
		cmd := m.handleDribbleRequestMsg(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	// Handle popup
	if m.popup != nil {
		m.popup, cmd = m.popup.Update(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	} else { // Handle focused
		updated.MainContent, cmd = updated.MainContent.Update(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	return updated, tea.Batch(cmds...)
}

// View implements tea.Model.
func (m Model) View() string {
	contentView := m.MainContent.View()
	// var helpView string
	// if m.showHelp {
	// 	helpView = lipgloss.NewStyle().Width(m.InnerWidth).Border(lipgloss.NormalBorder(), true, false, false).Render(m.help.View(m.helpKeyMap))
	// }

	// view := lipgloss.JoinVertical(lipgloss.Left, contentView, helpView)

	return contentView
}

// PopupView creates the view string used when a popup is present
func (m Model) PopupView(view string) string {
	return ""
}
