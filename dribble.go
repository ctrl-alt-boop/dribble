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
	executers     map[TargetName]*QueryExecuter
	loadedDrivers []DriverName

	onEvent func(eventType EventType, args any, err error)
}

func NewClient() *Client {
	return &Client{
		onEvent: func(eventType EventType, args any, err error) {

		},
		executers:     make(map[TargetName]*QueryExecuter),
		loadedDrivers: make([]DriverName, 0),
	}
}

func (c *Client) OnEvent(handler func(eventType EventType, args any, err error)) {
	c.onEvent = handler
}

func (c *Client) CreateExecuter(ctx context.Context, target *database.Target) error {
	executer, err := createQueryExecuter(target)
	if err != nil {
		return fmt.Errorf("error creating executer: %w", err)
	}
	executer.OnQueryExecuted(c.onQueryExecuted)
	c.executers[target.Name] = executer
	return nil
}

func (c *Client) UpdateTarget(ctx context.Context, targetName string, opts ...database.TargetOption) error {
	c.DeleteExecuter(context.TODO(), targetName)
	target := c.executers[targetName].Target().Copy(opts...)
	c.executers[targetName].SetTarget(target)
	err := c.executers[targetName].VerifyConnection()
	if err != nil {
		return fmt.Errorf("error updating target: %w", err)
	}
	return nil
}

func (c *Client) DeleteExecuter(ctx context.Context, targetName ...string) error {
	var err error
	for _, target := range targetName {
		executer, ok := c.executers[target]
		if !ok {
			err = errors.Join(err, fmt.Errorf("no executer found for target: %s", target))
			continue
		}
		err := executer.Close(context.TODO())
		if err != nil {
			err = errors.Join(err, fmt.Errorf("error closing executer for target: %s", target))
			continue
		}
		delete(c.executers, target)
	}
	if err != nil {
		return fmt.Errorf("error deleting executers: %w", err)
	}
	return nil
}

func (c *Client) Query(query *database.QueryIntent) (any, error) {
	return c.QueryContext(context.Background(), query)
}

var ErrNoTargetName = errors.New("no target name provided")

func (c *Client) QueryContext(ctx context.Context, query *database.QueryIntent) (any, error) {
	if query.TargetName == "" {
		return nil, ErrNoTargetName
	}
	executer, ok := c.executers[query.TargetName]
	if !ok {
		return nil, fmt.Errorf("no executer found for target: %s", query.TargetName)
	}
	queryCtx, queryCancel := context.WithCancel(ctx)
	defer queryCancel()

	return executer.QueryContext(queryCtx, query)
}

func (c *Client) onQueryExecuted(query string, err error) {
	// sqlLogger.Infof("Executed query: %s", query)
	if err != nil {
		// sqlLogger.Error("\t", err)
	}
}
