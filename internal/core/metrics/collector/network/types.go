package network

import (
	"time"

	"horizonx-server/internal/logger"
	"horizonx-server/pkg"
	"horizonx-server/pkg/types"
)

type Collector struct {
	log         logger.Logger
	lastRxBytes uint64
	lastTxBytes uint64
	lastTime    time.Time
	rxSpeedEMA  *pkg.EMA
	txSpeedEMA  *pkg.EMA
}

type NetworkMetric = types.NetworkMetric
