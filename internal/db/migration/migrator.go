package migration

import (
	"database/sql"
	"fp-dynamic-elements-manager-controller/internal/logging/structs"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type Migrator interface {
	RunMigration()
}

type DBMigrator struct {
	db     *sql.DB
	logger *structs.AppLogger
}

func NewDBMigrator(db *sql.DB, logger *structs.AppLogger) *DBMigrator {
	return &DBMigrator{
		db:     db,
		logger: logger,
	}
}

func (m *DBMigrator) RunMigration() {
	driver, _ := mysql.WithInstance(m.db, &mysql.Config{})
	mig, err := migrate.NewWithDatabaseInstance(
		"file:///db/migrations",
		"mysql",
		driver)

	if err != nil {
		m.logger.SystemLogger.Error(err, "error getting new db instance for migration")
	}

	err = mig.Up()

	if err != nil {
		m.logger.SystemLogger.Error(err, "running migration")
	}
}
