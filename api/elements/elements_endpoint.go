package elements

import (
	"encoding/json"
	"fp-dynamic-elements-manager-controller/api/util"
	"fp-dynamic-elements-manager-controller/internal/db/persistence"
	"fp-dynamic-elements-manager-controller/internal/export"
	"fp-dynamic-elements-manager-controller/internal/logging/structs"
	notificationfuncs "fp-dynamic-elements-manager-controller/internal/notification"
	"fp-dynamic-elements-manager-controller/internal/queue"
	structs2 "fp-dynamic-elements-manager-controller/internal/queue/structs"
	"github.com/rs/zerolog/log"
	"net/http"
	"strconv"
)

// DefaultPageSize defines the number of results to return to the client
const (
	DefaultPageSize = 20
)

// Handler handles all requests on the /elements route
func Handler(pusher queue.Pusher, dao *persistence.DataAccessObject, logger *structs.AppLogger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodOptions:
			w.Header().Set("Access-Control-Allow-Methods", "OPTIONS,GET,DELETE,PUT,POST")
			w.WriteHeader(http.StatusNoContent)
		case http.MethodGet:
			pageParam := r.URL.Query().Get("page")
			pageSizeParam := r.URL.Query().Get("pageSize")
			searchTerm := r.URL.Query().Get("searchterm")
			safeList := r.URL.Query().Get("safeList")
			if pageParam == "" {
				util.ReturnHTTPStatus(w, http.StatusBadRequest, "page number not specified")
				return
			}
			pageSize, err := strconv.Atoi(pageSizeParam)
			if err != nil {
				pageSize = DefaultPageSize
			}
			page, err := strconv.Atoi(pageParam)
			if err != nil {
				log.Error().Err(err).Msg("error parsing string to int")
				util.ReturnHTTPStatus(w, http.StatusInternalServerError, "could not parse page value")
				return
			}
			var safe = false
			if safeList != "" {
				safe, err = strconv.ParseBool(safeList)
			}
			json.NewEncoder(w).Encode(export.BuildPagedResults(page, pageSize, searchTerm, safe, dao))
		case http.MethodPost:
			item := structs2.ListElement{}
			err := json.NewDecoder(r.Body).Decode(&item)
			if err != nil {
				log.Error().Err(err).Msg("error decoding json entity")
				util.ReturnHTTPStatus(w, http.StatusNotAcceptable, "could not decode json into entity")
				return
			}
			err = queue.AddOne([]structs2.ListElement{item}, pusher, dao, logger)
			if err == persistence.ErrDuplicateValue {
				util.ReturnHTTPStatus(w, http.StatusConflict, "duplicate value")
				return
			}
			if err == queue.ErrInvalidFormat {
				util.ReturnHTTPStatus(w, http.StatusNotAcceptable, "invalid format")
				return
			}
			util.ReturnHTTPStatus(w, http.StatusOK, "success")
		case http.MethodPut:
			item := structs2.ListElement{}
			err := json.NewDecoder(r.Body).Decode(&item)
			if err != nil {
				log.Error().Err(err).Msg("error decoding json entity")
				util.ReturnHTTPStatus(w, http.StatusNotAcceptable, "could not decode json into entity")
				return
			}
			err = dao.ListElementRepo.UpdateListElement(item)
			if err == persistence.ErrDuplicateValue {
				logger.NotificationService.Send(notificationfuncs.Event{
					EventType: notificationfuncs.Error,
					Value:     "Cannot update, duplicate value",
				})
				util.ReturnHTTPStatus(w, http.StatusConflict, "duplicate value")
				return
			}
			if err != nil {
				util.ReturnHTTPStatus(w, http.StatusInternalServerError, "error updating value")
				return
			}
			util.ReturnHTTPStatus(w, http.StatusOK, "success")
		case http.MethodDelete:
			item := structs2.ListElement{}
			err := json.NewDecoder(r.Body).Decode(&item)
			if err != nil {
				log.Error().Err(err).Msg("error decoding json entity")
				util.ReturnHTTPStatus(w, http.StatusNotAcceptable, "could not decode json into entity")
				return
			}
			queue.Delete(item, pusher, dao)
			util.ReturnHTTPStatus(w, http.StatusOK, "success")
		}
		return
	})
}
