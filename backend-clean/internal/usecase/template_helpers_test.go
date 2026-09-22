package usecase

import (
	"errors"
	"testing"

	domainerr "immortal-architecture-notion/backend/internal/domain/errors"
)

func TestResolveParentPageID(t *testing.T) {
	tests := []struct {
		name      string
		rawURL    string
		want      string
		wantError error
	}{
		{
			name:   "[Success] empty url means not linked",
			rawURL: "",
			want:   "",
		},
		{
			name:   "[Success] whitespace only means not linked",
			rawURL: "   ",
			want:   "",
		},
		{
			name:   "[Success] plain url",
			rawURL: "https://notion.so/1429989fe8ac4effbc8f57f56486db54",
			want:   "1429989fe8ac4effbc8f57f56486db54",
		},
		{
			name:   "[Success] hyphenated id is normalized",
			rawURL: "https://notion.so/1429989f-e8ac-4eff-bc8f-57f56486db54",
			want:   "1429989fe8ac4effbc8f57f56486db54",
		},
		{
			name:   "[Success] title prefix and query string are ignored",
			rawURL: "https://www.notion.so/team/Daily-Note-1429989fe8ac4effbc8f57f56486db54?v=abcdef",
			want:   "1429989fe8ac4effbc8f57f56486db54",
		},
		{
			name:      "[Fail] url without a page id",
			rawURL:    "https://example.com/abc",
			wantError: domainerr.ErrInvalidNotionParentURL,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := resolveParentPageID(tt.rawURL)

			if tt.wantError != nil {
				if !errors.Is(err, tt.wantError) {
					t.Fatalf("expected error %v, got %v", tt.wantError, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Fatalf("resolveParentPageID() = %q, want %q", got, tt.want)
			}
		})
	}
}
