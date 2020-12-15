package persistence

import (
	"database/sql"
	"fmt"
	"fp-dynamic-elements-manager-controller/internal/logging/structs"
	"github.com/jmoiron/sqlx"
	"time"
)

const (
	BatchTable = "element_batches"
)

type ElementBatchRepo struct {
	db  *sqlx.DB
	log *structs.AppLogger
}

func NewElementBatchRepo(appDb *sqlx.DB, logger *structs.AppLogger) *ElementBatchRepo {
	return &ElementBatchRepo{db: appDb, log: logger}
}

func (e *ElementBatchRepo) InsertBatchElement() (res sql.Result, err error) {
	now := time.Now()

	smt := fmt.Sprintf("INSERT INTO %s (created_at, updated_at) VALUES (?,?)", BatchTable)
	tx, err := e.db.Begin()
	if err != nil {
		e.log.SystemLogger.Error(err, "Error starting transaction inserting batch element")
		return
	}
	res, err = tx.Exec(smt, now, now)
	if err != nil {
		e.log.SystemLogger.Error(err, "Error inserting batch element, rolling back")
		tx.Rollback()
		return
	}

	err = tx.Commit()

	if err != nil {
		e.log.SystemLogger.Error(err, "Insert batch element commit failed")
		return
	}

	return
}

func (e *ElementBatchRepo) GetBatchIds() (receiver []int64, err error) {
	err = e.db.Select(&receiver, fmt.Sprintf("SELECT id FROM %s", BatchTable))
	return
}
