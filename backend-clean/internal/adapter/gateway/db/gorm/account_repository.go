package gorm

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"immortal-architecture-cqrs/backend/internal/domain/account"
	domainerr "immortal-architecture-cqrs/backend/internal/domain/errors"
	"immortal-architecture-cqrs/backend/internal/port"
)

// AccountRepository implements account persistence using GORM.
type AccountRepository struct {
	db *gorm.DB
}

var _ port.AccountRepository = (*AccountRepository)(nil)

// NewAccountRepository creates AccountRepository with GORM.
func NewAccountRepository(db *gorm.DB) *AccountRepository {
	return &AccountRepository{
		db: db,
	}
}

// UpsertOAuthAccount inserts or updates an OAuth account using GORM.
func (r *AccountRepository) UpsertOAuthAccount(ctx context.Context, input account.OAuthAccountInput) (*account.Account, error) {
	var dbAccount Account

	// Check if account exists by email
	result := r.db.WithContext(ctx).Where("email = ?", input.Email).First(&dbAccount)

	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, result.Error
	}

	isNew := errors.Is(result.Error, gorm.ErrRecordNotFound)

	if isNew {
		// Create new account
		dbAccount = Account{
			Email:             input.Email,
			FirstName:         input.FirstName,
			LastName:          input.LastName,
			Provider:          input.Provider,
			ProviderAccountID: input.ProviderAccountID,
			Thumbnail:         input.Thumbnail,
		}

		if err := r.db.WithContext(ctx).Create(&dbAccount).Error; err != nil {
			return nil, err
		}
	} else {
		// Update existing account
		updates := map[string]interface{}{
			"first_name":          input.FirstName,
			"last_name":           input.LastName,
			"provider":            input.Provider,
			"provider_account_id": input.ProviderAccountID,
		}
		if input.Thumbnail != nil {
			updates["thumbnail"] = *input.Thumbnail
		}

		if err := r.db.WithContext(ctx).Model(&dbAccount).Updates(updates).Error; err != nil {
			return nil, err
		}

		// Reload to get updated values
		if err := r.db.WithContext(ctx).First(&dbAccount, "id = ?", dbAccount.ID).Error; err != nil {
			return nil, err
		}
	}

	return toDomainAccount(&dbAccount)
}

// GetByID retrieves an account by ID using GORM.
func (r *AccountRepository) GetByID(ctx context.Context, id string) (*account.Account, error) {
	var dbAccount Account

	if err := r.db.WithContext(ctx).First(&dbAccount, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainerr.ErrNotFound
		}
		return nil, err
	}

	return toDomainAccount(&dbAccount)
}

// GetByEmail retrieves an account by email using GORM.
func (r *AccountRepository) GetByEmail(ctx context.Context, emailAddr string) (*account.Account, error) {
	var dbAccount Account

	if err := r.db.WithContext(ctx).Where("email = ?", emailAddr).First(&dbAccount).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainerr.ErrNotFound
		}
		return nil, err
	}

	return toDomainAccount(&dbAccount)
}

// toDomainAccount converts GORM model to domain model.
func toDomainAccount(a *Account) (*account.Account, error) {
	email, err := account.ParseEmail(a.Email)
	if err != nil {
		return nil, err
	}

	thumbnail := ""
	if a.Thumbnail != nil {
		thumbnail = *a.Thumbnail
	}

	return &account.Account{
		ID:                a.ID,
		Email:             email,
		FirstName:         a.FirstName,
		LastName:          a.LastName,
		Provider:          a.Provider,
		ProviderAccountID: a.ProviderAccountID,
		Thumbnail:         thumbnail,
		LastLoginAt:       a.LastLoginAt,
		CreatedAt:         a.CreatedAt,
		UpdatedAt:         a.UpdatedAt,
	}, nil
}
