// Package controller contains HTTP controllers.
package controller

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"

	openapi "immortal-architecture-cqrs/backend/internal/adapter/http/generated/openapi"
	"immortal-architecture-cqrs/backend/internal/adapter/http/presenter"
	domainerr "immortal-architecture-cqrs/backend/internal/domain/errors"
	"immortal-architecture-cqrs/backend/internal/domain/template"
	"immortal-architecture-cqrs/backend/internal/port"
)

// TemplateController handles template HTTP endpoints.
type TemplateController struct {
	inputFactory  func(repo port.TemplateRepository, tx port.TxManager, output port.TemplateOutputPort) port.TemplateInputPort
	outputFactory func() *presenter.TemplatePresenter
	repoFactory   func() port.TemplateRepository
	txFactory     func() port.TxManager
}

// NewTemplateController creates TemplateController.
func NewTemplateController(
	inputFactory func(repo port.TemplateRepository, tx port.TxManager, output port.TemplateOutputPort) port.TemplateInputPort,
	outputFactory func() *presenter.TemplatePresenter,
	repoFactory func() port.TemplateRepository,
	txFactory func() port.TxManager,
) *TemplateController {
	return &TemplateController{
		inputFactory:  inputFactory,
		outputFactory: outputFactory,
		repoFactory:   repoFactory,
		txFactory:     txFactory,
	}
}

// List handles GET /templates.
func (c *TemplateController) List(ctx echo.Context, params openapi.TemplatesListTemplatesParams) error {
	filters := template.Filters{
		Query:   params.Q,
		OwnerID: params.OwnerId,
	}
	input, p := c.newIO()
	if err := input.List(ctx.Request().Context(), filters); err != nil {
		return handleError(ctx, err)
	}
	return ctx.JSON(http.StatusOK, p.Templates())
}

// GetByID handles GET /templates/:id.
func (c *TemplateController) GetByID(ctx echo.Context, templateID string) error {
	input, p := c.newIO()
	if err := input.Get(ctx.Request().Context(), templateID); err != nil {
		return handleError(ctx, err)
	}
	return ctx.JSON(http.StatusOK, p.Template())
}

// Create handles POST /templates.
func (c *TemplateController) Create(ctx echo.Context) error {
	var body openapi.ModelsCreateTemplateRequest
	if err := ctx.Bind(&body); err != nil {
		return ctx.JSON(http.StatusBadRequest, openapi.ModelsBadRequestError{Code: openapi.ModelsBadRequestErrorCodeBADREQUEST, Message: "invalid body"})
	}
	ownerID := body.OwnerId.String()
	fields := make([]template.Field, 0, len(body.Fields))
	for _, f := range body.Fields {
		fields = append(fields, template.Field{
			Label:      f.Label,
			Order:      int(f.Order),
			IsRequired: f.IsRequired,
		})
	}
	input, p := c.newIO()
	err := input.Create(ctx.Request().Context(), port.TemplateCreateInput{
		Name:    body.Name,
		OwnerID: ownerID,
		Fields:  fields,
	})
	if err != nil {
		return handleError(ctx, err)
	}
	return ctx.JSON(http.StatusOK, p.Template())
}

// Update handles PUT /templates/:id.
func (c *TemplateController) Update(ctx echo.Context, templateID string, params openapi.TemplatesUpdateTemplateParams) error {
	var body openapi.ModelsUpdateTemplateRequest
	if err := ctx.Bind(&body); err != nil {
		return ctx.JSON(http.StatusBadRequest, openapi.ModelsBadRequestError{Code: openapi.ModelsBadRequestErrorCodeBADREQUEST, Message: "invalid body"})
	}
	ownerID := strings.TrimSpace(params.OwnerId)
	if ownerID == "" {
		return handleError(ctx, domainerr.ErrUnauthorized)
	}
	fields := make([]template.Field, 0, len(body.Fields))
	for _, f := range body.Fields {
		fields = append(fields, template.Field{
			ID:         valueOrEmpty(f.Id),
			Label:      f.Label,
			Order:      int(f.Order),
			IsRequired: f.IsRequired,
		})
	}
	input, p := c.newIO()
	err := input.Update(ctx.Request().Context(), port.TemplateUpdateInput{
		ID:      templateID,
		Name:    body.Name,
		Fields:  fields,
		OwnerID: ownerID,
	})
	if err != nil {
		return handleError(ctx, err)
	}
	return ctx.JSON(http.StatusOK, p.Template())
}

// Delete handles DELETE /templates/:id.
func (c *TemplateController) Delete(ctx echo.Context, templateID string, params openapi.TemplatesDeleteTemplateParams) error {
	ownerID := strings.TrimSpace(params.OwnerId)
	if ownerID == "" {
		return handleError(ctx, domainerr.ErrUnauthorized)
	}
	input, p := c.newIO()
	if err := input.Delete(ctx.Request().Context(), templateID, ownerID); err != nil {
		return handleError(ctx, err)
	}
	return ctx.JSON(http.StatusOK, p.DeleteResponse())
}

func (c *TemplateController) newIO() (port.TemplateInputPort, *presenter.TemplatePresenter) {
	output := c.outputFactory()
	input := c.inputFactory(c.repoFactory(), c.txFactory(), output)
	return input, output
}
