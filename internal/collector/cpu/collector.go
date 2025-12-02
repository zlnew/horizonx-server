// Package cpu
package cpu

import (
	"context"
)

func NewCollector() *Collector {
	return &Collector{}
}

func (c *Collector) Collect(ctx context.Context) (any, error) {
	spec, _ := getSpec()
	usage, perCore, _ := getUsage()
	watt, _ := c.getWatt()
	temp, _ := getTemp()
	freq, _ := getFreq()

	return CPUMetric{
		Spec:      spec,
		Usage:     usage,
		PerCore:   perCore,
		Watt:      watt,
		Temp:      temp,
		Frequency: freq,
	}, nil
}
