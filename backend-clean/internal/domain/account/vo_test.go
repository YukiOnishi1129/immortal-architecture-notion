package account

import (
	"errors"
	"testing"
)

func TestParseEmail_VO(t *testing.T) {
	tests := []struct {
		name      string
		raw       string
		want      Email
		wantError error
	}{
		{name: "[Success] trims spaces", raw: "  user@example.com  ", want: "user@example.com"},
		{name: "[Fail] missing at", raw: "example.com", wantError: ErrInvalidEmail},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseEmail(tt.raw)
			if tt.wantError == nil && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.wantError != nil && !errors.Is(err, tt.wantError) {
				t.Fatalf("want %v, got %v", tt.wantError, err)
			}
			if err == nil && got != tt.want {
				t.Fatalf("want %s, got %s", tt.want, got)
			}
		})
	}
}

func TestEmail_String(t *testing.T) {
	t.Run("[Success] returns underlying string", func(t *testing.T) {
		e := Email("user@example.com")
		if e.String() != "user@example.com" {
			t.Fatalf("unexpected string: %s", e.String())
		}
	})
}
