// Package http
package http

import (
	"context"
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"syscall"
	"time"

	"github.com/gorilla/websocket"
	"zlnew/monitor-agent/internal/core"
	"zlnew/monitor-agent/internal/infra/config"
	"zlnew/monitor-agent/internal/infra/logger"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

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
	mux.HandleFunc("/ws", s.handleWs)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		mux.ServeHTTP(w, r)
	})

	s.srv = &http.Server{
		Addr:    s.cfg.Address,
		Handler: handler,
	}

	errCh := make(chan error, 1)
	go func() {
		s.log.Info("starting http server", "address", s.cfg.Address)
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

func (s *Server) handleWs(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.log.Error("upgrade", "error", err)
		return
	}
	defer conn.Close()

	ticker := time.NewTicker(s.cfg.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-r.Context().Done():
			return
		case <-ticker.C:
			data := s.store.Get()
			if err := conn.WriteJSON(data); err != nil {
				var netErr *net.OpError
				if errors.As(err, &netErr) && errors.Is(netErr.Err, syscall.EPIPE) {
					s.log.Debug("client disconnected", "remote_addr", conn.RemoteAddr())
					return
				}
				s.log.Error("write", "error", err)
				return
			}
		}
	}
}
