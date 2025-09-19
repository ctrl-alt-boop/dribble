package request

//go:generate stringer -type=Status

type (
	Status      int
	BatchStatus []Status
)

func (s Status) IsSuccess() bool {
	return s > 0
}

func (s Status) IsError() bool {
	return s < 0
}

func (s BatchStatus) AllSuccess() bool {
	for _, status := range s {
		if !status.IsSuccess() {
			return false
		}
	}
	return true
}

func (s BatchStatus) GetErrorIndices() []int {
	var indices []int
	for i, status := range s {
		if !status.IsSuccess() {
			indices = append(indices, i)
		}
	}
	return indices
}

const StatusUnknown Status = 0

const (
	_ Status = iota
	SuccessConnect
	SuccessReconnect
	SuccessDisconnect
	SuccessTargetOpen
	SuccessTargetUpdate
	SuccessTargetClose

	SuccessReadDatabaseSchema
	SuccessReadTableSchema
	SuccessReadColumnSchema

	SuccessReadDatabaseProperties
	SuccessReadTableProperties
	SuccessReadColumnProperties

	SuccessReadDatabaseList
	SuccessReadDBTableList
	SuccessReadDBColumnList

	SuccessReadCount

	SuccessCreate
	SuccessRead
	SuccessUpdate
	SuccessDelete

	SuccessExecute
	SuccessBatchExecute
)

const (
	_ Status = -iota
	ErrorConnect
	ErrorReconnect
	ErrorDisconnect
	ErrorTargetOpen
	ErrorTargetClose
	ErrorTargetUpdate

	ErrorReadDatabaseSchema
	ErrorReadTableSchema
	ErrorReadColumnSchema

	ErrorReadDatabaseProperties
	ErrorReadTableProperties
	ErrorReadColumnProperties

	ErrorReadDatabaseList
	ErrorReadDBTableList
	ErrorReadDBColumnList

	ErrorReadCount

	ErrorCreate
	ErrorRead
	ErrorUpdate
	ErrorDelete

	ErrorExecute
	ErrorBatchExecute
)
