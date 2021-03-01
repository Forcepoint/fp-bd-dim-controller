package queue

import (
	"errors"
	"fp-dynamic-elements-manager-controller/internal/db/persistence"
	structs2 "fp-dynamic-elements-manager-controller/internal/logging/structs"
	notificationfuncs "fp-dynamic-elements-manager-controller/internal/notification"
	"fp-dynamic-elements-manager-controller/internal/queue/structs"
	validation "fp-dynamic-elements-manager-controller/internal/util"
	"github.com/thoas/go-funk"
)

const MaxBatchSize = 5000

var ErrEmptySlice = errors.New("empty slice")
var ErrInvalidFormat = errors.New("invalid format")

func AddOne(items []structs.ListElement, pusher Pusher, dao *persistence.DataAccessObject, logger *structs2.AppLogger) error {
	if len(items) == 0 {
		return ErrEmptySlice
	}
	element := items[0]
	err := validateInput(element, logger)
	if err != nil {
		return err
	}
	res, err := dao.ElementBatchRepo.InsertBatchElement()
	if err != nil {
		logger.SystemLogger.Error(err, "Error inserting batch element in queue")
		return err
	}
	batchId, err := res.LastInsertId()
	if err != nil {
		logger.SystemLogger.Error(err, "Error retrieving last insert ID in queue")
		return err
	}
	for i := range items {
		items[i].UpdateBatchId = batchId
	}
	err = dao.ListElementRepo.InsertListElement(items[0])
	go func() {
		pusher.push()
	}()
	return err
}

func AddToQueue(items []structs.ListElement, pusher Pusher, dao *persistence.DataAccessObject, logger *structs2.AppLogger) {
	if len(items) == 0 {
		return
	}
	chunkedItems := funk.Chunk(items, MaxBatchSize)
	for _, chunk := range chunkedItems.([][]structs.ListElement) {
		res, err := dao.ElementBatchRepo.InsertBatchElement()
		if err != nil {
			logger.SystemLogger.Error(err, "Error inserting batch element in queue")
			return
		}
		batchId, err := res.LastInsertId()
		if err != nil {
			logger.SystemLogger.Error(err, "Error retrieving last insert ID in queue")
			return
		}

		for i := range chunk {
			chunk[i].UpdateBatchId = batchId
		}

		dao.ListElementRepo.BatchInsertListElements(chunk)
	}
	go func() {
		pusher.push()
	}()
}

func validateInput(element structs.ListElement, logger *structs2.AppLogger) error {
	switch element.Type {
	case structs.IP:
		if !validation.IsIpValid(element.Value) {
			logger.NotificationService.Send(notificationfuncs.Event{
				EventType: notificationfuncs.Error,
				Value:     "Invalid IP Format",
			})
			return ErrInvalidFormat
		}
	case structs.URL:
		if !validation.IsUrlValid(element.Value) {
			logger.NotificationService.Send(notificationfuncs.Event{
				EventType: notificationfuncs.Error,
				Value:     "Invalid URL Format",
			})
			return ErrInvalidFormat
		}
	case structs.DOMAIN:
		if !validation.IsDomainValid(element.Value) {
			logger.NotificationService.Send(notificationfuncs.Event{
				EventType: notificationfuncs.Error,
				Value:     "Invalid Domain Format",
			})
			return ErrInvalidFormat
		}
	case structs.RANGE:
		if !validation.IsRangeValid(element.Value) {
			logger.NotificationService.Send(notificationfuncs.Event{
				EventType: notificationfuncs.Error,
				Value:     "Invalid Range Format",
			})
			return ErrInvalidFormat
		}
	}
	return nil
}
