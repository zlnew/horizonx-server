// Package memory
package memory

import "context"

func NewCollector() *Collector {
	return &Collector{}
}

func (c *Collector) Collect(ctx context.Context) (any, error) {
	err := c.readMemInfo()
	if err != nil {
		return MemoryMetric{}, err
	}

	memTotal := c.getMemTotal()
	memAvailable := c.getMemAvailable()
	memUsed := c.getMemUsed()
	swapTotal := c.getSwapTotal()
	swapFree := c.getSwapFree()
	swapUsed := c.getSwapUsed()
	specs, _ := getSpecs()

	return MemoryMetric{
		MemTotal:     memTotal,
		MemAvailable: memAvailable,
		MemUsed:      memUsed,
		SwapTotal:    swapTotal,
		SwapFree:     swapFree,
		SwapUsed:     swapUsed,
		Specs:        specs,
	}, nil
}
