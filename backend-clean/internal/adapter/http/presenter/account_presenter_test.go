package presenter

import (
	"context"
	"testing"
	"time"

	"immortal-architecture-cqrs/backend/internal/domain/account"
)

func TestAccountPresenter_PresentAccount(t *testing.T) {
	now := time.Now()
	last := now.Add(-time.Hour)
	tests := []struct {
		name         string
		acc          *account.Account
		wantID       string
		wantEmail    string
		wantFullName string
		wantThumb    *string
	}{
		{
			name: "[Success] full info",
			acc: &account.Account{
				ID:          "acc-1",
				Email:       "user@example.com",
				FirstName:   "Taro",
				LastName:    "Yamada",
				Thumbnail:   "thumb",
				LastLoginAt: &last,
				CreatedAt:   now,
				UpdatedAt:   now,
			},
			wantID:       "acc-1",
			wantEmail:    "user@example.com",
			wantFullName: "Taro Yamada",
			wantThumb:    strPtr("thumb"),
		},
		{
			name: "[Success] empty names",
			acc: &account.Account{
				ID:    "acc-2",
				Email: "no@example.com",
			},
			wantID:    "acc-2",
			wantEmail: "no@example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewAccountPresenter()
			if err := p.PresentAccount(context.Background(), tt.acc); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			resp := p.Response()
			if resp == nil || resp.Id != tt.wantID || resp.Email != tt.wantEmail || resp.FullName != tt.wantFullName {
				t.Fatalf("unexpected response: %+v", resp)
			}
			if tt.wantThumb != nil {
				if resp.Thumbnail == nil || *resp.Thumbnail != *tt.wantThumb {
					t.Fatalf("thumbnail mismatch")
				}
			}
		})
	}
}
