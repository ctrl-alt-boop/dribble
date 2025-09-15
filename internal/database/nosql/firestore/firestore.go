package firestore

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/ctrl-alt-boop/dribble/database"
)

var _ database.Driver = &Firestore{}
var _ database.NoSQLClient = &Firestore{}

type Firestore struct {
	client *firestore.Client
}

func NewFirestoreDriver(target *database.Target) (*Firestore, error) {
	driver := &Firestore{}
	return driver, nil
}

// Capabilities implements database.Dialect.
func (m *Firestore) Capabilities() []database.Capabilities {
	return []database.Capabilities{
		database.SupportsSQL,
	}
}

// Open implements database.NoSQLClient.
func (m *Firestore) Open(ctx context.Context, target *database.Target) error { // I'll use target.Name for the time being
	var client *firestore.Client
	var err error
	if target.DBName == "" {
		client, err = firestore.NewClientWithDatabase(ctx, target.Name, target.DBName)
	} else {
		client, err = firestore.NewClient(ctx, target.Name)
	}
	if err != nil {
		return err
	}
	m.client = client
	return nil
}

// Close implements database.NoSQLClient.
func (m *Firestore) Close(_ context.Context) error {
	return m.client.Close()
}

// Read implements database.NoSQLClient.
func (m *Firestore) Read(any) {
	panic("unimplemented")
}

// ReadMany implements database.NoSQLClient.
func (m *Firestore) ReadMany(any) {
	panic("unimplemented")
}

// Create implements database.NoSQLClient.
func (m *Firestore) Create(any) {
	panic("unimplemented")
}

// Update implements database.NoSQLClient.
func (m *Firestore) Update(any) {
	panic("unimplemented")
}

// Delete implements database.NoSQLClient.
func (m *Firestore) Delete(any) {
	panic("unimplemented")
}

// ConnectionString implements database.Driver.
func (m *Firestore) ConnectionString(target *database.Target) string {
	if target.DBName != "" {
		return target.Name + "/" + target.DBName
	}
	return target.Name
}

func (m *Firestore) Dialect() database.Dialect {
	panic("unimplemented")
}

// RenderIntent implements database.Driver.
func (m *Firestore) RenderIntent(intent *database.Intent) (string, error) {
	panic("unimplemented")
}

// ResolveType implements database.Dialect.
func (m *Firestore) ResolveType(dbType string, value []byte) (any, error) {
	panic("unimplemented")
}
