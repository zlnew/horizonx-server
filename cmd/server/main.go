package main

import (
	"context"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"horizonx-server/internal/config"
	"horizonx-server/internal/core/auth"
	"horizonx-server/internal/core/job"
	"horizonx-server/internal/core/server"
	"horizonx-server/internal/core/user"
	"horizonx-server/internal/domain"
	"horizonx-server/internal/logger"
	"horizonx-server/internal/storage/postgres"
	"horizonx-server/internal/transport/rest"
	"horizonx-server/internal/transport/ws"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg := config.Load()
	log := logger.New(cfg)

	if cfg.JWTSecret == "" {
		panic("FATAL: JWT_SECRET is mandatory for Server!")
	}

	dbPool, err := postgres.InitDB(cfg.DatabaseURL, log)
	if err != nil {
		log.Error("failed to init DB", "error", err)
	}
	defer dbPool.Close()

	serverRepo := postgres.NewServerRepository(dbPool)
	userRepo := postgres.NewUserRepository(dbPool)
	jobRepo := postgres.NewJobRepository(dbPool)

	serverService := server.NewService(serverRepo)
	authService := auth.NewService(userRepo, cfg.JWTSecret, cfg.JWTExpiry)
	userService := user.NewService(userRepo)

	serverHandler := rest.NewServerHandler(serverService)
	authHandler := rest.NewAuthHandler(authService, cfg)
	userHandler := rest.NewUserHandler(userService)

	hub := ws.NewHub(ctx, log)
	wsHandler := ws.NewHandler(hub, log, cfg.JWTSecret, cfg.AllowedOrigins)

	agentHub := ws.NewAgentHub(ctx, log)

	retryCfg := domain.JobRetryConfig{
		MaxAttempts: 5,
		BaseDelay:   100 * time.Millisecond,
	}

	jobService := job.NewService(jobRepo, func(cmd *domain.WsAgentCommand, retryCfg domain.JobRetryConfig) {
		agentHub.SendCommand(cmd, retryCfg)
	}, retryCfg)

	jobHandler := rest.NewJobHandler(jobService)
	wsAgentHandler := ws.NewAgentHandler(agentHub, log, &ws.AgentHandlerDeps{
		Server: serverService,
		Job:    jobService,
	})

	go hub.Run()
	go agentHub.Run()

	router := rest.NewRouter(cfg, &rest.RouterDeps{
		WsWeb:   wsHandler,
		WsAgent: wsAgentHandler,
		Server:  serverHandler,
		Auth:    authHandler,
		User:    userHandler,
		Job:     jobHandler,

		ServerService: serverService,
	})

	srv := rest.NewServer(router, cfg.Address)

	errCh := make(chan error, 1)
	go func() {
		log.Info("http: starting server", "address", cfg.Address)
		errCh <- srv.ListenAndServe()
		close(errCh)
	}()

	select {
	case <-ctx.Done():
		hub.Stop()
		agentHub.Stop()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			log.Error("http: server shutdown error", "error", err)
		}

	case err := <-errCh:
		if err != nil && err != http.ErrServerClosed {
			log.Error("http: server error", "error", err)
		}
	}

	log.Info("server stopped")
}
