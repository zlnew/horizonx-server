// Package cpu
package cpu

import (
	"context"
)

func NewCollector() *Collector {
	return &Collector{}
}

func (c *Collector) Collect(ctx context.Context) (CPUMetric, error) {
	usage, perCore := getUsage()
	temperature := getTemperature()
	frequency := getFrequency()
	powerWatt := c.getPowerWatt()

	return CPUMetric{
		Usage:       usage,
		PerCore:     perCore,
		Temperature: temperature,
		Frequency:   frequency,
		PowerWatt:   powerWatt,
	}, nil
}
