package postgres

import (
	"fmt"
	"text/template"

	"github.com/ctrl-alt-boop/dribble/database"
	"github.com/ctrl-alt-boop/dribble/request"
	"github.com/google/uuid"

	_ "github.com/lib/pq"
)

func init() {
	database.DBTypes.Register("SQL", "postgres")
}

var _ database.SQLDialect = &Postgres{}

type Postgres struct{}

// Name implements database.SQLDialect.
func (p *Postgres) Name() string {
	return "postgres"
}

func NewPostgresDriver() (*Postgres, error) {
	driver := &Postgres{}
	return driver, nil
}

// Capabilities implements database.Dialect.
func (p *Postgres) Capabilities() []database.Capabilities {
	return []database.Capabilities{
		database.SupportsJSON,
		database.SupportsJSONB,
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

func (p *Postgres) Dialect() database.SQLDialect {
	return p
}

func (p *Postgres) ResolveType(dbType string, value []byte) (any, error) {
	switch dbType {
	case "UUID":
		return uuid.ParseBytes(value)
	default:
		return string(value), nil
	}
}

func (p *Postgres) GetPrefab(r database.Request) (string, []any, error) {
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

// func (p *Postgres) ConnectionString() string {
// 	var target int
// 	connString := ""
// 	connString += fmt.Sprintf("host=%s ", target.Properties.Addr)
// 	if target.Properties.Port == 0 {
// 		target.Properties.Port = 5432
// 	}
// 	connString += fmt.Sprintf("port=%d ", target.Properties.Port)
// 	connString += fmt.Sprintf("user=%s ", target.Properties.Username)
// 	connString += fmt.Sprintf("password=%s ", target.Properties.Password)
// 	if target.Properties.DBName == "" {
// 		target.Properties.DBName = "postgres"
// 	}
// 	connString += fmt.Sprintf("dbname=%s ", target.Properties.DBName)
// 	sslmode, ok := target.Properties.Additional["sslmode"]
// 	if !ok {
// 		sslmode = "disable"
// 	}
// 	connString += fmt.Sprintf("sslmode=%s ", sslmode)

// 	return connString
// }
