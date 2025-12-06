package metrics

import (
	"context"

	"horizonx-server/internal/core/metrics/collector/cpu"
	"horizonx-server/internal/core/metrics/collector/disk"
	"horizonx-server/internal/core/metrics/collector/gpu"
	"horizonx-server/internal/core/metrics/collector/memory"
	"horizonx-server/internal/core/metrics/collector/network"
	"horizonx-server/internal/core/metrics/collector/uptime"
	"horizonx-server/internal/logger"
	"horizonx-server/pkg/types"
)

type Sampler struct {
	cpu     *cpu.Collector
	gpu     *gpu.Collector
	memory  *memory.Collector
	disk    *disk.Collector
	network *network.Collector
	uptime  *uptime.Collector
	log     logger.Logger
}

func NewSampler(log logger.Logger) *Sampler {
	return &Sampler{
		cpu:     cpu.NewCollector(log),
		gpu:     gpu.NewCollector(log),
		memory:  memory.NewCollector(log),
		disk:    disk.NewCollector(log),
		network: network.NewCollector(log),
		uptime:  uptime.NewCollector(log),
		log:     log,
	}
}

func (s *Sampler) Collect(ctx context.Context) types.Metrics {
	var metrics types.Metrics

	if val, err := s.cpu.Collect(ctx); err != nil {
		s.log.Error("collector", "name", "cpu", "error", err)
	} else {
		metrics.CPU = val
	}

	if val, err := s.gpu.Collect(ctx); err != nil {
		s.log.Error("collector", "name", "gpu", "error", err)
	} else {
		metrics.GPU = val
	}

	if val, err := s.memory.Collect(ctx); err != nil {
		s.log.Error("collector", "name", "memory", "error", err)
	} else {
		metrics.Memory = val
	}

	if val, err := s.disk.Collect(ctx); err != nil {
		s.log.Error("collector", "name", "disk", "error", err)
	} else {
		metrics.Disk = val
	}

	if val, err := s.network.Collect(ctx); err != nil {
		s.log.Error("collector", "name", "network", "error", err)
	} else {
		metrics.Network = val
	}

	if val, err := s.uptime.Collect(ctx); err != nil {
		s.log.Error("collector", "name", "uptime", "error", err)
	} else {
		metrics.UptimeSeconds = val
	}

	return metrics
}
