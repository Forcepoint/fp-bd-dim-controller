package modules

import (
	"fp-dynamic-elements-manager-controller/internal/db/persistence"
	"fp-dynamic-elements-manager-controller/internal/health"
	"fp-dynamic-elements-manager-controller/internal/modules/structs"
	"github.com/rs/zerolog/log"
)

func GetModuleData(moduleType structs.ModuleType, dao *persistence.DataAccessObject) (receiver []structs.ModuleMetadata, err error) {
	if moduleType == "" {
		receiver, err = dao.ModuleMetadataRepo.GetAllModuleMetadata()
	} else {
		receiver, err = dao.ModuleMetadataRepo.GetAllOfType(moduleType)
	}

	if err != nil {
		log.Error().Err(err)
		return nil, err
	}

	health.ModuleHealthCheck(receiver, dao)

	return receiver, nil
}
