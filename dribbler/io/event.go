package io

import (
	"github.com/ctrl-alt-boop/dribble"
	"github.com/ctrl-alt-boop/dribble/database"
)

type (
	DribbleEventMsg struct {
		Type dribble.Response
		Args any
		Err  error
	}

	ConnectMsg struct {
		Target *database.Target
	}
)
