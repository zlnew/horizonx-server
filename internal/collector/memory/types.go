package memory

import (
	"zlnew/monitor-agent/internal/core"
	"zlnew/monitor-agent/internal/infra/logger"
)

type Collector struct {
	log          logger.Logger
	memTotal     uint64
	memAvailable uint64
	swapTotal    uint64
	swapFree     uint64
}

type MemoryMetric = core.MemoryMetric
