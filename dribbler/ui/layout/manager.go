package layout

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Manager interface {
	GetDefinition() layoutDefinition

	AddCenterPanel()

	UpdateLayout(opts ...layoutOption)
	UpdatePanel(index int, opts ...panelOption)

	Layout(models []tea.Model) []tea.Model
	View(models []tea.Model) string

	SetSize(width, height int)

	SetFocusPassThrough(v bool)
	GetFocusPassThrough() bool

	GetFocusable() []int
	IsNoFocusAllowed() bool
	SetFocusedIndex(index int)
	GetFocusedIndex() int
}

type managerBase struct {
	Width, Height int
	X, Y          int

	focusPassThrough bool
	layoutDefinition layoutDefinition
	focusedIndex     int
}

func (b *managerBase) UpdateLayout(opts ...layoutOption) {
	for _, opt := range opts {
		opt(&b.layoutDefinition)
	}
}

func (b *managerBase) GetFocusable() []int {
	focusable := []int{}
	for i, panel := range b.layoutDefinition.panels {
		if panel.focusable {
			focusable = append(focusable, i)
		}
	}
	return focusable
}

func (b *managerBase) IsNoFocusAllowed() bool {
	return b.layoutDefinition.allowNoFocus
}

func (b *managerBase) SetFocusedIndex(index int) {
	b.focusedIndex = index
}

func (b *managerBase) GetFocusedIndex() int {
	return b.focusedIndex
}

func (b *managerBase) AddCenterPanel() {
	b.layoutDefinition.Update(AddCenterPanel())
}

func (b *managerBase) GetDefinition() layoutDefinition {
	return b.layoutDefinition
}

func (b *managerBase) GetPanelDefinition(index int) panelDefinition {
	return b.layoutDefinition.panels[index]
}

// If position is not set, returns empty LayoutDefinition
func (b *managerBase) GetLayoutForPosition(position Position) panelDefinition {
	if index, ok := b.layoutDefinition.indexForPosition[position]; ok {
		return b.layoutDefinition.panels[index]
	}
	return panelDefinition{}
}

func (b *managerBase) SetLayout(index int, definition panelDefinition) {
	b.layoutDefinition.panels[index] = definition
}

func (b *managerBase) UpdatePanel(index int, opts ...panelOption) {
	if index >= len(b.layoutDefinition.panels) {
		panic("Index out of bounds")
	}
	for _, opt := range opts {
		opt(&b.layoutDefinition.panels[index])
	}
}

func (b *managerBase) GetFocusPassThrough() bool {
	return b.focusPassThrough
}

func (b *managerBase) SetFocusPassThrough(v bool) {
	b.focusPassThrough = v
}

func (b *managerBase) layout(models []tea.Model) []tea.Model {
	if len(models) == 0 {
		return models
	}

	updatedModels := models
	for i, def := range b.layoutDefinition.panels {
		var model tea.Model
		if i < len(models) {
			model = models[i]
			msg := tea.WindowSizeMsg{Width: def.actualWidth, Height: def.actualHeight}
			updatedModel, _ := model.Update(msg)
			updatedModels[i] = updatedModel
		}
	}
	b.layoutDefinition.updateBorders()
	return updatedModels
}

func (b *managerBase) getDefinitionStyle(index int) lipgloss.Style {
	style := b.layoutDefinition.normalStyle
	if b.layoutDefinition.hasFocusUnfocusStyles {
		if index == b.focusedIndex {
			style = b.layoutDefinition.focusedStyle
		}
	}
	if DebugBackgrounds {
		style = style.Background(lipgloss.Color(fmt.Sprintf("1%d3", index)))
	}
	definition := b.layoutDefinition.panels[index]

	border := b.layoutDefinition.getBorder(definition)
	style = style.Border(border, definition.topBorder, definition.rightBorder, definition.bottomBorder, definition.leftBorder).
		AlignHorizontal(definition.alignmentX).AlignVertical(definition.alignmentY)

	if definition.actualHeight == 0 || definition.actualWidth == 0 {
		return style.Width(0).Height(0)
	}
	return style.
		Width(definition.actualWidth - style.GetHorizontalFrameSize()).
		Height(definition.actualHeight - style.GetVerticalFrameSize())
}
