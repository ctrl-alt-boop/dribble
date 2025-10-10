package layout

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var _ Manager = (*UniformGridLayout)(nil)

type UniformGridLayout struct {
	managerBase
	Columns int
}

func NewUniformGridLayout(numColumns int, opts ...layoutOption) *UniformGridLayout {
	return &UniformGridLayout{
		managerBase: managerBase{
			layoutDefinition: New(
				[]panelDefinition{},
				opts...,
			),
			focusPassThrough: false,
		},
		Columns: numColumns,
	}
}

func (g *UniformGridLayout) SetSize(width, height int) {
	g.Width = width
	g.Height = height
}

func (g *UniformGridLayout) GetSize() (width, height int) {
	return g.Width, g.Height
}

func (g *UniformGridLayout) Layout(models []tea.Model) []tea.Model {
	if g.Columns <= 0 {
		return models
	}

	// Calculate the usable width for each column.
	baseCellWidth := g.Width / g.Columns
	widthRemainder := g.Width % g.Columns

	// Determine the number of rows needed
	numChildren := len(models)
	numRows := (numChildren + g.Columns - 1) / g.Columns

	baseCellHeight := g.Height / numRows
	heightRemainder := g.Height % numRows

	widths := make([]int, g.Columns)
	heights := make([]int, numRows)
	for i := range widths {
		if i < widthRemainder {
			widths[i] = baseCellWidth + 1
		} else {
			widths[i] = baseCellWidth
		}
	}
	for i := range heights {
		if i < heightRemainder {
			heights[i] = baseCellHeight + 1
		} else {
			heights[i] = baseCellHeight
		}
	}

	updatedDefinitions := make([]panelDefinition, len(models))
	currentY := 0
	for i := range numRows {
		currentX := 0
		for j := range g.Columns {
			index := i*g.Columns + j
			if index >= numChildren {
				break
			}
			updatedDefinitions[index].actualWidth = widths[j]
			updatedDefinitions[index].actualHeight = heights[i]
			updatedDefinitions[index].actualX = currentX
			updatedDefinitions[index].actualY = currentY

			currentX += widths[j]
		}
		currentY += heights[i]
	}

	g.layoutDefinition.panels = updatedDefinitions

	return g.layout(models)
}

func (g *UniformGridLayout) View(models []tea.Model) string {
	if g.Columns <= 0 || len(models) == 0 || g.Height == 0 || g.Width == 0 {
		return lipgloss.NewStyle().Width(g.Width).Height(g.Height).Render("")
	}

	var rows []string
	currentRow := 0
	// 1. Iterate over children and group them into rows
	for i := 0; i < len(models); i += g.Columns {
		end := min(i+g.Columns, len(models))

		// Get the View output for all models in the current row
		var cellViews []string
		for j, model := range models[i:end] {

			style := g.getDefinitionStyle(i + j)
			cellViews = append(cellViews, style.Render(model.View()))
		}

		// 2. Join the cells horizontally to form the row string
		rowString := lipgloss.JoinHorizontal(lipgloss.Top, cellViews...)

		rows = append(rows, rowString)
		currentRow++
	}

	// 3. Join the rows vertically to form the final grid
	return lipgloss.JoinVertical(lipgloss.Left, rows...)
}
