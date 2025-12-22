// Package application
package application

import (
	"context"
	"time"

	"horizonx-server/internal/domain"
	"horizonx-server/internal/event"
	"horizonx-server/internal/logger"
)

type Listener struct {
	svc domain.ApplicationService
	log logger.Logger
}

func NewListener(svc domain.ApplicationService, log logger.Logger) *Listener {
	return &Listener{
		svc: svc,
		log: log,
	}
}

func (l *Listener) Register(bus *event.Bus) {
	bus.Subscribe("job_started", l.handleJobStarted)
	bus.Subscribe("job_finished", l.handleJobFinished)
}

func (l *Listener) handleJobStarted(event any) {
	evt, ok := event.(domain.EventJobStarted)
	if !ok {
		l.log.Warn("invalid event payload for job_started", "event", event)
		return
	}

	if evt.ApplicationID == nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	switch evt.JobType {
	case domain.JobTypeDeployApp:
		_ = l.updateStatus(ctx, *evt.ApplicationID, domain.AppStatusDeploying)

	case domain.JobTypeStartApp:
		_ = l.updateStatus(ctx, *evt.ApplicationID, domain.AppStatusStarting)

	case domain.JobTypeStopApp:
		_ = l.updateStatus(ctx, *evt.ApplicationID, domain.AppStatusStopping)

	case domain.JobTypeRestartApp:
		_ = l.updateStatus(ctx, *evt.ApplicationID, domain.AppStatusRestarting)
	}
}

func (l *Listener) handleJobFinished(event any) {
	evt, ok := event.(domain.EventJobFinished)
	if !ok {
		l.log.Warn("invalid event payload for job_finished", "event", event)
		return
	}

	if evt.ApplicationID == nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if evt.Status == domain.JobFailed {
		_ = l.svc.UpdateStatus(ctx, *evt.ApplicationID, domain.AppStatusFailed)
		return
	}

	switch evt.JobType {
	case domain.JobTypeDeployApp:
		_ = l.updateStatus(ctx, *evt.ApplicationID, domain.AppStatusRunning)
		err := l.svc.UpdateLastDeployment(ctx, *evt.ApplicationID)
		if err != nil {
			l.log.Error("failed to update application last deployment", "app_id", evt.ApplicationID)
			return
		}
		l.log.Debug("application last deployment updated", "app_id", evt.ApplicationID)

	case domain.JobTypeStartApp:
		_ = l.updateStatus(ctx, *evt.ApplicationID, domain.AppStatusRunning)

	case domain.JobTypeStopApp:
		_ = l.updateStatus(ctx, *evt.ApplicationID, domain.AppStatusStopped)

	case domain.JobTypeRestartApp:
		_ = l.updateStatus(ctx, *evt.ApplicationID, domain.AppStatusRunning)
	}
}

func (l *Listener) updateStatus(ctx context.Context, appID int64, status domain.ApplicationStatus) error {
	err := l.svc.UpdateStatus(ctx, appID, status)
	if err != nil {
		l.log.Error("failed to update application status", "app_id", appID, "error", err)
		return err
	}

	l.log.Debug("application status updated", "app_id", appID)

	return nil
}
