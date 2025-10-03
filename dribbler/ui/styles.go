package ui

import "github.com/charmbracelet/lipgloss"

var (
	PanelBorder = lipgloss.Border{
		Top:         "─",
		Bottom:      "─",
		Left:        "│",
		Right:       "│",
		TopLeft:     "┌",
		TopRight:    "┬",
		BottomLeft:  "├",
		BottomRight: "┴",
	}

	PanelStyle = lipgloss.NewStyle().
			Border(PanelBorder, true, false, false, true).
			Align(lipgloss.Left, lipgloss.Top)
)

var (
	WorkspaceBorder = lipgloss.Border{
		Top:         "─",
		Bottom:      "─",
		Right:       "│",
		Left:        "│",
		TopLeft:     "┬",
		TopRight:    "┐",
		BottomLeft:  "┴",
		BottomRight: "┤",
	}

	WorkspaceStyle = lipgloss.NewStyle().
			Border(WorkspaceBorder, true, true, false, true).
			AlignHorizontal(lipgloss.Left).
			AlignVertical(lipgloss.Top)
)

var (
	PromptBorder = lipgloss.Border{
		Bottom:      "─",
		Top:         "─",
		Left:        "│",
		Right:       "│",
		TopLeft:     "├",
		TopRight:    "┤",
		BottomLeft:  "└",
		BottomRight: "┘",
	}

	PromptStyle = lipgloss.NewStyle().
			Border(PromptBorder, false, true, true, true).
			Align(lipgloss.Left, lipgloss.Center)
)

var (
	HelpStyle = lipgloss.NewStyle().
		Align(lipgloss.Left, lipgloss.Center).
		PaddingLeft(1)
)

var (
	PopupHandlerStyle = lipgloss.NewStyle().
				AlignHorizontal(lipgloss.Center).
				AlignVertical(lipgloss.Center)

	PopupStyle = lipgloss.NewStyle().
			Padding(1, 5).
			Border(lipgloss.RoundedBorder()).
			AlignHorizontal(lipgloss.Left).
			AlignVertical(lipgloss.Top)
)

const (
	DefaultHorizontalGutter = "\u2550"
	DefaultVerticalGutter   = "\u2551"
)

var (
	DefaultStyle  = lipgloss.NewStyle()
	DefaultBorder = lipgloss.Border{
		Top:         "─",
		Bottom:      "─",
		Left:        "│",
		Right:       "│",
		TopLeft:     "┌",
		TopRight:    "┐",
		BottomLeft:  "└",
		BottomRight: "┘",
	}
	DefaultBorderStyle = lipgloss.NewStyle().Border(DefaultBorder)
	DefaultGutter      = []rune{'\u2550', '\u2551'}

	DefaultBoxStyle = lipgloss.NewStyle().
			Border(DefaultBorder)
)

func DefaultBox(width, height int) lipgloss.Style {
	box := DefaultBoxStyle.Width(width).Height(height)
	return box.Width(width - box.GetHorizontalFrameSize()).Height(height - box.GetVerticalFrameSize())
}

func Inline(width, height int, render string) string {
	return DefaultStyle.Inline(true).MaxWidth(width).MaxHeight(height).Render(render)
}
