package template

import (
	"errors"
	"testing"

	domainerr "immortal-architecture-cqrs/backend/internal/domain/errors"
)

func TestNormalizeAndValidate(t *testing.T) {
	tests := []struct {
		name      string
		fields    []Field
		wantOrder []int
		wantError error
	}{
		{
			name: "[Success] fills order when zero",
			fields: []Field{
				{ID: "f1", Label: "Title", Order: 0, IsRequired: true},
				{ID: "f2", Label: "Body", Order: 0, IsRequired: false},
			},
			wantOrder: []int{1, 2},
		},
		{
			name:      "[Fail] empty fields",
			fields:    nil,
			wantError: domainerr.ErrFieldRequired,
		},
		{
			name: "[Fail] missing label",
			fields: []Field{
				{ID: "f1", Label: "", Order: 1},
			},
			wantError: domainerr.ErrFieldLabelRequired,
		},
		{
			name: "[Fail] duplicate order",
			fields: []Field{
				{ID: "f1", Label: "Title", Order: 1},
				{ID: "f2", Label: "Body", Order: 1},
			},
			wantError: domainerr.ErrFieldOrderInvalid,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, err := NormalizeAndValidate(tt.fields)
			if tt.wantError == nil && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.wantError != nil && !errors.Is(err, tt.wantError) {
				t.Fatalf("want %v, got %v", tt.wantError, err)
			}
			if tt.wantError == nil && tt.wantOrder != nil {
				if len(out) != len(tt.wantOrder) {
					t.Fatalf("unexpected length: %d", len(out))
				}
				for i, f := range out {
					if f.Order != tt.wantOrder[i] {
						t.Fatalf("order mismatch at %d: want %d, got %d", i, tt.wantOrder[i], f.Order)
					}
				}
			}
		})
	}
}

func TestValidateTemplate(t *testing.T) {
	valid := Template{
		ID:      "tpl-1",
		Name:    "Template",
		OwnerID: "owner-1",
		Fields: []Field{
			{ID: "f1", Label: "Title", Order: 1},
		},
	}

	tests := []struct {
		name      string
		tpl       Template
		wantError error
	}{
		{
			name: "[Success] valid template",
			tpl:  valid,
		},
		{
			name:      "[Fail] missing name",
			tpl:       Template{OwnerID: "owner-1", Fields: valid.Fields},
			wantError: domainerr.ErrTemplateNameRequired,
		},
		{
			name:      "[Fail] missing owner",
			tpl:       Template{Name: "Template", Fields: valid.Fields},
			wantError: domainerr.ErrTemplateOwnerRequired,
		},
		{
			name: "[Fail] invalid fields",
			tpl: Template{
				Name:    "Template",
				OwnerID: "owner-1",
				Fields: []Field{
					{ID: "f1", Label: "Title", Order: 1},
					{ID: "f2", Label: "DupOrder", Order: 1},
				},
			},
			wantError: domainerr.ErrFieldOrderInvalid,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTemplate(tt.tpl)
			if tt.wantError == nil && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.wantError != nil && !errors.Is(err, tt.wantError) {
				t.Fatalf("want %v, got %v", tt.wantError, err)
			}
		})
	}
}

func TestCanDeleteTemplate(t *testing.T) {
	tests := []struct {
		name      string
		isUsed    bool
		wantError error
	}{
		{name: "[Success] not used", isUsed: false},
		{name: "[Fail] used", isUsed: true, wantError: domainerr.ErrTemplateInUse},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CanDeleteTemplate(tt.isUsed)
			if tt.wantError == nil && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.wantError != nil && !errors.Is(err, tt.wantError) {
				t.Fatalf("want %v, got %v", tt.wantError, err)
			}
		})
	}
}

func TestValidateTemplateOwnership(t *testing.T) {
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
			actorID:   "other",
			wantError: domainerr.ErrUnauthorized,
		},
		{
			name:      "[Fail] empty owner",
			ownerID:   "",
			actorID:   "actor-1",
			wantError: domainerr.ErrTemplateOwnerRequired,
		},
		{
			name:      "[Fail] empty actor",
			ownerID:   "owner-1",
			actorID:   "",
			wantError: domainerr.ErrTemplateOwnerRequired,
		},
		{
			name:      "[Fail] whitespace owner",
			ownerID:   "  ",
			actorID:   "actor-1",
			wantError: domainerr.ErrTemplateOwnerRequired,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTemplateOwnership(tt.ownerID, tt.actorID)
			if tt.wantError == nil && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.wantError != nil && !errors.Is(err, tt.wantError) {
				t.Fatalf("want %v, got %v", tt.wantError, err)
			}
		})
	}
}
