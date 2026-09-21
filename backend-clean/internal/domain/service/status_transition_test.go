package service

import (
	"errors"
	"testing"

	domainerr "immortal-architecture-cqrs/backend/internal/domain/errors"
	"immortal-architecture-cqrs/backend/internal/domain/note"
)

func TestCanPublish(t *testing.T) {
	tests := []struct {
		name      string
		note      note.Note
		actorID   string
		wantError error
	}{
		{
			name:    "[Success] owner can publish draft",
			note:    note.Note{ID: "n1", OwnerID: "owner-1", Status: note.StatusDraft},
			actorID: "owner-1",
		},
		{
			name:      "[Fail] unauthorized actor",
			note:      note.Note{ID: "n1", OwnerID: "owner-1", Status: note.StatusDraft},
			actorID:   "other",
			wantError: domainerr.ErrUnauthorized,
		},
		{
			name:      "[Fail] invalid status value",
			note:      note.Note{ID: "n1", OwnerID: "owner-1", Status: note.NoteStatus("Invalid")},
			actorID:   "owner-1",
			wantError: domainerr.ErrInvalidStatus,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CanPublish(tt.note, tt.actorID)
			if tt.wantError == nil && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.wantError != nil && !errors.Is(err, tt.wantError) {
				t.Fatalf("want %v, got %v", tt.wantError, err)
			}
		})
	}
}

func TestCanUnpublish(t *testing.T) {
	tests := []struct {
		name      string
		note      note.Note
		actorID   string
		wantError error
	}{
		{
			name:    "[Success] owner can unpublish publish",
			note:    note.Note{ID: "n1", OwnerID: "owner-1", Status: note.StatusPublish},
			actorID: "owner-1",
		},
		{
			name:      "[Fail] unauthorized actor",
			note:      note.Note{ID: "n1", OwnerID: "owner-1", Status: note.StatusPublish},
			actorID:   "other",
			wantError: domainerr.ErrUnauthorized,
		},
		{
			name:      "[Fail] invalid status value",
			note:      note.Note{ID: "n1", OwnerID: "owner-1", Status: note.NoteStatus("Invalid")},
			actorID:   "owner-1",
			wantError: domainerr.ErrInvalidStatus,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CanUnpublish(tt.note, tt.actorID)
			if tt.wantError == nil && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.wantError != nil && !errors.Is(err, tt.wantError) {
				t.Fatalf("want %v, got %v", tt.wantError, err)
			}
		})
	}
}
