package gpu

import (
	"horizonx-server/internal/domain"
	"horizonx-server/internal/logger"
	"horizonx-server/internal/pkg"
)

type Collector struct {
	log         logger.Logger
	powerEMA    map[string]*pkg.EMA
	usageEMA    map[string]*pkg.EMA
	fanSpeedEMA map[string]*pkg.EMA
}

type GPUMetric = domain.GPUMetric
