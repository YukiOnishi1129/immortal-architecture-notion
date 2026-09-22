// Package sqlc implements gateway repositories using sqlc.
package sqlc

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"immortal-architecture-notion/backend/internal/adapter/gateway/db/sqlc/generated"
	domainerr "immortal-architecture-notion/backend/internal/domain/errors"
	"immortal-architecture-notion/backend/internal/domain/template"
	"immortal-architecture-notion/backend/internal/port"
)

func toTemplateOwner(ownerID pgtype.UUID, first, last string, thumb pgtype.Text) template.Owner {
	var thumbnail *string
	if thumb.Valid {
		s := thumb.String
		thumbnail = &s
	}
	return template.Owner{
		ID:        uuidToString(ownerID),
		FirstName: first,
		LastName:  last,
		Thumbnail: thumbnail,
	}
}

// TemplateRepository implements template persistence.
type TemplateRepository struct {
	pool    *pgxpool.Pool
	queries *generated.Queries
}

var _ port.TemplateRepository = (*TemplateRepository)(nil)

// NewTemplateRepository creates TemplateRepository.
func NewTemplateRepository(pool *pgxpool.Pool) *TemplateRepository {
	return &TemplateRepository{
		pool:    pool,
		queries: generated.New(pool),
	}
}

// List returns templates by filters.
func (r *TemplateRepository) List(ctx context.Context, filters template.Filters) ([]template.WithUsage, error) {
	params := &generated.ListTemplatesParams{}
	if filters.OwnerID != nil && *filters.OwnerID != "" {
		if id, err := toUUID(*filters.OwnerID); err == nil {
			params.Column1 = id
		}
	}
	if filters.Query != nil && *filters.Query != "" {
		params.Column2 = *filters.Query
	}

	rows, err := queriesForContext(ctx, r.queries).ListTemplates(ctx, params)
	if err != nil {
		return nil, err
	}

	result := make([]template.WithUsage, 0, len(rows))
	for _, row := range rows {
		fields, err := r.listFields(ctx, row.ID)
		if err != nil {
			return nil, err
		}
		owner := toTemplateOwner(row.OwnerID, row.OwnerFirstName, row.OwnerLastName, row.OwnerThumbnail)
		result = append(result, template.WithUsage{
			Template: template.Template{
				ID:                 uuidToString(row.ID),
				Name:               row.Name,
				OwnerID:            uuidToString(row.OwnerID),
				UpdatedAt:          timestamptzToTime(row.UpdatedAt),
				Fields:             fields,
				NotionParentPageID: nullableTextToString(row.NotionParentPageID),
			},
			IsUsed: row.IsUsed,
			Owner:  owner,
		})
	}
	return result, nil
}

// Get returns a template with usage and fields.
func (r *TemplateRepository) Get(ctx context.Context, id string) (*template.WithUsage, error) {
	pgID, err := toUUID(id)
	if err != nil {
		return nil, err
	}
	row, err := queriesForContext(ctx, r.queries).GetTemplateByID(ctx, pgID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainerr.ErrNotFound
		}
		return nil, err
	}
	fields, err := r.listFields(ctx, row.ID)
	if err != nil {
		return nil, err
	}
	owner := toTemplateOwner(row.OwnerID, row.OwnerFirstName, row.OwnerLastName, row.OwnerThumbnail)
	return &template.WithUsage{
		Template: template.Template{
			ID:                 uuidToString(row.ID),
			Name:               row.Name,
			OwnerID:            uuidToString(row.OwnerID),
			UpdatedAt:          timestamptzToTime(row.UpdatedAt),
			Fields:             fields,
			NotionParentPageID: nullableTextToString(row.NotionParentPageID),
		},
		IsUsed: row.IsUsed,
		Owner:  owner,
	}, nil
}

// Create inserts a template.
func (r *TemplateRepository) Create(ctx context.Context, tpl template.Template) (*template.Template, error) {
	owner, err := toUUID(tpl.OwnerID)
	if err != nil {
		return nil, err
	}
	row, err := queriesForContext(ctx, r.queries).CreateTemplate(ctx, &generated.CreateTemplateParams{
		Name:               tpl.Name,
		OwnerID:            owner,
		NotionParentPageID: pgTextFromString(tpl.NotionParentPageID),
	})
	if err != nil {
		return nil, err
	}
	return &template.Template{
		ID:                 uuidToString(row.ID),
		Name:               row.Name,
		OwnerID:            uuidToString(row.OwnerID),
		UpdatedAt:          timestamptzToTime(row.UpdatedAt),
		NotionParentPageID: nullableTextToString(row.NotionParentPageID),
	}, nil
}

// Update updates template name.
func (r *TemplateRepository) Update(ctx context.Context, tpl template.Template) (*template.Template, error) {
	pgID, err := toUUID(tpl.ID)
	if err != nil {
		return nil, err
	}
	row, err := queriesForContext(ctx, r.queries).UpdateTemplate(ctx, &generated.UpdateTemplateParams{
		ID:                 pgID,
		Name:               tpl.Name,
		NotionParentPageID: pgTextFromString(tpl.NotionParentPageID),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainerr.ErrNotFound
		}
		return nil, err
	}
	return &template.Template{
		ID:                 uuidToString(row.ID),
		Name:               row.Name,
		OwnerID:            uuidToString(row.OwnerID),
		UpdatedAt:          timestamptzToTime(row.UpdatedAt),
		NotionParentPageID: nullableTextToString(row.NotionParentPageID),
	}, nil
}

// Delete deletes a template.
func (r *TemplateRepository) Delete(ctx context.Context, id string) error {
	pgID, err := toUUID(id)
	if err != nil {
		return err
	}
	return queriesForContext(ctx, r.queries).DeleteTemplate(ctx, pgID)
}

// ReplaceFields replaces template fields.
// ReplaceFields makes the stored fields match the given ones.
//
// Fields are updated in place rather than deleted and re-created, because a
// field that a note already has a section for cannot be deleted: the section
// references it. Re-creating every field would therefore fail as soon as the
// template is in use, which also blocks edits that leave the fields alone,
// such as renaming the template or setting its Notion parent page.
func (r *TemplateRepository) ReplaceFields(ctx context.Context, templateID string, fields []template.Field) error {
	pgID, err := toUUID(templateID)
	if err != nil {
		return err
	}
	q := queriesForContext(ctx, r.queries)

	existing, err := q.ListFieldsByTemplate(ctx, pgID)
	if err != nil {
		return err
	}
	existingByID := make(map[string]*generated.Field, len(existing))
	for _, row := range existing {
		existingByID[uuidToString(row.ID)] = row
	}

	kept := make(map[string]bool, len(fields))
	for idx, f := range fields {
		order := f.Order
		if order == 0 {
			order = idx + 1
		}

		// A field carrying a known id is an existing one being edited.
		if f.ID != "" && existingByID[f.ID] != nil {
			fieldID, err := toUUID(f.ID)
			if err != nil {
				return err
			}
			if _, err := q.UpdateField(ctx, &generated.UpdateFieldParams{
				ID:         fieldID,
				Label:      f.Label,
				Order:      int32(order), //nolint:gosec
				IsRequired: f.IsRequired,
			}); err != nil {
				return err
			}
			kept[f.ID] = true
			continue
		}

		if _, err := q.CreateField(ctx, &generated.CreateFieldParams{
			TemplateID: pgID,
			Label:      f.Label,
			Order:      int32(order), //nolint:gosec
			IsRequired: f.IsRequired,
		}); err != nil {
			return err
		}
	}

	// Whatever the caller left out is removed. This still fails when a note
	// uses the field, which is the intended protection.
	for id, row := range existingByID {
		if kept[id] {
			continue
		}
		if err := q.DeleteField(ctx, row.ID); err != nil {
			return err
		}
	}
	return nil
}

func (r *TemplateRepository) listFields(ctx context.Context, templateID pgtype.UUID) ([]template.Field, error) {
	rows, err := queriesForContext(ctx, r.queries).ListFieldsByTemplate(ctx, templateID)
	if err != nil {
		return nil, err
	}
	fields := make([]template.Field, 0, len(rows))
	for _, f := range rows {
		fields = append(fields, template.Field{
			ID:         uuidToString(f.ID),
			Label:      f.Label,
			Order:      int(f.Order),
			IsRequired: f.IsRequired,
		})
	}
	return fields, nil
}

// UsedFieldIDs returns the fields of a template that notes reference.
func (r *TemplateRepository) UsedFieldIDs(ctx context.Context, templateID string) ([]string, error) {
	pgID, err := toUUID(templateID)
	if err != nil {
		return nil, err
	}
	rows, err := queriesForContext(ctx, r.queries).ListUsedFieldIDsByTemplate(ctx, pgID)
	if err != nil {
		return nil, err
	}
	ids := make([]string, 0, len(rows))
	for _, id := range rows {
		ids = append(ids, uuidToString(id))
	}
	return ids, nil
}
