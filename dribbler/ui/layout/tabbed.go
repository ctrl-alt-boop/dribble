package layout

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var _ Manager = (*TabbedLayout)(nil)

type TabbedLayout struct {
	managerBase
	ActiveIndex   int
	TabStyle      lipgloss.Style
	TabsSide      Direction // FIXME: Implement
	TabLabels     []string
	contentHeight int
}

func NewTabbedLayout(tabsSide Direction, opts ...layoutOption) *TabbedLayout {
	tabBorder := lipgloss.NormalBorder()

	return &TabbedLayout{
		managerBase: managerBase{
			layoutDefinition: New(
				[]panelDefinition{},
				opts...,
			),
			focusPassThrough: false,
			focusedIndex:     -1,
		},
		ActiveIndex: 0,
		TabStyle:    lipgloss.NewStyle().Border(tabBorder, true).BorderForeground(lipgloss.Color("63")).Padding(0, 1),
		TabsSide:    tabsSide,
	}
}

func (t *TabbedLayout) SetSize(width, height int) {
	t.Width = width
	t.Height = height
}

func (t *TabbedLayout) GetSize() (width, height int) {
	return t.Width, t.Height
}

// Layout implements Manager.
func (t *TabbedLayout) Layout(models []tea.Model) []tea.Model {
	if len(models) == 0 {
		return models
	}

	if t.ActiveIndex < 0 || t.ActiveIndex >= len(models) {
		t.ActiveIndex = 0
	}

	var tabLabels []string
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

		tabLabels = append(tabLabels, tabLabel)
	}
	t.TabLabels = tabLabels

	activeModel := models[t.ActiveIndex]

	// Calculate available space for the active model
	tabHeight := t.TabStyle.GetVerticalFrameSize() + 1 // Assuming one line for tabs
	t.contentHeight = t.Height - tabHeight

	msg := tea.WindowSizeMsg{Width: t.Width, Height: t.contentHeight}
	updatedModel, _ := activeModel.Update(msg)
	models[t.ActiveIndex] = updatedModel

	return models
}

// View implements Manager.
func (t *TabbedLayout) View(models []tea.Model) string {
	if len(models) == 0 {
		return lipgloss.NewStyle().Width(t.Width).Height(t.Height).Render("")
	}

	if t.ActiveIndex < 0 || t.ActiveIndex >= len(models) {
		t.ActiveIndex = 0
	}

	tabs := t.RenderTabs()

	style := t.layoutDefinition.normalStyle.Width(t.Width).Height(t.contentHeight)

	return lipgloss.JoinVertical(lipgloss.Left, tabs, style.Render(models[t.ActiveIndex].View()))
}

func (t *TabbedLayout) RenderTabs() string {

	tabViews := make([]string, len(t.TabLabels))
	for i, tabLabel := range t.TabLabels {
		style := t.TabStyle
		border := style.GetBorderStyle()

		if i == t.ActiveIndex {
			style = style.Bold(true)

			border.Bottom = " "
			border.BottomRight = t.TabStyle.GetBorderStyle().BottomLeft
			border.BottomLeft = t.TabStyle.GetBorderStyle().BottomRight
		} else {
			style = style.Faint(true)

			border.BottomRight = border.MiddleBottom
			border.BottomLeft = border.MiddleBottom
		}
		if i == 0 {
			border.BottomLeft = t.TabStyle.GetBorderStyle().BottomRight
		}
		tabViews[i] = style.BorderStyle(border).Render(tabLabel)
	}
	tabs := lipgloss.JoinHorizontal(lipgloss.Top, tabViews...)
	tabsWidth := lipgloss.Width(tabs)

	sepThing := t.TabStyle.BorderTop(false).BorderRight(false).BorderLeft(false).Width(t.Width - tabsWidth).Height(1).Render(" ")
	return lipgloss.JoinHorizontal(lipgloss.Bottom, tabs, sepThing)
}
