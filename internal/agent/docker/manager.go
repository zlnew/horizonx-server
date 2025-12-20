// Package docker
package docker

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"horizonx-server/internal/logger"
)

type Manager struct {
	log     logger.Logger
	workDir string
}

func NewManager(log logger.Logger, workDir string) *Manager {
	return &Manager{
		log:     log,
		workDir: workDir,
	}
}

func (m *Manager) Initialize() error {
	if err := os.Mkdir(m.workDir, 0o755); err != nil {
		return fmt.Errorf("failed to create work directory: %w", err)
	}

	m.log.Info("docker manager initialized", "work_dir", m.workDir)
	return nil
}

func (m *Manager) GetAppDir(appID int64) string {
	return filepath.Join(m.workDir, fmt.Sprintf("app-%d", appID))
}

func (m *Manager) ComposeUp(ctx context.Context, appID int64, detached, build bool) (string, error) {
	appDir := m.GetAppDir(appID)

	args := []string{"compose", "up"}
	if detached {
		args = append(args, "-d")
	}
	if build {
		args = append(args, "--build")
	}

	return m.runDockerCommand(ctx, appDir, args...)
}

func (m *Manager) ComposeDown(ctx context.Context, appID int64, removeVolumes bool) (string, error) {
	appDir := m.GetAppDir(appID)

	args := []string{"compose", "down"}
	if removeVolumes {
		args = append(args, "-v")
	}

	return m.runDockerCommand(ctx, appDir, args...)
}

func (m *Manager) ComposeStop(ctx context.Context, appID int64) (string, error) {
	appDir := m.GetAppDir(appID)
	return m.runDockerCommand(ctx, appDir, "compose", "stop")
}

func (m *Manager) ComposeStart(ctx context.Context, appID int64) (string, error) {
	appDir := m.GetAppDir(appID)
	return m.runDockerCommand(ctx, appDir, "compose", "start")
}

func (m *Manager) ComposeRestart(ctx context.Context, appID int64) (string, error) {
	appDir := m.GetAppDir(appID)
	return m.runDockerCommand(ctx, appDir, "compose", "restart")
}

func (m *Manager) ComposeLogs(ctx context.Context, appID int64, tail int) (string, error) {
	appDir := m.GetAppDir(appID)
	args := []string{"compose", "logs"}
	if tail > 0 {
		args = append(args, "--tail", fmt.Sprintf("%d", tail))
	}
	return m.runDockerCommand(ctx, appDir, args...)
}

func (m *Manager) ComposePs(ctx context.Context, appID int64) (string, error) {
	appDir := m.GetAppDir(appID)
	return m.runDockerCommand(ctx, appDir, "compose", "ps")
}

func (m *Manager) WriteDockerComposeFile(appID int64, content string) error {
	appDir := m.GetAppDir(appID)

	if err := os.MkdirAll(appDir, 0o755); err != nil {
		return fmt.Errorf("failed to create app directory: %w", err)
	}

	composePath := filepath.Join(appDir, "docker-compose.yml")
	if err := os.WriteFile(composePath, []byte(content), 0o644); err != nil {
		return fmt.Errorf("failed to write docker-compose.yml: %w", err)
	}

	m.log.Debug("docker-compose.yml written", "app_id", appID, "path", composePath)
	return nil
}

func (m *Manager) WriteEnvFile(appID int64, envVars map[string]string) error {
	if len(envVars) == 0 {
		return nil
	}

	appDir := m.GetAppDir(appID)
	envPath := filepath.Join(appDir, ".env")

	var buf bytes.Buffer
	for key, value := range envVars {
		escapedValue := strings.ReplaceAll(value, "\n", "\\n")
		escapedValue = strings.ReplaceAll(escapedValue, "\"", "\\\"")
		buf.WriteString(fmt.Sprintf("%s=\"%s\"\n", key, escapedValue))
	}

	if err := os.WriteFile(envPath, buf.Bytes(), 0o600); err != nil {
		return fmt.Errorf("failed to write .env file: %w", err)
	}

	m.log.Debug("environment file written", "app_id", appID, "vars_count", len(envVars))
	return nil
}

func (m *Manager) CleanupAppDir(appID int64) error {
	appDir := m.GetAppDir(appID)

	if err := os.RemoveAll(appDir); err != nil {
		return fmt.Errorf("failed to cleanup app directory: %w", err)
	}

	m.log.Info("app directory cleaned up", "app_id", appID)
	return nil
}

func (m *Manager) runDockerCommand(ctx context.Context, workDir string, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, "docker", args...)
	cmd.Dir = m.workDir

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	m.log.Debug("running docker command", "cmd", cmd.String(), "dir", workDir)

	err := cmd.Run()

	output := stdout.String()
	if stderr.Len() > 0 {
		output += "\n" + stderr.String()
	}

	if err != nil {
		m.log.Error("docker command failed",
			"cmd", cmd.String(),
			"error", err,
			"output", output,
		)
		return output, fmt.Errorf("docker command failed: %w", err)
	}

	m.log.Debug("docker command succeeded", "cmd", cmd.String())
	return output, nil
}

func (m *Manager) IsDockerInstalled() bool {
	cmd := exec.Command("docker", "--version")
	if err := cmd.Run(); err != nil {
		m.log.Warn("docker not found in PATH")
		return false
	}

	return true
}

func (m *Manager) IsDockerComposeAvailable() bool {
	cmd := exec.Command("docker", "compose", "version")
	if err := cmd.Run(); err != nil {
		m.log.Warn("docker compose not available")
		return false
	}

	return true
}
