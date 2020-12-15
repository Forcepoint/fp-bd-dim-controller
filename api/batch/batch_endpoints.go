package batch

import (
	"encoding/json"
	"fp-dynamic-elements-manager-controller/api/elements"
	"fp-dynamic-elements-manager-controller/api/util"
	"fp-dynamic-elements-manager-controller/internal/batch"
	"fp-dynamic-elements-manager-controller/internal/db/persistence"
	"fp-dynamic-elements-manager-controller/internal/queue/structs"
	"github.com/rs/zerolog/log"
	"net/http"
	"strconv"
	"strings"
)

func Handler(repo *persistence.UpdateStatusRepo) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodOptions:
			w.Header().Set("Access-Control-Allow-Methods", "OPTIONS,GET")
			w.WriteHeader(http.StatusNoContent)
		case http.MethodGet:
			query := r.URL.Query().Get("status")
			pageString := r.URL.Query().Get("page")
			if query == "" {
				util.ReturnHTTPStatus(w, http.StatusBadRequest, "status not specified")
				return
			}
			status := structs.Status(strings.ToLower(query))
			if pageString == "" {
				util.ReturnHTTPStatus(w, http.StatusBadRequest, "page number not specified")
				return
			}
			page, err := strconv.Atoi(pageString)
			if err != nil {
				log.Error().Err(err).Msg("error parsing string to int")
				util.ReturnHTTPStatus(w, http.StatusInternalServerError, "could not parse page number")
				return
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(batch.GetPaginatedBatchResults(page, elements.DefaultPageSize, status, repo))
		}
		return
	})
}
