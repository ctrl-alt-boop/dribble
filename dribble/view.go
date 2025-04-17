package dribble

import (
	"github.com/charmbracelet/lipgloss"
)

func (m AppModel) View() string {
	panelView := m.panel.View()
	promptBarView := m.prompt.View()
	workspaceView := m.workspace.View()
	helpFooterView := m.help.View()

	panelViewWidth := lipgloss.Width(panelView)
	promptBarViewWidth := lipgloss.Width(promptBarView)

	panelBorder := lipgloss.PlaceHorizontal(
		panelViewWidth-1,
		lipgloss.Left,
		"─", lipgloss.WithWhitespaceChars("─"))

	workspaceBorder := lipgloss.PlaceHorizontal(
		promptBarViewWidth-panelViewWidth-2,
		lipgloss.Left,
		"─", lipgloss.WithWhitespaceChars("─"))

	rightSeparatorCorner := "┤"
	if m.workspace.IsTableSet() {
		rightSeparatorCorner = "┬"
	}

	separator := lipgloss.JoinHorizontal(
		lipgloss.Left,
		"├",
		panelBorder,
		"┴",
		workspaceBorder,
		rightSeparatorCorner,
	)

	workspaceStyle := m.workspace.Style()

	popupView := m.popupHandler.View()
	if popupView != "" {
		workspaceView = popupView
		workspaceStyle = workspaceStyle.Align(lipgloss.Center, lipgloss.Center)
		// return lipgloss.Place(m.Width, m.Height,
		// 	lipgloss.Center, lipgloss.Center,
		// 	popupView,
		// 	lipgloss.WithWhitespaceChars(" "), lipgloss.WithWhitespaceBackground(lipgloss.Color("123")))
	} else {
		workspaceStyle = workspaceStyle.Align(lipgloss.Left, lipgloss.Top)
	}

	render := lipgloss.JoinVertical(
		lipgloss.Top,
		lipgloss.JoinHorizontal(
			lipgloss.Bottom,
			panelView,
			workspaceStyle.Render(workspaceView),
		),
		separator,
		promptBarView,
		helpFooterView,
	)

	return lipgloss.Place(m.Width, m.Height, lipgloss.Center, lipgloss.Center, render)
}
