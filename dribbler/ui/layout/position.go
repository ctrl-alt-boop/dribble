package layout

import (
	"slices"

	"github.com/charmbracelet/lipgloss"
)

//go:generate stringer -type=Position

type Position int

const (
	None Position = iota
	Center
	Top
	Bottom
	Left
	Right
)

const Middle = Center

func (p Position) AsLipglossHorizontal() lipgloss.Position {
	switch p {
	case Left:
		return lipgloss.Left
	case Right:
		return lipgloss.Right
	default:
		return lipgloss.Center
	}
}

func (p Position) AsLipglossVertical() lipgloss.Position {
	switch p {
	case Top:
		return lipgloss.Top
	case Bottom:
		return lipgloss.Bottom
	default:
		return lipgloss.Center
	}
}

// At the moment this one works mostly for the middle row in a layout and when bottom and top are full-width
func (p Position) ConnectedCorners(style lipgloss.Style, positions ...Position) lipgloss.Style {
	newStyle := style.BorderTop(false).BorderBottom(false).BorderLeft(false).BorderRight(false)
	connectedCornersBorder := CreateConnectedCorners(style)
	switch p {
	case Left, Right:
		newStyle = newStyle.BorderStyle(connectedCornersBorder) // setting the borderstyle may be unnecessary for left and right panels
		if slices.Contains(positions, Top) {
			newStyle = newStyle.BorderTop(true)
		}
		if slices.Contains(positions, Bottom) {
			newStyle = newStyle.BorderBottom(true)
		}
	case Center:
		newStyle = newStyle.BorderStyle(connectedCornersBorder)
		if slices.Contains(positions, Left) {
			newStyle = newStyle.BorderLeft(true)
		}
		if slices.Contains(positions, Right) {
			newStyle = newStyle.BorderRight(true)
		}
		if slices.Contains(positions, Top) {
			newStyle = newStyle.BorderTop(true)
		}
		if slices.Contains(positions, Bottom) {
			newStyle = newStyle.BorderBottom(true)
		}
	}
	return newStyle
}

func (p Position) CalculateTopLeft(width, height int) (x, y int) {
	switch p {
	case Center:
		x, y = width/2, height/2
	case Left:
		x, y = 0, height/2
	case Right:
		x, y = width-1, height/2
	case Top:
		x, y = width/2, 0
	case Bottom:
		x, y = width/2, height-1
	default:
		x, y = 0, 0
	}
	return
}

func (p Position) CalculateTopLeftInBox(boxWidth, boxHeight int) (lipgloss.Position, lipgloss.Position) {
	x, y := p.CalculateTopLeft(boxWidth, boxHeight)
	return lipgloss.Position(HorizontalPosInBox(x, boxWidth)), lipgloss.Position(VerticalPosInBox(y, boxHeight))
}
