package dribble

import (
	"github.com/ctrl-alt-boop/dribble/request"
	"github.com/ctrl-alt-boop/dribble/result"
	"github.com/ctrl-alt-boop/dribble/target"
)

type (
	EventHandler func(eventType request.Status, args any, err error)

	DatabaseListFetchData struct {
		Driver    string
		Databases []*target.Target
	}

	TableListFetchData struct {
		Database string
		Tables   []string
	}

	TableFetchData struct {
		TableName string
		Table     *result.Table
	}

	TableCountFetchData struct {
		TableName string
		Counts    map[string]int
		Errors    map[string]error
	}

	ExecuteQuery struct {
		Driver   string
		Database string
		Table    string
		Query    Query
	}

	QueryResult struct {
		Query Query
		Value any
		List  []any
		Table *result.Table
	}
)
