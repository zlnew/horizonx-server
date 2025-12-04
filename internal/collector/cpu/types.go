package cpu

import (
	"time"

	"zlnew/monitor-agent/internal/core"
	"zlnew/monitor-agent/internal/infra/logger"
)

type Collector struct {
	log        logger.Logger
	lastEnergy uint64
	lastTime   time.Time
	powerEMA   *core.EMA
}

type CPUMetric = core.CPUMetric
