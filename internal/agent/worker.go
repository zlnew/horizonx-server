// Package agent
package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"horizonx-server/internal/agent/executor"
	"horizonx-server/internal/config"
	"horizonx-server/internal/domain"
	"horizonx-server/internal/logger"
)

type JobWorker struct {
	cfg      *config.Config
	log      logger.Logger
	executor *executor.Executor
	client   *http.Client
}

func NewJobWorker(cfg *config.Config, log logger.Logger, workDir string) *JobWorker {
	return &JobWorker{
		cfg:      cfg,
		log:      log,
		executor: executor.NewExecutor(log, workDir),
		client: &http.Client{
			Timeout: 5 * time.Minute,
		},
	}
}

func (w *JobWorker) Initialize() error {
	return w.executor.Initialize()
}

func (w *JobWorker) Start(ctx context.Context) error {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	w.log.Info("job worker started, polling for jobs...")

	for {
		select {
		case <-ctx.Done():
			w.log.Info("job worker stopping...")
			return ctx.Err()

		case <-ticker.C:
			if err := w.pollAndExecuteJobs(ctx); err != nil {
				w.log.Error("failed to poll jobs", "error", err)
			}
		}
	}
}

func (w *JobWorker) pollAndExecuteJobs(ctx context.Context) error {
	jobs, err := w.fetchJobs(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch jobs: %w", err)
	}

	if len(jobs) == 0 {
		return nil
	}

	w.log.Info("received jobs", "count", len(jobs))

	for _, job := range jobs {
		if err := w.processJob(ctx, job); err != nil {
			w.log.Error("failed to process job", "job_id", job.ID, "error", err)
		}
	}

	return nil
}

func (w *JobWorker) fetchJobs(ctx context.Context) ([]domain.Job, error) {
	url := w.cfg.AgentTargetAPIURL + "/agent/jobs"

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+w.cfg.AgentServerID.String()+"."+w.cfg.AgentServerAPIToken)

	resp, err := w.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned status %d", resp.StatusCode)
	}

	var response struct {
		Data []domain.Job `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	var pendingJobs []domain.Job
	for _, job := range response.Data {
		if job.Status == string(domain.JobQueued) {
			pendingJobs = append(pendingJobs, job)
		}
	}

	return pendingJobs, nil
}

func (w *JobWorker) processJob(ctx context.Context, job domain.Job) error {
	w.log.Info("processing job", "id", job.ID, "type", job.JobType)

	if err := w.markJobRunning(ctx, job.ID); err != nil {
		w.log.Error("failed to mark job as running", "job_id", job.ID, "error", err)
		return err
	}

	output, execErr := w.executor.Execute(ctx, &job)

	status := domain.JobSuccess
	if execErr != nil {
		status = domain.JobFailed
		output = fmt.Sprintf("Job execution failed: %v\n\n%s", execErr, output)
		w.log.Error("job execution failed", "job_id", job.ID, "error", execErr)
	} else {
		w.log.Info("job executed successfully", "job_id", job.ID)
	}

	if err := w.markJobFinished(ctx, job.ID, status, output); err != nil {
		w.log.Error("failed to mark job as finished", "job_id", job.ID, "error", err)
		return err
	}

	return execErr
}

func (w *JobWorker) markJobRunning(ctx context.Context, jobID int64) error {
	url := fmt.Sprintf("%s/agent/jobs/%d/start", w.cfg.AgentTargetAPIURL, jobID)

	req, err := http.NewRequestWithContext(ctx, "POST", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+w.cfg.AgentServerID.String()+"."+w.cfg.AgentServerAPIToken)

	resp, err := w.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to mark job running, status: %d", resp.StatusCode)
	}

	return nil
}

func (w *JobWorker) markJobFinished(ctx context.Context, jobID int64, status domain.JobStatus, output string) error {
	url := fmt.Sprintf("%s/agent/jobs/%d/finish", w.cfg.AgentTargetAPIURL, jobID)

	payload := map[string]string{
		"status":     string(status),
		"output_log": output,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(body))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+w.cfg.AgentServerID.String()+"."+w.cfg.AgentServerAPIToken)

	resp, err := w.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to mark job finished, status: %d", resp.StatusCode)
	}

	return nil
}
