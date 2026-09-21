package controller

import (
	openapi "immortal-architecture-bad-api/backend/internal/generated/openapi"
	"immortal-architecture-bad-api/backend/internal/service"
)

// Controller wires OpenAPI-generated routes to the bloated domain controllers.
type Controller struct {
	accountService  service.AccountServicer
	templateService service.TemplateServicer
	noteService     service.NoteServicer
}

var _ openapi.ServerInterface = (*Controller)(nil)

// NewController creates a new OpenAPI controller facade.
func NewController(accountSvc *service.AccountService, templateSvc *service.TemplateService, noteSvc *service.NoteService) *Controller {
	return &Controller{
		accountService:  accountSvc,
		templateService: templateSvc,
		noteService:     noteSvc,
	}
}
