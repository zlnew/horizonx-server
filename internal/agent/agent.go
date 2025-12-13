// Package agent
package agent

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"horizonx-server/internal/domain"
	"horizonx-server/internal/logger"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 8192
)

type Agent struct {
	hub  *Hub
	conn *websocket.Conn

	send chan []byte

	log logger.Logger

	serverURL string
	serverID  int64
	token     string

	sessionCtxCh chan context.Context

	internalEvents chan *domain.WsInternalEvent
}

var ErrUnauthorized = errors.New("connection failed: unauthorized (check token)")

func IsFatalError(err error) bool {
	return errors.Is(err, ErrUnauthorized)
}

func NewAgent(hub *Hub, log logger.Logger, serverURL, token string) *Agent {
	return &Agent{
		hub:            hub,
		log:            log,
		send:           make(chan []byte, 256),
		serverURL:      serverURL,
		serverID:       0,
		token:          token,
		sessionCtxCh:   make(chan context.Context, 1),
		internalEvents: make(chan *domain.WsInternalEvent, 16),
	}
}

func (a *Agent) GetSessionContextChannel() chan context.Context {
	return a.sessionCtxCh
}

func (a *Agent) Run(ctx context.Context) error {
	reconnectInterval := 5 * time.Second
	attempt := 0

	for {
		select {
		case <-ctx.Done():
			a.log.Info("agent run loop received shutdown signal")
			return ctx.Err()
		default:
		}

		a.log.Info("attempting to start agent...", "attempt", attempt+1)

		err := a.start(ctx)

		if err != nil && errors.Is(err, ErrUnauthorized) {
			a.log.Error("FATAL: unauthorized token. exiting.", "error", err)
			return err
		}

		if err != nil {
			a.log.Error("agent session failed, retrying...", "error", err)
		} else {
			a.log.Info("agent session ended, attempting reconnect.")
		}

		attempt++
		a.log.Debug("waiting before next reconnection attempt")

		select {
		case <-time.After(reconnectInterval):
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (a *Agent) start(ctx context.Context) error {
	dialer := websocket.Dialer{
		HandshakeTimeout: 5 * time.Second,
	}

	header := make(http.Header)
	header.Set("Authorization", "Bearer "+a.token)

	conn, resp, err := dialer.DialContext(ctx, a.serverURL, header)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusUnauthorized {
			return ErrUnauthorized
		}
		return fmt.Errorf("ws dial failed: %w", err)
	}

	a.conn = conn
	a.log.Info("ws connected to server", "url", a.serverURL)

	sessionCtx, cancel := context.WithCancel(ctx)
	pumpDone := make(chan error, 1)

	defer func() {
		cancel()
		a.conn.Close()
		a.serverID = 0

		for {
			select {
			case <-a.internalEvents:
			default:
				goto cleanupDone
			}
		}

	cleanupDone:
		a.log.Info("ws connection closed and resources cleaned up")
	}()

	go a.forwardHubEvents(sessionCtx)

	go func() { pumpDone <- a.readPump(sessionCtx) }()
	go func() { pumpDone <- a.writePump(sessionCtx) }()

	initTimeout := 15 * time.Second
	initDone := make(chan struct{})
	go func() {
		for ev := range a.internalEvents {
			if ev == nil {
				continue
			}
			if ev.Event == domain.WsEventAgentReady {
				var sid int64
				switch v := ev.Payload.(type) {
				case float64:
					sid = int64(v)
				case int64:
					sid = v
				case int:
					sid = int64(v)
				case json.Number:
					if parsed, err := v.Int64(); err == nil {
						sid = parsed
					}
				default:
					a.log.Error("unexpected payload type for agent ready", "type", fmt.Sprintf("%T", ev.Payload))
				}
				if sid != 0 {
					a.serverID = sid
					a.log.Info("agent initialized with server id", "server_id", a.serverID)
					close(initDone)
					return
				}
			}
		}
	}()

	select {
	case <-initDone:
		select {
		case a.sessionCtxCh <- sessionCtx:
			a.log.Info("session context published to scheduler")
		default:
			a.log.Warn("session context channel is full, failed to publish")
		}
	case <-time.After(initTimeout):
		a.log.Warn("did not receive init from server within timeout; metrics will not start for this session")
	case <-ctx.Done():
		return ctx.Err()
	}

	var finalErr error
	select {
	case finalErr = <-pumpDone:
		a.log.Info("a pump has exited, shutting down agent session", "error", finalErr)
	case <-ctx.Done():
		finalErr = ctx.Err()
		a.log.Info("agent received external shutdown signal, closing session")
	}

	_ = a.conn.SetWriteDeadline(time.Now().Add(writeWait))
	_ = a.conn.WriteMessage(
		websocket.CloseMessage,
		websocket.FormatCloseMessage(websocket.CloseNormalClosure, "shutting down"),
	)

	select {
	case <-time.After(time.Millisecond * 100):
	case pumpErr := <-pumpDone:
		if finalErr == nil {
			finalErr = pumpErr
		}
	}

	return finalErr
}

func (a *Agent) forwardHubEvents(sessionCtx context.Context) {
	if a.hub == nil {
		return
	}

	for {
		select {
		case ev := <-a.hub.agentEvents:
			select {
			case a.internalEvents <- ev:
			default:
				a.log.Warn("internal events buffer full, dropping event", "event", ev.Event)
			}
		case <-sessionCtx.Done():
			return
		}
	}
}

func (a *Agent) SendMetrics(m domain.Metrics) {
	if a.serverID == 0 {
		a.log.Warn("agent not initialized, dropping metric")
		return
	}

	m.ServerID = a.serverID

	channel := domain.GetServerMetricsChannel(m.ServerID)
	event := domain.WsEventServerMetricsReport
	payload, err := json.Marshal(m)
	if err != nil {
		a.log.Error("failed to marchasl metrics payload", "error", err)
		return
	}

	a.hub.BroadcastToServer(&domain.WsClientMessage{
		Type:    domain.WsAgentReport,
		Channel: channel,
		Event:   event,
		Payload: payload,
	})
}
