// Package mock provides test mocks for sqlc repositories.
package mock

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"immortal-architecture-cqrs/backend/internal/adapter/gateway/db/sqlc/generated"
)

// TemplateDBTX is a lightweight mock for sqlc.DBTX used in template repository tests.
type TemplateDBTX struct {
	templateRow *generated.Template
	detailRow   *generated.GetTemplateByIDRow
	FieldRow    *generated.Field
	rowErr      error
	execErr     error
	QueryErr    error
}

// NewTemplateDBTX creates a mock DBTX that always returns the given row/err.
func NewTemplateDBTX(templateRow *generated.Template, detailRow *generated.GetTemplateByIDRow, rowErr, execErr error) *TemplateDBTX {
	return &TemplateDBTX{
		templateRow: templateRow,
		detailRow:   detailRow,
		rowErr:      rowErr,
		execErr:     execErr,
	}
}

// Exec implements sqlc.DBTX interface.
func (m *TemplateDBTX) Exec(_ context.Context, _ string, _ ...interface{}) (pgconn.CommandTag, error) {
	return pgconn.CommandTag{}, m.execErr
}

// Query implements sqlc.DBTX interface.
func (m *TemplateDBTX) Query(_ context.Context, _ string, _ ...interface{}) (pgx.Rows, error) {
	if m.QueryErr != nil {
		return nil, m.QueryErr
	}
	return &emptyRows{}, nil
}

// QueryRow implements sqlc.DBTX interface.
func (m *TemplateDBTX) QueryRow(_ context.Context, _ string, _ ...interface{}) pgx.Row {
	return &templateRow{templateRow: m.templateRow, detailRow: m.detailRow, fieldRow: m.FieldRow, err: m.rowErr}
}

type templateRow struct {
	templateRow *generated.Template
	detailRow   *generated.GetTemplateByIDRow
	fieldRow    *generated.Field
	err         error
}

func (m *templateRow) Scan(dest ...interface{}) error {
	if m.err != nil {
		return m.err
	}
	switch len(dest) {
	case 5: // Field
		if m.fieldRow == nil {
			return errors.New("fieldRow is nil")
		}
		setUUID(dest[0], m.fieldRow.ID)
		setUUID(dest[1], m.fieldRow.TemplateID)
		setString(dest[2], m.fieldRow.Label)
		setInt32Field(dest[3], m.fieldRow.Order)
		setBool(dest[4], m.fieldRow.IsRequired)
	case 4: // Template
		setUUID(dest[0], m.templateRow.ID)
		setString(dest[1], m.templateRow.Name)
		setUUID(dest[2], m.templateRow.OwnerID)
		setTimestamptz(dest[3], m.templateRow.UpdatedAt)
	case 8: // GetTemplateByIDRow
		setUUID(dest[0], m.detailRow.ID)
		setString(dest[1], m.detailRow.Name)
		setUUID(dest[2], m.detailRow.OwnerID)
		setTimestamptz(dest[3], m.detailRow.UpdatedAt)
		setString(dest[4], m.detailRow.OwnerFirstName)
		setString(dest[5], m.detailRow.OwnerLastName)
		setText(dest[6], m.detailRow.OwnerThumbnail)
		setBool(dest[7], m.detailRow.IsUsed)
	default:
		return errors.New("unexpected scan args")
	}
	return nil
}

func (m *templateRow) FieldDescriptions() []pgconn.FieldDescription { return nil }
func (m *templateRow) RawValues() [][]byte                          { return nil }
func (m *templateRow) Value(_ int) (interface{}, error)             { return nil, nil }
func (m *templateRow) Err() error                                   { return m.err }

type emptyRows struct{}

func (r *emptyRows) Close()                                       {}
func (r *emptyRows) Next() bool                                   { return false }
func (r *emptyRows) Err() error                                   { return nil }
func (r *emptyRows) CommandTag() pgconn.CommandTag                { return pgconn.CommandTag{} }
func (r *emptyRows) FieldDescriptions() []pgconn.FieldDescription { return nil }
func (r *emptyRows) Values() ([]interface{}, error)               { return nil, nil }
func (r *emptyRows) RawValues() [][]byte                          { return nil }
func (r *emptyRows) Scan(_ ...interface{}) error                  { return nil }
func (r *emptyRows) Conn() *pgx.Conn                              { return nil }

func setInt32Field(ptr interface{}, v int32) {
	if dest, ok := ptr.(*int32); ok {
		*dest = v
	}
}
