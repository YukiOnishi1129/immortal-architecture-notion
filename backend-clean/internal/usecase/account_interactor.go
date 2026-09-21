// Package usecase contains application use case implementations.
package usecase

import (
	"context"

	"immortal-architecture-cqrs/backend/internal/domain/account"
	"immortal-architecture-cqrs/backend/internal/port"
)

// AccountInteractor handles account use cases.
type AccountInteractor struct {
	repo   port.AccountRepository
	output port.AccountOutputPort
}

var _ port.AccountInputPort = (*AccountInteractor)(nil)

// NewAccountInteractor creates AccountInteractor.
func NewAccountInteractor(repo port.AccountRepository, output port.AccountOutputPort) *AccountInteractor {
	return &AccountInteractor{repo: repo, output: output}
}

// CreateOrGet handles upsert/get of OAuth account.
func (u *AccountInteractor) CreateOrGet(ctx context.Context, input account.OAuthAccountInput) error {
	email, err := account.ParseEmail(input.Email)
	if err != nil {
		return err
	}
	acc := account.Account{
		Email:             email,
		FirstName:         input.FirstName,
		LastName:          input.LastName,
		Provider:          input.Provider,
		ProviderAccountID: input.ProviderAccountID,
		Thumbnail:         valueOrEmpty(input.Thumbnail),
	}
	if err := account.Validate(acc); err != nil {
		return err
	}
	a, err := u.repo.UpsertOAuthAccount(ctx, input)
	if err != nil {
		return err
	}
	return u.output.PresentAccount(ctx, a)
}

// GetByID retrieves account by ID.
func (u *AccountInteractor) GetByID(ctx context.Context, id string) error {
	a, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	return u.output.PresentAccount(ctx, a)
}

// GetByEmail retrieves account by email.
func (u *AccountInteractor) GetByEmail(ctx context.Context, email string) error {
	// Validate email format
	_, err := account.ParseEmail(email)
	if err != nil {
		return err
	}

	a, err := u.repo.GetByEmail(ctx, email)
	if err != nil {
		return err
	}
	return u.output.PresentAccount(ctx, a)
}

func valueOrEmpty(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
