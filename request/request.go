package request

import (
	"github.com/ctrl-alt-boop/dribble/datasource"
)

var _ datasource.Request = (*Intent)(nil)

var (
	_ datasource.Request = (*BatchRequest)(nil)
	_ datasource.Request = (*ChainRequest)(nil)
)

type (
	Intent struct { // TODO: finalize this
		Type datasource.RequestType

		Operation any

		Args []any
	}

	BatchRequest []datasource.Request

	ChainRequest []datasource.Request
)

func Batch(requests ...datasource.Request) BatchRequest {
	return requests
}

func Chain(requests ...datasource.Request) ChainRequest {
	return requests
}

// IsPrefab implements database.Request.
func (i Intent) IsPrefab() bool {
	return false
}

// ResponseOnError implements database.Request.
func (i Intent) ResponseOnError() datasource.Response {
	switch i.Type {
	case datasource.Create:
		return Response{
			Status: ErrorCreate,
		}
	case datasource.Read:
		return Response{
			Status: ErrorRead,
		}
	case datasource.Update:
		return Response{
			Status: ErrorUpdate,
		}
	case datasource.Delete:
		return Response{
			Status: ErrorDelete,
		}
	default:
		panic("unknown request type")
	}
}

// ResponseOnSuccess implements database.Request.
func (i Intent) ResponseOnSuccess() datasource.Response {
	switch i.Type {
	case datasource.Create:
		return Response{
			Status: SuccessCreate,
		}
	case datasource.Read:
		return Response{
			Status: SuccessRead,
		}
	case datasource.Update:
		return Response{
			Status: SuccessUpdate,
		}
	case datasource.Delete:
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
func (c ChainRequest) ResponseOnError() datasource.Response {
	return Response{
		Status: ErrorChainExecute,
	}
}

// ResponseOnSuccess implements database.Request.
func (c ChainRequest) ResponseOnSuccess() datasource.Response {
	return Response{
		Status: SuccessChainExecute,
	}
}

// IsPrefab implements database.Request.
func (b BatchRequest) IsPrefab() bool {
	return false
}

// ResponseOnError implements database.Request.
func (b BatchRequest) ResponseOnError() datasource.Response {
	return Response{
		Status: ErrorBatchExecute,
	}
}

// ResponseOnSuccess implements database.Request.
func (b BatchRequest) ResponseOnSuccess() datasource.Response {
	return Response{
		Status: SuccessBatchExecute,
	}
}
