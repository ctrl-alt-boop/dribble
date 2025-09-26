package layout

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Layout implements Manager.
func (t *TabbedLayout) Layout(width int, height int, models []tea.Model) tea.Cmd {
	if len(models) == 0 {
		return nil
	}

	if t.ActiveIndex < 0 || t.ActiveIndex >= len(models) {
		t.ActiveIndex = 0
	}

	activeModel := models[t.ActiveIndex]

	// Calculate available space for the active model
	tabHeight := t.TabStyle.GetVerticalFrameSize() + 1 // Assuming one line for tabs
	contentHeight := height - tabHeight

	msg := tea.WindowSizeMsg{Width: width, Height: contentHeight}
	_, cmd := activeModel.Update(msg)

	return cmd
}

// View implements Manager.
func (t *TabbedLayout) View(models []tea.Model) string {
	if len(models) == 0 {
		return ""
	}

	if t.ActiveIndex < 0 || t.ActiveIndex >= len(models) {
		t.ActiveIndex = 0
	}

	var tabViews []string
	for i, model := range models {
		// Assuming models have a Name() method or similar for tab labels
		// For now, let's use a generic label or try to cast to a known type
		var tabLabel string
		switch v := model.(type) {
		case interface{ Name() string }:
			tabLabel = v.Name()
		default:
			tabLabel = "Tab " + string(rune('A'+i))
		}

		style := t.TabStyle
		if i == t.ActiveIndex {
			style = style.BorderBottom(false).Bold(true)
		} else {
			style = style.Faint(true)
		}
		tabViews = append(tabViews, style.Render(tabLabel))
	}

	tabs := lipgloss.JoinHorizontal(lipgloss.Top, tabViews...)

	content := models[t.ActiveIndex].View()

	return lipgloss.JoinVertical(lipgloss.Left, tabs, content)

}
