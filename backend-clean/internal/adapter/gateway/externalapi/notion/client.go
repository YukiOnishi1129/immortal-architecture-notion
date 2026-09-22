// Package notion implements the Notion API client.
package notion

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"time"

	"immortal-architecture-notion/backend/internal/port"
)

const (
	defaultBaseURL = "https://api.notion.com"
	notionVersion  = "2022-06-28"

	defaultTimeout    = 10 * time.Second
	defaultMaxRetries = 3
)

// Client talks to the Notion API over HTTP.
type Client struct {
	token      string
	baseURL    string
	httpClient *http.Client
	maxRetries int
	// sleep is swappable so tests do not wait for real backoff.
	sleep func(time.Duration)
}

var _ port.NotionClient = (*Client)(nil)

// Option customizes the client.
type Option func(*Client)

// WithBaseURL points the client at a different host. Used by tests.
func WithBaseURL(url string) Option {
	return func(c *Client) { c.baseURL = url }
}

// WithTimeout overrides the request timeout.
func WithTimeout(d time.Duration) Option {
	return func(c *Client) { c.httpClient.Timeout = d }
}

// WithMaxRetries overrides how many times a retryable request is repeated.
func WithMaxRetries(n int) Option {
	return func(c *Client) { c.maxRetries = n }
}

// WithSleep replaces the backoff sleep. Used by tests.
func WithSleep(fn func(time.Duration)) Option {
	return func(c *Client) { c.sleep = fn }
}

// NewClient creates a Notion API client.
// The HTTP client always has a timeout: the zero value of http.Client
// waits forever, which would hold up the caller indefinitely.
func NewClient(token string, opts ...Option) *Client {
	c := &Client{
		token:      token,
		baseURL:    defaultBaseURL,
		httpClient: &http.Client{Timeout: defaultTimeout},
		maxRetries: defaultMaxRetries,
		sleep:      time.Sleep,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// CreatePage creates a page under parentPageID.
func (c *Client) CreatePage(ctx context.Context, parentPageID, title string, sections []port.NotionSection) (*port.NotionPage, error) {
	body := map[string]any{
		"parent":     map[string]string{"page_id": parentPageID},
		"properties": titleProperty(title),
		"children":   toBlocks(sections),
	}

	var res pageResponse
	if err := c.do(ctx, http.MethodPost, "/v1/pages", body, &res); err != nil {
		return nil, err
	}
	return &port.NotionPage{PageID: res.ID, URL: res.URL}, nil
}

// UpdatePage replaces the title and content of an existing page.
//
// Notion has no "replace all content" call, so existing blocks are
// removed one by one before the new ones are appended.
func (c *Client) UpdatePage(ctx context.Context, pageID, title string, sections []port.NotionSection) (*port.NotionPage, error) {
	if err := c.do(ctx, http.MethodPatch, "/v1/pages/"+pageID,
		map[string]any{"properties": titleProperty(title)}, nil); err != nil {
		return nil, err
	}

	if err := c.clearBlocks(ctx, pageID); err != nil {
		return nil, err
	}

	if err := c.do(ctx, http.MethodPatch, "/v1/blocks/"+pageID+"/children",
		map[string]any{"children": toBlocks(sections)}, nil); err != nil {
		return nil, err
	}

	var res pageResponse
	if err := c.do(ctx, http.MethodGet, "/v1/pages/"+pageID, nil, &res); err != nil {
		return nil, err
	}
	return &port.NotionPage{PageID: res.ID, URL: res.URL}, nil
}

// Trash moves a page to the Notion trash. The page id stays valid.
func (c *Client) Trash(ctx context.Context, pageID string) error {
	return c.do(ctx, http.MethodPatch, "/v1/pages/"+pageID,
		map[string]any{"in_trash": true}, nil)
}

// Restore brings a trashed page back under its original id.
func (c *Client) Restore(ctx context.Context, pageID string) (*port.NotionPage, error) {
	var res pageResponse
	if err := c.do(ctx, http.MethodPatch, "/v1/pages/"+pageID,
		map[string]any{"in_trash": false}, &res); err != nil {
		return nil, err
	}
	return &port.NotionPage{PageID: res.ID, URL: res.URL}, nil
}

// clearBlocks deletes every child block of a page.
//
// The children endpoint is paginated, so the list is walked with
// next_cursor until has_more is false. Without this, a page holding
// more than one batch of blocks would keep its leftover content.
func (c *Client) clearBlocks(ctx context.Context, pageID string) error {
	cursor := ""

	for {
		path := "/v1/blocks/" + pageID + "/children?page_size=100"
		if cursor != "" {
			path += "&start_cursor=" + url.QueryEscape(cursor)
		}

		var listed blockChildrenResponse
		if err := c.do(ctx, http.MethodGet, path, nil, &listed); err != nil {
			return err
		}

		for _, b := range listed.Results {
			if err := c.do(ctx, http.MethodDelete, "/v1/blocks/"+b.ID, nil, nil); err != nil {
				return err
			}
		}

		if !listed.HasMore || listed.NextCursor == "" {
			return nil
		}
		cursor = listed.NextCursor
	}
}

// do sends one request, retrying on 429 and 5xx.
func (c *Client) do(ctx context.Context, method, path string, body any, out any) error {
	var payload []byte
	if body != nil {
		encoded, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("notion: encode request: %w", err)
		}
		payload = encoded
	}

	var lastErr error
	for attempt := 0; attempt <= c.maxRetries; attempt++ {
		if attempt > 0 {
			// exponential backoff: 1s, 2s, 4s
			wait := time.Duration(math.Pow(2, float64(attempt-1))) * time.Second
			if err := c.waitOrCancel(ctx, wait); err != nil {
				return err
			}
		}

		retryable, err := c.attempt(ctx, method, path, payload, out)
		if err == nil {
			return nil
		}
		lastErr = err
		if !retryable {
			return err
		}
	}
	return lastErr
}

// waitOrCancel sleeps for d, but returns early if the context is done.
// Without this, a canceled caller would still wait out the full backoff.
func (c *Client) waitOrCancel(ctx context.Context, d time.Duration) error {
	done := make(chan struct{})
	go func() {
		c.sleep(d)
		close(done)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-done:
		return nil
	}
}

// attempt performs a single HTTP round trip.
// The bool reports whether the failure is worth retrying.
func (c *Client) attempt(ctx context.Context, method, path string, payload []byte, out any) (bool, error) {
	var reader io.Reader
	if payload != nil {
		reader = bytes.NewReader(payload)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, reader)
	if err != nil {
		return false, fmt.Errorf("notion: build request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Notion-Version", notionVersion)
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		// Network failures and timeouts are worth another try.
		return true, fmt.Errorf("notion: request failed: %w", err)
	}
	defer res.Body.Close() //nolint:errcheck // nothing to do if closing fails

	raw, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		return true, fmt.Errorf("notion: read response: %w", readErr)
	}

	if res.StatusCode >= 200 && res.StatusCode < 300 {
		if out == nil {
			return false, nil
		}
		if err := json.Unmarshal(raw, out); err != nil {
			return false, fmt.Errorf("notion: decode response: %w", err)
		}
		return false, nil
	}

	apiErr := newAPIError(res.StatusCode, raw)
	// 429 and 5xx are transient; 4xx means the request itself is wrong.
	retryable := res.StatusCode == http.StatusTooManyRequests || res.StatusCode >= 500
	return retryable, apiErr
}

func titleProperty(title string) map[string]any {
	return map[string]any{
		"title": map[string]any{
			"title": []richText{{Type: "text", Text: textBody{Content: truncateTitle(title)}}},
		},
	}
}

type pageResponse struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

type blockChildrenResponse struct {
	Results []struct {
		ID string `json:"id"`
	} `json:"results"`
	HasMore    bool   `json:"has_more"`
	NextCursor string `json:"next_cursor"`
}
