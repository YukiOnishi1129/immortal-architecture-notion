package mock

import (
	"context"

	"immortal-architecture-cqrs/backend/internal/domain/account"
	"immortal-architecture-cqrs/backend/internal/port"
)

// AccountInputStub is a lightweight stub for account use case input.
type AccountInputStub struct {
	CreateErr error
	GetErr    error
	Output    port.AccountOutputPort
}

func (s *AccountInputStub) CreateOrGet(ctx context.Context, input account.OAuthAccountInput) error {
	if s.Output != nil && s.CreateErr == nil {
		_ = s.Output.PresentAccount(ctx, &account.Account{
			ID:        "acc-1",
			Email:     account.Email(input.Email),
			FirstName: input.FirstName,
			Provider:  input.Provider,
		})
	}
	return s.CreateErr
}

func (s *AccountInputStub) GetByID(ctx context.Context, id string) error {
	if s.Output != nil && s.GetErr == nil {
		_ = s.Output.PresentAccount(ctx, &account.Account{
			ID:        id,
			Email:     "user@example.com",
			FirstName: "Taro",
			Provider:  "google",
		})
	}
	return s.GetErr
}

func (s *AccountInputStub) GetByEmail(ctx context.Context, email string) error {
	if s.Output != nil && s.GetErr == nil {
		_ = s.Output.PresentAccount(ctx, &account.Account{
			ID:        "acc-1",
			Email:     account.Email(email),
			FirstName: "Taro",
			Provider:  "google",
		})
	}
	return s.GetErr
}
