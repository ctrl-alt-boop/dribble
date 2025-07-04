package database

import (
	"reflect"
)

type (
	Server struct { // This needs going over.
		Name string

		Engine    *Engine
		Databases Databases

		MaxConnections int
		Timeout        int
		SecurityModel  string
	}

	Engine struct {
		Name             string
		BuiltinFunctions []*Function

		Properties struct { // This needs going over.
			Version string

			TransactionSupport   []string
			ConcurrencyControl   string
			RecoveryModel        string
			StorageModel         string
			IntegrityEnforcement string

			ReplicationCapabilities []string
			IndexingCapabilities    []string
			Caching                 []string
		}
	}

	Database struct {
		Name        string
		Description string

		Engine *Engine

		Tables     []*Table
		Views      []*View
		Procedures []*Function
		Roles      []*Role

		Properties struct { // This needs going over.
			CharacterSet string
			Collation    string
			Locale       string

			ConnectionLimit int
			Owner           string

			DefaultStorageEngine string
			DefaultSchema        string
			DefaultNamespace     string

			SecurityModel                 string
			RecoveryModel, BackupStrategy string
		}
	}

	Table struct {
		Name        string
		Description string

		Engine *Engine

		Fields []*Field

		Properties struct { // This needs going over.
			PrimaryKeys []string

			Indexes []struct {
				Name        string
				Description string

				Fields    []string
				Type      string // B-tree, Hash, etc.
				Unique    bool
				Clustered bool
				Spatial   bool
			}
			ForeignKeys []struct {
				Name  string
				Field *Field

				ToTable string
				ToField string

				OnUpdate string
				OnDelete string
			}

			Partitioned bool
		}
	}

	Field struct {
		Name        string
		Description string

		Properties struct { // This needs going over.
			PrimaryKey bool

			Type      string
			Length    *int
			Precision *int
			Scale     *int

			Default   *string
			AutoInc   *string // Auto increment seed, nil if AutoIncrement false
			Virtual   bool
			Generated *string

			Unique        bool
			Index         bool
			Nullable      bool
			AllowedValues []string // name? for enum or set (MySQL)
			Encrypted     bool
			Overrides     struct {
				CharacterSet string
				Collation    string
			}
		}
	}

	View struct {
		Name        string
		Description string

		Fields []*Field

		Properties struct { // This needs going over.
			Definition   string
			Updatable    bool
			WithCheck    bool
			Materialized bool
		}
	}

	Function struct {
		Name        string
		Description string

		Properties struct { // This needs going over.
			Definition string
			Type       string
			Returns    *string
			Parameters []struct {
				Name        string
				Description string
				Type        string
				Direction   string // IN, OUT, INOUT
			}
		}
	}

	Role struct {
		Name        string
		Description string

		Permissions []string
		Inheritance []string
		Members     []string
	}

	// Collections
	Collection[T struct{}] []*T

	Servers   []*Server
	Engines   []*Engine
	Databases []*Database
	Tables    []*Table
	Views     []*View
	Roles     []*Role
	Functions []*Function
	Fields    []*Field
)

func (c Collection[T]) AsMap() map[string]*T { // Will try this sometime...
	items := make(map[string]*T)
	for _, item := range c {
		val := reflect.ValueOf(item).Elem()
		field := val.FieldByName("Name")
		if field.IsValid() {
			items[field.String()] = item
		}
	}
	return items
}

func (s *Server) DatabaseMap() map[string]*Database {
	return s.Databases.AsMap()
}

func (d Databases) AsMap() map[string]*Database {
	databases := make(map[string]*Database)
	for _, database := range d {
		databases[database.Name] = database
	}
	return databases
}

func (t Tables) AsMap() map[string]*Table {
	tables := make(map[string]*Table)
	for _, table := range t {
		tables[table.Name] = table
	}
	return tables
}

func (v Views) AsMap() map[string]*View {
	views := make(map[string]*View)
	for _, view := range v {
		views[view.Name] = view
	}
	return views
}

func (r Roles) AsMap() map[string]*Role {
	roles := make(map[string]*Role)
	for _, role := range r {
		roles[role.Name] = role
	}
	return roles
}

func (f Functions) AsMap() map[string]*Function {
	functions := make(map[string]*Function)
	for _, function := range f {
		functions[function.Name] = function
	}
	return functions
}

func (f Fields) AsMap() map[string]*Field {
	fields := make(map[string]*Field)
	for _, field := range f {
		fields[field.Name] = field
	}
	return fields
}
