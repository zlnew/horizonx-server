// Package network
package network

import (
	"context"
)

func NewCollector() *Collector {
	return &Collector{}
}

func (c *Collector) Collect(ctx context.Context) (any, error) {
	return c.collectMetric()
}
