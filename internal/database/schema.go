package database

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

	Servers   []*Server
	Engines   []*Engine
	Databases []*Database
	Tables    []*Table
	Views     []*View
	Roles     []*Role
	Functions []*Function
	Fields    []*Field
)

type Namer interface {
	GetName() string
}

func (s *Server) GetName() string   { return s.Name }
func (d *Database) GetName() string { return d.Name }
func (t *Table) GetName() string    { return t.Name }
func (v *View) GetName() string     { return v.Name }
func (r *Role) GetName() string     { return r.Name }
func (f *Function) GetName() string { return f.Name }
func (f *Field) GetName() string    { return f.Name }

func CollectionToMap[T Namer](collection []T) map[string]T {
	m := make(map[string]T, len(collection))
	for _, item := range collection {
		m[item.GetName()] = item
	}
	return m
}

func (s *Server) DatabaseMap() map[string]*Database {
	return CollectionToMap(s.Databases)
}

func (d Databases) AsMap() map[string]*Database {
	return CollectionToMap(d)
}

func (t Tables) AsMap() map[string]*Table {
	return CollectionToMap(t)
}

func (v Views) AsMap() map[string]*View {
	return CollectionToMap(v)
}

func (r Roles) AsMap() map[string]*Role {
	return CollectionToMap(r)
}

func (f Functions) AsMap() map[string]*Function {
	return CollectionToMap(f)
}

func (f Fields) AsMap() map[string]*Field {
	return CollectionToMap(f)
}
