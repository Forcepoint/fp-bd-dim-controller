package stats

import (
	"encoding/json"
	"fp-dynamic-elements-manager-controller/internal/db/persistence"
	"net/http"
)

func Handler(repo *persistence.ListElementRepo) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodOptions:
			w.Header().Set("Access-Control-Allow-Methods", "OPTIONS,GET")
			w.WriteHeader(http.StatusNoContent)
		case http.MethodGet:
			serviceName := r.URL.Query().Get("servicename")
			json.NewEncoder(w).Encode(repo.GetStats(serviceName))
		}
		return
	})
}
