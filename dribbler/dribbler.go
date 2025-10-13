package dribbler

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ctrl-alt-boop/dribble"
	"github.com/ctrl-alt-boop/dribbler/config"
	"github.com/ctrl-alt-boop/dribbler/widget"
)

const helpBarHeight = 1

type Model struct {
	Width, Height           int
	InnerWidth, InnerHeight int

	MainContent tea.Model

	showHelp     bool
	showFullHelp bool
	help         help.Model
	helpKeyMap   help.KeyMap

	popup tea.Model

	borderStyle lipgloss.Style

	dribbleClient *dribble.Client
}

func NewDribblerModel() Model {
	return Model{
		dribbleClient: dribble.NewClient(),
		MainContent:   CreateMainContent(),
		borderStyle:   lipgloss.NewStyle().Border(lipgloss.NormalBorder()),
		popup:         nil,
		help:          help.New(),
	}
}

func (m Model) Init() tea.Cmd {
	// config.LoadConfig()

	return m.MainContent.Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	updated := m
	var cmd tea.Cmd
	var cmds []tea.Cmd

	// Handle messages that can be received at any time, regardless of focus or popups.
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if key.Matches(msg, config.Keys.Quit) {
			return updated, tea.Quit
		} else if key.Matches(msg, config.Keys.Help) {
			updated.showFullHelp = !updated.showFullHelp
		}

	case tea.WindowSizeMsg:
		updated.Width, updated.Height = msg.Width, msg.Height

		updated.InnerWidth = msg.Width
		updated.InnerHeight = msg.Height

		updated.MainContent, cmd = updated.MainContent.Update(tea.WindowSizeMsg{Width: updated.InnerWidth, Height: updated.InnerHeight})
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
		return updated, tea.Batch(cmds...)
	}

	// Handle dribbler cmds
	switch msg := msg.(type) {
	case widget.UpdateHelpMsg:
		updated.helpKeyMap = msg.KeyMap
		return updated, nil
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

func (m Model) View() string {
	contentView := m.MainContent.View()
	// var helpView string
	// if m.showHelp {
	// 	helpView = lipgloss.NewStyle().Width(m.InnerWidth).Border(lipgloss.NormalBorder(), true, false, false).Render(m.help.View(m.helpKeyMap))
	// }

	// view := lipgloss.JoinVertical(lipgloss.Left, contentView, helpView)

	return contentView
}

func (m Model) PopupView() string {
	return ""
}
