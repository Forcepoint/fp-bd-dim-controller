package persistence

import (
	"fmt"
	structs2 "fp-dynamic-elements-manager-controller/internal/logging/structs"
	"fp-dynamic-elements-manager-controller/internal/modules/structs"
	"github.com/jmoiron/sqlx"
	"time"
)

const (
	ModuleTable = "module_metadata"
)

type ModuleRepo interface {
	DeleteByServiceName(string) error
}

type ModuleMetadataRepo struct {
	db       *sqlx.DB
	log      *structs2.AppLogger
	epRepo   *ModuleEndpointRepo
	typeRepo *ElementTypeRepo
}

func NewModuleMetadataRepo(appDb *sqlx.DB,
	logger *structs2.AppLogger,
	endpointRepo *ModuleEndpointRepo,
	elementTypeRepo *ElementTypeRepo) *ModuleMetadataRepo {
	return &ModuleMetadataRepo{db: appDb, log: logger, epRepo: endpointRepo, typeRepo: elementTypeRepo}
}

func (m *ModuleMetadataRepo) GetAllModuleMetadata() (receiver []structs.ModuleMetadata, err error) {
	err = m.db.Select(&receiver, "SELECT * FROM module_metadata")

	if err != nil {
		return nil, err
	}

	for i, module := range receiver {
		receiver[i].ModuleEndpoints, err = m.epRepo.GetModuleEndpointsForModule(module.ID)
	}

	return
}

func (m *ModuleMetadataRepo) UpsertModuleMetadata(item structs.ModuleMetadata) {
	now := time.Now()

	smt := `INSERT INTO %s (
id, created_at, updated_at, deleted_at, module_service_name, module_display_name, module_type,
module_description, inbound_route, internal_ip, internal_port, icon_url, configured, configurable,
last_ping) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?) 
ON DUPLICATE KEY UPDATE updated_at = ?, module_description = ?, icon_url = ?, configured = ?, last_ping = ?
`

	item.CreatedAt = now
	item.UpdatedAt = now

	smt = fmt.Sprintf(smt, ModuleTable)
	tx, err := m.db.Begin()
	if err != nil {
		m.log.SystemLogger.Error(err, "Error starting transaction to upsert module metadata")
		return
	}

	res, err := tx.Exec(smt, item.ID, item.CreatedAt, item.UpdatedAt, item.DeletedAt, item.ModuleServiceName,
		item.ModuleDisplayName, item.ModuleType, item.ModuleDescription, item.InboundRoute, item.InternalIP,
		item.InternalPort, item.IconURL, item.Configured, item.Configurable, item.LastPing, now, item.ModuleDescription,
		item.IconURL, item.Configured, item.LastPing)

	if err != nil {
		m.log.SystemLogger.Error(err, "Error upserting module metadata, rolling back")
		tx.Rollback()
		return
	}

	err = tx.Commit()

	if err != nil {
		m.log.SystemLogger.Error(err, "Error committing upsert module metadata")
		return
	}

	moduleId, err := res.LastInsertId()

	if err != nil {
		m.log.SystemLogger.Error(err, "Error retrieving last insert ID for module metadata")
		return
	}

	m.epRepo.BatchInsertModuleEndpoints(item.ModuleEndpoints, moduleId)

	m.typeRepo.InsertElementTypes(item.AcceptedElementTypes, moduleId)

	return
}

func (m *ModuleMetadataRepo) GetAll() (receiver []structs.ModuleMetadata, err error) {
	err = m.db.Select(&receiver, fmt.Sprintf("SELECT * FROM %s ORDER BY created_at DESC;", ModuleTable))
	return
}

func (m *ModuleMetadataRepo) GetAllOfType(moduleType structs.ModuleType) (receiver []structs.ModuleMetadata, err error) {
	err = m.db.Select(&receiver, fmt.Sprintf("SELECT * FROM %s WHERE module_type = ? ORDER BY created_at DESC;", ModuleTable), moduleType)
	return
}

func (m *ModuleMetadataRepo) DeleteByServiceName(serviceName string) error {
	smt := fmt.Sprintf(`DELETE FROM %s WHERE module_service_name = ?`, ModuleTable)
	tx, err := m.db.Begin()
	if err != nil {
		m.log.SystemLogger.Error(err, "Error starting transaction to delete module_metadata")
		return err
	}
	_, err = tx.Exec(smt, serviceName)
	if err != nil {
		m.log.SystemLogger.Error(err, "Error deleting module_metadata, rolling back")
		err = tx.Rollback()
		return err
	}

	err = tx.Commit()

	if err != nil {
		m.log.SystemLogger.Error(err, "Error committing delete module_metadata")
		return err
	}

	return nil
}
