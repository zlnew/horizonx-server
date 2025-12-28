// Package metrics
package metrics

import (
	"context"
	"sync"
	"time"

	"horizonx-server/internal/config"
	"horizonx-server/internal/domain"
	"horizonx-server/internal/logger"
	"horizonx-server/internal/metrics/collector/cpu"
	"horizonx-server/internal/metrics/collector/disk"
	"horizonx-server/internal/metrics/collector/gpu"
	"horizonx-server/internal/metrics/collector/memory"
	"horizonx-server/internal/metrics/collector/network"
	"horizonx-server/internal/metrics/collector/os"
	"horizonx-server/internal/metrics/collector/uptime"
)

type Collector struct {
	os      *os.Collector
	cpu     *cpu.Collector
	gpu     *gpu.Collector
	memory  *memory.Collector
	disk    *disk.Collector
	network *network.Collector
	uptime  *uptime.Collector

	cfg *config.Config
	log logger.Logger

	buffer    []domain.Metrics
	bufferMu  sync.Mutex
	batchSize int
	interval  time.Duration
}

func NewCollector(cfg *config.Config, log logger.Logger) *Collector {
	return &Collector{
		cfg: cfg,

		os:      os.NewCollector(log),
		cpu:     cpu.NewCollector(log),
		gpu:     gpu.NewCollector(log),
		memory:  memory.NewCollector(log),
		disk:    disk.NewCollector(log),
		network: network.NewCollector(log),
		uptime:  uptime.NewCollector(log),

		log: log,

		buffer:    make([]domain.Metrics, 0, 10),
		batchSize: 10,
		interval:  5 * time.Second,
	}
}

func (s *Collector) Start(ctx context.Context) error {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	s.log.Info("metrics collector started")

	for {
		select {
		case <-ctx.Done():
			s.log.Info("metrics collector stopping...")
			return ctx.Err()
		case <-ticker.C:
			s.collect(ctx)
		}
	}
}

func (s *Collector) Latest() domain.Metrics {
	s.bufferMu.Lock()
	defer s.bufferMu.Unlock()

	if len(s.buffer) == 0 {
		return domain.Metrics{}
	}
	return s.buffer[len(s.buffer)-1]
}

func (s *Collector) collect(ctx context.Context) {
	var metrics domain.Metrics
	metrics.ServerID = s.cfg.AgentServerID
	metrics.RecordedAt = time.Now().UTC()

	if val, err := s.os.Collect(ctx); err != nil {
		s.log.Error("collector", "name", "os", "error", err)
	} else {
		metrics.OSInfo = val
	}

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

	s.bufferMu.Lock()
	defer s.bufferMu.Unlock()

	if len(s.buffer) >= s.batchSize {
		s.buffer = s.buffer[1:]
	}
	s.buffer = append(s.buffer, metrics)
}
