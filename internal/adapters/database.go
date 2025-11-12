package adapters

import (
	"errors"

	"github.com/ctrl-alt-boop/dribble/datasource"
)

func init() {
}

type Database struct{}

func (d *Database) Path() []string {
	return []string{"Database"}
}

func (d *Database) DataSourceType() datasource.SourceType {
	return datasource.SourceType("Database")
}

func Create(sourceType datasource.SourceType) (func(datasource.Namer) datasource.DataSource, error) {
	if adapter, ok := datasource.GetAdapter(sourceType); ok {
		return adapter.FactoryFunc, nil
	}
	return nil, errors.New("unsupported datasource")
}
