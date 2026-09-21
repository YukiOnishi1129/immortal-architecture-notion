// Package presenter contains HTTP presenters that implement output ports.
package presenter

import (
	"context"

	openapi "immortal-architecture-cqrs/backend/internal/adapter/http/generated/openapi"
	"immortal-architecture-cqrs/backend/internal/domain/template"
	"immortal-architecture-cqrs/backend/internal/port"
)

// TemplatePresenter converts template domain models to OpenAPI responses.
type TemplatePresenter struct {
	template *openapi.ModelsTemplateResponse
	list     []openapi.ModelsTemplateResponse
	deleted  bool
}

var _ port.TemplateOutputPort = (*TemplatePresenter)(nil)

// NewTemplatePresenter creates a TemplatePresenter.
func NewTemplatePresenter() *TemplatePresenter {
	return &TemplatePresenter{}
}

// PresentTemplateList stores template list response.
func (p *TemplatePresenter) PresentTemplateList(_ context.Context, templates []template.WithUsage) error {
	res := make([]openapi.ModelsTemplateResponse, 0, len(templates))
	for _, t := range templates {
		res = append(res, toTemplateResponse(t))
	}
	p.list = res
	return nil
}

// PresentTemplate stores single template response.
func (p *TemplatePresenter) PresentTemplate(_ context.Context, tpl *template.WithUsage) error {
	resp := toTemplateResponse(*tpl)
	p.template = &resp
	return nil
}

// PresentTemplateDeleted marks delete success.
func (p *TemplatePresenter) PresentTemplateDeleted(_ context.Context) error {
	p.deleted = true
	return nil
}

// Template returns the last template response.
func (p *TemplatePresenter) Template() *openapi.ModelsTemplateResponse {
	return p.template
}

// Templates returns the template list response.
func (p *TemplatePresenter) Templates() []openapi.ModelsTemplateResponse {
	return p.list
}

// DeleteResponse returns deletion success response.
func (p *TemplatePresenter) DeleteResponse() openapi.ModelsSuccessResponse {
	return openapi.ModelsSuccessResponse{Success: p.deleted}
}

func toTemplateResponse(t template.WithUsage) openapi.ModelsTemplateResponse {
	fields := make([]openapi.ModelsField, 0, len(t.Template.Fields))
	for _, f := range t.Template.Fields {
		fields = append(fields, openapi.ModelsField{
			Id:         f.ID,
			Label:      f.Label,
			Order:      int32(f.Order), //nolint:gosec
			IsRequired: f.IsRequired,
		})
	}
	return openapi.ModelsTemplateResponse{
		Id:      t.Template.ID,
		Name:    t.Template.Name,
		OwnerId: t.Template.OwnerID,
		Owner: openapi.ModelsAccountSummary{
			Id:        t.Owner.ID,
			FirstName: t.Owner.FirstName,
			LastName:  t.Owner.LastName,
			Thumbnail: t.Owner.Thumbnail,
		},
		Fields:    fields,
		IsUsed:    t.IsUsed,
		UpdatedAt: t.Template.UpdatedAt,
	}
}
