package layout

import tea "github.com/charmbracelet/bubbletea"

type Manager interface {
	SetDefinition(definition RenderDefinition)
	GetDefinition() RenderDefinition

	AddLayout(definition LayoutDefinition)

	UpdateLayout(index int, opts ...LayoutOption)
	SetLayout(index int, definition LayoutDefinition)
	GetLayout(index int) LayoutDefinition

	// If position is not set, returns empty LayoutDefinition
	GetLayoutForPosition(position Position) LayoutDefinition

	Layout(models []tea.Model) []tea.Model
	View(models []tea.Model) string

	SetSize(width, height int)

	SetFocusPassThrough(v bool)
	GetFocusPassThrough() bool
}

type managerBase struct {
	focusPassThrough bool
	Width, Height    int
	X, Y             int

	renderDefinition RenderDefinition
}

func (b *managerBase) SetDefinition(definition RenderDefinition) {
	b.renderDefinition = definition
}

func (b *managerBase) AddLayout(definition LayoutDefinition) {
	b.renderDefinition.Definitions = append(b.renderDefinition.Definitions, definition)
	if _, ok := b.renderDefinition.indexForPosition[definition.Position]; ok {
		panic("Position already in use") // FIXME: temp panic
	}
	b.renderDefinition.indexForPosition[definition.Position] = len(b.renderDefinition.Definitions) - 1
}

func (b *managerBase) GetDefinition() RenderDefinition {
	return b.renderDefinition
}

func (b *managerBase) GetLayout(index int) LayoutDefinition {
	return b.renderDefinition.Definitions[index]
}

// If position is not set, returns empty LayoutDefinition
func (b *managerBase) GetLayoutForPosition(position Position) LayoutDefinition {
	if index, ok := b.renderDefinition.indexForPosition[position]; ok {
		return b.renderDefinition.Definitions[index]
	}
	return LayoutDefinition{}
}

func (b *managerBase) SetLayout(index int, definition LayoutDefinition) {
	b.renderDefinition.Definitions[index] = definition
}

func (b *managerBase) UpdateLayout(index int, opts ...LayoutOption) {
	for _, opt := range opts {
		opt(&b.renderDefinition.Definitions[index])
	}
}

// GetFocusPassThrough implements Manager.
func (b *managerBase) GetFocusPassThrough() bool {
	return b.focusPassThrough
}

// SetFocusPassThrough implements Manager.
func (b *managerBase) SetFocusPassThrough(v bool) {
	b.focusPassThrough = v
}
