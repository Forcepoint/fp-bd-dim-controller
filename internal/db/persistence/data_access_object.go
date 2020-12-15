package persistence

import (
	"fp-dynamic-elements-manager-controller/internal/logging/structs"
	"github.com/jmoiron/sqlx"
)

type DataAccessObject struct {
	HealthRepo         *HealthRepo
	UserRepo           *UserRepo
	ListElementRepo    *ListElementRepo
	ModuleMetadataRepo *ModuleMetadataRepo
	ElementTypeRepo    *ElementTypeRepo
	LogEntryRepo       *LogEntryRepo
	ElementBatchRepo   *ElementBatchRepo
	UpdateStatusRepo   *UpdateStatusRepo
}

func NewDataAccessObject(appDb *sqlx.DB, logger *structs.AppLogger) *DataAccessObject {
	elementTypeRepo := NewElementTypeRepo(appDb, logger)
	return &DataAccessObject{
		HealthRepo:      NewHealthRepo(appDb),
		UserRepo:        NewUserRepo(appDb, logger),
		ListElementRepo: NewListElementRepo(appDb, logger),
		ModuleMetadataRepo: NewModuleMetadataRepo(appDb, logger,
			NewModuleEndpointRepo(appDb, logger), elementTypeRepo),
		ElementTypeRepo:  elementTypeRepo,
		LogEntryRepo:     NewLogEntryRepo(appDb, logger),
		ElementBatchRepo: NewElementBatchRepo(appDb, logger),
		UpdateStatusRepo: NewUpdateStatusRepo(appDb, logger),
	}
}
