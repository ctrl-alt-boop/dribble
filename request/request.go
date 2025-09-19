package request

import (
	"fmt"
	"reflect"

	"github.com/ctrl-alt-boop/dribble/database"
)

var _ database.Request = Intent{}

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
		Status
		RequestID int64
		Body      any
		Error     error
	}
)

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
	return fmt.Sprintf("%d: %s", r.Code(), r.String())
}
