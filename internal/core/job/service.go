// Package job
package job

import (
	"context"

	"horizonx-server/internal/domain"

	"github.com/google/uuid"
)

type JobService struct {
	repo    domain.JobRepository
	publish func(cmd *domain.WsAgentCommand, retryCfg domain.JobRetryConfig)
	retry   domain.JobRetryConfig
}

func NewService(
	repo domain.JobRepository,
	publish func(cmd *domain.WsAgentCommand, retryCfg domain.JobRetryConfig),
	retry domain.JobRetryConfig,
) domain.JobService {
	return &JobService{
		repo:    repo,
		publish: publish,
		retry:   retry,
	}
}

func (s *JobService) Get(ctx context.Context) ([]domain.Job, error) {
	return s.repo.List(ctx)
}

func (s *JobService) Create(ctx context.Context, j *domain.Job) (*domain.Job, error) {
	return s.repo.Create(ctx, j)
}

func (s *JobService) Delete(ctx context.Context, jobID int64) error {
	return s.repo.Delete(ctx, jobID)
}

func (s *JobService) Start(ctx context.Context, jobID int64) (*domain.Job, error) {
	return s.repo.MarkRunning(ctx, jobID)
}

func (s *JobService) Finish(ctx context.Context, jobID int64, status domain.JobStatus, outputLog *string) (*domain.Job, error) {
	return s.repo.MarkFinished(ctx, jobID, status, outputLog)
}

func (s *JobService) InitAgent(ctx context.Context, serverID uuid.UUID) (*domain.Job, error) {
	job := &domain.Job{
		ServerID: serverID,
		JobType:  domain.WsCommandAgentInit,
	}

	job, err := s.Create(ctx, job)
	if err != nil {
		return nil, err
	}

	go func(job *domain.Job) {
		cmd := &domain.WsAgentCommand{
			TargetServerID: job.ServerID,
			CommandType:    job.JobType,
			Payload: domain.JobCommandPayload{
				JobID: job.ID,
			},
		}

		if s.publish != nil {
			s.publish(cmd, s.retry)
		}
	}(job)

	return job, nil
}
