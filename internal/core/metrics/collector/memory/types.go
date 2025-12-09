package memory

import (
	"horizonx-server/internal/domain"
	"horizonx-server/internal/logger"
)

type Collector struct {
	log          logger.Logger
	memTotal     uint64
	memAvailable uint64
	swapTotal    uint64
	swapFree     uint64
}

type MemoryMetric = domain.MemoryMetric
