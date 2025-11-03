// Package datastore contain widget to dribbler to io structs
package datastore

import (
	"context"

	"github.com/ctrl-alt-boop/dribble/database"
)

type (
	// DribbleRequestMsg should be used by widgets to command the main Model to request from the dribble api
	DribbleRequestMsg struct {
		TargetID int // -1 will temporarily mean all targets
		Request  database.Request
		Context  context.Context
	}
)
