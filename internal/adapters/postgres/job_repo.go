package postgres

import (
	"context"
	"errors"
	"fmt"

	"horizonx-server/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type JobRepository struct {
	db *pgxpool.Pool
}

func NewJobRepository(db *pgxpool.Pool) domain.JobRepository {
	return &JobRepository{db: db}
}

func (r *JobRepository) List(ctx context.Context) ([]domain.Job, error) {
	query := `
		SELECT
			id,
			server_id,
			application_id,
			deployment_id,
			job_type,
			status,
			output_log,
			queued_at,
			started_at,
			finished_at
		FROM server_jobs
		ORDER BY queued_at DESC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jobs []domain.Job

	for rows.Next() {
		var j domain.Job
		if err := rows.Scan(
			&j.ID,
			&j.ServerID,
			&j.ApplicationID,
			&j.DeploymentID,
			&j.JobType,
			&j.Status,
			&j.OutputLog,
			&j.QueuedAt,
			&j.StartedAt,
			&j.FinishedAt,
		); err != nil {
			return nil, err
		}

		jobs = append(jobs, j)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return jobs, nil
}

func (r *JobRepository) GetByID(ctx context.Context, jobID int64) (*domain.Job, error) {
	query := `
		SELECT
			id,
			server_id,
			application_id,
			deployment_id,
			job_type,
			command_payload,
			status,
			output_log,
			queued_at,
			started_at,
			finished_at
		FROM server_jobs
		WHERE id = $1 LIMIT 1
	`

	var j domain.Job
	err := r.db.QueryRow(ctx, query, jobID).Scan(
		&j.ID,
		&j.ServerID,
		&j.ApplicationID,
		&j.DeploymentID,
		&j.JobType,
		&j.CommandPayload,
		&j.Status,
		&j.OutputLog,
		&j.QueuedAt,
		&j.StartedAt,
		&j.FinishedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrJobNotFound
		}

		return nil, err
	}

	return &j, nil
}

func (r *JobRepository) Create(ctx context.Context, j *domain.Job) (*domain.Job, error) {
	query := `
		INSERT INTO server_jobs
		(server_id, application_id, deployment_id, job_type)
		VALUES ($1, $2, $3, $4)
		RETURNING id, queued_at
	`

	err := r.db.QueryRow(ctx, query,
		j.ServerID,
		j.ApplicationID,
		j.DeploymentID,
		j.JobType,
	).Scan(
		&j.ID,
		&j.QueuedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create job: %w", err)
	}

	return j, nil
}

func (r *JobRepository) Delete(ctx context.Context, jobID int64) error {
	query := `
		DELETE FROM server_jobs
		WHERE id = $1
		RETURNING id
	`

	var deletedID int64
	err := r.db.QueryRow(ctx, query, jobID).Scan(&deletedID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ErrJobNotFound
		}
		return err
	}

	return nil
}

func (r *JobRepository) MarkRunning(ctx context.Context, jobID int64) (*domain.Job, error) {
	query := `
		UPDATE server_jobs
		SET
			status = 'running',
			started_at = NOW()
		WHERE id = $1
		  AND status = 'queued'
		RETURNING
			id,
			server_id,
			application_id,
			deployment_id,
			job_type,
			command_payload,
			status,
			output_log,
			queued_at,
			started_at,
			finished_at
	`

	var job domain.Job
	err := r.db.QueryRow(ctx, query, jobID).Scan(
		&job.ID,
		&job.ServerID,
		&job.ApplicationID,
		&job.DeploymentID,
		&job.JobType,
		&job.CommandPayload,
		&job.Status,
		&job.OutputLog,
		&job.QueuedAt,
		&job.StartedAt,
		&job.FinishedAt,
	)

	if err == pgx.ErrNoRows {
		return r.GetByID(ctx, jobID)
	}
	if err != nil {
		return nil, err
	}

	return &job, nil
}

func (r *JobRepository) MarkFinished(
	ctx context.Context,
	jobID int64,
	status domain.JobStatus,
	outputLog *string,
) (*domain.Job, error) {
	query := `
		UPDATE server_jobs
		SET
			status = $1,
			output_log = $2,
			finished_at = NOW()
		WHERE id = $3
		  AND status = 'running'
		RETURNING
			id,
			server_id,
			application_id,
			deployment_id,
			job_type,
			command_payload,
			status,
			output_log,
			queued_at,
			started_at,
			finished_at
	`

	var job domain.Job
	err := r.db.QueryRow(ctx, query, status, outputLog, jobID).Scan(
		&job.ID,
		&job.ServerID,
		&job.ApplicationID,
		&job.DeploymentID,
		&job.JobType,
		&job.CommandPayload,
		&job.Status,
		&job.OutputLog,
		&job.QueuedAt,
		&job.StartedAt,
		&job.FinishedAt,
	)

	if err == pgx.ErrNoRows {
		return r.GetByID(ctx, jobID)
	}
	if err != nil {
		return nil, err
	}

	return &job, nil
}
