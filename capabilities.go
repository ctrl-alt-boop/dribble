package dribble

import (
	"github.com/ctrl-alt-boop/dribble/datasource"
	_ "github.com/ctrl-alt-boop/dribble/internal/adapters/sql/mysql"
	_ "github.com/ctrl-alt-boop/dribble/internal/adapters/sql/postgres"
	_ "github.com/ctrl-alt-boop/dribble/internal/adapters/sql/sqlite3"
)

func (c *Client) SupportedDataSources() []datasource.Adapter {
	return datasource.Adapters()
}

func (c *Client) SupportedSourceTypes() []datasource.SourceType {
	return datasource.AdapterTypes()
}
