// Package uptime
package uptime

import (
	"context"
	"os"
	"strconv"
	"strings"

	"zlnew/monitor-agent/internal/infra/logger"
)

type Collector struct {
	log logger.Logger
}

func NewCollector(log logger.Logger) *Collector {
	return &Collector{log: log}
}

func (c *Collector) Collect(ctx context.Context) (float64, error) {
	data, err := os.ReadFile("/proc/uptime")
	if err != nil {
		c.log.Error("failed to read /proc/uptime", "error", err)
		return float64(0), err
	}

	fields := strings.Fields(string(data))
	if len(fields) == 0 {
		c.log.Warn("/proc/uptime is empty")
		return float64(0), nil
	}

	seconds, err := strconv.ParseFloat(fields[0], 64)
	if err != nil {
		c.log.Warn("failed to parse uptime", "value", fields[0], "error", err)
		return float64(0), err
	}

	return seconds, nil
}
