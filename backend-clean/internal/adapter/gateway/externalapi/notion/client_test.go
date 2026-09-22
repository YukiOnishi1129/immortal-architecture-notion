package notion

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"immortal-architecture-notion/backend/internal/port"
)

// newTestClient points a client at a test server and removes the
// backoff wait so retry cases finish instantly.
func newTestClient(url string, opts ...Option) *Client {
	base := []Option{
		WithBaseURL(url),
		WithSleep(func(time.Duration) {}),
	}
	return NewClient("test-token", append(base, opts...)...)
}

func TestClient_CreatePage(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		body       string
		wantErr    bool
		wantCalls  int
	}{
		{
			name:       "[Success] page created",
			statusCode: http.StatusOK,
			body:       `{"id":"page-1","url":"https://notion.so/page-1"}`,
			wantCalls:  1,
		},
		{
			name:       "[Fail] 401 is not retried",
			statusCode: http.StatusUnauthorized,
			body:       `{"code":"unauthorized","message":"API token is invalid."}`,
			wantErr:    true,
			wantCalls:  1,
		},
		{
			name:       "[Fail] 404 is not retried",
			statusCode: http.StatusNotFound,
			body:       `{"code":"object_not_found","message":"Could not find page."}`,
			wantErr:    true,
			wantCalls:  1,
		},
		{
			name:       "[Fail] 400 is not retried",
			statusCode: http.StatusBadRequest,
			body:       `{"code":"validation_error","message":"body failed validation."}`,
			wantErr:    true,
			wantCalls:  1,
		},
		{
			name:       "[Fail] 429 is retried",
			statusCode: http.StatusTooManyRequests,
			body:       `{"code":"rate_limited","message":"Rate limited."}`,
			wantErr:    true,
			wantCalls:  4, // first attempt + 3 retries
		},
		{
			name:       "[Fail] 500 is retried",
			statusCode: http.StatusInternalServerError,
			body:       `{"code":"internal_server_error","message":"Unexpected error."}`,
			wantErr:    true,
			wantCalls:  4,
		},
		{
			name:       "[Fail] 503 is retried",
			statusCode: http.StatusServiceUnavailable,
			body:       `{"code":"service_unavailable","message":"Unavailable."}`,
			wantErr:    true,
			wantCalls:  4,
		},
		{
			name:       "[Fail] malformed json body",
			statusCode: http.StatusOK,
			body:       `{broken`,
			wantErr:    true,
			wantCalls:  1,
		},
		{
			name:       "[Fail] error body is not json",
			statusCode: http.StatusBadGateway,
			body:       `<html>gateway error</html>`,
			wantErr:    true,
			wantCalls:  4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calls := 0
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls++

				if got := r.Header.Get("Notion-Version"); got != notionVersion {
					t.Errorf("Notion-Version = %q, want %q", got, notionVersion)
				}
				if got := r.Header.Get("Authorization"); got != "Bearer test-token" {
					t.Errorf("Authorization = %q", got)
				}
				if got := r.URL.Path; got != "/v1/pages" {
					t.Errorf("path = %q, want /v1/pages", got)
				}

				w.WriteHeader(tt.statusCode)
				_, _ = w.Write([]byte(tt.body))
			}))
			defer srv.Close()

			client := newTestClient(srv.URL)
			page, err := client.CreatePage(context.Background(), "parent-1", "Title",
				[]port.NotionSection{{Label: "Field", Content: "content"}})

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected an error, got nil")
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if page.PageID != "page-1" {
					t.Errorf("PageID = %q, want page-1", page.PageID)
				}
			}

			if calls != tt.wantCalls {
				t.Errorf("request count = %d, want %d", calls, tt.wantCalls)
			}
		})
	}
}

func TestClient_CreatePage_RequestBody(t *testing.T) {
	var got map[string]any

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewDecoder(r.Body).Decode(&got)
		_, _ = w.Write([]byte(`{"id":"page-1","url":"https://notion.so/page-1"}`))
	}))
	defer srv.Close()

	client := newTestClient(srv.URL)
	_, err := client.CreatePage(context.Background(), "parent-1", "My Note",
		[]port.NotionSection{
			{Label: "背景", Content: "なぜ作るか"},
			{Label: "対策", Content: "どうするか"},
		})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	parent, ok := got["parent"].(map[string]any)
	if !ok || parent["page_id"] != "parent-1" {
		t.Fatalf("parent = %v, want page_id parent-1", got["parent"])
	}

	children, ok := got["children"].([]any)
	if !ok {
		t.Fatal("children is missing")
	}
	// two sections -> heading + paragraph each
	if len(children) != 4 {
		t.Fatalf("children count = %d, want 4", len(children))
	}

	first, _ := children[0].(map[string]any)
	if first["type"] != "heading_2" {
		t.Errorf("first block type = %v, want heading_2", first["type"])
	}
	second, _ := children[1].(map[string]any)
	if second["type"] != "paragraph" {
		t.Errorf("second block type = %v, want paragraph", second["type"])
	}
}

func TestClient_Timeout(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
	}))
	defer srv.Close()

	client := newTestClient(srv.URL,
		WithTimeout(20*time.Millisecond),
		WithMaxRetries(0),
	)

	_, err := client.CreatePage(context.Background(), "parent-1", "Title", nil)
	if err == nil {
		t.Fatal("expected a timeout error, got nil")
	}
}

func TestClient_Trash(t *testing.T) {
	var got map[string]any
	var gotPath string

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		_ = json.NewDecoder(r.Body).Decode(&got)
		_, _ = w.Write([]byte(`{"id":"page-1","url":"https://notion.so/page-1"}`))
	}))
	defer srv.Close()

	client := newTestClient(srv.URL)
	if err := client.Trash(context.Background(), "page-1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if gotPath != "/v1/pages/page-1" {
		t.Errorf("path = %q, want /v1/pages/page-1", gotPath)
	}
	if got["in_trash"] != true {
		t.Errorf("in_trash = %v, want true", got["in_trash"])
	}
}

func TestClient_Restore(t *testing.T) {
	var got map[string]any

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewDecoder(r.Body).Decode(&got)
		_, _ = w.Write([]byte(`{"id":"page-1","url":"https://notion.so/page-1"}`))
	}))
	defer srv.Close()

	client := newTestClient(srv.URL)
	page, err := client.Restore(context.Background(), "page-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got["in_trash"] != false {
		t.Errorf("in_trash = %v, want false", got["in_trash"])
	}
	if page.URL != "https://notion.so/page-1" {
		t.Errorf("URL = %q", page.URL)
	}
}

func TestClient_UpdatePage(t *testing.T) {
	var paths []string

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		paths = append(paths, r.Method+" "+r.URL.Path)

		switch {
		case r.Method == http.MethodGet && strings.HasSuffix(r.URL.Path, "/children"):
			_, _ = w.Write([]byte(`{"results":[{"id":"block-1"},{"id":"block-2"}]}`))
		default:
			_, _ = w.Write([]byte(`{"id":"page-1","url":"https://notion.so/page-1"}`))
		}
	}))
	defer srv.Close()

	client := newTestClient(srv.URL)
	page, err := client.UpdatePage(context.Background(), "page-1", "New Title",
		[]port.NotionSection{{Label: "Field", Content: "updated"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if page.PageID != "page-1" {
		t.Errorf("PageID = %q", page.PageID)
	}

	// title update -> list blocks -> delete each -> append -> re-read
	want := []string{
		"PATCH /v1/pages/page-1",
		"GET /v1/blocks/page-1/children",
		"DELETE /v1/blocks/block-1",
		"DELETE /v1/blocks/block-2",
		"PATCH /v1/blocks/page-1/children",
		"GET /v1/pages/page-1",
	}
	if len(paths) != len(want) {
		t.Fatalf("requests = %v, want %v", paths, want)
	}
	for i := range want {
		if paths[i] != want[i] {
			t.Errorf("request[%d] = %q, want %q", i, paths[i], want[i])
		}
	}
}

func TestAPIError(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		body       string
		wantCode   string
		check      func(*APIError) bool
	}{
		{
			name:       "404 is not found",
			statusCode: http.StatusNotFound,
			body:       `{"code":"object_not_found","message":"Could not find page."}`,
			wantCode:   "object_not_found",
			check:      (*APIError).IsNotFound,
		},
		{
			name:       "401 is unauthorized",
			statusCode: http.StatusUnauthorized,
			body:       `{"code":"unauthorized","message":"invalid token"}`,
			wantCode:   "unauthorized",
			check:      (*APIError).IsUnauthorized,
		},
		{
			name:       "403 is forbidden",
			statusCode: http.StatusForbidden,
			body:       `{"code":"restricted_resource","message":"no access"}`,
			wantCode:   "restricted_resource",
			check:      (*APIError).IsForbidden,
		},
		{
			name:       "429 is rate limited",
			statusCode: http.StatusTooManyRequests,
			body:       `{"code":"rate_limited","message":"slow down"}`,
			wantCode:   "rate_limited",
			check:      (*APIError).IsRateLimited,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := newAPIError(tt.statusCode, []byte(tt.body))

			if err.Code != tt.wantCode {
				t.Errorf("Code = %q, want %q", err.Code, tt.wantCode)
			}
			if !tt.check(err) {
				t.Error("status predicate returned false")
			}
			if err.Error() == "" {
				t.Error("Error() is empty")
			}
		})
	}
}

func TestAPIError_NonJSONBody(t *testing.T) {
	err := newAPIError(http.StatusBadGateway, []byte("<html>oops</html>"))

	if err.StatusCode != http.StatusBadGateway {
		t.Errorf("StatusCode = %d", err.StatusCode)
	}
	if err.Message != "Bad Gateway" {
		t.Errorf("Message = %q, want Bad Gateway", err.Message)
	}
}

func TestAPIError_DoesNotLeakToken(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"code":"unauthorized","message":"API token is invalid."}`))
	}))
	defer srv.Close()

	client := NewClient("super-secret-token", WithBaseURL(srv.URL), WithSleep(func(time.Duration) {}))
	_, err := client.CreatePage(context.Background(), "parent-1", "Title", nil)
	if err == nil {
		t.Fatal("expected an error")
	}

	if strings.Contains(err.Error(), "super-secret-token") {
		t.Fatalf("error message leaks the token: %v", err)
	}

	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *APIError, got %T", err)
	}
}

func TestClient_ContextCancellation(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
	}))
	defer srv.Close()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	client := newTestClient(srv.URL, WithMaxRetries(0))
	if _, err := client.CreatePage(ctx, "parent-1", "Title", nil); err == nil {
		t.Fatal("expected an error for a canceled context")
	}
}

func TestClient_UpdatePage_DeletesAllPages(t *testing.T) {
	// The first children response reports more pages, so the client
	// must follow next_cursor until has_more is false.
	var deleted []string
	var cursors []string
	listCalls := 0

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && strings.HasSuffix(r.URL.Path, "/children"):
			listCalls++
			cursors = append(cursors, r.URL.Query().Get("start_cursor"))

			if listCalls == 1 {
				_, _ = w.Write([]byte(`{
					"results":[{"id":"block-1"},{"id":"block-2"}],
					"has_more":true,
					"next_cursor":"cursor-abc"
				}`))
				return
			}
			_, _ = w.Write([]byte(`{
				"results":[{"id":"block-3"}],
				"has_more":false,
				"next_cursor":null
			}`))

		case r.Method == http.MethodDelete:
			deleted = append(deleted, strings.TrimPrefix(r.URL.Path, "/v1/blocks/"))
			_, _ = w.Write([]byte(`{"id":"block"}`))

		default:
			_, _ = w.Write([]byte(`{"id":"page-1","url":"https://notion.so/page-1"}`))
		}
	}))
	defer srv.Close()

	client := newTestClient(srv.URL)
	if _, err := client.UpdatePage(context.Background(), "page-1", "Title", nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if listCalls != 2 {
		t.Errorf("children list calls = %d, want 2", listCalls)
	}
	if len(cursors) == 2 && cursors[1] != "cursor-abc" {
		t.Errorf("second start_cursor = %q, want cursor-abc", cursors[1])
	}

	want := []string{"block-1", "block-2", "block-3"}
	if len(deleted) != len(want) {
		t.Fatalf("deleted = %v, want %v", deleted, want)
	}
	for i := range want {
		if deleted[i] != want[i] {
			t.Errorf("deleted[%d] = %q, want %q", i, deleted[i], want[i])
		}
	}
}

func TestClient_UpdatePage_StopsWhenNoMorePages(t *testing.T) {
	// has_more is false on the first response, so exactly one list call
	// should be made even though a cursor value is present.
	listCalls := 0

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet && strings.HasSuffix(r.URL.Path, "/children") {
			listCalls++
			_, _ = w.Write([]byte(`{"results":[],"has_more":false,"next_cursor":"unused"}`))
			return
		}
		_, _ = w.Write([]byte(`{"id":"page-1","url":"https://notion.so/page-1"}`))
	}))
	defer srv.Close()

	client := newTestClient(srv.URL)
	if _, err := client.UpdatePage(context.Background(), "page-1", "Title", nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if listCalls != 1 {
		t.Errorf("children list calls = %d, want 1", listCalls)
	}
}

func TestClient_BackoffIsCancellable(t *testing.T) {
	// A retryable status keeps the client in its backoff loop.
	// Cancelling the context must abort the wait instead of sleeping it out.
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
		_, _ = w.Write([]byte(`{"code":"rate_limited","message":"slow down"}`))
	}))
	defer srv.Close()

	ctx, cancel := context.WithCancel(context.Background())

	client := NewClient("test-token",
		WithBaseURL(srv.URL),
		// A long sleep would block for 10s if cancellation were ignored.
		WithSleep(func(time.Duration) { time.Sleep(10 * time.Second) }),
	)

	// Cancel while the client is waiting to retry.
	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()

	start := time.Now()
	_, err := client.CreatePage(ctx, "parent-1", "Title", nil)
	elapsed := time.Since(start)

	if err == nil {
		t.Fatal("expected an error after cancellation")
	}
	if !errors.Is(err, context.Canceled) {
		t.Errorf("error = %v, want context.Canceled", err)
	}
	if elapsed > 3*time.Second {
		t.Errorf("returned after %v; backoff did not abort on cancellation", elapsed)
	}
}
