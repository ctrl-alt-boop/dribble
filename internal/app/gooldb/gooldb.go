package gooldb

import (
	"errors"

	"github.com/ctrl-alt-boop/gooldb/internal/app/internal/database"
	"github.com/ctrl-alt-boop/gooldb/pkg/connection"
	"github.com/ctrl-alt-boop/gooldb/pkg/logging"
)

var sqlLogger = logging.NewLogger("sql.log")

var logger *logging.Logger

var ErrNoDatabaseContext = errors.New("no database context")

type GoolDb struct {
	logger *logging.Logger
	ip     string

	settings        *connection.Settings
	databaseContext *database.DatabaseContext
	loadedDrivers   []database.DriverName

	onEvent func(eventType EventType, args any, err error)
}

func New(log *logging.Logger, ip string) *GoolDb {
	logger = log
	return &GoolDb{
		logger: log,
		ip:     ip,
		onEvent: func(eventType EventType, args any, err error) {
			logger.Infof("unhandled event: %s, %v, %v", eventType, args, err)
		},
		loadedDrivers: make([]database.DriverName, 0),
	}
}

func (gool *GoolDb) OnEvent(handler func(eventType EventType, args any, err error)) {
	gool.onEvent = handler
}

func (gool *GoolDb) Connect(settings *connection.Settings) {
	gool.settings = settings
	context, err := database.Connect(gool.settings)
	if err != nil {
		logger.Warn(err)
		gool.onEvent(ConnectError, nil, err)
		return
	}
	context.OnQueryExecuted(gool.onQueryExecuted)
	gool.databaseContext = context
	gool.onEvent(Connected, nil, nil)
}

func (gool *GoolDb) Reconnect(opts ...connection.Option) {
	gool.Disconnect()
	if gool.settings == nil {
		gool.onEvent(ReconnectError, nil, errors.New("no existing settings"))
		return
	}
	gool.settings = gool.settings.Copy(opts...)
	context, err := database.Connect(gool.settings)
	if err != nil {
		logger.Warn(err)
		gool.onEvent(ReconnectError, nil, err)
		return
	}
	context.OnQueryExecuted(gool.onQueryExecuted)
	gool.databaseContext = context
	gool.onEvent(Reconnected, nil, nil)
}

func (gool *GoolDb) Disconnect() {
	if gool.databaseContext != nil {
		err := gool.databaseContext.DB.Close()
		if err != nil {
			logger.Warn(err)
			gool.onEvent(DisconnectError, nil, err)
			return
		}
		gool.onEvent(Disconnected, nil, nil)
	}
	gool.databaseContext = nil
}

func (gool *GoolDb) onQueryExecuted(query string, err error) {
	sqlLogger.Info(query)
	if err != nil {
		sqlLogger.Error("\t", err)
	}
}

func (gool *GoolDb) OpenDatabase(database string) {
	if gool.databaseContext == nil {
		logger.Warn(ErrNoDatabaseContext)
		gool.onEvent(DBOpenError, nil, ErrNoDatabaseContext)
		return
	}
	gool.Reconnect(connection.WithDb(database))
	gool.onEvent(DBOpened, nil, nil)
}
