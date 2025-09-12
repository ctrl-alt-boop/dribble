package database

type (
	// SQLStyle interface {
	// 	Select() QueryIntentBuilder
	// 	Insert() QueryIntentBuilder
	// 	Update() QueryIntentBuilder
	// 	Delete() QueryIntentBuilder
	// 	Function() QueryIntentBuilder
	// 	Procedure() QueryIntentBuilder
	// }
	// NoSQLStyle interface {
	// 	Find() QueryIntentBuilder
	// 	Update() QueryIntentBuilder
	// 	Delete() QueryIntentBuilder
	// }

	// Build a query in the style of a SQL query
	SQLStyleIntent struct {
	}

	// Build a query in the style of a NoSQL query
	NoSQLStyleIntent struct {
	}
	NoSQLStyleIntentBuilder struct {
	}
)

func NewSQLStyleIntent() *SQLStyleIntent {
	return &SQLStyleIntent{}
}

func NewNoSQLStyleQueryIntentBuilder() *NoSQLStyleIntent {
	return &NoSQLStyleIntent{}
}

func (s *SQLStyleIntent) Insert(values ...any) *SQLStyleIntent {
	return s
}

func (s *SQLStyleIntent) Update(table string) *SQLStyleIntent {
	return s
}

func (s *SQLStyleIntent) Delete() *SQLStyleIntent {
	return s
}

func (s *SQLStyleIntent) Function(funcName string, args ...any) *SQLStyleIntent {
	return s
}

func (s *SQLStyleIntent) Procedure(procName string, args ...any) *SQLStyleIntent {
	return s
}

func (n *NoSQLStyleIntent) Delete() *NoSQLStyleIntent {
	return n
}

func (n *NoSQLStyleIntent) Find() *NoSQLStyleIntent {
	return n
}

func (n *NoSQLStyleIntent) Update() *NoSQLStyleIntent {
	return n
}
