package update

import (
	"encoding/json"
	"fp-dynamic-elements-manager-controller/api/util"
	"fp-dynamic-elements-manager-controller/internal/db/persistence"
	"fp-dynamic-elements-manager-controller/internal/queue/structs"
	"github.com/rs/zerolog/log"
	"net/http"
)

// Handler handles all status updates for pushes of batches to egress modules,
// the update statuses can be "success" or "failed"
func Handler(repo *persistence.UpdateStatusRepo) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodOptions:
			w.Header().Set("Access-Control-Allow-Methods", "OPTIONS,GET")
			w.WriteHeader(http.StatusNoContent)
		case http.MethodPost:
			item := structs.UpdateStatus{}
			err := json.NewDecoder(r.Body).Decode(&item)
			if err != nil {
				log.Error().Err(err).Msg("error decoding json into entity")
				util.ReturnHTTPStatus(w, http.StatusNotAcceptable, "could not decode json into entity")
				return
			}
			repo.UpdateUpdateStatus(item)
			w.WriteHeader(http.StatusAccepted)
			json.NewEncoder(w).Encode(item)
		}
		return
	})
}
