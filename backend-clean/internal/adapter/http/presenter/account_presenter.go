// Package presenter contains HTTP presenters that implement output ports.
package presenter

import (
	"context"
	"strings"
	"time"

	openapi "immortal-architecture-cqrs/backend/internal/adapter/http/generated/openapi"
	"immortal-architecture-cqrs/backend/internal/domain/account"
	"immortal-architecture-cqrs/backend/internal/port"
)

// AccountPresenter converts domain account to OpenAPI response.
type AccountPresenter struct {
	account *openapi.ModelsAccountResponse
}

var _ port.AccountOutputPort = (*AccountPresenter)(nil)

// NewAccountPresenter creates a new AccountPresenter.
func NewAccountPresenter() *AccountPresenter {
	return &AccountPresenter{}
}

// PresentAccount stores converted account response.
func (p *AccountPresenter) PresentAccount(_ context.Context, a *account.Account) error {
	var lastLogin time.Time
	if a.LastLoginAt != nil {
		lastLogin = *a.LastLoginAt
	}
	p.account = &openapi.ModelsAccountResponse{
		Id:          a.ID,
		Email:       a.Email.String(),
		FirstName:   a.FirstName,
		LastName:    a.LastName,
		FullName:    strings.TrimSpace(a.FirstName + " " + a.LastName),
		Thumbnail:   strPtrOrNil(a.Thumbnail),
		LastLoginAt: lastLogin,
		CreatedAt:   a.CreatedAt,
		UpdatedAt:   a.UpdatedAt,
	}
	return nil
}

// Response returns the last account response.
func (p *AccountPresenter) Response() *openapi.ModelsAccountResponse {
	return p.account
}

func strPtrOrNil(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
