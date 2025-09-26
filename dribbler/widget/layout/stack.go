package layout

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Layout implements Manager.
func (s *StackLayout) Layout(width int, height int, models []tea.Model) tea.Cmd {
	if len(models) == 0 {
		return nil
	}

	var cmds []tea.Cmd
	for _, model := range models {
		msg := tea.WindowSizeMsg{Width: width, Height: height}
		_, cmd := model.Update(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}
	return tea.Batch(cmds...)
}

// View implements Manager.
func (s *StackLayout) View(models []tea.Model) string {
	if len(models) == 0 {
		return ""
	}

	var views []string
	for _, model := range models {
		views = append(views, model.View())
	}

	switch s.Direction {
	case Horizontal:
		return lipgloss.JoinHorizontal(lipgloss.Top, views...)
	case Vertical:
		return lipgloss.JoinVertical(lipgloss.Left, views...)
	default:
		return lipgloss.JoinHorizontal(lipgloss.Top, views...)
	}
}
