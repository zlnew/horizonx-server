package memory

import (
	"horizonx-server/internal/logger"
	"horizonx-server/pkg/types"
)

type Collector struct {
	log          logger.Logger
	memTotal     uint64
	memAvailable uint64
	swapTotal    uint64
	swapFree     uint64
}

type MemoryMetric = types.MemoryMetric
