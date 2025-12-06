package handler

import (
	"encoding/json"
	"net/http"

	"horizonx-server/internal/logger"
	"horizonx-server/internal/store"
	"horizonx-server/internal/transport/websocket"
)

func HandleMetrics(store *store.SnapshotStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := store.Get()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(data)
	}
}

func HandleWs(hub *websocket.Hub, log logger.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		websocket.ServeWs(hub, w, r, log)
	}
}
