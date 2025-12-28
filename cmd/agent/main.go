package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"golang.org/x/sync/errgroup"

	"horizonx-server/internal/agent"
	"horizonx-server/internal/config"
	"horizonx-server/internal/logger"
	"horizonx-server/internal/metrics"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("INFO: No .env file found, relying on system environment variables")
	}

	cfg := config.Load()
	appLog := logger.New(cfg)

	if cfg.AgentServerAPIToken == "" {
		log.Fatal("FATAL: HORIZONX_SERVER_API_TOKEN is missing in .env or system vars!")
	}

	if cfg.AgentServerID.String() == "00000000-0000-0000-0000-000000000000" {
		log.Fatal("FATAL: HORIZONX_SERVER_ID is missing or invalid in .env!")
	}

	appLog.Info("horizonx agent: starting...", "server_id", cfg.AgentServerID)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Initialize components
	ws := agent.NewAgent(cfg, appLog)
	mCollector := metrics.NewCollector(cfg, appLog)

	// Initialize job worker
	jWorker := agent.NewJobWorker(cfg, appLog, mCollector.Latest)
	if err := jWorker.Initialize(); err != nil {
		appLog.Error("failed to Initialize job worker", "error", err)
		log.Fatal(err)
	}

	g, gCtx := errgroup.WithContext(ctx)

	// WebSocket connection
	g.Go(func() error {
		return ws.Run(gCtx)
	})

	// Metrics collector
	g.Go(func() error {
		return mCollector.Start(gCtx)
	})

	// Job worker
	g.Go(func() error {
		return jWorker.Start(gCtx)
	})

	if err := g.Wait(); err != nil && err != context.Canceled && !agent.IsFatalError(err) {
		appLog.Error("agent failed unexpectedly", "error", err)
	} else if agent.IsFatalError(err) {
		appLog.Error("agent failed fatally, exiting", "error", err)
	}

	appLog.Info("agent stopped gracefully.")
}
