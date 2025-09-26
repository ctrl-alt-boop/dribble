package layout

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Layout implements Manager.
func (p *PrioritySplitLayout) Layout(width int, height int, models []tea.Model) tea.Cmd {
	if len(models) == 0 {
		return nil
	}

	if p.PrimaryIndex < 0 || p.PrimaryIndex >= len(models) {
		p.PrimaryIndex = 0
	}

	if p.Ratio <= 0 || p.Ratio >= 1 {
		p.Ratio = 0.5 // Default to 50/50 split
	}

	var cmds []tea.Cmd

	// Calculate sizes for primary and secondary models
	var primaryWidth, primaryHeight int
	var secondaryWidth, secondarySizeForOne int

	gutterSpace := p.GutterSize
	if len(models) <= 1 {
		gutterSpace = 0
	}

	secondaryCount := len(models) - 1

	switch p.Direction {
	case Horizontal:
		totalUsableWidth := width - gutterSpace
		primaryWidth = int(float64(totalUsableWidth) * p.Ratio)
		secondaryWidth = totalUsableWidth - primaryWidth
		if secondaryCount > 0 {
			secondarySizeForOne = secondaryWidth / secondaryCount
		}
		primaryHeight = height
	case Vertical:
		totalUsableHeight := height - gutterSpace
		primaryHeight = int(float64(totalUsableHeight) * p.Ratio)
		secondaryHeight := totalUsableHeight - primaryHeight
		if secondaryCount > 0 {
			secondarySizeForOne = secondaryHeight / secondaryCount
		}
		primaryWidth = width
		secondaryWidth = width
	}

	for i, model := range models {
		var msg tea.WindowSizeMsg
		if i == p.PrimaryIndex {
			msg = tea.WindowSizeMsg{Width: primaryWidth, Height: primaryHeight}
		} else {
			// msg = tea.WindowSizeMsg{Width: secondaryWidth, Height: secondaryHeight}
			if p.Direction == Horizontal {
				msg = tea.WindowSizeMsg{Width: secondarySizeForOne, Height: height}
			} else { // Vertical
				msg = tea.WindowSizeMsg{Width: width, Height: secondarySizeForOne}
			}
		}

		_, cmd := model.Update(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	return tea.Batch(cmds...)
}

// View implements Manager.
func (p *PrioritySplitLayout) View(models []tea.Model) string {
	if len(models) == 0 {
		return ""
	}

	if p.PrimaryIndex < 0 || p.PrimaryIndex >= len(models) {
		p.PrimaryIndex = 0
	}

	var views []string
	for _, model := range models {
		views = append(views, model.View())
	}

	if len(views) == 1 {
		return views[0]
	}

	primaryView := views[p.PrimaryIndex]
	secondaryViews := make([]string, 0, len(views)-1)
	for i, view := range views {
		if i != p.PrimaryIndex {
			secondaryViews = append(secondaryViews, view)
		}
	}

	var separator string
	if p.GutterSize > 0 {
		switch p.Direction {
		case Horizontal:
			separator = strings.Repeat(string(p.GutterRunes[0]), p.GutterSize)
		case Vertical:
			separator = strings.Repeat(string(p.GutterRunes[1]), p.GutterSize)
		}
	}

	// Join primary and secondary views
	var joinedSecondary string
	switch p.Direction {
	case Horizontal:
		joinedSecondary = lipgloss.JoinVertical(lipgloss.Left, secondaryViews...)
	case Vertical:
		joinedSecondary = lipgloss.JoinHorizontal(lipgloss.Top, secondaryViews...)
	}

	// Final join with the primary view
	if p.Direction == Horizontal {
		if p.PrimaryIndex == 0 {
			return lipgloss.JoinHorizontal(lipgloss.Top, primaryView, separator, joinedSecondary)
		}
		return lipgloss.JoinHorizontal(lipgloss.Top, joinedSecondary, separator, primaryView)
	} else { // Vertical
		if p.PrimaryIndex == 0 {
			return lipgloss.JoinVertical(lipgloss.Left, primaryView, separator, joinedSecondary)
		}
		return lipgloss.JoinVertical(lipgloss.Left, joinedSecondary, separator, primaryView)
	}
}
