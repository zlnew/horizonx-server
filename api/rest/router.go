// Package rest
package rest

import (
	"net/http"

	"horizonx-server/api/rest/handler"
	"horizonx-server/api/rest/middleware"
	"horizonx-server/internal/config"
	"horizonx-server/internal/logger"
	"horizonx-server/internal/storage/snapshot"
	"horizonx-server/internal/transport/websocket"
)

func NewRouter(cfg *config.Config, ms *snapshot.MetricsStore, hub *websocket.Hub, log logger.Logger) http.Handler {
	mux := http.NewServeMux()

	wsHandler := websocket.NewHandler(hub, log, cfg)

	mux.HandleFunc("/ws", wsHandler.Serve)
	mux.HandleFunc("/metrics", handler.HandleMetrics(ms))

	// Placeholder for new feature routes
	// mux.HandleFunc("/ssh", handler.HandleSSH)
	// mux.HandleFunc("/deploy", handler.HandleDeploy)

	mw := middleware.New()
	mw.Use(middleware.CORS(cfg))
	mw.Use(middleware.JWT(cfg))

	return mw.Apply(mux)
}
