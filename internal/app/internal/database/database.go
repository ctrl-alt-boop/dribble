package database

import (
	"fmt"

	"github.com/ctrl-alt-boop/gooldb/pkg/connection"
	"github.com/ctrl-alt-boop/gooldb/pkg/logging"
)

var logger = logging.NewLogger("db.log")

func Connect(settings *connection.Settings) (*DatabaseContext, error) {
	driver, err := createDriver(settings.DriverName)
	if err != nil {
		logger.ErrorF("%s", err.Error())
		return nil, fmt.Errorf("error creating %s driver: %w", settings.DriverName, err)
	}
	connString := driver.ConnectionString(settings)
	logger.Info(fmt.Sprintf("Connecting: %+v, (%s)", settings, connString))
	return createDatabaseContext(settings.DriverName, driver, connString)
}
