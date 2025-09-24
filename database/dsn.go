package database

type DataSourceNamer interface {

	// machine-readable data source string
	DSN() string

	// Type of data source, database, etc.
	Type() Type

	// Human-readable data source representation
	Info() string
}

type DataSourceNameRequest interface { // Something like that?
	DataSourceNamer
	Request
}
