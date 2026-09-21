// Package controller implements gRPC controllers.
package controller

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"immortal-architecture-cqrs/backend/internal/adapter/grpc/generated/accountpb"
	grpcpresenter "immortal-architecture-cqrs/backend/internal/adapter/grpc/presenter"
	"immortal-architecture-cqrs/backend/internal/domain/account"
	domainerr "immortal-architecture-cqrs/backend/internal/domain/errors"
	"immortal-architecture-cqrs/backend/internal/port"
)

// AccountController implements accountpb.AccountServiceServer.
type AccountController struct {
	accountpb.UnimplementedAccountServiceServer
	inputFactory  func(port.AccountRepository, port.AccountOutputPort) port.AccountInputPort
	outputFactory func() *grpcpresenter.AccountPresenter
	repoFactory   func() port.AccountRepository
}

// NewAccountController creates a new gRPC account controller.
func NewAccountController(
	inputFactory func(port.AccountRepository, port.AccountOutputPort) port.AccountInputPort,
	outputFactory func() *grpcpresenter.AccountPresenter,
	repoFactory func() port.AccountRepository,
) *AccountController {
	return &AccountController{
		inputFactory:  inputFactory,
		outputFactory: outputFactory,
		repoFactory:   repoFactory,
	}
}

// GetAccountByID retrieves an account by ID.
func (s *AccountController) GetAccountByID(ctx context.Context, req *accountpb.GetAccountByIdRequest) (*accountpb.AccountResponse, error) {
	presenter := s.outputFactory()
	input := s.inputFactory(s.repoFactory(), presenter)

	if err := input.GetByID(ctx, req.GetAccountId()); err != nil {
		return nil, handleError(err)
	}

	return presenter.Response(), nil
}

// GetAccountByEmail retrieves an account by email.
func (s *AccountController) GetAccountByEmail(ctx context.Context, req *accountpb.GetAccountByEmailRequest) (*accountpb.AccountResponse, error) {
	presenter := s.outputFactory()
	input := s.inputFactory(s.repoFactory(), presenter)

	if err := input.GetByEmail(ctx, req.GetEmail()); err != nil {
		return nil, handleError(err)
	}

	return presenter.Response(), nil
}

// CreateOrGetAccount creates or gets an OAuth account.
func (s *AccountController) CreateOrGetAccount(ctx context.Context, req *accountpb.CreateOrGetAccountRequest) (*accountpb.AccountResponse, error) {
	presenter := s.outputFactory()
	input := s.inputFactory(s.repoFactory(), presenter)

	thumbnail := req.GetThumbnail()
	var thumbnailPtr *string
	if thumbnail != "" {
		thumbnailPtr = &thumbnail
	}

	oauthInput := account.OAuthAccountInput{
		Email:             req.GetEmail(),
		FirstName:         req.GetFirstName(),
		LastName:          req.GetLastName(),
		Provider:          req.GetProvider(),
		ProviderAccountID: req.GetProviderAccountId(),
		Thumbnail:         thumbnailPtr,
	}

	if err := input.CreateOrGet(ctx, oauthInput); err != nil {
		return nil, handleError(err)
	}

	return presenter.Response(), nil
}

// handleError converts domain errors to gRPC status codes.
func handleError(err error) error {
	if errors.Is(err, domainerr.ErrNotFound) {
		return status.Error(codes.NotFound, err.Error())
	}
	if errors.Is(err, domainerr.ErrUnauthorized) {
		return status.Error(codes.Unauthenticated, err.Error())
	}
	if errors.Is(err, account.ErrInvalidEmail) {
		return status.Error(codes.InvalidArgument, err.Error())
	}
	if errors.Is(err, domainerr.ErrProviderRequired) {
		return status.Error(codes.InvalidArgument, err.Error())
	}
	return status.Error(codes.Internal, "internal server error")
}
