package panel

import lipgloss "charm.land/lipgloss/v2"

type layoutStyle struct {
	normalStyle, focusedStyle lipgloss.Style
	hasFocusUnfocusStyles     bool
}
