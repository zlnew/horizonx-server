package disk

import (
	"horizonx-server/internal/domain"
	"horizonx-server/internal/logger"
)

type Collector struct {
	log logger.Logger
}

type DiskMetric = domain.DiskMetric
type FilesystemUsage = domain.FilesystemUsage
