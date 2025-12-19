package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"horizonx-server/internal/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ApplicationRepository struct {
	db *pgxpool.Pool
}

func NewApplicationRepository(db *pgxpool.Pool) domain.ApplicationRepository {
	return &ApplicationRepository{db: db}
}

// ============================================================================
// APPLICATIONS
// ============================================================================

func (r *ApplicationRepository) List(ctx context.Context, serverID uuid.UUID) ([]domain.Application, error) {
	query := `
		SELECT id, server_id, name, repo_url, branch, docker_compose_raw, 
		       status, last_deployment_at, created_at, updated_at
		FROM applications
		WHERE server_id = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query, serverID)
	if err != nil {
		return nil, fmt.Errorf("failed to query applications: %w", err)
	}
	defer rows.Close()

	var apps []domain.Application
	for rows.Next() {
		var app domain.Application
		err := rows.Scan(
			&app.ID,
			&app.ServerID,
			&app.Name,
			&app.RepoURL,
			&app.Branch,
			&app.DockerComposeRaw,
			&app.Status,
			&app.LastDeploymentAt,
			&app.CreatedAt,
			&app.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan application: %w", err)
		}
		apps = append(apps, app)
	}

	return apps, nil
}

func (r *ApplicationRepository) GetByID(ctx context.Context, appID int64) (*domain.Application, error) {
	query := `
		SELECT id, server_id, name, repo_url, branch, docker_compose_raw, 
		       status, last_deployment_at, created_at, updated_at
		FROM applications
		WHERE id = $1 AND deleted_at IS NULL
	`

	var app domain.Application
	err := r.db.QueryRow(ctx, query, appID).Scan(
		&app.ID,
		&app.ServerID,
		&app.Name,
		&app.RepoURL,
		&app.Branch,
		&app.DockerComposeRaw,
		&app.Status,
		&app.LastDeploymentAt,
		&app.CreatedAt,
		&app.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrApplicationNotFound
		}
		return nil, fmt.Errorf("failed to get application: %w", err)
	}

	return &app, nil
}

func (r *ApplicationRepository) Create(ctx context.Context, app *domain.Application) (*domain.Application, error) {
	query := `
		INSERT INTO applications (server_id, name, repo_url, branch, docker_compose_raw, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at
	`

	now := time.Now().UTC()
	err := r.db.QueryRow(
		ctx, query,
		app.ServerID,
		app.Name,
		app.RepoURL,
		app.Branch,
		app.DockerComposeRaw,
		domain.AppStatusStopped,
		now,
		now,
	).Scan(&app.ID, &app.CreatedAt, &app.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create application: %w", err)
	}

	return app, nil
}

func (r *ApplicationRepository) Update(ctx context.Context, app *domain.Application, appID int64) error {
	query := `
		UPDATE applications
		SET name = $1, repo_url = $2, branch = $3, docker_compose_raw = $4, updated_at = $5
		WHERE id = $6 AND deleted_at IS NULL
	`

	now := time.Now().UTC()
	ct, err := r.db.Exec(ctx, query,
		app.Name,
		app.RepoURL,
		app.Branch,
		app.DockerComposeRaw,
		now,
		appID,
	)
	if err != nil {
		return fmt.Errorf("failed to update application: %w", err)
	}

	if ct.RowsAffected() == 0 {
		return domain.ErrApplicationNotFound
	}

	return nil
}

func (r *ApplicationRepository) Delete(ctx context.Context, appID int64) error {
	query := `UPDATE applications SET deleted_at = $1 WHERE id = $2 AND deleted_at IS NULL`

	ct, err := r.db.Exec(ctx, query, time.Now().UTC(), appID)
	if err != nil {
		return fmt.Errorf("failed to delete application: %w", err)
	}

	if ct.RowsAffected() == 0 {
		return domain.ErrApplicationNotFound
	}

	return nil
}

func (r *ApplicationRepository) UpdateStatus(ctx context.Context, appID int64, status domain.ApplicationStatus) error {
	query := `
		UPDATE applications 
		SET status = $1, updated_at = $2
		WHERE id = $3 AND deleted_at IS NULL
	`

	ct, err := r.db.Exec(ctx, query, status, time.Now().UTC(), appID)
	if err != nil {
		return fmt.Errorf("failed to update application status: %w", err)
	}

	if ct.RowsAffected() == 0 {
		return domain.ErrApplicationNotFound
	}

	return nil
}

// ============================================================================
// ENVIRONMENT VARIABLES
// ============================================================================

func (r *ApplicationRepository) ListEnvVars(ctx context.Context, appID int64) ([]domain.EnvironmentVariable, error) {
	query := `
		SELECT id, application_id, key, value, is_preview, created_at, updated_at
		FROM environment_variables
		WHERE application_id = $1
		ORDER BY key ASC
	`

	rows, err := r.db.Query(ctx, query, appID)
	if err != nil {
		return nil, fmt.Errorf("failed to query env vars: %w", err)
	}
	defer rows.Close()

	var envVars []domain.EnvironmentVariable
	for rows.Next() {
		var env domain.EnvironmentVariable
		err := rows.Scan(
			&env.ID,
			&env.ApplicationID,
			&env.Key,
			&env.Value,
			&env.IsPreview,
			&env.CreatedAt,
			&env.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan env var: %w", err)
		}
		envVars = append(envVars, env)
	}

	return envVars, nil
}

func (r *ApplicationRepository) CreateEnvVar(ctx context.Context, env *domain.EnvironmentVariable) error {
	query := `
		INSERT INTO environment_variables (application_id, key, value, is_preview, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (application_id, key) 
		DO UPDATE SET value = EXCLUDED.value, is_preview = EXCLUDED.is_preview, updated_at = EXCLUDED.updated_at
		RETURNING id, created_at, updated_at
	`

	now := time.Now().UTC()
	err := r.db.QueryRow(ctx, query,
		env.ApplicationID,
		env.Key,
		env.Value,
		env.IsPreview,
		now,
		now,
	).Scan(&env.ID, &env.CreatedAt, &env.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create env var: %w", err)
	}

	return nil
}

func (r *ApplicationRepository) UpdateEnvVar(ctx context.Context, env *domain.EnvironmentVariable) error {
	query := `
		UPDATE environment_variables
		SET value = $1, is_preview = $2, updated_at = $3
		WHERE application_id = $4 AND key = $5
	`

	ct, err := r.db.Exec(ctx, query,
		env.Value,
		env.IsPreview,
		time.Now().UTC(),
		env.ApplicationID,
		env.Key,
	)
	if err != nil {
		return fmt.Errorf("failed to update env var: %w", err)
	}

	if ct.RowsAffected() == 0 {
		return fmt.Errorf("env var not found")
	}

	return nil
}

func (r *ApplicationRepository) DeleteEnvVar(ctx context.Context, appID int64, key string) error {
	query := `DELETE FROM environment_variables WHERE application_id = $1 AND key = $2`

	ct, err := r.db.Exec(ctx, query, appID, key)
	if err != nil {
		return fmt.Errorf("failed to delete env var: %w", err)
	}

	if ct.RowsAffected() == 0 {
		return fmt.Errorf("env var not found")
	}

	return nil
}

// ============================================================================
// VOLUMES
// ============================================================================

func (r *ApplicationRepository) ListVolumes(ctx context.Context, appID int64) ([]domain.Volume, error) {
	query := `
		SELECT id, application_id, host_path, container_path, mode, created_at
		FROM volumes
		WHERE application_id = $1
		ORDER BY created_at ASC
	`

	rows, err := r.db.Query(ctx, query, appID)
	if err != nil {
		return nil, fmt.Errorf("failed to query volumes: %w", err)
	}
	defer rows.Close()

	var volumes []domain.Volume
	for rows.Next() {
		var vol domain.Volume
		err := rows.Scan(
			&vol.ID,
			&vol.ApplicationID,
			&vol.HostPath,
			&vol.ContainerPath,
			&vol.Mode,
			&vol.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan volume: %w", err)
		}
		volumes = append(volumes, vol)
	}

	return volumes, nil
}

func (r *ApplicationRepository) CreateVolume(ctx context.Context, vol *domain.Volume) error {
	query := `
		INSERT INTO volumes (application_id, host_path, container_path, mode, created_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at
	`

	now := time.Now().UTC()
	err := r.db.QueryRow(ctx, query,
		vol.ApplicationID,
		vol.HostPath,
		vol.ContainerPath,
		vol.Mode,
		now,
	).Scan(&vol.ID, &vol.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create volume: %w", err)
	}

	return nil
}

func (r *ApplicationRepository) DeleteVolume(ctx context.Context, volumeID int64) error {
	query := `DELETE FROM volumes WHERE id = $1`

	ct, err := r.db.Exec(ctx, query, volumeID)
	if err != nil {
		return fmt.Errorf("failed to delete volume: %w", err)
	}

	if ct.RowsAffected() == 0 {
		return fmt.Errorf("volume not found")
	}

	return nil
}
