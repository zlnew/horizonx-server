package network

import (
	"context"

	"zlnew/monitor-agent/internal/infra/logger"
)

func NewCollector(log logger.Logger) *Collector {
	return &Collector{log: log}
}

func (c *Collector) Collect(ctx context.Context) (NetworkMetric, error) {
	return c.collectMetric()
}
