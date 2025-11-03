package layout

import (
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/ctrl-alt-boop/dribbler/internal/panel"
)

// var DebugBackgrounds = false

// type layoutDefinition struct {
// 	panels []panel.Definition

// 	panelBorder               lipgloss.Border
// 	normalStyle, focusedStyle lipgloss.Style

// 	indexForPosition      map[panel.Position]int
// 	hasFocusUnfocusStyles bool
// 	customPanelBorder     bool
// 	emptyCenter           bool
// 	allowNoFocus          bool
// }

// func (r *layoutDefinition) Update(opts ...layoutOption) {
// 	for _, option := range opts {
// 		option(r)
// 	}
// }

// func New(panels []panel.Definition, opts ...layoutOption) layoutDefinition {
// 	definition := layoutDefinition{
// 		panels:       make([]panel.Definition, len(panels)),
// 		normalStyle:  lipgloss.NewStyle(),
// 		focusedStyle: lipgloss.NewStyle(),

// 		customPanelBorder: false,
// 		indexForPosition:  make(map[panel.Position]int),
// 	}

// 	copy(definition.panels, panels)

// 	for _, option := range opts {
// 		option(&definition)
// 	}

// 	if definition.customPanelBorder {
// 		if definition.panelBorder.MiddleTop == "" || definition.panelBorder.MiddleBottom == "" || definition.panelBorder.MiddleLeft == "" || definition.panelBorder.MiddleRight == "" {
// 			panic("custom panel border requires all Middle... borders")
// 		}
// 	} else {
// 		definition.panelBorder = lipgloss.DoubleBorder()
// 	}

// 	definition.emptyCenter = true
// 	for i, def := range definition.panels {
// 		if def.Position == panel.Center {
// 			definition.emptyCenter = false
// 		}
// 		definition.indexForPosition[def.Position] = i
// 	}

// 	if !definition.hasFocusUnfocusStyles {
// 		definition.normalStyle = definition.focusedStyle
// 	}

// 	return definition
// }

// type layoutOption func(*layoutDefinition)

// // The custom border use the Middle... fields for connecting corners
// func WithPanelBorder(border lipgloss.Border) layoutOption {
// 	return func(renderModel *layoutDefinition) {
// 		renderModel.panelBorder = border
// 		renderModel.customPanelBorder = true
// 	}
// }

// func WithStyle(style lipgloss.Style) layoutOption {
// 	return func(renderModel *layoutDefinition) {
// 		renderModel.hasFocusUnfocusStyles = false
// 		style = style.UnsetBorderStyle().UnsetBorderTop().UnsetBorderRight().UnsetBorderBottom().UnsetBorderLeft()
// 		renderModel.normalStyle = style
// 		renderModel.focusedStyle = style
// 	}
// }

// func WithFocusedStyle(focusedStyle lipgloss.Style) layoutOption {
// 	return func(renderModel *layoutDefinition) {
// 		renderModel.hasFocusUnfocusStyles = true
// 		focusedStyle = focusedStyle.UnsetBorderStyle().UnsetBorderTop().UnsetBorderRight().UnsetBorderBottom().UnsetBorderLeft()
// 		renderModel.focusedStyle = focusedStyle
// 	}
// }

// func WithDefaultUnfocusedStyle() layoutOption {
// 	return func(renderModel *layoutDefinition) {
// 		renderModel.hasFocusUnfocusStyles = true
// 		renderModel.focusedStyle = renderModel.normalStyle
// 		renderModel.normalStyle = renderModel.focusedStyle.Faint(true)
// 	}
// }

// func AddCenterPanel() layoutOption {
// 	return func(renderModel *layoutDefinition) {
// 		if renderModel.emptyCenter {
// 			renderModel.panels = append(renderModel.panels, panel.Definition{ // Do I really want this?
// 				Position:            panel.Center,
// 				ShouldFillRemaining: true,
// 			})
// 			renderModel.emptyCenter = false
// 		}
// 	}
// }

// func AllowNoFocus() layoutOption {
// 	return func(renderModel *layoutDefinition) {
// 		renderModel.allowNoFocus = true
// 	}
// }

// func (r layoutDefinition) getBorder(definition panel.Definition) lipgloss.Border {
// 	newBorder := r.panelBorder
// 	if definition.TopLeftChar != "" {
// 		newBorder.TopLeft = definition.TopLeftChar
// 	}
// 	if definition.TopRightChar != "" {
// 		newBorder.TopRight = definition.TopRightChar
// 	}
// 	if definition.BottomLeftChar != "" {
// 		newBorder.BottomLeft = definition.BottomLeftChar
// 	}
// 	if definition.BottomRightChar != "" {
// 		newBorder.BottomRight = definition.BottomRightChar
// 	}
// 	return newBorder
// }

func UpdateBorders(panelBorder lipgloss.Border, panels ...*panel.State) {
	boundingBoxes := panel.GetAllBoundingBoxes(panels...)

	for i := range panels {
		borderTop, borderRight, borderBottom, borderLeft := false, false, false, false
		boundingBox := boundingBoxes[i]
		directionToBB := map[panel.Direction]panel.BoundingBox{}
		for j := range panels {
			if i == j {
				continue
			}

			otherBoundingBox := boundingBoxes[j]

			borderSide := checkSidesTouchingAndShorter(boundingBox, otherBoundingBox, i, j)

			switch borderSide {
			case panel.North:
				borderTop = true
				directionToBB[panel.North] = otherBoundingBox
			case panel.East:
				borderRight = true
				directionToBB[panel.East] = otherBoundingBox
			case panel.South:
				borderBottom = true
				directionToBB[panel.South] = otherBoundingBox
			case panel.West:
				borderLeft = true
				directionToBB[panel.West] = otherBoundingBox
			}
		}
		topLeft, topRight, bottomRight, bottomLeft := boundingBox.GetCorners()
		topLeftChar, topRightChar, bottomRightChar, bottomLeftChar := "", "", "", ""
		if borderTop && borderLeft { // corner index 0, check West & North
			numTouchesW := checkCornerTouches(topLeft, panel.West, directionToBB[panel.West])
			numTouchesN := checkCornerTouches(topLeft, panel.North, directionToBB[panel.North])
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
		if borderTop && borderRight { // corner index 1, check panel.North & East
			numTouchesN := checkCornerTouches(topRight, panel.North, directionToBB[panel.North])
			numTouchesE := checkCornerTouches(topRight, panel.East, directionToBB[panel.East])
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
		if borderBottom && borderRight { // corner index 2, check panel.East & South
			numTouchesE := checkCornerTouches(bottomRight, panel.East, directionToBB[panel.East])
			numTouchesS := checkCornerTouches(bottomRight, panel.South, directionToBB[panel.South])
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
		if borderBottom && borderLeft { // corner index 3, check panel.South & panel.West
			numTouchesS := checkCornerTouches(bottomLeft, panel.South, directionToBB[panel.South])
			numTouchesW := checkCornerTouches(bottomLeft, panel.West, directionToBB[panel.West])
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
func checkSidesTouchingAndShorter(this, that panel.BoundingBox, tieBreakerThis, tieBreakerThat int) (direction panel.Direction) {
	// --- Vertical Adjacency (Horizontal Alignment) ---

	// Check if 'this' is directly above 'that'
	if this.BottomRight.Y+1 == that.TopLeft.Y {
		// Check for shared segment on X-axis
		if (this.TopLeft.X <= that.BottomRight.X) && (that.TopLeft.X <= this.BottomRight.X) {
			lengthThis := this.BottomRight.X - this.TopLeft.X
			lengthThat := that.BottomRight.X - that.TopLeft.X
			if lengthThis < lengthThat {
				return panel.South
			} else if lengthThis == lengthThat && tieBreakerThis < tieBreakerThat {
				return panel.South
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
				return panel.North
			} else if lengthThis == lengthThat && tieBreakerThis < tieBreakerThat {
				return panel.North
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
				return panel.East
			} else if lengthThis == lengthThat && tieBreakerThis < tieBreakerThat {
				return panel.East
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
				return panel.West
			} else if lengthThis == lengthThat && tieBreakerThis < tieBreakerThat {
				return panel.West
			}
		}
	}

	return 0
}

// numInside: 1 corner to corner touch, 2 corner to side touch
func checkCornerTouches(thisCorner panel.Coord, dir panel.Direction, that panel.BoundingBox) (numInside int) {
	pointInside := thisCorner.Move(dir)

	switch dir {
	case panel.West, panel.East:
		if isInsideRect(pointInside.Move(panel.North), that) {
			numInside++
		}
		if isInsideRect(pointInside.Move(panel.South), that) {
			numInside++
		}

	case panel.North, panel.South:
		if isInsideRect(pointInside.Move(panel.West), that) {
			numInside++
		}
		if isInsideRect(pointInside.Move(panel.East), that) {
			numInside++
		}
	}

	return
}

func isInsideRect(point panel.Coord, rect panel.BoundingBox) bool {
	return point.X >= rect.TopLeft.X && point.X <= rect.BottomRight.X &&
		point.Y >= rect.TopLeft.Y && point.Y <= rect.BottomRight.Y
}
