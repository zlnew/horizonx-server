package rest

import (
	"net/http"

	"horizonx-server/api/rest/handler"
	"horizonx-server/internal/logger"
	"horizonx-server/internal/store"
	"horizonx-server/internal/transport/websocket"
)

func NewRouter(store *store.SnapshotStore, hub *websocket.Hub, log logger.Logger) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/metrics", handler.HandleMetrics(store))
	mux.HandleFunc("/ws", handler.HandleWs(hub, log))

	// Placeholder for new feature routes
	// mux.HandleFunc("/ssh", handler.HandleSSH)
	// mux.HandleFunc("/deploy", handler.HandleDeploy)

	return mux
}
