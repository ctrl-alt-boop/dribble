package memory

import "fmt"

//go:generate stringer -type=NodeType -trimprefix=NodeType

type treeStore struct {
	*Node
	SelectionPath []int

	Width, Height int
}

type Node struct {
	Item

	Index     int
	IndexPath []int

	Children []*Node
	Parent   *Node
}

func newTree() *treeStore {
	return &treeStore{
		Node: &Node{
			Item: Item{
				name:  "root",
				value: nil,
			},
			Index:     0,
			IndexPath: make([]int, 0),
			Children:  make([]*Node, 0),
			Parent:    nil,
		},
		SelectionPath: make([]int, 0),
	}
}

func (n *Node) Get(index int) *Node {
	return n.Children[index]
}

func (n *Node) GetAt(index []int) *Node {
	if len(index) == 0 {
		return n
	}
	return n.Children[index[0]].GetAt(index[1:])
}

func (n *Node) AddChild(data any) *Node {
	newIndexPath := make([]int, len(n.IndexPath), len(n.IndexPath)+1)
	copy(newIndexPath, n.IndexPath)

	var name string
	if _, ok := data.(interface{ Name() string }); ok {
		name = data.(interface{ Name() string }).Name()
	} else {
		name = fmt.Sprintf("%T-%d", data, len(n.Children))
	}

	newNode := &Node{
		Item: Item{
			name:  name,
			value: data,
		},
		Index:     len(n.Children),
		IndexPath: append(newIndexPath, len(n.Children)),
		Children:  make([]*Node, 0),
		Parent:    n,
	}
	n.Children = append(n.Children, newNode)
	return newNode
}

func (t *treeStore) GetAt(index int) *Node {
	return t.Node.Get(index)
}

func (t *treeStore) GetAtPath(path []int) *Node {
	return t.Node.GetAt(path)
}
