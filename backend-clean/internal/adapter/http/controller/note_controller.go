// Package controller contains HTTP controllers.
package controller

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"

	openapi "immortal-architecture-notion/backend/internal/adapter/http/generated/openapi"
	"immortal-architecture-notion/backend/internal/adapter/http/presenter"
	domainerr "immortal-architecture-notion/backend/internal/domain/errors"
	"immortal-architecture-notion/backend/internal/domain/note"
	"immortal-architecture-notion/backend/internal/port"
)

// NoteController handles note HTTP endpoints using CQRS pattern.
type NoteController struct {
	// Command (write) factories
	commandInputFactory  func(noteRepo port.NoteRepository, readModelRepo port.NoteReadModelRepository, tplRepo port.TemplateRepository, tx port.TxManager, output port.NoteCommandOutputPort) port.NoteCommandInputPort
	commandOutputFactory func() *presenter.NoteCommandPresenter

	// Query (read) factories
	queryInputFactory  func(readModelRepo port.NoteReadModelRepository, output port.NoteQueryOutputPort) port.NoteQueryInputPort
	queryOutputFactory func() *presenter.NoteQueryPresenter

	// Shared repository factories
	noteRepoFactory      func() port.NoteRepository
	readModelRepoFactory func() port.NoteReadModelRepository
	tplRepoFactory       func() port.TemplateRepository
	txFactory            func() port.TxManager
}

// NewNoteController creates NoteController.
func NewNoteController(
	commandInputFactory func(noteRepo port.NoteRepository, readModelRepo port.NoteReadModelRepository, tplRepo port.TemplateRepository, tx port.TxManager, output port.NoteCommandOutputPort) port.NoteCommandInputPort,
	commandOutputFactory func() *presenter.NoteCommandPresenter,
	queryInputFactory func(readModelRepo port.NoteReadModelRepository, output port.NoteQueryOutputPort) port.NoteQueryInputPort,
	queryOutputFactory func() *presenter.NoteQueryPresenter,
	noteRepoFactory func() port.NoteRepository,
	readModelRepoFactory func() port.NoteReadModelRepository,
	tplRepoFactory func() port.TemplateRepository,
	txFactory func() port.TxManager,
) *NoteController {
	return &NoteController{
		commandInputFactory:  commandInputFactory,
		commandOutputFactory: commandOutputFactory,
		queryInputFactory:    queryInputFactory,
		queryOutputFactory:   queryOutputFactory,
		noteRepoFactory:      noteRepoFactory,
		readModelRepoFactory: readModelRepoFactory,
		tplRepoFactory:       tplRepoFactory,
		txFactory:            txFactory,
	}
}

// List handles GET /notes (Query side).
func (c *NoteController) List(ctx echo.Context, params openapi.NotesListNotesParams) error {
	var status *note.NoteStatus
	if params.Status != nil {
		s := note.NoteStatus(*params.Status)
		status = &s
	}
	filters := note.Filters{
		Status:     status,
		TemplateID: params.TemplateId,
		OwnerID:    params.OwnerId,
		Query:      params.Q,
	}
	input, p := c.newQueryIO()
	if err := input.List(ctx.Request().Context(), filters); err != nil {
		return handleError(ctx, err)
	}
	return ctx.JSON(http.StatusOK, p.Notes())
}

// GetByID handles GET /notes/:id (Query side).
func (c *NoteController) GetByID(ctx echo.Context, noteID string) error {
	input, p := c.newQueryIO()
	if err := input.Get(ctx.Request().Context(), noteID); err != nil {
		return handleError(ctx, err)
	}
	return ctx.JSON(http.StatusOK, p.Note())
}

// Create handles POST /notes (Command side).
func (c *NoteController) Create(ctx echo.Context) error {
	var body openapi.ModelsCreateNoteRequest
	if err := ctx.Bind(&body); err != nil {
		return ctx.JSON(http.StatusBadRequest, openapi.ModelsBadRequestError{Code: openapi.ModelsBadRequestErrorCodeBADREQUEST, Message: "invalid body"})
	}
	ownerID := body.OwnerId.String()
	sections := []port.SectionInput{}
	if body.Sections != nil {
		for _, s := range *body.Sections {
			sections = append(sections, port.SectionInput{
				FieldID: s.FieldId,
				Content: s.Content,
			})
		}
	}
	input, p := c.newCommandIO()
	err := input.Create(ctx.Request().Context(), port.NoteCreateInput{
		Title:      body.Title,
		TemplateID: body.TemplateId.String(),
		OwnerID:    ownerID,
		Sections:   sections,
	})
	if err != nil {
		return handleError(ctx, err)
	}
	return ctx.JSON(http.StatusOK, p.Note())
}

// Update handles PUT /notes/:id (Command side).
func (c *NoteController) Update(ctx echo.Context, noteID string, params openapi.NotesUpdateNoteParams) error {
	var body openapi.ModelsUpdateNoteRequest
	if err := ctx.Bind(&body); err != nil {
		return ctx.JSON(http.StatusBadRequest, openapi.ModelsBadRequestError{Code: openapi.ModelsBadRequestErrorCodeBADREQUEST, Message: "invalid body"})
	}
	ownerID := strings.TrimSpace(params.OwnerId)
	if ownerID == "" {
		return handleError(ctx, domainerr.ErrUnauthorized)
	}
	sections := make([]port.SectionUpdateInput, 0, len(body.Sections))
	for _, s := range body.Sections {
		sections = append(sections, port.SectionUpdateInput{
			SectionID: s.Id,
			Content:   s.Content,
		})
	}
	input, p := c.newCommandIO()
	err := input.Update(ctx.Request().Context(), port.NoteUpdateInput{
		ID:       noteID,
		Title:    body.Title,
		OwnerID:  ownerID,
		Sections: sections,
	})
	if err != nil {
		return handleError(ctx, err)
	}
	return ctx.JSON(http.StatusOK, p.Note())
}

// Delete handles DELETE /notes/:id (Command side).
func (c *NoteController) Delete(ctx echo.Context, noteID string, params openapi.NotesDeleteNoteParams) error {
	ownerID := strings.TrimSpace(params.OwnerId)
	if ownerID == "" {
		return handleError(ctx, domainerr.ErrUnauthorized)
	}
	input, p := c.newCommandIO()
	if err := input.Delete(ctx.Request().Context(), noteID, ownerID); err != nil {
		return handleError(ctx, err)
	}
	return ctx.JSON(http.StatusOK, p.DeleteResponse())
}

// Publish handles POST /notes/:id/publish (Command side).
func (c *NoteController) Publish(ctx echo.Context, noteID string, params openapi.NotesPublishNoteParams) error {
	ownerID := strings.TrimSpace(params.OwnerId)
	if ownerID == "" {
		return handleError(ctx, domainerr.ErrOwnerRequired)
	}
	input, p := c.newCommandIO()
	err := input.ChangeStatus(ctx.Request().Context(), port.NoteStatusChangeInput{
		ID:      noteID,
		Status:  note.StatusPublish,
		OwnerID: ownerID,
	})
	if err != nil {
		return handleError(ctx, err)
	}
	return ctx.JSON(http.StatusOK, p.Note())
}

// SyncToNotion handles POST /notes/:id/notion-sync (Command side).
func (c *NoteController) SyncToNotion(ctx echo.Context, noteID string, params openapi.NotesSyncNoteToNotionParams) error {
	ownerID := strings.TrimSpace(params.OwnerId)
	if ownerID == "" {
		return handleError(ctx, domainerr.ErrOwnerRequired)
	}
	input, p := c.newCommandIO()
	if err := input.SyncToNotion(ctx.Request().Context(), noteID, ownerID); err != nil {
		return handleError(ctx, err)
	}
	return ctx.JSON(http.StatusOK, p.Note())
}

// Unpublish handles POST /notes/:id/unpublish (Command side).
func (c *NoteController) Unpublish(ctx echo.Context, noteID string, params openapi.NotesUnpublishNoteParams) error {
	ownerID := strings.TrimSpace(params.OwnerId)
	if ownerID == "" {
		return handleError(ctx, domainerr.ErrOwnerRequired)
	}
	input, p := c.newCommandIO()
	err := input.ChangeStatus(ctx.Request().Context(), port.NoteStatusChangeInput{
		ID:      noteID,
		Status:  note.StatusDraft,
		OwnerID: ownerID,
	})
	if err != nil {
		return handleError(ctx, err)
	}
	return ctx.JSON(http.StatusOK, p.Note())
}

func (c *NoteController) newCommandIO() (port.NoteCommandInputPort, *presenter.NoteCommandPresenter) {
	output := c.commandOutputFactory()
	input := c.commandInputFactory(c.noteRepoFactory(), c.readModelRepoFactory(), c.tplRepoFactory(), c.txFactory(), output)
	return input, output
}

func (c *NoteController) newQueryIO() (port.NoteQueryInputPort, *presenter.NoteQueryPresenter) {
	output := c.queryOutputFactory()
	input := c.queryInputFactory(c.readModelRepoFactory(), output)
	return input, output
}
