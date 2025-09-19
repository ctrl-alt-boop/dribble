package target

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"

	"github.com/ctrl-alt-boop/dribble/database"
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
	Properties Properties

	database database.Database
	// mu               sync.Mutex
	nextRequestID    atomic.Int64
	pendingResponses map[int64]chan database.Response // TODO: Check if necessary
}

func New(name string, targetType Type, dialect database.SQLDialect, options ...Option) *Target {
	target := &Target{
		Name: name,
		Type: targetType,
		Properties: Properties{
			Dialect:    dialect,
			Ip:         "localhost",
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

	return target
}

func (t *Target) Initialize() error {
	executor, err := database.CreateClientForDialect(t.Properties.Dialect)
	if err != nil {
		return err
	}
	t.database = executor

	return nil
}

func (t *Target) Open(ctx context.Context) error {
	return t.database.Open(ctx)
}

func (t *Target) Ping(ctx context.Context) error {
	return t.database.Ping(ctx)
}

func (t *Target) Close(ctx context.Context) error {
	return t.database.Close(ctx)
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

				requestResult, err := t.database.Request(ctx, r)
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
