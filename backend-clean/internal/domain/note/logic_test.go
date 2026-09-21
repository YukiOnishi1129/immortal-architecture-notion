package note

import (
	"errors"
	"testing"

	domainerr "immortal-architecture-cqrs/backend/internal/domain/errors"
	"immortal-architecture-cqrs/backend/internal/domain/template"
)

func TestValidateNoteOwnership(t *testing.T) {
	tests := []struct {
		name      string
		ownerID   string
		actorID   string
		wantError error
	}{
		{
			name:    "[Success] owner matches",
			ownerID: "owner-1",
			actorID: "owner-1",
		},
		{
			name:      "[Fail] owner mismatch",
			ownerID:   "owner-1",
			actorID:   "actor-2",
			wantError: domainerr.ErrUnauthorized,
		},
		{
			name:      "[Fail] missing owner",
			ownerID:   "",
			actorID:   "actor-2",
			wantError: domainerr.ErrOwnerRequired,
		},
		{
			name:      "[Fail] missing actor",
			ownerID:   "owner-1",
			actorID:   "",
			wantError: domainerr.ErrOwnerRequired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateNoteOwnership(tt.ownerID, tt.actorID)
			if tt.wantError == nil && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.wantError != nil && !errors.Is(err, tt.wantError) {
				t.Fatalf("want %v, got %v", tt.wantError, err)
			}
		})
	}
}

func TestCanChangeStatus(t *testing.T) {
	tests := []struct {
		name      string
		from      NoteStatus
		to        NoteStatus
		wantError error
	}{
		{
			name: "[Success] draft to publish",
			from: StatusDraft,
			to:   StatusPublish,
		},
		{
			name: "[Success] publish to draft",
			from: StatusPublish,
			to:   StatusDraft,
		},
		{
			name: "[Success] no change",
			from: StatusDraft,
			to:   StatusDraft,
		},
		{
			name:      "[Fail] invalid transition",
			from:      NoteStatus("Invalid"),
			to:        StatusDraft,
			wantError: domainerr.ErrInvalidStatusChange,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CanChangeStatus(tt.from, tt.to)
			if tt.wantError == nil && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.wantError != nil && !errors.Is(err, tt.wantError) {
				t.Fatalf("want %v, got %v", tt.wantError, err)
			}
		})
	}
}

func TestValidateSections(t *testing.T) {
	tplFields := []template.Field{
		{ID: "f1", Label: "Title", Order: 1, IsRequired: true},
		{ID: "f2", Label: "Body", Order: 2, IsRequired: false},
	}

	t.Run("[Success] valid sections", func(t *testing.T) {
		sections := []Section{
			{FieldID: "f1", Content: "hello"},
			{FieldID: "f2", Content: ""},
		}
		if err := ValidateSections(tplFields, sections); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("[Fail] required field empty", func(t *testing.T) {
		sections := []Section{
			{FieldID: "f1", Content: ""},
			{FieldID: "f2", Content: "body"},
		}
		if err := ValidateSections(tplFields, sections); !errors.Is(err, domainerr.ErrRequiredFieldEmpty) {
			t.Fatalf("expected ErrRequiredFieldEmpty, got %v", err)
		}
	})

	t.Run("[Fail] missing template field", func(t *testing.T) {
		sections := []Section{
			{FieldID: "unknown", Content: "x"},
		}
		if err := ValidateSections(tplFields, sections); !errors.Is(err, domainerr.ErrSectionsMissing) {
			t.Fatalf("expected ErrSectionsMissing, got %v", err)
		}
	})

	t.Run("[Fail] duplicate field", func(t *testing.T) {
		sections := []Section{
			{FieldID: "f1", Content: "a"},
			{FieldID: "f1", Content: "b"},
		}
		if err := ValidateSections(tplFields, sections); !errors.Is(err, domainerr.ErrSectionsMissing) {
			t.Fatalf("expected ErrSectionsMissing, got %v", err)
		}
	})
}

func TestValidateNoteForCreate(t *testing.T) {
	validTpl := template.Template{
		ID:      "tpl-1",
		Name:    "Template",
		OwnerID: "owner-1",
		Fields: []template.Field{
			{ID: "f1", Label: "Title", Order: 1, IsRequired: true},
		},
	}
	sections := []Section{{FieldID: "f1", Content: "value"}}

	tests := []struct {
		name      string
		title     string
		tpl       template.Template
		sections  []Section
		wantError error
	}{
		{
			name:     "[Success] valid",
			title:    "Hello",
			tpl:      validTpl,
			sections: sections,
		},
		{
			name:      "[Fail] empty title",
			title:     "  ",
			tpl:       validTpl,
			sections:  sections,
			wantError: domainerr.ErrTitleRequired,
		},
		{
			name:      "[Fail] missing template owner",
			title:     "Hello",
			tpl:       template.Template{ID: "tpl-1", Name: "Template", Fields: validTpl.Fields},
			sections:  sections,
			wantError: domainerr.ErrOwnerRequired,
		},
		{
			name:      "[Fail] sections invalid",
			title:     "Hello",
			tpl:       validTpl,
			sections:  []Section{{FieldID: "f1", Content: ""}},
			wantError: domainerr.ErrRequiredFieldEmpty,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateNoteForCreate(tt.title, tt.tpl, tt.sections)
			if tt.wantError == nil && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.wantError != nil && !errors.Is(err, tt.wantError) {
				t.Fatalf("want %v, got %v", tt.wantError, err)
			}
		})
	}
}
