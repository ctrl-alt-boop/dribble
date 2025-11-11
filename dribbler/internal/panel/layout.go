package panel

import (
	lipgloss "charm.land/lipgloss/v2"
)

var DebugBackgrounds = true

type LayoutManager interface {
	Layout(pm *Manager)
	Render(pm Manager) *lipgloss.Canvas
}

type LayoutManagerBase struct {
	Width, Height int

	style layoutStyle
}

func (b *LayoutManagerBase) layout(panelManager *Manager) {
	// for i, panel := range panelManager.Panels {

	// }
}

// func (b *layoutManagerBase) getPanelStyle(index int) lipgloss.Style {
// 	style := b.style.normalStyle
// 	if b.style.hasFocusUnfocusStyles {
// 		// if index == b.focusedIndex {
// 		// 	style = b.layoutDefinition.focusedStyle
// 		// }
// 	}
// 	if DebugBackgrounds {
// 		style = style.Background(lipgloss.Color(fmt.Sprintf("1%d3", index)))
// 	}
// 	definition := b.layoutDefinition.panels[index]

// 	border := b.layoutDefinition.getBorder(definition)
// 	style = style.Border(border, definition.topBorder, definition.rightBorder, definition.bottomBorder, definition.leftBorder).
// 		AlignHorizontal(definition.alignmentX).AlignVertical(definition.alignmentY)

// 	if definition.actualHeight == 0 || definition.actualWidth == 0 {
// 		return style.Width(0).Height(0)
// 	}
// 	return style.
// 		Width(definition.actualWidth - style.GetHorizontalFrameSize()).
// 		Height(definition.actualHeight - style.GetVerticalFrameSize())
// }
