package account

import (
	"errors"
	"testing"

	domainerr "immortal-architecture-cqrs/backend/internal/domain/errors"
)

func TestParseEmail(t *testing.T) {
	tests := []struct {
		name      string
		raw       string
		wantError error
	}{
		{name: "[Success] valid email", raw: "user@example.com"},
		{name: "[Fail] missing at", raw: "invalid", wantError: ErrInvalidEmail},
		{name: "[Fail] empty", raw: "   ", wantError: ErrInvalidEmail},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseEmail(tt.raw)
			if tt.wantError == nil && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.wantError != nil && !errors.Is(err, tt.wantError) {
				t.Fatalf("want %v, got %v", tt.wantError, err)
			}
		})
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name      string
		acc       Account
		wantError error
	}{
		{
			name: "[Success] valid account",
			acc: Account{
				FirstName:         "Taro",
				LastName:          "Yamada",
				Provider:          "google",
				ProviderAccountID: "pid",
			},
		},
		{
			name: "[Fail] missing name",
			acc: Account{
				FirstName:         "",
				LastName:          "",
				Provider:          "google",
				ProviderAccountID: "pid",
			},
			wantError: ErrInvalidName,
		},
		{
			name: "[Fail] missing provider",
			acc: Account{
				FirstName:         "Taro",
				LastName:          "Yamada",
				Provider:          "",
				ProviderAccountID: "pid",
			},
			wantError: domainerr.ErrProviderRequired,
		},
		{
			name: "[Fail] missing provider account",
			acc: Account{
				FirstName: "Taro",
				LastName:  "Yamada",
				Provider:  "google",
			},
			wantError: domainerr.ErrProviderAccountRequired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate(tt.acc)
			if tt.wantError == nil && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.wantError != nil && !errors.Is(err, tt.wantError) {
				t.Fatalf("want %v, got %v", tt.wantError, err)
			}
		})
	}
}

func TestUpdateProfile(t *testing.T) {
	tests := []struct {
		name      string
		current   Account
		input     OAuthAccountInput
		wantEmail Email
		wantError error
	}{
		{
			name: "[Success] merge fields",
			current: Account{
				Email:             "old@example.com",
				FirstName:         "Old",
				LastName:          "Name",
				Provider:          "google",
				ProviderAccountID: "pid",
			},
			input: OAuthAccountInput{
				Email:     "new@example.com",
				FirstName: "New",
				LastName:  "Name",
				Thumbnail: ptr("http://thumb"),
			},
			wantEmail: "new@example.com",
		},
		{
			name: "[Fail] invalid email",
			current: Account{
				Email: "old@example.com",
			},
			input: OAuthAccountInput{
				Email: "invalid",
			},
			wantError: ErrInvalidEmail,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updated, err := UpdateProfile(tt.current, tt.input)
			if tt.wantError == nil && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.wantError != nil && !errors.Is(err, tt.wantError) {
				t.Fatalf("want %v, got %v", tt.wantError, err)
			}
			if err == nil && updated.Email != tt.wantEmail {
				t.Fatalf("email not updated: %s", updated.Email)
			}
			if err == nil && tt.input.Thumbnail != nil && updated.Thumbnail != *tt.input.Thumbnail {
				t.Fatalf("thumbnail not updated: %s", updated.Thumbnail)
			}
			if err == nil && updated.FirstName != tt.input.FirstName {
				t.Fatalf("first name not updated: %s", updated.FirstName)
			}
		})
	}
}

func ptr[T any](v T) *T { return &v }
