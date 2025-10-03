package layout

import "github.com/charmbracelet/lipgloss"

func (p Position) CalculateTopLeft(width, height int) (x, y int) {
	switch p {
	// case TopLeft:
	// 	x, y = 0, 0
	// case TopCenter:
	// 	x, y = width/2, 0
	// case TopRight:
	// 	x, y = width-1, 0
	// case MiddleLeft:
	// 	x, y = 0, height/2
	// case MiddleCenter:
	// 	x, y = width/2, height/2
	// case MiddleRight:
	// 	x, y = width-1, height/2
	// case BottomLeft:
	// 	x, y = 0, height-1
	// case BottomCenter:
	// 	x, y = width/2, height-1
	// case BottomRight:
	// 	x, y = width-1, height-1
	default:
		x, y = 0, 0
	}
	return
}

type Coord struct {
	X, Y int
}
type BoundingBox struct {
	TopLeft, BottomRight Coord
}

func (b BoundingBox) WillCollide(other BoundingBox) bool {
	return b.TopLeft.X < other.BottomRight.X && b.BottomRight.X > other.TopLeft.X && b.TopLeft.Y < other.BottomRight.Y && b.BottomRight.Y > other.TopLeft.Y
}

func (p Position) CalculateBoundingBox(width, height int) BoundingBox {
	topLeftX, topLeftY := p.CalculateTopLeft(width, height)
	bottomRightX, bottomRightY := topLeftX+width-1, topLeftY+height-1
	return BoundingBox{
		TopLeft:     Coord{X: topLeftX, Y: topLeftY},
		BottomRight: Coord{X: bottomRightX, Y: bottomRightY},
	}
}

func (p Position) CalculateTopLeftInBox(boxWidth, boxHeight int) (lipgloss.Position, lipgloss.Position) {
	x, y := p.CalculateTopLeft(boxWidth, boxHeight)
	return lipgloss.Position(HorizontalPosInBox(x, boxWidth)), lipgloss.Position(VerticalPosInBox(y, boxHeight))
}

func HorizontalPosInBox(x, boxWidth int) float64 {
	return float64(x) / float64(boxWidth)
}

func VerticalPosInBox(y, boxHeight int) float64 {
	return float64(y) / float64(boxHeight)
}

func WillMaxFit(l1, l2 LayoutDefinition) bool {
	l1MaxBoundingBox := l1.Position.CalculateBoundingBox(l1.MaxWidth, l1.MaxHeight)
	l2MaxBoundingBox := l2.Position.CalculateBoundingBox(l2.MaxWidth, l2.MaxHeight)

	return !l1MaxBoundingBox.WillCollide(l2MaxBoundingBox)
}

func WillMinFit(l1, l2 LayoutDefinition) bool {
	l1MinBoundingBox := l1.Position.CalculateBoundingBox(l1.MinWidth, l1.MinHeight)
	l2MinBoundingBox := l2.Position.CalculateBoundingBox(l2.MinWidth, l2.MinHeight)

	return !l1MinBoundingBox.WillCollide(l2MinBoundingBox)
}
