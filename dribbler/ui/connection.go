package ui

import (
	"cmp"
	"fmt"
	"slices"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/list"
	"github.com/charmbracelet/lipgloss/tree"
	"github.com/ctrl-alt-boop/dribble/database"
	"github.com/ctrl-alt-boop/dribble/target"
	"github.com/ctrl-alt-boop/dribbler/config"
	"github.com/ctrl-alt-boop/dribbler/logging"
)

type (
	ConnectionItem struct {
		Name string
		Type target.Type
		DSN  database.DataSourceNamer
	}
)

func (item ConnectionItem) FilterValue() string { return "" }
func (item ConnectionItem) Title() string       { return item.Name }
func (item ConnectionItem) Description() string { return "" }
func (item ConnectionItem) Inspect() string     { return item.Name }

func GetSavedConfigsSorted() []*ConnectionItem {
	items := make([]*ConnectionItem, 0, len(config.SavedConfigs))
	for name, dataSourceNamer := range config.SavedConfigs {
		items = append(items, &ConnectionItem{
			Name: name,
			DSN:  dataSourceNamer,
		})
	}

	sortFunc := func(left, right *ConnectionItem) int {
		return cmp.Compare(left.Name, right.Name)
	}
	slices.SortFunc(items, sortFunc)
	logging.GlobalLogger().Infof("GetSavedConfigsSorted: %+v", items)
	return items
}

func SettingsToConnectionItems(targets []*target.Target) []*ConnectionItem {
	listString := ""
	for _, target := range targets {
		listString += target.Name + "\n"
	}
	logging.GlobalLogger().Infof("SettingsToConnectionItems: %+v", listString)
	items := make([]*ConnectionItem, 0, len(targets))
	for i, target := range targets {
		name := target.Name
		if target.Name == "" {
			name = "Unnamed Connection " + fmt.Sprint(i)
		}
		items = append(items, &ConnectionItem{
			Name: name,
			Type: target.Type,
		})
	}
	return items
}

func CreateNestedList() *list.List {
	types := map[target.Type]*list.List{
		target.TypeDriver:   list.New(),
		target.TypeServer:   list.New(),
		target.TypeDatabase: list.New(),
		target.TypeTable:    list.New(),
		target.TypeUnknown:  list.New(),
	}
	configs := GetSavedConfigsSorted()
	for _, item := range configs {
		types[item.Type] = types[item.Type].Item(item)
	}
	nested := list.New()
	for key, value := range types {
		nested.Items(key, value)
	}
	return nested
}

type (
	Tree struct {
		*tree.Tree
	}

	TreeNodeChildren []*TreeNode

	TreeNode struct {
		Item     *ConnectionItem
		Type     target.Type
		Name     string
		children TreeNodeChildren
		hidden   bool
	}
)

// At implements tree.Children.
func (t TreeNodeChildren) At(index int) tree.Node {
	return t[index]
}

// Length implements tree.Children.
func (t TreeNodeChildren) Length() int {
	return len(t)
}

func NewTree() *Tree {
	tree := tree.Root("")
	tree.RootStyle(lipgloss.NewStyle())
	tree.ItemStyle(lipgloss.NewStyle())
	return &Tree{
		Tree: tree,
	}
}

func NewCategoryNode(categoryType target.Type, children TreeNodeChildren) *TreeNode {
	return &TreeNode{
		Item:     nil,
		Type:     categoryType,
		Name:     categoryType.String(),
		children: children,
		hidden:   false,
	}
}

func NewConnectionNode(nodeType target.Type, connectionItem *ConnectionItem) *TreeNode {
	var name string
	switch nodeType {
	case target.TypeDriver:
		name = connectionItem.Type.String()
	case target.TypeServer:
		name = connectionItem.DSN.DSN()
	case target.TypeDatabase:
		name = connectionItem.DSN.DSN()
	case target.TypeTable:
		name = "{table}"
	default:
		name = "ERR"
	}
	return &TreeNode{
		Item:     connectionItem,
		Type:     nodeType,
		Name:     name,
		children: nil,
		hidden:   false,
	}
}

// Children implements tree.Node.
func (t *TreeNode) Children() tree.Children {
	return t.children
}

// Hidden implements tree.Node.
func (t *TreeNode) Hidden() bool {
	return false
}

// SetHidden implements tree.Node.
func (t *TreeNode) SetHidden(value bool) {
	t.hidden = value
}

// SetValue implements tree.Node.
func (t *TreeNode) SetValue(value any) {
	// t.ConnectionItem = value.(*ConnectionItem)
}

func (t *TreeNode) SetConnectionItem(connectionItem *ConnectionItem) {
	t.Item = connectionItem
}

// String implements tree.Node.
func (t *TreeNode) String() string {
	if t.Item == nil {
		return t.Name
	}
	return t.Item.Name
	// return t.Name
}

// Value implements tree.Node.
func (t *TreeNode) Value() string {
	return t.Name
}
