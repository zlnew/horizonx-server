package domain

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrApplicationNotFound  = errors.New("application not found")
	ErrInvalidDockerCompose = errors.New("invalid docker compose configuration")
)

type ApplicationStatus string

const (
	AppStatusDeploying  ApplicationStatus = "deploying"
	AppStatusStarting   ApplicationStatus = "starting"
	AppStatusStopping   ApplicationStatus = "stopping"
	AppStatusRestarting ApplicationStatus = "restarting"
	AppStatusRunning    ApplicationStatus = "running"
	AppStatusStopped    ApplicationStatus = "stopped"
	AppStatusFailed     ApplicationStatus = "failed"
)

type Application struct {
	ID               int64             `json:"id"`
	ServerID         uuid.UUID         `json:"server_id"`
	Name             string            `json:"name"`
	RepoURL          string            `json:"repo_url,omitempty"`
	Branch           string            `json:"branch"`
	Status           ApplicationStatus `json:"status"`
	LastDeploymentAt *time.Time        `json:"last_deployment_at,omitempty"`
	CreatedAt        time.Time         `json:"created_at"`
	UpdatedAt        time.Time         `json:"updated_at"`

	EnvVars *[]EnvironmentVariable `json:"env_vars,omitempty"`
}

type EnvironmentVariable struct {
	ID            int64     `json:"id"`
	ApplicationID int64     `json:"application_id"`
	Key           string    `json:"key"`
	Value         string    `json:"value"`
	IsPreview     bool      `json:"is_preview"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type ApplicationCreateRequest struct {
	ServerID uuid.UUID `json:"server_id" validate:"required"`
	Name     string    `json:"name" validate:"required,min=3,max=100"`
	RepoURL  string    `json:"repo_url" validate:"required"`
	Branch   string    `json:"branch" validate:"required"`

	EnvVars []EnvironmentVariableRequest `json:"env_vars" validate:"omitempty,dive"`
}

type ApplicationUpdateRequest struct {
	Name    string `json:"name" validate:"required,min=3,max=100"`
	RepoURL string `json:"repo_url" validate:"required"`
	Branch  string `json:"branch" validate:"required"`

	EnvVars []EnvironmentVariableRequest `json:"env_vars" validate:"omitempty,dive"`
}

type EnvironmentVariableRequest struct {
	Key       string `json:"key" validate:"required"`
	Value     string `json:"value" validate:"required"`
	IsPreview bool   `json:"is_preview"`
}

type ApplicationRepository interface {
	// Applications
	List(ctx context.Context, serverID uuid.UUID) ([]Application, error)
	GetByID(ctx context.Context, appID int64) (*Application, error)
	Create(ctx context.Context, app *Application) (*Application, error)
	Update(ctx context.Context, app *Application, appID int64) error
	Delete(ctx context.Context, appID int64) error
	UpdateStatus(ctx context.Context, appID int64, status ApplicationStatus) error
	UpdateLastDeployment(ctx context.Context, appID int64) error

	// Environment Variables
	SyncEnvVars(ctx context.Context, appID int64, envVars []EnvironmentVariable) error
	ListEnvVars(ctx context.Context, appID int64) ([]EnvironmentVariable, error)
	CreateEnvVar(ctx context.Context, env *EnvironmentVariable) error
	UpdateEnvVar(ctx context.Context, env *EnvironmentVariable) error
	DeleteEnvVar(ctx context.Context, appID int64, key string) error
}

type ApplicationService interface {
	// Applications
	List(ctx context.Context, serverID uuid.UUID) ([]Application, error)
	GetByID(ctx context.Context, appID int64) (*Application, error)
	Create(ctx context.Context, req ApplicationCreateRequest) (*Application, error)
	Update(ctx context.Context, req ApplicationUpdateRequest, appID int64) error
	Delete(ctx context.Context, appID int64) error
	UpdateStatus(ctx context.Context, appID int64, status ApplicationStatus) error
	UpdateLastDeployment(ctx context.Context, appID int64) error

	// Actions
	Deploy(ctx context.Context, appID int64, deployedBy int64) (*Deployment, error)
	Start(ctx context.Context, appID int64) error
	Stop(ctx context.Context, appID int64) error
	Restart(ctx context.Context, appID int64) error

	// Environment Variables
	ListEnvVars(ctx context.Context, appID int64) ([]EnvironmentVariable, error)
	AddEnvVar(ctx context.Context, appID int64, req EnvironmentVariableRequest) error
	UpdateEnvVar(ctx context.Context, appID int64, key string, req EnvironmentVariableRequest) error
	DeleteEnvVar(ctx context.Context, appID int64, key string) error
}
