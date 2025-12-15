package network

import (
	"time"

	"horizonx-server/internal/domain"
	"horizonx-server/internal/logger"
	"horizonx-server/internal/pkg"
)

type Collector struct {
	log         logger.Logger
	lastRxBytes uint64
	lastTxBytes uint64
	lastTime    time.Time
	rxSpeedEMA  *pkg.EMA
	txSpeedEMA  *pkg.EMA
}

type NetworkMetric = domain.NetworkMetric
