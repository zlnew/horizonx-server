package cpu

import (
	"time"

	"horizonx-server/internal/domain"
	"horizonx-server/internal/logger"
	"horizonx-server/internal/pkg"
)

type Collector struct {
	log          logger.Logger
	lastEnergy   uint64
	lastTime     time.Time
	powerEMA     *pkg.EMA
	prevCPUStats map[string]cpuStat
	usageEMA     *pkg.EMA
	perCoreEMA   []*pkg.EMA
}

type cpuStat struct {
	total uint64
	idle  uint64
}

type CPUMetric = domain.CPUMetric
