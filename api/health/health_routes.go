package health

import (
	"encoding/json"
	"fp-dynamic-elements-manager-controller/internal/db/persistence"
	"fp-dynamic-elements-manager-controller/internal/health"
	"net/http"
)

// Handler returns the health status of the controller and the connected database
// There are 3 possible states, Healthy, Unhealthy and Down
func Handler(dao *persistence.DataAccessObject) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodOptions:
			w.Header().Set("Access-Control-Allow-Methods", "OPTIONS,GET")
			w.WriteHeader(http.StatusNoContent)
		case http.MethodGet:
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(health.GetControllerHealth(dao))
		}
		return
	})
}
