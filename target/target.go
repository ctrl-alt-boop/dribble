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
	TypeTable
)

var ErrNoRequests = errors.New("no requests provided")

type Target struct {
	Name       string
	Type       Type
	DBType     database.Type
	Properties database.ConnectionProperties

	executor database.Database
	// mu               sync.Mutex
	nextRequestID atomic.Int64
}

func New(name string, targetType Type, dbType database.Type, options ...Option) (*Target, error) {
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
	}

	for _, option := range options {
		option(target)
	}

	executor, err := internal.CreateClientForType(target.DBType)
	if err != nil {
		return nil, err
	}
	executor.SetConnectionProperties(target.Properties)
	target.executor = executor

	return target, nil
}

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
// RequestWithHandler sends requests and processes responses synchronously using a handler.
// It blocks until all responses are received and handled or the context is cancelled.
func (t *Target) RequestWithHandler(ctx context.Context, handler database.ResponseHandler, requests ...database.Request) error {
	resultChan, err := t.Request(ctx, requests...)
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case result, ok := <-resultChan:
			if !ok {
				// Channel is closed, all responses have been processed.
				return nil
			}
			// Assuming ResponseHandler is of type `func(database.Response)`.
			// The `result` object itself contains any error information.
			handler(result)
		}
	}
}
