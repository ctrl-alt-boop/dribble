package layout

import (
	lipgloss "charm.land/lipgloss/v2"
	"github.com/ctrl-alt-boop/dribbler/internal/panel"
)

var _ panel.Composer = (*TabbedLayoutComposer)(nil)

// TabbedLayoutComposer composes panels in a tabbed layout.
type TabbedLayoutComposer struct {
	panel.ComposerBase

	definitions []panel.Definition

	TabStyle lipgloss.Style
	TabsSide panel.Position
	// TabLabels     []string
	contentHeight int
}

// NewTabbedLayout creates a new TabbedLayoutComposer.
func NewTabbedLayout(tabsSide panel.Position, panels panel.DefinitionList, opts ...panel.ComposerOption) *TabbedLayoutComposer {
	tabBorder := lipgloss.NormalBorder()
	composer := &TabbedLayoutComposer{
		definitions: panels,
		TabStyle:    lipgloss.NewStyle().Border(tabBorder, true).BorderForeground(lipgloss.Color("63")).Padding(0, 1),
		TabsSide:    tabsSide,
	}

	for _, opt := range opts {
		opt(composer)
	}

	return composer
}

// Compose implements panel.Composer.
func (t *TabbedLayoutComposer) Compose(width int, height int) *panel.Composition {
	// compose tabs panel
	tabPanel := t.ComposeTabPanel(width, height)
	contentWidth, contentHeight := width, height
	contentX, contentY := 0, 0
	switch t.TabsSide {
	case panel.Top, panel.Bottom:
		contentHeight = height - tabPanel.GetHeight()
		contentY = tabPanel.GetHeight()
	case panel.Left, panel.Right:
		contentWidth = width - tabPanel.GetWidth()
		contentX = tabPanel.GetWidth()
	}

	panelState := &panel.State{
		Width:  contentWidth,
		Height: contentHeight,
		X:      contentX,
		Y:      contentY,
	}
	composition := t.BuildComposition(panelState)

	// add tabs panel
	composition.PreRendered = append(composition.PreRendered, tabPanel)

	return composition
}

// ComposeTabPanel composes the tab panel.
func (t *TabbedLayoutComposer) ComposeTabPanel(width int, height int) *lipgloss.Layer {
	var tabLabels []string
	var longestLabel int

	for i, def := range t.definitions {

		var tabLabel string
		if def.Name == "" {
			tabLabel = "Tab " + string(rune('A'+i))
		} else {
			tabLabel = def.Name
		}
		if len(tabLabel) > longestLabel {
			longestLabel = len(tabLabel)
		}

		tabLabels = append(tabLabels, tabLabel)
	}
	var x, y int
	var panelWidth, panelHeight int
	switch t.TabsSide {
	case panel.Top:
		x, y = 0, 0
		panelHeight = 3
	case panel.Bottom:
		x, y = 0, height-3
		panelHeight = 3
	case panel.Left:
		x, y = 0, 0
		panelWidth = longestLabel
		panelHeight = height
	case panel.Right:
		x, y = width-3, 0
		panelWidth = longestLabel
		panelHeight = height
	}

	tabPanel := lipgloss.NewCanvas()

	var offset int
	for _, label := range tabLabels {
		tabLayer := lipgloss.NewLayer(label).X(offset).Y(0) // FIXME: fix vertical tab panel
		offset += lipgloss.Width(label)
		tabPanel.AddLayers(tabLayer)
	}
	return lipgloss.NewLayer(tabPanel.Render()).Width(panelWidth).Height(panelHeight).X(x).Y(y)
}
