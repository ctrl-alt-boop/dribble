package datasource

import "sync"

// Adapter describes a registered data source.
type Adapter struct {
	Name         string
	Type         SourceType
	Properties   map[string]string
	Capabilities Capabilities
	Metadata     Metadata
	FactoryFunc  func(Namer) DataSource
}

var (
	registry = make(map[SourceType]Adapter)
	mu       sync.RWMutex
)

// Register makes a data source available by name.
// If Register is called twice with the same name or if driver is nil,
// it panics.
func Register(adapter Adapter) {
	if adapter.Name == "" {
		panic("Register adapter with no name")
	}
	if adapter.Type == "" {
		panic("Register adapter with no type")
	}
	if adapter.FactoryFunc == nil {
		panic("Register adapter with no factory function")
	}

	mu.Lock()
	defer mu.Unlock()
	if _, dup := registry[adapter.Type]; dup {
		panic("Register called twice for adapter " + adapter.Name)
	}
	registry[adapter.Type] = adapter
}

// Adapters returns a list of the registered data source adapters.
func Adapters() []Adapter {
	mu.RLock()
	defer mu.RUnlock()
	list := make([]Adapter, 0, len(registry))
	for _, adapter := range registry {
		list = append(list, adapter)
	}
	return list
}

// AdapterTypes returns a list of the registered data source adapter types.
func AdapterTypes() []SourceType {
	mu.RLock()
	defer mu.RUnlock()
	list := make([]SourceType, 0, len(registry))
	for t := range registry {
		list = append(list, t)
	}
	return list
}

func GetAdapter(sourceType SourceType) (Adapter, bool) {
	mu.RLock()
	defer mu.RUnlock()
	adapter, ok := registry[sourceType]
	return adapter, ok
}
