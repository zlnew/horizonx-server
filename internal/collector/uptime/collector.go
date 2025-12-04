// Package uptime
package uptime

import (
	"context"
	"os"
	"strconv"
	"strings"
)

type Collector struct{}

func NewCollector() *Collector {
	return &Collector{}
}

func (c *Collector) Collect(ctx context.Context) (float64, error) {
	data, err := os.ReadFile("/proc/uptime")
	if err != nil {
		return float64(0), err
	}

	fields := strings.Fields(string(data))
	if len(fields) == 0 {
		return float64(0), nil
	}

	seconds, err := strconv.ParseFloat(fields[0], 64)
	if err != nil {
		return float64(0), err
	}

	return seconds, nil
}
