// Package network
package network

import (
	"context"
)

func NewCollector() *Collector {
	return &Collector{}
}

func (c *Collector) Collect(ctx context.Context) (NetworkMetric, error) {
	return c.collectMetric()
}
