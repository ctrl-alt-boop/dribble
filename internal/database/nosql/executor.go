package nosql

import (
	"context"
	"fmt"

	"github.com/ctrl-alt-boop/dribble/database"
	"github.com/ctrl-alt-boop/dribble/target"
)

const (
	MongoDB = "mongo"
)

// var SupportedDrivers []string = []string{
// 	MongoDB,
// }

var SupportedDialects []target.Dialect = []target.Dialect{
	target.MongoDB,
}

var Defaults = map[string]*target.Target{
	MongoDB: {
		Name: "mongo",
		Type: target.TypeDriver,
		Properties: target.Properties{
			Dialect: target.MongoDB,
			Ip:      "127.0.0.1",
			Port:    27017,
		},
	},
}

type Method string

func (s Method) String() string {
	return string(s)
}

const (
	MethodFind Method = "FIND"
)

const DefaultFindLimit int = 10 // Just a safeguard

var NoSQLMethods = []Method{MethodFind}

func CreateDriverFromTarget(target *target.Target) (database.NoSQL, error) {
	switch target.Dialect {
	// case MongoDB:
	// return mongodb.NewMongoDBDriver(target)
	default:
		return nil, fmt.Errorf("unknown or unsupported driver: %s", target.Dialect)
	}
}

var _ database.Executor = &Executor{}

type (
	IntentHandler func(intent *database.Intent, err error)
	ResultHandler func(result any, err error)
)

type Executor struct {
	client database.NoSQL
	target *target.Target
	// driver database.Driver

	onBefore IntentHandler
	onAfter  IntentHandler
	onResult ResultHandler
}

func NewExecutor(target *target.Target) *Executor {
	client, err := CreateDriverFromTarget(target)
	if err != nil {
		panic(err)
	}
	return &Executor{
		target: target,
		client: client,
	}
}

func (e *Executor) Open(ctx context.Context) error {

	return nil
}

func (e *Executor) Close(_ context.Context) error {
	return nil
}

// Driver implements database.Executor.
func (e *Executor) Driver() database.SQL {
	return nil
}

// Execute implements database.Executor.
func (e *Executor) Execute(ctx context.Context, intent *database.Intent) (any, error) {
	panic("unimplemented")
}

// ExecutePrefab implements database.Executor.
func (e *Executor) ExecutePrefab(ctx context.Context, prefabType database.PrefabType, args ...any) (any, error) {
	panic("unimplemented")
}

// ExecuteAndHandle implements database.Executor.
func (e *Executor) ExecuteWithHandler(ctx context.Context, intent *database.Intent, handler func(result any, err error)) {
	panic("unimplemented")
}

// ExecuteWithChannel implements database.Executor.
func (e *Executor) ExecuteWithChannel(ctx context.Context, intent *database.Intent, eventChannel chan any) {
	panic("unimplemented")
}

// Ping implements database.Executor.
func (e *Executor) Ping(ctx context.Context) error {
	panic("unimplemented")
}

// SetDriver implements database.Executor.
func (e *Executor) SetDriver(driver database.SQL) {

}

// SetTarget implements database.Executor.
func (e *Executor) SetTarget(target *target.Target) {
	e.target = target
}

// Target implements database.Executor.
func (e *Executor) Target() *target.Target {
	return e.target
}

// OnResult implements database.Executor.
func (e *Executor) OnBefore(f func(intent *database.Intent, err error)) {
	e.onBefore = f
}

// OnResult implements database.Executor.
func (e *Executor) OnAfter(f func(intent *database.Intent, err error)) {
	e.onAfter = f
}

// OnResult implements database.Executor.
func (e *Executor) OnResult(f func(result any, err error)) {
	e.onResult = f
}

// SetEventChannel implements database.Executor.
func (e *Executor) SetEventChannel(...chan any) {
	panic("unimplemented")
}

// EventChannel implements database.Executor.
func (e *Executor) EventChannel() []chan any {
	panic("unimplemented")
}
