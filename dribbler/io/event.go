package io

import (
	"github.com/ctrl-alt-boop/dribble/database"
	"github.com/ctrl-alt-boop/dribble/target"
)

type (
	DribbleEventMsg struct {
		Response database.Response
		Args     any
		Err      error
	}

	ConnectMsg struct {
		Target *target.Target
		DSN    database.DataSourceNamer
	}
)
