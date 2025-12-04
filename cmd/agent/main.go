package main

import (
	"context"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"zlnew/monitor-agent/internal/agent"
	"zlnew/monitor-agent/internal/infra/config"
	"zlnew/monitor-agent/internal/infra/logger"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg := config.Load()
	mode := cfg.Mode
	if len(os.Args) > 1 {
		mode = os.Args[1]
	}
	mode = strings.ToLower(mode)
	cfg.Mode = mode

	log := logger.New(cfg)

	a := agent.New(log, cfg)
	if err := a.Run(ctx); err != nil {
		log.Fatal("agent stopped with error:", err)
	}
}
