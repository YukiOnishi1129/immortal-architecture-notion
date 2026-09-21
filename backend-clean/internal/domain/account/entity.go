// Package account holds account domain models.
package account

import "time"

// Account is the aggregate root representing a user.
type Account struct {
	ID                string
	Email             Email
	FirstName         string
	LastName          string
	IsActive          bool
	Provider          string
	ProviderAccountID string
	Thumbnail         string
	LastLoginAt       *time.Time
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

// OAuthAccountInput describes account info from OAuth provider.
type OAuthAccountInput struct {
	Email             string
	FirstName         string
	LastName          string
	Provider          string
	ProviderAccountID string
	Thumbnail         *string
}
