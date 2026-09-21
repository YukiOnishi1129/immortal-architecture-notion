package note

import (
	"testing"

	domainerr "immortal-architecture-cqrs/backend/internal/domain/errors"
)

func TestNote_ReplaceSections(t *testing.T) {
	tests := []struct {
		name      string
		noteID    string
		sections  []Section
		wantError error
	}{
		{
			name:   "[Success] assigns note id when empty",
			noteID: "note-1",
			sections: []Section{
				{FieldID: "field-1", Content: "c1"},
				{FieldID: "field-2", Content: "c2"},
			},
		},
		{
			name:   "[Success] keeps existing matching note id",
			noteID: "note-1",
			sections: []Section{
				{NoteID: "note-1", FieldID: "field-1", Content: "c1"},
			},
		},
		{
			name:   "[Fail] mismatch note id",
			noteID: "note-1",
			sections: []Section{
				{NoteID: "other", FieldID: "field-1", Content: "c1"},
			},
			wantError: domainerr.ErrUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := Note{ID: tt.noteID}
			err := n.ReplaceSections(tt.sections)
			if tt.wantError == nil && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.wantError != nil && err != tt.wantError {
				t.Fatalf("want %v, got %v", tt.wantError, err)
			}
			if err == nil {
				for _, sec := range n.Sections {
					if sec.NoteID != tt.noteID {
						t.Fatalf("section noteID not set: %s", sec.NoteID)
					}
				}
			}
		})
	}
}
