package database

type SQLStyleSelectBuilder struct {
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

func Select(fields ...string) *SQLStyleFrom {
	return &SQLStyleFrom{
		distinct: false,
		fields:   fields,
	}
}
func DistinctSelect(fields ...string) *SQLStyleFrom {
	return &SQLStyleFrom{
		distinct: true,
		fields:   fields,
	}
}

func SelectAll() *SQLStyleFrom {
	return &SQLStyleFrom{
		distinct: false,
		fields:   []string{"*"},
	}
}

func DistinctSelectAll() *SQLStyleFrom {
	return &SQLStyleFrom{
		distinct: true,
		fields:   []string{"*"},
	}
}

type SQLStyleFrom struct {
	distinct bool
	fields   []string
}

func (b SQLStyleFrom) From(table string, joins ...joinClause) *SQLStyleSelectBuilder {
	return &SQLStyleSelectBuilder{
		asDistinct:    b.distinct,
		fields:        b.fields,
		table:         table,
		joins:         joins,
		params:        []any{},
		orderByClause: []orderByClause{},
	}
}

func (s *SQLStyleSelectBuilder) Copy() *SQLStyleSelectBuilder {
	return &SQLStyleSelectBuilder{
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

func Count(field, table string) *SQLStyleSelectBuilder {
	return &SQLStyleSelectBuilder{
		isCount: true,
		fields:  []string{"COUNT(" + field + ")"},
		table:   table,
	}
}

func (s *SQLStyleSelectBuilder) ShouldReturn() int {
	if s.isCount {
		return 0 // Scalar
	}
	if s.fields[0] == "*" || len(s.fields) > 1 {
		return 2 // Table
	}

	return len(s.fields)
}

func (s *SQLStyleSelectBuilder) Where(expr ...Expr) *SQLStyleSelectBuilder {
	s.whereClause, s.params = Exprs(expr).ToSql()
	return s
}

type orderByClause struct {
	Field string
	Desc  bool
}

func (s *SQLStyleSelectBuilder) OrderBy(field string, desc bool) *SQLStyleSelectBuilder { // FIXME: Something funky, is broken
	s.orderByClause = append(s.orderByClause, orderByClause{
		Field: field,
		Desc:  desc,
	})
	return s
}

func (s *SQLStyleSelectBuilder) GroupBy(fields ...string) *SQLStyleSelectBuilder {
	s.groupByClause = fields
	return s
}

func (s *SQLStyleSelectBuilder) Having(having string) *SQLStyleSelectBuilder {
	s.havingClause = having
	return s
}

func (s *SQLStyleSelectBuilder) Limit(limit int) *SQLStyleSelectBuilder {
	s.limitClause = &limit
	return s
}

func (s *SQLStyleSelectBuilder) Offset(offset int) *SQLStyleSelectBuilder {
	s.offsetClause = &offset
	return s
}

func (s *SQLStyleSelectBuilder) ToQuery() *QueryIntent {
	return &QueryIntent{
		Type:       ReadQuery,
		QueryStyle: SQL,
		SQLQuery: &SQLSelectQuery{
			AsDistinct:    s.asDistinct,
			IsCount:       s.isCount,
			Fields:        s.fields,
			Table:         s.table,
			Joins:         s.joins,
			WhereClause:   s.whereClause,
			args:          s.params,
			GroupByClause: s.groupByClause,
			HavingClause:  s.havingClause,
			OrderByClause: s.orderByClause,
			LimitClause:   s.limitClause,
			OffsetClause:  s.offsetClause,
		},
		Args: s.params,
	}
}

func (s *SQLStyleSelectBuilder) ToQueryOn(targetName string) *QueryIntent {
	return &QueryIntent{
		Type:       ReadQuery,
		QueryStyle: SQL,
		TargetName: targetName,
		SQLQuery: &SQLSelectQuery{
			AsDistinct:    s.asDistinct,
			IsCount:       s.isCount,
			Fields:        s.fields,
			Table:         s.table,
			Joins:         s.joins,
			WhereClause:   s.whereClause,
			args:          s.params,
			GroupByClause: s.groupByClause,
			HavingClause:  s.havingClause,
			OrderByClause: s.orderByClause,
			LimitClause:   s.limitClause,
			OffsetClause:  s.offsetClause,
		},
		Args: s.params,
	}
}

func (s *SQLStyleSelectBuilder) Parameters() []any {
	return s.params
}
