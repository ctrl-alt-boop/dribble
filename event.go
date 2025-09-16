package dribble

import (
	"github.com/ctrl-alt-boop/dribble/database"
	"github.com/ctrl-alt-boop/dribble/result"
)

//go:generate stringer -type=Success
//go:generate stringer -type=Error

type (
	Status      uint
	Success     Status
	Error       Status
	BatchStatus []Status
)

const (
	SuccessConnect Success = iota
	SuccessReconnect
	SuccessDisconnect
	SuccessTargetOpen
	SuccessTargetUpdate
	SuccessTargetClose

	SuccessReadDatabaseSchema
	SuccessReadTableSchema
	SuccessReadColumnSchema

	SuccessReadDatabaseProperties
	SuccessReadTableProperties
	SuccessReadColumnProperties

	SuccessReadDatabaseList
	SuccessReadDBTableList
	SuccessReadDBColumnList

	SuccessReadCount

	SuccessReadTable

	SuccessExecute
	SuccessBatchExecute
)

const (
	ErrorConnect Error = iota
	ErrorReconnect
	ErrorDisconnect
	ErrorTargetOpen
	ErrorTargetClose
	ErrorTargetUpdate

	ErrorReadDatabaseSchema
	ErrorReadTableSchema
	ErrorReadColumnSchema

	ErrorReadDatabaseProperties
	ErrorReadTableProperties
	ErrorReadColumnProperties

	ErrorReadDatabaseList
	ErrorReadDBTableList
	ErrorReadDBColumnList

	ErrorReadCount

	ErrorReadTable

	ErrorExecute
	ErrorBatchExecute
)

type (
	EventHandler func(eventType Status, args any, err error)

	DatabaseListFetchData struct {
		Driver    string
		Databases []*database.Target
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
