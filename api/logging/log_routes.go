package logging

import (
	"encoding/json"
	"fp-dynamic-elements-manager-controller/api/util"
	"fp-dynamic-elements-manager-controller/internal/db/persistence"
	"fp-dynamic-elements-manager-controller/internal/logging"
	"fp-dynamic-elements-manager-controller/internal/logging/structs"
	"github.com/rs/zerolog/log"
	"net/http"
	"strconv"
)

// PageSize defines the number of results to return to the client
var pageSize = 10

// Handler returns the logs from the database based on a pagination system.
// It takes 3 query parameters, Page (the offset used to query the database), Level (the log level filter), and Module Name.
func Handler(repo *persistence.LogEntryRepo) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodOptions:
			w.Header().Set("Access-Control-Allow-Methods", "OPTIONS,GET,POST")
			w.WriteHeader(http.StatusNoContent)
		case http.MethodGet:
			param1 := r.URL.Query().Get("page")
			logLevel := r.URL.Query().Get("level")
			moduleName := r.URL.Query().Get("modulename")
			if param1 == "" {
				util.ReturnHTTPStatus(w, http.StatusBadRequest, "page number not specified")
				return
			}
			page, err := strconv.Atoi(param1)
			if err != nil {
				log.Error().Err(err).Msg("error parsing string to int")
				util.ReturnHTTPStatus(w, http.StatusInternalServerError, "could not parse page number")
				return
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(logging.BuildLogResults(page, pageSize, moduleName, logLevel, repo))
		case http.MethodPost:
			item := structs.LogEntry{}
			err := json.NewDecoder(r.Body).Decode(&item)
			if err != nil {
				log.Error().Err(err).Msg("error decoding json into entity")
				util.ReturnHTTPStatus(w, http.StatusNotAcceptable, "could not decode json into entity")
				return
			}
			go repo.InsertLogEntry(&item)
		}
		return
	})
}
