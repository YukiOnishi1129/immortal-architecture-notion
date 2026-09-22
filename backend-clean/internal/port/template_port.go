// Package port defines application ports (interfaces).
package port

import (
	"context"

	"immortal-architecture-notion/backend/internal/domain/template"
)

// TemplateInputPort defines template use case inputs.
type TemplateInputPort interface {
	List(ctx context.Context, filters template.Filters) error
	Get(ctx context.Context, id string) error
	Create(ctx context.Context, input TemplateCreateInput) error
	Update(ctx context.Context, input TemplateUpdateInput) error
	Delete(ctx context.Context, id, ownerID string) error
}

// TemplateOutputPort defines template presenters.
type TemplateOutputPort interface {
	PresentTemplateList(ctx context.Context, templates []template.WithUsage) error
	PresentTemplate(ctx context.Context, template *template.WithUsage) error
	PresentTemplateDeleted(ctx context.Context) error
}

// TemplateRepository abstracts template persistence.
type TemplateRepository interface {
	List(ctx context.Context, filters template.Filters) ([]template.WithUsage, error)
	Get(ctx context.Context, id string) (*template.WithUsage, error)
	Create(ctx context.Context, tpl template.Template) (*template.Template, error)
	Update(ctx context.Context, tpl template.Template) (*template.Template, error)
	Delete(ctx context.Context, id string) error
	// ReplaceFields deletes every field and writes the given ones back.
	// Only valid while no note references them.
	ReplaceFields(ctx context.Context, templateID string, fields []template.Field) error

	// SyncFields updates the fields in place, keeping the ids notes point at.
	SyncFields(ctx context.Context, templateID string, fields []template.Field) error

	// UsedFieldIDs returns the ids of fields that notes already store content
	// for. Those fields cannot be renamed or removed.
	UsedFieldIDs(ctx context.Context, templateID string) ([]string, error)
}

// TemplateCreateInput is input for creating templates.
type TemplateCreateInput struct {
	Name    string
	OwnerID string
	Fields  []template.Field

	// NotionParentPageURL is the Notion page notes will be created under.
	// Empty means the template is not linked to Notion.
	NotionParentPageURL string
}

// TemplateUpdateInput is input for updating templates.
type TemplateUpdateInput struct {
	ID      string
	Name    string
	Fields  []template.Field
	OwnerID string

	// NotionParentPageURL is the Notion page notes will be created under.
	// nil means the field was not sent and the current link is kept;
	// an empty string is an explicit request to clear it.
	NotionParentPageURL *string
}
