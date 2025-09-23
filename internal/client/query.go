package client

// import "time"

// type SQLMethod string

// func (s SQLMethod) String() string {
// 	return string(s)
// }

// const (
// 	MethodSelect SQLMethod = "SELECT"
// 	MethodInsert SQLMethod = "INSERT"
// 	MethodUpdate SQLMethod = "UPDATE"
// 	MethodDelete SQLMethod = "DELETE"
// )

// const DefaultSelectLimit int = 10 // Just a safeguard

// var SQLMethods = []SQLMethod{MethodSelect, MethodInsert, MethodUpdate, MethodDelete}

// type Dialect interface {
// 	SelectTemplate() string
// 	InsertTemplate() string
// 	UpdateTemplate() string
// 	DeleteTemplate() string
// 	SupportsJsonResult() bool
// 	IsFile() bool
// 	RenderPlaceholder(index int) string
// 	RenderLimit(limit int) string
// 	RenderOffset(offset int) string
// 	RenderOrderBy(column string, desc bool) string
// 	RenderTypeCast() string
// 	RenderCurrentTimestamp() string
// 	RenderDateFormatting(date time.Time, format string) string
// 	Quote(value string) string
// 	QuoteRune() rune
// }

// type Statement struct {
// 	Method  SQLMethod
// 	Table   string
// 	Columns []string
// 	Values  []any
// 	Set     []struct {
// 		Column string
// 		Value  any
// 	}
// 	Where []struct {
// 		Column   string
// 		Operator string
// 		Value    any
// 	}
// 	OrderBy struct {
// 		Column string
// 		Desc   bool
// 	}
// 	Limit  int
// 	Offset int
// }
