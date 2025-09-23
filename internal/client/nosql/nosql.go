package nosql

import (
	"context"
	"fmt"

	"github.com/ctrl-alt-boop/dribble/database"
	"github.com/ctrl-alt-boop/dribble/internal/client/nosql/firestore"
	"github.com/ctrl-alt-boop/dribble/internal/client/nosql/mongodb"
)

// var Defaults = map[string]*target.Target{
// 	MongoDB: {
// 		Name: "mongo",
// 		Type: target.TypeDriver,
// 		Properties: target.Properties{
// 			Dialect: target.MongoDB,
// 			Addr:    "127.0.0.1",
// 			Port:    27017,
// 		},
// 	},
// }

var _ database.NoSQL = &Executor{}

type Method string

func (s Method) String() string {
	return string(s)
}

const (
	MethodFind Method = "FIND"
)

const DefaultFindLimit int = 10 // Just a safeguard

var NoSQLMethods = []Method{MethodFind}

type Executor struct {
	client database.NoSQLClient
}

func NewExecutor(modelType database.NoSQLModelType) (*Executor, error) {
	var client database.NoSQLClient
	var err error
	switch modelType {
	case database.MongoDB:
		client, err = mongodb.NewMongoDBClient()
	case database.Firestore:
		client, err = firestore.NewFirestoreClient()
	default:
		return nil, fmt.Errorf("unknown or unsupported database model: %s", modelType)
	}
	if err != nil {
		panic(err)
	}
	return &Executor{
		client: client,
	}, nil
}

func New(dsn database.DataSourceNamer) (*Executor, error) {
	var client database.NoSQLClient
	var err error
	switch dsn.Type() {
	case database.MongoDB:
		client, err = mongodb.NewMongoDBClient()
	case database.Firestore:
		client, err = firestore.NewFirestoreClient()
	case database.Redis:
		client, err = nil, fmt.Errorf("redis not implemented")
	default:
		return nil, fmt.Errorf("unknown or unsupported database model: %s", dsn.Type())
	}
	if err != nil {
		panic(err)
	}
	return &Executor{
		client: client,
	}, nil
}

// Client implements database.NoSQL.
func (e *Executor) Client() database.NoSQLClient {
	panic("unimplemented")
}

// Close implements database.NoSQL.
func (e *Executor) Close(ctx context.Context) error {
	panic("unimplemented")
}

// Open implements database.NoSQL.
func (e *Executor) Open(ctx context.Context) error {
	panic("unimplemented")
}

// Ping implements database.NoSQL.
func (e *Executor) Ping(ctx context.Context) error {
	panic("unimplemented")
}

// Request implements database.NoSQL.
func (e *Executor) Request(ctx context.Context, requests ...database.Request) (any, error) {
	panic("unimplemented")
}

// RequestWithHandler implements database.NoSQL.
func (e *Executor) RequestWithHandler(ctx context.Context, handler func(response database.Response, err error), requests ...database.Request) error {
	panic("unimplemented")
}

// Type implements database.NoSQL.
func (e *Executor) Type() database.Type {
	panic("unimplemented")
}
