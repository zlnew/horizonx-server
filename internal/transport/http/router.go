// Package http
package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"zlnew/monitor-agent/internal/core"
	"zlnew/monitor-agent/internal/infra/config"
	"zlnew/monitor-agent/internal/infra/logger"
)

type Server struct {
	store *core.SnapshotStore
	log   logger.Logger
	cfg   *config.Config
	srv   *http.Server
}

func NewServer(cfg *config.Config, store *core.SnapshotStore, log logger.Logger) *Server {
	return &Server{cfg: cfg, store: store, log: log}
}

func (s *Server) Start(ctx context.Context) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/metrics", s.handleMetrics)

	s.srv = &http.Server{
		Addr:    s.cfg.Address,
		Handler: mux,
	}

	errCh := make(chan error, 1)
	go func() {
		s.log.Info("starting http server on " + s.cfg.Address)
		if err := s.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
		close(errCh)
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := s.srv.Shutdown(shutdownCtx); err != nil {
			return err
		}
		return nil
	case err := <-errCh:
		return err
	}
}

func (s *Server) handleMetrics(w http.ResponseWriter, r *http.Request) {
	data := s.store.Get()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
