package logging

import (
	"fp-dynamic-elements-manager-controller/internal/db/persistence"
	structs2 "fp-dynamic-elements-manager-controller/internal/logging/structs"
	"github.com/rs/zerolog/log"
	"github.com/sirupsen/logrus"
	"math"
	"strings"
)

func BuildLogResults(page, pageSize int, moduleName, level string, repo *persistence.LogEntryRepo) structs2.LogEvents {
	logs := structs2.LogEvents{
		Events:         []structs2.LogEntry{},
		TotalPageCount: 0,
		PageNumber:     0,
	}

	// Create the offset used to query the table for the correct results to match the page number
	offset := (page - 1) * pageSize

	// Query the table using the given parameters in the where clause
	var err error
	logs.Events, err = repo.GetAllForLevels(moduleName, pageSize, offset, buildSearchLevels(strings.ToLower(level)))

	if err != nil {
		log.Error().Err(err).Msg("Error retrieving paginated log entries")
		return logs
	}

	var totalRows int64
	totalRows, err = repo.GetTotalCount()

	if err != nil {
		log.Error().Err(err).Msg("Error retrieving total count log entries")
		return logs
	}

	// Set the current page number and use some division to find the total page count
	logs.PageNumber = page
	logs.TotalPageCount = int(math.Round(float64(totalRows) / float64(pageSize)))

	return logs
}

func buildSearchLevels(levelParam string) []string {
	var levelsToSearch []string
	var addLevel = true
	for _, level := range logrus.AllLevels {
		if addLevel {
			levelsToSearch = append(levelsToSearch, level.String())
		}
		if levelParam == level.String() {
			addLevel = false
		}
	}

	return levelsToSearch
}
