package ws

import (
	"context"
	"encoding/json"
	"math"
	"time"

	"horizonx-server/internal/domain"
	"horizonx-server/internal/logger"

	"github.com/google/uuid"
)

type AgentHub struct {
	ctx    context.Context
	cancel context.CancelFunc

	agents map[uuid.UUID]*Agent

	register   chan *Agent
	unregister chan *Agent
	commands   chan *domain.WsAgentCommand

	log logger.Logger
}

func NewAgentHub(parent context.Context, log logger.Logger) *AgentHub {
	ctx, cancel := context.WithCancel(parent)

	return &AgentHub{
		ctx:        ctx,
		cancel:     cancel,
		agents:     make(map[uuid.UUID]*Agent),
		register:   make(chan *Agent, 64),
		unregister: make(chan *Agent, 64),
		commands:   make(chan *domain.WsAgentCommand, 1024),
		log:        log,
	}
}

func (h *AgentHub) Run() {
	for {
		select {
		case <-h.ctx.Done():
			h.log.Info("ws: agent hub shutting down...")
			for _, agent := range h.agents {
				close(agent.send)
			}
			return

		case a := <-h.register:
			h.agents[a.ID] = a
			a.log.Info("ws: agent registered", "id", a.ID)

		case a := <-h.unregister:
			agent, ok := h.agents[a.ID]
			if !ok {
				continue
			}

			delete(h.agents, a.ID)
			close(agent.send)
			h.log.Info("ws: agent unregistered", "id", a.ID)

		case cmd := <-h.commands:
			h.handleCommand(cmd)
		}
	}
}

func (h *AgentHub) Stop() {
	h.cancel()
}

func (h *AgentHub) SendCommand(cmd *domain.WsAgentCommand, retryCfg domain.JobRetryConfig) {
	go func() {
		for attempt := 0; attempt < retryCfg.MaxAttempts; attempt++ {
			select {
			case h.commands <- cmd:
				return
			case <-h.ctx.Done():
				h.log.Warn("ws: hub stopped, abort send command", "command", cmd.CommandType)
				return
			default:
				delay := retryCfg.BaseDelay * time.Duration(math.Pow(2, float64(attempt)))
				h.log.Warn("ws: command buffer full, retrying", "command", cmd.CommandType, "attempt", attempt+1, "delay_ms", delay.Milliseconds())
				time.Sleep(delay)
			}
		}
		h.log.Error("ws: command dropped after retries", "command", cmd.CommandType)
	}()
}

func (h *AgentHub) handleCommand(cmd *domain.WsAgentCommand) {
	agent, ok := h.agents[cmd.TargetServerID]
	if !ok {
		h.log.Warn("ws: target agent not connected", "server_id", cmd.TargetServerID)
		return
	}

	message, err := json.Marshal(cmd)
	if err != nil {
		h.log.Error("ws: failed to marshal server command", "error", err)
		return
	}

	select {
	case agent.send <- message:
		h.log.Info("ws: command sent to agent", "server_id", agent.ID, "command", cmd.CommandType)
	default:
		h.log.Warn("ws: agent channel full, force unregister", "server_id", agent.ID)
		h.unregister <- agent
	}
}
