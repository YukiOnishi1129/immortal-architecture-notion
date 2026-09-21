// Package sqlc implements gateway repositories using sqlc.
package sqlc

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"immortal-architecture-cqrs/backend/internal/adapter/gateway/db/sqlc/generated"
	domainerr "immortal-architecture-cqrs/backend/internal/domain/errors"
	"immortal-architecture-cqrs/backend/internal/domain/note"
	"immortal-architecture-cqrs/backend/internal/port"
)

// NoteRepository implements note persistence.
type NoteRepository struct {
	pool    *pgxpool.Pool
	queries *generated.Queries
}

var _ port.NoteRepository = (*NoteRepository)(nil)

// NewNoteRepository creates NoteRepository.
func NewNoteRepository(pool *pgxpool.Pool) *NoteRepository {
	return &NoteRepository{
		pool:    pool,
		queries: generated.New(pool),
	}
}

// List returns notes by filters.
func (r *NoteRepository) List(ctx context.Context, filters note.Filters) ([]note.WithMeta, error) {
	params := &generated.ListNotesParams{}
	if filters.Status != nil {
		params.Column1 = string(*filters.Status)
	}
	if filters.TemplateID != nil && *filters.TemplateID != "" {
		if id, err := toUUID(*filters.TemplateID); err == nil {
			params.Column2 = id
		}
	}
	if filters.OwnerID != nil && *filters.OwnerID != "" {
		if id, err := toUUID(*filters.OwnerID); err == nil {
			params.Column3 = id
		}
	}
	if filters.Query != nil && *filters.Query != "" {
		params.Column4 = *filters.Query
	}

	rows, err := queriesForContext(ctx, r.queries).ListNotes(ctx, params)
	if err != nil {
		return nil, err
	}

	result := make([]note.WithMeta, 0, len(rows))
	for _, row := range rows {
		sections, err := r.listSections(ctx, row.ID)
		if err != nil {
			return nil, err
		}
		var thumbnail *string
		if row.OwnerThumbnail.Valid {
			s := row.OwnerThumbnail.String
			thumbnail = &s
		}
		result = append(result, note.WithMeta{
			Note: note.Note{
				ID:         uuidToString(row.ID),
				Title:      row.Title,
				TemplateID: uuidToString(row.TemplateID),
				OwnerID:    uuidToString(row.OwnerID),
				Status:     note.NoteStatus(row.Status),
				CreatedAt:  timestamptzToTime(row.CreatedAt),
				UpdatedAt:  timestamptzToTime(row.UpdatedAt),
			},
			TemplateName:   row.TemplateName,
			OwnerFirstName: row.FirstName,
			OwnerLastName:  row.LastName,
			OwnerThumbnail: thumbnail,
			Sections:       sections,
		})
	}
	return result, nil
}

// Get returns a note with sections.
func (r *NoteRepository) Get(ctx context.Context, id string) (*note.WithMeta, error) {
	pgID, err := toUUID(id)
	if err != nil {
		return nil, err
	}
	row, err := queriesForContext(ctx, r.queries).GetNoteByID(ctx, pgID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainerr.ErrNotFound
		}
		return nil, err
	}
	sections, err := r.listSections(ctx, row.ID)
	if err != nil {
		return nil, err
	}
	var thumbnail *string
	if row.OwnerThumbnail.Valid {
		s := row.OwnerThumbnail.String
		thumbnail = &s
	}
	return &note.WithMeta{
		Note: note.Note{
			ID:         uuidToString(row.ID),
			Title:      row.Title,
			TemplateID: uuidToString(row.TemplateID),
			OwnerID:    uuidToString(row.OwnerID),
			Status:     note.NoteStatus(row.Status),
			CreatedAt:  timestamptzToTime(row.CreatedAt),
			UpdatedAt:  timestamptzToTime(row.UpdatedAt),
		},
		TemplateName:   row.TemplateName,
		OwnerFirstName: row.FirstName,
		OwnerLastName:  row.LastName,
		OwnerThumbnail: thumbnail,
		Sections:       sections,
	}, nil
}

// Create inserts a note.
func (r *NoteRepository) Create(ctx context.Context, n note.Note) (*note.Note, error) {
	templateID, err := toUUID(n.TemplateID)
	if err != nil {
		return nil, err
	}
	ownerID, err := toUUID(n.OwnerID)
	if err != nil {
		return nil, err
	}
	row, err := queriesForContext(ctx, r.queries).CreateNote(ctx, &generated.CreateNoteParams{
		Title:      n.Title,
		TemplateID: templateID,
		OwnerID:    ownerID,
		Status:     string(n.Status),
	})
	if err != nil {
		return nil, err
	}
	return &note.Note{
		ID:         uuidToString(row.ID),
		Title:      row.Title,
		TemplateID: uuidToString(row.TemplateID),
		OwnerID:    uuidToString(row.OwnerID),
		Status:     note.NoteStatus(row.Status),
		CreatedAt:  timestamptzToTime(row.CreatedAt),
		UpdatedAt:  timestamptzToTime(row.UpdatedAt),
	}, nil
}

// Update updates a note title.
func (r *NoteRepository) Update(ctx context.Context, n note.Note) (*note.Note, error) {
	pgID, err := toUUID(n.ID)
	if err != nil {
		return nil, err
	}
	row, err := queriesForContext(ctx, r.queries).UpdateNote(ctx, &generated.UpdateNoteParams{
		ID:    pgID,
		Title: n.Title,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainerr.ErrNotFound
		}
		return nil, err
	}
	return &note.Note{
		ID:         uuidToString(row.ID),
		Title:      row.Title,
		TemplateID: uuidToString(row.TemplateID),
		OwnerID:    uuidToString(row.OwnerID),
		Status:     note.NoteStatus(row.Status),
		CreatedAt:  timestamptzToTime(row.CreatedAt),
		UpdatedAt:  timestamptzToTime(row.UpdatedAt),
	}, nil
}

// UpdateStatus updates note status.
func (r *NoteRepository) UpdateStatus(ctx context.Context, id string, status note.NoteStatus) (*note.Note, error) {
	pgID, err := toUUID(id)
	if err != nil {
		return nil, err
	}
	row, err := queriesForContext(ctx, r.queries).UpdateNoteStatus(ctx, &generated.UpdateNoteStatusParams{
		ID:     pgID,
		Status: string(status),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainerr.ErrNotFound
		}
		return nil, err
	}
	return &note.Note{
		ID:         uuidToString(row.ID),
		Title:      row.Title,
		TemplateID: uuidToString(row.TemplateID),
		OwnerID:    uuidToString(row.OwnerID),
		Status:     note.NoteStatus(row.Status),
		CreatedAt:  timestamptzToTime(row.CreatedAt),
		UpdatedAt:  timestamptzToTime(row.UpdatedAt),
	}, nil
}

// Delete deletes a note.
func (r *NoteRepository) Delete(ctx context.Context, id string) error {
	pgID, err := toUUID(id)
	if err != nil {
		return err
	}
	return queriesForContext(ctx, r.queries).DeleteNote(ctx, pgID)
}

// ReplaceSections replaces note sections.
func (r *NoteRepository) ReplaceSections(ctx context.Context, noteID string, sections []note.Section) error {
	nID, err := toUUID(noteID)
	if err != nil {
		return err
	}
	q := queriesForContext(ctx, r.queries)

	for _, s := range sections {
		if s.ID != "" {
			secID, err := toUUID(s.ID)
			if err != nil {
				return err
			}
			if _, err := q.UpdateSectionContent(ctx, &generated.UpdateSectionContentParams{
				ID:      secID,
				Content: s.Content,
			}); err != nil {
				return err
			}
			continue
		}
		fieldID, err := toUUID(s.FieldID)
		if err != nil {
			return err
		}
		if _, err := q.CreateSection(ctx, &generated.CreateSectionParams{
			NoteID:  nID,
			FieldID: fieldID,
			Content: s.Content,
		}); err != nil {
			return err
		}
	}
	return nil
}

func (r *NoteRepository) listSections(ctx context.Context, noteID pgtype.UUID) ([]note.SectionWithField, error) {
	rows, err := queriesForContext(ctx, r.queries).ListSectionsByNote(ctx, noteID)
	if err != nil {
		return nil, err
	}
	sections := make([]note.SectionWithField, 0, len(rows))
	for _, row := range rows {
		sections = append(sections, note.SectionWithField{
			Section: note.Section{
				ID:      uuidToString(row.ID),
				NoteID:  uuidToString(row.NoteID),
				FieldID: uuidToString(row.FieldID),
				Content: row.Content,
			},
			FieldLabel: row.Label,
			FieldOrder: int(row.Order),
			IsRequired: row.IsRequired,
		})
	}
	return sections, nil
}
