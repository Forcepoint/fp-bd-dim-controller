package notification

import (
	"fp-dynamic-elements-manager-controller/internal/notification"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
	"net/http"
)

// Handler upgrades the http(s) connection to a WebSocket connection and adds it to the Notification Service
func Handler(upgrader websocket.Upgrader, ns notification.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Error().Err(err).Msg("error upgrading connection")
			return
		}
		notification.AddClient(ns.Hub(), conn, make(chan notification.Event, 10))
		return
	})
}
