package persistence

import (
	"fmt"
	structs2 "fp-dynamic-elements-manager-controller/internal/logging/structs"
	"fp-dynamic-elements-manager-controller/internal/queue/structs"
	"github.com/jmoiron/sqlx"
	"time"
)

const (
	UpdateStatusTable = "update_statuses"
)

type UpdateStatusRepo struct {
	db  *sqlx.DB
	log *structs2.AppLogger
}

func NewUpdateStatusRepo(appDb *sqlx.DB, logger *structs2.AppLogger) *UpdateStatusRepo {
	return &UpdateStatusRepo{db: appDb, log: logger}
}

func (u *UpdateStatusRepo) InsertUpdateStatus(item structs.UpdateStatus) {
	if u.exists(item.ServiceName, item.UpdateBatchId) {
		return
	}

	now := time.Now()

	smt := fmt.Sprintf(`INSERT INTO %s (id, created_at, updated_at, deleted_at, service_name, status, update_batch_id, module_metadata_id) VALUES (?,?,?,?,?,?,?,?)`, UpdateStatusTable)
	tx, err := u.db.Begin()
	if err != nil {
		u.log.SystemLogger.Error(err, "Error starting transaction to insert update status")
		return
	}
	_, err = tx.Exec(smt, item.ID, now, now, item.DeletedAt, item.ServiceName, item.Status, item.UpdateBatchId, item.ModuleMetadataId)
	if err != nil {
		tx.Rollback()
		u.log.SystemLogger.Error(err, "Error inserting update status, rolling back")
		return
	}

	err = tx.Commit()

	if err != nil {
		u.log.SystemLogger.Error(err, "Error committing insert update status")
		return
	}

	return
}

func (u *UpdateStatusRepo) UpdateUpdateStatus(item structs.UpdateStatus) {
	now := time.Now()

	smt := fmt.Sprintf("UPDATE %s SET updated_at = ?, status = ? WHERE service_name = ? AND update_batch_id = ?", UpdateStatusTable)
	tx, err := u.db.Begin()
	if err != nil {
		u.log.SystemLogger.Error(err, "Error starting transaction to update update status")
		return
	}
	_, err = tx.Exec(smt, now, item.Status, item.ServiceName, item.UpdateBatchId)
	if err != nil {
		tx.Rollback()
		u.log.SystemLogger.Error(err, "Error updating update status, rolling back")
		return
	}

	err = tx.Commit()

	if err != nil {
		u.log.SystemLogger.Error(err, "Error committing update update status")
		return
	}

	return
}

func (u *UpdateStatusRepo) GetAll(receiver []structs.UpdateStatus) error {
	return u.db.Select(&receiver, fmt.Sprintf("SELECT * FROM %s ORDER BY created_at DESC;", UpdateStatusTable))
}

func (u *UpdateStatusRepo) GetAllWithStatusForModule(status structs.Status, moduleId int64, safe bool) (receiver []int64, err error) {
	err = u.db.Select(&receiver, fmt.Sprintf("SELECT update_batch_id FROM %s WHERE safe = ? status = ? AND module_metadata_id = ? ORDER BY created_at DESC;", UpdateStatusTable), safe, status.String(), moduleId)
	return
}

func (u *UpdateStatusRepo) GetAllWithStatusPaginated(offset, pageSize int, status []string) (receiver []structs.UpdateStatus, err error) {
	smt := fmt.Sprintf("SELECT * FROM %s WHERE status IN (?) ORDER BY created_at DESC LIMIT ? OFFSET ?;", UpdateStatusTable)
	query, args, err := sqlx.In(smt, status, pageSize, offset)
	if err != nil {
		u.log.SystemLogger.Error(err, "Error binding query parameters update status")
	}

	query = u.db.Rebind(query)

	err = u.db.Select(&receiver, query, args...)
	return
}

func (u *UpdateStatusRepo) GetTotalCount() (total int64, err error) {
	err = u.db.Get(&total, fmt.Sprintf("SELECT count(1) FROM %s", UpdateStatusTable))
	return
}

func (u *UpdateStatusRepo) GetLatestUpdate(moduleId int64) (receiver structs.UpdateStatus, err error) {
	err = u.db.Get(&receiver, fmt.Sprintf("SELECT * FROM %s WHERE module_metadata_id = ? AND (status = ? OR status = ?) ORDER BY id DESC LIMIT 1;", UpdateStatusTable), moduleId, structs.PENDING.String(), structs.SUCCESS.String())
	return
}

func (u *UpdateStatusRepo) GetAllFailed(moduleId int64) (receiver []structs.UpdateStatus, err error) {
	err = u.db.Get(&receiver, fmt.Sprintf("SELECT * FROM %s WHERE module_metadata_id = ? AND status = ? ORDER BY id DESC;", UpdateStatusTable), moduleId, structs.FAILED.String())
	return
}

func (u *UpdateStatusRepo) exists(svcName string, updateBatchId int64) bool {
	var status structs.UpdateStatus
	return u.db.Get(&status, fmt.Sprintf("SELECT * FROM %s WHERE update_batch_id = ? AND service_name = ? LIMIT 1;", UpdateStatusTable), updateBatchId, svcName) == nil
}
