package ui

import (
	"cmp"
	"slices"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/list"
	"github.com/charmbracelet/lipgloss/tree"
	"github.com/ctrl-alt-boop/gooldb/dribble/config"
	"github.com/ctrl-alt-boop/gooldb/pkg/connection"
)

type (
	ConnectionItem struct {
		*connection.Settings
		Name string
	}
)

func (item ConnectionItem) FilterValue() string { return "" }
func (item ConnectionItem) Title() string       { return item.Name }
func (item ConnectionItem) Description() string { return "" }
func (item ConnectionItem) Inspect() string     { return item.AsString() }

func GetSavedConfigsSorted() []*ConnectionItem {
	items := make([]*ConnectionItem, 0, len(config.SavedConfigs))
	for name, settings := range config.SavedConfigs {
		items = append(items, &ConnectionItem{
			Name:     name,
			Settings: settings,
		})
	}

	sortFunc := func(left, right *ConnectionItem) int {
		return cmp.Compare(left.Name, right.Name)
	}
	slices.SortFunc(items, sortFunc)
	logger.Infof("GetSavedConfigsSorted: %+v", items)
	return items
}

func SettingsToConnectionItems(settings []*connection.Settings) []*ConnectionItem {
	listString := ""
	for _, setting := range settings {
		listString += setting.AsString() + "\n"
	}
	logger.Infof("SettingsToConnectionItems: %+v", listString)
	items := make([]*ConnectionItem, 0, len(settings))
	for _, setting := range settings {
		if setting.SettingsName == "" {
			setting.SettingsName = setting.DbName
		}
		items = append(items, &ConnectionItem{
			Name:     setting.SettingsName,
			Settings: setting,
		})
	}
	return items
}

func CreateNestedList() *list.List {
	types := map[connection.Type]*list.List{
		connection.Driver:      list.New(),
		connection.Server:      list.New(),
		connection.Database:    list.New(),
		connection.Table:       list.New(),
		connection.TypeUnknown: list.New(),
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
		Type     connection.Type
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

func NewConnectionNode(nodeType connection.Type, connectionItem *ConnectionItem) *TreeNode {
	var name string
	switch nodeType {
	case connection.Driver:
		name = connectionItem.DriverName
	case connection.Server:
		name = connectionItem.Ip
	case connection.Database:
		name = connectionItem.DbName
	case connection.Table:
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
	return t.Item.AsString()
	// return t.Name
}

// Value implements tree.Node.
func (t *TreeNode) Value() string {
	return t.Name
}
