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
// ReplaceFields deletes every field and writes the given ones back.
//
// Notes reference fields by id, so this is only valid while no note uses them.
// The caller decides; see SyncFields for the in-place variant.
func (r *TemplateRepository) ReplaceFields(ctx context.Context, templateID string, fields []template.Field) error {
	pgID, err := toUUID(templateID)
	if err != nil {
		return err
	}
	q := queriesForContext(ctx, r.queries)

	if err := q.DeleteFieldsByTemplate(ctx, pgID); err != nil {
		return err
	}
	for idx, f := range fields {
		if _, err := q.CreateField(ctx, &generated.CreateFieldParams{
			TemplateID: pgID,
			Label:      f.Label,
			Order:      int32(orderOrIndex(f, idx)), //nolint:gosec
			IsRequired: f.IsRequired,
		}); err != nil {
			return err
		}
	}
	return nil
}

// SyncFields makes the stored fields match the given ones without recreating
// them, so the ids that notes point at survive.
//
// Reordering happens in two passes: (template_id, "order") is unique, so
// moving a field onto a position another one still holds would clash
// mid-update. Every updated field is first parked above any position in use,
// then moved into place. The column also has CHECK (order > 0), so the
// parking positions stay positive.
func (r *TemplateRepository) SyncFields(ctx context.Context, templateID string, fields []template.Field) error {
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
	updates := make([]generated.UpdateFieldParams, 0, len(fields))
	creates := make([]generated.CreateFieldParams, 0, len(fields))

	for idx, f := range fields {
		order := orderOrIndex(f, idx)

		if f.ID != "" && existingByID[f.ID] != nil {
			fieldID, err := toUUID(f.ID)
			if err != nil {
				return err
			}
			updates = append(updates, generated.UpdateFieldParams{
				ID:         fieldID,
				Label:      f.Label,
				Order:      int32(order), //nolint:gosec
				IsRequired: f.IsRequired,
			})
			kept[f.ID] = true
			continue
		}

		creates = append(creates, generated.CreateFieldParams{
			TemplateID: pgID,
			Label:      f.Label,
			Order:      int32(order), //nolint:gosec
			IsRequired: f.IsRequired,
		})
	}

	parkFrom := int32(len(existing) + len(fields) + 1) //nolint:gosec
	for i, u := range updates {
		parked := u
		parked.Order = parkFrom + int32(i) //nolint:gosec
		if _, err := q.UpdateField(ctx, &parked); err != nil {
			return err
		}
	}

	// New rows go in while the existing ones are parked, so their positions
	// are free.
	for i := range creates {
		if _, err := q.CreateField(ctx, &creates[i]); err != nil {
			return err
		}
	}
	for i := range updates {
		if _, err := q.UpdateField(ctx, &updates[i]); err != nil {
			return err
		}
	}

	// Anything the caller left out is removed. This still fails when a note
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

// orderOrIndex falls back to the position in the slice when no order is set.
func orderOrIndex(f template.Field, idx int) int {
	if f.Order == 0 {
		return idx + 1
	}
	return f.Order
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
