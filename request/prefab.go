package request

import "github.com/ctrl-alt-boop/dribble/database"

var (
	_ database.Request = ReadDatabaseSchema{}
	_ database.Request = ReadTableSchema{}
	_ database.Request = ReadColumnSchema{}
	_ database.Request = ReadDatabaseProperties{}
	_ database.Request = ReadTableProperties{}
	_ database.Request = ReadColumnProperties{}
	_ database.Request = ReadDatabaseNames{}
	_ database.Request = ReadTableNames{}
	_ database.Request = ReadColumnNames{}
	_ database.Request = ReadCount{}
	_ database.Request = ReadAllCounts{}
	_ database.Request = Execute{}
	_ database.Request = BatchExecute{}
)

type PrefabRequest struct{}

func (p PrefabRequest) IsPrefab() bool {
	return true
}

// BatchExecute and BatchRequest are not the same thing
// BatchRequest will send through multiple channels while BatchExecute will send all results in one channel
// Each channel results in one Response for each request in the BatchRequest
type batchRequest []database.Request

func BatchRequest(requests ...database.Request) batchRequest {
	return requests
}

type RequestChain []database.Request

func ChainRequest(requests ...database.Request) RequestChain {
	return requests
}

type (
	// args = [database_name]
	ReadDatabaseSchema struct {
		PrefabRequest
		DatabaseName string
	}

	// args = [database name, table name]
	ReadTableSchema struct {
		PrefabRequest
		DatabaseName string
		TableName    string
	}

	// args = [database name, table name, column name]
	ReadColumnSchema struct {
		PrefabRequest
		DatabaseName string
		TableName    string
		ColumnName   string
	}

	// args = [database name]
	ReadDatabaseProperties struct {
		PrefabRequest
		DatabaseName string
	}

	// args = [database name, table name]
	ReadTableProperties struct {
		PrefabRequest
		DatabaseName string
		TableName    string
	}

	// args = [database name, table name, column name]
	ReadColumnProperties struct {
		PrefabRequest
		DatabaseName string
		TableName    string
		ColumnName   string
	}

	// args = [target name]
	ReadDatabaseNames struct {
		PrefabRequest
		TargetName string
	}

	// args = [database name]
	ReadTableNames struct {
		PrefabRequest
		DatabaseName string
	}

	// args = [database name, table name]
	ReadColumnNames struct {
		PrefabRequest
		DatabaseName string
		TableName    string
	}

	// args = [database name, table name]
	ReadCount struct {
		PrefabRequest
		DatabaseName string
		TableName    string
	}

	// args = [[database name, table name], [database name, table name], ...]
	ReadAllCounts struct {
		PrefabRequest
		ReadCount
		DatabaseName string
		TableNames   []string
	}

	// args = []
	Execute struct {
	}

	// BatchExecute and BatchRequest are not the same thing
	// BatchExecute will send all results in one channel while a BatchRequest will send through multiple channels
	// This results in a BatchResponse
	BatchExecute struct {
	}
)

// ResponseOnError implements database.Request.
func (r ReadDatabaseSchema) ResponseOnError() database.Response {
	return Response{
		Status: ErrorReadDatabaseSchema,
	}
}

// ResponseOnSuccess implements database.Request.
func (r ReadDatabaseSchema) ResponseOnSuccess() database.Response {
	return Response{
		Status: SuccessReadDatabaseSchema,
	}
}

// ResponseOnError implements database.Request.
func (r ReadTableSchema) ResponseOnError() database.Response {
	return Response{
		Status: ErrorReadTableSchema,
	}
}

// ResponseOnSuccess implements database.Request.
func (r ReadTableSchema) ResponseOnSuccess() database.Response {
	return Response{
		Status: SuccessReadTableSchema,
	}
}

// ResponseOnError implements database.Request.
func (r ReadColumnSchema) ResponseOnError() database.Response {
	return Response{
		Status: ErrorReadColumnSchema,
	}
}

// ResponseOnSuccess implements database.Request.
func (r ReadColumnSchema) ResponseOnSuccess() database.Response {
	return Response{
		Status: SuccessReadColumnSchema,
	}
}

// ResponseOnError implements database.Request.
func (r ReadDatabaseProperties) ResponseOnError() database.Response {
	return Response{
		Status: ErrorReadDatabaseProperties,
	}
}

// ResponseOnSuccess implements database.Request.
func (r ReadDatabaseProperties) ResponseOnSuccess() database.Response {
	return Response{
		Status: SuccessReadDatabaseProperties,
	}
}

// ResponseOnError implements database.Request.
func (r ReadTableProperties) ResponseOnError() database.Response {
	return Response{
		Status: ErrorReadTableProperties,
	}
}

// ResponseOnSuccess implements database.Request.
func (r ReadTableProperties) ResponseOnSuccess() database.Response {
	return Response{
		Status: SuccessReadTableProperties,
	}
}

// ResponseOnError implements database.Request.
func (r ReadColumnProperties) ResponseOnError() database.Response {
	return Response{
		Status: ErrorReadColumnProperties,
	}
}

// ResponseOnSuccess implements database.Request.
func (r ReadColumnProperties) ResponseOnSuccess() database.Response {
	return Response{
		Status: SuccessReadColumnProperties,
	}
}

// ResponseOnError implements database.Request.
func (r ReadDatabaseNames) ResponseOnError() database.Response {
	return Response{
		Status: ErrorReadDatabaseList,
	}
}

// ResponseOnSuccess implements database.Request.
func (r ReadDatabaseNames) ResponseOnSuccess() database.Response {
	return Response{
		Status: SuccessReadDatabaseList,
	}
}

// ResponseOnError implements database.Request.
func (r ReadTableNames) ResponseOnError() database.Response {
	return Response{
		Status: ErrorReadDBTableList,
	}
}

// ResponseOnSuccess implements database.Request.
func (r ReadTableNames) ResponseOnSuccess() database.Response {
	return Response{
		Status: SuccessReadDBTableList,
	}
}

// ResponseOnError implements database.Request.
func (r ReadColumnNames) ResponseOnError() database.Response {
	return Response{
		Status: ErrorReadDBColumnList,
	}
}

// ResponseOnSuccess implements database.Request.
func (r ReadColumnNames) ResponseOnSuccess() database.Response {
	return Response{
		Status: SuccessReadDBColumnList,
	}
}

// ResponseOnError implements database.Request.
func (r ReadCount) ResponseOnError() database.Response {
	return Response{
		Status: ErrorReadCount,
	}
}

// ResponseOnSuccess implements database.Request.
func (r ReadCount) ResponseOnSuccess() database.Response {
	return Response{
		Status: SuccessReadCount,
	}
}

// ResponseOnError implements database.Request.
func (e Execute) ResponseOnError() database.Response {
	return Response{
		Status: ErrorExecute,
	}
}

// ResponseOnSuccess implements database.Request.
func (e Execute) ResponseOnSuccess() database.Response {
	return Response{
		Status: SuccessExecute,
	}
}

// ResponseOnError implements database.Request.
func (b BatchExecute) ResponseOnError() database.Response {
	return Response{
		Status: ErrorBatchExecute,
	}
}

// ResponseOnSuccess implements database.Request.
func (b BatchExecute) ResponseOnSuccess() database.Response {
	return Response{
		Status: SuccessBatchExecute,
	}
}

// IsPrefab implements database.Request.
func (b BatchExecute) IsPrefab() bool {
	return false
}

// IsPrefab implements database.Request.
func (e Execute) IsPrefab() bool {
	return false
}
