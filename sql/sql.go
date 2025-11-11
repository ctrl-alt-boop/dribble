package sql

import (
	"errors"
	"strings"

	"github.com/ctrl-alt-boop/dribble/datasource"
	"github.com/ctrl-alt-boop/dribble/request"
)

var (
	ErrUnknownOperation   = errors.New("unknown operation")
	ErrMultipleStatements = errors.New("multiple statements in query")
)

func FromString(query string, args ...any) (datasource.Request, error) {
	if strings.ContainsAny(query, ";") {
		return nil, ErrMultipleStatements
	}

	query = strings.TrimSpace(query)
	operationType := datasource.NoOp

	switch {
	case strings.HasPrefix(query, "SELECT"):
		operationType = datasource.Read
	case strings.HasPrefix(query, "INSERT"):
		operationType = datasource.Create
	case strings.HasPrefix(query, "UPDATE"):
		operationType = datasource.Update
	case strings.HasPrefix(query, "DELETE"):
		operationType = datasource.Delete
	default:
		return nil, ErrUnknownOperation
	}

	return &request.Intent{
		Type:      operationType,
		Operation: query,
		Args:      args,
	}, nil
}

// type ConnectionProperties struct {
// 	Addr     string
// 	Port     int
// 	DBName   string
// 	Username string
// 	Password string

// 	Extra map[string]string
// }
