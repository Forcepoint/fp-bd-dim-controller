package persistence

import (
	"database/sql"
	"errors"
	"fmt"
	structs3 "fp-dynamic-elements-manager-controller/internal/logging/structs"
	"fp-dynamic-elements-manager-controller/internal/queue/structs"
	structs2 "fp-dynamic-elements-manager-controller/internal/stats/structs"
	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"strings"
	"time"
)

const (
	ElementsTable          = "list_elements"
	MySqlErrDuplicateValue = 1062
)

var ErrDuplicateValue = errors.New("duplicate value")

type ElementRepo interface {
	GetTotalElementCount() (int64, error)
}

type ListElementRepo struct {
	db  *sqlx.DB
	log *structs3.AppLogger
}

func NewListElementRepo(appDb *sqlx.DB, logger *structs3.AppLogger) *ListElementRepo {
	return &ListElementRepo{db: appDb, log: logger}
}

func (l *ListElementRepo) BatchInsertListElements(items []structs.ListElement) {
	var valueStrings []string
	var valueArgs []interface{}
	for _, element := range items {
		valueStrings = append(valueStrings, "(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")

		valueArgs = append(valueArgs, element.ID)
		valueArgs = append(valueArgs, time.Now())
		valueArgs = append(valueArgs, time.Now())
		valueArgs = append(valueArgs, element.DeletedAt)
		valueArgs = append(valueArgs, element.Source)
		valueArgs = append(valueArgs, element.ServiceName)
		valueArgs = append(valueArgs, element.Type)
		valueArgs = append(valueArgs, element.Value)
		valueArgs = append(valueArgs, element.Safe)
		valueArgs = append(valueArgs, element.UpdateBatchId)
	}

	smt := `INSERT INTO %s 
					(id, created_at, updated_at, deleted_at, source, service_name, type, value, safe, update_batch_id) 
					VALUES %s 
					ON DUPLICATE KEY UPDATE updated_at = VALUES(updated_at)`
	smt = fmt.Sprintf(smt, ElementsTable, strings.Join(valueStrings, ","))
	tx, err := l.db.Begin()
	if err != nil {
		l.log.SystemLogger.Error(err, "Error starting transaction to batch insert list elements")
		return
	}
	_, err = tx.Exec(smt, valueArgs...)
	if err != nil {
		l.log.SystemLogger.Error(err, "Error batch inserting list elements, rolling back")
		tx.Rollback()
		return
	}

	err = tx.Commit()

	if err != nil {
		l.log.SystemLogger.Error(err, "Error committing batch insert list elements")
		return
	}

	return
}

func (l *ListElementRepo) InsertListElement(item structs.ListElement) error {
	if l.exists(item.Value) {
		return ErrDuplicateValue
	}
	var valueArgs []interface{}

	valueArgs = append(valueArgs, item.ID)
	valueArgs = append(valueArgs, time.Now())
	valueArgs = append(valueArgs, time.Now())
	valueArgs = append(valueArgs, item.DeletedAt)
	valueArgs = append(valueArgs, item.Source)
	valueArgs = append(valueArgs, item.ServiceName)
	valueArgs = append(valueArgs, item.Type)
	valueArgs = append(valueArgs, item.Value)
	valueArgs = append(valueArgs, item.Safe)
	valueArgs = append(valueArgs, item.UpdateBatchId)

	smt := `INSERT INTO %s (id, created_at, updated_at, deleted_at, source, service_name, type, value, safe, update_batch_id)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	smt = fmt.Sprintf(smt, ElementsTable)
	tx, err := l.db.Begin()
	if err != nil {
		l.log.SystemLogger.Error(err, "Error starting transaction to insert list element")
		return err
	}
	_, err = tx.Exec(smt, valueArgs...)
	if err != nil {
		l.log.SystemLogger.Error(err, "Error inserting list element, rolling back")
		tx.Rollback()
		return err
	}

	err = tx.Commit()

	if err != nil {
		l.log.SystemLogger.Error(err, "Error committing insert list element")
		return err
	}

	return err
}

func (l *ListElementRepo) UpdateListElement(item structs.ListElement) error {
	var valueArgs []interface{}

	valueArgs = append(valueArgs, time.Now())
	valueArgs = append(valueArgs, item.Value)
	valueArgs = append(valueArgs, item.Safe)
	valueArgs = append(valueArgs, item.ID)

	smt := `UPDATE %s SET updated_at = ?, value = ?, safe = ? WHERE id = ?`
	smt = fmt.Sprintf(smt, ElementsTable)
	tx, err := l.db.Begin()
	if err != nil {
		l.log.SystemLogger.Error(err, "Error starting transaction to update list element")
		if err.(*mysql.MySQLError).Number == MySqlErrDuplicateValue {
			return ErrDuplicateValue
		}
		return err
	}
	_, err = tx.Exec(smt, valueArgs...)
	if err != nil {
		l.log.SystemLogger.Error(err, "Error updating list element, rolling back")
		tx.Rollback()
		if err.(*mysql.MySQLError).Number == MySqlErrDuplicateValue {
			return ErrDuplicateValue
		}
		return err
	}

	err = tx.Commit()

	if err != nil {
		l.log.SystemLogger.Error(err, "Error committing update list element")
		if err.(*mysql.MySQLError).Number == MySqlErrDuplicateValue {
			return ErrDuplicateValue
		}
		return err
	}

	return nil
}

func (l *ListElementRepo) DeleteByValue(value string) {
	smt := fmt.Sprintf(`DELETE FROM %s WHERE value = ?`, ElementsTable)
	tx, err := l.db.Begin()
	if err != nil {
		l.log.SystemLogger.Error(err, "Error starting transaction to delete list element")
		return
	}
	_, err = tx.Exec(smt, value)
	if err != nil {
		l.log.SystemLogger.Error(err, "Error deleting list element, rolling back")
		err = tx.Rollback()
		return
	}

	err = tx.Commit()

	if err != nil {
		l.log.SystemLogger.Error(err, "Error committing delete list element")
		return
	}
}

func (l *ListElementRepo) GetById(id int) (receiver []structs.ListElement, err error) {
	err = l.db.Select(&receiver, fmt.Sprintf("SELECT * FROM %s WHERE id = ?;", ElementsTable), id)
	return
}

func (l *ListElementRepo) GetAll() (receiver []structs.ListElement, err error) {
	err = l.db.Select(&receiver, fmt.Sprintf("SELECT * FROM %s ORDER BY created_at DESC;", ElementsTable))
	return
}

func (l *ListElementRepo) GetAllByBatchId(batchId int64, types []structs.ElementType, safe bool) (receiver []structs.ListElement, err error) {
	smt := fmt.Sprintf("SELECT * FROM %s WHERE update_batch_id = ? AND type IN (?) AND safe = ? ORDER BY created_at DESC;", ElementsTable)

	query, args, err := sqlx.In(smt, batchId, types, safe)

	if err != nil {
		l.log.SystemLogger.Error(err, "Error binding args to query")
	}

	query = l.db.Rebind(query)

	err = l.db.Select(&receiver, query, args...)
	return
}

func (l *ListElementRepo) GetAllPaginated(offset, pageSize int, safe bool) (receiver []structs.ListElement, err error) {
	err = l.db.Select(&receiver, fmt.Sprintf("SELECT * FROM %s WHERE safe = ? ORDER BY created_at DESC LIMIT ? OFFSET ?;", ElementsTable), safe, pageSize, offset)
	return
}

func (l *ListElementRepo) GetAllLike(offset, pageSize int, searchTerm string, safe bool) (receiver []structs.ListElement, err error) {
	err = l.db.Select(&receiver, fmt.Sprintf("SELECT * FROM %s WHERE value LIKE ? AND safe = ? ORDER BY created_at DESC LIMIT ? OFFSET ?;", ElementsTable), "%"+searchTerm+"%", safe, pageSize, offset)
	return
}

func (l *ListElementRepo) GetTotalCountWhereLike(safe bool, like string) (total int64, err error) {
	err = l.db.Get(&total, fmt.Sprintf("SELECT count(1) FROM %s WHERE value LIKE ? AND safe = ?", ElementsTable), "%"+like+"%", safe)
	return
}

func (l *ListElementRepo) GetTotalCount(safe bool) (total int64, err error) {
	err = l.db.Get(&total, fmt.Sprintf("SELECT count(1) FROM %s WHERE safe = ?", ElementsTable), safe)
	return
}

func (l *ListElementRepo) GetTotalElementCount() (total int64, err error) {
	err = l.db.Get(&total, fmt.Sprintf("SELECT count(1) FROM %s", ElementsTable))
	return
}

func (l *ListElementRepo) GetLatestUpdate(svcName string) (receiver structs.ListElement, err error) {
	err = l.db.Get(&receiver, fmt.Sprintf("SELECT * FROM %s WHERE service_name = ? ORDER BY created_at DESC LIMIT 1;", ElementsTable), svcName)
	return
}

func (l *ListElementRepo) GetUnpushedBatchIds(moduleId int64, types []structs.ElementType) (receiver []int64, err error) {
	smt := fmt.Sprintf("SELECT DISTINCT(update_batch_id) FROM %s WHERE type IN (?) AND update_batch_id NOT IN (SELECT update_batch_id FROM update_statuses WHERE module_metadata_id = ?);", ElementsTable)
	query, args, err := sqlx.In(smt, types, moduleId)

	if err != nil {
		l.log.SystemLogger.Error(err, "Error binding args to query")
	}

	query = l.db.Rebind(query)

	err = l.db.Select(&receiver, query, args...)
	return
}

func (l *ListElementRepo) GetStats(serviceName string) (result structs2.Stats) {
	var smt string
	if serviceName == "" {
		smt = fmt.Sprintf(`SELECT 
				count(1) 						AS total,
				sum(IF(type = 'IP', 1, 0))     AS ip,
				sum(IF(type = 'DOMAIN', 1, 0)) AS domain,
				sum(IF(type = 'URL', 1, 0))    AS url
				FROM %s
				GROUP BY service_name is not null`, ElementsTable)

	} else {
		smt = fmt.Sprintf(`SELECT 
				count(1) 						AS total,
				sum(IF(type = 'IP', 1, 0))     AS ip,
				sum(IF(type = 'DOMAIN', 1, 0)) AS domain,
				sum(IF(type = 'URL', 1, 0))    AS url
				FROM %s
				WHERE service_name = '%s'
				GROUP BY service_name`, ElementsTable, serviceName)
	}

	stats := structs2.Stats{}
	err := l.db.Get(&stats, smt)
	if err == sql.ErrNoRows {
		return
	}
	if err != nil {
		l.log.SystemLogger.Error(err, "Error getting stats for list elements")
		return
	}

	return stats
}

func (l *ListElementRepo) exists(value string) bool {
	var element structs.ListElement
	return l.db.Get(&element, fmt.Sprintf("SELECT * FROM %s WHERE value = ? LIMIT 1;", ElementsTable), value) == nil
}
