package sqlite3

import (
	"fmt"
	"text/template"

	"github.com/ctrl-alt-boop/dribble/database"
	"github.com/ctrl-alt-boop/dribble/request"
	_ "github.com/mattn/go-sqlite3"
)

func init() {
	database.DBTypes.Register("SQL", "sqlite3")
}

var _ database.SQLDialect = &SQLite3{}

type SQLite3 struct{}

func NewSQLite3Driver() (*SQLite3, error) {
	driver := &SQLite3{}
	return driver, nil
}

// Capabilities implements database.Dialect.
func (s *SQLite3) Capabilities() []database.Capabilities {
	return []database.Capabilities{
		database.IsFile,
	}
}

// Name implements database.SQLDialect.
func (s *SQLite3) Name() string {
	return "sqlite3"
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

// GetPrefab implements database.SQLDialect.
func (s *SQLite3) GetPrefab(r database.Request) (string, []any, error) {
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
