package dribble

import (
	"github.com/ctrl-alt-boop/dribble/database"
	"github.com/ctrl-alt-boop/dribble/result"
)

//go:generate stringer -type=EventType

type EventType uint

const (
	Connected    EventType = iota //EventType = "ConnectionSuccess"
	ConnectError                  //EventType = "ConnectionError"

	Reconnected    //EventType = "ReconnectSuccess"
	ReconnectError //EventType = "ReconnectError"

	Disconnected    //EventType = "DisconnectSuccess"
	DisconnectError //EventType = "DisconnectError"

	DriverLoadError //EventType = "DriverLoadError"

	TargetOpened    //EventType = "TargetOpened"
	TargetOpenError //EventType = "TargetOpenError"

	TargetClosed     //EventType = "TargetClosed"
	TargetCloseError //EventType = "TargetCloseError"

	TargetUpdated     //EventType = "TargetUpdated"
	TargetUpdateError //EventType = "TargetUpdateError"

	DBOpened    //EventType = "DatabaseConnectSuccess"
	DBOpenError //EventType = "DatabaseConnectError"

	DatabaseListFetched    //EventType = "DatabaseListFetchSuccess"
	DatabaseListFetchError //EventType = "DatabaseListFetchError"

	DBTableListFetched    //EventType = "TableListFetchSuccess"
	DBTableListFetchError //EventType = "TableListFetchError"

	TableSelected    //EventType = "TableSelectSuccess"
	TableSelectError //EventType = "TableSelectError"

	TableFetched    //EventType = "TableFetchSuccess"
	TableFetchError //EventType = "TableFetchError"

	TableCountFetched

	QueryExecuted
	QueryExecuteError
)

type (
	EventHandler func(any, error)

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

	// QueryData struct {
	// 	Query Query
	// 	Value any
	// 	List  []any
	// 	Table *result.Table
	// }
)
