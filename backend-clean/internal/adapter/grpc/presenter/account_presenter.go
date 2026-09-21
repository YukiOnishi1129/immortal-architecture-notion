// Package presenter implements gRPC output ports.
package presenter

import (
	"context"
	"sync"

	"google.golang.org/protobuf/types/known/timestamppb"

	"immortal-architecture-cqrs/backend/internal/adapter/grpc/generated/accountpb"
	"immortal-architecture-cqrs/backend/internal/domain/account"
	"immortal-architecture-cqrs/backend/internal/port"
)

// AccountPresenter implements port.AccountOutputPort for gRPC.
type AccountPresenter struct {
	mu       sync.RWMutex
	response *accountpb.AccountResponse
}

var _ port.AccountOutputPort = (*AccountPresenter)(nil)

// NewAccountPresenter creates a new gRPC account presenter.
func NewAccountPresenter() *AccountPresenter {
	return &AccountPresenter{}
}

// PresentAccount converts domain account to gRPC response and stores it.
func (p *AccountPresenter) PresentAccount(_ context.Context, acc *account.Account) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	thumbnail := ""
	if acc.Thumbnail != "" {
		thumbnail = acc.Thumbnail
	}

	fullName := acc.FirstName + " " + acc.LastName

	var lastLoginAt *timestamppb.Timestamp
	if acc.LastLoginAt != nil {
		lastLoginAt = timestamppb.New(*acc.LastLoginAt)
	}

	p.response = &accountpb.AccountResponse{
		Id:          acc.ID,
		Email:       string(acc.Email),
		FirstName:   acc.FirstName,
		LastName:    acc.LastName,
		FullName:    fullName,
		Thumbnail:   &thumbnail,
		LastLoginAt: lastLoginAt,
		CreatedAt:   timestamppb.New(acc.CreatedAt),
		UpdatedAt:   timestamppb.New(acc.UpdatedAt),
	}

	return nil
}

// Response returns the stored gRPC response.
func (p *AccountPresenter) Response() *accountpb.AccountResponse {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.response
}
