package agent

import (
	"context"
	"time"
	"zlnew/monitor-agent/internal/core"
)

func (a *Agent) Run(ctx context.Context) error {
	sched := core.NewScheduler(a.reg, time.Second*1, a.log)
	go sched.Start(ctx)

	return a.http.Start(ctx)
}
