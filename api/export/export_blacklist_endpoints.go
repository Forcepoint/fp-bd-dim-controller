package export

import (
	"encoding/json"
	"fp-dynamic-elements-manager-controller/api/util"
	"fp-dynamic-elements-manager-controller/internal/db/persistence"
	"fp-dynamic-elements-manager-controller/internal/export/structs"
	"github.com/rs/zerolog/log"
	"net/http"
)

// Handler returns the entire ListElements table as a JSON array
// TODO update this to use pagination
func Handler(repo *persistence.ListElementRepo) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodOptions:
			w.Header().Set("Access-Control-Allow-Methods", "OPTIONS,GET")
			w.WriteHeader(http.StatusNoContent)
		case http.MethodGet:
			res, err := repo.GetAll()
			if err != nil {
				log.Error().Err(err).Msg("error retrieving list elements")
				util.ReturnHTTPStatus(w, http.StatusInternalServerError, "could not return requested resource")
				return
			}
			json.NewEncoder(w).Encode(structs.JsonExportResults{Results: res})
		}
		return
	})
}

func LookupHandler(repo *persistence.ListElementRepo) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodOptions:
			w.Header().Set("Access-Control-Allow-Methods", "OPTIONS,GET")
			w.WriteHeader(http.StatusNoContent)
		case http.MethodGet:
			values := r.URL.Query()
			searchTerm := values.Get("key")
			if searchTerm == "" {
				util.ReturnHTTPStatus(w, http.StatusFailedDependency, "no search term provided")
			}
			res, err := repo.GetAllEquals(searchTerm)
			if err != nil {
				log.Error().Err(err).Msg("error retrieving list elements")
				util.ReturnHTTPStatus(w, http.StatusInternalServerError, "could not return requested resource")
				return
			}
			json.NewEncoder(w).Encode(res)
		}
		return
	})
}
