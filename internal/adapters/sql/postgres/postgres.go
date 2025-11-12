package postgres

import (
	"fmt"
	"text/template"

	"github.com/ctrl-alt-boop/dribble/datasource"

	"github.com/ctrl-alt-boop/dribble/internal/adapters/sql"
	"github.com/ctrl-alt-boop/dribble/request"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

const SourceType datasource.SourceType = "postgres"

func init() {
	datasource.Register(datasource.Adapter{
		Name: "PostgreSQL",
		Type: SourceType,
		Properties: map[string]string{
			"dialect": "postgres",
		},
		Capabilities: []datasource.Capability{
			datasource.SupportsJSON,
			datasource.SupportsJSONB,
		},
		Metadata: datasource.Metadata{
			SourceType:  datasource.SourceTypeSQL,
			StorageType: datasource.IsDatabase,
		},
		FactoryFunc: New,
	})
}

var _ datasource.DataSource = (*Postgres)(nil)

type Postgres struct {
	sql.Base
}

func New(dsn datasource.Namer) datasource.DataSource {
	p := &Postgres{
		Base: sql.NewBase(dsn),
	}
	p.Self = p
	return p
}

// Name implements datasource.DataSource.
func (p *Postgres) Name() string {
	return "PostgreSQL"
}

// GoName implements datasource.DataSource.
func (p *Postgres) GoName() string {
	return "postgres"
}

// DriverName implements datasource.DataSource.
func (p *Postgres) DriverName() string {
	return "postgres"
}

// ModelType implements datasource.DataSource.
func (p *Postgres) ModelType() datasource.ModelType {
	return "Postgres"
}

// Capabilities implements datasource.SQLDialect.
func (p *Postgres) Capabilities() []datasource.Capability {
	return []datasource.Capability{
		datasource.SupportsJSON,
		datasource.SupportsJSONB,
	}
}

const connectionStringTemplate = "host={{.Addr}} port={{.Port}} user={{.Username}}{{if .Password}} password={{.Password}}{{end}}{{if .DBName}} dbname={{.DBName}}{{end}}{{with .Additional.sslmode}} sslmode={{.}}{{else}} sslmode=disable{{end}}"

func (p *Postgres) ConnectionStringTemplate() *template.Template {
	tmpl, err := template.New("connectionString").Parse(connectionStringTemplate)
	if err != nil {
		panic(err)
	}

	return tmpl
}

func (p *Postgres) ResolveType(dbType string, value []byte) (any, error) {
	switch dbType {
	case "UUID":
		return uuid.ParseBytes(value)
	default:
		return string(value), nil
	}
}

func (p *Postgres) GetPrefab(r datasource.Request) (string, []any, error) {
	switch r := r.(type) {
	case request.ReadDatabaseNames:
		return PrefabDatabases, nil, nil
	case request.ReadTableNames:
		return PrefabTables, nil, nil
	case request.ReadColumnNames:
		return PrefabColumns, []any{r.TableName}, nil
	case request.ReadCount:
		if r.DatabaseName != "" {
			return fmt.Sprintf(PrefabCountDBFormat, r.DatabaseName, r.TableName), nil, nil
		}
		return fmt.Sprintf(PrefabCountFormat, r.TableName), nil, nil
	default:
		return "", nil, fmt.Errorf("unknown prefab request: %T", r)
	}
}

// GetTemplate implements database.Dialect.
func (p *Postgres) GetTemplate(queryType datasource.RequestType) string {
	switch queryType {
	case datasource.Read:
		return ""
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
func (p *Postgres) IncreamentPlaceholder() string {
	panic("unimplemented")
}

// Quote implements database.Dialect.
func (p *Postgres) Quote(value string) string {
	return fmt.Sprintf(`"%s"`, value)
}

// QuoteRune implements database.Dialect.
func (p *Postgres) QuoteRune() rune {
	return '"'
}

// RenderCurrentTimestamp implements database.Dialect.
func (p *Postgres) RenderCurrentTimestamp() string {
	return "NOW()"
}

// RenderPlaceholder implements database.Dialect.
func (p *Postgres) RenderPlaceholder(index int) string {
	return fmt.Sprintf("$%d", index)
}

// RenderTypeCast implements database.Dialect.
func (p *Postgres) RenderTypeCast() string { // FIXME
	return "::"
}

// RenderValue implements database.Dialect.
func (p *Postgres) RenderValue(value any) string {
	return fmt.Sprintf("%v", value)
}

const (
	PrefabCurrentDatabase = "SELECT current_database()"
	PrefabDatabases       = "SELECT datname FROM pg_database"
	PrefabTables          = "SELECT table_name FROM information_schema.tables WHERE table_schema = 'public'"
	PrefabColumns         = "SELECT column_name FROM information_schema.columns WHERE table_name = $1"
	PrefabColumnsFormat   = "SELECT column_name FROM information_schema.columns WHERE table_name = '%s'"
	PrefabCountFormat     = "SELECT COUNT(*) FROM %s"
	PrefabCountDBFormat   = "SELECT COUNT(*) FROM %s.%s"
)
