package main

import (
	"context"
	"errors"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"horizonx-server/internal/config"
	"horizonx-server/internal/core"
	"horizonx-server/internal/core/auth"
	"horizonx-server/internal/core/metrics"
	"horizonx-server/internal/core/user"
	"horizonx-server/internal/domain"
	"horizonx-server/internal/logger"
	"horizonx-server/internal/storage/snapshot"
	"horizonx-server/internal/storage/sqlite"
	"horizonx-server/internal/transport/rest"
	"horizonx-server/internal/transport/websocket"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg := config.Load()
	log := logger.New(cfg)

	ms := snapshot.NewMetricsStore()
	hub := websocket.NewHub(log)
	sampler := metrics.NewSampler(log)

	sched := core.NewScheduler(cfg.Interval, log, sampler.Collect, func(m domain.Metrics) {
		ms.Set(m)
		hub.Emit("metrics", "metrics.updated", m)
	})
	go sched.Start(ctx)
	go hub.Run()

	db, err := sqlite.NewSqliteDB(cfg.DBPath, log)
	if err != nil {
		log.Error("sqlite", "connect", err)
		return
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Error("sqlite", "close", err)
		}
	}()

	userRepo := sqlite.NewUserRepository(db)
	userService := user.NewService(userRepo, cfg)
	authService := auth.NewService(userRepo, cfg)

	wsHandler := websocket.NewHandler(hub, cfg, log)
	metricsHandler := rest.NewMetricsHandler(ms)
	authHandler := rest.NewAuthHandler(authService, cfg)
	userHandler := rest.NewUserHandler(userService, cfg)

	router := rest.NewRouter(cfg, &rest.RouterDeps{
		WS:      wsHandler,
		Metrics: metricsHandler,
		Auth:    authHandler,
		User:    userHandler,
	})

	srv := rest.NewServer(router, cfg.Address)

	errCh := make(chan error, 1)
	go func() {
		log.Info("starting http server", "address", cfg.Address)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
		close(errCh)
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			log.Error("http server shutdown error", "error", err)
		}

	case err := <-errCh:
		log.Error("http server error", "error", err)
	}

	log.Info("server stopped")
}
