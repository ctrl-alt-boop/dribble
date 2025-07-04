package nosql

import (
	"context"
	"fmt"

	"github.com/ctrl-alt-boop/dribble/database"
	"go.mongodb.org/mongo-driver/mongo"
)

var _ database.Driver = &MongoDB{}

type MongoDB struct {
	Client *mongo.Client
	target *database.Target

	FetchLimit       int
	FetchLimitOffset int
}

func NewMongoDBDriver(target *database.Target) (*MongoDB, error) {
	if target.DriverName != "mongodb" {
		return nil, fmt.Errorf("invalid driver name: %s", target.DriverName)
	}
	driver := &MongoDB{
		target: target,
	}
	return driver, nil

}

func (m *MongoDB) Close(ctx context.Context) error {
	return m.Client.Disconnect(ctx)
}

func (m *MongoDB) Open(ctx context.Context) error {
	client, err := mongo.Connect(ctx, nil)
	if err != nil {
		return err
	}
	m.Client = client
	return nil
}

func (m *MongoDB) Ping(ctx context.Context) error {
	return m.Client.Ping(ctx, nil)
}

func (m *MongoDB) Dialect() database.Dialect {
	return m
}

func (m *MongoDB) Query(query *database.QueryIntent) (any, error) {
	return m.QueryContext(context.Background(), query)
}

func (m *MongoDB) QueryContext(ctx context.Context, query *database.QueryIntent) (any, error) {
	panic("unimplemented")
}

func (m *MongoDB) SetTarget(target *database.Target) {
	m.target = target
}

func (m *MongoDB) Target() *database.Target {
	return m.target
}

// Capabilities implements database.Dialect.
func (m *MongoDB) Capabilities() []database.DialectProperties {
	return nil
}

// GetTemplate implements database.Dialect.
func (m *MongoDB) GetTemplate(queryType database.QueryType) string {
	switch queryType {
	case database.ReadQuery:
		return "" // MongoDBSelectTemplate
	case database.CreateQuery:
		return ""
	case database.UpdateQuery:
		return ""
	case database.DeleteQuery:
		return ""
	default:
		return ""
	}
}

// Quote implements database.Dialect.
func (m *MongoDB) Quote(value string) string {
	panic("unimplemented")
}

// QuoteRune implements database.Dialect.
func (m *MongoDB) QuoteRune() rune {
	panic("unimplemented")
}

// RenderCurrentTimestamp implements database.Dialect.
func (m *MongoDB) RenderCurrentTimestamp() string {
	panic("unimplemented")
}

// RenderPlaceholder implements database.Dialect.
func (m *MongoDB) RenderPlaceholder(index int) string {
	panic("unimplemented")
}

// RenderTypeCast implements database.Dialect.
func (m *MongoDB) RenderTypeCast() string {
	panic("unimplemented")
}

// RenderValue implements database.Dialect.
func (m *MongoDB) RenderValue(value any) string {
	panic("unimplemented")
}

// ResolveType implements database.Dialect.
func (m *MongoDB) ResolveType(dbType string, value []byte) (any, error) {
	panic("unimplemented")
}
