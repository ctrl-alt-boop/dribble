package dribble

// I will temporarily assume that args will be in the same order as query parameter placeholders
type request interface {
	StatusOnSuccess() Success
	StatusOnError() Error
}

var (
	_ request = ReadDatabaseSchema{}
	_ request = ReadTableSchema{}
	_ request = ReadColumnSchema{}
	_ request = ReadDatabaseProperties{}
	_ request = ReadTableProperties{}
	_ request = ReadColumnProperties{}
	_ request = ReadDatabaseNames{}
	_ request = ReadTableNames{}
	_ request = ReadColumnNames{}
	_ request = ReadCount{}
	_ request = ReadAllCounts{}
	_ request = Execute{}
	_ request = BatchExecute{}
)

// BatchExecute and BatchRequest are not the same thing
// BatchRequest will send through multiple channels while BatchExecute will send all results in one channel
// Each channel results in one Response for each request in the BatchRequest
type batchRequest []request

func BatchRequest(requests ...request) batchRequest {
	return requests
}

type (
	// args = [database_name]
	ReadDatabaseSchema struct {
		DatabaseName string
	}

	// args = [database name, table name]
	ReadTableSchema struct {
		DatabaseName string
		TableName    string
	}

	// args = [database name, table name, column name]
	ReadColumnSchema struct {
		DatabaseName string
		TableName    string
		ColumnName   string
	}

	// args = [database name]
	ReadDatabaseProperties struct {
		DatabaseName string
	}

	// args = [database name, table name]
	ReadTableProperties struct {
		DatabaseName string
		TableName    string
	}

	// args = [database name, table name, column name]
	ReadColumnProperties struct {
		DatabaseName string
		TableName    string
		ColumnName   string
	}

	// args = [target name]
	ReadDatabaseNames struct {
		TargetName string
	}

	// args = [database name]
	ReadTableNames struct {
		DatabaseName string
	}

	// args = [database name, table name]
	ReadColumnNames struct {
		DatabaseName string
		TableName    string
	}

	// args = [database name, table name]
	ReadCount struct {
		DatabaseName string
		TableName    string
	}

	// args = [[database name, table name], [database name, table name], ...]
	ReadAllCounts struct {
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

// StatusOnError implements request.
func (r ReadDatabaseSchema) StatusOnError() Error {
	return ErrorReadDatabaseSchema
}

// StatusOnSuccess implements request.
func (r ReadDatabaseSchema) StatusOnSuccess() Success {
	return SuccessReadDatabaseSchema
}

// StatusOnError implements request.
func (r ReadTableSchema) StatusOnError() Error {
	return ErrorReadTableSchema
}

// StatusOnSuccess implements request.
func (r ReadTableSchema) StatusOnSuccess() Success {
	return SuccessReadTableSchema
}

// StatusOnError implements request.
func (r ReadColumnSchema) StatusOnError() Error {
	return ErrorReadColumnSchema
}

// StatusOnSuccess implements request.
func (r ReadColumnSchema) StatusOnSuccess() Success {
	return SuccessReadColumnSchema
}

// StatusOnError implements request.
func (r ReadDatabaseProperties) StatusOnError() Error {
	return ErrorReadDatabaseProperties
}

// StatusOnSuccess implements request.
func (r ReadDatabaseProperties) StatusOnSuccess() Success {
	return SuccessReadDatabaseProperties
}

// StatusOnError implements request.
func (r ReadTableProperties) StatusOnError() Error {
	return ErrorReadTableProperties
}

// StatusOnSuccess implements request.
func (r ReadTableProperties) StatusOnSuccess() Success {
	return SuccessReadTableProperties
}

// StatusOnError implements request.
func (r ReadColumnProperties) StatusOnError() Error {
	return ErrorReadColumnProperties
}

// StatusOnSuccess implements request.
func (r ReadColumnProperties) StatusOnSuccess() Success {
	return SuccessReadColumnProperties
}

// StatusOnError implements request.
func (r ReadDatabaseNames) StatusOnError() Error {
	return ErrorReadDatabaseList
}

// StatusOnSuccess implements request.
func (r ReadDatabaseNames) StatusOnSuccess() Success {
	return SuccessReadDatabaseList
}

// StatusOnError implements request.
func (r ReadTableNames) StatusOnError() Error {
	return ErrorReadDBTableList
}

// StatusOnSuccess implements request.
func (r ReadTableNames) StatusOnSuccess() Success {
	return SuccessReadDBTableList
}

// StatusOnError implements request.
func (r ReadColumnNames) StatusOnError() Error {
	return ErrorReadDBColumnList
}

// StatusOnSuccess implements request.
func (r ReadColumnNames) StatusOnSuccess() Success {
	return SuccessReadDBColumnList
}

// StatusOnError implements request.
func (r ReadCount) StatusOnError() Error {
	return ErrorReadCount
}

// StatusOnSuccess implements request.
func (r ReadCount) StatusOnSuccess() Success {
	return SuccessReadCount
}

// StatusOnError implements request.
func (e Execute) StatusOnError() Error {
	return ErrorExecute
}

// StatusOnSuccess implements request.
func (e Execute) StatusOnSuccess() Success {
	return SuccessExecute
}

// StatusOnError implements request.
func (b BatchExecute) StatusOnError() Error {
	return ErrorBatchExecute
}

// StatusOnSuccess implements request.
func (b BatchExecute) StatusOnSuccess() Success {
	return SuccessBatchExecute
}
