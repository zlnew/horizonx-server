// Package git
package git

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"horizonx-server/internal/logger"
)

type Manager struct {
	log logger.Logger
}

func NewManager(log logger.Logger) *Manager {
	return &Manager{log: log}
}

func (m *Manager) Clone(ctx context.Context, repoURL, branch, destDir string) (string, error) {
	if _, err := os.Stat(destDir); err == nil {
		m.log.Debug("repository directory exists, pulling instead", "dir", destDir)
		return m.Pull(ctx, destDir, branch)
	}

	parent := filepath.Dir(destDir)
	if err := os.MkdirAll(parent, 0o755); err != nil {
		return "", fmt.Errorf("failed to create parent directory: %w", err)
	}

	args := []string{"clone", "--branch", branch, "--depth", "1", repoURL, destDir}
	output, err := m.runGitCommand(ctx, parent, args...)
	if err != nil {
		return output, fmt.Errorf("git clone failed: %w", err)
	}

	m.log.Info("repository cloned", "repo", repoURL, "branch", branch, "dest", destDir)
	return output, nil
}

func (m *Manager) Pull(ctx context.Context, repoDir, branch string) (string, error) {
	if _, err := m.runGitCommand(ctx, repoDir, "checkout", branch); err != nil {
		m.log.Warn("failed to checkout branch, continuing", "branch", branch, "error", err)
	}

	output, err := m.runGitCommand(ctx, repoDir, "pull", "origin", branch)
	if err != nil {
		return output, fmt.Errorf("git pull failed: %w", err)
	}

	m.log.Info("repository updated", "branch", branch, "dir", repoDir)
	return output, nil
}

func (m *Manager) GetCurrentCommit(ctx context.Context, repoDir string) (string, error) {
	output, err := m.runGitCommand(ctx, repoDir, "rev-parse", "HEAD")
	if err != nil {
		return "", fmt.Errorf("failed to get commit hash: %w", err)
	}
	return output, nil
}

func (m *Manager) GetCommitMessage(ctx context.Context, repoDir string) (string, error) {
	output, err := m.runGitCommand(ctx, repoDir, "log", "-1", "--pretty=%B")
	if err != nil {
		return "", fmt.Errorf("failed to get commit message: %w", err)
	}
	return output, nil
}

func (m *Manager) IsGitRepo(dir string) bool {
	gitDir := filepath.Join(dir, ".git")
	_, err := os.Stat(gitDir)
	return err == nil
}

func (m *Manager) runGitCommand(ctx context.Context, workDir string, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, "git", args...)
	cmd.Dir = workDir

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	m.log.Debug("running git command", "cmd", cmd.String(), "dir", workDir)

	err := cmd.Run()

	output := stdout.String()
	if stderr.Len() > 0 {
		output += "\n" + stderr.String()
	}

	if err != nil {
		m.log.Error("git command failed",
			"cmd", cmd.String(),
			"error", err,
			"output", output,
		)
		return output, fmt.Errorf("git command failed: %w", err)
	}

	m.log.Debug("git command succeeded", "cmd", cmd.String())
	return output, nil
}

func (m *Manager) IsGitInstalled() bool {
	cmd := exec.Command("git", "--version")
	if err := cmd.Run(); err != nil {
		m.log.Warn("git not found in PATH")
		return false
	}
	return true
}
