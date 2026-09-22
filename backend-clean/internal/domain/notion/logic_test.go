package notion

import (
	"errors"
	"testing"
	"time"

	domainerr "immortal-architecture-notion/backend/internal/domain/errors"
)

func strPtr(s string) *string { return &s }

func TestSync_IsFirstSync(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name string
		sync *Sync
		want bool
	}{
		{
			name: "[True] no sync record",
			sync: nil,
			want: true,
		},
		{
			name: "[True] page id is nil",
			sync: &Sync{NoteID: "note-1"},
			want: true,
		},
		{
			name: "[True] page id is empty",
			sync: &Sync{NoteID: "note-1", PageID: strPtr("")},
			want: true,
		},
		{
			name: "[False] page id exists",
			sync: &Sync{
				NoteID:   "note-1",
				PageID:   strPtr("1429989fe8ac4effbc8f57f56486db54"),
				PageURL:  strPtr("https://notion.so/1429989fe8ac4effbc8f57f56486db54"),
				SyncedAt: &now,
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.sync.IsFirstSync(); got != tt.want {
				t.Fatalf("IsFirstSync() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSync_PageIDOrEmpty(t *testing.T) {
	tests := []struct {
		name string
		sync *Sync
		want string
	}{
		{name: "[Empty] nil receiver", sync: nil, want: ""},
		{name: "[Empty] page id is nil", sync: &Sync{NoteID: "note-1"}, want: ""},
		{
			name: "[Value] page id exists",
			sync: &Sync{PageID: strPtr("abc123")},
			want: "abc123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.sync.PageIDOrEmpty(); got != tt.want {
				t.Fatalf("PageIDOrEmpty() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestExtractPageID(t *testing.T) {
	tests := []struct {
		name      string
		rawURL    string
		want      string
		wantError error
	}{
		{
			name:   "[Success] plain 32 characters",
			rawURL: "https://notion.so/1429989fe8ac4effbc8f57f56486db54",
			want:   "1429989fe8ac4effbc8f57f56486db54",
		},
		{
			name:   "[Success] with workspace name",
			rawURL: "https://notion.so/myworkspace/1429989fe8ac4effbc8f57f56486db54",
			want:   "1429989fe8ac4effbc8f57f56486db54",
		},
		{
			name:   "[Success] with page title prefix",
			rawURL: "https://www.notion.so/myteam/Daily-Note-1429989fe8ac4effbc8f57f56486db54",
			want:   "1429989fe8ac4effbc8f57f56486db54",
		},
		{
			name:   "[Success] hyphenated id is normalized",
			rawURL: "https://notion.so/1429989f-e8ac-4eff-bc8f-57f56486db54",
			want:   "1429989fe8ac4effbc8f57f56486db54",
		},
		{
			name:   "[Success] query string is ignored",
			rawURL: "https://notion.so/1429989fe8ac4effbc8f57f56486db54?v=aaaaaaaabbbbccccddddeeeeeeeeeeee",
			want:   "1429989fe8ac4effbc8f57f56486db54",
		},
		{
			name:   "[Success] uppercase is normalized",
			rawURL: "https://notion.so/1429989FE8AC4EFFBC8F57F56486DB54",
			want:   "1429989fe8ac4effbc8f57f56486db54",
		},
		{
			name:   "[Success] surrounding spaces are trimmed",
			rawURL: "  https://notion.so/1429989fe8ac4effbc8f57f56486db54  ",
			want:   "1429989fe8ac4effbc8f57f56486db54",
		},
		{
			name:   "[Success] bare id without url",
			rawURL: "1429989fe8ac4effbc8f57f56486db54",
			want:   "1429989fe8ac4effbc8f57f56486db54",
		},
		{
			name:      "[Fail] empty string",
			rawURL:    "",
			wantError: domainerr.ErrNotionParentNotSet,
		},
		{
			name:      "[Fail] whitespace only",
			rawURL:    "   ",
			wantError: domainerr.ErrNotionParentNotSet,
		},
		{
			name:      "[Fail] url without an id",
			rawURL:    "https://notion.so/myworkspace",
			wantError: domainerr.ErrNotionParentNotSet,
		},
		{
			name:      "[Fail] id is too short",
			rawURL:    "https://notion.so/1429989fe8ac",
			wantError: domainerr.ErrNotionParentNotSet,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExtractPageID(tt.rawURL)

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
				t.Fatalf("ExtractPageID() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestValidateParentPageID(t *testing.T) {
	tests := []struct {
		name      string
		pageID    string
		wantError error
	}{
		{
			name:   "[Success] valid id",
			pageID: "1429989fe8ac4effbc8f57f56486db54",
		},
		{
			name:   "[Success] hyphenated id",
			pageID: "1429989f-e8ac-4eff-bc8f-57f56486db54",
		},
		{
			name:      "[Fail] empty",
			pageID:    "",
			wantError: domainerr.ErrNotionParentNotSet,
		},
		{
			name:      "[Fail] whitespace only",
			pageID:    "  ",
			wantError: domainerr.ErrNotionParentNotSet,
		},
		{
			name:      "[Fail] not a page id at all",
			pageID:    "abc",
			wantError: domainerr.ErrNotionParentNotSet,
		},
		{
			name:      "[Fail] too short",
			pageID:    "1429989fe8ac",
			wantError: domainerr.ErrNotionParentNotSet,
		},
		{
			name:      "[Fail] contains non-hex characters",
			pageID:    "zzzz989fe8ac4effbc8f57f56486db54",
			wantError: domainerr.ErrNotionParentNotSet,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateParentPageID(tt.pageID)

			if tt.wantError != nil {
				if !errors.Is(err, tt.wantError) {
					t.Fatalf("expected error %v, got %v", tt.wantError, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}
