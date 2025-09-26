package layout

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Manager interface {
	Layout(width, height int, models []tea.Model) tea.Cmd
	View(models []tea.Model) string
}

var _ Manager = (*SimpleLayout)(nil)
var _ Manager = (*UniformGridLayout)(nil)
var _ Manager = (*PrioritySplitLayout)(nil)
var _ Manager = (*StackLayout)(nil)
var _ Manager = (*TabbedLayout)(nil)

type Direction int

const (
	Horizontal Direction = iota
	Vertical
)

func (d Direction) ToLipglossPosition() lipgloss.Position {
	switch d {
	case Horizontal:
		return lipgloss.Left
	case Vertical:
		return lipgloss.Top
	default:
		return lipgloss.Left
	}
}

type (
	SimpleLayout struct{}

	UniformGridLayout struct {
		Columns int

		// I'm not entierly sure how I want to do this, either letting the content decide or Layout decide...
		// One option is to reset the style of Children and using the one the parent decides
		GutterSize  int
		GutterRunes [2]rune // Horizontal, Vertical
	}

	PrioritySplitLayout struct {
		PrimaryIndex int
		Ratio        float64
		Direction    Direction

		// I'm not entierly sure how I want to do this, either letting the content decide or Layout decide...
		// One option is to reset the style of Children and using the one the parent decides
		GutterSize  int     // Only between areas
		GutterRunes [2]rune // Horizontal, Vertical
	}

	StackLayout struct {
		Direction Direction
	}

	TabbedLayout struct {
		ActiveIndex int

		TabStyle lipgloss.Style
	}
)
