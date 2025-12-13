package agent

import (
	"context"
	"encoding/json"
	"time"

	"horizonx-server/internal/domain"

	"github.com/gorilla/websocket"
)

func (a *Agent) readPump(ctx context.Context) error {
	a.conn.SetReadLimit(maxMessageSize)
	a.conn.SetReadDeadline(time.Now().Add(pongWait))
	a.conn.SetPongHandler(func(string) error {
		a.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		default:
			_, message, err := a.conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					a.log.Error("ws read error (unexpected close)", "error", err)
					return err
				} else {
					a.log.Info("ws read finished (normal closure or ping/pong timeout)")
					return nil
				}
			}

			var command domain.WsAgentCommand
			if err := json.Unmarshal(message, &command); err != nil {
				a.log.Error("invalid command payload received", "error", err)
				continue
			}

			select {
			case a.hub.commands <- &command:
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}
}

func (a *Agent) writePump(ctx context.Context) error {
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case message, ok := <-a.send:
			a.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				a.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return nil
			}

			w, err := a.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return err
			}

			_, err = w.Write(message)
			if err != nil {
				w.Close()
				return err
			}

			if err := w.Close(); err != nil {
				return err
			}

		case <-ticker.C:
			a.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := a.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return err
			}
		}
	}
}
