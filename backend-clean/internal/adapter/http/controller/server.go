package controller

import (
	"github.com/labstack/echo/v4"

	openapi "immortal-architecture-cqrs/backend/internal/adapter/http/generated/openapi"
)

// Server implements the OpenAPI ServerInterface by delegating to domain-specific controllers.
type Server struct {
	account  *AccountController
	note     *NoteController
	template *TemplateController
}

// NewServer wires controller dependencies to generated ServerInterface.
func NewServer(ac *AccountController, nc *NoteController, tc *TemplateController) *Server {
	return &Server{account: ac, note: nc, template: tc}
}

// AccountsCreateOrGetAccount handles POST /api/accounts/auth.
func (s *Server) AccountsCreateOrGetAccount(ctx echo.Context) error {
	return s.account.CreateOrGet(ctx)
}

// AccountsGetCurrentAccount handles GET /api/accounts/me.
func (s *Server) AccountsGetCurrentAccount(ctx echo.Context) error {
	return s.account.GetCurrent(ctx)
}

// AccountsGetAccountById handles GET /api/accounts/:id.
func (s *Server) AccountsGetAccountById(ctx echo.Context, accountId string) error { //nolint:revive
	return s.account.GetByID(ctx, accountId)
}

// AccountsGetAccountByEmail handles GET /api/accounts/by-email.
func (s *Server) AccountsGetAccountByEmail(ctx echo.Context, params openapi.AccountsGetAccountByEmailParams) error {
	return s.account.GetAccountByEmail(ctx, params)
}

// NotesListNotes handles GET /api/notes.
func (s *Server) NotesListNotes(ctx echo.Context, params openapi.NotesListNotesParams) error {
	return s.note.List(ctx, params)
}

// NotesCreateNote handles POST /api/notes.
func (s *Server) NotesCreateNote(ctx echo.Context) error {
	return s.note.Create(ctx)
}

// NotesDeleteNote handles DELETE /api/notes/:id.
func (s *Server) NotesDeleteNote(ctx echo.Context, noteId string, params openapi.NotesDeleteNoteParams) error { //nolint:revive
	return s.note.Delete(ctx, noteId, params)
}

// NotesGetNoteById handles GET /api/notes/:id.
func (s *Server) NotesGetNoteById(ctx echo.Context, noteId string) error { //nolint:revive
	return s.note.GetByID(ctx, noteId)
}

// NotesUpdateNote handles PUT /api/notes/:noteId.
// NotesUpdateNote handles PUT /api/notes/:id.
func (s *Server) NotesUpdateNote(ctx echo.Context, noteId string, params openapi.NotesUpdateNoteParams) error { //nolint:revive
	return s.note.Update(ctx, noteId, params)
}

// NotesPublishNote handles POST /api/notes/:noteId/publish.
// NotesPublishNote handles POST /api/notes/:id/publish.
func (s *Server) NotesPublishNote(ctx echo.Context, noteId string, params openapi.NotesPublishNoteParams) error { //nolint:revive
	return s.note.Publish(ctx, noteId, params)
}

// NotesUnpublishNote handles POST /api/notes/:noteId/unpublish.
// NotesUnpublishNote handles POST /api/notes/:id/unpublish.
func (s *Server) NotesUnpublishNote(ctx echo.Context, noteId string, params openapi.NotesUnpublishNoteParams) error { //nolint:revive
	return s.note.Unpublish(ctx, noteId, params)
}

// TemplatesListTemplates handles GET /api/templates.
func (s *Server) TemplatesListTemplates(ctx echo.Context, params openapi.TemplatesListTemplatesParams) error {
	return s.template.List(ctx, params)
}

// TemplatesCreateTemplate handles POST /api/templates.
func (s *Server) TemplatesCreateTemplate(ctx echo.Context) error {
	return s.template.Create(ctx)
}

// TemplatesDeleteTemplate handles DELETE /api/templates/:id.
func (s *Server) TemplatesDeleteTemplate(ctx echo.Context, templateId string, params openapi.TemplatesDeleteTemplateParams) error { //nolint:revive
	return s.template.Delete(ctx, templateId, params)
}

// TemplatesGetTemplateById handles GET /api/templates/:id.
func (s *Server) TemplatesGetTemplateById(ctx echo.Context, templateId string) error { //nolint:revive
	return s.template.GetByID(ctx, templateId)
}

// TemplatesUpdateTemplate handles PUT /api/templates/:templateId.
// TemplatesUpdateTemplate handles PUT /api/templates/:id.
func (s *Server) TemplatesUpdateTemplate(ctx echo.Context, templateId string, params openapi.TemplatesUpdateTemplateParams) error { //nolint:revive
	return s.template.Update(ctx, templateId, params)
}
