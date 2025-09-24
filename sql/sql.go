package sql

import (
	"errors"
	"strings"

	"github.com/ctrl-alt-boop/dribble/database"
	"github.com/ctrl-alt-boop/dribble/request"
)

var (
	ErrUnknownOperation   = errors.New("unknown operation")
	ErrMultipleStatements = errors.New("multiple statements in query")
)

func FromString(query string, args ...any) (database.Request, error) {
	if strings.ContainsAny(query, ";") {
		return nil, ErrMultipleStatements
	}

	query = strings.TrimSpace(query)
	operationType := database.NoOp

	switch {
	case strings.HasPrefix(query, "SELECT"):
		operationType = database.Read
	case strings.HasPrefix(query, "INSERT"):
		operationType = database.Create
	case strings.HasPrefix(query, "UPDATE"):
		operationType = database.Update
	case strings.HasPrefix(query, "DELETE"):
		operationType = database.Delete
	default:
		return nil, ErrUnknownOperation
	}

	return &request.Intent{
		Type:      operationType,
		Operation: query,
		Args:      args,
	}, nil
}

type ConnectionProperties struct {
	Addr     string
	Port     int
	DBName   string
	Username string
	Password string

	Extra map[string]string
}
