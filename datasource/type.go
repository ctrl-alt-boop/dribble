package datasource

import (
	"context"
)

type Type interface {
	BaseType() Type
}

type DataSourceType int

func (t DataSourceType) BaseType() Type {
	return Type(SourceTypeUndefined)
}

//go:generate stringer -type=DataSourceType,SQLDialectType,NoSQLType,GraphType,TimeSeriesType -trimprefix=Type

const (
	SourceTypeSQL DataSourceType = iota
	SourceTypeNoSQL
	SourceTypeGraph
	SourceTypeTimeSeries

	// file?
	// url?
	// filesystem?

	SourceTypeUndefined DataSourceType = -1 // undefined
)

type SQLDialectType int

func (t SQLDialectType) BaseType() Type {
	return SourceTypeSQL
}

const (
	PostgreSQL SQLDialectType = iota // postgres
	MySQL                            // mysql
	SQLite3                          // sqlite3

	NumSupportedSQLDialects
)

type NoSQLType int

func (t NoSQLType) BaseType() Type {
	return SourceTypeNoSQL
}

const (
	MongoDB   NoSQLType = iota // mongo
	Firestore                  // firestore
	Redis                      // redis

	NumSupportedNoSQLModels
)

type GraphType int

func (t GraphType) BaseType() Type {
	return SourceTypeGraph
}

const (
	NumSupportedGraphModels GraphType = iota
)

type TimeSeriesType int

func (t TimeSeriesType) BaseType() Type {
	return SourceTypeTimeSeries
}

const (
	NumSupportedTimeSeries TimeSeriesType = iota
)

type (
	DataSource interface {
		Name() string
		GoName() string

		// Init(dsn Namer) error

		Open(context.Context) error
		Ping(context.Context) error
		Close(context.Context) error

		IsClosed() bool

		Request(ctx context.Context, req Request) (any, error)

		DataSourceType() SourceType

		// e.g. SQL, Postgres
		// // Or do I actually want Database, SQL, Postgres?
		// I'll temporarily do string, DST, string
		// Ooor, DataSourceType, ExecutorType, string?
		Path() []string
	}

	Executor interface {
		DataSource
		Name() string
		GoName() string
		ExecutorType() ExecutorType
	}

	Model interface {
		Executor
		Name() string
		GoName() string
		ModelType() ModelType
	}

	SourceType   string
	ExecutorType string
	ModelType    string
)
