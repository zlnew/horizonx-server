package os

import (
	"horizonx-server/internal/domain"
	"horizonx-server/internal/logger"
)

type Collector struct {
	log logger.Logger
}

type OSInfo = domain.OSInfo
