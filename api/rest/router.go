// Package rest
package rest

import (
	"net/http"

	"horizonx-server/api/rest/handler"
	"horizonx-server/internal/config"
	"horizonx-server/internal/logger"
	"horizonx-server/internal/store"
	"horizonx-server/internal/transport/websocket"
)

func NewRouter(cfg *config.Config, store *store.SnapshotStore, hub *websocket.Hub, log logger.Logger) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/metrics", handler.HandleMetrics(store))
	mux.HandleFunc("/ws", handler.HandleWs(hub, log))

	// Placeholder for new feature routes
	// mux.HandleFunc("/ssh", handler.HandleSSH)
	// mux.HandleFunc("/deploy", handler.HandleDeploy)

	return applyCORS(cfg, mux)
}

func applyCORS(cfg *config.Config, next http.Handler) http.Handler {
	allowed := make(map[string]bool, len(cfg.AllowedOrigins))
	for _, origin := range cfg.AllowedOrigins {
		allowed[origin] = true
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		if allowed[origin] {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Methods", "GET,POST,OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")
		}

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// func applyCORS(cfg *config.Config, next http.Handler) http.Handler {
// 	allowed := make(map[string]bool, len(cfg.AllowedOrigins))
// 	for _, origin := range cfg.AllowedOrigins {
// 		allowed[origin] = true
// 	}
//
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		origin := r.Header.Get("Origin")
//
// 		if allowed[origin] {
// 			w.Header().Set("Access-Control-Allow-Origin", origin)
// 			w.Header().Set("Vary", "Origin")
// 			w.Header().Set("Access-Control-Allow-Credentials", "true")
// 			w.Header().Set("Access-Control-Allow-Methods", "GET,POST,OPTIONS")
// 			w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")
// 		}
//
// 		if r.Method == http.MethodOptions {
// 			w.WriteHeader(http.StatusOK)
// 			return
// 		}
//
// 		next.ServeHTTP(w, r)
// 	})
// }
