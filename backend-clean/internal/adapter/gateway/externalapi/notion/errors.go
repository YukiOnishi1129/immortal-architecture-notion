package notion

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// APIError is a non-2xx response from the Notion API.
//
// The message comes from Notion's own error body. The request payload
// and the API token are never included, so logging this error cannot
// leak credentials.
type APIError struct {
	StatusCode int
	Code       string
	Message    string
}

func (e *APIError) Error() string {
	if e.Code != "" {
		return fmt.Sprintf("notion: %d %s: %s", e.StatusCode, e.Code, e.Message)
	}
	return fmt.Sprintf("notion: %d: %s", e.StatusCode, e.Message)
}

// IsNotFound reports whether the page is missing or not shared with
// the integration. Notion answers both cases with 404 object_not_found.
func (e *APIError) IsNotFound() bool {
	return e.StatusCode == http.StatusNotFound
}

// IsUnauthorized reports whether the token is invalid or expired.
func (e *APIError) IsUnauthorized() bool {
	return e.StatusCode == http.StatusUnauthorized
}

// IsForbidden reports whether the integration lacks permission.
func (e *APIError) IsForbidden() bool {
	return e.StatusCode == http.StatusForbidden
}

// IsRateLimited reports whether the request hit the rate limit.
func (e *APIError) IsRateLimited() bool {
	return e.StatusCode == http.StatusTooManyRequests
}

// newAPIError builds an APIError from a response body.
// A body that is not valid JSON still produces a usable error.
func newAPIError(statusCode int, body []byte) *APIError {
	var parsed struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	}

	if err := json.Unmarshal(body, &parsed); err != nil || parsed.Message == "" {
		return &APIError{
			StatusCode: statusCode,
			Message:    http.StatusText(statusCode),
		}
	}

	return &APIError{
		StatusCode: statusCode,
		Code:       parsed.Code,
		Message:    parsed.Message,
	}
}
