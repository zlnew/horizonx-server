package workers

import (
	"context"

	"horizonx-server/internal/domain"
	"horizonx-server/internal/logger"
)

type MetricsCleanupWorker struct {
	svc domain.MetricsService
	log logger.Logger
}

func NewMetricsCleanupWorker(svc domain.MetricsService, log logger.Logger) Worker {
	return &MetricsCleanupWorker{
		svc: svc,
		log: log,
	}
}

func (w *MetricsCleanupWorker) Name() string {
	return "metrics_cleanup"
}

func (w *MetricsCleanupWorker) Run(ctx context.Context) error {
	return nil
}
