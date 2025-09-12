package nosql

import "github.com/ctrl-alt-boop/dribble/database"

type FindQuery struct {
	Collection string

	ConditionsClause string
	Conditions       database.Exprs
	args             []any

	LimitClause  *int
	OffsetClause *int
}

type FindBuilder struct {
	collection string

	conditions database.Exprs
	args       []any

	limitClause  *int
	offsetClause *int
}

func Find() *FindBuilder {
	return &FindBuilder{}
}

func (n *FindBuilder) Cond(expr ...database.Expr) *FindBuilder {
	n.conditions = expr
	return n
}

func (n *FindBuilder) Limit(limit int) *FindBuilder {
	n.limitClause = &limit
	return n
}

func (n *FindBuilder) Offset(offset int) *FindBuilder {
	n.offsetClause = &offset
	return n
}

func (n *FindBuilder) ToIntent() *database.Intent {
	return nil
	// return &database.QueryIntent{
	// 	Type:       ReadQuery,
	// 	QueryStyle: NoSQL,
	// 	NoSQLQuery: &NoSQLSelectQuery{
	// 		Collection:   n.collection,
	// 		Conditions:   n.conditions,
	// 		args:         n.args,
	// 		LimitClause:  n.limitClause,
	// 		OffsetClause: n.offsetClause,
	// 	},
	// 	Args: n.args,
	// }
}

func (n *FindBuilder) Parameters() []any {
	return n.args
}
