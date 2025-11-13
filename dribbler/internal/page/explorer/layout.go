package explorer

import lipgloss "charm.land/lipgloss/v2"

var (
	unstyled  = lipgloss.NewStyle().Padding(0, 0)
	baseStyle = lipgloss.NewStyle().Padding(1, 2)
)

const (
	defaultSidebarWidthRatio = 0.2
	defaultSidebarMinWidth   = 50

	defaultCommandlineHeight = 1
	defaultHelpHeight        = 1
)

func getTabStyle(width int, current bool) lipgloss.Style {
	border := lipgloss.NormalBorder()
	if current {
		border.Bottom = " "
		return unstyled.AlignHorizontal(lipgloss.Center).Bold(true).BorderStyle(border).Width(width)
	}
	return unstyled.AlignHorizontal(lipgloss.Center).BorderStyle(border).Width(width)
}

type rect struct {
	x, y          int
	width, height int
}

func newRect(x, y, width, height int) *rect {
	return &rect{
		x:      x,
		y:      y,
		width:  width,
		height: height,
	}
}
