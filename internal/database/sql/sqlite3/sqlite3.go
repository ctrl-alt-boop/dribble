package sqlite3

import (
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

const connectionStringTemplate = "Data Source={{.Addr}};Version=3{{if .Password}};Password={{.Password}}{{end}}"

// ConnectionString implements database.Driver.
// Data Source=c:\mydb.db;Version=3;Password=myPassword
func (s *SQLite3) ConnectionStringTemplate() *template.Template {
	tmpl, err := template.New("connectionString").Parse(connectionStringTemplate)
	if err != nil {
		panic(err)
	}
	return tmpl
}

// GetPrefab implements database.SQLDialect.
func (s *SQLite3) GetPrefab(request database.Request) (string, []any, error) {
	panic("unimplemented")
}

// RenderRequest implements database.SQLDialect.
func (s *SQLite3) RenderRequest(request database.Request) (string, []any, error) {
	panic("unimplemented")
}

// RenderIntent implements database.Driver.
func (s *SQLite3) RenderIntent(intent *request.Intent) (string, error) {
	panic("unimplemented")
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
