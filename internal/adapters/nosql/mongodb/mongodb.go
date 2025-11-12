package mongodb

import (
	"context"
	"fmt"
	"strconv"

	"github.com/ctrl-alt-boop/dribble/datasource"
	"github.com/ctrl-alt-boop/dribble/nosql"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var _ datasource.NoSQLAdapter = (*MongoDB)(nil)

type MongoDB struct {
	client *mongo.Client

	clientProperties *nosql.MongoDBClientProperties
}

func NewMongoDBClient() *MongoDB {
	return &MongoDB{}
}

// Capabilities implements database.Dialect.
func (m *MongoDB) Capabilities() []datasource.Capability {
	return []datasource.Capability{
		datasource.SupportsBSON,
	}
}

// SetConnectionProperties implements database.NoSQLClient.
func (m *MongoDB) SetConnectionProperties(props map[string]string) {
	port, err := strconv.Atoi(props["port"])
	if err != nil {
		port = 27017 // Default MongoDB port
	}
	m.clientProperties = &nosql.MongoDBClientProperties{
		Ip:   props["ip"],
		Port: port,
	}
}

// Open implements database.NoSQLClient.
func (m *MongoDB) Open(ctx context.Context) error {
	connString := fmt.Sprintf("mongodb://%s:%d", m.clientProperties.Ip, m.clientProperties.Port)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connString))
	if err != nil {
		return err
	}
	m.client = client
	return nil
}

// Ping implements database.NoSQLClient.
func (m *MongoDB) Ping(ctx context.Context) error {
	err := m.client.Ping(ctx, nil)
	if err != nil {
		return err
	}
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

// Client implements database.NoSQL.
func (m *MongoDB) Client() datasource.NoSQLAdapter {
	return m
}

// Request implements database.NoSQL.
func (m *MongoDB) Request(ctx context.Context, requests ...datasource.Request) (any, error) {
	panic("unimplemented")
}

// RequestWithHandler implements database.NoSQL.
func (m *MongoDB) RequestWithHandler(ctx context.Context, handler func(response datasource.Response, err error), requests ...datasource.Request) error {
	panic("unimplemented")
}

// Type implements database.NoSQL.
func (m *MongoDB) Type() datasource.Type {
	panic("unimplemented")
}

// Create implements database.NoSQLClient.
func (m *MongoDB) Create(any) {
	panic("unimplemented")
}

// Delete implements database.NoSQLClient.
func (m *MongoDB) Delete(any) {
	panic("unimplemented")
}

// Read implements database.NoSQLClient.
func (m *MongoDB) Read(any) {
	panic("unimplemented")
}

// ReadMany implements database.NoSQLClient.
func (m *MongoDB) ReadMany(any) {
	panic("unimplemented")
}

// Update implements database.NoSQLClient.
func (m *MongoDB) Update(any) {
	panic("unimplemented")
}
