package database

import (
	"context"
	"reflect"
)

type OperationType int

const (
	Read OperationType = iota
	Create
	Update
	Delete
	Execute
	// Meta?

	NoOp OperationType = -1
)

var OperationTypes = []OperationType{
	Read,
	Create,
	Update,
	Delete,
	Execute,
}

type PrefabType int

const (
	PrefabCurrentDatabase PrefabType = iota
	PrefabDatabases
	PrefabTables
	PrefabColumns
)

type Capabilities string

const (
	SupportsJSON    Capabilities = "json"
	SupportsJSONB   Capabilities = "jsonb"
	IsFile          Capabilities = "file"
	SupportsSQLLike Capabilities = "sql"
	SupportsBSON    Capabilities = "json"
)

type DatabaseType int

const (
	SQL DatabaseType = iota
	NoSQL
)

type (
	Dialect interface {
		GetTemplate(operationType OperationType) string
		GetPrefab(prefabType PrefabType) (string, bool)

		Capabilities() []Capabilities

		RenderPlaceholder(index int) string
		IncreamentPlaceholder() string

		RenderTypeCast() string
		RenderCurrentTimestamp() string
		RenderValue(value any) string
		QuoteRune() rune
		Quote(value string) string

		ResolveType(dbType string, value []byte) (any, error)
	}

	Driver interface {
		Dialect() Dialect

		ConnectionString(target *Target) string
		RenderIntent(intent *Intent) (string, error)
	}

	NoSQLClient interface {
		Dialect() Dialect
		Open(ctx context.Context, target *Target) error
		Ping(ctx context.Context) error
		Close(ctx context.Context) error

		Read(any)
		ReadMany(any)
		Create(any)
		Update(any)
		Delete(any)
		// Execute()
	}

	Executor interface {
		SetTarget(target *Target)
		Target() *Target

		SetDriver(driver Driver)
		Driver() Driver

		Open(ctx context.Context) error
		Ping(ctx context.Context) error
		Close(ctx context.Context) error

		Execute(ctx context.Context, intent *Intent) (any, error)
		ExecutePrefab(ctx context.Context, prefabType PrefabType, args ...any) (any, error)
		ExecuteWithHandler(ctx context.Context, intent *Intent, handler func(result any, err error))
		ExecuteWithChannel(ctx context.Context, intent *Intent, eventChannel chan any)
	}

	Intent struct {
		Target *Target

		Type OperationType

		OperationKind reflect.Kind
		Operation     any

		Args []any
	}

	IntentBatch []*Intent
)

func NewReadIntent(target *Target, operation any, args ...any) *Intent {
	return &Intent{
		Target:        target,
		Type:          Read,
		OperationKind: reflect.TypeOf(operation).Kind(),
		Operation:     operation,
		Args:          args,
	}
}

func NewCreateIntent(target *Target, operation any, args ...any) *Intent {
	return &Intent{
		Target:        target,
		Type:          Create,
		OperationKind: reflect.TypeOf(operation).Kind(),
		Operation:     operation,
		Args:          args,
	}
}

func NewUpdateIntent(target *Target, operation any, args ...any) *Intent {
	return &Intent{
		Target:        target,
		Type:          Update,
		OperationKind: reflect.TypeOf(operation).Kind(),
		Operation:     operation,
		Args:          args,
	}
}

func NewDeleteIntent(target *Target, operation any, args ...any) *Intent {
	return &Intent{
		Target:        target,
		Type:          Delete,
		OperationKind: reflect.TypeOf(operation).Kind(),
		Operation:     operation,
		Args:          args,
	}
}
