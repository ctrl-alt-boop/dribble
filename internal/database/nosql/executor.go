package nosql

import (
	"context"
	"fmt"

	"github.com/ctrl-alt-boop/dribble/database"
)

const (
	MongoDB = "mongo"
)

var SupportedDrivers []string = []string{
	MongoDB,
}

var Defaults = map[string]*database.Target{
	MongoDB: {
		Type:       database.DBDriver,
		DriverName: MongoDB,
		Ip:         "127.0.0.1",
		Port:       27017,
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

func CreateDriverFromTarget(target *database.Target) (database.NoSQLClient, error) {
	switch target.DriverName {
	// case MongoDB:
	// return mongodb.NewMongoDBDriver(target)
	default:
		return nil, fmt.Errorf("unknown or unsupported driver: %s", target.DriverName)
	}
}

var _ database.Executor = &Executor{}

type (
	IntentHandler func(intent *database.Intent, err error)
	ResultHandler func(result any, err error)
)

type Executor struct {
	client database.NoSQLClient
	target *database.Target
	// driver database.Driver

	onBefore IntentHandler
	onAfter  IntentHandler
	onResult ResultHandler
}

func NewExecutor(target *database.Target) *Executor {
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
func (e *Executor) Driver() database.Driver {
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
func (e *Executor) SetDriver(driver database.Driver) {

}

// SetTarget implements database.Executor.
func (e *Executor) SetTarget(target *database.Target) {
	e.target = target
}

// Target implements database.Executor.
func (e *Executor) Target() *database.Target {
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
