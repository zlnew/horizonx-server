// Package system
package system

import "context"

func NewCollector() *Collector {
	return &Collector{}
}

func (c *Collector) Collect(ctx context.Context) (any, error) {
	hostname := getHostname()
	OS := getOSName()
	kernel := getKernelVersion()
	uptime := getUptime()

	return SystemMetric{
		Hostname: hostname,
		OS:       OS,
		Kernel:   kernel,
		Uptime:   uptime,
	}, nil
}
