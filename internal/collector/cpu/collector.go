package cpu

import (
	"context"

	"zlnew/monitor-agent/internal/core"
	"zlnew/monitor-agent/internal/infra/logger"
)

func NewCollector(log logger.Logger) *Collector {
	return &Collector{
		log:      log,
		powerEMA: core.NewEMA(0.3),
	}
}

func (c *Collector) Collect(ctx context.Context) (CPUMetric, error) {
	usage, perCore := c.readUsage()
	temperature := c.readTemperature()
	frequency := c.readFrequency()
	powerWatt := c.readPowerWatt()

	return CPUMetric{
		Usage:       usage,
		PerCore:     perCore,
		Temperature: temperature,
		Frequency:   frequency,
		PowerWatt:   powerWatt,
	}, nil
}
