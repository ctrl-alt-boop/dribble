package gooldb

import (
	"github.com/ctrl-alt-boop/gooldb/pkg/data"
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
		Databases []string
	}

	DBTableListFetchData struct {
		Database string
		Tables   []string
	}

	TableFetchData struct {
		TableName string
		Table     data.Table
	}

	TableCountFetchData struct {
		TableName string
		Counts    map[string]int
		Errors    map[string]error
	}

	Query struct {
		Driver   string
		Database string
		Table    string
		Query    string
	}

	QueryData struct {
		Query Query
		Data  data.Table
	}
)
