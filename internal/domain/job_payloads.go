package domain

type DeployAppPayload struct {
	ApplicationID    int64             `json:"application_id"`
	RepoURL          *string           `json:"repo_url"`
	Branch           string            `json:"branch"`
	DockerComposeRaw string            `json:"docker_compose_raw"`
	EnvVars          map[string]string `json:"env_vars,omitempty"`
	Volumes          []VolumeMount     `json:"volumes,omitempty"`
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

type VolumeMount struct {
	HostPath      string `json:"host_path"`
	ContainerPath string `json:"container_path"`
	Mode          string `json:"mode"`
}

const (
	JobTypeDeployApp  = "deploy_app"
	JobTypeStartApp   = "start_app"
	JobTypeStopApp    = "stop_app"
	JobTypeRestartApp = "restart_app"
	JobTypeGetLogs    = "get_logs"
)
