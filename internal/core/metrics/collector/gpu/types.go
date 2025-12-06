package gpu

import (
	"horizonx-server/internal/logger"
	"horizonx-server/pkg"
	"horizonx-server/pkg/types"
)

type Collector struct {
	log         logger.Logger
	powerEMA    map[string]*pkg.EMA
	usageEMA    map[string]*pkg.EMA
	fanSpeedEMA map[string]*pkg.EMA
}

type GPUMetric = types.GPUMetric
