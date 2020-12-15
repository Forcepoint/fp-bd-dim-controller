package main

import (
	"fmt"
	"fp-dynamic-elements-manager-controller/api"
	"fp-dynamic-elements-manager-controller/internal/backup"
	"fp-dynamic-elements-manager-controller/internal/config"
	"fp-dynamic-elements-manager-controller/internal/db"
	"fp-dynamic-elements-manager-controller/internal/db/migration"
	"fp-dynamic-elements-manager-controller/internal/db/persistence"
	docker2 "fp-dynamic-elements-manager-controller/internal/docker"
	"fp-dynamic-elements-manager-controller/internal/logging"
	"fp-dynamic-elements-manager-controller/internal/logging/structs"
	"fp-dynamic-elements-manager-controller/internal/notification"
	"fp-dynamic-elements-manager-controller/internal/queue"
	"github.com/gammazero/workerpool"
	"github.com/sirupsen/logrus"
	"os"
)

// main is used to set up dependencies for the rest of the project,
// all deps are injected into the funcs/constructors that require them.
// main also starts the server
func main() {
	// This channel is used to signal that the DB is up and can be pinged healthily
	// this triggers the running of DB dependent functionality within the server logic
	dbReadyChan := make(chan struct{}, 1)

	database := db.NewAppDatabase(dbReadyChan)
	defer database.SqlDatabase.Close()

	// Set up the config and logging services
	config.InitViper("./config", "config", ".yml")
	logging.InitUserLogger(os.Getenv("LOG_LEVEL"), os.Getenv("LOG_FILE"))
	logging.InitInternalLogger()

	// Set up the NotificationService which allows the system to disseminate notifications to clients
	notificationService := notification.NewNotificationsService(notification.NewHub())

	// Set up the AppLogger which contains an instance of the UserLogger, SystemLogger, and NotificationService
	// This is used to send logs to the DB, Console and WebSocket, respectively.
	logger := structs.NewAppLogger(notificationService)

	// Set up and start the DB migration tool, on first run this creates all of the tables and constraints.
	// On subsequent runs it can alter tables if changes are defined in the .sql files in the /db folder
	migration.NewDBMigrator(database.SqlDatabase.DB, logger).RunMigration()

	// Set up the DAO, this is a holder for pointers to each of the separate entity repositories
	dao := persistence.NewDataAccessObject(database.SqlDatabase, logger)

	// Set up the hook for logrus which watches the UserLogger for events above a certain threshold and
	// writes them to the DB
	logrus.AddHook(logging.NewDatabaseHook(dao.LogEntryRepo))

	// Set up the pushing mechanism which pushes list elements to all egress modules
	pusher := queue.NewTableDataPusher(dao, logger)

	// Set up a worker pool of goroutines
	wp := workerpool.New(5)

	// Set up the connection to the docker socket on the host machine using the Docker CLI
	docker, err := docker2.NewDocker(
		os.Getenv("DOCKER_USER"),
		os.Getenv("DOCKER_PASSWORD"),
		fmt.Sprintf("https://%s", os.Getenv("DOCKER_REGISTRY")),
		logger,
	)

	if err != nil {
		logger.SystemLogger.Error(err, "error creating new docker")
	}

	// Set up the handler for incoming docker commands from the client
	handler := docker2.NewCommandHandler(docker, dao.ModuleMetadataRepo)

	// Set up the Backup/Restore provider
	provider := backup.NewDatabaseBackupProvider(
		docker,
		backup.NewGitController(logger),
		logger,
		dao.ListElementRepo)

	// Set up and start our server
	api.NewServer(logger, dbReadyChan, dao, pusher, wp, handler, provider).StartServer()
}
