package layout

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
