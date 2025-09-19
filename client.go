package dribble

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/ctrl-alt-boop/dribble/database"
	"github.com/ctrl-alt-boop/dribble/target"
)

const Version = "0.0.1"

func ErrUpdateTarget(targetName string, err error) error {
	return fmt.Errorf("error updating target %s: %w", targetName, err)
}

func ErrTargetNotFound(targetName string) error {
	return fmt.Errorf("target not found: %s", targetName)
}

var (
	ErrNoTarget   = errors.New("no target provided")
	ErrNoRequests = errors.New("no requests provided")
)

type Client struct {
	mu        sync.Mutex
	executors map[string]database.Executor
	targets   map[string]*target.Target

	nextRequestID atomic.Int64
	// requestChannels map[string]chan ...

}

func NewClient() *Client {
	return &Client{
		executors: make(map[string]database.Executor),
		// requestChannels: make(map[string]chan ...),

	}
}

func (c *Client) OnEvent(handler EventHandler) {

}

func createExecutor(ctx context.Context, target *target.Target) (database.Executor, error) {
	executor, err := CreateExecutorFromTarget(target)
	if err != nil {
		return nil, fmt.Errorf("error creating executor: %w", err)
	}

	return executor, nil
}

func (c *Client) OpenTarget(ctx context.Context, target *target.Target) error {
	executor, err := createExecutor(ctx, target)
	if err != nil {

		return err
	}

	if ctx.Err() != nil {
		return ctx.Err()
	}
	err = executor.Open(ctx)
	if err != nil {
		return fmt.Errorf("error opening target: %w", err)
	}
	c.executors[target.Name] = executor

	return nil
}

func (c *Client) PingTarget(ctx context.Context, targetName string) error {

	if targetName == "" {
		return ErrNoTarget
	}
	executor, ok := c.executors[targetName]
	if !ok {
		return fmt.Errorf("executor not found: %s", targetName)
	}
	if ctx.Err() != nil {
		return ctx.Err()
	}
	err := executor.Ping(ctx)
	if err != nil {
		return fmt.Errorf("error pinging target: %w", err)
	}
	return nil
}

// func (c *Client) UpdateTarget(ctx context.Context, targetName string, opts ...target.Option) error {
// 	if ctx.Err() != nil {
// 		return ctx.Err()
// 	}
// 	if targetName == "" {
// 		return ErrNoTarget
// 	}
// 	c.CloseTarget(ctx, targetName)

// 	target := c.executors[targetName].Target().Copy(opts...)
// 	c.executors[targetName].SetTarget(target)

// 	if err := c.executors[targetName].Open(ctx); err != nil {

// 		return ErrUpdateTarget(targetName, err)
// 	}
// 	if err := c.executors[targetName].Ping(ctx); err != nil {

// 		return ErrUpdateTarget(targetName, err)
// 	}

// 	return nil
// }

func (c *Client) CloseTarget(ctx context.Context, targetName ...string) error {
	var err error
	for _, target := range targetName {
		executor, ok := c.executors[target]
		if !ok {
			err = errors.Join(err, fmt.Errorf("no executor found for target: %s", target))

			continue
		}
		err := executor.Close(ctx)
		if err != nil {
			err = errors.Join(err, fmt.Errorf("error closing executor for target: %s", target))

			continue
		}
		delete(c.executors, target)
	}
	if err != nil {
		return fmt.Errorf("error deleting executors: %w", err)
	}

	return nil
}

func (c *Client) Request(ctx context.Context, targetName string, requests ...database.Request) (chan database.Response, error) {
	requestTarget, ok := c.targets[targetName]
	if !ok {
		return nil, ErrTargetNotFound(targetName)
	}
	numRequests := len(requests)
	if numRequests == 0 {
		return nil, ErrNoRequests
	}

	return requestTarget.Request(ctx, requests...)
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
