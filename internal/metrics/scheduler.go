package metrics

import (
	"context"
	"time"

	"horizonx-server/internal/domain"
	"horizonx-server/internal/logger"
)

type Scheduler struct {
	interval time.Duration
	log      logger.Logger
	sample   func(context.Context) domain.Metrics
	sink     func(domain.Metrics)
}

func NewScheduler(interval time.Duration, log logger.Logger, sample func(context.Context) domain.Metrics, sink func(domain.Metrics)) *Scheduler {
	return &Scheduler{
		interval: interval,
		log:      log,
		sample:   sample,
		sink:     sink,
	}
}

func (s *Scheduler) Start(ctx context.Context) {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	go s.tick(ctx)

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

	timeoutCtx, cancel := context.WithTimeout(ctx, s.interval)
	defer cancel()

	m := s.sample(timeoutCtx)
	s.sink(m)
}
