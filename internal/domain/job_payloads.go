package domain

import "github.com/google/uuid"

type DeployAppPayload struct {
	ApplicationID int64             `json:"application_id"`
	DeploymentID  int64             `json:"deployment_id"`
	RepoURL       string            `json:"repo_url"`
	Branch        string            `json:"branch"`
	EnvVars       map[string]string `json:"env_vars,omitempty"`
}

type StartAppPayload struct {
	ApplicationID int64 `json:"application_id"`
}

type StopAppPayload struct {
	ApplicationID int64 `json:"application_id"`
}

type RestartAppPayload struct {
	ApplicationID int64 `json:"application_id"`
}

type AppHealthCheckPayload struct {
	ServerID        uuid.UUID `json:"server_id"`
	ApplicationsIDs []int64   `json:"application_ids"`
}
