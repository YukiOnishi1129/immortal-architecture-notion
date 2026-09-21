package service

import (
	"context"
	"errors"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"

	sqldb "immortal-architecture-bad-api/backend/internal/db/sqlc"
	openapi "immortal-architecture-bad-api/backend/internal/generated/openapi"
)

var (
	// ErrNoteNotFound indicates note missing.
	ErrNoteNotFound = errors.New("note not found")
	// ErrInvalidNoteID invalid UUID.
	ErrInvalidNoteID = errors.New("invalid note id")
)

// NoteService bundles note logic.
type NoteService struct {
	pool    *pgxpool.Pool
	queries *sqldb.Queries
}

// NewNoteService creates service.
func NewNoteService(pool *pgxpool.Pool) *NoteService {
	return &NoteService{
		pool:    pool,
		queries: sqldb.New(pool),
	}
}

// NoteFilters for listing notes.
type NoteFilters struct {
	Status     *string
	TemplateID *string
	OwnerID    *string
	Query      *string
}

// ListNotes returns note list with filters.
func (s *NoteService) ListNotes(ctx echo.Context, filters NoteFilters) ([]*openapi.ModelsNoteResponse, error) {
	dbCtx := ctx.Request().Context()
	params := &sqldb.ListNotesParams{}

	if filters.Status != nil && *filters.Status != "" {
		params.Column1 = *filters.Status
	}
	if filters.TemplateID != nil && *filters.TemplateID != "" {
		if id, err := parseUUID(*filters.TemplateID); err == nil {
			params.Column2 = id
		}
	}
	if filters.OwnerID != nil && *filters.OwnerID != "" {
		if id, err := parseUUID(*filters.OwnerID); err == nil {
			params.Column3 = id
		}
	}
	if filters.Query != nil && *filters.Query != "" {
		params.Column4 = *filters.Query
	}

	rows, err := s.queries.ListNotes(dbCtx, params)
	if err != nil {
		return nil, err
	}

	result := make([]*openapi.ModelsNoteResponse, 0, len(rows))
	for _, row := range rows {
		res, err := s.composeNoteResponse(dbCtx, noteResponseInput{
			ID:           row.ID,
			Title:        row.Title,
			TemplateID:   row.TemplateID,
			OwnerID:      row.OwnerID,
			Status:       row.Status,
			CreatedAt:    row.CreatedAt,
			UpdatedAt:    row.UpdatedAt,
			TemplateName: row.TemplateName,
		})
		if err != nil {
			return nil, err
		}
		result = append(result, res)
	}
	return result, nil
}

// GetNote returns a note by ID.
func (s *NoteService) GetNote(ctx echo.Context, id string) (*openapi.ModelsNoteResponse, error) {
	noteID, err := parseUUID(id)
	if err != nil {
		return nil, ErrInvalidNoteID
	}

	note, err := s.queries.GetNoteByID(ctx.Request().Context(), noteID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNoteNotFound
		}
		return nil, err
	}

	return s.composeNoteResponse(ctx.Request().Context(), noteResponseInput{
		ID:           note.ID,
		Title:        note.Title,
		TemplateID:   note.TemplateID,
		OwnerID:      note.OwnerID,
		Status:       note.Status,
		CreatedAt:    note.CreatedAt,
		UpdatedAt:    note.UpdatedAt,
		TemplateName: note.TemplateName,
	})
}

// CreateNote creates note and sections.
func (s *NoteService) CreateNote(ctx echo.Context, ownerID, templateID, title string, sections map[string]string) (*openapi.ModelsNoteResponse, error) {
	ownerUUID, err := parseUUID(ownerID)
	if err != nil {
		return nil, ErrInvalidAccountID
	}
	templateUUID, err := parseUUID(templateID)
	if err != nil {
		return nil, ErrInvalidTemplateID
	}

	tx, err := s.pool.Begin(ctx.Request().Context())
	if err != nil {
		return nil, err
	}
	txQueries := s.queries.WithTx(tx)

	note, err := txQueries.CreateNote(ctx.Request().Context(), &sqldb.CreateNoteParams{
		Title:      strings.TrimSpace(title),
		TemplateID: templateUUID,
		OwnerID:    ownerUUID,
		Status:     "Draft",
	})
	if err != nil {
		if rbErr := tx.Rollback(ctx.Request().Context()); rbErr != nil {
			return nil, rbErr
		}
		return nil, err
	}

	for fieldID, content := range sections {
		fid, parseErr := parseUUID(fieldID)
		if parseErr != nil {
			if rbErr := tx.Rollback(ctx.Request().Context()); rbErr != nil {
				return nil, rbErr
			}
			return nil, parseErr
		}
		if _, err = txQueries.CreateSection(ctx.Request().Context(), &sqldb.CreateSectionParams{
			NoteID:  note.ID,
			FieldID: fid,
			Content: content,
		}); err != nil {
			if rbErr := tx.Rollback(ctx.Request().Context()); rbErr != nil {
				return nil, rbErr
			}
			return nil, err
		}
	}

	if err := tx.Commit(ctx.Request().Context()); err != nil {
		if rbErr := tx.Rollback(ctx.Request().Context()); rbErr != nil {
			return nil, rbErr
		}
		return nil, err
	}

	return s.composeNoteResponse(ctx.Request().Context(), noteResponseInput{
		ID:         note.ID,
		Title:      note.Title,
		TemplateID: note.TemplateID,
		OwnerID:    note.OwnerID,
		Status:     note.Status,
		CreatedAt:  note.CreatedAt,
		UpdatedAt:  note.UpdatedAt,
	})
}

// UpdateNote updates title and sections.
func (s *NoteService) UpdateNote(ctx echo.Context, id, title string, sections map[string]string) (*openapi.ModelsNoteResponse, error) {
	noteID, err := parseUUID(id)
	if err != nil {
		return nil, ErrInvalidNoteID
	}

	tx, err := s.pool.Begin(ctx.Request().Context())
	if err != nil {
		return nil, err
	}
	txQueries := s.queries.WithTx(tx)

	note, err := txQueries.UpdateNote(ctx.Request().Context(), &sqldb.UpdateNoteParams{
		ID:    noteID,
		Title: strings.TrimSpace(title),
	})
	if err != nil {
		if rbErr := tx.Rollback(ctx.Request().Context()); rbErr != nil {
			return nil, rbErr
		}
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNoteNotFound
		}
		return nil, err
	}

	for sectionID, content := range sections {
		secUUID, parseErr := parseUUID(sectionID)
		if parseErr != nil {
			if rbErr := tx.Rollback(ctx.Request().Context()); rbErr != nil {
				return nil, rbErr
			}
			return nil, parseErr
		}
		if _, err := txQueries.UpdateSectionContent(ctx.Request().Context(), &sqldb.UpdateSectionContentParams{
			ID:      secUUID,
			Content: content,
		}); err != nil {
			if rbErr := tx.Rollback(ctx.Request().Context()); rbErr != nil {
				return nil, rbErr
			}
			return nil, err
		}
	}

	if err := tx.Commit(ctx.Request().Context()); err != nil {
		if rbErr := tx.Rollback(ctx.Request().Context()); rbErr != nil {
			return nil, rbErr
		}
		return nil, err
	}

	return s.composeNoteResponse(ctx.Request().Context(), noteResponseInput{
		ID:         note.ID,
		Title:      note.Title,
		TemplateID: note.TemplateID,
		OwnerID:    note.OwnerID,
		Status:     note.Status,
		CreatedAt:  note.CreatedAt,
		UpdatedAt:  note.UpdatedAt,
	})
}

// ChangeStatus toggles note status.
func (s *NoteService) ChangeStatus(ctx echo.Context, id, status string) (*openapi.ModelsNoteResponse, error) {
	noteID, err := parseUUID(id)
	if err != nil {
		return nil, ErrInvalidNoteID
	}

	if status != "Draft" && status != "Publish" {
		return nil, errors.New("invalid status")
	}

	tx, err := s.pool.Begin(ctx.Request().Context())
	if err != nil {
		return nil, err
	}
	txQueries := s.queries.WithTx(tx)

	note, err := txQueries.UpdateNoteStatus(ctx.Request().Context(), &sqldb.UpdateNoteStatusParams{
		ID:     noteID,
		Status: status,
	})
	if err != nil {
		if rbErr := tx.Rollback(ctx.Request().Context()); rbErr != nil {
			return nil, rbErr
		}
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNoteNotFound
		}
		return nil, err
	}
	if err := tx.Commit(ctx.Request().Context()); err != nil {
		if rbErr := tx.Rollback(ctx.Request().Context()); rbErr != nil {
			return nil, rbErr
		}
		return nil, err
	}
	return s.composeNoteResponse(ctx.Request().Context(), noteResponseInput{
		ID:         note.ID,
		Title:      note.Title,
		TemplateID: note.TemplateID,
		OwnerID:    note.OwnerID,
		Status:     note.Status,
		CreatedAt:  note.CreatedAt,
		UpdatedAt:  note.UpdatedAt,
	})
}

// DeleteNote removes note and sections.
func (s *NoteService) DeleteNote(ctx echo.Context, id string) error {
	noteID, err := parseUUID(id)
	if err != nil {
		return ErrInvalidNoteID
	}

	if err := s.queries.DeleteNote(ctx.Request().Context(), noteID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNoteNotFound
		}
		return err
	}
	return nil
}

type noteResponseInput struct {
	ID           pgtype.UUID
	Title        string
	TemplateID   pgtype.UUID
	OwnerID      pgtype.UUID
	Status       string
	CreatedAt    pgtype.Timestamptz
	UpdatedAt    pgtype.Timestamptz
	TemplateName string
}

func (s *NoteService) composeNoteResponse(ctx context.Context, input noteResponseInput) (*openapi.ModelsNoteResponse, error) {
	templateName := input.TemplateName
	if strings.TrimSpace(templateName) == "" {
		template, err := s.queries.GetTemplateByID(ctx, input.TemplateID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, ErrInvalidTemplateID
			}
			return nil, err
		}
		templateName = template.Name
	}

	owner, err := s.queries.GetAccountByID(ctx, input.OwnerID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrAccountNotFound
		}
		return nil, err
	}

	sections, err := s.queries.ListSectionsByNote(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	sectionResponses := make([]openapi.ModelsSection, 0, len(sections))
	for _, section := range sections {
		sectionResponses = append(sectionResponses, openapi.ModelsSection{
			Id:         uuidToString(section.ID),
			FieldId:    uuidToString(section.FieldID),
			FieldLabel: section.Label,
			Content:    section.Content,
			IsRequired: section.IsRequired,
		})
	}

	status := openapi.ModelsNoteStatus(input.Status)
	if status != openapi.ModelsNoteStatusDraft && status != openapi.ModelsNoteStatusPublish {
		status = openapi.ModelsNoteStatusDraft
	}
	response := &openapi.ModelsNoteResponse{
		Id:           uuidToString(input.ID),
		Title:        input.Title,
		TemplateId:   uuidToString(input.TemplateID),
		TemplateName: templateName,
		OwnerId:      uuidToString(input.OwnerID),
		Owner: openapi.ModelsAccountSummary{
			Id:        uuidToString(owner.ID),
			FirstName: owner.FirstName,
			LastName:  owner.LastName,
			Thumbnail: textToPointer(owner.Thumbnail),
		},
		Status:    status,
		Sections:  sectionResponses,
		CreatedAt: input.CreatedAt.Time,
		UpdatedAt: input.UpdatedAt.Time,
	}

	return response, nil
}
