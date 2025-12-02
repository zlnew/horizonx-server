package network

import (
	"time"

	"zlnew/monitor-agent/internal/core"
)

type Collector struct {
	lastRxBytes uint64
	lastTxBytes uint64
	lastTime    time.Time
}

type NetworkMetric = core.NetworkMetric
