// Package service bundles business logic for the API.
package service

import (
	"errors"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"

	sqldb "immortal-architecture-bad-api/backend/internal/db/sqlc"
)

var (
	// ErrAccountNotFound indicates the account could not be located.
	ErrAccountNotFound = errors.New("account not found")
	// ErrInvalidAccountID indicates the provided ID is not a UUID.
	ErrInvalidAccountID = errors.New("invalid account id")
)

// AccountService bundles all account-related operations.
type AccountService struct {
	pool    *pgxpool.Pool
	queries *sqldb.Queries
}

// NewAccountService creates a new service instance.
func NewAccountService(pool *pgxpool.Pool) *AccountService {
	return &AccountService{
		pool:    pool,
		queries: sqldb.New(pool),
	}
}

// CreateOrGetAccountInput represents the payload required to fetch or create an account.
type CreateOrGetAccountInput struct {
	Email             string
	FullName          string
	Provider          string
	ProviderAccountID string
	Thumbnail         *string
}

// CreateOrGetAccount upserts an account based on OAuth payload.
func (s *AccountService) CreateOrGetAccount(ctx echo.Context, input CreateOrGetAccountInput) (*sqldb.Account, error) {
	firstName, lastName := splitName(input.FullName)

	params := &sqldb.UpsertAccountParams{
		Email:             strings.TrimSpace(input.Email),
		FirstName:         firstName,
		LastName:          lastName,
		Provider:          strings.TrimSpace(input.Provider),
		ProviderAccountID: strings.TrimSpace(input.ProviderAccountID),
		LastLoginAt: pgtype.Timestamptz{
			Time:  time.Now().UTC(),
			Valid: true,
		},
	}

	if input.Thumbnail != nil && strings.TrimSpace(*input.Thumbnail) != "" {
		params.Thumbnail = pgtype.Text{String: strings.TrimSpace(*input.Thumbnail), Valid: true}
	}

	tx, err := s.pool.Begin(ctx.Request().Context())
	if err != nil {
		return nil, err
	}

	queries := s.queries.WithTx(tx)
	account, err := queries.UpsertAccount(ctx.Request().Context(), params)
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

	return account, nil
}

// GetAccountByID fetches an account using its UUID string.
func (s *AccountService) GetAccountByID(ctx echo.Context, id string) (*sqldb.Account, error) {
	pgID, err := parseUUID(id)
	if err != nil {
		return nil, ErrInvalidAccountID
	}

	account, err := s.queries.GetAccountByID(ctx.Request().Context(), pgID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrAccountNotFound
		}
		return nil, err
	}

	return account, nil
}

func splitName(full string) (string, string) {
	full = strings.TrimSpace(full)
	if full == "" {
		return "Unknown", ""
	}

	parts := strings.Fields(full)
	if len(parts) == 1 {
		return parts[0], ""
	}

	return parts[0], strings.Join(parts[1:], " ")
}
