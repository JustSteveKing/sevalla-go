package sevalla

import (
	"context"
	"fmt"
)

// ApplicationsService handles communication with the application-related
// methods of the Sevalla API.
type ApplicationsService struct {
	client *Client
}

// CreateApplicationRequest represents a request to create a new application
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

// UpdateApplicationRequest represents a request to update an application
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

// ScaleApplicationRequest represents a request to scale an application
type ScaleApplicationRequest struct {
	Replicas int   `json:"replicas"`
	Plan     *Plan `json:"pod_size,omitempty"`
}

// AddDomainRequest represents a request to add a custom domain
type AddDomainRequest struct {
	Domain string `json:"domain"`
}

// CDNSettingsRequest represents CDN configuration
type CDNSettingsRequest struct {
	Enabled bool `json:"enabled"`
}

// List returns all applications
func (s *ApplicationsService) List(ctx context.Context, opts *ListOptions) ([]*Application, *Response, error) {
	u := "applications"
	req, err := s.client.NewRequestWithQuery(ctx, "GET", u, opts)
	if err != nil {
		return nil, nil, err
	}

	var apps []*Application
	resp, err := s.client.Do(req, &apps)
	if err != nil {
		return nil, resp, err
	}

	return apps, resp, nil
}

// Get returns a single application by ID
func (s *ApplicationsService) Get(ctx context.Context, id string) (*Application, *Response, error) {
	u := fmt.Sprintf("applications/%s", id)
	req, err := s.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	app := new(Application)
	resp, err := s.client.Do(req, &app)
	if err != nil {
		return nil, resp, err
	}

	return app, resp, nil
}

// Create creates a new application
func (s *ApplicationsService) Create(ctx context.Context, createReq *CreateApplicationRequest) (*Application, *Response, error) {
	u := "applications"
	req, err := s.client.NewRequest(ctx, "POST", u, createReq)
	if err != nil {
		return nil, nil, err
	}

	app := new(Application)
	resp, err := s.client.Do(req, &app)
	if err != nil {
		return nil, resp, err
	}

	return app, resp, nil
}

// Update updates an existing application
func (s *ApplicationsService) Update(ctx context.Context, id string, updateReq *UpdateApplicationRequest) (*Application, *Response, error) {
	u := fmt.Sprintf("applications/%s", id)
	req, err := s.client.NewRequest(ctx, "PATCH", u, updateReq)
	if err != nil {
		return nil, nil, err
	}

	app := new(Application)
	resp, err := s.client.Do(req, &app)
	if err != nil {
		return nil, resp, err
	}

	return app, resp, nil
}

// Delete deletes an application
func (s *ApplicationsService) Delete(ctx context.Context, id string) (*Response, error) {
	u := fmt.Sprintf("applications/%s", id)
	req, err := s.client.NewRequest(ctx, "DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// Scale scales an application's resources
func (s *ApplicationsService) Scale(ctx context.Context, id string, scaleReq *ScaleApplicationRequest) (*Application, *Response, error) {
	u := fmt.Sprintf("applications/%s/scale", id)
	req, err := s.client.NewRequest(ctx, "POST", u, scaleReq)
	if err != nil {
		return nil, nil, err
	}

	app := new(Application)
	resp, err := s.client.Do(req, &app)
	if err != nil {
		return nil, resp, err
	}

	return app, resp, nil
}

// Deploy triggers a new deployment for an application
func (s *ApplicationsService) Deploy(ctx context.Context, id string) (*Deployment, *Response, error) {
	u := fmt.Sprintf("applications/%s/deployments", id)
	req, err := s.client.NewRequest(ctx, "POST", u, nil)
	if err != nil {
		return nil, nil, err
	}

	deployment := new(Deployment)
	resp, err := s.client.Do(req, &deployment)
	if err != nil {
		return nil, resp, err
	}

	return deployment, resp, nil
}

// Restart restarts an application
func (s *ApplicationsService) Restart(ctx context.Context, id string) (*Response, error) {
	u := fmt.Sprintf("applications/%s/restart", id)
	req, err := s.client.NewRequest(ctx, "POST", u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// Stop stops an application
func (s *ApplicationsService) Stop(ctx context.Context, id string) (*Response, error) {
	u := fmt.Sprintf("applications/%s/stop", id)
	req, err := s.client.NewRequest(ctx, "POST", u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// Start starts a stopped application
func (s *ApplicationsService) Start(ctx context.Context, id string) (*Response, error) {
	u := fmt.Sprintf("applications/%s/start", id)
	req, err := s.client.NewRequest(ctx, "POST", u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// GetLogs retrieves application logs
func (s *ApplicationsService) GetLogs(ctx context.Context, id string, lines int) (string, *Response, error) {
	u := fmt.Sprintf("applications/%s/logs", id)
	if lines > 0 {
		u = fmt.Sprintf("%s?lines=%d", u, lines)
	}

	req, err := s.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return "", nil, err
	}

	var result struct {
		Logs string `json:"logs"`
	}
	resp, err := s.client.Do(req, &result)
	if err != nil {
		return "", resp, err
	}

	return result.Logs, resp, nil
}

// ListDeployments lists all deployments for an application
func (s *ApplicationsService) ListDeployments(ctx context.Context, id string, opts *ListOptions) ([]*Deployment, *Response, error) {
	u := fmt.Sprintf("applications/%s/deployments", id)
	req, err := s.client.NewRequestWithQuery(ctx, "GET", u, opts)
	if err != nil {
		return nil, nil, err
	}

	var deployments []*Deployment
	resp, err := s.client.Do(req, &deployments)
	if err != nil {
		return nil, resp, err
	}

	return deployments, resp, nil
}

// GetDeployment gets a specific deployment for an application
func (s *ApplicationsService) GetDeployment(ctx context.Context, appID, deploymentID string) (*Deployment, *Response, error) {
	u := fmt.Sprintf("applications/%s/deployments/%s", appID, deploymentID)
	req, err := s.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	deployment := new(Deployment)
	resp, err := s.client.Do(req, &deployment)
	if err != nil {
		return nil, resp, err
	}

	return deployment, resp, nil
}

// CancelDeployment cancels a deployment
func (s *ApplicationsService) CancelDeployment(ctx context.Context, appID, deploymentID string) (*Response, error) {
	u := fmt.Sprintf("applications/%s/deployments/%s/cancel", appID, deploymentID)
	req, err := s.client.NewRequest(ctx, "POST", u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// AddCustomDomain adds a custom domain to an application
func (s *ApplicationsService) AddCustomDomain(ctx context.Context, id string, domain string) (*Response, error) {
	u := fmt.Sprintf("applications/%s/domains", id)
	req, err := s.client.NewRequest(ctx, "POST", u, &AddDomainRequest{Domain: domain})
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// RemoveCustomDomain removes a custom domain from an application
func (s *ApplicationsService) RemoveCustomDomain(ctx context.Context, id string, domain string) (*Response, error) {
	u := fmt.Sprintf("applications/%s/domains/%s", id, domain)
	req, err := s.client.NewRequest(ctx, "DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// UpdateCDNSettings updates CDN settings for an application
func (s *ApplicationsService) UpdateCDNSettings(ctx context.Context, id string, enabled bool) (*Response, error) {
	u := fmt.Sprintf("applications/%s/cdn", id)
	req, err := s.client.NewRequest(ctx, "PUT", u, &CDNSettingsRequest{Enabled: enabled})
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// GetUsage retrieves usage metrics for an application
func (s *ApplicationsService) GetUsage(ctx context.Context, id string, period string) (*Usage, *Response, error) {
	u := fmt.Sprintf("applications/%s/usage", id)
	if period != "" {
		u = fmt.Sprintf("%s?period=%s", u, period)
	}

	req, err := s.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	usage := new(Usage)
	resp, err := s.client.Do(req, &usage)
	if err != nil {
		return nil, resp, err
	}

	return usage, resp, nil
}

// SetEnvironmentVariables sets environment variables for an application
func (s *ApplicationsService) SetEnvironmentVariables(ctx context.Context, id string, vars map[string]string) (*Response, error) {
	u := fmt.Sprintf("applications/%s/env", id)
	req, err := s.client.NewRequest(ctx, "PUT", u, vars)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// GetEnvironmentVariables gets environment variables for an application
func (s *ApplicationsService) GetEnvironmentVariables(ctx context.Context, id string) (map[string]string, *Response, error) {
	u := fmt.Sprintf("applications/%s/env", id)
	req, err := s.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	vars := make(map[string]string)
	resp, err := s.client.Do(req, &vars)
	if err != nil {
		return nil, resp, err
	}

	return vars, resp, nil
}

// Rollback rolls back to a previous deployment
func (s *ApplicationsService) Rollback(ctx context.Context, appID, deploymentID string) (*Deployment, *Response, error) {
	u := fmt.Sprintf("applications/%s/rollback/%s", appID, deploymentID)
	req, err := s.client.NewRequest(ctx, "POST", u, nil)
	if err != nil {
		return nil, nil, err
	}

	deployment := new(Deployment)
	resp, err := s.client.Do(req, &deployment)
	if err != nil {
		return nil, resp, err
	}

	return deployment, resp, nil
}
