package sevalla

import "time"

// Region represents a deployment region
type Region string

// Available deployment regions
const (
	RegionUSCentral  Region = "us-central1"
	RegionUSEast     Region = "us-east1"
	RegionEuropeWest Region = "europe-west1"
	RegionAsiaSouth  Region = "asia-south1"
)

// Plan represents a pricing/compute plan tier
type Plan string

// Available plan tiers
const (
	PlanHobby      Plan = "hobby"
	PlanStarter    Plan = "starter"
	PlanPro        Plan = "pro"
	PlanBusiness   Plan = "business"
	PlanEnterprise Plan = "enterprise"
)

// ApplicationState represents the state of an application
type ApplicationState string

// Application states
const (
	StateRunning   ApplicationState = "running"
	StateDeploying ApplicationState = "deploying"
	StateFailed    ApplicationState = "failed"
	StateStopped   ApplicationState = "stopped"
	StatePending   ApplicationState = "pending"
	StateBuilding  ApplicationState = "building"
)

// Engine represents a database engine type
type Engine string

// Available database engines
const (
	EnginePostgreSQL Engine = "postgresql"
	EngineMySQL      Engine = "mysql"
	EngineMongoDB    Engine = "mongodb"
	EngineRedis      Engine = "redis"
)

// Status represents a deployment status
type Status string

// Deployment statuses
const (
	StatusQueued    Status = "queued"
	StatusBuilding  Status = "building"
	StatusDeploying Status = "deploying"
	StatusSuccess   Status = "success"
	StatusFailed    Status = "failed"
	StatusCancelled Status = "cancelled"
)

// Application represents a Sevalla application
type Application struct {
	ID               string                 `json:"id"`
	Name             string                 `json:"name"`
	State            ApplicationState       `json:"state"`
	RepositoryURL    string                 `json:"repository_url"`
	Branch           string                 `json:"branch"`
	Region           Region                 `json:"location"`
	Plan             Plan                   `json:"pod_size"`
	Replicas         int                    `json:"replicas"`
	EnvironmentVars  map[string]string      `json:"environment_variables,omitempty"`
	BuildCommand     string                 `json:"build_command,omitempty"`
	StartCommand     string                 `json:"start_command,omitempty"`
	Port             int                    `json:"port,omitempty"`
	URL              string                 `json:"url,omitempty"`
	CustomDomains    []string               `json:"custom_domains,omitempty"`
	AutoDeploy       bool                   `json:"auto_deploy"`
	CDNEnabled       bool                   `json:"cdn_enabled"`
	SSLEnabled       bool                   `json:"ssl_enabled"`
	CreatedAt        time.Time              `json:"created_at"`
	UpdatedAt        time.Time              `json:"updated_at"`
	LastDeploymentID string                 `json:"last_deployment_id,omitempty"`
	Metadata         map[string]interface{} `json:"metadata,omitempty"`
}

// Database represents a Sevalla database
type Database struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Type        Engine                 `json:"type"`
	Version     string                 `json:"version"`
	Region      Region                 `json:"location"`
	Size        string                 `json:"size"`
	Storage     int                    `json:"storage_gb"`
	Backups     bool                   `json:"backups_enabled"`
	PublicURL   string                 `json:"public_url,omitempty"`
	InternalURL string                 `json:"internal_url,omitempty"`
	Username    string                 `json:"username,omitempty"`
	Password    string                 `json:"password,omitempty"`
	SSLEnabled  bool                   `json:"ssl_enabled"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// StaticSite represents a Sevalla static site
type StaticSite struct {
	ID               string            `json:"id"`
	Name             string            `json:"name"`
	State            ApplicationState  `json:"state"`
	RepositoryURL    string            `json:"repository_url"`
	Branch           string            `json:"branch"`
	Region           Region            `json:"location"`
	BuildCommand     string            `json:"build_command,omitempty"`
	OutputDirectory  string            `json:"output_directory,omitempty"`
	EnvironmentVars  map[string]string `json:"environment_variables,omitempty"`
	URL              string            `json:"url,omitempty"`
	CustomDomains    []string          `json:"custom_domains,omitempty"`
	AutoDeploy       bool              `json:"auto_deploy"`
	CDNEnabled       bool              `json:"cdn_enabled"`
	SSLEnabled       bool              `json:"ssl_enabled"`
	CreatedAt        time.Time         `json:"created_at"`
	UpdatedAt        time.Time         `json:"updated_at"`
	LastDeploymentID string            `json:"last_deployment_id,omitempty"`
}

// Deployment represents a deployment
type Deployment struct {
	ID            string     `json:"id"`
	ApplicationID string     `json:"application_id,omitempty"`
	StaticSiteID  string     `json:"static_site_id,omitempty"`
	State         Status     `json:"state"`
	CommitSHA     string     `json:"commit_sha,omitempty"`
	CommitMessage string     `json:"commit_message,omitempty"`
	Branch        string     `json:"branch,omitempty"`
	BuildLogs     string     `json:"build_logs,omitempty"`
	RuntimeLogs   string     `json:"runtime_logs,omitempty"`
	ErrorMessage  string     `json:"error_message,omitempty"`
	StartedAt     time.Time  `json:"started_at"`
	CompletedAt   *time.Time `json:"completed_at,omitempty"`
	Duration      int        `json:"duration_seconds,omitempty"`
}

// Pipeline represents a CI/CD pipeline
type Pipeline struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Enabled     bool                   `json:"enabled"`
	Trigger     string                 `json:"trigger"`
	Branch      string                 `json:"branch"`
	Steps       []PipelineStep         `json:"steps"`
	Environment map[string]string      `json:"environment,omitempty"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// PipelineStep represents a step in a pipeline
type PipelineStep struct {
	Name      string   `json:"name"`
	Command   string   `json:"command"`
	Image     string   `json:"image,omitempty"`
	Timeout   int      `json:"timeout_seconds,omitempty"`
	Retries   int      `json:"retries,omitempty"`
	DependsOn []string `json:"depends_on,omitempty"`
}

// PipelineRun represents an execution of a pipeline
type PipelineRun struct {
	ID          string            `json:"id"`
	PipelineID  string            `json:"pipeline_id"`
	State       Status            `json:"state"`
	Branch      string            `json:"branch"`
	CommitSHA   string            `json:"commit_sha,omitempty"`
	StartedAt   time.Time         `json:"started_at"`
	CompletedAt *time.Time        `json:"completed_at,omitempty"`
	Duration    int               `json:"duration_seconds,omitempty"`
	Steps       []PipelineRunStep `json:"steps,omitempty"`
}

// PipelineRunStep represents the execution of a pipeline step
type PipelineRunStep struct {
	Name         string     `json:"name"`
	State        Status     `json:"state"`
	Output       string     `json:"output,omitempty"`
	ErrorMessage string     `json:"error_message,omitempty"`
	StartedAt    time.Time  `json:"started_at"`
	CompletedAt  *time.Time `json:"completed_at,omitempty"`
}

// ListOptions represents options for listing resources
type ListOptions struct {
	Page    int    `url:"page,omitempty"`
	PerPage int    `url:"per_page,omitempty"`
	Sort    string `url:"sort,omitempty"`
	Order   string `url:"order,omitempty"`
}

// Backup represents a database backup
type Backup struct {
	ID         string    `json:"id"`
	DatabaseID string    `json:"database_id"`
	Type       string    `json:"type"`
	Size       int64     `json:"size_bytes"`
	Status     string    `json:"status"`
	URL        string    `json:"download_url,omitempty"`
	ExpiresAt  time.Time `json:"expires_at,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}

// Usage represents resource usage metrics
type Usage struct {
	ApplicationID string    `json:"application_id,omitempty"`
	DatabaseID    string    `json:"database_id,omitempty"`
	Period        string    `json:"period"`
	CPUUsage      float64   `json:"cpu_usage"`
	MemoryUsage   float64   `json:"memory_usage"`
	StorageUsage  float64   `json:"storage_usage"`
	BandwidthUsed int64     `json:"bandwidth_bytes"`
	RequestCount  int64     `json:"request_count"`
	StartTime     time.Time `json:"start_time"`
	EndTime       time.Time `json:"end_time"`
}
