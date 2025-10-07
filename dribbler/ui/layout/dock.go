package layout

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var _ Manager = (*DockLayout)(nil)

type DockLayout struct {
	managerBase

	usableWidth, usableHeight int
	usableX, usableY          int
}

func NewDockLayout() *DockLayout {
	return &DockLayout{
		managerBase: managerBase{
			focusPassThrough: false,
		},
	}
}

func (d *DockLayout) SetSize(width, height int) {
	d.Width = width
	d.Height = height
	d.usableWidth = width
	d.usableHeight = height
	d.usableX = 0
	d.usableY = 0
}

func (d *DockLayout) Layout(models []tea.Model) []tea.Model {
	if len(models) == 0 {
		return models
	}

	updatedModels := models
	updatedDefinitions := make([]LayoutDefinition, len(d.renderDefinition.Definitions))
	var center *LayoutDefinition
	fillRemainingIndex := -1

	for i, definition := range d.renderDefinition.Definitions {
		pos := definition.Position
		if pos == None {
			continue
		}
		updatedDefinition := d.Allocate(definition)

		if updatedDefinition.FillRemaining {
			center = &updatedDefinition
			fillRemainingIndex = i
		}
		updatedDefinitions[i] = updatedDefinition
	}
	if center == nil {
		panic("No fill remaining") // FIXME: temp panic
	}

	style := d.renderDefinition.FocusedStyle // TODO: Should we check and use focused and unfocused
	if d.renderDefinition.connectCorners {
		style = updatedDefinitions[fillRemainingIndex].Position.ConnectedCorners(d.renderDefinition.FocusedStyle, d.renderDefinition.PositionsInUse()...)
	}
	updatedDefinitions[fillRemainingIndex].actualWidth = d.usableWidth - style.GetHorizontalFrameSize()
	updatedDefinitions[fillRemainingIndex].actualHeight = d.usableHeight - style.GetVerticalFrameSize()

	d.renderDefinition.Definitions = updatedDefinitions

	for i, def := range updatedDefinitions {
		var model tea.Model
		if i < len(models) {
			model = models[i]
			msg := tea.WindowSizeMsg{Width: def.actualWidth, Height: def.actualHeight}
			updatedModel, _ := model.Update(msg)
			updatedModels[i] = updatedModel
		}
	}

	return updatedModels
}

func (d *DockLayout) Allocate(definition LayoutDefinition) LayoutDefinition {
	style := definition.Position.ConnectedCorners(d.renderDefinition.FocusedStyle, d.renderDefinition.PositionsInUse()...)
	updated := definition

	switch definition.Position {
	case Top:
		height := definition.MinHeight
		if height == 0 {
			height = definition.MaxHeight
		}
		height = min(height, d.usableHeight)

		width := d.usableWidth

		updated.actualX = d.usableX
		updated.actualY = d.usableY
		updated.actualWidth = width - style.GetHorizontalFrameSize()
		updated.actualHeight = height - style.GetVerticalFrameSize()

		d.usableY += height
		d.usableHeight -= height
	case Bottom:
		height := definition.MinHeight
		if height == 0 {
			height = definition.MaxHeight
		}
		height = min(height, d.usableHeight)

		width := d.usableWidth

		updated.actualX = d.usableX
		updated.actualY = d.usableY + d.usableY - height
		updated.actualWidth = width - style.GetHorizontalFrameSize()
		updated.actualHeight = height - style.GetVerticalFrameSize()

		d.usableHeight -= height
	case Left:
		width := definition.MinWidth
		if width == 0 {
			width = definition.MaxWidth
		}
		width = min(width, d.usableWidth)

		height := d.usableHeight

		updated.actualX = d.usableX
		updated.actualY = d.usableY
		updated.actualWidth = width - style.GetHorizontalFrameSize()
		updated.actualHeight = height - style.GetVerticalFrameSize()

		d.usableX += width
		d.usableWidth -= width
	case Right:
		width := definition.MinWidth
		if width == 0 {
			width = definition.MaxWidth
		}
		width = min(width, d.usableWidth)

		height := d.usableHeight

		updated.actualX = d.usableX
		updated.actualY = d.usableY
		updated.actualWidth = width - style.GetHorizontalFrameSize()
		updated.actualHeight = height - style.GetVerticalFrameSize()

		d.usableWidth -= width
	default:
		updated.FillRemaining = true
	}

	return updated
}

func (d *DockLayout) View(models []tea.Model) string {
	if len(models) == 0 {
		return ""
	}

	top, bottom := "", ""
	middle := make([]string, 3)

	for i, model := range models {
		definition := d.renderDefinition.Definitions[i]
		style := definition.Position.ConnectedCorners(d.renderDefinition.FocusedStyle, d.renderDefinition.PositionsInUse()...)
		render := style.
			Width(definition.actualWidth - style.GetHorizontalFrameSize()).Height(definition.actualHeight - style.GetVerticalFrameSize()).
			Render(model.View())
		switch definition.Position {
		case Top:
			top = render
		case Bottom:
			bottom = render
		case Left:
			middle[0] = render
		case Right:
			middle[2] = render
		default:
			middle[1] = render
		}
	}

	joinedMiddle := lipgloss.JoinHorizontal(lipgloss.Left, middle...) // This is fine since no newlines is done when horizontally joining empty string
	columns := []string{}
	if top != "" {
		columns = append(columns, top)
	}
	if joinedMiddle != "" {
		columns = append(columns, joinedMiddle)
	}
	if bottom != "" {
		columns = append(columns, bottom)
	}

	composed := lipgloss.JoinVertical(lipgloss.Top, columns...) // However, here we can't have an empty string

	return composed
}
