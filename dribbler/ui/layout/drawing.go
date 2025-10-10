package layout

type Coord struct {
	X, Y int
}

func (c Coord) MoveOne(direction Direction) Coord {
	switch direction {
	case North:
		return Coord{X: c.X, Y: c.Y - 1}
	case South:
		return Coord{X: c.X, Y: c.Y + 1}
	case West:
		return Coord{X: c.X - 1, Y: c.Y}
	case East:
		return Coord{X: c.X + 1, Y: c.Y}
	default:
		return c
	}
}

func (c Coord) Move(directions ...Direction) Coord {
	newCoord := c
	for _, direction := range directions {
		newCoord = newCoord.MoveOne(direction)
	}
	return newCoord
}

type BoundingBox struct {
	TopLeft, BottomRight Coord
}

// top left, top right, bottom right, bottom left
func (b BoundingBox) AllCorners() []Coord {
	corners := make([]Coord, 4)
	corners[0] = b.TopLeft
	corners[1] = Coord{X: b.BottomRight.X, Y: b.TopLeft.Y}
	corners[2] = b.BottomRight
	corners[3] = Coord{X: b.TopLeft.X, Y: b.BottomRight.Y}
	return corners
}

func (b BoundingBox) GetCorners() (topLeft, topRight, bottomRight, bottomLeft Coord) {
	topLeft = b.TopLeft
	topRight = Coord{X: b.BottomRight.X, Y: b.TopLeft.Y}
	bottomRight = b.BottomRight
	bottomLeft = Coord{X: b.TopLeft.X, Y: b.BottomRight.Y}
	return
}
