package sevalla // import "github.com/juststeveking/sevalla-go"

Package sevalla provides a Go client library for the Sevalla API

CONSTANTS

const (
	// BaseURL is the default base URL for the Sevalla API
	BaseURL = "https://api.sevalla.com/v2"

	// DefaultTimeout is the default timeout for HTTP requests
	DefaultTimeout = 30 * time.Second

	// Version is the SDK version
	Version = "0.1.0"

	// UserAgent is the default user agent
	UserAgent = "sevalla-go/" + Version
)

FUNCTIONS

func Bool(v bool) *bool
    Bool is a helper function that allocates a new bool value

func CheckResponse(r *http.Response) error
    CheckResponse checks the API response for errors

func Int(v int) *int
    Int is a helper function that allocates a new int value

func IsBadRequest(err error) bool
    IsBadRequest returns true if the error is a 400 Bad Request

func IsClientError(err error) bool
    IsClientError returns true if the error is a 4xx client error

func IsConflict(err error) bool
    IsConflict returns true if the error is a 409 Conflict

func IsForbidden(err error) bool
    IsForbidden returns true if the error is a 403 Forbidden

func IsNotFound(err error) bool
    IsNotFound returns true if the error is a 404 Not Found

func IsRateLimited(err error) bool
    IsRateLimited returns true if the error is a 429 Too Many Requests

func IsServerError(err error) bool
    IsServerError returns true if the error is a 5xx server error

func IsUnauthorized(err error) bool
    IsUnauthorized returns true if the error is a 401 Unauthorized

func IsUnprocessableEntity(err error) bool
    IsUnprocessableEntity returns true if the error is a 422 Unprocessable
    Entity

func String(v string) *string
    String is a helper function that allocates a new string value


TYPES

type AddDomainRequest struct {
	Domain string `json:"domain"`
}
    AddDomainRequest represents a request to add a custom domain

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
    Application represents a Sevalla application

type ApplicationState string
    ApplicationState represents the state of an application

const (
	StateRunning   ApplicationState = "running"
	StateDeploying ApplicationState = "deploying"
	StateFailed    ApplicationState = "failed"
	StateStopped   ApplicationState = "stopped"
	StatePending   ApplicationState = "pending"
	StateBuilding  ApplicationState = "building"
)
    Application states

type ApplicationsService struct {
	// Has unexported fields.
}
    ApplicationsService handles communication with the application-related
    methods of the Sevalla API.

func (s *ApplicationsService) AddCustomDomain(ctx context.Context, id string, domain string) (*Response, error)
    AddCustomDomain adds a custom domain to an application

func (s *ApplicationsService) CancelDeployment(ctx context.Context, appID, deploymentID string) (*Response, error)
    CancelDeployment cancels a deployment

func (s *ApplicationsService) Create(ctx context.Context, createReq *CreateApplicationRequest) (*Application, *Response, error)
    Create creates a new application

func (s *ApplicationsService) Delete(ctx context.Context, id string) (*Response, error)
    Delete deletes an application

func (s *ApplicationsService) Deploy(ctx context.Context, id string) (*Deployment, *Response, error)
    Deploy triggers a new deployment for an application

func (s *ApplicationsService) Get(ctx context.Context, id string) (*Application, *Response, error)
    Get returns a single application by ID

func (s *ApplicationsService) GetDeployment(ctx context.Context, appID, deploymentID string) (*Deployment, *Response, error)
    GetDeployment gets a specific deployment for an application

func (s *ApplicationsService) GetEnvironmentVariables(ctx context.Context, id string) (map[string]string, *Response, error)
    GetEnvironmentVariables gets environment variables for an application

func (s *ApplicationsService) GetLogs(ctx context.Context, id string, lines int) (string, *Response, error)
    GetLogs retrieves application logs

func (s *ApplicationsService) GetUsage(ctx context.Context, id string, period string) (*Usage, *Response, error)
    GetUsage retrieves usage metrics for an application

func (s *ApplicationsService) List(ctx context.Context, opts *ListOptions) ([]*Application, *Response, error)
    List returns all applications

func (s *ApplicationsService) ListDeployments(ctx context.Context, id string, opts *ListOptions) ([]*Deployment, *Response, error)
    ListDeployments lists all deployments for an application

func (s *ApplicationsService) RemoveCustomDomain(ctx context.Context, id string, domain string) (*Response, error)
    RemoveCustomDomain removes a custom domain from an application

func (s *ApplicationsService) Restart(ctx context.Context, id string) (*Response, error)
    Restart restarts an application

func (s *ApplicationsService) Rollback(ctx context.Context, appID, deploymentID string) (*Deployment, *Response, error)
    Rollback rolls back to a previous deployment

func (s *ApplicationsService) Scale(ctx context.Context, id string, scaleReq *ScaleApplicationRequest) (*Application, *Response, error)
    Scale scales an application's resources

func (s *ApplicationsService) SetEnvironmentVariables(ctx context.Context, id string, vars map[string]string) (*Response, error)
    SetEnvironmentVariables sets environment variables for an application

func (s *ApplicationsService) Start(ctx context.Context, id string) (*Response, error)
    Start starts a stopped application

func (s *ApplicationsService) Stop(ctx context.Context, id string) (*Response, error)
    Stop stops an application

func (s *ApplicationsService) Update(ctx context.Context, id string, updateReq *UpdateApplicationRequest) (*Application, *Response, error)
    Update updates an existing application

func (s *ApplicationsService) UpdateCDNSettings(ctx context.Context, id string, enabled bool) (*Response, error)
    UpdateCDNSettings updates CDN settings for an application

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
    Backup represents a database backup

type CDNSettingsRequest struct {
	Enabled bool `json:"enabled"`
}
    CDNSettingsRequest represents CDN configuration

type Client struct {

	// Services
	Applications *ApplicationsService
	Databases    *DatabasesService
	StaticSites  *StaticSitesService
	Deployments  *DeploymentsService
	Pipelines    *PipelinesService
	// Has unexported fields.
}
    Client manages communication with the Sevalla API

func NewClient(opts ...ClientOption) *Client
    NewClient creates a new Sevalla API client

func (c *Client) Do(req *http.Request, v interface{}) (*Response, error)
    Do executes an API request and returns the response

func (c *Client) NewRequest(ctx context.Context, method, urlStr string, body interface{}) (*http.Request, error)
    NewRequest creates an API request

func (c *Client) NewRequestWithQuery(ctx context.Context, method, urlStr string, opts interface{}) (*http.Request, error)
    NewRequestWithQuery creates an API request with query parameters

type ClientOption func(*Client)
    ClientOption is a function that configures a Client

func WithAPIKey(key string) ClientOption
    WithAPIKey sets the API key for authentication

func WithBaseURL(baseURL string) ClientOption
    WithBaseURL sets a custom base URL for the API

func WithHTTPClient(client *http.Client) ClientOption
    WithHTTPClient sets a custom HTTP client

func WithUserAgent(ua string) ClientOption
    WithUserAgent sets a custom user agent

type CreateApplicationRequest struct {
	Name            string            `json:"name"`
	RepositoryURL   string            `json:"repository_url"`
	Branch          string            `json:"branch,omitempty"`
	Region          Region            `json:"location,omitempty"`
	Plan            Plan              `json:"pod_size,omitempty"`
	Replicas        int               `json:"replicas,omitempty"`
	EnvironmentVars map[string]string `json:"environment_variables,omitempty"`
	BuildCommand    string            `json:"build_command,omitempty"`
	StartCommand    string            `json:"start_command,omitempty"`
	Port            int               `json:"port,omitempty"`
	AutoDeploy      bool              `json:"auto_deploy,omitempty"`
	CDNEnabled      bool              `json:"cdn_enabled,omitempty"`
	SSLEnabled      bool              `json:"ssl_enabled,omitempty"`
}
    CreateApplicationRequest represents a request to create a new application

type CreateBackupRequest struct {
	Type        string `json:"type,omitempty"` // "manual" or "automatic"
	Description string `json:"description,omitempty"`
}
    CreateBackupRequest represents a request to create a database backup

type CreateDatabaseRequest struct {
	Name       string `json:"name"`
	Type       Engine `json:"type"`
	Version    string `json:"version,omitempty"`
	Region     Region `json:"location,omitempty"`
	Size       string `json:"size,omitempty"`
	Storage    int    `json:"storage_gb,omitempty"`
	Backups    bool   `json:"backups_enabled,omitempty"`
	SSLEnabled bool   `json:"ssl_enabled,omitempty"`
}
    CreateDatabaseRequest represents a request to create a new database

type CreatePipelineRequest struct {
	Name        string                 `json:"name"`
	Enabled     bool                   `json:"enabled"`
	Trigger     string                 `json:"trigger"`
	Branch      string                 `json:"branch"`
	Steps       []PipelineStep         `json:"steps"`
	Environment map[string]string      `json:"environment,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}
    CreatePipelineRequest represents a request to create a pipeline

type CreateStaticSiteRequest struct {
	Name            string            `json:"name"`
	RepositoryURL   string            `json:"repository_url"`
	Branch          string            `json:"branch,omitempty"`
	Region          Region            `json:"location,omitempty"`
	BuildCommand    string            `json:"build_command,omitempty"`
	OutputDirectory string            `json:"output_directory,omitempty"`
	EnvironmentVars map[string]string `json:"environment_variables,omitempty"`
	AutoDeploy      bool              `json:"auto_deploy,omitempty"`
	CDNEnabled      bool              `json:"cdn_enabled,omitempty"`
	SSLEnabled      bool              `json:"ssl_enabled,omitempty"`
}
    CreateStaticSiteRequest represents a request to create a new static site

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
    Database represents a Sevalla database

type DatabasesService struct {
	// Has unexported fields.
}
    DatabasesService handles communication with the database-related methods of
    the Sevalla API.

func (s *DatabasesService) Create(ctx context.Context, createReq *CreateDatabaseRequest) (*Database, *Response, error)
    Create creates a new database

func (s *DatabasesService) CreateBackup(ctx context.Context, id string, backupReq *CreateBackupRequest) (*Backup, *Response, error)
    CreateBackup creates a new backup for a database

func (s *DatabasesService) Delete(ctx context.Context, id string) (*Response, error)
    Delete deletes a database

func (s *DatabasesService) DeleteBackup(ctx context.Context, dbID, backupID string) (*Response, error)
    DeleteBackup deletes a backup

func (s *DatabasesService) DisablePublicAccess(ctx context.Context, id string) (*Database, *Response, error)
    DisablePublicAccess disables public access to a database

func (s *DatabasesService) EnablePublicAccess(ctx context.Context, id string) (*Database, *Response, error)
    EnablePublicAccess enables public access to a database

func (s *DatabasesService) Get(ctx context.Context, id string) (*Database, *Response, error)
    Get returns a single database by ID

func (s *DatabasesService) GetBackup(ctx context.Context, dbID, backupID string) (*Backup, *Response, error)
    GetBackup gets a specific backup

func (s *DatabasesService) GetCredentials(ctx context.Context, id string) (*Database, *Response, error)
    GetCredentials retrieves database connection credentials

func (s *DatabasesService) GetUsage(ctx context.Context, id string, period string) (*Usage, *Response, error)
    GetUsage retrieves usage metrics for a database

func (s *DatabasesService) List(ctx context.Context, opts *ListOptions) ([]*Database, *Response, error)
    List returns all databases

func (s *DatabasesService) ListBackups(ctx context.Context, id string, opts *ListOptions) ([]*Backup, *Response, error)
    ListBackups lists all backups for a database

func (s *DatabasesService) ResetPassword(ctx context.Context, id string) (*Database, *Response, error)
    ResetPassword resets the database password

func (s *DatabasesService) RestoreFromBackup(ctx context.Context, id string, restoreReq *RestoreBackupRequest) (*Response, error)
    RestoreFromBackup restores a database from a backup

func (s *DatabasesService) Update(ctx context.Context, id string, updateReq *UpdateDatabaseRequest) (*Database, *Response, error)
    Update updates an existing database

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
    Deployment represents a deployment

type DeploymentsService struct {
	// Has unexported fields.
}
    DeploymentsService handles communication with the deployment-related methods
    of the Sevalla API.

func (s *DeploymentsService) Cancel(ctx context.Context, id string) (*Response, error)
    Cancel cancels a deployment

func (s *DeploymentsService) Get(ctx context.Context, id string) (*Deployment, *Response, error)
    Get returns a single deployment by ID

func (s *DeploymentsService) GetLogs(ctx context.Context, id string) (string, *Response, error)
    GetLogs retrieves deployment logs

func (s *DeploymentsService) List(ctx context.Context, opts *ListOptions) ([]*Deployment, *Response, error)
    List returns all deployments

type Engine string
    Engine represents a database engine type

const (
	EnginePostgreSQL Engine = "postgresql"
	EngineMySQL      Engine = "mysql"
	EngineMongoDB    Engine = "mongodb"
	EngineRedis      Engine = "redis"
)
    Available database engines

type ErrorDetail struct {
	Field   string `json:"field,omitempty"`
	Code    string `json:"code,omitempty"`
	Message string `json:"message"`
}
    ErrorDetail represents a detailed error message

type ErrorResponse struct {
	Response  *http.Response `json:"-"`
	Message   string         `json:"message"`
	Code      string         `json:"code,omitempty"`
	RequestID string         `json:"request_id,omitempty"`
	Errors    []ErrorDetail  `json:"errors,omitempty"`
}
    ErrorResponse represents an error response from the Sevalla API

func (e *ErrorResponse) Error() string
    Error returns the error message

type ListOptions struct {
	Page      int    `url:"page,omitempty"`
	PerPage   int    `url:"per_page,omitempty"`
	Sort      string `url:"sort,omitempty"`
	Order     string `url:"order,omitempty"`
	CompanyID string `url:"company_id,omitempty"`
}
    ListOptions represents options for listing resources

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
    Pipeline represents a CI/CD pipeline

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
    PipelineRun represents an execution of a pipeline

type PipelineRunStep struct {
	Name         string     `json:"name"`
	State        Status     `json:"state"`
	Output       string     `json:"output,omitempty"`
	ErrorMessage string     `json:"error_message,omitempty"`
	StartedAt    time.Time  `json:"started_at"`
	CompletedAt  *time.Time `json:"completed_at,omitempty"`
}
    PipelineRunStep represents the execution of a pipeline step

type PipelineStep struct {
	Name      string   `json:"name"`
	Command   string   `json:"command"`
	Image     string   `json:"image,omitempty"`
	Timeout   int      `json:"timeout_seconds,omitempty"`
	Retries   int      `json:"retries,omitempty"`
	DependsOn []string `json:"depends_on,omitempty"`
}
    PipelineStep represents a step in a pipeline

type PipelinesService struct {
	// Has unexported fields.
}
    PipelinesService handles communication with pipeline-related endpoints

func (s *PipelinesService) CancelRun(ctx context.Context, pipelineID string, runID string) (*Response, error)
    CancelRun cancels a pipeline run

func (s *PipelinesService) Create(ctx context.Context, createReq *CreatePipelineRequest) (*Pipeline, *Response, error)
    Create creates a new pipeline

func (s *PipelinesService) Delete(ctx context.Context, id string) (*Response, error)
    Delete deletes a pipeline

func (s *PipelinesService) Get(ctx context.Context, id string) (*Pipeline, *Response, error)
    Get retrieves a single pipeline by ID

func (s *PipelinesService) GetRun(ctx context.Context, pipelineID string, runID string) (*PipelineRun, *Response, error)
    GetRun retrieves a single pipeline run

func (s *PipelinesService) GetRunLogs(ctx context.Context, pipelineID string, runID string) (string, *Response, error)
    GetRunLogs retrieves logs for a pipeline run

func (s *PipelinesService) List(ctx context.Context, opts *ListOptions) ([]*Pipeline, *Response, error)
    List retrieves all pipelines

func (s *PipelinesService) ListRuns(ctx context.Context, pipelineID string, opts *ListOptions) ([]*PipelineRun, *Response, error)
    ListRuns retrieves all runs for a pipeline

func (s *PipelinesService) RetryRun(ctx context.Context, pipelineID string, runID string) (*PipelineRun, *Response, error)
    RetryRun retries a failed pipeline run

func (s *PipelinesService) Run(ctx context.Context, id string) (*PipelineRun, *Response, error)
    Run triggers a pipeline run

func (s *PipelinesService) Update(ctx context.Context, id string, updateReq *UpdatePipelineRequest) (*Pipeline, *Response, error)
    Update updates an existing pipeline

type Plan string
    Plan represents a pricing/compute plan tier

const (
	PlanHobby      Plan = "hobby"
	PlanStarter    Plan = "starter"
	PlanPro        Plan = "pro"
	PlanBusiness   Plan = "business"
	PlanEnterprise Plan = "enterprise"
)
    Available plan tiers

type Rate struct {
	Limit     int
	Remaining int
	Reset     time.Time
}
    Rate represents the rate limit information

type RateLimitError struct {
	*ErrorResponse
	RetryAfter int // Seconds to wait before retrying
}
    RateLimitError represents a rate limit error

func (e *RateLimitError) Error() string
    Error returns the rate limit error message

type Region string
    Region represents a deployment region

const (
	RegionUSCentral  Region = "us-central1"
	RegionUSEast     Region = "us-east1"
	RegionEuropeWest Region = "europe-west1"
	RegionAsiaSouth  Region = "asia-south1"
)
    Available deployment regions

type Response struct {
	*http.Response

	// Pagination
	NextPage  int
	PrevPage  int
	FirstPage int
	LastPage  int

	// Rate limiting
	Rate Rate
}
    Response wraps the standard HTTP response and includes pagination
    information

type RestoreBackupRequest struct {
	BackupID string `json:"backup_id"`
}
    RestoreBackupRequest represents a request to restore a database from backup

type ScaleApplicationRequest struct {
	Replicas int   `json:"replicas"`
	Plan     *Plan `json:"pod_size,omitempty"`
}
    ScaleApplicationRequest represents a request to scale an application

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
    StaticSite represents a Sevalla static site

type StaticSitesService struct {
	// Has unexported fields.
}
    StaticSitesService handles communication with the static site-related
    methods of the Sevalla API.

func (s *StaticSitesService) Create(ctx context.Context, createReq *CreateStaticSiteRequest) (*StaticSite, *Response, error)
    Create creates a new static site

func (s *StaticSitesService) Delete(ctx context.Context, id string) (*Response, error)
    Delete deletes a static site

func (s *StaticSitesService) Deploy(ctx context.Context, id string) (*Deployment, *Response, error)
    Deploy triggers a new deployment for a static site

func (s *StaticSitesService) Get(ctx context.Context, id string) (*StaticSite, *Response, error)
    Get returns a single static site by ID

func (s *StaticSitesService) List(ctx context.Context, opts *ListOptions) ([]*StaticSite, *Response, error)
    List returns all static sites

type Status string
    Status represents a deployment status

const (
	StatusQueued    Status = "queued"
	StatusBuilding  Status = "building"
	StatusDeploying Status = "deploying"
	StatusSuccess   Status = "success"
	StatusFailed    Status = "failed"
	StatusCancelled Status = "cancelled"
)
    Deployment statuses

type UpdateApplicationRequest struct {
	Name            *string           `json:"name,omitempty"`
	Branch          *string           `json:"branch,omitempty"`
	Plan            *Plan             `json:"pod_size,omitempty"`
	Replicas        *int              `json:"replicas,omitempty"`
	EnvironmentVars map[string]string `json:"environment_variables,omitempty"`
	BuildCommand    *string           `json:"build_command,omitempty"`
	StartCommand    *string           `json:"start_command,omitempty"`
	Port            *int              `json:"port,omitempty"`
	AutoDeploy      *bool             `json:"auto_deploy,omitempty"`
}
    UpdateApplicationRequest represents a request to update an application

type UpdateDatabaseRequest struct {
	Name       *string `json:"name,omitempty"`
	Size       *string `json:"size,omitempty"`
	Storage    *int    `json:"storage_gb,omitempty"`
	Backups    *bool   `json:"backups_enabled,omitempty"`
	SSLEnabled *bool   `json:"ssl_enabled,omitempty"`
}
    UpdateDatabaseRequest represents a request to update a database

type UpdatePipelineRequest struct {
	Name        *string                `json:"name,omitempty"`
	Enabled     *bool                  `json:"enabled,omitempty"`
	Trigger     *string                `json:"trigger,omitempty"`
	Branch      *string                `json:"branch,omitempty"`
	Steps       []PipelineStep         `json:"steps,omitempty"`
	Environment map[string]string      `json:"environment,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}
    UpdatePipelineRequest represents a request to update a pipeline

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
    Usage represents resource usage metrics

type ValidationError struct {
	Field   string
	Message string
}
    ValidationError represents a validation error

func (e *ValidationError) Error() string
    Error returns the validation error message

