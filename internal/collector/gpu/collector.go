// Package gpu
package gpu

import "context"

func NewCollector() *Collector {
	return &Collector{}
}

func (c *Collector) Collect(ctx context.Context) (any, error) {
	usage := getUsage()
	temp, _ := getTemp()
	vramTotal := getVramTotal()
	vramUsed := getVramUsed()
	watt, _ := getPower()
	spec, _ := getSpec()

	return GPUMetric{
		Spec:      spec,
		Usage:     usage,
		Temp:      temp,
		VramTotal: vramTotal,
		VramUsed:  vramUsed,
		Watt:      watt,
	}, nil
}
