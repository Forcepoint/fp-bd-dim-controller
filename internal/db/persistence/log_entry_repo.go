package persistence

import (
	"fmt"
	"fp-dynamic-elements-manager-controller/internal/logging/structs"
	"github.com/jmoiron/sqlx"
	"time"
)

const (
	LogTable = "log_entries"
)

type LogEntryRepo struct {
	db  *sqlx.DB
	log *structs.AppLogger
}

func NewLogEntryRepo(appDb *sqlx.DB, logger *structs.AppLogger) *LogEntryRepo {
	return &LogEntryRepo{db: appDb, log: logger}
}

func (l *LogEntryRepo) InsertLogEntry(item *structs.LogEntry) {
	now := time.Now()

	smt := fmt.Sprintf("INSERT INTO %s (id, created_at, updated_at, deleted_at, module_name, level, message, caller, time) VALUES (?,?,?,?,?,?,?,?,?)", LogTable)
	tx, err := l.db.Begin()
	if err != nil {
		l.log.SystemLogger.Error(err, "Error starting transaction to insert log entry")
		return
	}
	_, err = tx.Exec(smt, item.ID, now, now, item.DeletedAt, item.ModuleName, item.Level, item.Message, item.Caller, item.Time)
	if err != nil {
		l.log.SystemLogger.Error(err, "Error inserting log entry, rolling back")
		tx.Rollback()
		return
	}

	err = tx.Commit()

	if err != nil {
		l.log.SystemLogger.Error(err, "Error committing insert log entry")
		return
	}

	return
}

func (l *LogEntryRepo) GetAll() (receiver []structs.LogEntry, err error) {
	err = l.db.Select(&receiver, fmt.Sprintf("SELECT * FROM %s ORDER BY created_at DESC;", LogTable))
	return
}

func (l *LogEntryRepo) GetAllForLevels(moduleName string, pageSize, offset int, levels []string) (receiver []structs.LogEntry, err error) {
	var smt string
	if moduleName == "" {
		smt = fmt.Sprintf("SELECT * FROM %s WHERE level IN (?) ORDER BY created_at DESC LIMIT ? OFFSET ?;", LogTable)
	} else {
		smt = fmt.Sprintf("SELECT * FROM %s WHERE level IN (?) AND module_name = ? ORDER BY created_at DESC LIMIT ? OFFSET ?;", LogTable)
	}

	query, args, err := sqlx.In(smt, levels, pageSize, offset)
	if err != nil {
		l.log.SystemLogger.Error(err, "Error binding args to query")
	}

	query = l.db.Rebind(query)

	err = l.db.Select(&receiver, query, args...)

	return
}

func (l *LogEntryRepo) GetTotalCount() (total int64, err error) {
	err = l.db.Get(&total, fmt.Sprintf("SELECT count(1) FROM %s", LogTable))
	return
}
