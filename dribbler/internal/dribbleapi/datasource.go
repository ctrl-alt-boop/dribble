// Package dribbleapi contain
//
//   - messages for dribble requests
//   - error types for dribble requests
//   - a fluent api for dribble requests
package dribbleapi

import (
	"context"
	"fmt"

	tea "charm.land/bubbletea/v2"
	"github.com/ctrl-alt-boop/dribble/datasource"
	"github.com/ctrl-alt-boop/dribble/request"
	"github.com/ctrl-alt-boop/dribble/target"
)

func GetSupportedDataSourceNames() []string { // temp dribble.NewClient()
	supported := []string{}
	for _, s := range datasource.Adapters() {
		supported = append(supported, s.Name)
	}
	return supported
}

func GetSupportedSourceTypes() []datasource.SourceType { // temp dribble.NewClient()
	return datasource.AdapterTypes()
}

func GetSupportedDataSources() []datasource.Adapter { // temp dribble.NewClient()
	return datasource.Adapters()
}

// TargetingError for dribble requests
type TargetingError string

// Error implements error
func (e TargetingError) Error() string {
	return string(e)
}

// APIRequestErrorMsg for dribble requests
type APIRequestErrorMsg error

// NewTargetingError creates a new TargetingError
func NewTargetingError(targetName string) error {
	if targetName == "" {
		return TargetingError("Target name cannot be empty")
	}
	msg := fmt.Sprintf("Target %s not found", targetName)
	return TargetingError(msg)
}

// NewAPIError creates a new DribbleRequestError
func NewAPIError(err error) tea.Cmd {
	return func() tea.Msg {
		return APIRequestErrorMsg(err)
	}
}

type (
	DribbleApiMsg struct {
		Request tea.Msg
	}
	DribbleResponseMsg struct {
		Response tea.Msg
	}

	TargetOpenedMsg struct {
		Target *target.Target
	}

	DSNOpenTargetMsg struct {
		Name      string
		DSN       datasource.Namer
		OnSuccess []datasource.Request
	}

	AdapterOpenTargetMsg struct {
		Name    string
		Adapter datasource.Adapter
	}
)

func TargetOpened(target *target.Target) tea.Msg {
	return DribbleResponseMsg{
		Response: TargetOpenedMsg{
			Target: target,
		},
	}
}

func DSNOpen(name string, dsn datasource.Namer, onSuccess ...datasource.Request) tea.Msg {
	return DribbleApiMsg{
		Request: DSNOpenTargetMsg{
			Name:      name,
			DSN:       dsn,
			OnSuccess: onSuccess,
		},
	}
}

func AdapterOpen(name string, adapter datasource.Adapter) tea.Msg {
	return DribbleApiMsg{
		Request: AdapterOpenTargetMsg{
			Name:    name,
			Adapter: adapter,
		},
	}
}

type (
	// DataSourceRequestMsg should be used by widgets to command the main Model to request from the dribble api
	DataSourceRequestMsg struct {
		TargetID   *int
		TargetName *string
		Request    datasource.Request
		Context    context.Context
	}
)

// Request is used to create a tea.Cmd that will run a request against the dribble api
func Request(opts ...requestOption) tea.Cmd {
	return func() tea.Msg {
		rm := &DataSourceRequestMsg{
			Context: context.Background(), // Default context
		}
		for _, opt := range opts {
			opt(rm)
		}
		return *rm
	}
}

type requestOption func(*DataSourceRequestMsg)

func WithTargetID(id int) requestOption {
	return func(rm *DataSourceRequestMsg) {
		rm.TargetID = &id
	}
}

func WithTargetName(name string) requestOption {
	return func(rm *DataSourceRequestMsg) {
		rm.TargetName = &name
	}
}

func WithRequest(req datasource.Request) requestOption {
	return func(rm *DataSourceRequestMsg) {
		rm.Request = req
	}
}

func WithContext(ctx context.Context) requestOption {
	return func(rm *DataSourceRequestMsg) {
		rm.Context = ctx
	}
}

// Requester is used to fluently enable requests to dribble api
type Requester struct {
	Target *target.Target
}

// Request is used to create a tea.Cmd that will run a request against the dribble api
func (d Requester) Request(ctx context.Context, request datasource.Request) tea.Cmd {
	responseChan, err := d.Target.Request(ctx, request)
	if err != nil {
		return func() tea.Msg {
			return NewAPIError(err)
		}
	}

	return func() tea.Msg {
		return ResponseMsg{
			Channel: responseChan,
		}
	}
}

// ResponseMsg is used to either get the channel or collect it all
type ResponseMsg struct {
	Channel chan *request.Response
}

// All collects all Responses into a list
func (d ResponseMsg) All() []*request.Response {
	responses := []*request.Response{}
	for response := range d.Channel {
		responses = append(responses, response)
	}
	return responses
}
