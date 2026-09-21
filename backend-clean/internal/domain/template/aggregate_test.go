package template

import (
	"errors"
	"testing"

	domainerr "immortal-architecture-cqrs/backend/internal/domain/errors"
)

func TestTemplate_ReplaceFields(t *testing.T) {
	tests := []struct {
		name      string
		fields    []Field
		wantError error
	}{
		{
			name: "[Success] normalize and replace",
			fields: []Field{
				{ID: "f1", Label: "Title", Order: 0},
				{ID: "f2", Label: "Body", Order: 2},
			},
		},
		{
			name: "[Fail] invalid fields",
			fields: []Field{
				{ID: "f1", Label: "Title", Order: 1},
				{ID: "f2", Label: "", Order: 2},
			},
			wantError: domainerr.ErrFieldLabelRequired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tpl := Template{}
			err := tpl.ReplaceFields(tt.fields)
			if tt.wantError == nil && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.wantError != nil && !errors.Is(err, tt.wantError) {
				t.Fatalf("want %v, got %v", tt.wantError, err)
			}
			if err == nil && len(tpl.Fields) != len(tt.fields) {
				t.Fatalf("fields not replaced: %+v", tpl.Fields)
			}
		})
	}
}

