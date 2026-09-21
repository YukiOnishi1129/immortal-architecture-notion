package sqlc

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"immortal-architecture-notion/backend/internal/adapter/gateway/db/sqlc/generated"
	domainerr "immortal-architecture-notion/backend/internal/domain/errors"
	"immortal-architecture-notion/backend/internal/domain/note"
	"immortal-architecture-notion/backend/internal/port"
)

// NoteReadModelRepository implements query (read) persistence using the read model table.
type NoteReadModelRepository struct {
	pool    *pgxpool.Pool
	queries *generated.Queries
}

var _ port.NoteReadModelRepository = (*NoteReadModelRepository)(nil)

// NewNoteReadModelRepository creates NoteReadModelRepository.
func NewNoteReadModelRepository(pool *pgxpool.Pool) *NoteReadModelRepository {
	return &NoteReadModelRepository{
		pool:    pool,
		queries: generated.New(pool),
	}
}

// List returns notes from the read model table without JOINs.
func (r *NoteReadModelRepository) List(ctx context.Context, filters note.Filters) ([]note.ReadModel, error) {
	params := &generated.ListNoteReadModelsParams{}
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

	rows, err := queriesForContext(ctx, r.queries).ListNoteReadModels(ctx, params)
	if err != nil {
		return nil, err
	}

	result := make([]note.ReadModel, 0, len(rows))
	for _, row := range rows {
		rm, err := toNoteReadModel(row)
		if err != nil {
			return nil, err
		}
		result = append(result, rm)
	}
	return result, nil
}

// Get returns a single note from the read model table.
func (r *NoteReadModelRepository) Get(ctx context.Context, id string) (*note.ReadModel, error) {
	pgID, err := toUUID(id)
	if err != nil {
		return nil, err
	}
	row, err := queriesForContext(ctx, r.queries).GetNoteReadModel(ctx, pgID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainerr.ErrNotFound
		}
		return nil, err
	}
	rm, err := toNoteReadModel(row)
	if err != nil {
		return nil, err
	}
	return &rm, nil
}

// Upsert inserts or updates a read model entry (called from the command side).
func (r *NoteReadModelRepository) Upsert(ctx context.Context, model note.ReadModel) error {
	id, err := toUUID(model.ID)
	if err != nil {
		return err
	}
	templateID, err := toUUID(model.TemplateID)
	if err != nil {
		return err
	}
	ownerID, err := toUUID(model.OwnerID)
	if err != nil {
		return err
	}

	sectionsJSON, err := json.Marshal(model.Sections)
	if err != nil {
		return err
	}

	return queriesForContext(ctx, r.queries).UpsertNoteReadModel(ctx, &generated.UpsertNoteReadModelParams{
		ID:             id,
		Title:          model.Title,
		Status:         string(model.Status),
		TemplateID:     templateID,
		TemplateName:   model.TemplateName,
		OwnerID:        ownerID,
		OwnerFirstName: model.OwnerFirstName,
		OwnerLastName:  model.OwnerLastName,
		OwnerThumbnail: pgNullableText(model.OwnerThumbnail),
		SectionsJson:   sectionsJSON,
		CreatedAt:      timeToPgTimestamptz(model.CreatedAt),
		UpdatedAt:      timeToPgTimestamptz(model.UpdatedAt),
	})
}

// Delete removes a read model entry (called from the command side).
func (r *NoteReadModelRepository) Delete(ctx context.Context, id string) error {
	pgID, err := toUUID(id)
	if err != nil {
		return err
	}
	return queriesForContext(ctx, r.queries).DeleteNoteReadModel(ctx, pgID)
}

func toNoteReadModel(row *generated.NoteReadModel) (note.ReadModel, error) {
	var sections []note.SectionReadModel
	if len(row.SectionsJson) > 0 {
		if err := json.Unmarshal(row.SectionsJson, &sections); err != nil {
			return note.ReadModel{}, err
		}
	}

	var thumbnail *string
	if row.OwnerThumbnail.Valid {
		s := row.OwnerThumbnail.String
		thumbnail = &s
	}

	return note.ReadModel{
		ID:             uuidToString(row.ID),
		Title:          row.Title,
		Status:         note.NoteStatus(row.Status),
		TemplateID:     uuidToString(row.TemplateID),
		TemplateName:   row.TemplateName,
		OwnerID:        uuidToString(row.OwnerID),
		OwnerFirstName: row.OwnerFirstName,
		OwnerLastName:  row.OwnerLastName,
		OwnerThumbnail: thumbnail,
		Sections:       sections,
		CreatedAt:      timestamptzToTime(row.CreatedAt),
		UpdatedAt:      timestamptzToTime(row.UpdatedAt),
	}, nil
}
