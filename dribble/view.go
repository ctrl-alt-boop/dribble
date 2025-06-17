package dribble

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/ctrl-alt-boop/gooldb/dribble/ui"
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
	// if m.workspace.IsTableSet() && m.workspace.ViewportWidth() > m.workspace.Width {
	// 	rightSeparatorCorner = "┬"
	// }

	separator := lipgloss.JoinHorizontal(
		lipgloss.Top,
		"├",
		panelBorder,
		"┴",
		workspaceBorder,
		rightSeparatorCorner,
	)

	workspaceRender := ui.WorkspaceStyle.
		Width(m.workspace.ContentWidth).Height(m.workspace.ContentHeight).
		Render(workspaceView)

	popupView := m.popupHandler.View()
	if popupView != "" {
		// workspaceRender = ui.PlaceOverlayTest( // Actually not sure if this works just as good
		// 	lipgloss.Center,
		// 	lipgloss.Center,
		// 	popupView,
		// 	workspaceRender)
		workspaceRender = ui.PlaceOverlay(
			lipgloss.Center,
			lipgloss.Center,
			popupView,
			workspaceRender)
	}

	render := lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.JoinHorizontal(
			lipgloss.Top,
			panelView,
			workspaceRender,
		),
		separator,
		promptBarView,
		helpFooterView,
	)

	return lipgloss.Place(m.Width, m.Height, lipgloss.Left, lipgloss.Top, render)
}
