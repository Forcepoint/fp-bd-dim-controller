package registration

import (
	"encoding/json"
	"fp-dynamic-elements-manager-controller/api/util"
	"fp-dynamic-elements-manager-controller/internal/modules/structs"
	"github.com/rs/zerolog/log"
	"net/http"
)

// Handler takes ModuleMetadata via a POST request from a child module to register a new module
// It passes the struct through a channel to a separately running goroutine that handles processing and adding that module
func Handler(addRoute chan<- structs.ModuleMetadata) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodOptions:
			w.Header().Set("Access-Control-Allow-Methods", "OPTIONS,GET,POST")
			w.WriteHeader(http.StatusNoContent)
		case http.MethodPost:
			metadata := structs.ModuleMetadata{}
			err := json.NewDecoder(r.Body).Decode(&metadata)
			if err != nil {
				log.Error().Err(err).Msg("error decoding json into entity")
				util.ReturnHTTPStatus(w, http.StatusNotAcceptable, "could not decode json into entity")
				return
			}
			addRoute <- metadata
			w.WriteHeader(http.StatusAccepted)
		}
		return
	})
}
