// Package deployment
package deployment

import (
	"horizonx-server/internal/domain"
	"horizonx-server/internal/event"
	"horizonx-server/internal/logger"
)

type Listener struct {
	svc domain.DeploymentService
	log logger.Logger
}

func NewListener(svc domain.DeploymentService, log logger.Logger) *Listener {
	return &Listener{
		svc: svc,
		log: log,
	}
}

func (l *Listener) Register(bus *event.Bus) {
	bus.Subscribe("job_started", l.handleJobStarted)
	bus.Subscribe("job_finished", l.handleJobFinished)
}

func (l *Listener) handleJobStarted(event any) {}

func (l *Listener) handleJobFinished(event any) {}
