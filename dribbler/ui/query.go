package ui

import (
	"github.com/charmbracelet/huh"
	"github.com/ctrl-alt-boop/dribble/database"
)

type QueryForm struct {
	huh.Form

	intent *database.QueryIntent
}

func CreateQueryForm(method database.QueryType, target string) *QueryForm {
	return &QueryForm{
		intent: &database.QueryIntent{
			Type:       method,
			TargetName: target,
		},
	}
}
