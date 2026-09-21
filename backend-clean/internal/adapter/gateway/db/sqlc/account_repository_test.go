package sqlc

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"

	mockdb "immortal-architecture-cqrs/backend/internal/adapter/gateway/db/sqlc/mock"
	"immortal-architecture-cqrs/backend/internal/adapter/gateway/db/sqlc/generated"
	"immortal-architecture-cqrs/backend/internal/domain/account"
)

func TestToDomainAccount(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	tests := []struct {
		name    string
		row     *generated.Account
		wantErr bool
	}{
		{
			name: "[Success] maps nullable fields",
			row: &generated.Account{
				ID:                pgtype.UUID{Bytes: [16]byte{1}, Valid: true},
				Email:             "user@example.com",
				FirstName:         "Taro",
				LastName:          "Yamada",
				IsActive:          true,
				Provider:          "google",
				ProviderAccountID: "pid",
				Thumbnail:         pgtype.Text{String: "thumb", Valid: true},
				LastLoginAt:       pgtype.Timestamptz{Time: now, Valid: true},
				CreatedAt:         pgtype.Timestamptz{Time: now, Valid: true},
				UpdatedAt:         pgtype.Timestamptz{Time: now, Valid: true},
			},
		},
		{
			name: "[Fail] invalid email",
			row: &generated.Account{
				ID:        pgtype.UUID{Bytes: [16]byte{1}, Valid: true},
				Email:     "bad-email",
				CreatedAt: pgtype.Timestamptz{Time: now, Valid: true},
				UpdatedAt: pgtype.Timestamptz{Time: now, Valid: true},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			acc, err := toDomainAccount(tt.row)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if acc.Email != account.Email(tt.row.Email) {
				t.Fatalf("email = %s, want %s", acc.Email, tt.row.Email)
			}
			if acc.Thumbnail != tt.row.Thumbnail.String {
				t.Fatalf("thumb = %s, want %s", acc.Thumbnail, tt.row.Thumbnail.String)
			}
			if acc.LastLoginAt == nil || !acc.LastLoginAt.Equal(now) {
				t.Fatalf("lastLoginAt = %+v, want %v", acc.LastLoginAt, now)
			}
		})
	}
}

func TestAccountRepository_UpsertOAuthAccount(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	baseRow := &generated.Account{
		ID:                pgtype.UUID{Bytes: [16]byte{1}, Valid: true},
		Email:             "user@example.com",
		FirstName:         "Taro",
		LastName:          "Yamada",
		IsActive:          true,
		Provider:          "google",
		ProviderAccountID: "pid",
		Thumbnail:         pgtype.Text{String: "thumb", Valid: true},
		LastLoginAt:       pgtype.Timestamptz{Time: now, Valid: true},
		CreatedAt:         pgtype.Timestamptz{Time: now, Valid: true},
		UpdatedAt:         pgtype.Timestamptz{Time: now, Valid: true},
	}
	tests := []struct {
		name    string
		row     *generated.Account
		rowErr  error
		wantErr bool
	}{
		{name: "[Success] upsert returns domain", row: baseRow},
		{name: "[Fail] invalid email", row: func() *generated.Account { r := *baseRow; r.Email = "bad"; return &r }(), wantErr: true},
		{name: "[Fail] query error", rowErr: errors.New("db error"), wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := mockdb.NewAccountDBTX(tt.row, tt.rowErr)
			repo := &AccountRepository{queries: generated.New(mock)}
			input := account.OAuthAccountInput{
				Email:             baseRow.Email,
				FirstName:         baseRow.FirstName,
				LastName:          baseRow.LastName,
				Provider:          baseRow.Provider,
				ProviderAccountID: baseRow.ProviderAccountID,
				Thumbnail:         nil,
			}
			acc, err := repo.UpsertOAuthAccount(context.Background(), input)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if acc.Email != account.Email(baseRow.Email) {
				t.Fatalf("email = %s, want %s", acc.Email, baseRow.Email)
			}
		})
	}
}

func TestAccountRepository_GetByID(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	row := &generated.Account{
		ID:                pgtype.UUID{Bytes: [16]byte{1}, Valid: true},
		Email:             "user@example.com",
		FirstName:         "Taro",
		LastName:          "Yamada",
		IsActive:          true,
		Provider:          "google",
		ProviderAccountID: "pid",
		Thumbnail:         pgtype.Text{String: "thumb", Valid: true},
		LastLoginAt:       pgtype.Timestamptz{Time: now, Valid: true},
		CreatedAt:         pgtype.Timestamptz{Time: now, Valid: true},
		UpdatedAt:         pgtype.Timestamptz{Time: now, Valid: true},
	}

	tests := []struct {
		name    string
		id      string
		row     *generated.Account
		rowErr  error
		wantErr bool
	}{
		{name: "[Success] GetByID returns domain", id: row.ID.String(), row: row},
		{name: "[Fail] GetByID invalid uuid", id: "not-uuid", row: row, wantErr: true},
		{name: "[Fail] GetByID query error", id: row.ID.String(), rowErr: errors.New("db error"), wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := mockdb.NewAccountDBTX(tt.row, tt.rowErr)
			repo := &AccountRepository{queries: generated.New(mock)}
			acc, err := repo.GetByID(context.Background(), tt.id)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if acc.Email != account.Email(row.Email) {
				t.Fatalf("email = %s, want %s", acc.Email, row.Email)
			}
		})
	}
}

func TestAccountRepository_GetByEmail(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	row := &generated.Account{
		ID:                pgtype.UUID{Bytes: [16]byte{1}, Valid: true},
		Email:             "user@example.com",
		FirstName:         "Taro",
		LastName:          "Yamada",
		IsActive:          true,
		Provider:          "google",
		ProviderAccountID: "pid",
		Thumbnail:         pgtype.Text{String: "thumb", Valid: true},
		LastLoginAt:       pgtype.Timestamptz{Time: now, Valid: true},
		CreatedAt:         pgtype.Timestamptz{Time: now, Valid: true},
		UpdatedAt:         pgtype.Timestamptz{Time: now, Valid: true},
	}

	tests := []struct {
		name    string
		email   string
		row     *generated.Account
		rowErr  error
		wantErr bool
	}{
		{name: "[Success] GetByEmail returns domain", email: "user@example.com", row: row},
		{name: "[Fail] GetByEmail invalid email", email: "user@example.com", row: func() *generated.Account { r := *row; r.Email = "bad"; return &r }(), wantErr: true},
		{name: "[Fail] GetByEmail query error", email: "user@example.com", rowErr: errors.New("db error"), wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := mockdb.NewAccountDBTX(tt.row, tt.rowErr)
			repo := &AccountRepository{queries: generated.New(mock)}
			acc, err := repo.GetByEmail(context.Background(), tt.email)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if acc.Email != account.Email(row.Email) {
				t.Fatalf("email = %s, want %s", acc.Email, row.Email)
			}
		})
	}
}
