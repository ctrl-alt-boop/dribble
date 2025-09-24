package mysql

import (
	"fmt"
	"text/template"

	"github.com/ctrl-alt-boop/dribble/database"
	"github.com/ctrl-alt-boop/dribble/request"
	_ "github.com/go-sql-driver/mysql"
)

func init() {
	database.DBTypes.Register("SQL", "mysql")
}

var _ database.SQLDialect = (*MySQL)(nil)

type MySQL struct{}

func NewMySQLDriver() (*MySQL, error) {
	driver := &MySQL{}
	return driver, nil
}

// // Capabilities implements database.Dialect.
func (m *MySQL) Capabilities() []database.Capabilities {
	return []database.Capabilities{}
}

const connectionStringTemplate = "{{if .Username}}{{.Username}}{{if .Password}}:{{.Password}}{{end}}@{{end}}tcp({{.Addr}}{{if .Port}}:{{.Port}}{{end}})/{{.DBName}}?{{with .Additional.allowCleartextPasswords}}allowCleartextPasswords={{.}}{{else}}&allowCleartextPasswords=false{{end}}"

// ConnectionStringTemplate implements database.SQLDialect.
// user:password@tcp(host:port)/dbname
func (m *MySQL) ConnectionStringTemplate() *template.Template {
	tmpl, err := template.New("connectionString").Parse(connectionStringTemplate)
	if err != nil {
		panic(err)
	}
	return tmpl
}

// GetPrefab implements database.SQLDialect.
func (m *MySQL) GetPrefab(r database.Request) (string, []any, error) {
	switch r := r.(type) {
	case request.ReadDatabaseNames, *request.ReadDatabaseNames:
		return PrefabDatabases, nil, nil
	case request.ReadTableNames:
		return PrefabTables, []any{r.DatabaseName}, nil
	case *request.ReadTableNames:
		return PrefabTables, []any{r.DatabaseName}, nil
	case request.ReadColumnNames:
		return PrefabColumns, []any{r.TableName}, nil
	case *request.ReadColumnNames:
		return PrefabColumns, []any{r.TableName}, nil
	default:
		return "", nil, fmt.Errorf("unknown prefab request: %T", r)
	}
}

// Name implements database.SQLDialect.
func (m *MySQL) Name() string {
	return "mysql"
}

// ResolveType implements database.SQLDialect.
func (m *MySQL) ResolveType(dbType string, value []byte) (any, error) {
	panic("unimplemented")
}
