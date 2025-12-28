package workers

import (
	"context"
	"time"

	"horizonx-server/internal/logger"
)

type DailySchedule struct {
	Hour   int
	Minute int
}

type Scheduler struct {
	log logger.Logger
}

func NewScheduler(log logger.Logger) *Scheduler {
	return &Scheduler{log: log}
}

func (s *Scheduler) RunByDuration(ctx context.Context, dur time.Duration, worker Worker) {
	go func() {
		ticker := time.NewTicker(dur)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				start := time.Now()

				err := worker.Run(ctx)
				if err != nil {
					s.log.Error("worker failed", "name", worker.Name(), "error", err)
				}

				s.log.Debug("worker finished", "name", worker.Name(), "time", time.Since(start))
			}
		}
	}()
}

func (s *Scheduler) RunDaily(ctx context.Context, schedule DailySchedule, worker Worker) {
	go func() {
		now := time.Now()

		next := time.Date(
			now.Year(),
			now.Month(),
			now.Day(),
			schedule.Hour,
			schedule.Minute,
			0,
			0,
			now.Location(),
		)

		if next.Before(now) {
			next = next.Add(24 * time.Hour)
		}

		timer := time.NewTimer(time.Until(next))
		defer timer.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-timer.C:
				start := time.Now()

				err := worker.Run(ctx)
				if err != nil {
					s.log.Error("worker failed", "name", worker.Name(), "error", err)
				}

				s.log.Debug("worker finished", "name", worker.Name(), "time", time.Since(start))

				timer.Reset(24 * time.Hour)
			}
		}
	}()
}
