package io

import (
	"github.com/ctrl-alt-boop/gooldb/internal/app/gooldb"
	"github.com/ctrl-alt-boop/gooldb/pkg/connection"
)

type (
	GoolDbEventMsg struct {
		Type gooldb.EventType
		Args any
		Err  error
	}

	ConnectMsg struct {
		Settings *connection.Settings
	}
)
