package disk

import (
	"zlnew/monitor-agent/internal/core"
	"zlnew/monitor-agent/internal/infra/logger"
)

type Collector struct {
	log logger.Logger
}

type DiskMetric = core.DiskMetric
type FilesystemUsage = core.FilesystemUsage
