// Package agent
package agent

import (
	"horizonx-server/internal/config"
	"horizonx-server/internal/core"
	"horizonx-server/internal/logger"
	"horizonx-server/internal/transport/http"
)

type Agent struct {
	log     logger.Logger
	cfg     *config.Config
	sampler *Sampler
	store   *core.SnapshotStore
	http    *http.Server
	hub     *http.Hub
}

func New(log logger.Logger, cfg *config.Config, hub *http.Hub) *Agent {
	sam := NewSampler(log)
	store := core.NewSnapshotStore()
	httpServer := http.NewServer(cfg, store, log, hub)

	return &Agent{
		log:     log,
		cfg:     cfg,
		sampler: sam,
		store:   store,
		http:    httpServer,
		hub:     hub,
	}
}
