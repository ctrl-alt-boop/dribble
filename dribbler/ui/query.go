package ui

import (
	"github.com/charmbracelet/huh"
	"github.com/ctrl-alt-boop/dribble/playbook/database"
)

type QueryForm struct {
	huh.Form

	statement *database.Statement
}

func CreateQueryForm(method database.SqlMethod, table string) *QueryForm {
	return &QueryForm{
		statement: &database.Statement{
			Method: method,
			Table:  table,
		},
	}
}
