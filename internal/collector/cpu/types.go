package cpu

import (
	"time"
	"zlnew/monitor-agent/internal/core"
)

type Collector struct {
	lastEnergy uint64
	lastTime   time.Time
}

type CPUMetric = core.CPUMetric
