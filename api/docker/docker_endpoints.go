package docker

import (
	"encoding/json"
	"fp-dynamic-elements-manager-controller/api/util"
	"fp-dynamic-elements-manager-controller/internal/docker"
	"fp-dynamic-elements-manager-controller/internal/docker/structs"
	notifications "fp-dynamic-elements-manager-controller/internal/notification"
	"github.com/rs/zerolog/log"
	"net/http"
	"time"
)

func Handler(handler *docker.CommandHandler, ns notifications.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodOptions:
			w.Header().Set("Access-Control-Allow-Methods", "OPTIONS,GET,POST")
			w.WriteHeader(http.StatusNoContent)
		case http.MethodGet:
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(handler.MapContainerNames())
		case http.MethodPost:
			item := structs.ContainerDetailsWrapper{}
			err := json.NewDecoder(r.Body).Decode(&item)
			if err != nil {
				log.Error().Err(err).Msg("error decoding details wrapper")
				util.ReturnHTTPStatus(w, http.StatusNotAcceptable, "could not decode json into entity")
				return
			}
			go captureResultAsync(item, handler, ns)
			util.ReturnHTTPStatus(w, http.StatusOK, "commands added to queue")
		}
		return
	})
}

func captureResultAsync(item structs.ContainerDetailsWrapper, handler *docker.CommandHandler, ns notifications.Service) {
	doneCh, evtCh, errCh := handler.RunCommands(item)
	for {
		select {
		case err := <-errCh:
			if err != nil {
				log.Error().Err(err)
				ns.Send(notifications.Event{
					EventType: notifications.Error,
					Value:     err.Error(),
				})
			}
		case evt := <-evtCh:
			// send done over socket with a slight delay for modules which are slow to startup
			if evt.Context.State == notifications.Started || evt.Context.State == notifications.Created {
				time.Sleep(1 * time.Second)
			}
			ns.Send(evt)
		case <-doneCh:
			log.Info().Msg("commands completed successfully")
			return
		}
	}
}
