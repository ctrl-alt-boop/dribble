package ui

import (
	"cmp"
	"slices"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/list"
	"github.com/charmbracelet/lipgloss/tree"
	"github.com/ctrl-alt-boop/dribble/database"
	"github.com/ctrl-alt-boop/dribbler/config"
)

type (
	ConnectionItem struct {
		*database.Target
		Name string
	}
)

func (item ConnectionItem) FilterValue() string { return "" }
func (item ConnectionItem) Title() string       { return item.Name }
func (item ConnectionItem) Description() string { return "" }
func (item ConnectionItem) Inspect() string     { return item.Name }

func GetSavedConfigsSorted() []*ConnectionItem {
	items := make([]*ConnectionItem, 0, len(config.SavedConfigs))
	for name, settings := range config.SavedConfigs {
		items = append(items, &ConnectionItem{
			Name:   name,
			Target: settings,
		})
	}

	sortFunc := func(left, right *ConnectionItem) int {
		return cmp.Compare(left.Name, right.Name)
	}
	slices.SortFunc(items, sortFunc)
	logger.Infof("GetSavedConfigsSorted: %+v", items)
	return items
}

func SettingsToConnectionItems(targets []*database.Target) []*ConnectionItem {
	listString := ""
	for _, target := range targets {
		listString += target.Name + "\n"
	}
	logger.Infof("SettingsToConnectionItems: %+v", listString)
	items := make([]*ConnectionItem, 0, len(targets))
	for _, target := range targets {
		if target.Name == "" {
			target.Name = target.DBName
		}
		items = append(items, &ConnectionItem{
			Name:   target.Name,
			Target: target,
		})
	}
	return items
}

func CreateNestedList() *list.List {
	types := map[database.TargetType]*list.List{
		database.DBDriver: list.New(),
		database.DBServer: list.New(),
		database.Database: list.New(),
		database.DBTable:  list.New(),
		database.Unknown:  list.New(),
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
		Type     database.TargetType
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

func NewCategoryNode(categoryName string, children TreeNodeChildren) *TreeNode {
	return &TreeNode{
		Item:     nil,
		Type:     "Category",
		Name:     categoryName,
		children: children,
		hidden:   false,
	}
}

func NewConnectionNode(nodeType database.TargetType, connectionItem *ConnectionItem) *TreeNode {
	var name string
	switch nodeType {
	case database.DBDriver:
		name = connectionItem.DriverName
	case database.DBServer:
		name = connectionItem.Ip
	case database.Database:
		name = connectionItem.DBName
	case database.DBTable:
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
