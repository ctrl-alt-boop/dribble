package layout

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func NewUniformGridLayout(name string, numColumns int) *UniformGridLayout {
	return &UniformGridLayout{
		Columns:    numColumns,
		GutterSize: 1,
		GutterRunes: [2]rune{
			'\u2550',
			'\u2551',
		},
	}
}

func (g *UniformGridLayout) Layout(width, height int, models []tea.Model) tea.Cmd {
	if g.Columns <= 0 {
		return nil
	}

	// Calculate the usable width for each column.
	// We subtract space for the gutters (g.Columns - 1)
	usableWidth := width - (g.Columns-1)*g.GutterSize
	cellWidth := usableWidth / g.Columns

	// Determine the number of rows needed
	numChildren := len(models)
	numRows := (numChildren + (g.Columns-1)*g.GutterSize) / g.Columns
	cellHeight := height / numRows // Simple uniform height distribution

	// Send a WindowSizeMsg to all children to inform them of their new size
	var cmds []tea.Cmd
	msg := tea.WindowSizeMsg{Width: cellWidth, Height: cellHeight}

	for _, model := range models {
		// Use tea.Model to send the message to the child
		_, cmd := model.Update(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	return tea.Batch(cmds...)
}

func (g *UniformGridLayout) View(models []tea.Model) string {
	if g.Columns <= 0 || len(models) == 0 {
		return ""
	}

	var rows []string

	// 1. Iterate over children and group them into rows
	for i := 0; i < len(models); i += g.Columns {
		end := min(i+g.Columns, len(models))

		// Get the View output for all models in the current row
		var cellViews []string
		for _, model := range models[i:end] {
			cellViews = append(cellViews, model.View())
		}

		// 2. Join the cells horizontally to form the row string
		rowString := lipgloss.JoinHorizontal(lipgloss.Top, cellViews...)
		rows = append(rows, rowString)
	}

	// 3. Join the rows vertically to form the final grid
	return lipgloss.JoinVertical(lipgloss.Left, rows...)
}
