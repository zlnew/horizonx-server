package agentws

import (
	"context"

	"horizonx-server/internal/logger"

	"github.com/google/uuid"
)

type Router struct {
	ctx    context.Context
	cancel context.CancelFunc

	agents map[uuid.UUID]*Client

	register   chan *Client
	unregister chan *Client

	log logger.Logger
}

func NewRouter(parent context.Context, log logger.Logger) *Router {
	ctx, cancel := context.WithCancel(parent)

	return &Router{
		ctx:        ctx,
		cancel:     cancel,
		agents:     make(map[uuid.UUID]*Client),
		register:   make(chan *Client, 64),
		unregister: make(chan *Client, 64),
		log:        log,
	}
}

func (r *Router) Run() {
	for {
		select {
		case <-r.ctx.Done():
			r.log.Info("ws: agent router shutting down...")
			for _, agent := range r.agents {
				close(agent.send)
			}
			return

		case a := <-r.register:
			r.agents[a.ID] = a
			a.log.Info("ws: agent registered", "id", a.ID)

		case a := <-r.unregister:
			agent, ok := r.agents[a.ID]
			if !ok {
				continue
			}

			delete(r.agents, a.ID)
			close(agent.send)
			r.log.Info("ws: agent unregistered", "id", a.ID)
		}
	}
}

func (r *Router) Stop() {
	r.cancel()
}
