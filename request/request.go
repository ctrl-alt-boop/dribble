package request

import (
	"fmt"
	"reflect"

	"github.com/ctrl-alt-boop/dribble/database"
)

var _ database.Request = Intent{}
var _ database.Response = Response{}

type (
	Intent struct { // TODO: finalize this
		Type database.RequestType

		// Was this supposed to be Select, Insert, Find etc.?
		OperationKind reflect.Kind

		Operation any

		Args []any
	}

	IntentBatch []*Intent

	Response struct {
		Status    Status
		RequestID int64
		Body      any
		Error     error
	}
)

// IsPrefab implements database.Request.
func (i Intent) IsPrefab() bool {
	return false
}

// ResponseOnError implements database.Request.
func (i Intent) ResponseOnError() database.Response {
	switch i.Type {
	case database.Create:
		return Response{
			Status: ErrorCreate,
		}
	case database.Read:
		return Response{
			Status: ErrorRead,
		}
	case database.Update:
		return Response{
			Status: ErrorUpdate,
		}
	case database.Delete:
		return Response{
			Status: ErrorDelete,
		}
	default:
		panic("unknown request type")
	}
}

// ResponseOnSuccess implements database.Request.
func (i Intent) ResponseOnSuccess() database.Response {
	switch i.Type {
	case database.Create:
		return Response{
			Status: SuccessCreate,
		}
	case database.Read:
		return Response{
			Status: SuccessRead,
		}
	case database.Update:
		return Response{
			Status: SuccessUpdate,
		}
	case database.Delete:
		return Response{
			Status: SuccessDelete,
		}
	default:
		panic("unknown request type")
	}
}

func (r Response) Code() int {
	return int(r.Status)
}

func (r Response) Message() string {
	return fmt.Sprintf("%d: %s", r.Code(), r.Status.String())
}

type RequestBatch []database.Request

func Batch(requests ...database.Request) RequestBatch {
	return requests
}

type RequestChain []database.Request

func Chain(requests ...database.Request) RequestChain {
	return requests
}
