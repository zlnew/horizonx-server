// Package executor
package executor

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	"horizonx-server/internal/agent/docker"
	"horizonx-server/internal/agent/git"
	"horizonx-server/internal/domain"
	"horizonx-server/internal/logger"
)

type Executor struct {
	log    logger.Logger
	docker *docker.Manager
	git    *git.Manager
}

func NewExecutor(log logger.Logger, workDir string) *Executor {
	return &Executor{
		log:    log,
		docker: docker.NewManager(log, workDir),
		git:    git.NewManager(log),
	}
}

func (e *Executor) Initialize() error {
	if !e.docker.IsDockerInstalled() {
		return fmt.Errorf("docker is not installed")
	}

	if !e.docker.IsDockerComposeAvailable() {
		return fmt.Errorf("docker compose is not available")
	}

	if !e.git.IsGitInstalled() {
		return fmt.Errorf("git is not installed")
	}

	if err := e.docker.Initialize(); err != nil {
		return err
	}

	e.log.Info("job executor initialized successfully")
	return nil
}

func (e *Executor) Execute(ctx context.Context, job *domain.Job) (string, error) {
	e.log.Info("executing job", "id", job.ID, "type", job.JobType)

	switch job.JobType {
	case domain.JobTypeDeployApp:
		return e.executeDeployApp(ctx, job)
	case domain.JobTypeStartApp:
		return e.executeStartApp(ctx, job)
	case domain.JobTypeStopApp:
		return e.executeStopApp(ctx, job)
	case domain.JobTypeRestartApp:
		return e.executeRestartApp(ctx, job)
	default:
		return "", fmt.Errorf("unknown job type: %s", job.JobType)
	}
}

func (e *Executor) executeDeployApp(ctx context.Context, job *domain.Job) (string, error) {
	var payload domain.DeployAppPayload
	payloadBytes, err := json.Marshal(job.CommandPayload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal payload: %w", err)
	}
	if err := json.Unmarshal(payloadBytes, &payload); err != nil {
		return "", fmt.Errorf("failed to parse deploy payload: %w", err)
	}

	appID := payload.ApplicationID
	var logs strings.Builder
	logs.WriteString(fmt.Sprintf("=== Deploying Application %d ===\n\n", appID))

	if payload.RepoURL != nil && *payload.RepoURL != "" {
		logs.WriteString("Step 1: Fetching source code...\n")

		repoDir := filepath.Join(e.docker.GetAppDir(appID), "source")
		var gitOutput string

		if e.git.IsGitRepo(repoDir) {
			e.log.Debug("repository exists, pulling latest", "app_id", appID)
			gitOutput, err = e.git.Pull(ctx, repoDir, payload.Branch)
		} else {
			e.log.Debug("cloning repository", "app_id", appID, "repo", *payload.RepoURL)
			gitOutput, err = e.git.Clone(ctx, *payload.RepoURL, payload.Branch, repoDir)
		}

		if err != nil {
			logs.WriteString(fmt.Sprintf("❌ Git operation failed: %v\n", err))
			return logs.String(), err
		}

		logs.WriteString(gitOutput)
		logs.WriteString("✓ Source code ready\n\n")

		commitHash, _ := e.git.GetCurrentCommit(ctx, repoDir)
		commitMsg, _ := e.git.GetCommitMessage(ctx, repoDir)
		if commitHash != "" {
			logs.WriteString(fmt.Sprintf("Commit: %s\n", commitHash[:8]))
			logs.WriteString(fmt.Sprintf("Message: %s\n\n", commitMsg))
		}
	}

	// FIX: check if repo has docker-compose.yml, if not exists then return error
	// logs.WriteString("Step 2: Writing docker-compose.yml...\n")
	// if err := e.docker.WriteDockerComposeFile(appID, payload.DockerComposeRaw); err != nil {
	// 	logs.WriteString(fmt.Sprintf("❌ Failed to write compose file: %v\n", err))
	// 	return logs.String(), err
	// }
	// logs.WriteString("✓ Compose file written\n\n")

	if len(payload.EnvVars) > 0 {
		logs.WriteString(fmt.Sprintf("Step 3: Writing environment variables (%d vars)...\n", len(payload.EnvVars)))
		if err := e.docker.WriteEnvFile(appID, payload.EnvVars); err != nil {
			logs.WriteString(fmt.Sprintf("❌ Failed to write env file: %v\n", err))
			return logs.String(), err
		}
		logs.WriteString("✓ Environment configured\n\n")
	}

	logs.WriteString("Step 4: Stopping existing containers...\n")
	stopOutput, err := e.docker.ComposeDown(ctx, appID, false)
	if err != nil {
		e.log.Warn("failed to stop existing containers (may not exist)", "error", err)
	}
	if stopOutput != "" {
		logs.WriteString(stopOutput)
		logs.WriteString("\n")
	}

	logs.WriteString("Step 5: Building and starting containers...\n")
	upOutput, err := e.docker.ComposeUp(ctx, appID, true, true)
	if err != nil {
		logs.WriteString(fmt.Sprintf("❌ Docker compose up failed: %v\n", err))
		logs.WriteString(upOutput)
		return logs.String(), err
	}
	logs.WriteString(upOutput)
	logs.WriteString("\n✓ Containers started\n\n")

	logs.WriteString("Step 6: Verifying deployment...\n")
	psOutput, err := e.docker.ComposePs(ctx, appID)
	if err != nil {
		logs.WriteString(fmt.Sprintf("⚠ Failed to check container status: %v\n", err))
	} else {
		logs.WriteString(psOutput)
		logs.WriteString("\n")
	}

	logs.WriteString("\n=== Deployment Complete ===\n")
	return logs.String(), nil
}

func (e *Executor) executeStartApp(ctx context.Context, job *domain.Job) (string, error) {
	var payload domain.StartAppPayload
	payloadBytes, err := json.Marshal(job.CommandPayload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal payload: %w", err)
	}
	if err := json.Unmarshal(payloadBytes, &payload); err != nil {
		return "", fmt.Errorf("failed to parse start payload: %w", err)
	}

	appID := payload.ApplicationID
	e.log.Info("starting application", "app_id", appID)

	output, err := e.docker.ComposeStart(ctx, appID)
	if err != nil {
		return output, fmt.Errorf("failed to start application: %w", err)
	}

	return fmt.Sprintf("Application %d started successfully\n%s", appID, output), nil
}

func (e *Executor) executeStopApp(ctx context.Context, job *domain.Job) (string, error) {
	var payload domain.StopAppPayload
	payloadBytes, err := json.Marshal(job.CommandPayload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal payload: %w", err)
	}
	if err := json.Unmarshal(payloadBytes, &payload); err != nil {
		return "", fmt.Errorf("failed to parse stop payload: %w", err)
	}

	appID := payload.ApplicationID
	e.log.Info("stopping application", "app_id", appID)

	output, err := e.docker.ComposeStop(ctx, appID)
	if err != nil {
		return output, fmt.Errorf("failed to stop application: %w", err)
	}

	return fmt.Sprintf("Application %d stopped successfully\n%s", appID, output), nil
}

func (e *Executor) executeRestartApp(ctx context.Context, job *domain.Job) (string, error) {
	var payload domain.RestartAppPayload
	payloadBytes, err := json.Marshal(job.CommandPayload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal payload: %w", err)
	}
	if err := json.Unmarshal(payloadBytes, &payload); err != nil {
		return "", fmt.Errorf("failed to parse restart payload: %w", err)
	}

	appID := payload.ApplicationID
	e.log.Info("restarting application", "app_id", appID)

	output, err := e.docker.ComposeRestart(ctx, appID)
	if err != nil {
		return output, fmt.Errorf("failed to restart application: %w", err)
	}

	return fmt.Sprintf("Application %d restarted successfully\n%s", appID, output), nil
}
