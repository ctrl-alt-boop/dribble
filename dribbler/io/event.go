package io

import (
	"github.com/ctrl-alt-boop/dribble"
	"github.com/ctrl-alt-boop/dribble/internal/connection"
)

type (
	DribbleEventMsg struct {
		Type dribble.EventType
		Args any
		Err  error
	}

	ConnectMsg struct {
		Settings *connection.Settings
	}
)
