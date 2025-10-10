package layout

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var _ Manager = (*TabbedLayout)(nil)

type TabbedLayout struct {
	managerBase
	ActiveIndex int
	TabStyle    lipgloss.Style
	TabsSide    Direction // FIXME: Implement
}

func NewTabbedLayout(tabsSide Direction, opts ...layoutOption) *TabbedLayout {
	return &TabbedLayout{
		managerBase: managerBase{
			layoutDefinition: New(
				[]panelDefinition{},
				opts...,
			),
			focusPassThrough: false,
		},
		ActiveIndex: 0,
		TabStyle:    lipgloss.NewStyle().Border(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("63")).Padding(0, 1),
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

	activeModel := models[t.ActiveIndex]

	// Calculate available space for the active model
	tabHeight := t.TabStyle.GetVerticalFrameSize() + 1 // Assuming one line for tabs
	contentHeight := t.Height - tabHeight

	msg := tea.WindowSizeMsg{Width: t.Width, Height: contentHeight}
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

	tabsStyle := lipgloss.NewStyle().Align(lipgloss.Left, lipgloss.Center)
	tabsRender := tabsStyle.Render(tabs)
	tabsHeight := lipgloss.Height(tabsRender)

	contentStyle := lipgloss.NewStyle().Height(t.Height-tabsHeight).Width(t.Width).Align(lipgloss.Center, lipgloss.Center)

	return lipgloss.JoinVertical(lipgloss.Left, tabsRender, contentStyle.Render(models[t.ActiveIndex].View()))
}
