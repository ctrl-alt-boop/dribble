package panel

import (
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"
	"github.com/ctrl-alt-boop/dribbler/internal/core/util"
	"github.com/ctrl-alt-boop/dribbler/keys"
	"github.com/ctrl-alt-boop/dribbler/logging"
)

type Manager struct {
	composer    Composer
	composition *Composition

	Panels []Model

	Width, Height int

	focusRing *util.Ring[int]

	normalStyle  lipgloss.Style
	focusedStyle lipgloss.Style
	panelBorder  lipgloss.Border
}

func NewPanelManager(composer Composer, panels ...Model) *Manager {
	composer.SetNumPanels(len(panels))
	return &Manager{
		composer:     composer,
		Panels:       panels,
		focusRing:    util.NewRing(0, 1, 2), // FIXME: Fix focus ring
		normalStyle:  composer.GetStyle().Border(composer.GetPanelBorder()),
		focusedStyle: composer.GetFocusedStyle().Border(composer.GetPanelBorder()),
		panelBorder:  composer.GetPanelBorder(),
	}
}

func (p Manager) GetFocused() int {
	return p.focusRing.Value()
}

func (p Manager) Init() tea.Cmd {
	return nil
}

func (p *Manager) Update(msg tea.Msg) (*Manager, tea.Cmd) {
	var cmd tea.Cmd
	updated := p
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		updated.Layout(msg.Width, msg.Height)
		return updated, nil

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Map.CycleViewPrev):
			updated.focusRing.Backward()
			return updated, nil
		case key.Matches(msg, keys.Map.CycleViewNext):
			updated.focusRing.Forward()
			return updated, nil
		}
	}
	for i, panel := range p.Panels {
		updated.Panels[i], cmd = panel.Update(msg)
		if cmd != nil {
			return updated, cmd
		}
	}

	return updated, cmd
}

func (p Manager) Render() string {
	if p.composition == nil || len(p.composition.Layers) == 0 {
		return ""
	}
	render := lipgloss.NewCanvas()

	for i, panl := range p.Panels {
		style := p.normalStyle
		z := 0
		if i == p.focusRing.Value() {
			style = p.focusedStyle
			z = 1
		}
		panelBorder := p.panelBorder
		panl.SetBorderStyle(panelBorder)
		logging.GlobalLogger().Infof("panelBorder: %+v", panelBorder)

		panelRender := style.
			Width(p.composition.Layers[i].GetWidth()).
			Height(p.composition.Layers[i].GetHeight()).
			BorderStyle(panelBorder).
			Render(panl.Canvas().Render())

		layer := p.composition.Layers[i].
			SetContent(panelRender).
			Z(z)

		render.AddLayers(layer)
	}

	render.AddLayers(p.composition.PreRendered...)

	return render.Render()
}

func (p Manager) View() tea.View {
	return tea.NewView(p.Render())
}

func (p *Manager) Layout(width, height int) {
	p.Width, p.Height = width, height
	p.composition = p.composer.Compose(width, height)
	p.updateBorders()
}

func (p *Manager) String() string {
	return p.Render()
}

func (p *Manager) updateBorders() {
	boundingBoxes := GetAllBoundingBoxes(p.composition.states...)

	for i := range p.composition.states {
		borderTop, borderRight, borderBottom, borderLeft := false, false, false, false
		boundingBox := boundingBoxes[i]
		directionToBB := map[Direction]BoundingBox{}
		for j := range p.composition.states {
			if i == j {
				continue
			}

			otherBoundingBox := boundingBoxes[j]

			borderSide := checkSidesTouchingAndShorter(boundingBox, otherBoundingBox)

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
		topLeftChar, topRightChar, bottomRightChar, bottomLeftChar := p.panelBorder.TopLeft, p.panelBorder.TopRight, p.panelBorder.BottomRight, p.panelBorder.BottomLeft
		if borderTop && borderLeft { // corner index 0, check West & North
			numTouchesW := checkCornerTouches(topLeft, West, directionToBB[West])
			numTouchesN := checkCornerTouches(topLeft, North, directionToBB[North])
			if numTouchesW == 2 && numTouchesN == 2 { // Can't happen...
				panic("numTouchesW == 2 && numTouchesN == 2")
			} else if numTouchesW == 2 {
				topLeftChar = p.panelBorder.MiddleLeft
			} else if numTouchesN == 2 {
				topLeftChar = p.panelBorder.MiddleTop
			} else {
				topLeftChar = p.panelBorder.Middle
			}
		}
		if borderTop && borderRight { // corner index 1, check North & East
			numTouchesN := checkCornerTouches(topRight, North, directionToBB[North])
			numTouchesE := checkCornerTouches(topRight, East, directionToBB[East])
			if numTouchesE == 2 && numTouchesN == 2 { // Can't happen...
				panic("numTouchesN == 2 && numTouchesE == 2")
			} else if numTouchesN == 2 {
				topRightChar = p.panelBorder.MiddleTop
			} else if numTouchesE == 2 {
				topRightChar = p.panelBorder.MiddleRight
			} else {
				topRightChar = p.panelBorder.Middle
			}
		}
		if borderBottom && borderRight { // corner index 2, check East & South
			numTouchesE := checkCornerTouches(bottomRight, East, directionToBB[East])
			numTouchesS := checkCornerTouches(bottomRight, South, directionToBB[South])
			if numTouchesE == 2 && numTouchesS == 2 { // Can't happen...
				panic("numTouchesE == 2 && numTouchesS == 2")
			} else if numTouchesE == 2 {
				bottomRightChar = p.panelBorder.MiddleRight
			} else if numTouchesS == 2 {
				bottomRightChar = p.panelBorder.MiddleBottom
			} else {
				bottomRightChar = p.panelBorder.Middle
			}
		}
		if borderBottom && borderLeft { // corner index 3, check South & West
			numTouchesS := checkCornerTouches(bottomLeft, South, directionToBB[South])
			numTouchesW := checkCornerTouches(bottomLeft, West, directionToBB[West])
			if numTouchesW == 2 && numTouchesS == 2 { // Can't happen...
				panic("numTouchesS == 2 && numTouchesW == 2")
			} else if numTouchesS == 2 {
				bottomLeftChar = p.panelBorder.MiddleBottom
			} else if numTouchesW == 2 {
				bottomLeftChar = p.panelBorder.MiddleLeft
			} else {
				bottomLeftChar = p.panelBorder.Middle
			}
		}
		borderStyle := lipgloss.Border{
			Top:    p.panelBorder.Top,
			Bottom: p.panelBorder.Bottom,
			Left:   p.panelBorder.Left,
			Right:  p.panelBorder.Right,

			TopLeft:     topLeftChar,
			TopRight:    topRightChar,
			BottomLeft:  bottomLeftChar,
			BottomRight: bottomRightChar,

			MiddleLeft:   p.panelBorder.MiddleLeft,
			MiddleRight:  p.panelBorder.MiddleRight,
			MiddleTop:    p.panelBorder.MiddleTop,
			MiddleBottom: p.panelBorder.MiddleBottom,
			Middle:       p.panelBorder.Middle,
		}
		p.Panels[i].SetBorderStyle(borderStyle)
	}
}

// checkSidesTouchingAndShorter checks if 'this' and 'that' share an edge (1-unit overlap).
// It returns the direction to 'that' (the other panel) from 'this', ignoring tiebreakers.
func checkSidesTouchingAndShorter(this, that BoundingBox) (direction Direction) {
	// --- Vertical Adjacency (X-Alignment) ---

	// Case 1: 'this' is immediately above 'that' (that is to the South)
	if this.BottomRight.Y == that.TopLeft.Y {
		// Check for shared segment on X-axis (horizontal overlap)
		if (this.TopLeft.X <= that.BottomRight.X) && (that.TopLeft.X <= this.BottomRight.X) {
			return South
		}
	}

	// Case 2: 'this' is immediately below 'that' (that is to the North)
	if this.TopLeft.Y == that.BottomRight.Y {
		// Check for shared segment on X-axis
		if (this.TopLeft.X <= that.BottomRight.X) && (that.TopLeft.X <= this.BottomRight.X) {
			return North
		}
	}

	// --- Horizontal Adjacency (Y-Alignment) ---

	// Case 3: 'this' is directly to the left of 'that' (that is to the East)
	if this.BottomRight.X == that.TopLeft.X {
		// Check for shared segment on Y-axis (vertical overlap)
		if (this.TopLeft.Y <= that.BottomRight.Y) && (that.TopLeft.Y <= this.BottomRight.Y) {
			return East
		}
	}

	// Case 4: 'this' is directly to the right of 'that' (that is to the West)
	if this.TopLeft.X == that.BottomRight.X {
		// Check for shared segment on Y-axis
		if (this.TopLeft.Y <= that.BottomRight.Y) && (that.TopLeft.Y <= this.BottomRight.Y) {
			return West
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
