package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"

	"immortal-architecture-cqrs/backend/internal/domain/account"
	domainerr "immortal-architecture-cqrs/backend/internal/domain/errors"
	uc "immortal-architecture-cqrs/backend/internal/usecase"
	mockusecase "immortal-architecture-cqrs/backend/internal/usecase/mock"
)

func TestAccountInteractor_CreateOrGet(t *testing.T) {
	tests := []struct {
		name      string
		input     account.OAuthAccountInput
		repoAcc   *account.Account
		repoErr   error
		wantError error
	}{
		{
			name: "[Success] upsert account",
			input: account.OAuthAccountInput{
				Email:             "user@example.com",
				FirstName:         "Taro",
				LastName:          "Yamada",
				Provider:          "google",
				ProviderAccountID: "pid",
			},
			repoAcc: &account.Account{ID: "acc-1"},
		},
		{
			name: "[Success] return existing",
			input: account.OAuthAccountInput{
				Email:             "user@example.com",
				FirstName:         "Hanako",
				LastName:          "Yamada",
				Provider:          "google",
				ProviderAccountID: "pid",
			},
			repoAcc: &account.Account{ID: "acc-1", FirstName: "Hanako"},
		},
		{
			name: "[Fail] invalid email",
			input: account.OAuthAccountInput{
				Email:             "invalid",
				Provider:          "google",
				ProviderAccountID: "pid",
			},
			wantError: account.ErrInvalidEmail,
		},
		{
			name: "[Fail] missing provider",
			input: account.OAuthAccountInput{
				Email:             "user@example.com",
				FirstName:         "Taro",
				Provider:          "",
				ProviderAccountID: "pid",
			},
			wantError: domainerr.ErrProviderRequired,
		},
		{
			name: "[Fail] repo error",
			input: account.OAuthAccountInput{
				Email:             "user@example.com",
				FirstName:         "Taro",
				Provider:          "google",
				ProviderAccountID: "pid",
			},
			repoErr:   errors.New("repo err"),
			wantError: errors.New("repo err"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mockusecase.NewMockAccountRepository(ctrl)
			out := mockusecase.NewMockAccountOutputPort(ctrl)

			shouldCallRepo := tt.wantError == nil || (tt.repoErr != nil && tt.wantError.Error() == tt.repoErr.Error())
			if shouldCallRepo {
				repo.EXPECT().UpsertOAuthAccount(gomock.Any(), tt.input).Return(tt.repoAcc, tt.repoErr)
			}
			if tt.wantError == nil && tt.repoErr == nil {
				out.EXPECT().PresentAccount(gomock.Any(), tt.repoAcc).Return(nil)
			}

			interactor := uc.NewAccountInteractor(repo, out)
			err := interactor.CreateOrGet(context.Background(), tt.input)

			if tt.wantError == nil && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.wantError != nil && (err == nil || tt.wantError.Error() != err.Error()) {
				t.Fatalf("want %v, got %v", tt.wantError, err)
			}
		})
	}
}

func TestAccountInteractor_GetByID(t *testing.T) {
	tests := []struct {
		name      string
		id        string
		repoAcc   *account.Account
		repoErr   error
		wantError error
	}{
		{
			name:    "[Success] get by id",
			id:      "acc-1",
			repoAcc: &account.Account{ID: "acc-1"},
		},
		{
			name:      "[Fail] not found",
			id:        "missing",
			repoErr:   domainerr.ErrNotFound,
			wantError: domainerr.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mockusecase.NewMockAccountRepository(ctrl)
			out := mockusecase.NewMockAccountOutputPort(ctrl)

			repo.EXPECT().GetByID(gomock.Any(), tt.id).Return(tt.repoAcc, tt.repoErr)
			if tt.repoErr == nil {
				out.EXPECT().PresentAccount(gomock.Any(), tt.repoAcc).Return(nil)
			}

			interactor := uc.NewAccountInteractor(repo, out)
			err := interactor.GetByID(context.Background(), tt.id)

			if tt.wantError == nil && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.wantError != nil && !errors.Is(err, tt.wantError) {
				t.Fatalf("want %v, got %v", tt.wantError, err)
			}
		})
	}
}

func TestAccountInteractor_GetByEmail(t *testing.T) {
	email, _ := account.ParseEmail("user@example.com")
	tests := []struct {
		name      string
		email     string
		repoAcc   *account.Account
		repoErr   error
		wantError error
	}{
		{
			name:  "[Success] get by email",
			email: "user@example.com",
			repoAcc: &account.Account{
				ID:    "acc-1",
				Email: email,
			},
		},
		{
			name:      "[Fail] invalid email format",
			email:     "invalid-email",
			wantError: account.ErrInvalidEmail,
		},
		{
			name:      "[Fail] not found",
			email:     "missing@example.com",
			repoErr:   domainerr.ErrNotFound,
			wantError: domainerr.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mockusecase.NewMockAccountRepository(ctrl)
			out := mockusecase.NewMockAccountOutputPort(ctrl)

			shouldCallRepo := tt.wantError == nil || tt.wantError == domainerr.ErrNotFound
			if shouldCallRepo {
				repo.EXPECT().GetByEmail(gomock.Any(), tt.email).Return(tt.repoAcc, tt.repoErr)
			}
			if tt.wantError == nil && tt.repoErr == nil {
				out.EXPECT().PresentAccount(gomock.Any(), tt.repoAcc).Return(nil)
			}

			interactor := uc.NewAccountInteractor(repo, out)
			err := interactor.GetByEmail(context.Background(), tt.email)

			if tt.wantError == nil && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.wantError != nil && !errors.Is(err, tt.wantError) {
				t.Fatalf("want %v, got %v", tt.wantError, err)
			}
		})
	}
}
