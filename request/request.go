package request

import (
	"github.com/ctrl-alt-boop/dribble/database"
)

var _ database.Request = (*Intent)(nil)

var _ database.Request = (*BatchRequest)(nil)
var _ database.Request = (*ChainRequest)(nil)

type (
	Intent struct { // TODO: finalize this
		Type database.RequestType

		Operation any

		Args []any
	}

	BatchRequest []database.Request

	ChainRequest []database.Request
)

func Batch(requests ...database.Request) BatchRequest {
	return requests
}

func Chain(requests ...database.Request) ChainRequest {
	return requests
}

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

// IsPrefab implements database.Request.
func (c ChainRequest) IsPrefab() bool {
	return false
}

// ResponseOnError implements database.Request.
func (c ChainRequest) ResponseOnError() database.Response {
	return Response{
		Status: ErrorChainExecute,
	}
}

// ResponseOnSuccess implements database.Request.
func (c ChainRequest) ResponseOnSuccess() database.Response {
	return Response{
		Status: SuccessChainExecute,
	}
}

// IsPrefab implements database.Request.
func (b BatchRequest) IsPrefab() bool {
	return false
}

// ResponseOnError implements database.Request.
func (b BatchRequest) ResponseOnError() database.Response {
	return Response{
		Status: ErrorBatchExecute,
	}
}

// ResponseOnSuccess implements database.Request.
func (b BatchRequest) ResponseOnSuccess() database.Response {
	return Response{
		Status: SuccessBatchExecute,
	}
}
