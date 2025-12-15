// Package cpu
package cpu

import (
	"context"

	"horizonx-server/internal/logger"
	"horizonx-server/internal/pkg"
)

func NewCollector(log logger.Logger) *Collector {
	return &Collector{
		log:          log,
		powerEMA:     pkg.NewEMA(0.3),
		prevCPUStats: make(map[string]cpuStat),
		usageEMA:     pkg.NewEMA(0.5),
	}
}

func (c *Collector) Collect(ctx context.Context) (CPUMetric, error) {
	usage, perCore := c.readUsage()
	temperature := c.readTemperature()
	frequency := c.readFrequency()
	powerWatt := c.readPowerWatt()

	c.usageEMA.Add(usage)
	if len(c.perCoreEMA) != len(perCore) {
		c.perCoreEMA = make([]*pkg.EMA, len(perCore))
		for i := range c.perCoreEMA {
			c.perCoreEMA[i] = pkg.NewEMA(0.5)
		}
	}
	for i, coreUsage := range perCore {
		c.perCoreEMA[i].Add(coreUsage)
	}

	smoothedPerCore := make([]float64, len(c.perCoreEMA))
	for i, ema := range c.perCoreEMA {
		smoothedPerCore[i] = ema.Value()
	}

	return CPUMetric{
		Usage:       c.usageEMA.Value(),
		PerCore:     smoothedPerCore,
		Temperature: temperature,
		Frequency:   frequency,
		PowerWatt:   powerWatt,
	}, nil
}
