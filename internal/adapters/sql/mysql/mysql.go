package mysql

import (
	_ "embed"
	"fmt"
	"text/template"

	"github.com/ctrl-alt-boop/dribble/datasource"

	"github.com/ctrl-alt-boop/dribble/internal/adapters/sql"
	"github.com/ctrl-alt-boop/dribble/request"
	_ "github.com/go-sql-driver/mysql"
)

const SourceType datasource.SourceType = "mysql"

func init() {
	datasource.Register(datasource.Adapter{
		Name: "MySQL",
		Type: SourceType,
		Properties: map[string]string{
			"dialect": "mysql",
		},
		Capabilities: []datasource.Capability{},
		Metadata: datasource.Metadata{
			SourceType:  datasource.SourceTypeSQL,
			StorageType: datasource.IsDatabase,
		},
		FactoryFunc: New,
	})
}

var _ datasource.DataSource = (*MySQL)(nil)

type MySQL struct {
	sql.Base
}

func New(dsn datasource.Namer) datasource.DataSource {
	m := &MySQL{
		Base: sql.NewBase(dsn),
	}
	m.Self = m
	return m
}

// Name implements datasource.DataSource.
func (m *MySQL) Name() string {
	return "MySQL"
}

// DriverName implements datasource.DataSource.
func (m *MySQL) DriverName() string {
	return "mysql"
}

// ModelType implements datasource.DataSource.
func (m *MySQL) ModelType() datasource.ModelType {
	return datasource.ModelType("MySQL")
}

func NewDriver() *MySQL {
	return &MySQL{}
}

// // Capabilities implements database.Dialect.
func (m *MySQL) Capabilities() []datasource.Capability {
	return []datasource.Capability{}
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

const (
	PrefabDatabases = "SHOW DATABASES"
	PrefabTables    = "SHOW TABLES FROM ?"
	PrefabColumns   = "SHOW COLUMNS FROM ?"
)

// GetPrefab implements database.SQLDialect.
func (m *MySQL) GetPrefab(r datasource.Request) (string, []any, error) {
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

// ResolveType implements database.SQLDialect.
func (m *MySQL) ResolveType(dbType string, value []byte) (any, error) {
	panic("unimplemented")
}

//go:embed templates/select.tmpl
var selectQueryTemplate string

// GetTemplate implements database.Dialect.
func (m *MySQL) GetTemplate(queryType datasource.RequestType) string {
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
func (m *MySQL) IncreamentPlaceholder() string {
	panic("unimplemented")
}

// Quote implements database.Dialect.
func (m *MySQL) Quote(value string) string {
	return fmt.Sprintf(`"%s"`, value)
}

// QuoteRune implements database.Dialect.
func (m *MySQL) QuoteRune() rune {
	return '"'
}

// RenderCurrentTimestamp implements database.Dialect.
func (m *MySQL) RenderCurrentTimestamp() string {
	return "NOW()"
}

// RenderPlaceholder implements database.Dialect.
func (m *MySQL) RenderPlaceholder(index int) string {
	return "?"
}

// RenderTypeCast implements database.Dialect.
func (m *MySQL) RenderTypeCast() string { // FIXME
	return "::"
}

// RenderValue implements database.Dialect.
func (m *MySQL) RenderValue(value any) string {
	return fmt.Sprintf("%v", value)
}
