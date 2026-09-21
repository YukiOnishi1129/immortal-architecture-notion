// Package port defines application ports (interfaces).
package port

import (
	"context"

	"immortal-architecture-cqrs/backend/internal/domain/template"
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
	ReplaceFields(ctx context.Context, templateID string, fields []template.Field) error
}

// TemplateCreateInput is input for creating templates.
type TemplateCreateInput struct {
	Name    string
	OwnerID string
	Fields  []template.Field
}

// TemplateUpdateInput is input for updating templates.
type TemplateUpdateInput struct {
	ID      string
	Name    string
	Fields  []template.Field
	OwnerID string
}
