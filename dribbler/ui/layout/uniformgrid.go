package layout

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ctrl-alt-boop/dribbler/ui"
)

var _ Manager = (*UniformGridLayout)(nil)

type UniformGridLayout struct {
	Columns int

	CellWidth, CellHeight int
	Width, Height         int
	X, Y                  int
	// I'm not entierly sure how I want to do this, either letting the content decide or Layout decide...
	// One option is to reset the style of Children and using the one the parent decides
	VerticalGutter, HorizontalGutter string

	renderDefinition RenderDefinition
}

func NewUniformGridLayout(numColumns int) *UniformGridLayout {
	return &UniformGridLayout{
		Columns:          numColumns,
		HorizontalGutter: ui.DefaultHorizontalGutter,
		VerticalGutter:   ui.DefaultVerticalGutter,
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
	// We subtract space for the gutters (g.Columns - 1)
	numGuttersX := g.Columns - 1
	totalHorizontalGutterWidth := numGuttersX * 1
	usableWidth := g.Width - totalHorizontalGutterWidth
	cellWidth := usableWidth / g.Columns

	// Determine the number of rows needed
	numChildren := len(models)
	numRows := (numChildren + g.Columns - 1) / g.Columns

	numGuttersY := numRows - 1
	totalVerticalGutterHeight := numGuttersY * 1
	usableHeight := g.Height - totalVerticalGutterHeight
	cellHeight := usableHeight / numRows

	// Send a WindowSizeMsg to all children to inform them of their new size
	msg := tea.WindowSizeMsg{Width: cellWidth, Height: cellHeight}

	updatedModels := make([]tea.Model, len(models))
	for i, model := range models {
		// Use tea.Model to send the message to the updatedModel
		updatedModel, _ := model.Update(msg)

		updatedModels[i] = updatedModel
	}

	g.CellWidth = cellWidth
	g.CellHeight = cellHeight

	return updatedModels
}

func (g *UniformGridLayout) View(models []tea.Model) string {
	if g.Columns <= 0 || len(models) == 0 {
		return ""
	}

	var rows []string
	currentRow := 0
	// 1. Iterate over children and group them into rows
	for i := 0; i < len(models); i += g.Columns {
		end := min(i+g.Columns, len(models))

		// Get the View output for all models in the current row
		var cellViews []string
		for j, model := range models[i:end] {
			style := lipgloss.NewStyle().Margin(1, 2)

			if j > 0 {
				cellViews = append(cellViews, lipgloss.PlaceVertical(
					g.CellHeight,
					lipgloss.Center,
					g.VerticalGutter, lipgloss.WithWhitespaceChars(g.VerticalGutter)))
			}
			cellViews = append(cellViews, style.Width(g.CellWidth-style.GetHorizontalFrameSize()).Height(g.CellHeight-style.GetVerticalFrameSize()).Render(model.View()))
		}

		// 2. Join the cells horizontally to form the row string
		rowString := lipgloss.JoinHorizontal(lipgloss.Top, cellViews...)
		style := lipgloss.NewStyle()
		if currentRow > 0 {
			rows = append(rows, lipgloss.PlaceHorizontal(
				g.CellWidth*g.Columns+(g.Columns-1),
				lipgloss.Center,
				g.HorizontalGutter, lipgloss.WithWhitespaceChars(g.HorizontalGutter)))
		}
		rows = append(rows, style.Render(rowString))
		currentRow++
	}

	// 3. Join the rows vertically to form the final grid
	return lipgloss.JoinVertical(lipgloss.Left, rows...)
}

func (g *UniformGridLayout) AddLayout(definition LayoutDefinition) {
	g.renderDefinition.Definitions = append(g.renderDefinition.Definitions, definition)
}

func (g *UniformGridLayout) GetDefinition() RenderDefinition {
	return g.renderDefinition
}

func (g *UniformGridLayout) GetLayout(index int) LayoutDefinition {
	return g.renderDefinition.Definitions[index]
}

// If position is not set, returns empty LayoutDefinition
func (g *UniformGridLayout) GetLayoutForPosition(position Position) LayoutDefinition {
	if index, ok := g.renderDefinition.indexForPosition[position]; ok {
		return g.renderDefinition.Definitions[index]
	}
	return LayoutDefinition{}
}

func (g *UniformGridLayout) SetDefinition(definition RenderDefinition) {
	g.renderDefinition = definition
}

func (g *UniformGridLayout) SetLayout(index int, definition LayoutDefinition) {
	g.renderDefinition.Definitions[index] = definition
}

func (g *UniformGridLayout) UpdateLayout(index int, opts ...LayoutOption) {
	for _, opt := range opts {
		opt(&g.renderDefinition.Definitions[index])
	}
}
