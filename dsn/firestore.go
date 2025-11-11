package dsn

import (
	"fmt"

	"github.com/ctrl-alt-boop/dribble/datasource"
)

var _ datasource.Namer = (*Firestore)(nil)

type Firestore struct {
	ProjectID string `json:"project_id"`
	Database  string `json:"database"`
}

// SourceType implements datasource.Namer.
func (f *Firestore) SourceType() datasource.SourceType {
	panic("unimplemented")
}

// Info implements database.DataSourceNamer.
func (f Firestore) Info() string {
	dbString := ""
	if f.Database != "" {
		dbString = fmt.Sprintf(" (%s)", f.Database)
	}
	return fmt.Sprintf("Firestore: %s%s", f.ProjectID, dbString)
}

// Type implements database.DataSourceNamer.
func (f Firestore) Type() datasource.Type {
	return datasource.Firestore
}

func (f Firestore) DSN() string {
	// Firestore DSN format: projects/<PROJECT_ID>/databases/<DATABASE_ID>
	if f.ProjectID == "" {
		return ""
	}
	dsn := fmt.Sprintf("projects/%s", f.ProjectID)
	if f.Database != "" {
		dsn += fmt.Sprintf("/databases/%s", f.Database)
	} else {
		dsn += "/databases/(default)"
	}
	return dsn
}

// FirestoreOption defines a function that configures a Firestore DSN.
type FirestoreOption func(*Firestore)

// FirestoreDSN creates a new FirestoreDSN with the given options.
func FirestoreDSN(projectID string, opts ...FirestoreOption) *Firestore {
	dsn := &Firestore{
		ProjectID: projectID,
	}
	for _, opt := range opts {
		opt(dsn)
	}
	return dsn
}

// FirestoreDatabase sets the database for the DSN.
func FirestoreDatabase(database string) FirestoreOption {
	return func(f *Firestore) {
		f.Database = database
	}
}
