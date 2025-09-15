package dribble

import (
	"context"
	"errors"
	"fmt"

	"github.com/ctrl-alt-boop/dribble/database"
)

const Version = "0.0.1"

type TargetName = string

type Client struct {
	executors map[TargetName]database.Executor

	onEvent EventHandler
}

func NewClient() *Client {
	return &Client{
		onEvent:   func(eventType EventType, args any, err error) {},
		executors: make(map[TargetName]database.Executor),
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
		c.onEvent(TargetOpenError, target, err)
		return fmt.Errorf("error creating executor: %w", err)
	}
	executor.OnBefore(c.onExecuteEvent)
	executor.OnAfter(c.onExecuteEvent)
	executor.OnResult(c.onExecuteResult)
	c.executors[target.Name] = executor

	c.onEvent(TargetOpened, target, nil)
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
		c.onEvent(TargetUpdateError, targetName, err)
		return fmt.Errorf("error updating target: %w", err)
	}
	if err := c.executors[targetName].Ping(ctx); err != nil {
		c.onEvent(TargetUpdateError, targetName, err)
		return fmt.Errorf("error updating target: %w", err)
	}

	c.onEvent(TargetUpdated, target, nil)
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
			c.onEvent(TargetCloseError, targetName, err)
			continue
		}
		err := executor.Close(ctx)
		if err != nil {
			err = errors.Join(err, fmt.Errorf("error closing executor for target: %s", target))
			c.onEvent(TargetCloseError, targetName, err)
			continue
		}
		delete(c.executors, target)
	}
	if err != nil {
		return fmt.Errorf("error deleting executors: %w", err)
	}
	c.onEvent(TargetClosed, targetName, nil)
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

func (c *Client) Execute(ctx context.Context, intent *database.Intent) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	if intent.Target == nil {
		return ErrNoTarget
	}
	executor, ok := c.executors[intent.Target.Name]
	if !ok {
		return fmt.Errorf("executor not found: %s", intent.Target.Name)
	}
	err := c.executors[intent.Target.Name].Ping(ctx)
	if err != nil {
		return fmt.Errorf("executor connection error: %w", err)
	}

	go func() {
		err := executor.Execute(ctx, intent)
		if err != nil {
			c.onEvent(QueryExecuteError, intent, err)
			return
		}

		c.onEvent(QueryExecuted, intent, nil)
	}()

	return nil
}

func (c *Client) FetchDatabaseNames(ctx context.Context, targetName string) error {
	if targetName == "" {
		return ErrNoTarget
	}
	executor, ok := c.executors[targetName]
	if !ok {
		return fmt.Errorf("target not open: %s", targetName)
	}

	go func() {
		err := executor.ExecutePrefab(ctx, database.PrefabDatabases)
		if err != nil {
			c.onEvent(DatabaseListFetchError, targetName, err)
			return
		}

		// c.onEvent(DatabaseListFetched, targetName, nil)
	}()

	return nil
}

func (c *Client) FetchTableNames(ctx context.Context, targetName string) error { // FIXME: redo this
	if targetName == "" {
		return ErrNoTarget
	}
	executor, ok := c.executors[targetName]
	if !ok {
		return fmt.Errorf("target not open: %s", targetName)
	}

	go func() {
		err := executor.ExecutePrefab(ctx, database.PrefabTables)
		if err != nil {
			c.onEvent(DBTableListFetchError, targetName, err)
			return
		}
		// c.onEvent(DBTableListFetched, targetName, nil)
	}()

	return nil
}

func (c *Client) FetchColumnNames(ctx context.Context, targetName string) error { // FIXME: redo this
	if targetName == "" {
		return ErrNoTarget
	}
	executor, ok := c.executors[targetName]
	if !ok {
		return fmt.Errorf("target not open: %s", targetName)
	}

	go func() {
		err := executor.ExecutePrefab(ctx, database.PrefabColumns)
		if err != nil {
			c.onEvent(DBTableListFetchError, targetName, err)
			return
		}
		// c.onEvent(DBTableListFetched, targetName, nil)
	}()

	return nil
}

func (c *Client) onExecuteEvent(intent *database.Intent, err error) {
	if err != nil {
		c.onEvent(QueryExecuteError, nil, err)
		return
	}
	c.onEvent(QueryExecuted, nil, nil)
}

func (c *Client) onExecuteResult(result any, err error) {
	if err != nil {
		c.onEvent(QueryExecuteError, nil, err)
		return
	}
	c.onEvent(QueryExecuted, result, nil)
}
