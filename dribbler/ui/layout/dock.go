package layout

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var _ Manager = (*DockLayout)(nil)

type DockLayout struct {
	managerBase

	usableWidth, usableHeight int
	usableX, usableY          int
}

func NewDockLayout(panels panelsDefinition, opts ...layoutOption) *DockLayout {
	return &DockLayout{
		managerBase: managerBase{
			layoutDefinition: New(
				panels,
				opts...,
			),
			focusPassThrough: false,
		},
	}
}

func (d *DockLayout) SetSize(width, height int) {
	d.Width = width - d.layoutDefinition.normalStyle.GetHorizontalFrameSize()
	d.Height = height - d.layoutDefinition.normalStyle.GetVerticalFrameSize()
	d.usableWidth = width
	d.usableHeight = height
	d.usableX = 0
	d.usableY = 0
}

func (d *DockLayout) Layout(models []tea.Model) []tea.Model {
	if len(models) == 0 {
		return models
	}

	updatedDefinitions := make([]panelDefinition, len(d.layoutDefinition.panels))
	fillRemainingIndex := -1

	for i, definition := range d.layoutDefinition.panels {
		pos := definition.position
		if pos == None {
			continue
		}
		updatedDefinition := d.Allocate(definition)

		if updatedDefinition.fillRemaining {
			fillRemainingIndex = i
		}
		updatedDefinitions[i] = updatedDefinition
	}
	if fillRemainingIndex != -1 {
		updatedDefinitions[fillRemainingIndex].actualWidth = d.usableWidth
		updatedDefinitions[fillRemainingIndex].actualHeight = d.usableHeight
		updatedDefinitions[fillRemainingIndex].actualX = d.usableX
		updatedDefinitions[fillRemainingIndex].actualY = d.usableY
	}

	d.layoutDefinition.panels = updatedDefinitions

	return d.layout(models)
}

func (d *DockLayout) Allocate(definition panelDefinition) panelDefinition {
	updated := definition

	if definition.widthRatio > 0.0 {
		definition.width = int(float64(d.Width) * definition.widthRatio)
	}
	if definition.heightRatio > 0.0 {
		definition.height = int(float64(d.Height) * definition.heightRatio)
	}

	switch definition.position {
	case Top:
		height := min(definition.height, d.usableHeight)

		width := d.usableWidth

		updated.actualWidth = width
		updated.actualHeight = height
		updated.actualX = d.usableX
		updated.actualY = 0

		d.usableY += height
		d.usableHeight -= height

	case Bottom:
		height := min(definition.height, d.usableHeight)

		width := d.usableWidth

		updated.actualWidth = width
		updated.actualHeight = height
		updated.actualX = d.usableX
		updated.actualY = d.Height - updated.actualHeight

		d.usableHeight -= height

	case Left:
		width := min(definition.width, d.usableWidth)

		height := d.usableHeight

		updated.actualWidth = width
		updated.actualHeight = height
		updated.actualX = 0
		updated.actualY = d.usableY

		d.usableX += width
		d.usableWidth -= width

	case Right:
		width := min(definition.width, d.usableWidth)

		height := d.usableHeight

		updated.actualWidth = width
		updated.actualHeight = height
		updated.actualX = d.Width - updated.actualWidth
		updated.actualY = d.usableY

		d.usableWidth -= width

	default:
		updated.fillRemaining = true
	}

	return updated
}

func (d *DockLayout) View(models []tea.Model) string {
	if len(models) == 0 || d.Height == 0 || d.Width == 0 {
		return lipgloss.NewStyle().Width(d.Width).Height(d.Height).Render("")
	}

	compositeLines := make([]string, d.Height)

	for i := range d.layoutDefinition.getXYOrderedIndices() {
		if i >= len(models) {
			break
		}

		definition := d.layoutDefinition.panels[i]
		model := models[i]

		style := d.getDefinitionStyle(i)
		if style.GetWidth() == 0 || style.GetHeight() == 0 {
			continue
		}

		render := style.Render(model.View())

		lineIndex := 0
		for line := range strings.Lines(render) {
			renderLine := strings.TrimRight(line, "\n")

			if len(compositeLines[definition.actualY+lineIndex]) > 0 && definition.actualX == 0 {
				renderLine = "\n" + renderLine
			}
			compositeLines[definition.actualY+lineIndex] += renderLine
			lineIndex++
		}
	}

	return strings.Join(compositeLines, "\n")
}
