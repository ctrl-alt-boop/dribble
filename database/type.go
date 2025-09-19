package database

func createBaseTree() *DBTypeNode {
	return &DBTypeNode{
		Name: "",
		Children: []*DBTypeNode{
			{
				Name:     "SQL",
				Children: []*DBTypeNode{},
			},
			{
				Name:     "NoSQL",
				Children: []*DBTypeNode{},
			},
		},
	}
}

type DBTypeNode struct {
	Name     string
	TrueName string

	Properties map[string]string
	Children   []*DBTypeNode
}

func (n *DBTypeNode) Register(path ...string) {
	if len(path) == 0 {
		return
	}
	if n.Children == nil {
		n.Children = make([]*DBTypeNode, 0)
	}
	for _, child := range n.Children {
		if child.Name == path[0] {
			child.Register(path[1:]...)
			return
		}
	}
	newNode := &DBTypeNode{
		Name: path[0],
	}
	n.Children = append(n.Children, newNode)
	newNode.Register(path[1:]...)
}
