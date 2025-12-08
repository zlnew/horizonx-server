package websocket

import (
	"errors"
	"net/http"
	"slices"

	"horizonx-server/internal/config"
	"horizonx-server/internal/logger"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
)

type Handler struct {
	hub      *Hub
	upgrader websocket.Upgrader
	log      logger.Logger
	verify   func(token string) error
}

func NewHandler(hub *Hub, log logger.Logger, cfg *config.Config) *Handler {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			origin := r.Header.Get("Origin")

			allowed := slices.Contains(cfg.AllowedOrigins, origin)
			if !allowed {
				log.Warn("websocket origin rejected", "origin", origin)
				return false
			}

			return allowed
		},
	}

	verify := func(token string) error {
		if token == "" {
			return errors.New("empty token")
		}

		parsed, err := jwt.Parse(token, func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(cfg.JWTSecret), nil
		})
		if err != nil {
			log.Warn("websocket jwt parse failed", "error", err)
			return err
		}

		if !parsed.Valid {
			log.Warn("websocket jwt invalid token")
			return errors.New("invalid token")
		}

		return nil
	}

	return &Handler{
		hub:      hub,
		upgrader: upgrader,
		log:      log,
		verify:   verify,
	}
}

func (h *Handler) Serve(w http.ResponseWriter, r *http.Request) {
	token := ""
	authHeader := r.Header.Get("Authorization")
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		token = authHeader[7:]
	}

	if token == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	if err := h.verify(token); err != nil {
		h.log.Warn("jwt verification failed", "error", err)
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}

	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.log.Error("upgrade failed", "error", err)
		return
	}

	client := NewClient(h.hub, conn, h.log)
	go client.writePump()
	go client.readPump()

	h.log.Info("client connected", "remote_addr", conn.RemoteAddr())
}
