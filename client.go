package dribble

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"strings"

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
	targets map[string]*target.Target
}

func NewClient() *Client {
	return &Client{
		targets: make(map[string]*target.Target),
	}
}

func (c *Client) String() string {
	targets := maps.Values(c.targets)
	targetStrings := []string{}
	for t := range targets {
		targetStrings = append(targetStrings, fmt.Sprint(t))
	}
	return fmt.Sprintf("dribble version: %s \ntargets:\n%s\n", Version, strings.Join(targetStrings, "\n"))
}

func (c *Client) OnEvent(handler EventHandler) {

}

func (c *Client) Target(targetName string) *target.Target {
	return c.targets[targetName]
}

func (c *Client) OpenTarget(ctx context.Context, t *target.Target) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	err := t.Open(ctx)
	if err != nil {
		return fmt.Errorf("error opening target: %w", err)
	}
	c.targets[t.Name] = t

	return nil
}

func (c *Client) OpenTargets(ctx context.Context, targets ...*target.Target) error {
	var errs error
	if ctx.Err() != nil {
		return ctx.Err()
	}
	for _, t := range targets {
		err := t.Open(ctx)
		if err != nil {
			errs = errors.Join(errs, err)
			continue
		}
		c.targets[t.Name] = t
	}
	return errs
}

func (c *Client) PingTarget(ctx context.Context, targetName string) error {
	if targetName == "" {
		return ErrNoTarget
	}
	t, ok := c.targets[targetName]
	if !ok {
		return ErrTargetNotFound(targetName)
	}
	return t.Ping(ctx)
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

func (c *Client) CloseTarget(ctx context.Context, targetName string) error {
	if targetName == "" {
		return ErrNoTarget
	}
	t, ok := c.targets[targetName]
	if !ok {
		return ErrTargetNotFound(targetName)
	}
	err := t.Close(ctx)
	if err != nil {
		return fmt.Errorf("error closing executor for target: %s", targetName)
	}
	delete(c.targets, targetName)
	return nil
}

func (c *Client) CloseTargets(ctx context.Context, targetName ...string) error {
	var errs error
	for _, target := range targetName {
		t, ok := c.targets[target]
		if !ok {
			errs = errors.Join(errs, ErrTargetNotFound(target))
			continue
		}
		err := t.Close(ctx)
		if err != nil {
			errs = errors.Join(errs, fmt.Errorf("error closing executor for target: %s", target))
		}
		delete(c.targets, target)
		if err != nil {
			errs = errors.Join(errs, err)
		}
	}

	return errs
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

// // Blocks until result is ready
// func (c *Client) GetResult(ctx context.Context, requestID int64) (any, error) {
// 	if ctx.Err() != nil {
// 		return nil, ctx.Err()
// 	}
// 	resultChan, ok := c.pending[requestID]
// 	if !ok {
// 		return nil, fmt.Errorf("no result found for request id: %d", requestID)
// 	}
// 	select {
// 	case result := <-resultChan:
// 		return result.Result, result.Err
// 	case <-ctx.Done():
// 		return nil, ctx.Err()
// 	}
// }
