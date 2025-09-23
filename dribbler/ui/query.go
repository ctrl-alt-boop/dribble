package ui

import (
	"github.com/charmbracelet/huh"
	"github.com/ctrl-alt-boop/dribble/database"
	"github.com/ctrl-alt-boop/dribble/request"
)

type QueryForm struct {
	huh.Form

	intent database.Request
}

func CreateQueryForm(req database.Request, target string) *QueryForm {
	return &QueryForm{
		intent: &request.Intent{},
	}
}
