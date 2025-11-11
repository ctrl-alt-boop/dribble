package firestore

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/ctrl-alt-boop/dribble/datasource"
	"github.com/ctrl-alt-boop/dribble/nosql"
)

var _ datasource.NoSQLClient = (*Firestore)(nil)

type Firestore struct {
	client *firestore.Client

	clientProperties *nosql.FirestoreClientProperties
}

func NewFirestoreClient() *Firestore {
	return &Firestore{}
}

// Capabilities implements database.Dialect.
func (f *Firestore) Capabilities() []datasource.Capability {
	return []datasource.Capability{
		datasource.SupportsSQLLike,
	}
}

// SetConnectionProperties implements database.NoSQLClient.
func (f *Firestore) SetConnectionProperties(props map[string]string) {
	f.clientProperties = &nosql.FirestoreClientProperties{
		Name:         props["name"],
		DatabaseName: props["databaseName"],
	}
}

// Open implements database.NoSQLClient.
func (f *Firestore) Open(ctx context.Context) error { // I'll use target.Name for the time being
	var client *firestore.Client
	var err error
	if f.clientProperties.DatabaseName != "" {
		client, err = firestore.NewClientWithDatabase(ctx, f.clientProperties.Name, f.clientProperties.DatabaseName)
	} else {
		client, err = firestore.NewClient(ctx, f.clientProperties.Name)
	}
	if err != nil {
		return err
	}
	f.client = client
	return nil
}

// Client implements database.NoSQL.
func (f *Firestore) Client() datasource.NoSQLClient {
	panic("unimplemented")
}

// Close implements database.NoSQL.
func (f *Firestore) Close(ctx context.Context) error {
	panic("unimplemented")
}

// Ping implements database.NoSQL.
func (f *Firestore) Ping(ctx context.Context) error {
	panic("unimplemented")
}

// Request implements database.NoSQL.
func (f *Firestore) Request(ctx context.Context, requests ...datasource.Request) (any, error) {
	panic("unimplemented")
}

// RequestWithHandler implements database.NoSQL.
func (f *Firestore) RequestWithHandler(ctx context.Context, handler func(response datasource.Response, err error), requests ...datasource.Request) error {
	panic("unimplemented")
}

// Type implements database.NoSQL.
func (f *Firestore) Type() datasource.Type {
	panic("unimplemented")
}

// Create implements database.NoSQLClient.
func (f *Firestore) Create(any) {
	panic("unimplemented")
}

// Delete implements database.NoSQLClient.
func (f *Firestore) Delete(any) {
	panic("unimplemented")
}

// Read implements database.NoSQLClient.
func (f *Firestore) Read(any) {
	panic("unimplemented")
}

// ReadMany implements database.NoSQLClient.
func (f *Firestore) ReadMany(any) {
	panic("unimplemented")
}

// Update implements database.NoSQLClient.
func (f *Firestore) Update(any) {
	panic("unimplemented")
}
