package persistence

import (
	"fmt"
	structs2 "fp-dynamic-elements-manager-controller/internal/logging/structs"
	"fp-dynamic-elements-manager-controller/internal/modules/structs"
	"github.com/jmoiron/sqlx"
	"strings"
	"time"
)

const (
	ModuleEndpointTablename = "module_endpoints"
)

type ModuleEndpointRepo struct {
	db  *sqlx.DB
	log *structs2.AppLogger
}

func NewModuleEndpointRepo(appDb *sqlx.DB, logger *structs2.AppLogger) *ModuleEndpointRepo {
	return &ModuleEndpointRepo{db: appDb, log: logger}
}

func (m *ModuleEndpointRepo) GetModuleEndpointsForModule(moduleId int64) (receiver []structs.ModuleEndpoint, err error) {
	err = m.db.Select(&receiver, fmt.Sprintf("SELECT * FROM %s WHERE module_metadata_id = ? ORDER BY created_at DESC LIMIT 4", ModuleEndpointTablename), moduleId)
	return
}

func (m *ModuleEndpointRepo) BatchInsertModuleEndpoints(endpoints []structs.ModuleEndpoint, moduleId int64) {
	if len(endpoints) == 0 {
		return
	}

	now := time.Now()

	for i, _ := range endpoints {
		endpoints[i].ModuleMetadataId = uint(moduleId)
		endpoints[i].CreatedAt = now
		endpoints[i].UpdatedAt = now
	}

	var valueStrings []string
	var valueArgs []interface{}
	for _, ep := range endpoints {
		valueStrings = append(valueStrings, "(?, ?, ?, ?, ?, ?, ?)")

		valueArgs = append(valueArgs, ep.ID)
		valueArgs = append(valueArgs, time.Now())
		valueArgs = append(valueArgs, time.Now())
		valueArgs = append(valueArgs, ep.DeletedAt)
		valueArgs = append(valueArgs, ep.Secure)
		valueArgs = append(valueArgs, ep.Endpoint)
		valueArgs = append(valueArgs, ep.ModuleMetadataId)
	}

	smt := `INSERT INTO %s (id, created_at, updated_at, deleted_at, secure, endpoint, module_metadata_id) VALUES %s`
	smt = fmt.Sprintf(smt, ModuleEndpointTablename, strings.Join(valueStrings, ","))
	tx, err := m.db.Begin()
	if err != nil {
		m.log.SystemLogger.Error(err, "Error starting transaction to insert module endpoints")
		return
	}
	_, err = tx.Exec(smt, valueArgs...)
	if err != nil {
		m.log.SystemLogger.Error(err, "Error inserting module endpoints, rolling back")
		tx.Rollback()
		return
	}

	err = tx.Commit()

	if err != nil {
		m.log.SystemLogger.Error(err, "Error committing insert module endpoints")
		return
	}

	return
}
