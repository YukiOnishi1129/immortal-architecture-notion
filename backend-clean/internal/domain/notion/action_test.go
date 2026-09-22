package notion

import (
	"testing"

	"immortal-architecture-notion/backend/internal/domain/note"
)

func TestActionForStatusChange(t *testing.T) {
	synced := &Sync{NoteID: "n1", PageID: strPtr("1429989fe8ac4effbc8f57f56486db54")}

	tests := []struct {
		name string
		to   note.NoteStatus
		sync *Sync
		want Action
	}{
		{
			name: "[Create] first publish creates a page",
			to:   note.StatusPublish, sync: nil, want: ActionCreate,
		},
		{
			name: "[Create] page id is empty",
			to:   note.StatusPublish, sync: &Sync{NoteID: "n1", PageID: strPtr("")}, want: ActionCreate,
		},
		{
			name: "[Restore] publishing again restores the same page",
			to:   note.StatusPublish, sync: synced, want: ActionRestore,
		},
		{
			name: "[Trash] unpublishing trashes the page",
			to:   note.StatusDraft, sync: synced, want: ActionTrash,
		},
		{
			name: "[None] unpublishing a never-synced note does nothing",
			to:   note.StatusDraft, sync: nil, want: ActionNone,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ActionForStatusChange(tt.to, tt.sync); got != tt.want {
				t.Fatalf("ActionForStatusChange() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestActionForEdit(t *testing.T) {
	synced := &Sync{NoteID: "n1", PageID: strPtr("1429989fe8ac4effbc8f57f56486db54")}

	tests := []struct {
		name   string
		status note.NoteStatus
		sync   *Sync
		want   Action
	}{
		{
			name:   "[Update] editing a published synced note updates the page",
			status: note.StatusPublish, sync: synced, want: ActionUpdate,
		},
		{
			name:   "[None] editing a draft does nothing",
			status: note.StatusDraft, sync: synced, want: ActionNone,
		},
		{
			name:   "[None] published but never synced",
			status: note.StatusPublish, sync: nil, want: ActionNone,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ActionForEdit(tt.status, tt.sync); got != tt.want {
				t.Fatalf("ActionForEdit() = %v, want %v", got, tt.want)
			}
		})
	}
}
