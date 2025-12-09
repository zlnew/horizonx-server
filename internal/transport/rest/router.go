// Package rest
package rest

import (
	"net/http"

	"horizonx-server/internal/config"
	"horizonx-server/internal/transport/rest/middleware"
	"horizonx-server/internal/transport/websocket"
)

type RouterDeps struct {
	WS      *websocket.Handler
	Metrics *MetricsHandler
	Auth    *AuthHandler
	User    *UserHandler
}

func NewRouter(cfg *config.Config, deps *RouterDeps) http.Handler {
	mux := http.NewServeMux()

	globalMw := middleware.New()
	globalMw.Use(middleware.CORS(cfg))
	globalMw.Use(middleware.CSRF(cfg))

	authMw := middleware.New()
	authMw.Use(middleware.JWT(cfg))

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	mux.HandleFunc("GET /ws", deps.WS.Serve)

	mux.Handle("GET /metrics", authMw.Then(http.HandlerFunc(deps.Metrics.Get)))

	mux.HandleFunc("POST /auth/register", deps.Auth.Register)
	mux.HandleFunc("POST /auth/login", deps.Auth.Login)
	mux.Handle("POST /auth/logout", authMw.Then(http.HandlerFunc(deps.Auth.Logout)))

	mux.Handle("GET /users", authMw.Then(http.HandlerFunc(deps.User.Index)))
	mux.Handle("POST /users", authMw.Then(http.HandlerFunc(deps.User.Store)))
	mux.Handle("PUT /users/{id}", authMw.Then(http.HandlerFunc(deps.User.Update)))
	mux.Handle("DELETE /users/{id}", authMw.Then(http.HandlerFunc(deps.User.Destroy)))

	// Placeholder for new feature routes
	// mux.HandleFunc("/ssh", handler.HandleSSH)
	// mux.HandleFunc("/deploy", handler.HandleDeploy)

	return globalMw.Apply(mux)
}
