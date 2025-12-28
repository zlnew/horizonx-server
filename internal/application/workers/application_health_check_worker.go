package workers

import (
	"context"

	"horizonx-server/internal/domain"
	"horizonx-server/internal/logger"
)

type ApplicationHealthCheckWorker struct {
	svc domain.ApplicationService
	log logger.Logger
}

func NewApplicationHealthCheckWorker(svc domain.ApplicationService, log logger.Logger) Worker {
	return &ApplicationHealthCheckWorker{
		svc: svc,
		log: log,
	}
}

func (w *ApplicationHealthCheckWorker) Name() string {
	return "application_health_check"
}

func (w *ApplicationHealthCheckWorker) Run(ctx context.Context) error {
	return nil
}
