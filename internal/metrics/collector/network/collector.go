// Package network
package network

import (
	"context"

	"horizonx-server/internal/logger"
	"horizonx-server/internal/pkg"
)

func NewCollector(log logger.Logger) *Collector {
	return &Collector{
		log:        log,
		rxSpeedEMA: pkg.NewEMA(0.5),
		txSpeedEMA: pkg.NewEMA(0.5),
	}
}
func (c *Collector) Collect(ctx context.Context) (NetworkMetric, error) {
	return c.collectMetric()
}
