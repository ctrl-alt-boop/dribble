package sql

import "github.com/ctrl-alt-boop/dribble/database"

type SelectQuery struct {
	AsDistinct bool
	IsCount    bool

	Fields []string
	Table  string

	Joins []joinClause

	WhereClause string
	args        []any

	GroupByClause []string
	HavingClause  string
	OrderByClause []orderByClause

	LimitClause  *int
	OffsetClause *int
}

type SelectBuilder struct {
	asDistinct bool
	isCount    bool

	fields []string
	table  string

	joins []joinClause

	whereClause string
	params      []any

	groupByClause []string
	havingClause  string
	orderByClause []orderByClause

	limitClause  *int
	offsetClause *int
}

func Select(fields ...string) *From {
	return &From{
		distinct: false,
		fields:   fields,
	}
}
func DistinctSelect(fields ...string) *From {
	return &From{
		distinct: true,
		fields:   fields,
	}
}

func SelectAll() *From {
	return &From{
		distinct: false,
		fields:   []string{"*"},
	}
}

func DistinctSelectAll() *From {
	return &From{
		distinct: true,
		fields:   []string{"*"},
	}
}

type From struct {
	distinct bool
	fields   []string
}

func (b From) From(table string, joins ...joinClause) *SelectBuilder {
	return &SelectBuilder{
		asDistinct:    b.distinct,
		fields:        b.fields,
		table:         table,
		joins:         joins,
		params:        []any{},
		orderByClause: []orderByClause{},
	}
}

func (s *SelectBuilder) Copy() *SelectBuilder {
	return &SelectBuilder{
		asDistinct:    s.asDistinct,
		fields:        s.fields,
		table:         s.table,
		joins:         s.joins,
		whereClause:   s.whereClause,
		params:        s.params,
		groupByClause: s.groupByClause,
		havingClause:  s.havingClause,
		orderByClause: s.orderByClause,
		limitClause:   s.limitClause,
		offsetClause:  s.offsetClause,
	}
}

type JoinType string

const (
	JoinTypeInner JoinType = "INNER"
	JoinTypeLeft  JoinType = "LEFT"
	JoinTypeRight JoinType = "RIGHT"
	JoinTypeFull  JoinType = "FULL"
)

type joinClause struct {
	Type  JoinType
	Table string
	On    string
}

func InnerJoin(table, on string) joinClause {
	return joinClause{
		Type:  JoinTypeInner,
		Table: table,
		On:    on,
	}
}

func LeftJoin(table, on string) joinClause {
	return joinClause{
		Type:  JoinTypeLeft,
		Table: table,
		On:    on,
	}
}

func RightJoin(table, on string) joinClause {
	return joinClause{
		Type:  JoinTypeRight,
		Table: table,
		On:    on,
	}
}

func FullJoin(table, on string) joinClause {
	return joinClause{
		Type:  JoinTypeFull,
		Table: table,
		On:    on,
	}
}

func Count(field, table string) *SelectBuilder {
	return &SelectBuilder{
		isCount: true,
		fields:  []string{"COUNT(" + field + ")"},
		table:   table,
	}
}

func (s *SelectBuilder) ShouldReturn() int {
	if s.isCount {
		return 0 // Scalar
	}
	if s.fields[0] == "*" || len(s.fields) > 1 {
		return 2 // Table
	}

	return len(s.fields)
}

func (s *SelectBuilder) Where(expr ...database.Expr) *SelectBuilder {
	s.whereClause, s.params = database.Exprs(expr).ToSql()
	return s
}

type orderByClause struct {
	Field string
	Desc  bool
}

func (s *SelectBuilder) OrderBy(field string, desc bool) *SelectBuilder { // FIXME: Something funky, is broken
	s.orderByClause = append(s.orderByClause, orderByClause{
		Field: field,
		Desc:  desc,
	})
	return s
}

func (s *SelectBuilder) GroupBy(fields ...string) *SelectBuilder {
	s.groupByClause = fields
	return s
}

func (s *SelectBuilder) Having(having string) *SelectBuilder {
	s.havingClause = having
	return s
}

func (s *SelectBuilder) Limit(limit int) *SelectBuilder {
	s.limitClause = &limit
	return s
}

func (s *SelectBuilder) Offset(offset int) *SelectBuilder {
	s.offsetClause = &offset
	return s
}

func (s *SelectBuilder) ToIntent() *database.Intent { // FIXME
	return nil
	// return &database.QueryIntent{
	// 	Type:       database.ReadQuery,
	// 	QueryStyle: database.SQL,
	// 	SQLQuery: &SQLSelectQuery{
	// 		AsDistinct:    s.asDistinct,
	// 		IsCount:       s.isCount,
	// 		Fields:        s.fields,
	// 		Table:         s.table,
	// 		Joins:         s.joins,
	// 		WhereClause:   s.whereClause,
	// 		args:          s.params,
	// 		GroupByClause: s.groupByClause,
	// 		HavingClause:  s.havingClause,
	// 		OrderByClause: s.orderByClause,
	// 		LimitClause:   s.limitClause,
	// 		OffsetClause:  s.offsetClause,
	// 	},
	// 	Args: s.params,
	// }
}

func (s *SelectBuilder) ToQueryOn(targetName string) *database.Intent { // FIXME
	return nil
	// return &database.QueryIntent{
	// 	Type:       database.ReadQuery,
	// 	QueryStyle: database.SQL,
	// 	TargetName: targetName,
	// 	SQLQuery: &SQLSelectQuery{
	// 		AsDistinct:    s.asDistinct,
	// 		IsCount:       s.isCount,
	// 		Fields:        s.fields,
	// 		Table:         s.table,
	// 		Joins:         s.joins,
	// 		WhereClause:   s.whereClause,
	// 		args:          s.params,
	// 		GroupByClause: s.groupByClause,
	// 		HavingClause:  s.havingClause,
	// 		OrderByClause: s.orderByClause,
	// 		LimitClause:   s.limitClause,
	// 		OffsetClause:  s.offsetClause,
	// 	},
	// 	Args: s.params,
	// }
}

func (s *SelectBuilder) Parameters() []any {
	return s.params
}
