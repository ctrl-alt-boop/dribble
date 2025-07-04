package database

func NewDatabaseNamesQuery() *QueryIntent {
	builder := Select("db_name").From("information_schema.schemata")
	return builder.ToQuery()
}
