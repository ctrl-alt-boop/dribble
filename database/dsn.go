package database

type DataSourceNamer interface {
	DSN() string
	Type() Type
}

type DataSourceNamerRequest interface { // Something like that?
	DataSourceNamer
	Request
}
