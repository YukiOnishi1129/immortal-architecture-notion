// Package mock provides test mocks for sqlc repositories.
package mock

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"

	"immortal-architecture-cqrs/backend/internal/adapter/gateway/db/sqlc/generated"
)

// AccountDBTX is a lightweight mock for sqlc.DBTX to test account repository.
type AccountDBTX struct {
	row *generated.Account
	err error
}

// NewAccountDBTX creates a mock DBTX that always returns the given row/err.
func NewAccountDBTX(row *generated.Account, err error) *AccountDBTX {
	return &AccountDBTX{row: row, err: err}
}

// Exec implements sqlc.DBTX interface.
func (m *AccountDBTX) Exec(_ context.Context, _ string, _ ...interface{}) (pgconn.CommandTag, error) {
	return pgconn.CommandTag{}, nil
}

// Query implements sqlc.DBTX interface.
func (m *AccountDBTX) Query(_ context.Context, _ string, _ ...interface{}) (pgx.Rows, error) {
	return nil, nil
}

// QueryRow implements sqlc.DBTX interface.
func (m *AccountDBTX) QueryRow(_ context.Context, _ string, _ ...interface{}) pgx.Row {
	return &accountRow{row: m.row, err: m.err}
}

// accountRow mocks pgx.Row for account queries.
type accountRow struct {
	row *generated.Account
	err error
}

func (m *accountRow) Scan(dest ...interface{}) error {
	if m.err != nil {
		return m.err
	}
	if len(dest) != 11 {
		return errors.New("unexpected scan args")
	}
	setUUID(dest[0], m.row.ID)
	setString(dest[1], m.row.Email)
	setString(dest[2], m.row.FirstName)
	setString(dest[3], m.row.LastName)
	setBool(dest[4], m.row.IsActive)
	setString(dest[5], m.row.Provider)
	setString(dest[6], m.row.ProviderAccountID)
	setText(dest[7], m.row.Thumbnail)
	setTimestamptz(dest[8], m.row.LastLoginAt)
	setTimestamptz(dest[9], m.row.CreatedAt)
	setTimestamptz(dest[10], m.row.UpdatedAt)
	return nil
}

// pgx.Row interface other methods (unused).
func (m *accountRow) FieldDescriptions() []pgconn.FieldDescription { return nil }
func (m *accountRow) RawValues() [][]byte                          { return nil }
func (m *accountRow) Value(_ int) (interface{}, error)             { return nil, nil }
func (m *accountRow) Err() error                                   { return m.err }

func setUUID(ptr interface{}, v pgtype.UUID) {
	if dest, ok := ptr.(*pgtype.UUID); ok {
		*dest = v
	}
}

func setString(ptr interface{}, v string) {
	if dest, ok := ptr.(*string); ok {
		*dest = v
	}
}

func setBool(ptr interface{}, v bool) {
	if dest, ok := ptr.(*bool); ok {
		*dest = v
	}
}

func setText(ptr interface{}, v pgtype.Text) {
	if dest, ok := ptr.(*pgtype.Text); ok {
		*dest = v
	}
}

func setTimestamptz(ptr interface{}, v pgtype.Timestamptz) {
	if dest, ok := ptr.(*pgtype.Timestamptz); ok {
		*dest = v
	}
}
