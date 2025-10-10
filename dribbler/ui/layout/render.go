package layout

import (
	"iter"
	"slices"

	"github.com/charmbracelet/lipgloss"
)

var DebugBackgrounds = false

type layoutDefinition struct {
	panels []panelDefinition

	panelBorder               lipgloss.Border
	normalStyle, focusedStyle lipgloss.Style

	indexForPosition      map[Position]int
	hasFocusUnfocusStyles bool
	customPanelBorder     bool
	emptyCenter           bool
}

func (r *layoutDefinition) Update(opts ...layoutOption) {
	for _, option := range opts {
		option(r)
	}
}

func (r layoutDefinition) getXOrderedIndices() iter.Seq[int] {
	orderedIndices := make([]int, len(r.panels))

	for i := range orderedIndices {
		orderedIndices[i] = i
	}
	slices.SortFunc(orderedIndices, func(i, j int) int {
		return r.panels[i].actualX - r.panels[j].actualX
	})
	return func(yield func(int) bool) {
		for _, index := range orderedIndices {
			if !yield(index) {
				break
			}
		}
	}
}

func New(panels []panelDefinition, opts ...layoutOption) layoutDefinition {
	definition := layoutDefinition{
		panels:       make([]panelDefinition, len(panels)),
		normalStyle:  lipgloss.NewStyle().Border(lipgloss.NormalBorder()),
		focusedStyle: lipgloss.NewStyle().Border(lipgloss.NormalBorder()),

		customPanelBorder: false,
		indexForPosition:  make(map[Position]int),
	}

	copy(definition.panels, panels)

	for _, option := range opts {
		option(&definition)
	}

	if definition.customPanelBorder {
		if definition.panelBorder.MiddleTop == "" || definition.panelBorder.MiddleBottom == "" || definition.panelBorder.MiddleLeft == "" || definition.panelBorder.MiddleRight == "" {
			panic("custom panel border requires all Middle... borders")
		}
	} else {
		definition.panelBorder = definition.normalStyle.GetBorderStyle()
	}

	definition.emptyCenter = true
	for i, def := range definition.panels {
		if def.position == Center {
			definition.emptyCenter = false
		}
		definition.indexForPosition[def.position] = i
	}

	if !definition.hasFocusUnfocusStyles {
		definition.normalStyle = definition.focusedStyle
	}

	return definition
}

type layoutOption func(*layoutDefinition)

// The custom border use the Middle... fields for connecting corners
func WithPanelBorder(border lipgloss.Border) layoutOption {
	return func(renderModel *layoutDefinition) {
		renderModel.panelBorder = border
		renderModel.customPanelBorder = true
	}
}

func WithStyle(style lipgloss.Style) layoutOption {
	return func(renderModel *layoutDefinition) {
		renderModel.hasFocusUnfocusStyles = false
		renderModel.normalStyle = style
		renderModel.focusedStyle = style
	}
}

func WithFocusedStyle(focusedStyle lipgloss.Style) layoutOption {
	return func(renderModel *layoutDefinition) {
		renderModel.hasFocusUnfocusStyles = true
		renderModel.focusedStyle = focusedStyle
	}
}

func WithDefaultUnfocusedStyle() layoutOption {
	return func(renderModel *layoutDefinition) {
		renderModel.hasFocusUnfocusStyles = true
		renderModel.focusedStyle = renderModel.normalStyle
		renderModel.normalStyle = renderModel.focusedStyle.Faint(true)
	}
}

func AddCenterPanel() layoutOption {
	return func(renderModel *layoutDefinition) {
		if renderModel.emptyCenter {
			renderModel.panels = append(renderModel.panels, panelDefinition{ // Do I really want this?
				position:      Center,
				fillRemaining: true,
			})
			renderModel.emptyCenter = false
		}
	}
}

func (r layoutDefinition) getBorder(definition panelDefinition) lipgloss.Border {
	newBorder := r.panelBorder
	if definition.topLeftChar != "" {
		newBorder.TopLeft = definition.topLeftChar
	}
	if definition.topRightChar != "" {
		newBorder.TopRight = definition.topRightChar
	}
	if definition.bottomLeftChar != "" {
		newBorder.BottomLeft = definition.bottomLeftChar
	}
	if definition.bottomRightChar != "" {
		newBorder.BottomRight = definition.bottomRightChar
	}
	return newBorder
}

func (r *layoutDefinition) updateBorders() {
	boundingBoxes := r.allBoundingBoxes()

	for i := range r.panels {

		borderTop, borderRight, borderBottom, borderLeft := false, false, false, false
		boundingBox := boundingBoxes[i]
		directionToBB := map[Direction]BoundingBox{}
		for j := range r.panels {
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
				topLeftChar = r.panelBorder.MiddleLeft
			} else if numTouchesN == 2 {
				topLeftChar = r.panelBorder.MiddleTop
			} else {
				topLeftChar = r.panelBorder.Middle
			}
		}
		if borderTop && borderRight { // corner index 1, check North & East
			numTouchesN := checkCornerTouches(topRight, North, directionToBB[North])
			numTouchesE := checkCornerTouches(topRight, East, directionToBB[East])
			if numTouchesE == 2 && numTouchesN == 2 { // Can't happen...
				panic("numTouchesN == 2 && numTouchesE == 2")
			} else if numTouchesN == 2 {
				topRightChar = r.panelBorder.MiddleTop
			} else if numTouchesE == 2 {
				topRightChar = r.panelBorder.MiddleRight
			} else {
				topRightChar = r.panelBorder.Middle
			}
		}
		if borderBottom && borderRight { // corner index 2, check East & South
			numTouchesE := checkCornerTouches(bottomRight, East, directionToBB[East])
			numTouchesS := checkCornerTouches(bottomRight, South, directionToBB[South])
			if numTouchesE == 2 && numTouchesS == 2 { // Can't happen...
				panic("numTouchesE == 2 && numTouchesS == 2")
			} else if numTouchesE == 2 {
				bottomRightChar = r.panelBorder.MiddleRight
			} else if numTouchesS == 2 {
				bottomRightChar = r.panelBorder.MiddleBottom
			} else {
				bottomRightChar = r.panelBorder.Middle
			}
		}
		if borderBottom && borderLeft { // corner index 3, check South & West
			numTouchesS := checkCornerTouches(bottomLeft, South, directionToBB[South])
			numTouchesW := checkCornerTouches(bottomLeft, West, directionToBB[West])
			if numTouchesW == 2 && numTouchesS == 2 { // Can't happen...
				panic("numTouchesS == 2 && numTouchesW == 2")
			} else if numTouchesS == 2 {
				bottomLeftChar = r.panelBorder.MiddleBottom
			} else if numTouchesW == 2 {
				bottomLeftChar = r.panelBorder.MiddleLeft
			} else {
				bottomLeftChar = r.panelBorder.Middle
			}
		}
		r.panels[i].topLeftChar = topLeftChar
		r.panels[i].topRightChar = topRightChar
		r.panels[i].bottomRightChar = bottomRightChar
		r.panels[i].bottomLeftChar = bottomLeftChar
		r.panels[i].topBorder = borderTop
		r.panels[i].rightBorder = borderRight
		r.panels[i].bottomBorder = borderBottom
		r.panels[i].leftBorder = borderLeft
	}
}

func (r layoutDefinition) allBoundingBoxes() []BoundingBox {
	boundingBoxes := make([]BoundingBox, len(r.panels))
	for i, definition := range r.panels {
		boundingBoxes[i] = definition.GetBoundingBox()
	}
	return boundingBoxes
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
