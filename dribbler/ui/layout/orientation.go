package layout

import (
	"github.com/charmbracelet/lipgloss"
)

//go:generate stringer -type=Direction
type Direction int

const (
	North Direction = iota + 1
	East
	South
	West
)

//go:generate stringer -type=Position
type Position int

const (
	None Position = iota

	Top
	Right
	Bottom
	Left

	Center
	Middle = Center
)

func (d Direction) AsLipglossHorizontal() lipgloss.Position {
	switch d {
	case West:
		return lipgloss.Left
	case East:
		return lipgloss.Right
	default:
		return lipgloss.Center
	}
}

func (d Direction) AsLipglossVertical() lipgloss.Position {
	switch d {
	case North:
		return lipgloss.Top
	case South:
		return lipgloss.Bottom
	default:
		return lipgloss.Center
	}
}

func (d Direction) Inverse() Direction {
	switch d {
	case North:
		return South
	case South:
		return North
	case West:
		return East
	case East:
		return West
	default:
		return d
	}
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

func (p Position) Opposite() Position {
	switch p {
	case Top:
		return Bottom
	case Right:
		return Left
	case Bottom:
		return Top
	case Left:
		return Right
	default:
		return None
	}
}
