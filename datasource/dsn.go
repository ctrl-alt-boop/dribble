package datasource

type Namer interface {
	// Type of data source, database, etc.
	Type() Type

	SourceType() SourceType

	// machine-readable data source string
	DSN() string

	// Human-readable data source representation
	Info() string
}

type DataSourceNameRequest interface { // Something like that?
	Namer
	Request
}
