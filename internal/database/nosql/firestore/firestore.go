package firestore

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/ctrl-alt-boop/dribble/database"
	"github.com/ctrl-alt-boop/dribble/target"
)

func init() {
	database.DBTypes.Register("NoSQL", "firestore")
}

var _ database.NoSQL = &Firestore{}

type Firestore struct {
	client *firestore.Client
}

func NewFirestoreDriver(target *target.Target) (*Firestore, error) {
	driver := &Firestore{}
	return driver, nil
}

// Capabilities implements database.Dialect.
func (f *Firestore) Capabilities() []database.Capabilities {
	return []database.Capabilities{
		database.SupportsSQLLike,
	}
}

// Open implements database.NoSQLClient.
func (f *Firestore) Open(ctx context.Context, target *target.Target) error { // I'll use target.Name for the time being
	var client *firestore.Client
	var err error
	if target.DBName != "" {
		client, err = firestore.NewClientWithDatabase(ctx, target.Name, target.DBName)
	} else {
		client, err = firestore.NewClient(ctx, target.Name)
	}
	if err != nil {
		return err
	}
	f.client = client
	return nil
}

// Ping implements database.NoSQLClient.
func (f *Firestore) Ping(ctx context.Context) error {
	return nil
}

// Close implements database.NoSQLClient.
func (f *Firestore) Close(_ context.Context) error {
	return f.client.Close()
}

// Read implements database.NoSQLClient.
func (f *Firestore) Read(any) {
	panic("unimplemented")
}

// ReadMany implements database.NoSQLClient.
func (f *Firestore) ReadMany(any) {
	panic("unimplemented")
}

// Create implements database.NoSQLClient.
func (f *Firestore) Create(any) {
	panic("unimplemented")
}

// Update implements database.NoSQLClient.
func (f *Firestore) Update(any) {
	panic("unimplemented")
}

// Delete implements database.NoSQLClient.
func (f *Firestore) Delete(any) {
	panic("unimplemented")
}

// ConnectionString implements database.Driver.
func (f *Firestore) ConnectionString(target *target.Target) string {
	if target.DBName != "" {
		return target.Name + "/" + target.DBName
	}
	return target.Name
}

func (f *Firestore) Dialect() database.SQLDialect {
	panic("unimplemented")
}

// RenderIntent implements database.Driver.
func (f *Firestore) RenderIntent(intent *database.Intent) (string, error) {
	panic("unimplemented")
}

// ResolveType implements database.Dialect.
func (f *Firestore) ResolveType(dbType string, value []byte) (any, error) {
	panic("unimplemented")
}
