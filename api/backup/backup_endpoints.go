package backup

import (
	"encoding/json"
	"fp-dynamic-elements-manager-controller/api/util"
	"fp-dynamic-elements-manager-controller/internal/backup"
	"fp-dynamic-elements-manager-controller/internal/backup/structs"
	notificationfuncs "fp-dynamic-elements-manager-controller/internal/notification"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"net/http"
)

// Handler handles requests to the backup and restore provider, it accepts PostedCommand's
// Which contain the command to run (backup/restore) and if the command is restore,
// it expects a commit hash to restore to
func Handler(provider backup.Provider, ns notificationfuncs.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodOptions:
			w.Header().Set("Access-Control-Allow-Methods", "OPTIONS,GET,POST,PUT")
			w.WriteHeader(http.StatusNoContent)
		case http.MethodGet:
			history, err := provider.List()
			if err != nil {
				log.Error().Err(err).Msg("error retrieving git history")
				util.ReturnHTTPStatus(w, http.StatusInternalServerError, "error retrieving git history")
				return
			}
			if len(history) > 10 {
				history = history[:10]
			}
			sch := backup.Schedule{
				DayOfWeek: viper.GetString("dayofweek"),
				TimeOfDay: viper.GetString("timeofday"),
			}
			s := DetailResponse{
				History:  history,
				Schedule: sch,
			}
			json.NewEncoder(w).Encode(s)
		case http.MethodPost:
			handlePOST(r, w, provider, ns)
			util.ReturnHTTPStatus(w, http.StatusOK, "command executed successfully")
		case http.MethodPut:
			sched := backup.Schedule{}
			err := json.NewDecoder(r.Body).Decode(&sched)
			if err != nil {
				log.Error().Err(err).Msg("error decoding schedule wrapper")
				util.ReturnHTTPStatus(w, http.StatusNotAcceptable, "could not decode json into entity")
				return
			}
			go provider.StartAutoBackup(sched)
			util.ReturnHTTPStatus(w, http.StatusOK, "command executed successfully")
		}
		return
	})
}

func handlePOST(r *http.Request, w http.ResponseWriter, provider backup.Provider, ns notificationfuncs.Service) {
	item := PostedCommand{}
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		log.Error().Err(err).Msg("error decoding details wrapper")
		util.ReturnHTTPStatus(w, http.StatusNotAcceptable, "could not decode json into entity")
		return
	}
	switch item.Cmd {
	case backup.Backup:
		if err := provider.Backup("Manual"); err != nil {
			sendStatus(ns, notificationfuncs.Error, "Error running backup")
			return
		}
		sendStatus(ns, notificationfuncs.Success, "Backed up successfully")
	case backup.Restore:
		if err := provider.Restore(item.Hash); err != nil {
			sendStatus(ns, notificationfuncs.Error, "Error running restore")
			return
		}
		sendStatus(ns, notificationfuncs.Success, "Restored successfully")
	default:
		util.ReturnHTTPStatus(w, http.StatusBadRequest, "command not recognised")
		return
	}
}

type PostedCommand struct {
	Cmd  backup.Command `json:"cmd"`
	Hash string         `json:"hash"`
}

type DetailResponse struct {
	History  []structs.History `json:"history"`
	Schedule backup.Schedule   `json:"schedule"`
}

func sendStatus(ns notificationfuncs.Service, status notificationfuncs.EventType, msg string) {
	ns.Send(notificationfuncs.Event{
		EventType: status,
		Value:     msg,
	})
}
