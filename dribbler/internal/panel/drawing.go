package panel

import "github.com/charmbracelet/lipgloss/v2"

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

func UpdateBorders(panelBorder lipgloss.Border, panels ...*State) {
	boundingBoxes := GetAllBoundingBoxes(panels...)

	for i := range panels {
		borderTop, borderRight, borderBottom, borderLeft := false, false, false, false
		boundingBox := boundingBoxes[i]
		directionToBB := map[Direction]BoundingBox{}
		for j := range panels {
			if i == j {
				continue
			}

			otherBoundingBox := boundingBoxes[j]

			borderSide := checkSidesTouchingAndShorter(boundingBox, otherBoundingBox, i, j)

			switch borderSide {
			case North:
				borderTop = true
				directionToBB[North] = otherBoundingBox
			case East:
				borderRight = true
				directionToBB[East] = otherBoundingBox
			case South:
				borderBottom = true
				directionToBB[South] = otherBoundingBox
			case West:
				borderLeft = true
				directionToBB[West] = otherBoundingBox
			}
		}
		topLeft, topRight, bottomRight, bottomLeft := boundingBox.GetCorners()
		topLeftChar, topRightChar, bottomRightChar, bottomLeftChar := "", "", "", ""
		if borderTop && borderLeft { // corner index 0, check West & North
			numTouchesW := checkCornerTouches(topLeft, West, directionToBB[West])
			numTouchesN := checkCornerTouches(topLeft, North, directionToBB[North])
			if numTouchesW == 2 && numTouchesN == 2 { // Can't happen...
				panic("numTouchesW == 2 && numTouchesN == 2")
			} else if numTouchesW == 2 {
				topLeftChar = panelBorder.MiddleLeft
			} else if numTouchesN == 2 {
				topLeftChar = panelBorder.MiddleTop
			} else {
				topLeftChar = panelBorder.Middle
			}
		}
		if borderTop && borderRight { // corner index 1, check North & East
			numTouchesN := checkCornerTouches(topRight, North, directionToBB[North])
			numTouchesE := checkCornerTouches(topRight, East, directionToBB[East])
			if numTouchesE == 2 && numTouchesN == 2 { // Can't happen...
				panic("numTouchesN == 2 && numTouchesE == 2")
			} else if numTouchesN == 2 {
				topRightChar = panelBorder.MiddleTop
			} else if numTouchesE == 2 {
				topRightChar = panelBorder.MiddleRight
			} else {
				topRightChar = panelBorder.Middle
			}
		}
		if borderBottom && borderRight { // corner index 2, check East & South
			numTouchesE := checkCornerTouches(bottomRight, East, directionToBB[East])
			numTouchesS := checkCornerTouches(bottomRight, South, directionToBB[South])
			if numTouchesE == 2 && numTouchesS == 2 { // Can't happen...
				panic("numTouchesE == 2 && numTouchesS == 2")
			} else if numTouchesE == 2 {
				bottomRightChar = panelBorder.MiddleRight
			} else if numTouchesS == 2 {
				bottomRightChar = panelBorder.MiddleBottom
			} else {
				bottomRightChar = panelBorder.Middle
			}
		}
		if borderBottom && borderLeft { // corner index 3, check South & West
			numTouchesS := checkCornerTouches(bottomLeft, South, directionToBB[South])
			numTouchesW := checkCornerTouches(bottomLeft, West, directionToBB[West])
			if numTouchesW == 2 && numTouchesS == 2 { // Can't happen...
				panic("numTouchesS == 2 && numTouchesW == 2")
			} else if numTouchesS == 2 {
				bottomLeftChar = panelBorder.MiddleBottom
			} else if numTouchesW == 2 {
				bottomLeftChar = panelBorder.MiddleLeft
			} else {
				bottomLeftChar = panelBorder.Middle
			}
		}
		panels[i].TopLeftChar = topLeftChar
		panels[i].TopRightChar = topRightChar
		panels[i].BottomRightChar = bottomRightChar
		panels[i].BottomLeftChar = bottomLeftChar
		panels[i].TopBorder = borderTop
		panels[i].RightBorder = borderRight
		panels[i].BottomBorder = borderBottom
		panels[i].LeftBorder = borderLeft
	}
}

// Checks if this and that are one unit adjacent on sides
func checkSidesTouchingAndShorter(this, that BoundingBox, tieBreakerThis, tieBreakerThat int) (direction Direction) {
	// --- Vertical Adjacency (Horizontal Alignment) ---

	// Check if 'this' is directly above 'that'
	if this.BottomRight.Y+1 == that.TopLeft.Y {
		// Check for shared segment on X-axis
		if (this.TopLeft.X <= that.BottomRight.X) && (that.TopLeft.X <= this.BottomRight.X) {
			lengthThis := this.BottomRight.X - this.TopLeft.X
			lengthThat := that.BottomRight.X - that.TopLeft.X
			if lengthThis < lengthThat {
				return South
			} else if lengthThis == lengthThat && tieBreakerThis < tieBreakerThat {
				return South
			}
		}
	}

	// Check if 'this' is directly below 'that'
	if this.TopLeft.Y-1 == that.BottomRight.Y {
		// Check for shared segment on X-axis
		if (this.TopLeft.X <= that.BottomRight.X) && (that.TopLeft.X <= this.BottomRight.X) {
			lengthThis := this.BottomRight.X - this.TopLeft.X
			lengthThat := that.BottomRight.X - that.TopLeft.X
			if lengthThis < lengthThat {
				return North
			} else if lengthThis == lengthThat && tieBreakerThis < tieBreakerThat {
				return North
			}
		}
	}

	// --- Horizontal Adjacency (Vertical Alignment) ---

	// Check if 'this' is directly to the left of 'that'
	if this.BottomRight.X+1 == that.TopLeft.X {
		// Check for shared segment on Y-axis
		if (this.TopLeft.Y <= that.BottomRight.Y) && (that.TopLeft.Y <= this.BottomRight.Y) {
			lengthThis := this.BottomRight.Y - this.TopLeft.Y
			lengthThat := that.BottomRight.Y - that.TopLeft.Y
			if lengthThis < lengthThat {
				return East
			} else if lengthThis == lengthThat && tieBreakerThis < tieBreakerThat {
				return East
			}
		}
	}

	// Check if 'this' is directly to the right of 'that'
	if this.TopLeft.X-1 == that.BottomRight.X {
		// Check for shared segment on Y-axis
		if (this.TopLeft.Y <= that.BottomRight.Y) && (that.TopLeft.Y <= this.BottomRight.Y) {
			lengthThis := this.BottomRight.Y - this.TopLeft.Y
			lengthThat := that.BottomRight.Y - that.TopLeft.Y
			if lengthThis < lengthThat {
				return West
			} else if lengthThis == lengthThat && tieBreakerThis < tieBreakerThat {
				return West
			}
		}
	}

	return 0
}

// numInside: 1 corner to corner touch, 2 corner to side touch
func checkCornerTouches(thisCorner Coord, dir Direction, that BoundingBox) (numInside int) {
	pointInside := thisCorner.Move(dir)

	switch dir {
	case West, East:
		if isInsideRect(pointInside.Move(North), that) {
			numInside++
		}
		if isInsideRect(pointInside.Move(South), that) {
			numInside++
		}

	case North, South:
		if isInsideRect(pointInside.Move(West), that) {
			numInside++
		}
		if isInsideRect(pointInside.Move(East), that) {
			numInside++
		}
	}

	return
}

func isInsideRect(point Coord, rect BoundingBox) bool {
	return point.X >= rect.TopLeft.X && point.X <= rect.BottomRight.X &&
		point.Y >= rect.TopLeft.Y && point.Y <= rect.BottomRight.Y
}
