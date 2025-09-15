package ui

import (
	"github.com/charmbracelet/huh"
	"github.com/ctrl-alt-boop/dribble/database"
)

type QueryForm struct {
	huh.Form

	intent *database.Intent
}

func CreateQueryForm(method database.OperationType, target string) *QueryForm {
	return &QueryForm{
		intent: &database.Intent{
			Type: method,
		},
	}
}
