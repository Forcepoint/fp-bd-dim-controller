package export

import (
	"fp-dynamic-elements-manager-controller/internal/db/persistence"
	"fp-dynamic-elements-manager-controller/internal/queue/structs"
	"github.com/rs/zerolog/log"
	"math"
)

func BuildPagedResults(page, pageSize int, searchTerm string, safeList bool, dao *persistence.DataAccessObject) structs.PaginatedElements {
	paginatedResults := structs.PaginatedElements{}
	// Create the offset used to query the table for the correct results to match the page number
	offset := (page - 1) * pageSize

	// Query the table using the given parameters in the where clause
	var err error
	var totalCount int64
	if searchTerm == "" {
		paginatedResults.Elements, err = dao.ListElementRepo.GetAllPaginated(offset, pageSize, safeList)

		if err != nil {
			log.Error().Err(err).Msg("Error retrieving paged list elements")
			return paginatedResults
		}

		totalCount, err = dao.ListElementRepo.GetTotalCount(safeList)

		if err != nil {
			log.Error().Err(err).Msg("Error retrieving total count list elements")
			return paginatedResults
		}
	} else {
		paginatedResults.Elements, err = dao.ListElementRepo.GetAllLike(offset, pageSize, searchTerm, safeList)

		if err != nil {
			log.Error().Err(err).Msg("Error retrieving paged list elements with search term")
			return paginatedResults
		}

		totalCount, err = dao.ListElementRepo.GetTotalCountWhereLike(safeList, searchTerm)

		if err != nil {
			log.Error().Err(err).Msg("Error retrieving total count list elements with search term")
			return paginatedResults
		}
	}

	// Set the current page number and use some division to find the total page count
	paginatedResults.PageNumber = page
	paginatedResults.TotalPageCount = int(math.Round(float64(totalCount) / float64(pageSize)))

	return paginatedResults
}
