package dribble

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/ctrl-alt-boop/dribble/database"
)

const Version = "0.0.1"

type TargetName = string

type (
	Execution struct {
		RequestID int64
	}
	ExecutionResult struct {
		RequestID int64
		Status    Status
		Result    any
		Err       error
	}
)

type Client struct {
	mu        sync.Mutex
	executors map[TargetName]database.Executor

	nextRequestID atomic.Int64
	pending       map[int64]chan ExecutionResult

	onEvent EventHandler
}

func NewClient() *Client {
	return &Client{
		executors: make(map[TargetName]database.Executor),
		pending:   make(map[int64]chan ExecutionResult),
		onEvent:   func(eventType Status, args any, err error) {},
	}
}

func (c *Client) OnEvent(handler EventHandler) {
	c.onEvent = handler
}

func createExecutor(ctx context.Context, target *database.Target) (database.Executor, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	executor, err := CreateExecutorFromTarget(target)
	if err != nil {
		return nil, fmt.Errorf("error creating executor: %w", err)
	}

	err = executor.Open(ctx)
	if err != nil {
		return nil, fmt.Errorf("error opening target: %w", err)
	}

	return executor, nil
}

func (c *Client) OpenTarget(ctx context.Context, target *database.Target) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	executor, err := createExecutor(ctx, target)
	if err != nil {
		c.onEvent(ErrorTargetOpen, target, err)
		return fmt.Errorf("error creating executor: %w", err)
	}
	c.executors[target.Name] = executor

	c.onEvent(SuccessTargetOpen, target, nil)
	return nil
}

func (c *Client) PingTarget(ctx context.Context, targetName string) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	if targetName == "" {
		return ErrNoTarget
	}
	executor, ok := c.executors[targetName]
	if !ok {
		return fmt.Errorf("executor not found: %s", targetName)
	}
	err := executor.Ping(ctx)
	if err != nil {
		return fmt.Errorf("error pinging target: %w", err)
	}
	return nil
}

func (c *Client) UpdateTarget(ctx context.Context, targetName string, opts ...database.TargetOption) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	if targetName == "" {
		return ErrNoTarget
	}
	c.CloseTarget(ctx, targetName)

	target := c.executors[targetName].Target().Copy(opts...)
	c.executors[targetName].SetTarget(target)

	if err := c.executors[targetName].Open(ctx); err != nil {
		c.onEvent(ErrorTargetUpdate, targetName, err)
		return fmt.Errorf("error updating target: %w", err)
	}
	if err := c.executors[targetName].Ping(ctx); err != nil {
		c.onEvent(ErrorTargetUpdate, targetName, err)
		return fmt.Errorf("error updating target: %w", err)
	}

	c.onEvent(SuccessTargetUpdate, target, nil)
	return nil
}

func (c *Client) CloseTarget(ctx context.Context, targetName ...string) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	var err error
	for _, target := range targetName {
		executor, ok := c.executors[target]
		if !ok {
			err = errors.Join(err, fmt.Errorf("no executor found for target: %s", target))
			c.onEvent(ErrorTargetClose, targetName, err)
			continue
		}
		err := executor.Close(ctx)
		if err != nil {
			err = errors.Join(err, fmt.Errorf("error closing executor for target: %s", target))
			c.onEvent(ErrorTargetClose, targetName, err)
			continue
		}
		delete(c.executors, target)
	}
	if err != nil {
		return fmt.Errorf("error deleting executors: %w", err)
	}
	c.onEvent(SuccessTargetClose, targetName, nil)
	return nil
}

func (c *Client) GetExecutor(targetName string) (database.Executor, error) {
	executor, ok := c.executors[targetName]
	if !ok {
		return nil, fmt.Errorf("executor not found: %s", targetName)
	}
	return executor, nil
}

var ErrNoTarget = errors.New("no target provided")

func (c *Client) Execute(ctx context.Context, intent *database.Intent) (int64, error) {
	if ctx.Err() != nil {
		return 0, ctx.Err()
	}
	if intent.Target == nil {
		return 0, ErrNoTarget
	}

	executor, ok := c.executors[intent.Target.Name]
	if !ok {
		return 0, fmt.Errorf("executor not found: %s", intent.Target.Name)
	}

	requestID := c.nextRequestID.Add(1)
	resultChan := make(chan ExecutionResult, 1)

	c.mu.Lock()
	c.pending[requestID] = resultChan
	c.mu.Unlock()

	go func() {
		defer func() {
			c.mu.Lock()
			delete(c.pending, requestID)
			c.mu.Unlock()
			close(resultChan) // Close the channel after deletion
		}()
		err := c.executors[intent.Target.Name].Ping(ctx)
		if err != nil {
			resultChan <- ExecutionResult{
				RequestID: requestID,
				Status:    ErrorExecute,
				Err:       fmt.Errorf("executor connection error: %w", err),
			}
			return
		}

		result, err := executor.Execute(ctx, intent)
		resultChan <- ExecutionResult{
			RequestID: requestID,
			Status:    SuccessExecute,
			Result:    result,
			Err:       err,
		}
	}()

	return requestID, nil
}

func (c *Client) Request(ctx context.Context, request request) (int64, error) {
	panic("not implemented")
}

// Blocks until result is ready
func (c *Client) GetResult(ctx context.Context, requestID int64) (any, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	resultChan, ok := c.pending[requestID]
	if !ok {
		return nil, fmt.Errorf("no result found for request id: %d", requestID)
	}
	select {
	case result := <-resultChan:
		return result.Result, result.Err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// func (c *Client) FetchDatabaseNames(ctx context.Context, targetName string) (int64, error) {
// 	if ctx.Err() != nil {
// 		return 0, ctx.Err()
// 	}
// 	if targetName == "" {
// 		return 0, ErrNoTarget
// 	}

// 	executor, ok := c.executors[targetName]
// 	if !ok {
// 		return 0, fmt.Errorf("executor not found: %s", targetName)
// 	}

// 	requestID := c.nextRequestID.Add(1)
// 	resultChan := make(chan ExecutionResult, 1)

// 	c.mu.Lock()
// 	c.pending[requestID] = resultChan
// 	c.mu.Unlock()

// 	go func() {
// 		defer func() {
// 			c.mu.Lock()
// 			delete(c.pending, requestID)
// 			c.mu.Unlock()
// 			close(resultChan) // Close the channel after deletion
// 		}()
// 		err := c.executors[targetName].Ping(ctx)
// 		if err != nil {
// 			resultChan <- ExecutionResult{
// 				Err: fmt.Errorf("executor connection error: %w", err),
// 			}
// 			return
// 		}
// 		result, err := executor.ExecutePrefab(ctx, database.PrefabDatabases)
// 		resultChan <- ExecutionResult{
// 			Result: result,
// 			Err:    err,
// 		}
// 	}()

// 	return requestID, nil
// }

// func (c *Client) FetchTableNames(ctx context.Context, targetName string) (int64, error) { // FIXME: redo this
// 	if ctx.Err() != nil {
// 		return 0, ctx.Err()
// 	}
// 	if targetName == "" {
// 		return 0, ErrNoTarget
// 	}

// 	executor, ok := c.executors[targetName]
// 	if !ok {
// 		return 0, fmt.Errorf("executor not found: %s", targetName)
// 	}

// 	requestID := c.nextRequestID.Add(1)
// 	resultChan := make(chan ExecutionResult, 1)

// 	c.mu.Lock()
// 	c.pending[requestID] = resultChan
// 	c.mu.Unlock()

// 	go func() {
// 		defer func() {
// 			c.mu.Lock()
// 			delete(c.pending, requestID)
// 			c.mu.Unlock()
// 			close(resultChan) // Close the channel after deletion
// 		}()
// 		err := c.executors[targetName].Ping(ctx)
// 		if err != nil {
// 			resultChan <- ExecutionResult{
// 				Err: fmt.Errorf("executor connection error: %w", err),
// 			}
// 			return
// 		}
// 		result, err := executor.ExecutePrefab(ctx, database.PrefabTables)
// 		resultChan <- ExecutionResult{
// 			Result: result,
// 			Err:    err,
// 		}
// 	}()

// 	return requestID, nil
// }

// func (c *Client) FetchColumnNames(ctx context.Context, targetName string) (int64, error) { // FIXME: redo this
// 	if ctx.Err() != nil {
// 		return 0, ctx.Err()
// 	}
// 	if targetName == "" {
// 		return 0, ErrNoTarget
// 	}

// 	executor, ok := c.executors[targetName]
// 	if !ok {
// 		return 0, fmt.Errorf("executor not found: %s", targetName)
// 	}

// 	requestID := c.nextRequestID.Add(1)
// 	resultChan := make(chan ExecutionResult, 1)

// 	c.mu.Lock()
// 	c.pending[requestID] = resultChan
// 	c.mu.Unlock()

// 	go func() {
// 		defer func() {
// 			c.mu.Lock()
// 			delete(c.pending, requestID)
// 			c.mu.Unlock()
// 			close(resultChan) // Close the channel after deletion
// 		}()
// 		err := c.executors[targetName].Ping(ctx)
// 		if err != nil {
// 			resultChan <- ExecutionResult{
// 				Err: fmt.Errorf("executor connection error: %w", err),
// 			}
// 			return
// 		}
// 		result, err := executor.ExecutePrefab(ctx, database.PrefabColumns)
// 		resultChan <- ExecutionResult{
// 			Result: result,
// 			Err:    err,
// 		}
// 	}()

// 	return requestID, nil
// }
