package query

import (
	"github.com/ctrl-alt-boop/dribble/internal/database"
)

type QueryRenderer interface {
	Render(dialect database.Dialect) string
}
