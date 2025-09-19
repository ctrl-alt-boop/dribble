package database_test

import (
	"testing"

	"github.com/ctrl-alt-boop/dribble/database"
)

func TestTypeNode(t *testing.T) {
	root := &database.DBTypeNode{
		Name: "",
		Children: []*database.DBTypeNode{
			{
				Name:     "SQL",
				Children: []*database.DBTypeNode{},
			},
			{
				Name:     "NoSQL",
				Children: []*database.DBTypeNode{},
			},
		},
	}

	root.Register("SQL", "MySQL")
	root.Register("Graph", "Neo4j")

	t.Logf("%+v", root)
}
