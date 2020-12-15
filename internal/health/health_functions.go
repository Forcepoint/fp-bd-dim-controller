package health

import (
	"database/sql"
	"fmt"
	"fp-dynamic-elements-manager-controller/internal/db/persistence"
	"fp-dynamic-elements-manager-controller/internal/health/structs"
	structs2 "fp-dynamic-elements-manager-controller/internal/modules/structs"
	"fp-dynamic-elements-manager-controller/internal/util"
	"github.com/rs/zerolog/log"
	"net/http"
	"sync"
	"time"
)

func GetControllerHealth(dao *persistence.DataAccessObject) structs.Health {
	controller := structs.ModuleHealth{
		ModuleName: "master-controller",
		Status:     structs.Healthy,
		StatusCode: http.StatusOK,
	}

	database := structs.ModuleHealth{
		ModuleName: "master-database",
		StatusCode: http.StatusOK,
	}

	if up := dao.HealthRepo.PingDB(); up {
		database.Status = structs.Healthy
	} else {
		database.Status = structs.Down
	}

	return structs.Health{
		Modules: []structs.ModuleHealth{controller, database},
	}
}

func ModuleHealthCheck(modules []structs2.ModuleMetadata, dao *persistence.DataAccessObject) {
	wg := sync.WaitGroup{}
	wg.Add(len(modules))
	for i := range modules {
		go func(wait *sync.WaitGroup, module *structs2.ModuleMetadata) {
			defer wg.Done()
			getModuleHealth(module, dao)
		}(&wg, &(modules)[i])
	}
	wg.Wait()
}

func getModuleHealth(module *structs2.ModuleMetadata, dao *persistence.DataAccessObject) {
	resp, err := util.DimHTTPClient.Get(fmt.Sprintf("http://%s:%s/health", module.ModuleServiceName, module.InternalPort))

	if err != nil {
		log.Error().Err(err).Msg("Error retrieving module http health")
		module.ModuleHealth.Status = structs.Down
		module.ModuleHealth.StatusCode = 0
		return
	}

	module.ModuleHealth.ModuleName = module.ModuleDisplayName
	module.ModuleHealth.StatusCode = resp.StatusCode

	switch resp.StatusCode {
	case http.StatusOK, http.StatusTeapot:
		module.ModuleHealth.Status = structs.Healthy
		module.LastPing = time.Now()
		module.ModuleEndpoints = nil
		go dao.ModuleMetadataRepo.UpsertModuleMetadata(*module)
	case http.StatusNotImplemented:
		module.ModuleHealth.Status = structs.Unhealthy
	default:
		module.ModuleHealth.Status = structs.Down
	}

	switch module.ModuleType {
	case structs2.INGRESS:
		item, err := dao.ListElementRepo.GetLatestUpdate(module.ModuleServiceName)
		if err == sql.ErrNoRows {
			return
		}
		if err != nil {
			log.Error().Err(err).Msg("Error retrieving latest update ingress module health")
			return
		}
		if !item.CreatedAt.IsZero() {
			module.ModuleHealth.LastUpdate = item.CreatedAt.String()
		}
	case structs2.EGRESS:
		updateStatus, err := dao.UpdateStatusRepo.GetLatestUpdate(module.ID)
		if err == sql.ErrNoRows {
			return
		}
		if err != nil {
			log.Error().Err(err).Msg("Error retrieving latest update egress module health")
			return
		}
		if !updateStatus.CreatedAt.IsZero() {
			module.ModuleHealth.LastUpdate = updateStatus.CreatedAt.String()
		}
	}
}

func IsUp(svcName, port string) bool {
	resp, err := util.DimHTTPClient.Get(fmt.Sprintf("http://%s:%s/health", svcName, port))
	if err != nil {
		log.Error().Err(err).Msg("Error retrieving module isUp")
		return false
	}
	if resp.StatusCode == http.StatusOK {
		return true
	}
	return false
}
