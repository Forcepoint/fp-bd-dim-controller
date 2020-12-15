package batch

import (
	"fp-dynamic-elements-manager-controller/internal/db/persistence"
	"fp-dynamic-elements-manager-controller/internal/queue/structs"
	"github.com/rs/zerolog/log"
	"math"
)

func GetPaginatedBatchResults(page, pageSize int, status structs.Status, repo *persistence.UpdateStatusRepo) (items structs.PaginatedStatus) {
	// Create the offset used to query the table for the correct results to match the page number
	offset := (page - 1) * pageSize

	// Query the table using the given parameters
	var err error
	switch status {
	case structs.INCOMPLETE:
		items.Statuses, err = repo.GetAllWithStatusPaginated(offset, pageSize, []string{structs.PENDING.String(), structs.FAILED.String()})
	case structs.PENDING:
		items.Statuses, err = repo.GetAllWithStatusPaginated(offset, pageSize, []string{structs.PENDING.String()})
	case structs.FAILED:
		items.Statuses, err = repo.GetAllWithStatusPaginated(offset, pageSize, []string{structs.FAILED.String()})
	case structs.SUCCESS:
		items.Statuses, err = repo.GetAllWithStatusPaginated(offset, pageSize, []string{structs.SUCCESS.String()})
	}

	if err != nil {
		log.Error().Err(err).Msg("Error retrieving paginated update statuses")
		return items
	}

	totalRows, err := repo.GetTotalCount()

	if err != nil {
		log.Error().Err(err).Msg("Error getting total count of update statuses")
		return items
	}

	// Set the current page number and use some division to find the total page count
	items.PageNumber = page
	items.TotalPageCount = int(math.Round(float64(totalRows) / float64(pageSize)))

	return
}
