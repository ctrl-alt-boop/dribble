package target

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/ctrl-alt-boop/dribble/database"
	"github.com/ctrl-alt-boop/dribble/internal/datasource"
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
	nextRequestID atomic.Int64
	Name          string
	Type          Type
	DSN           database.DataSourceNamer
	DBType        database.Type

	dataSource database.Database
}

func New(name string, dsn database.DataSourceNamer) (*Target, error) {
	target := &Target{
		Name:   name,
		Type:   TypeDriver,
		DSN:    dsn,
		DBType: dsn.Type(),
	}

	executor, err := datasource.Create(dsn)
	if err != nil {
		return nil, err
	}
	target.dataSource = executor

	return target, nil
}

func (t *Target) String() string {
	return fmt.Sprintf("%s: %s", t.Name, t.DSN.Info())
}

func (t *Target) Open(ctx context.Context) error {
	return t.dataSource.Open(ctx)
}

func (t *Target) Ping(ctx context.Context) error {
	return t.dataSource.Ping(ctx)
}

func (t *Target) Close(ctx context.Context) error {
	return t.dataSource.Close(ctx)
}

func (t *Target) Request(ctx context.Context, req database.Request) (chan *request.Response, error) {
	switch r := req.(type) {
	case request.ChainRequest:
		return t.chainedRequest(ctx, r)
	case request.BatchRequest:
		return t.batchRequest(ctx, r)
	default:
		return t.simpleRequest(ctx, req)
	}
}

func (t *Target) simpleRequest(ctx context.Context, req database.Request) (chan *request.Response, error) {
	requestID := t.nextRequestID.Add(1)
	resultChan := make(chan *request.Response, 1)

	go func() {
		defer close(resultChan)

		requestResult, err := t.dataSource.Request(ctx, req)
		var resp database.Response
		if err != nil {
			resp = req.ResponseOnError()
		} else {
			resp = req.ResponseOnSuccess()
		}
		resultChan <- &request.Response{
			Status:        request.Status(resp.Code()),
			RequestID:     requestID,
			RequestTarget: t.Name,
			Body:          requestResult,
			Error:         err,
		}

	}()

	return resultChan, nil
}

func (t *Target) chainedRequest(ctx context.Context, requestChain request.ChainRequest) (chan *request.Response, error) {
	if len(requestChain) == 0 {
		return nil, ErrNoRequests
	}
	requestID := t.nextRequestID.Add(1)
	resultChan := make(chan *request.Response, 1)
	responses := make([]*request.Response, len(requestChain))

	go func() {
		defer func() {
			resultChan <- &request.Response{
				Status:        request.Status(requestChain.ResponseOnError().Code()),
				RequestID:     requestID,
				RequestTarget: t.Name,
				Body:          responses,
				Error:         nil,
			}
			close(resultChan)
		}()
		for i, req := range requestChain {
			requestResult, err := t.dataSource.Request(ctx, req)
			var resp database.Response
			if err != nil {
				resp = req.ResponseOnError()
			} else {
				resp = req.ResponseOnSuccess()
			}
			responses[i] = &request.Response{
				Status:        request.Status(resp.Code()),
				RequestID:     int64(i),
				RequestTarget: t.Name,
				Body:          requestResult,
				Error:         err,
			}

			if err != nil {
				// If an error occurs in a chained request, stop the chain.
				return
			}
		}
	}()

	// requestID := t.nextRequestID.Add(1)
	return resultChan, nil
}

func (t *Target) batchRequest(ctx context.Context, requestBatch request.BatchRequest) (chan *request.Response, error) {
	numRequests := len(requestBatch)
	if numRequests == 0 {
		return nil, ErrNoRequests
	}

	resultChan := make(chan *request.Response, len(requestBatch))

	go func() {
		defer close(resultChan)

		var wg sync.WaitGroup
		wg.Add(numRequests)

		for _, req := range requestBatch {
			go func(r database.Request) {
				defer wg.Done()
				requestID := t.nextRequestID.Add(1)
				requestResult, err := t.dataSource.Request(ctx, r)
				var resp database.Response
				if err != nil {
					resp = r.ResponseOnError()
				} else {
					resp = r.ResponseOnSuccess()
				}
				resultChan <- &request.Response{
					Status:        request.Status(resp.Code()),
					RequestID:     requestID,
					RequestTarget: t.Name,
					Body:          requestResult,
					Error:         err,
				}
			}(req)
		}
		wg.Wait()
	}()
	return resultChan, nil
}

// Blocking
// PerformWithHandler sends requests and processes responses synchronously using a handler.
// It blocks until all responses are received and handled or the context is cancelled.
func (t *Target) PerformWithHandler(ctx context.Context, handler func(*request.Response), req database.Request) error {
	resultChan, err := t.Request(ctx, req)
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
			handler(result)
			if result.Error != nil {
				// The handler has already processed this error response.
				// We return the error to signal that the operation was not fully successful.
				return result.Error
			}
		}
	}
}

// Non-Blocking
// RequestWithHandler sends requests and processes responses asynchronously using a handler.
func (t *Target) RequestWithHandler(ctx context.Context, handler func(*request.Response), req database.Request) error {
	resultChan, err := t.Request(ctx, req)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case result, ok := <-resultChan:
				if !ok {
					// Channel is closed, all responses have been processed.
					return
				}
				// The `result` object itself contains any error information.
				handler(result)
			}
		}
	}()

	return nil
}

func (t *Target) Update(ctx context.Context, opts ...TargetOption) error {
	for _, opt := range opts {
		if err := opt(ctx, t); err != nil {
			return err
		}
	}
	return nil
}

type TargetOption func(context.Context, *Target) error

func WithDataSource(dsn database.DataSourceNamer) TargetOption {
	return func(ctx context.Context, t *Target) error {
		if dsn == nil {
			return errors.New("dsn cannot be nil")
		}
		if !t.dataSource.IsClosed() {
			if err := t.dataSource.Close(ctx); err != nil {
				return err
			}
		}

		t.DSN = dsn
		t.DBType = dsn.Type()
		executor, err := datasource.Create(dsn)
		if err != nil {
			return err
		}
		t.dataSource = executor
		return nil
	}
}

func WithName(newName string) TargetOption {
	return func(ctx context.Context, t *Target) error {
		t.Name = newName
		return nil
	}
}
