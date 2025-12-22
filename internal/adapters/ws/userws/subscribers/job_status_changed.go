package subscribers

import (
	"fmt"

	"horizonx-server/internal/adapters/ws/userws"
	"horizonx-server/internal/domain"
)

type JobStatusChanged struct {
	hub *userws.Hub
}

func NewJobStatusChanged(hub *userws.Hub) *JobStatusChanged {
	return &JobStatusChanged{hub: hub}
}

func (s *JobStatusChanged) Handle(event any) {
	evt, ok := event.(domain.EventJobStatusChanged)
	if !ok {
		return
	}

	s.hub.Broadcast(&domain.WsServerEvent{
		Channel: fmt.Sprintf("job:%d", evt.JobID),
		Event:   "job_status_changed",
		Payload: evt,
	})

	s.hub.Broadcast(&domain.WsServerEvent{
		Channel: "jobs",
		Event:   "job_status_changed",
		Payload: evt,
	})
}
