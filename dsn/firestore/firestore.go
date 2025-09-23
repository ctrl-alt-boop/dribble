package firestore

import (
	"fmt"

	"github.com/ctrl-alt-boop/dribble/database"
)

type FirestoreDSN struct {
	ProjectID string `json:"project_id"`
	Database  string `json:"database"`
}

// Type implements database.DataSourceNamer.
func (f FirestoreDSN) Type() database.Type {
	return database.Firestore
}

func (f FirestoreDSN) DSN() string {
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

// DSNOption defines a function that configures a FirestoreDSN.
type DSNOption func(*FirestoreDSN)

// NewDSN creates a new FirestoreDSN with the given options.
func NewDSN(projectID string, opts ...DSNOption) *FirestoreDSN {
	dsn := &FirestoreDSN{
		ProjectID: projectID,
	}
	for _, opt := range opts {
		opt(dsn)
	}
	return dsn
}

// WithDatabase sets the database for the DSN.
func WithDatabase(database string) DSNOption {
	return func(f *FirestoreDSN) {
		f.Database = database
	}
}
