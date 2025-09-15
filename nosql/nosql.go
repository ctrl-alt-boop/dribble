package nosql

import (
	"reflect"

	"github.com/ctrl-alt-boop/dribble/database"
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

func (n *FindBuilder) ToIntent() *database.Intent { // TODO: Ofcourse this needs implementation for the other operation types
	intent := &FindQuery{
		Collection:   n.collection,
		Conditions:   n.conditions,
		args:         n.args,
		LimitClause:  n.limitClause,
		OffsetClause: n.offsetClause,
	}
	return &database.Intent{
		Type:      database.Read,
		QueryType: reflect.TypeOf(intent),
		Operation: intent,
		Args:      n.args,
	}
}

func (n *FindBuilder) ToIntentOn(target *database.Target) *database.Intent { // TODO: Ofcourse this needs implementation for the other operation types
	intent := &FindQuery{
		Collection:   n.collection,
		Conditions:   n.conditions,
		args:         n.args,
		LimitClause:  n.limitClause,
		OffsetClause: n.offsetClause,
	}
	return &database.Intent{
		Target:    target,
		Type:      database.Read,
		QueryType: reflect.TypeOf(intent),
		Operation: intent,
		Args:      n.args,
	}
}

func (n *FindBuilder) Parameters() []any {
	return n.args
}
