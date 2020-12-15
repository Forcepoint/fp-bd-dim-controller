package db

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"os"
	"time"
)

// AppDatabase holds a reference to our sqlx.Db instance which in turn wraps the standard sql.DB instance
type AppDatabase struct {
	SqlDatabase *sqlx.DB
}

func NewAppDatabase(ready chan struct{}) *AppDatabase {
	appDb := &AppDatabase{}

	appDb.openDatabase()

	// Check to see if the DB is up and healthy, if not, panic and crash the container so it will restart and try again
	if !appDb.IsHealthy() {
		panic("Database not up, panicking to restart container")
	}

	// Indicate that the DB is up and healthy so dependent functionality can start
	ready <- struct{}{}

	return appDb
}

// openDatabase connects to the DB, sets the default values and sets the value of the sql DB in the AppDatabase struct
func (a *AppDatabase) openDatabase() {
	db, err := sqlx.Connect(
		"mysql",
		os.ExpandEnv("${MYSQL_USER}:${MYSQL_PASSWORD}@(mariadb)/${MYSQL_DATABASE}?charset=utf8&parseTime=True&loc=Local&multiStatements=true"),
	)

	if err != nil {
		panic(err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	a.SqlDatabase = db
}

// IsHealthy simply pings the sql DB for a basic health check
func (a *AppDatabase) IsHealthy() bool {
	return a.SqlDatabase.Ping() == nil
}
