package service

import (
	"context"
	"errors"
	"math"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"

	sqldb "immortal-architecture-bad-api/backend/internal/db/sqlc"
	openapi "immortal-architecture-bad-api/backend/internal/generated/openapi"
)

var (
	// ErrTemplateNotFound indicates the template does not exist.
	ErrTemplateNotFound = errors.New("template not found")
	// ErrTemplateInUse indicates the template is referenced by notes.
	ErrTemplateInUse = errors.New("template is in use")
	// ErrInvalidTemplateID invalid UUID.
	ErrInvalidTemplateID = errors.New("invalid template id")
)

// TemplateService handles template CRUD.
type TemplateService struct {
	pool    *pgxpool.Pool
	queries *sqldb.Queries
}

// NewTemplateService creates a new service.
func NewTemplateService(pool *pgxpool.Pool) *TemplateService {
	return &TemplateService{
		pool:    pool,
		queries: sqldb.New(pool),
	}
}

// TemplateFilters filters template list query.
type TemplateFilters struct {
	OwnerID *string
	Query   *string
}

// ListTemplates returns templates optionally filtered by owner or search query.
func (s *TemplateService) ListTemplates(ctx echo.Context, filters TemplateFilters) ([]*openapi.ModelsTemplateResponse, error) {
	dbCtx := ctx.Request().Context()
	params := &sqldb.ListTemplatesParams{}

	if filters.OwnerID != nil && *filters.OwnerID != "" {
		if id, err := parseUUID(*filters.OwnerID); err == nil {
			params.Column1 = id
		}
	}

	if filters.Query != nil && *filters.Query != "" {
		params.Column2 = *filters.Query
	}

	rows, err := s.queries.ListTemplates(dbCtx, params)
	if err != nil {
		return nil, err
	}

	responses := make([]*openapi.ModelsTemplateResponse, 0, len(rows))
	for _, tpl := range rows {
		response, err := s.composeTemplateResponse(dbCtx, tpl.ID, tpl.Name, tpl.OwnerID, tpl.UpdatedAt, tpl.IsUsed)
		if err != nil {
			return nil, err
		}
		responses = append(responses, response)
	}

	return responses, nil
}

// GetTemplate returns a template by ID along with is_used flag.
func (s *TemplateService) GetTemplate(ctx echo.Context, id string) (*openapi.ModelsTemplateResponse, error) {
	pgID, err := parseUUID(id)
	if err != nil {
		return nil, ErrInvalidTemplateID
	}

	template, err := s.queries.GetTemplateByID(ctx.Request().Context(), pgID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrTemplateNotFound
		}
		return nil, err
	}
	return s.composeTemplateResponse(ctx.Request().Context(), template.ID, template.Name, template.OwnerID, template.UpdatedAt, template.IsUsed)
}

// CreateTemplate creates a template with its fields.
func (s *TemplateService) CreateTemplate(ctx echo.Context, ownerID string, name string, fields []sqldb.Field) (*openapi.ModelsTemplateResponse, error) {
	pgOwner, err := parseUUID(ownerID)
	if err != nil {
		return nil, ErrInvalidAccountID
	}

	tx, err := s.pool.Begin(ctx.Request().Context())
	if err != nil {
		return nil, err
	}
	txQueries := s.queries.WithTx(tx)

	template, err := txQueries.CreateTemplate(ctx.Request().Context(), &sqldb.CreateTemplateParams{
		Name:    name,
		OwnerID: pgOwner,
	})
	if err != nil {
		if rbErr := tx.Rollback(ctx.Request().Context()); rbErr != nil {
			return nil, rbErr
		}
		return nil, err
	}

	for idx, field := range fields {
		if idx > math.MaxInt32-1 {
			return nil, errors.New("too many fields")
		}
		order := field.Order
		if order == 0 {
			if idx >= math.MaxInt32 {
				return nil, errors.New("too many fields")
			}
			order = int32(idx + 1) //nolint:gosec // idx is bounded above
		}
		_, err = txQueries.CreateField(ctx.Request().Context(), &sqldb.CreateFieldParams{
			TemplateID: template.ID,
			Label:      field.Label,
			Order:      order,
			IsRequired: field.IsRequired,
		})
		if err != nil {
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

	return s.composeTemplateResponse(ctx.Request().Context(), template.ID, template.Name, template.OwnerID, template.UpdatedAt, false)
}

// UpdateTemplate updates template metadata and reorders fields.
func (s *TemplateService) UpdateTemplate(ctx echo.Context, id, name string, fields []openapi.ModelsUpdateFieldRequest) (*openapi.ModelsTemplateResponse, error) {
	templateID, err := parseUUID(id)
	if err != nil {
		return nil, ErrInvalidTemplateID
	}

	tx, err := s.pool.Begin(ctx.Request().Context())
	if err != nil {
		return nil, err
	}
	txQueries := s.queries.WithTx(tx)

	template, err := txQueries.UpdateTemplate(ctx.Request().Context(), &sqldb.UpdateTemplateParams{
		ID:   templateID,
		Name: name,
	})
	if err != nil {
		if rbErr := tx.Rollback(ctx.Request().Context()); rbErr != nil {
			return nil, rbErr
		}
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrTemplateNotFound
		}
		return nil, err
	}

	if fields != nil {
		if len(fields) == 0 {
			return nil, errors.New("at least one field is required")
		}
		if err := s.syncTemplateFields(ctx.Request().Context(), txQueries, templateID, fields); err != nil {
			if rbErr := tx.Rollback(ctx.Request().Context()); rbErr != nil {
				return nil, rbErr
			}
			return nil, err
		}
	}

	isUsed, err := txQueries.CheckTemplateInUse(ctx.Request().Context(), templateID)
	if err != nil {
		if rbErr := tx.Rollback(ctx.Request().Context()); rbErr != nil {
			return nil, rbErr
		}
		return nil, err
	}

	if err := tx.Commit(ctx.Request().Context()); err != nil {
		if rbErr := tx.Rollback(ctx.Request().Context()); rbErr != nil {
			return nil, rbErr
		}
		return nil, err
	}

	return s.composeTemplateResponse(ctx.Request().Context(), template.ID, template.Name, template.OwnerID, template.UpdatedAt, isUsed)
}

// DeleteTemplate removes template if not in use.
func (s *TemplateService) DeleteTemplate(ctx echo.Context, id string) error {
	templateID, err := parseUUID(id)
	if err != nil {
		return ErrInvalidTemplateID
	}

	tx, err := s.pool.Begin(ctx.Request().Context())
	if err != nil {
		return err
	}
	txQueries := s.queries.WithTx(tx)

	inUse, err := txQueries.CheckTemplateInUse(ctx.Request().Context(), templateID)
	if err != nil {
		if rbErr := tx.Rollback(ctx.Request().Context()); rbErr != nil {
			return rbErr
		}
		return err
	}
	if inUse {
		if rbErr := tx.Rollback(ctx.Request().Context()); rbErr != nil {
			return rbErr
		}
		return ErrTemplateInUse
	}

	if err := txQueries.DeleteTemplate(ctx.Request().Context(), templateID); err != nil {
		if rbErr := tx.Rollback(ctx.Request().Context()); rbErr != nil {
			return err
		}
		return err
	}
	if err := tx.Commit(ctx.Request().Context()); err != nil {
		if rbErr := tx.Rollback(ctx.Request().Context()); rbErr != nil {
			return err
		}
		return err
	}
	return nil
}

func (s *TemplateService) syncTemplateFields(ctx context.Context, queries *sqldb.Queries, templateID pgtype.UUID, fields []openapi.ModelsUpdateFieldRequest) error {
	existing, err := queries.ListFieldsByTemplate(ctx, templateID)
	if err != nil {
		return err
	}

	needsTrim := len(fields) < len(existing)
	if needsTrim {
		inUse, err := queries.CheckTemplateInUse(ctx, templateID)
		if err != nil {
			return err
		}
		if inUse {
			return ErrTemplateInUse
		}
	}

	for idx, field := range fields {
		if idx >= math.MaxInt32 {
			return errors.New("too many fields")
		}

		order64 := idx + 1
		if order64 >= math.MaxInt32 {
			return errors.New("field order overflow")
		}
		order := int32(order64) //nolint:gosec // bounded above by math.MaxInt32
		if idx < len(existing) {
			if _, err := s.queries.UpdateField(ctx, &sqldb.UpdateFieldParams{
				ID:         existing[idx].ID,
				Label:      field.Label,
				Order:      order,
				IsRequired: field.IsRequired,
			}); err != nil {
				return err
			}
			continue
		}

		if _, err := s.queries.CreateField(ctx, &sqldb.CreateFieldParams{
			TemplateID: templateID,
			Label:      field.Label,
			Order:      order,
			IsRequired: field.IsRequired,
		}); err != nil {
			return err
		}
	}

	if needsTrim {
		for _, obsolete := range existing[len(fields):] {
			if err := queries.DeleteField(ctx, obsolete.ID); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *TemplateService) composeTemplateResponse(ctx context.Context, templateID pgtype.UUID, name string, ownerID pgtype.UUID, updatedAt pgtype.Timestamptz, isUsed bool) (*openapi.ModelsTemplateResponse, error) {
	owner, err := s.queries.GetAccountByID(ctx, ownerID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrAccountNotFound
		}
		return nil, err
	}

	fieldRows, err := s.queries.ListFieldsByTemplate(ctx, templateID)
	if err != nil {
		return nil, err
	}

	fields := make([]openapi.ModelsField, 0, len(fieldRows))
	for _, field := range fieldRows {
		fields = append(fields, openapi.ModelsField{
			Id:         uuidToString(field.ID),
			Label:      field.Label,
			Order:      field.Order,
			IsRequired: field.IsRequired,
		})
	}

	response := &openapi.ModelsTemplateResponse{
		Id:      uuidToString(templateID),
		Name:    name,
		OwnerId: uuidToString(ownerID),
		Owner: openapi.ModelsAccountSummary{
			Id:        uuidToString(owner.ID),
			FirstName: owner.FirstName,
			LastName:  owner.LastName,
			Thumbnail: textToPointer(owner.Thumbnail),
		},
		Fields:    fields,
		UpdatedAt: updatedAt.Time,
		IsUsed:    isUsed,
	}

	return response, nil
}
