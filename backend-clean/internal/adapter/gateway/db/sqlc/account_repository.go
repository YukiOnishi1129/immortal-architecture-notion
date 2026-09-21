// Package sqlc implements gateway repositories using sqlc.
package sqlc

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"immortal-architecture-cqrs/backend/internal/adapter/gateway/db/sqlc/generated"
	"immortal-architecture-cqrs/backend/internal/domain/account"
	domainerr "immortal-architecture-cqrs/backend/internal/domain/errors"
	"immortal-architecture-cqrs/backend/internal/port"
)

// AccountRepository implements account persistence.
type AccountRepository struct {
	pool    *pgxpool.Pool
	queries *generated.Queries
}

var _ port.AccountRepository = (*AccountRepository)(nil)

// NewAccountRepository creates AccountRepository.
func NewAccountRepository(pool *pgxpool.Pool) *AccountRepository {
	return &AccountRepository{
		pool:    pool,
		queries: generated.New(pool),
	}
}

// UpsertOAuthAccount inserts or updates an OAuth account.
func (r *AccountRepository) UpsertOAuthAccount(ctx context.Context, input account.OAuthAccountInput) (*account.Account, error) {
	q := queriesForContext(ctx, r.queries)

	row, err := q.UpsertAccount(ctx, &generated.UpsertAccountParams{
		Email:             input.Email,
		FirstName:         input.FirstName,
		LastName:          input.LastName,
		Provider:          input.Provider,
		ProviderAccountID: input.ProviderAccountID,
		Thumbnail:         pgNullableText(input.Thumbnail),
		LastLoginAt:       pgNullableTime(nil),
	})
	if err != nil {
		return nil, err
	}
	return toDomainAccount(row)
}

// GetByID fetches account by ID.
func (r *AccountRepository) GetByID(ctx context.Context, id string) (*account.Account, error) {
	q := queriesForContext(ctx, r.queries)
	uuid, err := toUUID(id)
	if err != nil {
		return nil, err
	}
	row, err := q.GetAccountByID(ctx, uuid)
	if err != nil {
		return nil, err
	}
	return toDomainAccount(row)
}

// GetByEmail fetches account by email.
func (r *AccountRepository) GetByEmail(ctx context.Context, email string) (*account.Account, error) {
	q := queriesForContext(ctx, r.queries)
	row, err := q.GetAccountByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainerr.ErrNotFound
		}
		return nil, err
	}
	return toDomainAccount(row)
}

func toDomainAccount(a *generated.Account) (*account.Account, error) {
	var lastLogin *time.Time
	if a.LastLoginAt.Valid {
		t := timestamptzToTime(a.LastLoginAt)
		lastLogin = &t
	}
	email, err := account.ParseEmail(a.Email)
	if err != nil {
		return nil, err
	}
	return &account.Account{
		ID:                uuidToString(a.ID),
		Email:             email,
		FirstName:         a.FirstName,
		LastName:          a.LastName,
		IsActive:          a.IsActive,
		Provider:          a.Provider,
		ProviderAccountID: a.ProviderAccountID,
		Thumbnail:         nullableTextToString(a.Thumbnail),
		LastLoginAt:       lastLogin,
		CreatedAt:         timestamptzToTime(a.CreatedAt),
		UpdatedAt:         timestamptzToTime(a.UpdatedAt),
	}, nil
}
