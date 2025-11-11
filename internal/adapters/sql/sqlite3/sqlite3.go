package sqlite3

import (
	"fmt"
	"text/template"

	_ "embed"

	"github.com/ctrl-alt-boop/dribble/datasource"

	"github.com/ctrl-alt-boop/dribble/internal/adapters/sql"
	"github.com/ctrl-alt-boop/dribble/request"
	_ "github.com/mattn/go-sqlite3"
)

const SourceType datasource.SourceType = "sqlite3"

func init() {
	datasource.Register(datasource.Adapter{
		Name: "SQLite3",
		Type: SourceType,
		Properties: map[string]string{
			"dialect": "sqlite3",
		},
		Capabilities: []datasource.Capability{},
		Metadata: datasource.Metadata{
			SourceType:  datasource.SourceTypeSQL,
			StorageType: datasource.IsFile,
		},
		FactoryFunc: func(dsn datasource.Namer) datasource.DataSource {
			return &SQLite3{
				Base: sql.Base{
					DSN: dsn,
				},
			}
		},
	})
}

var _ datasource.DataSource = (*SQLite3)(nil)

type SQLite3 struct {
	sql.Base
}

// Name implements datasource.DataSource.
func (s *SQLite3) Name() string {
	return "SQLite3"
}

// GoName implements datasource.DataSource.
func (s *SQLite3) GoName() string {
	return "sqlite3"
}

// ModelType implements datasource.DataSource.
func (s *SQLite3) ModelType() datasource.ModelType {
	return datasource.ModelType("SQLite3")
}

// Capabilities implements database.Dialect.
func (s *SQLite3) Capabilities() []datasource.Capability {
	return []datasource.Capability{}
}

// example: file:test.db?cache=shared&mode=memory
// _auth_user
// _auth_pass
// _auth_crypt
// mode ro, rw, rwc, memory
// cache shared, private
const connectionStringTemplate = "file:{{.Addr}}{{if .DBName}}?_auth_user={{.Username}}&_auth_pass={{.Password}}&_auth_crypt={{.DBName}}{{end}}{{with .Additional.mode}}&mode={{.}}{{else}}{{end}}{{with .Additional.cache}}&cache={{.}}{{else}}{{end}}"

// ConnectionString implements database.Driver.
func (s *SQLite3) ConnectionStringTemplate() *template.Template {
	tmpl, err := template.New("connectionString").Parse(connectionStringTemplate)
	if err != nil {
		panic(err)
	}
	return tmpl
}

const (
	PrefabDatabases = "PRAGMA database_list"
	PrefabTables    = "SELECT name FROM sqlite_master WHERE type='table'"
	PrefabColumns   = "PRAGMA table_info(?)"
)

// GetPrefab implements database.SQLDialect.
func (s *SQLite3) GetPrefab(r datasource.Request) (string, []any, error) {
	switch r := r.(type) {
	case request.ReadDatabaseNames, *request.ReadDatabaseNames:
		return PrefabDatabases, nil, nil
	case request.ReadTableNames, *request.ReadTableNames:
		return PrefabTables, nil, nil
	case request.ReadColumnNames:
		return PrefabColumns, []any{r.TableName}, nil
	case *request.ReadColumnNames:
		return PrefabColumns, []any{r.TableName}, nil
	default:
		return "", nil, fmt.Errorf("unknown prefab request: %T", r)
	}
}

func (s *SQLite3) ResolveType(dbType string, value []byte) (any, error) {
	// |go        | sqlite3           |
	// |----------|-------------------|
	// |nil       | null              |
	// |int       | integer           |
	// |int64     | integer           |
	// |float64   | float             |
	// |bool      | integer           |
	// |[]byte    | blob              |
	// |string    | text              |
	// |time.Time | timestamp/datetime|
	return string(value), nil
}

//go:embed templates/select.tmpl
var selectQueryTemplate string

// GetTemplate implements database.Dialect.
func (s *SQLite3) GetTemplate(queryType datasource.RequestType) string {
	switch queryType {
	case datasource.Read:
		return selectQueryTemplate
	case datasource.Create:
		return "" // DefaultSQLInsertTemplate
	case datasource.Update:
		return "" // DefaultSQLUpdateTemplate
	case datasource.Delete:
		return "" // DefaultSQLDeleteTemplate
	default:
		return ""
	}
}

// IncreamentPlaceholder implements database.Dialect.
func (s *SQLite3) IncreamentPlaceholder() string {
	panic("unimplemented")
}

// Quote implements database.Dialect.
func (s *SQLite3) Quote(value string) string {
	return fmt.Sprintf(`"%s"`, value)
}

// QuoteRune implements database.Dialect.
func (s *SQLite3) QuoteRune() rune {
	return '"'
}

// RenderCurrentTimestamp implements database.Dialect.
func (s *SQLite3) RenderCurrentTimestamp() string {
	return "NOW()"
}

// RenderPlaceholder implements database.Dialect.
func (s *SQLite3) RenderPlaceholder(index int) string {
	return "?"
}

// RenderTypeCast implements database.Dialect.
func (s *SQLite3) RenderTypeCast() string { // FIXME
	return "::"
}

// RenderValue implements database.Dialect.
func (s *SQLite3) RenderValue(value any) string {
	return fmt.Sprintf("%v", value)
}
