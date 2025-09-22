package mongodb

import (
	"context"
	"fmt"

	"strconv"

	"github.com/ctrl-alt-boop/dribble/database"
	"github.com/ctrl-alt-boop/dribble/nosql"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func init() {
	database.DBTypes.Register("NoSQL", "mongo")
}

var _ database.NoSQLClient = &MongoDB{}

type MongoDB struct {
	client *mongo.Client

	clientProperties *nosql.MongoDBClientProperties
}

func NewMongoDBClient() (*MongoDB, error) {
	driver := &MongoDB{
		// clientProperties: clientProperties,
	}
	return driver, nil
}

// Capabilities implements database.Dialect.
func (m *MongoDB) Capabilities() []database.Capabilities {
	return []database.Capabilities{
		database.SupportsBSON,
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
func (m *MongoDB) Client() database.NoSQLClient {
	return m
}

// Request implements database.NoSQL.
func (m *MongoDB) Request(ctx context.Context, requests ...database.Request) (any, error) {
	panic("unimplemented")
}

// RequestWithHandler implements database.NoSQL.
func (m *MongoDB) RequestWithHandler(ctx context.Context, handler func(response database.Response, err error), requests ...database.Request) error {
	panic("unimplemented")
}

// Type implements database.NoSQL.
func (m *MongoDB) Type() database.Type {
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
