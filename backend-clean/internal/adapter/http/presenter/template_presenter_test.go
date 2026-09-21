package presenter

import (
	"context"
	"testing"
	"time"

	"immortal-architecture-cqrs/backend/internal/domain/template"
)

func TestTemplatePresenter_TableDriven(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name       string
		action     string
		single     *template.WithUsage
		list       []template.WithUsage
		wantID     string
		wantOwner  string
		wantCount  int
		expectUsed bool
	}{
		{
			name:   "[Success] single",
			action: "single",
			single: &template.WithUsage{
				Template: template.Template{
					ID:      "tpl-1",
					Name:    "Template",
					OwnerID: "owner-1",
					Fields:  []template.Field{{ID: "f1", Label: "Title", Order: 2, IsRequired: true}},
					UpdatedAt: now,
				},
				Owner:  template.Owner{ID: "owner-1", FirstName: "Taro", LastName: "Yamada"},
				IsUsed: true,
			},
			wantID:     "tpl-1",
			wantOwner:  "owner-1",
			expectUsed: true,
		},
		{
			name:      "[Success] list",
			action:    "list",
			list:      []template.WithUsage{{Template: template.Template{ID: "tpl-1"}}, {Template: template.Template{ID: "tpl-2"}}},
			wantCount: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewTemplatePresenter()
			var err error
			switch tt.action {
			case "single":
				err = p.PresentTemplate(context.Background(), tt.single)
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				resp := p.Template()
				if resp == nil || resp.Id != tt.wantID || resp.OwnerId != tt.wantOwner {
					t.Fatalf("unexpected response: %+v", resp)
				}
				if len(resp.Fields) != 1 || resp.Fields[0].Order != 2 {
					t.Fatalf("fields not converted correctly: %+v", resp.Fields)
				}
				if resp.IsUsed != tt.expectUsed {
					t.Fatalf("IsUsed mismatch")
				}
				if resp.UpdatedAt.IsZero() {
					t.Fatalf("UpdatedAt not set")
				}
			case "list":
				err = p.PresentTemplateList(context.Background(), tt.list)
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if len(p.Templates()) != tt.wantCount {
					t.Fatalf("want %d templates, got %d", tt.wantCount, len(p.Templates()))
				}
			}
		})
	}
}

func TestTemplatePresenter_PresentTemplateDeleted(t *testing.T) {
	p := NewTemplatePresenter()
	_ = p.PresentTemplateDeleted(context.Background())
	if !p.DeleteResponse().Success {
		t.Fatalf("delete flag not set")
	}
}
