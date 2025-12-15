// Package memory
package memory

import (
	"context"

	"horizonx-server/internal/logger"
)

func NewCollector(log logger.Logger) *Collector {
	return &Collector{log: log}
}

func (c *Collector) Collect(ctx context.Context) (MemoryMetric, error) {
	err := c.readMemInfo()
	if err != nil {
		return MemoryMetric{}, err
	}

	memTotalGB := c.getMemTotalGB()
	memAvailableGB := c.getMemAvailableGB()
	memUsedGB := c.getMemUsedGB()
	swapTotalGB := c.getSwapTotalGB()
	swapFreeGB := c.getSwapFreeGB()
	swapUsedGB := c.getSwapUsedGB()

	usagePercent := 0.0
	if memTotalGB > 0 {
		usagePercent = (1 - (memAvailableGB / memTotalGB)) * 100
	}

	return MemoryMetric{
		TotalGB:      memTotalGB,
		AvailableGB:  memAvailableGB,
		UsedGB:       memUsedGB,
		UsagePercent: usagePercent,
		SwapTotalGB:  swapTotalGB,
		SwapFreeGB:   swapFreeGB,
		SwapUsedGB:   swapUsedGB,
	}, nil
}
