// Package factory provides constructors for driver-level wiring.
package factory

import (
	"github.com/jackc/pgx/v5/pgxpool"

	"immortal-architecture-notion/backend/internal/adapter/gateway/db/sqlc"
	// "immortal-architecture-notion/backend/internal/adapter/gateway/db/gorm"
	"immortal-architecture-notion/backend/internal/port"
)

// NewAccountRepoFactory returns a factory that creates AccountRepository.
//
// To switch ORM implementation (e.g., from sqlc to GORM):
// 1. Change the import (uncomment gorm import above)
// 2. Update this factory function
// 3. Pass gorm.DB instead of pgxpool.Pool
//
// All domain, use case, and adapter layers (HTTP/gRPC controllers, presenters)
// remain unchanged. This demonstrates Clean Architecture's changeability.
func NewAccountRepoFactory(pool *pgxpool.Pool) func() port.AccountRepository {
	return func() port.AccountRepository {
		// Current: sqlc implementation
		return sqlc.NewAccountRepository(pool)

		// To switch to GORM, replace above with:
		// return gorm.NewAccountRepository(db)
	}
}

// NewTemplateRepoFactory returns a factory that creates TemplateRepository.
func NewTemplateRepoFactory(pool *pgxpool.Pool) func() port.TemplateRepository {
	return func() port.TemplateRepository {
		return sqlc.NewTemplateRepository(pool)
	}
}

// NewNoteRepoFactory returns a factory that creates NoteRepository.
func NewNoteRepoFactory(pool *pgxpool.Pool) func() port.NoteRepository {
	return func() port.NoteRepository {
		return sqlc.NewNoteRepository(pool)
	}
}

// NewNoteReadModelRepoFactory returns a factory that creates NoteReadModelRepository.
func NewNoteReadModelRepoFactory(pool *pgxpool.Pool) func() port.NoteReadModelRepository {
	return func() port.NoteReadModelRepository {
		return sqlc.NewNoteReadModelRepository(pool)
	}
}
