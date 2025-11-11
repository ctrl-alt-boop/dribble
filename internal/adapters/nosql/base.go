package nosql

import (
	"context"
	"fmt"

	"github.com/ctrl-alt-boop/dribble/datasource"
	"github.com/ctrl-alt-boop/dribble/internal/adapters"
	"github.com/ctrl-alt-boop/dribble/request"
)

var SupportedModels = map[datasource.SourceType]func(datasource.Namer) datasource.DataSource{}

// const SourceType datasource.SourceType = "sql"

func init() {}

type Base struct {
	adapters.BaseDatabase
	DSN  datasource.Namer
	Self datasource.NoSQL
	DB   datasource.NoSQLClient
}

// Type implements datasource.Database.
func (b *Base) Type() datasource.Type {
	return b.DSN.Type()
}

// Open implements datasource.Database.
func (b *Base) Open(ctx context.Context) error {
	return b.DB.Open(ctx)
}

// Ping implements datasource.Database.
func (b *Base) Ping(ctx context.Context) error {
	return b.DB.Ping(ctx)
}

// Close implements datasource.Database.
func (b *Base) Close(ctx context.Context) error {
	return b.DB.Close(ctx)
}

// IsClosed implements datasource.Database.
func (b *Base) IsClosed() bool {
	// TODO: Implement this for NoSQL clients
	return false
}

// Client implements datasource.NoSQL.
func (b *Base) Client() datasource.NoSQLClient {
	return b.DB
}

// Request implements datasource.Database.
func (b *Base) Request(ctx context.Context, req datasource.Request) (any, error) {
	return b.handlePrefab(ctx, req)
}

func (b *Base) handlePrefab(ctx context.Context, req datasource.Request) (any, error) {
	switch r := req.(type) {
	case request.ReadDatabaseNames, *request.ReadDatabaseNames:
		// For NoSQL, we might not have a direct "database list" concept like SQL.
		// This would typically return the current database name if connected, or a list of accessible databases.
		// For now, return the DBName from the DSN if available.
		if b.DSN.Info() != "" { // Assuming Info() contains the DBName or a relevant identifier
			return []string{b.DSN.Info()}, nil
		}
		return []string{}, nil
	case request.ReadTableNames, *request.ReadTableNames:
		// For NoSQL, "tables" are usually "collections".
		// This would require a client-specific implementation to list collections.
		return nil, fmt.Errorf("listing collections not yet implemented for %s", fmt.Sprintf("%s: %v", b.DSN.Type(), r))
	case request.ReadColumnNames, *request.ReadColumnNames:
		// For NoSQL, "columns" are "fields" within documents.
		// This would require inspecting a sample document from a collection.
		return nil, fmt.Errorf("listing fields not yet implemented for %s", fmt.Sprintf("%s: %v", b.DSN.Type(), r))
	default:
		return nil, fmt.Errorf("unknown prefab request for NoSQL: %T", req)
	}
}
