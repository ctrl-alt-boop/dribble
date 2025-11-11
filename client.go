package dribble

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"strings"

	"github.com/ctrl-alt-boop/dribble/datasource"
	"github.com/ctrl-alt-boop/dribble/request"
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

// func (c *Client) Target(targetName string) *target.Target {
// 	return c.targets[targetName]
// }

func (c *Client) Target(targetName string) (*target.Target, bool) {
	target, ok := c.targets[targetName]
	return target, ok
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
	for _, t := range targets {
		err := c.OpenTarget(ctx, t)
		if err != nil {
			errs = errors.Join(errs, err)
		}
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

func (c *Client) UpdateTarget(ctx context.Context, targetName string, opts ...target.TargetOption) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	if targetName == "" {
		return ErrNoTarget
	}
	if _, ok := c.targets[targetName]; !ok {
		return ErrTargetNotFound(targetName)
	}
	c.targets[targetName].Update(ctx, opts...)

	if err := c.targets[targetName].Open(ctx); err != nil {
		return ErrUpdateTarget(targetName, err)
	}
	if err := c.targets[targetName].Ping(ctx); err != nil {
		return ErrUpdateTarget(targetName, err)
	}

	return nil
}

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

func (c *Client) CloseTargets(ctx context.Context, targets ...string) error {
	var errs error
	for _, targetName := range targets {
		err := c.CloseTarget(ctx, targetName)
		if err != nil {
			errs = errors.Join(errs, err)
		}
	}

	return errs
}

type ResponseHandler func(*request.Response)

func (c *Client) Request(ctx context.Context, targetName string, request datasource.Request) (chan *request.Response, error) {
	requestTarget, ok := c.targets[targetName]
	if !ok || requestTarget == nil {
		return nil, ErrTargetNotFound(targetName)
	}
	return requestTarget.Request(ctx, request)
}

// Non-Blocking
func (c *Client) RequestWithHandler(ctx context.Context, handler ResponseHandler, targetName string, req datasource.Request) error {
	requestTarget, ok := c.targets[targetName]
	if !ok {
		return ErrTargetNotFound(targetName)
	}

	return requestTarget.RequestWithHandler(ctx, handler, req)
}

// no targets = NoOp
func (c *Client) RequestForAll(ctx context.Context, req datasource.Request) (chan *request.Response, error) {
	if len(c.targets) == 0 {
		return nil, nil
	}
	var errs error
	responseChan := make(chan *request.Response, len(c.targets))
	defer close(responseChan)

	for _, target := range c.targets {
		respChan, err := target.Request(ctx, req)
		if err != nil {
			errs = errors.Join(errs, err)
			continue
		}
		responseChan <- (<-respChan)
	}
	return responseChan, errs
}

// Blocking
func (c *Client) PerformWithHandler(ctx context.Context, handler ResponseHandler, targetName string, req datasource.Request) error {
	requestTarget, ok := c.targets[targetName]
	if !ok {
		return ErrTargetNotFound(targetName)
	}

	return requestTarget.PerformWithHandler(ctx, handler, req)
}

// no targets = NoOp
func (c *Client) PerformForAll(ctx context.Context, handler ResponseHandler, req datasource.Request) error {
	if len(c.targets) == 0 {
		return nil
	}
	var errs error
	responseChan := make(chan *request.Response, len(c.targets))
	defer close(responseChan)

	for _, target := range c.targets {
		respChan, err := target.Request(ctx, req)
		if err != nil {
			errs = errors.Join(errs, err)
			continue
		}
		responseChan <- (<-respChan)
	}
	for response := range responseChan { // Maybe...
		handler(response)
	}
	return errs
}
