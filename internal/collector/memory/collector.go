// Package memory
package memory

import "context"

func NewCollector() *Collector {
	return &Collector{}
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

	return MemoryMetric{
		TotalGB:     memTotalGB,
		AvailableGB: memAvailableGB,
		UsedGB:      memUsedGB,
		SwapTotalGB: swapTotalGB,
		SwapFreeGB:  swapFreeGB,
		SwapUsedGB:  swapUsedGB,
	}, nil
}
