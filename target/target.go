package target

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"

	"github.com/ctrl-alt-boop/dribble/database"
	internal "github.com/ctrl-alt-boop/dribble/internal/database"
	"github.com/ctrl-alt-boop/dribble/request"
)

//go:generate stringer -type=Type -output=target_type_string.go
type Type int

const (
	TypeUnknown Type = iota - 1
	TypeDriver
	TypeServer
	TypeDatabase
	TableTable
)

var ErrNoRequests = errors.New("no requests provided")

type Target struct {
	Name       string
	Type       Type
	DBType     database.Type
	Properties database.ConnectionProperties

	executor database.Database
	// mu               sync.Mutex
	nextRequestID    atomic.Int64
	pendingResponses map[int64]chan database.Response // TODO: Check if necessary
}

func New(name string, targetType Type, dbType database.Type, options ...Option) *Target {
	target := &Target{
		Name:   name,
		Type:   targetType,
		DBType: dbType,
		Properties: database.ConnectionProperties{
			Addr:       "localhost",
			Port:       0,
			DBName:     "",
			Username:   "",
			Password:   "",
			Additional: make(map[string]string),
		},
		pendingResponses: make(map[int64]chan database.Response),
	}

	for _, option := range options {
		option(target)
	}

	executor, err := internal.CreateClientForType(target.DBType)
	if err != nil {
		panic(err)
	}
	executor.SetConnectionProperties(target.Properties)
	target.executor = executor

	return target
}

// func (t *Target) Initialize() error {

// 	executor, err := internal.CreateClientForType(t.Properties.Type)
// 	if err != nil {
// 		return err
// 	}
// 	t.database = executor

// 	return nil
// }

func (t *Target) Open(ctx context.Context) error {
	return t.executor.Open(ctx)
}

func (t *Target) Ping(ctx context.Context) error {
	return t.executor.Ping(ctx)
}

func (t *Target) Close(ctx context.Context) error {
	return t.executor.Close(ctx)
}

func (t *Target) Request(ctx context.Context, requests ...database.Request) (chan database.Response, error) {
	numRequests := len(requests)
	if numRequests == 0 {
		return nil, ErrNoRequests
	}

	requestID := t.nextRequestID.Add(1)
	resultChan := make(chan database.Response, numRequests)

	go func() {
		defer close(resultChan)

		var wg sync.WaitGroup
		wg.Add(numRequests)

		for _, req := range requests {
			go func(r database.Request) {
				defer wg.Done()

				requestResult, err := t.executor.Request(ctx, r)
				resultChan <- &request.Response{
					RequestID: requestID,
					Status:    request.Status(r.ResponseOnSuccess().Code()),
					Body:      requestResult,
					Error:     err,
				}
			}(req)
		}
		wg.Wait()
	}()

	return resultChan, nil
}

// Blocks
func (t *Target) RequestWithHandler(ctx context.Context, handler database.ResponseHandler, requests ...database.Request) error { // TODO
	resultChan, err := t.Request(ctx, requests...)
	if err != nil {
		return err
	}

	for result := range resultChan {
		handler(result, nil)
	}

	return nil
}
