package database

type NoSQLStyleSelectBuilder struct {
	collection string

	conditions Exprs
	args       []any

	limitClause  *int
	offsetClause *int
}

func Find() *NoSQLStyleSelectBuilder {
	return &NoSQLStyleSelectBuilder{}
}

func (n *NoSQLStyleSelectBuilder) Cond(expr ...Expr) *NoSQLStyleSelectBuilder {
	n.conditions = expr
	return n
}

func (n *NoSQLStyleSelectBuilder) Limit(limit int) *NoSQLStyleSelectBuilder {
	n.limitClause = &limit
	return n
}

func (n *NoSQLStyleSelectBuilder) Offset(offset int) *NoSQLStyleSelectBuilder {
	n.offsetClause = &offset
	return n
}

func (n *NoSQLStyleSelectBuilder) ToQuery() *QueryIntent {
	return &QueryIntent{
		Type:       ReadQuery,
		QueryStyle: NoSQL,
		NoSQLQuery: &NoSQLSelectQuery{
			Collection:   n.collection,
			Conditions:   n.conditions,
			args:         n.args,
			LimitClause:  n.limitClause,
			OffsetClause: n.offsetClause,
		},
		Args: n.args,
	}
}

func (n *NoSQLStyleSelectBuilder) Parameters() []any {
	return n.args
}
