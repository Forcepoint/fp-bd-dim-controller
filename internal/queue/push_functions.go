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
	"net/http"
)

type Pusher interface {
	pushToModules()
}

type TableDataPusher struct {
	dao    *persistence.DataAccessObject
	logger *structs3.AppLogger
}

func NewTableDataPusher(dao *persistence.DataAccessObject, logger *structs3.AppLogger) *TableDataPusher {
	return &TableDataPusher{
		dao:    dao,
		logger: logger,
	}
}

func (t TableDataPusher) pushToModules() {
	// Enter function, get a list of all modules capable of consuming intelligence
	res, err := getAllEgressModules(t.dao)

	if err != nil {
		t.logger.SystemLogger.Error(err, "error getting egress modules")
		return
	}
	// Iterate over the modules and check if they are up and configured, if not skip them to avoid congesting the network unnecessarily
	for _, module := range res {
		if !module.Configured {
			continue
		}
		if !health.IsUp(module.ModuleServiceName, module.InternalPort) {
			t.logger.UserLogger.Debug(fmt.Sprintf("%s is not up, cannot push", module.ModuleServiceName))
			continue
		}

		acceptedTypes, err := t.dao.ElementTypeRepo.GetAllForModule(module.ID)

		if err != nil {
			t.logger.SystemLogger.Error(err, "error retrieving types for module")
			return
		}

		batchIds, err := getAllUnpushedBatchIds(module.ID, acceptedTypes, t.dao)

		if err != nil {
			t.logger.SystemLogger.Error(err, "error retrieving batch ids for module")
			return
		}

		// Iterate through the batch IDs and push them to the module
		for _, val := range batchIds {
			queryBatchAndPush(val, module.ID, module.ModuleServiceName, module.InternalPort, acceptedTypes, false, t.dao, t.logger, false)
			queryBatchAndPush(val, module.ID, module.ModuleServiceName, module.InternalPort, acceptedTypes, true, t.dao, t.logger, false)
		}

		err = pushFailedBatches(t.dao, module, acceptedTypes, t.logger)

		if err != nil {
			t.logger.SystemLogger.Error(err, "error pushing failed batches")
			continue
		}
	}
}

func pushFailedBatches(dao *persistence.DataAccessObject, module structs2.ModuleMetadata, types []structs.ElementType, logger *structs3.AppLogger) error {
	failedBatches, err := getFailedBatchIds(dao, module.ID)

	if err != nil {
		return err
	}

	if len(failedBatches) == 0 {
		return nil
	}

	logger.UserLogger.Info(fmt.Sprintf("Pushing failed batches for %s", module.ModuleServiceName))

	for _, val := range failedBatches {
		queryBatchAndPush(val, module.ID, module.ModuleServiceName, module.InternalPort, types, false, dao, logger, true)
		queryBatchAndPush(val, module.ID, module.ModuleServiceName, module.InternalPort, types, true, dao, logger, true)
	}

	return nil
}

func getAllUnpushedBatchIds(moduleId int64, types []structs.ElementType, dao *persistence.DataAccessObject) ([]int64, error) {
	return dao.ListElementRepo.GetUnpushedBatchIds(moduleId, types)
}

func getFailedBatchIds(dao *persistence.DataAccessObject, moduleId int64) (batchIds []int64, err error) {
	batchIds, err = dao.UpdateStatusRepo.GetAllWithStatusForModule(structs.FAILED, moduleId)
	return
}

func getNextBatch(batchId int64, types []structs.ElementType, safe bool, dao *persistence.DataAccessObject) (receiver []structs.ListElement, err error) {
	receiver, err = dao.ListElementRepo.GetAllByBatchId(batchId, types, safe)
	return
}

func getAllEgressModules(dao *persistence.DataAccessObject) (receiver []structs2.ModuleMetadata, err error) {
	receiver, err = dao.ModuleMetadataRepo.GetAllOfType(structs2.EGRESS)
	return
}

func queryBatchAndPush(batchId int64, moduleId int64, svcName, svcPort string, types []structs.ElementType, safe bool, dao *persistence.DataAccessObject, logger *structs3.AppLogger, failed bool) {
	// Get the next batch for a module using a provided batch ID
	updateBatch, err := getNextBatch(batchId, types, safe, dao)

	if err != nil {
		logger.SystemLogger.Error(err, "Error retrieving next batch for pushing")
	}

	if len(updateBatch) == 0 {
		return
	}

	logger.UserLogger.Info(fmt.Sprintf("Pushing to %s", svcName))

	wrappedBatch := structs.ProcessedItems{SafeList: safe, Items: updateBatch, BatchId: batchId}

	resp, err := pushData(svcName, svcPort, wrappedBatch, logger)

	if err != nil {
		logger.SystemLogger.Error(err, "Http: Error pushing next batch")
		if !failed {
			dao.UpdateStatusRepo.InsertUpdateStatus(structs.UpdateStatus{
				ServiceName:      svcName,
				Status:           structs.FAILED,
				UpdateBatchId:    batchId,
				ModuleMetadataId: moduleId,
			})
		}
		return
	}

	logger.UserLogger.Info(fmt.Sprintf("Pushed batch %d to %s", batchId, svcName))
	if resp.StatusCode != http.StatusAccepted {
		logger.UserLogger.Info(fmt.Sprintf("Pushing failed to %s", svcName))
		if !failed {
			dao.UpdateStatusRepo.InsertUpdateStatus(structs.UpdateStatus{
				ServiceName:      svcName,
				Status:           structs.FAILED,
				UpdateBatchId:    batchId,
				ModuleMetadataId: moduleId,
			})
		}
	} else {
		logger.UserLogger.Info(fmt.Sprintf("Pushing succeeded to %s", svcName))
		if !failed {
			dao.UpdateStatusRepo.InsertUpdateStatus(structs.UpdateStatus{
				ServiceName:      svcName,
				Status:           structs.PENDING,
				UpdateBatchId:    batchId,
				ModuleMetadataId: moduleId,
			})
		} else {
			dao.UpdateStatusRepo.UpdateUpdateStatus(structs.UpdateStatus{
				ServiceName:      svcName,
				Status:           structs.PENDING,
				UpdateBatchId:    batchId,
				ModuleMetadataId: moduleId,
			})
		}
	}
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

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	return resp, err
}
