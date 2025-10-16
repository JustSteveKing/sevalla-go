package sevalla

import (
	"context"
	"fmt"
)

// StaticSitesService handles communication with the static site-related
// methods of the Sevalla API.
type StaticSitesService struct {
	client *Client
}

// CreateStaticSiteRequest represents a request to create a new static site
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

// List returns all static sites
func (s *StaticSitesService) List(ctx context.Context, opts *ListOptions) ([]*StaticSite, *Response, error) {
	u := "static-sites"
	req, err := s.client.NewRequestWithQuery(ctx, "GET", u, opts)
	if err != nil {
		return nil, nil, err
	}

	var sites []*StaticSite
	resp, err := s.client.Do(req, &sites)
	if err != nil {
		return nil, resp, err
	}

	return sites, resp, nil
}

// Get returns a single static site by ID
func (s *StaticSitesService) Get(ctx context.Context, id string) (*StaticSite, *Response, error) {
	u := fmt.Sprintf("static-sites/%s", id)
	req, err := s.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	site := new(StaticSite)
	resp, err := s.client.Do(req, &site)
	if err != nil {
		return nil, resp, err
	}

	return site, resp, nil
}

// Create creates a new static site
func (s *StaticSitesService) Create(ctx context.Context, createReq *CreateStaticSiteRequest) (*StaticSite, *Response, error) {
	u := "static-sites"
	req, err := s.client.NewRequest(ctx, "POST", u, createReq)
	if err != nil {
		return nil, nil, err
	}

	site := new(StaticSite)
	resp, err := s.client.Do(req, &site)
	if err != nil {
		return nil, resp, err
	}

	return site, resp, nil
}

// Delete deletes a static site
func (s *StaticSitesService) Delete(ctx context.Context, id string) (*Response, error) {
	u := fmt.Sprintf("static-sites/%s", id)
	req, err := s.client.NewRequest(ctx, "DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// Deploy triggers a new deployment for a static site
func (s *StaticSitesService) Deploy(ctx context.Context, id string) (*Deployment, *Response, error) {
	u := fmt.Sprintf("static-sites/%s/deployments", id)
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
