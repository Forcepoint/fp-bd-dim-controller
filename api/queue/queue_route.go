package queue

import (
	"encoding/json"
	"fmt"
	"fp-dynamic-elements-manager-controller/api/util"
	"fp-dynamic-elements-manager-controller/internal/db/persistence"
	structs2 "fp-dynamic-elements-manager-controller/internal/logging/structs"
	"fp-dynamic-elements-manager-controller/internal/queue"
	"fp-dynamic-elements-manager-controller/internal/queue/structs"
	"github.com/rs/zerolog/log"
	"net/http"
)

// Handler handles all requests to add elements to the queue from external modules to be written to the DB
func Handler(pusher queue.Pusher, dao *persistence.DataAccessObject, logger *structs2.AppLogger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodOptions:
			w.Header().Set("Access-Control-Allow-Methods", "OPTIONS,GET")
			w.WriteHeader(http.StatusNoContent)
		case http.MethodPost:
			items := structs.ProcessedItems{}
			err := json.NewDecoder(r.Body).Decode(&items)
			if err != nil {
				log.Error().Err(err).Msg("error decoding json into entity")
				util.ReturnHTTPStatus(w, http.StatusNotAcceptable, "could not decode json into entity")
				return
			}
			go queue.AddToQueue(items.Items, pusher, dao, logger)
			util.ReturnHTTPStatus(w, http.StatusAccepted, fmt.Sprintf("Success: %d items uploaded", len(items.Items)))
		}
		return
	})
}
