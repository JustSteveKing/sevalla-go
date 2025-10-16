package sevalla

import (
	"context"
	"fmt"
)

// PipelinesService handles communication with pipeline-related endpoints
type PipelinesService struct {
	client *Client
}

// CreatePipelineRequest represents a request to create a pipeline
type CreatePipelineRequest struct {
	Name        string                 `json:"name"`
	Enabled     bool                   `json:"enabled"`
	Trigger     string                 `json:"trigger"`
	Branch      string                 `json:"branch"`
	Steps       []PipelineStep         `json:"steps"`
	Environment map[string]string      `json:"environment,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// UpdatePipelineRequest represents a request to update a pipeline
type UpdatePipelineRequest struct {
	Name        *string                `json:"name,omitempty"`
	Enabled     *bool                  `json:"enabled,omitempty"`
	Trigger     *string                `json:"trigger,omitempty"`
	Branch      *string                `json:"branch,omitempty"`
	Steps       []PipelineStep         `json:"steps,omitempty"`
	Environment map[string]string      `json:"environment,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// List retrieves all pipelines
func (s *PipelinesService) List(ctx context.Context, opts *ListOptions) ([]*Pipeline, *Response, error) {
	u := "pipelines"
	req, err := s.client.NewRequestWithQuery(ctx, "GET", u, opts)
	if err != nil {
		return nil, nil, err
	}

	var pipelines []*Pipeline
	resp, err := s.client.Do(req, &pipelines)
	if err != nil {
		return nil, resp, err
	}

	return pipelines, resp, nil
}

// Get retrieves a single pipeline by ID
func (s *PipelinesService) Get(ctx context.Context, id string) (*Pipeline, *Response, error) {
	u := fmt.Sprintf("pipelines/%s", id)
	req, err := s.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var pipeline Pipeline
	resp, err := s.client.Do(req, &pipeline)
	if err != nil {
		return nil, resp, err
	}

	return &pipeline, resp, nil
}

// Create creates a new pipeline
func (s *PipelinesService) Create(ctx context.Context, createReq *CreatePipelineRequest) (*Pipeline, *Response, error) {
	req, err := s.client.NewRequest(ctx, "POST", "pipelines", createReq)
	if err != nil {
		return nil, nil, err
	}

	var pipeline Pipeline
	resp, err := s.client.Do(req, &pipeline)
	if err != nil {
		return nil, resp, err
	}

	return &pipeline, resp, nil
}

// Update updates an existing pipeline
func (s *PipelinesService) Update(ctx context.Context, id string, updateReq *UpdatePipelineRequest) (*Pipeline, *Response, error) {
	u := fmt.Sprintf("pipelines/%s", id)
	req, err := s.client.NewRequest(ctx, "PUT", u, updateReq)
	if err != nil {
		return nil, nil, err
	}

	var pipeline Pipeline
	resp, err := s.client.Do(req, &pipeline)
	if err != nil {
		return nil, resp, err
	}

	return &pipeline, resp, nil
}

// Delete deletes a pipeline
func (s *PipelinesService) Delete(ctx context.Context, id string) (*Response, error) {
	u := fmt.Sprintf("pipelines/%s", id)
	req, err := s.client.NewRequest(ctx, "DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// Run triggers a pipeline run
func (s *PipelinesService) Run(ctx context.Context, id string) (*PipelineRun, *Response, error) {
	u := fmt.Sprintf("pipelines/%s/runs", id)
	req, err := s.client.NewRequest(ctx, "POST", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var run PipelineRun
	resp, err := s.client.Do(req, &run)
	if err != nil {
		return nil, resp, err
	}

	return &run, resp, nil
}

// ListRuns retrieves all runs for a pipeline
func (s *PipelinesService) ListRuns(ctx context.Context, pipelineID string, opts *ListOptions) ([]*PipelineRun, *Response, error) {
	u := fmt.Sprintf("pipelines/%s/runs", pipelineID)
	req, err := s.client.NewRequestWithQuery(ctx, "GET", u, opts)
	if err != nil {
		return nil, nil, err
	}

	var runs []*PipelineRun
	resp, err := s.client.Do(req, &runs)
	if err != nil {
		return nil, resp, err
	}

	return runs, resp, nil
}

// GetRun retrieves a single pipeline run
func (s *PipelinesService) GetRun(ctx context.Context, pipelineID string, runID string) (*PipelineRun, *Response, error) {
	u := fmt.Sprintf("pipelines/%s/runs/%s", pipelineID, runID)
	req, err := s.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var run PipelineRun
	resp, err := s.client.Do(req, &run)
	if err != nil {
		return nil, resp, err
	}

	return &run, resp, nil
}

// CancelRun cancels a pipeline run
func (s *PipelinesService) CancelRun(ctx context.Context, pipelineID string, runID string) (*Response, error) {
	u := fmt.Sprintf("pipelines/%s/runs/%s/cancel", pipelineID, runID)
	req, err := s.client.NewRequest(ctx, "POST", u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// GetRunLogs retrieves logs for a pipeline run
func (s *PipelinesService) GetRunLogs(ctx context.Context, pipelineID string, runID string) (string, *Response, error) {
	u := fmt.Sprintf("pipelines/%s/runs/%s/logs", pipelineID, runID)
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

// RetryRun retries a failed pipeline run
func (s *PipelinesService) RetryRun(ctx context.Context, pipelineID string, runID string) (*PipelineRun, *Response, error) {
	u := fmt.Sprintf("pipelines/%s/runs/%s/retry", pipelineID, runID)
	req, err := s.client.NewRequest(ctx, "POST", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var run PipelineRun
	resp, err := s.client.Do(req, &run)
	if err != nil {
		return nil, resp, err
	}

	return &run, resp, nil
}
