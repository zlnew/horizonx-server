// Package websocket
package websocket

import (
	"encoding/json"

	"horizonx-server/internal/logger"
)

type Hub struct {
	rooms      map[string]map[*Client]bool
	register   chan *Subscription
	unregister chan *Subscription
	events     chan *ServerEvent
	log        logger.Logger
}

type Subscription struct {
	client  *Client
	channel string
}

type ServerEvent struct {
	Channel string
	Event   string
	Payload any
}

func NewHub(log logger.Logger) *Hub {
	return &Hub{
		rooms:      make(map[string]map[*Client]bool),
		register:   make(chan *Subscription),
		unregister: make(chan *Subscription),
		events:     make(chan *ServerEvent),
		log:        log,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case sub := <-h.register:
			if _, ok := h.rooms[sub.channel]; !ok {
				h.rooms[sub.channel] = make(map[*Client]bool)
			}
			h.rooms[sub.channel][sub.client] = true
			h.log.Info("client subscribed", "channel", sub.channel, "remote_addr", sub.client.conn.RemoteAddr())

		case sub := <-h.unregister:
			if clients, ok := h.rooms[sub.channel]; ok {
				if _, ok := clients[sub.client]; ok {
					delete(clients, sub.client)
					if len(clients) == 0 {
						delete(h.rooms, sub.channel)
						h.log.Info("channel closed", "channel", sub.channel)
					}
					h.log.Info("client unsubscribed", "channel", sub.channel, "remote_addr", sub.client.conn.RemoteAddr())
				}
			}

		case evt := <-h.events:
			data := map[string]any{
				"type":    "event",
				"event":   evt.Event,
				"channel": evt.Channel,
				"payload": evt.Payload,
			}

			bytes, _ := json.Marshal(data)

			if clients, ok := h.rooms[evt.Channel]; ok {
				for client := range clients {
					select {
					case client.send <- bytes:
					default:
						h.log.Warn("client send channel full, skipping", "remote_addr", client.conn.RemoteAddr())
					}
				}
			}
		}
	}
}

func (h *Hub) Emit(channel, event string, payload any) {
	h.events <- &ServerEvent{
		Channel: channel,
		Event:   event,
		Payload: payload,
	}
}
