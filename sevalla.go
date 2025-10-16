// Package sevalla provides a Go client library for the Sevalla API
package sevalla

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/google/go-querystring/query"
)

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

// Client manages communication with the Sevalla API
type Client struct {
	// HTTP client for API requests
	client *http.Client

	// Base URL for API requests
	baseURL *url.URL

	// API key for authentication
	apiKey string

	// User agent for requests
	userAgent string

	// Services
	Applications *ApplicationsService
	Databases    *DatabasesService
	StaticSites  *StaticSitesService
	Deployments  *DeploymentsService
	Pipelines    *PipelinesService
}

// ClientOption is a function that configures a Client
type ClientOption func(*Client)

// WithAPIKey sets the API key for authentication
func WithAPIKey(key string) ClientOption {
	return func(c *Client) {
		c.apiKey = key
	}
}

// WithHTTPClient sets a custom HTTP client
func WithHTTPClient(client *http.Client) ClientOption {
	return func(c *Client) {
		c.client = client
	}
}

// WithBaseURL sets a custom base URL for the API
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) {
		if u, err := url.Parse(baseURL); err == nil {
			c.baseURL = u
		}
	}
}

// WithUserAgent sets a custom user agent
func WithUserAgent(ua string) ClientOption {
	return func(c *Client) {
		c.userAgent = ua
	}
}

// NewClient creates a new Sevalla API client
func NewClient(opts ...ClientOption) *Client {
	baseURL, _ := url.Parse(BaseURL)

	c := &Client{
		client:    &http.Client{Timeout: DefaultTimeout},
		baseURL:   baseURL,
		userAgent: UserAgent,
	}

	// Apply options
	for _, opt := range opts {
		opt(c)
	}

	// Initialize services
	c.Applications = &ApplicationsService{client: c}
	c.Databases = &DatabasesService{client: c}
	c.StaticSites = &StaticSitesService{client: c}
	c.Deployments = &DeploymentsService{client: c}
	c.Pipelines = &PipelinesService{client: c}

	return c
}

// NewRequest creates an API request
func (c *Client) NewRequest(ctx context.Context, method, urlStr string, body interface{}) (*http.Request, error) {
	u, err := c.baseURL.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	var buf io.ReadWriter
	if body != nil {
		buf = &bytes.Buffer{}
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		if err := enc.Encode(body); err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	// Set authentication
	if c.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
	}

	// Set user agent
	req.Header.Set("User-Agent", c.userAgent)

	return req, nil
}

// NewRequestWithQuery creates an API request with query parameters
func (c *Client) NewRequestWithQuery(ctx context.Context, method, urlStr string, opts interface{}) (*http.Request, error) {
	u, err := c.baseURL.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	if opts != nil {
		v, err := query.Values(opts)
		if err != nil {
			return nil, err
		}
		u.RawQuery = v.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), nil)
	if err != nil {
		return nil, err
	}

	// Set authentication
	if c.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
	}

	// Set user agent
	req.Header.Set("User-Agent", c.userAgent)

	return req, nil
}

// Do executes an API request and returns the response
func (c *Client) Do(req *http.Request, v interface{}) (*Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	response := &Response{Response: resp}
	response.populatePageValues()

	// Check for errors
	if err := CheckResponse(resp); err != nil {
		return response, err
	}

	// Decode response body if v is provided
	if v != nil && resp.StatusCode != http.StatusNoContent {
		if w, ok := v.(io.Writer); ok {
			_, err = io.Copy(w, resp.Body)
		} else {
			decErr := json.NewDecoder(resp.Body).Decode(v)
			if decErr == io.EOF {
				decErr = nil // Ignore EOF errors from empty response
			}
			if decErr != nil {
				err = decErr
			}
		}
	}

	return response, err
}

// Response wraps the standard HTTP response and includes pagination information
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

// populatePageValues populates the pagination values from Link header
func (r *Response) populatePageValues() {
	if links := r.Header.Get("Link"); links != "" {
		for _, link := range strings.Split(links, ",") {
			segments := strings.Split(strings.TrimSpace(link), ";")
			if len(segments) != 2 {
				continue
			}

			// Parse URL
			urlStr := strings.Trim(segments[0], "<>")
			u, err := url.Parse(urlStr)
			if err != nil {
				continue
			}

			// Get page number
			page := u.Query().Get("page")
			if page == "" {
				continue
			}

			pageNum, err := strconv.Atoi(page)
			if err != nil {
				continue
			}

			// Get rel value
			rel := strings.Trim(strings.Split(segments[1], "=")[1], "\"")

			switch rel {
			case "next":
				r.NextPage = pageNum
			case "prev":
				r.PrevPage = pageNum
			case "first":
				r.FirstPage = pageNum
			case "last":
				r.LastPage = pageNum
			}
		}
	}
}

// Rate represents the rate limit information
type Rate struct {
	Limit     int
	Remaining int
	Reset     time.Time
}

// CheckResponse checks the API response for errors
func CheckResponse(r *http.Response) error {
	if c := r.StatusCode; c >= 200 && c <= 299 {
		return nil
	}

	errorResponse := &ErrorResponse{Response: r}
	data, err := io.ReadAll(r.Body)
	if err == nil && len(data) > 0 {
		if err := json.Unmarshal(data, errorResponse); err != nil {
			errorResponse.Message = string(data)
		}
	}

	return errorResponse
}

// Bool is a helper function that allocates a new bool value
func Bool(v bool) *bool { return &v }

// Int is a helper function that allocates a new int value
func Int(v int) *int { return &v }

// String is a helper function that allocates a new string value
func String(v string) *string { return &v }
