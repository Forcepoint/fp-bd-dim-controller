package queue

import (
	"bytes"
	"encoding/json"
	"fmt"
	"fp-dynamic-elements-manager-controller/internal/db/persistence"
	"fp-dynamic-elements-manager-controller/internal/health"
	structs3 "fp-dynamic-elements-manager-controller/internal/logging/structs"
	structs2 "fp-dynamic-elements-manager-controller/internal/modules/structs"
	"fp-dynamic-elements-manager-controller/internal/queue/structs"
	"github.com/pkg/errors"
	"net/http"
	"time"
)

type Pusher interface {
	push()
}

type DataPusher struct {
	dao    *persistence.DataAccessObject
	logger *structs3.AppLogger
}

func NewDataPusher(dao *persistence.DataAccessObject, logger *structs3.AppLogger) *DataPusher {
	return &DataPusher{
		dao:    dao,
		logger: logger,
	}
}

func (t *DataPusher) push() {
	// Enter function, get a list of all modules capable of consuming intelligence
	modules, err := egressModules(t.dao.ModuleMetadataRepo)

	if err != nil {
		t.logger.SystemLogger.Error(err, "error getting egress modules")
		return
	}
	// Iterate over the modules and check if they are up and configured, if not skip them to avoid congesting the network unnecessarily
	for _, module := range modules {

		if !module.Configured {
			continue
		}

		// If the module is not up and healthy, skip it and move on to the next one in the list
		if !health.IsUp(module.ModuleServiceName, module.InternalPort) {
			t.logger.UserLogger.Debug(fmt.Sprintf("%s is not up, cannot push", module.ModuleServiceName))
			continue
		}

		// Get a list of the accepted element types for the specific module (IP, Domain, Range, etc.)
		acceptedTypes, err := t.dao.ElementTypeRepo.GetAllForModule(module.ID)

		if err != nil {
			t.logger.SystemLogger.Error(err, "error retrieving types for module")
			return
		}

		// Get batch IDs for the blocklist items that have not been pushed to this module before
		blocklistBatchIds, err := unpushedBatchIds(module.ID, false, acceptedTypes, t.dao.ListElementRepo)

		if err != nil {
			t.logger.SystemLogger.Error(err, "error retrieving batch ids for module")
			return
		}

		// Iterate through the batch IDs and push them to the module
		for _, val := range blocklistBatchIds {
			err := queryBatchAndPush(
				val,
				module,
				acceptedTypes,
				false,
				t.dao.ListElementRepo,
				t.dao.UpdateStatusRepo,
				t.logger,
				false)

			if err != nil {
				t.logger.SystemLogger.Error(err, fmt.Sprintf("Safelist: %v pushing to module ID: %d", false, val))
			}
		}

		// Get batch IDs for the safelist items that have not been pushed to this module before
		safelistBatchIds, err := unpushedBatchIds(module.ID, true, acceptedTypes, t.dao.ListElementRepo)
		if err != nil {
			t.logger.SystemLogger.Error(err, "error retrieving batch ids for module")
			return
		}

		// Iterate through the batch IDs and push them to the module
		for _, val := range safelistBatchIds {
			err := queryBatchAndPush(
				val,
				module,
				acceptedTypes,
				true,
				t.dao.ListElementRepo,
				t.dao.UpdateStatusRepo,
				t.logger,
				false)

			if err != nil {
				t.logger.SystemLogger.Error(err, fmt.Sprintf("Safelist: %v pushing to module ID: %d", true, val))
			}
		}

		go func(moduleData structs2.ModuleMetadata, types []structs.ElementType) {
			time.Sleep(60 * time.Second)
			// Check if there are any failed batches for the current module and if so, query them and push them
			err = pushFailedBatches(t.dao, moduleData, types, t.logger)
			if err != nil {
				t.logger.SystemLogger.Error(err, "error pushing failed batches")
			}
			return
		}(module, acceptedTypes)
	}
}

func pushFailedBatches(dao *persistence.DataAccessObject, module structs2.ModuleMetadata, types []structs.ElementType, logger *structs3.AppLogger) error {
	// Get batch IDs for the blocklist items that have failed for this module before (status in the table of FAILED)
	blocklistFailedBatches, err := failedBatchIds(dao.UpdateStatusRepo, module.ID, false)

	if err != nil {
		return err
	}

	if len(blocklistFailedBatches) == 0 {
		return nil
	}

	logger.UserLogger.Info(fmt.Sprintf("Pushing failed batches for %s", module.ModuleServiceName))

	for _, val := range blocklistFailedBatches {
		err := queryBatchAndPush(
			val,
			module,
			types,
			false,
			dao.ListElementRepo,
			dao.UpdateStatusRepo,
			logger,
			true)

		if err != nil {
			logger.SystemLogger.Error(err, fmt.Sprintf("Safelist: %v Failed: %v pushing to module ID: %d", false, true, val))
		}
	}

	// Get batch IDs for the safelist items that have failed for this module before (status in the table of FAILED)
	safelistFailedBatches, err := failedBatchIds(dao.UpdateStatusRepo, module.ID, true)

	if err != nil {
		return err
	}

	if len(safelistFailedBatches) == 0 {
		return nil
	}

	logger.UserLogger.Info(fmt.Sprintf("Pushing failed batches for %s", module.ModuleServiceName))

	for _, val := range safelistFailedBatches {
		err := queryBatchAndPush(
			val,
			module,
			types,
			true,
			dao.ListElementRepo,
			dao.UpdateStatusRepo,
			logger,
			true)

		if err != nil {
			logger.SystemLogger.Error(err, fmt.Sprintf("Safelist: %v Failed: %v pushing to module ID: %d", true, true, val))
		}
	}

	return nil
}

func unpushedBatchIds(moduleId int64, safe bool, types []structs.ElementType, repo *persistence.ListElementRepo) ([]int64, error) {
	return repo.GetUnpushedBatchIds(moduleId, safe, types)
}

func failedBatchIds(repo *persistence.UpdateStatusRepo, moduleId int64, safe bool) (batchIds []int64, err error) {
	batchIds, err = repo.GetAllWithStatusForModule(structs.FAILED, moduleId, safe)
	return
}

func nextBatch(batchId int64, types []structs.ElementType, repo *persistence.ListElementRepo) (receiver []structs.ListElement, err error) {
	receiver, err = repo.GetAllByBatchId(batchId, types)
	return
}

func egressModules(repo *persistence.ModuleMetadataRepo) (receiver []structs2.ModuleMetadata, err error) {
	receiver, err = repo.GetAllOfType(structs2.EGRESS)
	return
}

func queryBatchAndPush(batchId int64, module structs2.ModuleMetadata,
	types []structs.ElementType, safe bool, listElementRepo *persistence.ListElementRepo,
	updateStatusRepo *persistence.UpdateStatusRepo, logger *structs3.AppLogger, failed bool) error {

	// Get the next batch for a module using a provided batch ID
	updateBatch, err := nextBatch(batchId, types, listElementRepo)

	if err != nil {
		return errors.Wrap(err, "Error retrieving next batch for pushing")
	}

	if len(updateBatch) == 0 {
		return nil
	}

	logger.UserLogger.Info(fmt.Sprintf("Pushing to %s", module.ModuleServiceName))

	wrappedBatch := structs.ProcessedItems{SafeList: safe, Items: updateBatch, BatchId: batchId}

	resp, err := pushData(module.ModuleServiceName, module.InternalPort, wrappedBatch, logger)

	if err != nil {
		if !failed {
			updateStatusRepo.InsertUpdateStatus(structs.UpdateStatus{
				ServiceName:      module.ModuleServiceName,
				Status:           structs.FAILED,
				UpdateBatchId:    batchId,
				ModuleMetadataId: module.ID,
			})
		}
		return errors.Wrap(err, "Http: Error pushing next batch")
	}

	logger.UserLogger.Info(fmt.Sprintf("Pushed batch %d to %s", batchId, module.ModuleServiceName))
	if resp.StatusCode != http.StatusAccepted {
		logger.UserLogger.Info(fmt.Sprintf("Pushing failed to %s", module.ModuleServiceName))
		if !failed {
			updateStatusRepo.InsertUpdateStatus(structs.UpdateStatus{
				ServiceName:      module.ModuleServiceName,
				Status:           structs.FAILED,
				UpdateBatchId:    batchId,
				ModuleMetadataId: module.ID,
			})
		}
	} else {
		logger.UserLogger.Info(fmt.Sprintf("Pushing succeeded to %s", module.ModuleServiceName))
		if !failed {
			updateStatusRepo.InsertUpdateStatus(structs.UpdateStatus{
				ServiceName:      module.ModuleServiceName,
				Status:           structs.PENDING,
				UpdateBatchId:    batchId,
				ModuleMetadataId: module.ID,
			})
		} else {
			updateStatusRepo.UpdateUpdateStatus(structs.UpdateStatus{
				ServiceName:      module.ModuleServiceName,
				Status:           structs.PENDING,
				UpdateBatchId:    batchId,
				ModuleMetadataId: module.ID,
			})
		}
	}
	return nil
}

func pushData(moduleName, modulePort, data interface{}, logger *structs3.AppLogger) (*http.Response, error) {
	jsonData, err := json.Marshal(data)

	if err != nil {
		logger.SystemLogger.Error(err, "Error marshalling batch into JSON")
		return nil, err
	}

	registerUrl := fmt.Sprintf("http://%s:%s/run", moduleName, modulePort)

	req, err := http.NewRequest("POST", registerUrl, bytes.NewBuffer(jsonData))

	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Timeout: 120 * time.Second,
	}
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	return resp, err
}
