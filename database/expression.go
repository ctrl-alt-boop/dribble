package database

import (
	"fmt"
	"strings"
)

type ResultType int

const (
	Result ResultType = iota - 1
	ResultScalar
	ResultList
	ResultTable
	ResultSet
)

var _ Expr = &compExpr{}
var _ Expr = &logicExpr{}
var _ Expr = &notExpr{}

var Parameterized = false

type (
	Expr interface {
		ToSql() (string, []any)
	}
	Exprs []Expr

	compExpr struct {
		column string
		op     string
		value  any
		wordOp bool
	}

	logicExpr struct {
		op    string
		exprs []Expr
	}

	notExpr struct {
		expr *compExpr
	}
)

func (e Exprs) ToSql() (string, []any) {
	parts := make([]string, 0, len(e))
	var params []any
	for _, expr := range e {
		part, partParams := expr.ToSql()
		parts = append(parts, part)
		params = append(params, partParams...)
	}
	return strings.Join(parts, " AND "), params
}

func (n *notExpr) ToSql() (string, []any) {
	if n.expr.wordOp {
		return fmt.Sprintf("%s NOT %s ?", n.expr.column, n.expr.op), []any{n.expr.value}
	}
	return fmt.Sprintf("%s !%s ?", n.expr.column, n.expr.op), nil
}

func (c *compExpr) ToSql() (string, []any) {
	if Parameterized {
		return fmt.Sprintf("%s %s ?", c.column, c.op), []any{c.value}
	}
	value := c.resolveValueType()
	return fmt.Sprintf("%s %s %v", c.column, c.op, value), nil
}

func (c *compExpr) resolveValueType() string {
	switch c.value.(type) {
	case string:
		return fmt.Sprintf("'%s'", c.value)
	case bool:
		return fmt.Sprintf("%t", c.value)
	default:
		return fmt.Sprint(c.value)
	}
}

func (l *logicExpr) ToSql() (string, []any) {
	parts := make([]string, 0, len(l.exprs))
	var params []any
	for _, expr := range l.exprs {
		part, partParams := expr.ToSql()
		parts = append(parts, part)
		params = append(params, partParams...)
	}
	return fmt.Sprintf("(%s)", strings.Join(parts, fmt.Sprintf(" %s ", l.op))), params
}

func Eq(column string, value any) *compExpr { return &compExpr{column: column, op: "=", value: value} }
func Ne(column string, value any) *compExpr { return &compExpr{column: column, op: "!=", value: value} }
func Gt(column string, value any) *compExpr { return &compExpr{column: column, op: ">", value: value} }
func Lt(column string, value any) *compExpr { return &compExpr{column: column, op: "<", value: value} }
func Like(column string, value any) *compExpr {
	return &compExpr{column: column, op: "LIKE", value: value, wordOp: true}
}
func Null(column string) *compExpr {
	return &compExpr{column: column, op: "IS", value: "NULL", wordOp: true}
}

func And(exprs ...Expr) *logicExpr { return &logicExpr{op: "AND", exprs: exprs} }
func Or(exprs ...Expr) *logicExpr  { return &logicExpr{op: "OR", exprs: exprs} }

func Not(expr *compExpr) *notExpr { return &notExpr{expr: expr} }
