package controller

import (
	"errors"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"

	sqldb "immortal-architecture-bad-api/backend/internal/db/sqlc"
	openapi "immortal-architecture-bad-api/backend/internal/generated/openapi"
	"immortal-architecture-bad-api/backend/internal/service"
)

// TemplatesListTemplates returns template list.
func (c *Controller) TemplatesListTemplates(ctx echo.Context, params openapi.TemplatesListTemplatesParams) error {
	filters := service.TemplateFilters{}
	if params.OwnerId != nil {
		filters.OwnerID = params.OwnerId
	}
	if params.Q != nil {
		filters.Query = params.Q
	}

	templates, err := c.templateService.ListTemplates(ctx, filters)
	if err != nil {
		return respondError(ctx, http.StatusInternalServerError, "failed to list templates")
	}
	return ctx.JSON(http.StatusOK, templates)
}

// TemplatesCreateTemplate creates template.
func (c *Controller) TemplatesCreateTemplate(ctx echo.Context) error {
	var body openapi.TemplatesCreateTemplateJSONRequestBody
	if err := ctx.Bind(&body); err != nil {
		return respondError(ctx, http.StatusBadRequest, "invalid payload")
	}

	fields := make([]sqldb.Field, len(body.Fields))
	for i, field := range body.Fields {
		fields[i] = sqldb.Field{
			Label:      field.Label,
			Order:      field.Order,
			IsRequired: field.IsRequired,
		}
	}

	template, err := c.templateService.CreateTemplate(ctx, body.OwnerId.String(), body.Name, fields)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidAccountID):
			return respondError(ctx, http.StatusBadRequest, "invalid owner id")
		default:
			return respondError(ctx, http.StatusInternalServerError, "failed to create template")
		}
	}
	return ctx.JSON(http.StatusCreated, template)
}

// TemplatesDeleteTemplate deletes template.
func (c *Controller) TemplatesDeleteTemplate(ctx echo.Context, templateID string) error {
	if err := c.templateService.DeleteTemplate(ctx, templateID); err != nil {
		switch {
		case errors.Is(err, service.ErrTemplateNotFound):
			return respondError(ctx, http.StatusNotFound, "template not found")
		case errors.Is(err, service.ErrTemplateInUse):
			return respondError(ctx, http.StatusConflict, "template in use")
		case errors.Is(err, service.ErrInvalidTemplateID):
			return respondError(ctx, http.StatusBadRequest, "invalid template id")
		default:
			return respondError(ctx, http.StatusInternalServerError, "failed to delete template")
		}
	}
	return ctx.JSON(http.StatusOK, openapi.ModelsSuccessResponse{Success: true})
}

// TemplatesGetTemplateById returns template detail.
// revive:disable-next-line:var-naming // Method name fixed by generated interface.
func (c *Controller) TemplatesGetTemplateById(ctx echo.Context, templateID string) error {
	template, err := c.templateService.GetTemplate(ctx, templateID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrTemplateNotFound):
			return respondError(ctx, http.StatusNotFound, "template not found")
		case errors.Is(err, service.ErrInvalidTemplateID):
			return respondError(ctx, http.StatusBadRequest, "invalid template id")
		default:
			return respondError(ctx, http.StatusInternalServerError, "failed to fetch template")
		}
	}
	return ctx.JSON(http.StatusOK, template)
}

// TemplatesUpdateTemplate updates template name.
func (c *Controller) TemplatesUpdateTemplate(ctx echo.Context, templateID string) error {
	var body openapi.TemplatesUpdateTemplateJSONRequestBody
	if err := ctx.Bind(&body); err != nil {
		return respondError(ctx, http.StatusBadRequest, "invalid payload")
	}
	if strings.TrimSpace(body.Name) == "" {
		return respondError(ctx, http.StatusBadRequest, "name is required")
	}
	if len(body.Fields) == 0 {
		return respondError(ctx, http.StatusBadRequest, "at least one field is required")
	}

	template, err := c.templateService.UpdateTemplate(ctx, templateID, body.Name, body.Fields)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrTemplateNotFound):
			return respondError(ctx, http.StatusNotFound, "template not found")
		case errors.Is(err, service.ErrInvalidTemplateID):
			return respondError(ctx, http.StatusBadRequest, "invalid template id")
		case errors.Is(err, service.ErrTemplateInUse):
			return respondError(ctx, http.StatusConflict, "template is used by existing notes")
		default:
			return respondError(ctx, http.StatusInternalServerError, "failed to update template")
		}
	}
	return ctx.JSON(http.StatusOK, template)
}
