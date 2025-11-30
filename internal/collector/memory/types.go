package memory

import "zlnew/monitor-agent/internal/core"

type Collector struct {
	memTotal uint64
	memAvailable uint64
	swapTotal uint64
	swapFree uint64
}

type MemoryMetric = core.MemoryMetric
