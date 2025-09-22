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

type Type interface {
	BaseType() Type
}

type DBType int

func (t DBType) BaseType() Type {
	return Type(TypeUndefined)
}

//go:generate stringer -type=DBType,SQLDialectType,NoSQLModelType,GraphType,TimeSeriesType -trimprefix=Type

const (
	TypeSQL DBType = iota
	TypeNoSQL
	TypeGraph
	TypeTimeSeries

	TypeUndefined DBType = -1 // undefined
)

type SQLDialectType int

func (t SQLDialectType) BaseType() Type {
	return TypeSQL
}

const (
	PostgreSQL SQLDialectType = iota // postgres
	MySQL                            // mysql
	SQLite3                          // sqlite3

	NumSupportedSQLDialects
)

type NoSQLModelType int

func (t NoSQLModelType) BaseType() Type {
	return TypeNoSQL
}

const (
	MongoDB   NoSQLModelType = iota // mongo
	Firestore                       // firestore
	Redis                           // redis

	NumSupportedNoSQLModels
)

type GraphType int

func (t GraphType) BaseType() Type {
	return TypeGraph
}

const (
	NumSupportedGraphModels GraphType = iota
)

type TimeSeriesType int

func (t TimeSeriesType) BaseType() Type {
	return TypeTimeSeries
}

const (
	NumSupportedTimeSeries TimeSeriesType = iota
)
