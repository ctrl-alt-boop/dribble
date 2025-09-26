package layout

import tea "github.com/charmbracelet/bubbletea"

func (s SimpleLayout) Layout(width, height int, models []tea.Model) tea.Cmd {
	if len(models) == 0 {
		return nil
	}

	// Issue a warning if there are too many children (optional, but helpful)
	if len(models) > 1 {
		// TODO: LOGGING
	}

	child := models[0]

	msg := tea.WindowSizeMsg{Width: width, Height: height}

	_, cmd := child.Update(msg)

	return cmd
}

func (s *SimpleLayout) View(models []tea.Model) string {
	if len(models) == 0 {
		return ""
	}

	// Only render the first model
	return models[0].View()
}
