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
	executors     map[TargetName]*database.Executor
	loadedDrivers []DriverName

	onEvent func(eventType EventType, args any, err error)
}

func NewClient() *Client {
	return &Client{
		onEvent: func(eventType EventType, args any, err error) {

		},
		executors:     make(map[TargetName]*database.Executor),
		loadedDrivers: make([]DriverName, 0),
	}
}

func (c *Client) OnEvent(handler func(eventType EventType, args any, err error)) {
	c.onEvent = handler
}

func (c *Client) OpenTarget(ctx context.Context, target *database.Target) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	executor, err := createQueryExecutor(ctx, target)
	if err != nil {
		c.onEvent(TargetOpenError, target, err)
		return fmt.Errorf("error creating executor: %w", err)
	}
	executor.OnQueryExecuted(c.onExecutorEvent)
	c.executors[target.Name] = executor

	c.onEvent(TargetOpened, target, nil)
	return nil
}

func (c *Client) UpdateTarget(ctx context.Context, targetName string, opts ...database.TargetOption) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	if targetName == "" {
		return ErrNoTargetName
	}
	c.CloseTarget(ctx, targetName)

	target := c.executors[targetName].Target().Copy(opts...)
	c.executors[targetName].SetTarget(target)

	err := c.executors[targetName].VerifyConnection(ctx)
	if err != nil {
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

var ErrNoTargetName = errors.New("no target name provided")

func (c *Client) Query(ctx context.Context, query *database.Intent) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	if query.TargetName == "" {
		return ErrNoTargetName
	}
	executor, ok := c.executors[query.TargetName]
	if !ok {
		return fmt.Errorf("target not open: %s", query.TargetName)
	}

	go func() {
		data, err := executor.Query(ctx, query)
		if err != nil {
			c.onEvent(QueryExecuteError, query, err)
			return
		}
		// TODO: Maybe something like a data verify
		c.onEvent(QueryExecuted, data, nil)
	}()

	return nil
}

func (c *Client) FetchDatabaseNames(ctx context.Context, targetName string) error {
	if targetName == "" {
		return ErrNoTargetName
	}
	executor, ok := c.executors[targetName]
	if !ok {
		return fmt.Errorf("target not open: %s", targetName)
	}

	go func() {
		data, err := executor.ExecutePrefab(ctx, database.PrefabTypeDatabases)
		if err != nil {
			c.onEvent(DatabaseListFetchError, targetName, err)
			return
		}
		// TODO: Maybe something like a data verify
		c.onEvent(DatabaseListFetched, data, nil)
	}()

	return nil
}

func (c *Client) FetchTableNames(ctx context.Context, targetName string) error {
	if targetName == "" {
		return ErrNoTargetName
	}
	executor, ok := c.executors[targetName]
	if !ok {
		return fmt.Errorf("target not open: %s", targetName)
	}

	go func() {
		data, err := executor.ExecutePrefab(ctx, database.PrefabTypeTables)
		if err != nil {
			c.onEvent(DBTableListFetchError, targetName, err)
			return
		}
		// TODO: Maybe something like a data verify
		c.onEvent(DBTableListFetched, data, nil)
	}()

	return nil
}

func (c *Client) FetchColumnNames(ctx context.Context, targetName string) error {
	if targetName == "" {
		return ErrNoTargetName
	}
	executor, ok := c.executors[targetName]
	if !ok {
		return fmt.Errorf("target not open: %s", targetName)
	}

	go func() {
		data, err := executor.ExecutePrefab(ctx, database.PrefabTypeColumns)
		if err != nil {
			c.onEvent(DBTableListFetchError, targetName, err)
			return
		}
		// TODO: Maybe something like a data verify
		c.onEvent(DBTableListFetched, data, nil)
	}()

	return nil
}

func (c *Client) onExecutorEvent(query string, err error) {
	if err != nil {
		c.onEvent(QueryExecuteError, nil, err)
		return
	}
	c.onEvent(QueryExecuted, nil, nil)
}
