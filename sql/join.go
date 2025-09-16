package sql

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
