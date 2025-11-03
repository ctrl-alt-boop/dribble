package dribbler

import (
	"context"
	"fmt"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/ctrl-alt-boop/dribble/database"
	"github.com/ctrl-alt-boop/dribble/request"
	"github.com/ctrl-alt-boop/dribble/target"
	"github.com/ctrl-alt-boop/dribbler/datastore"
)

// TargetingError for dribble requests
type TargetingError string

// Error implements error
func (e TargetingError) Error() string {
	return string(e)
}

// DribbleRequestError for dribble requests
type DribbleRequestError error

// NewTargetingError creates a new TargetingError
func NewTargetingError(targetName string) error {
	if targetName == "" {
		return TargetingError("Target name cannot be empty")
	}
	msg := fmt.Sprintf("Target %s not found", targetName)
	return TargetingError(msg)
}

// NewDribbleRequestError creates a new DribbleRequestError
func NewDribbleRequestError(err error) tea.Cmd {
	return func() tea.Msg {
		return DribbleRequestError(err)
	}
}

// Target tries to get target by name from the dribble client
func (m Dribbler) Target(targetName string) (*DribbleRequester, error) {
	target, ok := m.dribbleClient.Target(targetName)
	if !ok {
		return nil, NewTargetingError(targetName)
	}
	return &DribbleRequester{
		target: target,
	}, nil
}

// DribbleResponse is used to either get the channel or collect it all to an iterator
type DribbleResponse struct {
	responseChan chan *request.Response
}

// All collects all Responses into a list
func (d DribbleResponse) All() []*request.Response {
	responses := []*request.Response{}
	for response := range d.responseChan {
		responses = append(responses, response)
	}
	return responses
}

// Channel gets the response channel
func (d DribbleResponse) Channel() chan *request.Response {
	return d.responseChan
}

// DribbleRequester is used to fluently enable requests to dribble api
type DribbleRequester struct {
	target *target.Target
}

// Request is used to create a tea.Cmd that will run a request against the dribble api
func (d DribbleRequester) Request(ctx context.Context, request database.Request) tea.Cmd {
	responseChan, err := d.target.Request(ctx, request)
	if err != nil {
		return func() tea.Msg {
			return NewDribbleRequestError(err)
		}
	}

	return func() tea.Msg {
		return DribbleResponse{
			responseChan: responseChan,
		}
	}
}

func (m Dribbler) handleDribbleRequestMsg(msg datastore.DribbleRequestMsg) tea.Cmd {
	if msg.TargetID == 0 {
		return NewDribbleRequestError(NewTargetingError("0"))
	}
	var responseChan chan *request.Response
	var err error

	if msg.TargetID == -1 { // Target all
		responseChan, err = m.dribbleClient.RequestForAll(msg.Context, msg.Request)
	} else {
		responseChan, err = m.dribbleClient.Request(msg.Context, fmt.Sprint(msg.TargetID), msg.Request) // FIXME: not fmt.Sprint
	}

	if err != nil {
		return NewDribbleRequestError(err)
	}

	return func() tea.Msg {
		return DribbleResponse{
			responseChan: responseChan,
		}
	}
}
