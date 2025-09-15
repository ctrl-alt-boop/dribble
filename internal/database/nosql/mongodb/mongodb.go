package mongodb

import (
	"context"
	"fmt"

	"github.com/ctrl-alt-boop/dribble/database"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var _ database.Driver = &MongoDB{}
var _ database.NoSQLClient = &MongoDB{}

type MongoDB struct {
	client *mongo.Client
}

func NewMongoDBDriver(target *database.Target) (*MongoDB, error) {
	driver := &MongoDB{}
	return driver, nil
}

// Capabilities implements database.Dialect.
func (m *MongoDB) Capabilities() []database.Capabilities {
	return []database.Capabilities{}
}

// Open implements database.NoSQLClient.
func (m *MongoDB) Open(ctx context.Context, target *database.Target) error {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(m.ConnectionString(target)))
	if err != nil {
		return err
	}
	m.client = client
	return nil
}

// Close implements database.NoSQLClient.
func (m *MongoDB) Close(ctx context.Context) error {
	err := m.client.Disconnect(ctx)
	if err != nil {
		return err
	}
	return nil
}

// ConnectionString implements database.Driver.
func (m *MongoDB) ConnectionString(target *database.Target) string {
	return fmt.Sprintf("mongodb://%s:%d", target.Ip, target.Port)
}

// Read implements database.NoSQLClient.
func (m *MongoDB) Read(any) {
	panic("unimplemented")
}

// ReadMany implements database.NoSQLClient.
func (m *MongoDB) ReadMany(any) {
	panic("unimplemented")
}

// Create implements database.NoSQLClient.
func (m *MongoDB) Create(any) {
	panic("unimplemented")
}

// Update implements database.NoSQLClient.
func (m *MongoDB) Update(any) {
	panic("unimplemented")
}

// Delete implements database.NoSQLClient.
func (m *MongoDB) Delete(any) {
	panic("unimplemented")
}

func (m *MongoDB) Dialect() database.Dialect {
	return m
}

// RenderIntent implements database.Driver.
func (m *MongoDB) RenderIntent(intent *database.Intent) (string, error) {
	panic("unimplemented")
}

// ResolveType implements database.Dialect.
func (m *MongoDB) ResolveType(dbType string, value []byte) (any, error) {
	panic("unimplemented")
}
