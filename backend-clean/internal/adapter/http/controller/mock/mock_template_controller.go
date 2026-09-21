package mock

import (
	"context"

	"immortal-architecture-cqrs/backend/internal/domain/template"
	"immortal-architecture-cqrs/backend/internal/port"
)

// TemplateInputStub is a lightweight stub for template use case input.
type TemplateInputStub struct {
	Err    error
	Output port.TemplateOutputPort
}

func (s *TemplateInputStub) List(ctx context.Context, filters template.Filters) error { return s.Err }

func (s *TemplateInputStub) Get(ctx context.Context, id string) error {
	if s.Output != nil && s.Err == nil {
		_ = s.Output.PresentTemplate(ctx, &template.WithUsage{Template: template.Template{ID: id}})
	}
	return s.Err
}

func (s *TemplateInputStub) Create(ctx context.Context, input port.TemplateCreateInput) error {
	if s.Output != nil && s.Err == nil {
		_ = s.Output.PresentTemplate(ctx, &template.WithUsage{Template: template.Template{ID: "tpl-1", Name: input.Name, OwnerID: input.OwnerID}})
	}
	return s.Err
}

func (s *TemplateInputStub) Update(ctx context.Context, input port.TemplateUpdateInput) error {
	if s.Output != nil && s.Err == nil {
		_ = s.Output.PresentTemplate(ctx, &template.WithUsage{Template: template.Template{ID: input.ID, Name: input.Name, OwnerID: input.OwnerID}})
	}
	return s.Err
}

func (s *TemplateInputStub) Delete(ctx context.Context, id, ownerID string) error {
	if s.Output != nil && s.Err == nil {
		_ = s.Output.PresentTemplateDeleted(ctx)
	}
	return s.Err
}
