package postgres

import (
	"fmt"
	"html/template"
	"strings"

	"github.com/ctrl-alt-boop/dribble/database"
	"github.com/ctrl-alt-boop/dribble/sql"
	"github.com/google/uuid"

	_ "github.com/lib/pq"
)

var _ database.Driver = &Postgres{}

type Postgres struct{}

func NewPostgresDriver(target *database.Target) (*Postgres, error) {
	driver := &Postgres{}
	return driver, nil
}

// Capabilities implements database.Dialect.
func (p *Postgres) Capabilities() []database.Capabilities {
	return []database.Capabilities{
		database.SupportsJson,
		database.SupportsJsonB,
	}
}

// ConnectionString implements database.Driver.
func (p *Postgres) ConnectionString(target *database.Target) string {
	connString := ""
	connString += fmt.Sprintf("host=%s ", target.Ip)
	if target.Port == 0 {
		target.Port = 5432
	}
	connString += fmt.Sprintf("port=%d ", target.Port)
	connString += fmt.Sprintf("user=%s ", target.Username)
	connString += fmt.Sprintf("password=%s ", target.Password)
	if target.DBName == "" {
		target.DBName = "postgres"
	}
	connString += fmt.Sprintf("dbname=%s ", target.DBName)
	sslmode, ok := target.AdditionalSettings["sslmode"]
	if !ok {
		sslmode = "disable"
	}
	connString += fmt.Sprintf("sslmode=%s ", sslmode)

	return connString
}

func (p *Postgres) Dialect() database.Dialect {
	return p
}

func (p *Postgres) RenderIntent(intent *database.Intent) (string, error) {
	var queryStringTemplate string
	operation := intent.Operation
	switch operation.(type) {
	case sql.SelectQuery:
		queryStringTemplate = p.GetTemplate(intent.Type)

	default:
		return "", fmt.Errorf("intent type %d not supported for postgres", intent.Type)
	}
	tmpl, err := template.New("query").Parse(queryStringTemplate)
	if err != nil {
		return "", fmt.Errorf("error parsing query template: %w", err)
	}
	var sb strings.Builder
	err = tmpl.Execute(&sb, operation)
	if err != nil {
		return "", fmt.Errorf("error executing query template: %w", err)
	}
	queryString := strings.TrimSpace(sb.String())
	return queryString, nil
}

// ResolveType implements database.Dialect.
func (p *Postgres) ResolveType(dbType string, value []byte) (any, error) {
	switch dbType {
	case "UUID":
		return uuid.ParseBytes(value)
	default:
		return string(value), nil
	}
}
