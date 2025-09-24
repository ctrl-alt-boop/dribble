package nosql

import (
	"github.com/ctrl-alt-boop/dribble/database"
	"github.com/ctrl-alt-boop/dribble/request"
)

type FindQuery struct {
	Collection string

	ConditionsClause string
	Conditions       Exprs
	args             []any

	LimitClause  *int
	OffsetClause *int
}

type FindBuilder struct {
	collection string

	conditions Exprs
	args       []any

	limitClause  *int
	offsetClause *int
}

func Find() *FindBuilder {
	return &FindBuilder{}
}

func (n *FindBuilder) Cond(expr ...Expr) *FindBuilder {
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

func (n *FindBuilder) ToIntent() database.Request { // TODO: Ofcourse this needs implementation for the other operation types
	intent := &FindQuery{
		Collection:   n.collection,
		Conditions:   n.conditions,
		args:         n.args,
		LimitClause:  n.limitClause,
		OffsetClause: n.offsetClause,
	}
	return &request.Intent{
		Type:      database.Read,
		Operation: intent,
		Args:      n.args,
	}
}

func (n *FindBuilder) Parameters() []any {
	return n.args
}
