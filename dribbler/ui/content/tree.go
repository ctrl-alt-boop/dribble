package content

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

//go:generate stringer -type=NodeType -trimprefix=NodeType

type NodeType int

const (
	NodeTypeServer NodeType = iota
	NodeTypeDatabase
	NodeTypeTable
	NodeTypeColumn

	NodeTypeSchema
	NodeTypeSchemaDatabase
	NodeTypeSchemaTable
	NodeTypeSchemaColumn

	NodeTypeProperties
	NodeTypePropery

	NodeTypeRequest
	NodeTypeIntent
	NodeTypeResponse

	NodeTypeCustom
)

var _ Content[[]string] = (*Tree)(nil)
var _ Selection = (*Tree)(nil)
var _ tea.Model = (*Tree)(nil)

type Tree struct {
	*Node
	SelectionPath []int

	Width, Height int
}

type Node struct {
	Index     int
	IndexPath []int
	Type      NodeType
	Name      string
	Children  []*Node
	Data      any
	Parent    *Node
	Expanded  bool
	Selected  bool

	// checked bool // maybe for multi selection
}

func (n *Node) ExpandAll() {
	for _, child := range n.Children {
		child.ExpandAll()
	}
	n.Expanded = true
}

func (n *Node) CollapseAll() {
	for _, child := range n.Children {
		child.CollapseAll()
	}
	n.Expanded = false
}

func (n *Node) GetAt(index int) *Node {
	return n.Children[index]
}

func (n *Node) GetAtPath(index []int) *Node {
	if len(index) == 0 {
		return n
	}
	return n.Children[index[0]].GetAtPath(index[1:])
}

func (n *Node) NewChild(name string, data any) *Node {
	newIndexPath := make([]int, len(n.IndexPath), len(n.IndexPath)+1)
	copy(newIndexPath, n.IndexPath)
	newNode := &Node{
		Index:     len(n.Children),
		IndexPath: append(newIndexPath, len(n.Children)),
		Name:      name,
		Children:  make([]*Node, 0),
		Data:      data,
		Parent:    n,
	}
	n.Children = append(n.Children, newNode)
	return newNode
}

// const indentSybmol = "--"
const indentSybmol = "\u258F "
const expandedSymbol = "\u2bc6"
const collapsedSymbol = "\u2bc7"

func (n *Node) View() string {
	var childrenStringBuilder strings.Builder
	if n.Expanded {
		for _, child := range n.Children {
			childrenStringBuilder.WriteString(strings.Repeat(indentSybmol, len(child.IndexPath)))
			childrenStringBuilder.WriteString(child.View())
		}
		return n.Name + expandedSymbol + "\n" + childrenStringBuilder.String()
	}
	return n.Name + collapsedSymbol + "\n" + childrenStringBuilder.String()
}

func NewTree() *Tree {
	return &Tree{
		Node: &Node{
			Index:     0,
			IndexPath: make([]int, 0),
			Name:      "root",
			Children:  make([]*Node, 0),
			Data:      nil,
			Parent:    nil,
			Expanded:  true,
		},
		SelectionPath: make([]int, 0),
	}
}

func (t *Tree) CollapseAll() {
	t.Node.CollapseAll()
}

func (t *Tree) ExpandAll() {
	t.Node.ExpandAll()
}

func (t *Tree) GetSelected() any {
	return t.Node.GetAtPath(t.SelectionPath)
}

func (t *Tree) GetAt(index int) *Node {
	return t.Node.GetAt(index)
}

func (t *Tree) GetAtPath(path []int) *Node {
	return t.Node.GetAtPath(path)
}

func (t *Tree) Data() any {
	return t.Node.Data
}

func (t *Tree) UpdateSize(width int, height int) {
	t.Width, t.Height = width, height
}

func (t *Tree) Get() []string {
	return strings.Split(t.Node.View(), "\n")
}

func (t *Tree) View() string {
	return t.Node.View()
}

// Init implements tea.Model.
func (t *Tree) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model.
func (t *Tree) Update(tea.Msg) (tea.Model, tea.Cmd) {
	return t, nil
}

func (t *Tree) Cursor() (int, int) {
	return t.CursorX(), t.CursorY()
}

func (t *Tree) CursorX() int {
	// X cursor is always at the start of the line in a tree view.
	// The indentation is handled by the View() method.
	return len(t.SelectionPath) * len(indentSybmol)
}

func (t *Tree) CursorY() int {
	return t.calculateCursorY()
}

func (t *Tree) MoveCursor(_ int, dY int) {
	if dY > 0 {
		t.MoveCursorDown(dY)
	} else if dY < 0 {
		t.MoveCursorUp(dY)
	}
}

func (t *Tree) SetCursor(_ int, _ int) {}

func (t *Tree) MoveCursorUp(y ...int) {
	delta := 1
	if len(y) > 0 {
		delta = y[0]
	}

	if len(t.SelectionPath) == 0 {
		return
	}

	// Check if we are moving to a parent or a sibling.
	isMovingToParent := t.SelectionPath[len(t.SelectionPath)-1] == 0

	if t.SelectionPath[len(t.SelectionPath)-1] > 0 {
		t.SelectionPath[len(t.SelectionPath)-1] -= delta
	} else {
		t.SelectionPath = t.SelectionPath[:len(t.SelectionPath)-1]
	}

	// If we moved to a previous sibling, and it's expanded, navigate to its last visible child.
	if !isMovingToParent && len(t.SelectionPath) > 0 {
		node := t.GetAtPath(t.SelectionPath)
		for node.Expanded && len(node.Children) > 0 {
			lastChildIndex := len(node.Children) - 1
			t.SelectionPath = append(t.SelectionPath, lastChildIndex)
			node = node.Children[lastChildIndex]
		}
	}
}

func (t *Tree) MoveCursorDown(y ...int) {
	delta := 1
	if len(y) > 0 {
		delta = y[0]
	}

	if len(t.SelectionPath) == 0 {
		if len(t.Node.Children) > 0 {
			t.SelectionPath = []int{0}
		}
		return
	}

	currentNode := t.GetAtPath(t.SelectionPath)
	if currentNode.Expanded && len(currentNode.Children) > 0 {
		t.SelectionPath = append(t.SelectionPath, 0)
		return
	}

	// Walk up the tree from the current selection
	tempPath := t.SelectionPath
	for len(tempPath) > 0 {
		parentPath := tempPath[:len(tempPath)-1]
		parent := t.GetAtPath(parentPath)
		currentIndex := tempPath[len(tempPath)-1]

		// Try to move to the next sibling
		if currentIndex+delta < len(parent.Children) {
			t.SelectionPath = append(parentPath, currentIndex+delta)
			return
		}

		// No sibling, go up one level and try again
		tempPath = parentPath
	}
}

func (t *Tree) MoveCursorLeft(x ...int) {
	if len(t.SelectionPath) == 0 {
		return
	}

	node := t.GetAtPath(t.SelectionPath)
	if node.Expanded {
		// If expanded, collapse it.
		node.Expanded = false
	} else if len(t.SelectionPath) > 0 {
		// If collapsed, move to parent.
		t.SelectionPath = t.SelectionPath[:len(t.SelectionPath)-1]
	}
}

func (t *Tree) MoveCursorRight(x ...int) {
	if len(t.SelectionPath) == 0 {
		return
	}

	node := t.GetAtPath(t.SelectionPath)
	if !node.Expanded && len(node.Children) > 0 {
		// If collapsed and has children, expand it.
		node.Expanded = true
	}
}

func (t *Tree) ToggleExpand() {
	if len(t.SelectionPath) == 0 {
		return
	}
	node := t.GetAtPath(t.SelectionPath)
	if len(node.Children) > 0 {
		node.Expanded = !node.Expanded
	}
}

// calculateCursorY traverses the visible tree to find the line number of the selection.
func (t *Tree) calculateCursorY() int {
	if len(t.SelectionPath) == 0 {
		return -1 // No selection
	}

	visibleNodes := make([]*Node, 0)
	var collectVisible func(*Node)
	collectVisible = func(n *Node) {
		if n != t.Node { // Don't include the invisible root node itself
			visibleNodes = append(visibleNodes, n)
		}
		if n.Expanded {
			for _, child := range n.Children {
				collectVisible(child)
			}
		}
	}

	// Start collection from the root's children
	for _, child := range t.Node.Children {
		collectVisible(child)
	}

	selectedNode := t.GetAtPath(t.SelectionPath)
	for i, node := range visibleNodes {
		if node == selectedNode {
			return i
		}
	}

	panic("Selected node not found in visible nodes")
}
