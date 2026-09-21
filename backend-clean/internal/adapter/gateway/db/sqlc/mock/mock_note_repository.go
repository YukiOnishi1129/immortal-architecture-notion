// Package mock provides test mocks for sqlc repositories.
package mock

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"immortal-architecture-cqrs/backend/internal/adapter/gateway/db/sqlc/generated"
)

// NoteDBTX is a lightweight mock for sqlc.DBTX used in note repository tests.
type NoteDBTX struct {
	row        *generated.Note
	getRow     *generated.GetNoteByIDRow
	sectionRow *generated.Section
	rowErr     error
	execErr    error
	queryErr   error
	listNotes  []*generated.ListNotesRow
	sections   []*generated.Section
}

// NewNoteDBTX creates a mock DBTX that always returns the given row/err.
func NewNoteDBTX(row *generated.Note, rowErr, execErr error) *NoteDBTX {
	return &NoteDBTX{row: row, rowErr: rowErr, execErr: execErr}
}

// WithList allows configuring rows returned by ListNotes/ListSections.
func (m *NoteDBTX) WithList(notes []*generated.ListNotesRow, sections []*generated.Section, queryErr error) *NoteDBTX {
	m.listNotes = notes
	m.sections = sections
	m.queryErr = queryErr
	return m
}

// WithGetRow sets a GetNoteByIDRow for QueryRow scans requiring 11 columns.
func (m *NoteDBTX) WithGetRow(row *generated.GetNoteByIDRow) *NoteDBTX {
	m.getRow = row
	return m
}

// WithSectionRow sets Section row for UpdateSectionContent scans.
func (m *NoteDBTX) WithSectionRow(row *generated.Section) *NoteDBTX {
	m.sectionRow = row
	return m
}

// Exec implements sqlc.DBTX interface.
func (m *NoteDBTX) Exec(_ context.Context, _ string, _ ...interface{}) (pgconn.CommandTag, error) {
	return pgconn.CommandTag{}, m.execErr
}

// Query implements sqlc.DBTX interface.
func (m *NoteDBTX) Query(_ context.Context, _ string, args ...interface{}) (pgx.Rows, error) {
	if m.queryErr != nil {
		return nil, m.queryErr
	}
	// Heuristic: ListNotes has 4 args, ListSectionsByNote has 1 arg.
	if len(args) == 4 {
		return &noteRows{items: m.listNotes}, nil
	}
	return &sectionRows{items: m.sections}, nil
}

// QueryRow implements sqlc.DBTX interface.
func (m *NoteDBTX) QueryRow(_ context.Context, _ string, _ ...interface{}) pgx.Row {
	return &noteRow{row: m.row, getRow: m.getRow, secRow: m.sectionRow, err: m.rowErr}
}

type noteRow struct {
	row    *generated.Note
	getRow *generated.GetNoteByIDRow
	secRow *generated.Section
	err    error
}

func (m *noteRow) Scan(dest ...interface{}) error {
	if m.err != nil {
		return m.err
	}
	switch len(dest) {
	case 11:
		if m.getRow == nil {
			return errors.New("getRow is nil")
		}
		setUUID(dest[0], m.getRow.ID)
		setString(dest[1], m.getRow.Title)
		setUUID(dest[2], m.getRow.TemplateID)
		setUUID(dest[3], m.getRow.OwnerID)
		setString(dest[4], m.getRow.Status)
		setTimestamptz(dest[5], m.getRow.CreatedAt)
		setTimestamptz(dest[6], m.getRow.UpdatedAt)
		setString(dest[7], m.getRow.TemplateName)
		setString(dest[8], m.getRow.FirstName)
		setString(dest[9], m.getRow.LastName)
		setText(dest[10], m.getRow.OwnerThumbnail)
		return nil
	case 7:
		if m.row == nil {
			return errors.New("row is nil")
		}
		setUUID(dest[0], m.row.ID)
		setString(dest[1], m.row.Title)
		setUUID(dest[2], m.row.TemplateID)
		setUUID(dest[3], m.row.OwnerID)
		setString(dest[4], m.row.Status)
		setTimestamptz(dest[5], m.row.CreatedAt)
		setTimestamptz(dest[6], m.row.UpdatedAt)
		return nil
	case 4:
		if m.secRow == nil {
			return errors.New("sectionRow is nil")
		}
		setUUID(dest[0], m.secRow.ID)
		setUUID(dest[1], m.secRow.NoteID)
		setUUID(dest[2], m.secRow.FieldID)
		setString(dest[3], m.secRow.Content)
		return nil
	default:
		return errors.New("unexpected scan args")
	}
}

func (m *noteRow) FieldDescriptions() []pgconn.FieldDescription { return nil }
func (m *noteRow) RawValues() [][]byte                          { return nil }
func (m *noteRow) Value(_ int) (interface{}, error)             { return nil, nil }
func (m *noteRow) Err() error                                   { return m.err }

type noteRows struct {
	items []*generated.ListNotesRow
	idx   int
	err   error
}

func (r *noteRows) Close()                                       {}
func (r *noteRows) Next() bool                                   { r.idx++; return r.idx <= len(r.items) }
func (r *noteRows) Err() error                                   { return r.err }
func (r *noteRows) CommandTag() pgconn.CommandTag                { return pgconn.CommandTag{} }
func (r *noteRows) FieldDescriptions() []pgconn.FieldDescription { return nil }
func (r *noteRows) Values() ([]interface{}, error)               { return nil, nil }
func (r *noteRows) RawValues() [][]byte                          { return nil }
func (r *noteRows) Scan(dest ...interface{}) error {
	if r.idx == 0 || r.idx > len(r.items) {
		return errors.New("scan called out of range")
	}
	item := r.items[r.idx-1]
	if len(dest) != 11 {
		return errors.New("unexpected scan args")
	}
	setUUID(dest[0], item.ID)
	setString(dest[1], item.Title)
	setUUID(dest[2], item.TemplateID)
	setUUID(dest[3], item.OwnerID)
	setString(dest[4], item.Status)
	setTimestamptz(dest[5], item.CreatedAt)
	setTimestamptz(dest[6], item.UpdatedAt)
	setString(dest[7], item.TemplateName)
	setString(dest[8], item.FirstName)
	setString(dest[9], item.LastName)
	setText(dest[10], item.OwnerThumbnail)
	return nil
}
func (r *noteRows) Conn() *pgx.Conn { return nil }

type sectionRows struct {
	items []*generated.Section
	idx   int
	err   error
}

func (r *sectionRows) Close()                                       {}
func (r *sectionRows) Next() bool                                   { r.idx++; return r.idx <= len(r.items) }
func (r *sectionRows) Err() error                                   { return r.err }
func (r *sectionRows) CommandTag() pgconn.CommandTag                { return pgconn.CommandTag{} }
func (r *sectionRows) FieldDescriptions() []pgconn.FieldDescription { return nil }
func (r *sectionRows) Values() ([]interface{}, error)               { return nil, nil }
func (r *sectionRows) RawValues() [][]byte                          { return nil }
func (r *sectionRows) Scan(dest ...interface{}) error {
	if r.idx == 0 || r.idx > len(r.items) {
		return errors.New("scan called out of range")
	}
	item := r.items[r.idx-1]
	if len(dest) != 7 {
		return errors.New("unexpected scan args")
	}
	setUUID(dest[0], item.ID)
	setUUID(dest[1], item.NoteID)
	setUUID(dest[2], item.FieldID)
	setString(dest[3], item.Content)
	setString(dest[4], "")      // label
	setInt32(dest[5], int32(0)) // order
	setBool(dest[6], false)     // is_required
	return nil
}
func (r *sectionRows) Conn() *pgx.Conn { return nil }

func setInt32(ptr interface{}, v int32) {
	if dest, ok := ptr.(*int32); ok {
		*dest = v
	}
}
