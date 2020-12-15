package modules

import (
	"encoding/json"
	"fp-dynamic-elements-manager-controller/api/util"
	"fp-dynamic-elements-manager-controller/internal/db/persistence"
	"fp-dynamic-elements-manager-controller/internal/modules"
	"fp-dynamic-elements-manager-controller/internal/modules/structs"
	"github.com/rs/zerolog/log"
	"net/http"
)

func Handler(dao *persistence.DataAccessObject) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodOptions:
			w.Header().Set("Access-Control-Allow-Methods", "OPTIONS,GET,DELETE")
			w.WriteHeader(http.StatusNoContent)
		case http.MethodGet:
			moduleType := r.URL.Query().Get("moduleType")
			mods, err := modules.GetModuleData(structs.ModuleType(moduleType), dao)
			if err != nil {
				log.Error().Err(err).Msg("error retrieving module data")
				util.ReturnHTTPStatus(w, http.StatusInternalServerError, "could not retrieve requested resource")
				return
			}
			json.NewEncoder(w).Encode(mods)
		}
		return
	})
}
