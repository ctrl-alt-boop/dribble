package mysql

import (
	"text/template"

	"github.com/ctrl-alt-boop/dribble/database"
	_ "github.com/go-sql-driver/mysql"
)

func init() {
	database.DBTypes.Register("SQL", "mysql")
}

var _ database.SQLDialect = &MySQL{}

type MySQL struct{}

func NewMySQLDriver() (*MySQL, error) {
	driver := &MySQL{}
	return driver, nil
}

// // Capabilities implements database.Dialect.
func (m *MySQL) Capabilities() []database.Capabilities {
	return []database.Capabilities{}
}

const connectionStringTemplate = "{{if .Username}}{{.Username}}{{if .Password}}:{{.Password}}{{end}}@{{end}}tcp({{.Addr}}{{if .Port}}:{{.Port}}{{end}})/{{.DBName}}"

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
func (m *MySQL) GetPrefab(request database.Request) (string, []any, error) {
	panic("unimplemented")
}

// Name implements database.SQLDialect.
func (m *MySQL) Name() string {
	return "mysql"
}

// RenderRequest implements database.SQLDialect.
func (m *MySQL) RenderRequest(request database.Request) (string, []any, error) {
	panic("unimplemented")
}

// ResolveType implements database.SQLDialect.
func (m *MySQL) ResolveType(dbType string, value []byte) (any, error) {
	panic("unimplemented")
}
