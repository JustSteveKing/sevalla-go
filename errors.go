package sevalla

import (
	"fmt"
	"net/http"
)

// ErrorResponse represents an error response from the Sevalla API
type ErrorResponse struct {
	Response  *http.Response `json:"-"`
	Message   string         `json:"message"`
	Code      string         `json:"code,omitempty"`
	RequestID string         `json:"request_id,omitempty"`
	Errors    []ErrorDetail  `json:"errors,omitempty"`
}

// ErrorDetail represents a detailed error message
type ErrorDetail struct {
	Field   string `json:"field,omitempty"`
	Code    string `json:"code,omitempty"`
	Message string `json:"message"`
}

// Error returns the error message
func (e *ErrorResponse) Error() string {
	if e.Code != "" {
		return fmt.Sprintf("sevalla: %s (%d) - %s", e.Code, e.Response.StatusCode, e.Message)
	}

	if e.RequestID != "" {
		return fmt.Sprintf("sevalla: %d - %s (request_id: %s)", e.Response.StatusCode, e.Message, e.RequestID)
	}

	return fmt.Sprintf("sevalla: %d - %s", e.Response.StatusCode, e.Message)
}

// IsNotFound returns true if the error is a 404 Not Found
func IsNotFound(err error) bool {
	if e, ok := err.(*ErrorResponse); ok {
		return e.Response.StatusCode == http.StatusNotFound
	}
	return false
}

// IsBadRequest returns true if the error is a 400 Bad Request
func IsBadRequest(err error) bool {
	if e, ok := err.(*ErrorResponse); ok {
		return e.Response.StatusCode == http.StatusBadRequest
	}
	return false
}

// IsUnauthorized returns true if the error is a 401 Unauthorized
func IsUnauthorized(err error) bool {
	if e, ok := err.(*ErrorResponse); ok {
		return e.Response.StatusCode == http.StatusUnauthorized
	}
	return false
}

// IsForbidden returns true if the error is a 403 Forbidden
func IsForbidden(err error) bool {
	if e, ok := err.(*ErrorResponse); ok {
		return e.Response.StatusCode == http.StatusForbidden
	}
	return false
}

// IsConflict returns true if the error is a 409 Conflict
func IsConflict(err error) bool {
	if e, ok := err.(*ErrorResponse); ok {
		return e.Response.StatusCode == http.StatusConflict
	}
	return false
}

// IsUnprocessableEntity returns true if the error is a 422 Unprocessable Entity
func IsUnprocessableEntity(err error) bool {
	if e, ok := err.(*ErrorResponse); ok {
		return e.Response.StatusCode == http.StatusUnprocessableEntity
	}
	return false
}

// IsRateLimited returns true if the error is a 429 Too Many Requests
func IsRateLimited(err error) bool {
	if e, ok := err.(*ErrorResponse); ok {
		return e.Response.StatusCode == http.StatusTooManyRequests
	}
	return false
}

// IsServerError returns true if the error is a 5xx server error
func IsServerError(err error) bool {
	if e, ok := err.(*ErrorResponse); ok {
		return e.Response.StatusCode >= 500 && e.Response.StatusCode < 600
	}
	return false
}

// IsClientError returns true if the error is a 4xx client error
func IsClientError(err error) bool {
	if e, ok := err.(*ErrorResponse); ok {
		return e.Response.StatusCode >= 400 && e.Response.StatusCode < 500
	}
	return false
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string
	Message string
}

// Error returns the validation error message
func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error: %s - %s", e.Field, e.Message)
}

// RateLimitError represents a rate limit error
type RateLimitError struct {
	*ErrorResponse
	RetryAfter int // Seconds to wait before retrying
}

// Error returns the rate limit error message
func (e *RateLimitError) Error() string {
	return fmt.Sprintf("rate limited: retry after %d seconds - %s", e.RetryAfter, e.Message)
}
