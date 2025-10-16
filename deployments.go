package sevalla

import (
	"context"
	"fmt"
)

// DeploymentsService handles communication with the deployment-related
// methods of the Sevalla API.
type DeploymentsService struct {
	client *Client
}

// Get returns a single deployment by ID
func (s *DeploymentsService) Get(ctx context.Context, id string) (*Deployment, *Response, error) {
	u := fmt.Sprintf("deployments/%s", id)
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

// List returns all deployments
func (s *DeploymentsService) List(ctx context.Context, opts *ListOptions) ([]*Deployment, *Response, error) {
	u := "deployments"
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

// GetLogs retrieves deployment logs
func (s *DeploymentsService) GetLogs(ctx context.Context, id string) (string, *Response, error) {
	u := fmt.Sprintf("deployments/%s/logs", id)
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

// Cancel cancels a deployment
func (s *DeploymentsService) Cancel(ctx context.Context, id string) (*Response, error) {
	u := fmt.Sprintf("deployments/%s/cancel", id)
	req, err := s.client.NewRequest(ctx, "POST", u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}
