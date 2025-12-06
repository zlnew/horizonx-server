package agent

import (
	"context"
	"encoding/json"
	"os"

	"horizonx-server/internal/config"
	"horizonx-server/internal/core"
)

func (a *Agent) Run(ctx context.Context) error {
	switch a.cfg.Mode {
	case config.ModeSnapshot:
		return a.runSnapshot(ctx)
	case config.ModeStream:
		return a.runStream(ctx)
	case config.ModeServe:
		return a.runHTTP(ctx)
	default:
		a.log.Info("unknown mode, defaulting to serve", "mode", a.cfg.Mode)
		return a.runHTTP(ctx)
	}
}

func (a *Agent) runSnapshot(ctx context.Context) error {
	metrics := a.sampler.Collect(ctx)
	return json.NewEncoder(os.Stdout).Encode(metrics)
}

func (a *Agent) runStream(ctx context.Context) error {
	encoder := json.NewEncoder(os.Stdout)
	sched := core.NewScheduler(a.cfg.Interval, a.log, a.sampler.Collect, func(m core.Metrics) {
		if err := encoder.Encode(m); err != nil {
			a.log.Error("stream encode", "error", err)
		}
	})

	sched.Start(ctx)
	return nil
}

func (a *Agent) runHTTP(ctx context.Context) error {
	sched := core.NewScheduler(a.cfg.Interval, a.log, a.sampler.Collect, func(m core.Metrics) {
		a.store.Set(m)
		a.hub.BroadcastMetrics(m)
	})
	go sched.Start(ctx)

	return a.http.Start(ctx)
}
