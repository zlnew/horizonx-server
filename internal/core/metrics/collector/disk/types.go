package disk

import (
	"horizonx-server/internal/logger"
	"horizonx-server/pkg/types"
)

type Collector struct {
	log logger.Logger
}

type DiskMetric = types.DiskMetric
type FilesystemUsage = types.FilesystemUsage
