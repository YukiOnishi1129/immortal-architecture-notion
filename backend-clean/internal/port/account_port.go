// Package port defines application ports (interfaces).
package port

import (
	"context"

	"immortal-architecture-cqrs/backend/internal/domain/account"
)

// AccountInputPort defines account use case input methods.
type AccountInputPort interface {
	CreateOrGet(ctx context.Context, input account.OAuthAccountInput) error
	GetByID(ctx context.Context, id string) error
	GetByEmail(ctx context.Context, email string) error
}

// AccountOutputPort defines presenter for accounts.
type AccountOutputPort interface {
	PresentAccount(ctx context.Context, account *account.Account) error
}

// AccountRepository abstracts account persistence.
type AccountRepository interface {
	UpsertOAuthAccount(ctx context.Context, input account.OAuthAccountInput) (*account.Account, error)
	GetByID(ctx context.Context, id string) (*account.Account, error)
	GetByEmail(ctx context.Context, email string) (*account.Account, error)
}
