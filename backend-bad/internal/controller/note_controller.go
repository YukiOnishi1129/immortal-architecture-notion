package controller

import (
	"errors"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"

	openapi "immortal-architecture-bad-api/backend/internal/generated/openapi"
	"immortal-architecture-bad-api/backend/internal/service"
)

// NotesListNotes returns note list.
func (c *Controller) NotesListNotes(ctx echo.Context, params openapi.NotesListNotesParams) error {
	filters := service.NoteFilters{}
	if params.Q != nil {
		filters.Query = params.Q
	}
	if params.OwnerId != nil {
		filters.OwnerID = params.OwnerId
	}
	if params.TemplateId != nil {
		filters.TemplateID = params.TemplateId
	}
	if params.Status != nil {
		status := string(*params.Status)
		filters.Status = &status
	}

	notes, err := c.noteService.ListNotes(ctx, filters)
	if err != nil {
		return respondError(ctx, http.StatusInternalServerError, "failed to list notes")
	}
	return ctx.JSON(http.StatusOK, notes)
}

// NotesCreateNote creates note.
func (c *Controller) NotesCreateNote(ctx echo.Context) error {
	var body openapi.NotesCreateNoteJSONRequestBody
	if err := ctx.Bind(&body); err != nil {
		return respondError(ctx, http.StatusBadRequest, "invalid payload")
	}

	sections := map[string]string{}
	if body.Sections != nil {
		for _, section := range *body.Sections {
			sections[section.FieldId] = section.Content
		}
	}

	note, err := c.noteService.CreateNote(
		ctx,
		body.OwnerId.String(),
		body.TemplateId.String(),
		body.Title,
		sections,
	)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidAccountID), errors.Is(err, service.ErrInvalidTemplateID):
			return respondError(ctx, http.StatusBadRequest, err.Error())
		default:
			return respondError(ctx, http.StatusInternalServerError, "failed to create note")
		}
	}
	return ctx.JSON(http.StatusCreated, note)
}

// NotesDeleteNote deletes note.
func (c *Controller) NotesDeleteNote(ctx echo.Context, noteID string) error {
	if err := c.noteService.DeleteNote(ctx, noteID); err != nil {
		switch {
		case errors.Is(err, service.ErrNoteNotFound):
			return respondError(ctx, http.StatusNotFound, "note not found")
		case errors.Is(err, service.ErrInvalidNoteID):
			return respondError(ctx, http.StatusBadRequest, "invalid note id")
		default:
			return respondError(ctx, http.StatusInternalServerError, "failed to delete note")
		}
	}
	return ctx.JSON(http.StatusOK, openapi.ModelsSuccessResponse{Success: true})
}

// NotesGetNoteById returns note detail.
// revive:disable-next-line:var-naming // Method name fixed by generated interface.
func (c *Controller) NotesGetNoteById(ctx echo.Context, noteID string) error {
	note, err := c.noteService.GetNote(ctx, noteID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrNoteNotFound):
			return respondError(ctx, http.StatusNotFound, "note not found")
		case errors.Is(err, service.ErrInvalidNoteID):
			return respondError(ctx, http.StatusBadRequest, "invalid note id")
		default:
			return respondError(ctx, http.StatusInternalServerError, "failed to fetch note")
		}
	}
	return ctx.JSON(http.StatusOK, note)
}

// NotesUpdateNote updates note.
func (c *Controller) NotesUpdateNote(ctx echo.Context, noteID string) error {
	var body openapi.NotesUpdateNoteJSONRequestBody
	if err := ctx.Bind(&body); err != nil {
		return respondError(ctx, http.StatusBadRequest, "invalid payload")
	}
	if strings.TrimSpace(body.Title) == "" {
		return respondError(ctx, http.StatusBadRequest, "title is required")
	}
	sections := make(map[string]string, len(body.Sections))
	for _, section := range body.Sections {
		sections[section.Id] = section.Content
	}

	note, err := c.noteService.UpdateNote(ctx, noteID, body.Title, sections)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrNoteNotFound):
			return respondError(ctx, http.StatusNotFound, "note not found")
		case errors.Is(err, service.ErrInvalidNoteID):
			return respondError(ctx, http.StatusBadRequest, "invalid note id")
		default:
			return respondError(ctx, http.StatusInternalServerError, "failed to update note")
		}
	}
	return ctx.JSON(http.StatusOK, note)
}

// NotesPublishNote publishes note.
func (c *Controller) NotesPublishNote(ctx echo.Context, noteID string) error {
	return c.changeNoteStatus(ctx, noteID, "Publish")
}

// NotesUnpublishNote unpublishes note.
func (c *Controller) NotesUnpublishNote(ctx echo.Context, noteID string) error {
	return c.changeNoteStatus(ctx, noteID, "Draft")
}

func (c *Controller) changeNoteStatus(ctx echo.Context, noteID, status string) error {
	note, err := c.noteService.ChangeStatus(ctx, noteID, status)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrNoteNotFound):
			return respondError(ctx, http.StatusNotFound, "note not found")
		case errors.Is(err, service.ErrInvalidNoteID):
			return respondError(ctx, http.StatusBadRequest, "invalid note id")
		default:
			return respondError(ctx, http.StatusInternalServerError, "failed to change status")
		}
	}
	return ctx.JSON(http.StatusOK, note)
}
