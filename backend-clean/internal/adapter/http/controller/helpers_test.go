package controller

import (
	"net/http/httptest"
	"strings"
	"testing"
)

// assertStatusBody checks HTTP status and optional body substring.
func assertStatusBody(t *testing.T, rec *httptest.ResponseRecorder, wantStatus int, wantBody string) {
	t.Helper()
	if rec.Code != wantStatus {
		t.Fatalf("status = %d, want %d", rec.Code, wantStatus)
	}
	if wantBody != "" && !strings.Contains(rec.Body.String(), wantBody) {
		t.Fatalf("body = %q, want to contain %q", rec.Body.String(), wantBody)
	}
}
