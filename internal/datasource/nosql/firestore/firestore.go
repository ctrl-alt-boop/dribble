package firestore

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/ctrl-alt-boop/dribble/database"
	"github.com/ctrl-alt-boop/dribble/nosql"
)

func init() {
	database.DBTypes.Register("NoSQL", "firestore")
}

var _ database.NoSQLClient = (*Firestore)(nil)

type Firestore struct {
	client *firestore.Client

	clientProperties *nosql.FirestoreClientProperties
}

func NewFirestoreClient() (*Firestore, error) {
	// if firestoreClientProperties.Name == "" {
	// 	return nil, errors.New("firestore target client name is required")
	// }
	driver := &Firestore{
		// clientProperties: firestoreClientProperties,
	}
	return driver, nil
}

// Capabilities implements database.Dialect.
func (f *Firestore) Capabilities() []database.Capabilities {
	return []database.Capabilities{
		database.SupportsSQLLike,
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
func (f *Firestore) Client() database.NoSQLClient {
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
func (f *Firestore) Request(ctx context.Context, requests ...database.Request) (any, error) {
	panic("unimplemented")
}

// RequestWithHandler implements database.NoSQL.
func (f *Firestore) RequestWithHandler(ctx context.Context, handler func(response database.Response, err error), requests ...database.Request) error {
	panic("unimplemented")
}

// Type implements database.NoSQL.
func (f *Firestore) Type() database.Type {
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
