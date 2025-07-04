package sql

const (
	Select            = "SELECT"
	Insert            = "INSERT"
	Update            = "UPDATE"
	Delete            = "DELETE"
	From              = "FROM"
	Where             = "WHERE"
	Set               = "SET"
	Values            = "VALUES"
	OrderBy           = "ORDER BY"
	Asc               = "ASC"
	Desc              = "DESC"
	Limit             = "LIMIT"
	Offset            = "OFFSET"
	GroupBy           = "GROUP BY"
	Having            = "HAVING"
	Join              = "JOIN"
	LeftJoin          = "LEFT JOIN"
	RightJoin         = "RIGHT JOIN"
	FullJoin          = "FULL JOIN"
	CrossJoin         = "CROSS JOIN"
	On                = "ON"
	As                = "AS"
	InnerJoin         = "INNER JOIN"
	OuterJoin         = "OUTER JOIN"
	Union             = "UNION"
	Intersect         = "INTERSECT"
	Except            = "EXCEPT"
	UnionAll          = "UNION ALL"
	IntersectAll      = "INTERSECT ALL"
	ExceptAll         = "EXCEPT ALL"
	Not               = "NOT"
	In                = "IN"
	Between           = "BETWEEN"
	And               = "AND"
	Or                = "OR"
	IsNull            = "IS NULL"
	IsNotNull         = "IS NOT NULL"
	IsTrue            = "IS TRUE"
	IsFalse           = "IS FALSE"
	IsUnknown         = "IS UNKNOWN"
	IsDistinctFrom    = "IS DISTINCT FROM"
	IsNotDistinctFrom = "IS NOT DISTINCT FROM"
	Like              = "LIKE"
	NotLike           = "NOT LIKE"
	Ilike             = "ILIKE"
	NotIlike          = "NOT ILIKE"
	Any               = "ANY"
	All               = "ALL"
	Exists            = "EXISTS"
	Some              = "SOME"
	Unique            = "UNIQUE"
	PrimaryKey        = "PRIMARY KEY"
	ForeignKey        = "FOREIGN KEY"
	Check             = "CHECK"
	Default           = "DEFAULT"
	Null              = "NULL"
	True              = "TRUE"
	False             = "FALSE"
	Unknown           = "UNKNOWN"
)

type Method string

func (s Method) String() string {
	return string(s)
}

const (
	MethodSelect  Method = "SELECT"
	MethodInsert  Method = "INSERT"
	MethodUpdate  Method = "UPDATE"
	MethodDelete  Method = "DELETE"
	MethodCall    Method = "CALL"
	MethodExec    Method = "EXEC"
	MethodExecute Method = "EXECUTE"
)

const DefaultSelectLimit int = 10 // Just a safeguard

var SQLMethods = []Method{MethodSelect, MethodInsert, MethodUpdate, MethodDelete}

const DefaultSQLSelectTemplate = `SELECT {{if .AsDistinct}}DISTINCT{{end}}
{{- range $i, $field := .Fields -}}
    {{if $i}}, {{end}}{{$field}}
{{- end}}
FROM {{.Table}}
{{- range .Joins}}
    {{.Type}} JOIN {{.Table}} ON {{.On}}
{{- end}}
{{- if .WhereClause}}
WHERE {{.WhereClause}}
{{- end}}
{{- if .GroupByClause}}
GROUP BY {{range $i, $field := .GroupByClause}}{{if $i}}, {{end}}{{$field}}{{end}}
{{- end}}
{{- if .HavingClause}}
HAVING {{.HavingClause}}
{{- end}}
{{- if .OrderByClause}}
ORDER BY {{range $i, $field := .OrderByClause}}{{if $i}}, {{end}}{{$field}}{{if .DescClause}} DESC{{end}}{{end}}
{{- end}}
{{- if .LimitClause}}
LIMIT {{.LimitClause}}
{{- end}}
{{- if .OffsetClause}}
OFFSET {{.OffsetClause}}
{{- end}}`
