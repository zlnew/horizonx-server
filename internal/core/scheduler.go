package core

import (
	"context"
	"time"

	"zlnew/monitor-agent/internal/infra/logger"
)

type Scheduler struct {
	interval time.Duration
	log      logger.Logger
	sample   func(context.Context) Metrics
	sink     func(Metrics)
}

func NewScheduler(interval time.Duration, log logger.Logger, sample func(context.Context) Metrics, sink func(Metrics)) *Scheduler {
	return &Scheduler{interval: interval, log: log, sample: sample, sink: sink}
}

func (s *Scheduler) Start(ctx context.Context) {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	s.tick(ctx)

	for {
		select {
		case <-ticker.C:
			s.tick(ctx)
		case <-ctx.Done():
			return
		}
	}
}

func (s *Scheduler) tick(ctx context.Context) {
	if s.sample == nil || s.sink == nil {
		return
	}

	metrics := s.sample(ctx)
	s.sink(metrics)
}
