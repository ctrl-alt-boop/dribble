package sql

import (
	"reflect"
	"strings"

	"github.com/ctrl-alt-boop/dribble/database"
)

func FromString(query string, args ...any) *database.Intent {
	if strings.ContainsAny(query, ";") {
		panic("query cannot contain a semicolon/multiple statements") // FIXME: return error?
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
		panic("statement not supported") // FIXME: return error?
	}

	return &database.Intent{
		Type:          operationType,
		OperationKind: reflect.TypeOf(query).Kind(),
		Operation:     query,
		Args:          args,
	}
}
