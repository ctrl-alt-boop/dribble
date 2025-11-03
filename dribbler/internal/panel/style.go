package panel

import "github.com/charmbracelet/lipgloss/v2"

type layoutStyle struct {
	normalStyle, focusedStyle lipgloss.Style
	hasFocusUnfocusStyles     bool
}
