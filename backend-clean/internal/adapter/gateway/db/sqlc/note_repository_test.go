package sqlc

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	mockdb "immortal-architecture-cqrs/backend/internal/adapter/gateway/db/sqlc/mock"
	"immortal-architecture-cqrs/backend/internal/adapter/gateway/db/sqlc/generated"
	domainerr "immortal-architecture-cqrs/backend/internal/domain/errors"
	"immortal-architecture-cqrs/backend/internal/domain/note"
)

func TestNoteRepository_UpdateStatus(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	baseRow := &generated.Note{
		ID:         pgtype.UUID{Bytes: [16]byte{1}, Valid: true},
		Title:      "t",
		TemplateID: pgtype.UUID{Bytes: [16]byte{2}, Valid: true},
		OwnerID:    pgtype.UUID{Bytes: [16]byte{3}, Valid: true},
		Status:     string(note.StatusDraft),
		CreatedAt:  pgtype.Timestamptz{Time: now, Valid: true},
		UpdatedAt:  pgtype.Timestamptz{Time: now, Valid: true},
	}
	tests := []struct {
		name    string
		id      string
		row     *generated.Note
		rowErr  error
		wantErr error
	}{
		{name: "[Success] update status", id: baseRow.ID.String(), row: baseRow},
		{name: "[Fail] invalid uuid", id: "bad-uuid", row: baseRow, wantErr: errors.New("invalid")},
		{name: "[Fail] not found", id: baseRow.ID.String(), rowErr: pgx.ErrNoRows, wantErr: domainerr.ErrNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := mockdb.NewNoteDBTX(tt.row, tt.rowErr, nil)
			repo := &NoteRepository{queries: generated.New(mock)}
			_, err := repo.UpdateStatus(context.Background(), tt.id, note.StatusPublish)
			if tt.wantErr == nil {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				return
			}
			if err == nil {
				t.Fatalf("expected error, got nil")
			}
			// no strict match on invalid uuid; just ensure error exists
			if tt.wantErr == domainerr.ErrNotFound && !errors.Is(err, domainerr.ErrNotFound) {
				t.Fatalf("want ErrNotFound, got %v", err)
			}
		})
	}
}

func TestNoteRepository_Delete(t *testing.T) {
	baseID := pgtype.UUID{Bytes: [16]byte{1}, Valid: true}
	tests := []struct {
		name    string
		id      string
		execErr error
		wantErr bool
	}{
		{name: "[Success] delete note", id: baseID.String()},
		{name: "[Fail] invalid uuid", id: "bad-uuid", wantErr: true},
		{name: "[Fail] exec error", id: baseID.String(), execErr: errors.New("db error"), wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := mockdb.NewNoteDBTX(nil, nil, tt.execErr)
			repo := &NoteRepository{queries: generated.New(mock)}
			err := repo.Delete(context.Background(), tt.id)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
			} else if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestNoteRepository_Create(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	row := &generated.Note{
		ID:         pgtype.UUID{Bytes: [16]byte{1}, Valid: true},
		Title:      "t",
		TemplateID: pgtype.UUID{Bytes: [16]byte{2}, Valid: true},
		OwnerID:    pgtype.UUID{Bytes: [16]byte{3}, Valid: true},
		Status:     string(note.StatusDraft),
		CreatedAt:  pgtype.Timestamptz{Time: now, Valid: true},
		UpdatedAt:  pgtype.Timestamptz{Time: now, Valid: true},
	}
	tests := []struct {
		name      string
		note      note.Note
		row       *generated.Note
		rowErr    error
		wantErr   bool
		wantTitle string
	}{
		{
			name:      "[Success] create note",
			note:      note.Note{Title: "t", TemplateID: row.TemplateID.String(), OwnerID: row.OwnerID.String(), Status: note.StatusDraft},
			row:       row,
			wantTitle: "t",
		},
		{
			name:    "[Fail] invalid template uuid",
			note:    note.Note{Title: "t", TemplateID: "bad-uuid", OwnerID: row.OwnerID.String(), Status: note.StatusDraft},
			wantErr: true,
		},
		{
			name:    "[Fail] query error",
			note:    note.Note{Title: "t", TemplateID: row.TemplateID.String(), OwnerID: row.OwnerID.String(), Status: note.StatusDraft},
			rowErr:  errors.New("db error"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := mockdb.NewNoteDBTX(tt.row, tt.rowErr, nil)
			repo := &NoteRepository{queries: generated.New(mock)}
			n, err := repo.Create(context.Background(), tt.note)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if n.Title != tt.wantTitle {
				t.Fatalf("title = %s, want %s", n.Title, tt.wantTitle)
			}
		})
	}
}

func TestNoteRepository_Get(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	baseRow := &generated.Note{
		ID:         pgtype.UUID{Bytes: [16]byte{1}, Valid: true},
		Title:      "t",
		TemplateID: pgtype.UUID{Bytes: [16]byte{2}, Valid: true},
		OwnerID:    pgtype.UUID{Bytes: [16]byte{3}, Valid: true},
		Status:     string(note.StatusDraft),
		CreatedAt:  pgtype.Timestamptz{Time: now, Valid: true},
		UpdatedAt:  pgtype.Timestamptz{Time: now, Valid: true},
	}
	sections := []*generated.Section{
		{ID: pgtype.UUID{Bytes: [16]byte{9}, Valid: true}, NoteID: baseRow.ID, FieldID: pgtype.UUID{Bytes: [16]byte{8}, Valid: true}, Content: "c"},
	}
	detail := &generated.GetNoteByIDRow{
		ID:             baseRow.ID,
		Title:          baseRow.Title,
		TemplateID:     baseRow.TemplateID,
		OwnerID:        baseRow.OwnerID,
		Status:         baseRow.Status,
		CreatedAt:      baseRow.CreatedAt,
		UpdatedAt:      baseRow.UpdatedAt,
		TemplateName:   "tpl",
		FirstName:      "Taro",
		LastName:       "Yamada",
		OwnerThumbnail: pgtype.Text{String: "thumb", Valid: true},
	}
	tests := []struct {
		name      string
		id        string
		row       *generated.Note
		getRow    *generated.GetNoteByIDRow
		rowErr    error
		queryErr  error
		wantErr   error
		wantTitle string
	}{
		{name: "[Success] get note", id: baseRow.ID.String(), row: baseRow, getRow: detail, wantTitle: baseRow.Title},
		{name: "[Fail] invalid uuid", id: "bad-uuid", row: baseRow, wantErr: errors.New("invalid")},
		{name: "[Fail] not found", id: baseRow.ID.String(), rowErr: pgx.ErrNoRows, wantErr: domainerr.ErrNotFound},
		{name: "[Fail] list sections error", id: baseRow.ID.String(), row: baseRow, getRow: detail, queryErr: errors.New("query"), wantErr: errors.New("query")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := mockdb.NewNoteDBTX(tt.row, tt.rowErr, nil).WithGetRow(tt.getRow).WithList(nil, sections, tt.queryErr)
			repo := &NoteRepository{queries: generated.New(mock)}
			got, err := repo.Get(context.Background(), tt.id)
			if tt.wantErr == nil {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if got.Note.Title != tt.wantTitle {
					t.Fatalf("title = %s, want %s", got.Note.Title, tt.wantTitle)
				}
				return
			}
			if err == nil {
				t.Fatalf("expected error, got nil")
			}
			if tt.wantErr == domainerr.ErrNotFound && !errors.Is(err, domainerr.ErrNotFound) {
				t.Fatalf("want ErrNotFound, got %v", err)
			}
		})
	}
}

func TestNoteRepository_Update(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	row := &generated.Note{
		ID:         pgtype.UUID{Bytes: [16]byte{1}, Valid: true},
		Title:      "t2",
		TemplateID: pgtype.UUID{Bytes: [16]byte{2}, Valid: true},
		OwnerID:    pgtype.UUID{Bytes: [16]byte{3}, Valid: true},
		Status:     string(note.StatusDraft),
		CreatedAt:  pgtype.Timestamptz{Time: now, Valid: true},
		UpdatedAt:  pgtype.Timestamptz{Time: now, Valid: true},
	}
	tests := []struct {
		name    string
		note    note.Note
		row     *generated.Note
		rowErr  error
		wantErr error
	}{
		{name: "[Success] update note", note: note.Note{ID: row.ID.String(), Title: "t2"}, row: row},
		{name: "[Fail] invalid uuid", note: note.Note{ID: "bad-uuid", Title: "t2"}, wantErr: errors.New("invalid")},
		{name: "[Fail] not found", note: note.Note{ID: row.ID.String(), Title: "t2"}, rowErr: pgx.ErrNoRows, wantErr: domainerr.ErrNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := mockdb.NewNoteDBTX(tt.row, tt.rowErr, nil)
			repo := &NoteRepository{queries: generated.New(mock)}
			got, err := repo.Update(context.Background(), tt.note)
			if tt.wantErr == nil {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if got.Title != tt.note.Title {
					t.Fatalf("title = %s, want %s", got.Title, tt.note.Title)
				}
				return
			}
			if err == nil {
				t.Fatalf("expected error, got nil")
			}
			if tt.wantErr == domainerr.ErrNotFound && !errors.Is(err, domainerr.ErrNotFound) {
				t.Fatalf("want ErrNotFound, got %v", err)
			}
		})
	}
}

func TestNoteRepository_List(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	noteRow := &generated.ListNotesRow{
		ID:             pgtype.UUID{Bytes: [16]byte{1}, Valid: true},
		Title:          "t",
		TemplateID:     pgtype.UUID{Bytes: [16]byte{2}, Valid: true},
		OwnerID:        pgtype.UUID{Bytes: [16]byte{3}, Valid: true},
		Status:         string(note.StatusDraft),
		CreatedAt:      pgtype.Timestamptz{Time: now, Valid: true},
		UpdatedAt:      pgtype.Timestamptz{Time: now, Valid: true},
		TemplateName:   "tpl",
		FirstName:      "Taro",
		LastName:       "Yamada",
		OwnerThumbnail: pgtype.Text{String: "thumb", Valid: true},
	}
	sections := []*generated.Section{
		{ID: pgtype.UUID{Bytes: [16]byte{9}, Valid: true}, NoteID: noteRow.ID, FieldID: pgtype.UUID{Bytes: [16]byte{8}, Valid: true}, Content: "c"},
	}
	tests := []struct {
		name     string
		notes    []*generated.ListNotesRow
		sections []*generated.Section
		queryErr error
		wantErr  bool
	}{
		{name: "[Success] list notes", notes: []*generated.ListNotesRow{noteRow}, sections: sections},
		{name: "[Fail] query error", queryErr: errors.New("db error"), wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := mockdb.NewNoteDBTX(nil, nil, nil).WithList(tt.notes, tt.sections, tt.queryErr)
			repo := &NoteRepository{queries: generated.New(mock)}
			_, err := repo.List(context.Background(), note.Filters{})
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
			} else if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestNoteRepository_ReplaceSections(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	noteID := pgtype.UUID{Bytes: [16]byte{1}, Valid: true}
	existingSection := &generated.Section{
		ID:      pgtype.UUID{Bytes: [16]byte{9}, Valid: true},
		NoteID:  noteID,
		FieldID: pgtype.UUID{Bytes: [16]byte{8}, Valid: true},
		Content: "old",
	}
	newSection := &generated.Section{
		ID:      pgtype.UUID{Bytes: [16]byte{7}, Valid: true},
		NoteID:  noteID,
		FieldID: pgtype.UUID{Bytes: [16]byte{6}, Valid: true},
		Content: "new",
	}

	tests := []struct {
		name     string
		noteID   string
		sections []note.Section
		secRow   *generated.Section
		rowErr   error
		wantErr  bool
	}{
		{
			name:   "[Success] update existing and create new",
			noteID: noteID.String(),
			sections: []note.Section{
				{ID: existingSection.ID.String(), Content: "updated"},
				{FieldID: newSection.FieldID.String(), Content: "created"},
			},
			secRow: existingSection,
		},
		{
			name:    "[Fail] invalid note uuid",
			noteID:  "bad-uuid",
			wantErr: true,
		},
		{
			name:   "[Fail] invalid section uuid",
			noteID: noteID.String(),
			sections: []note.Section{
				{ID: "bad-sec", Content: "updated"},
			},
			secRow:  existingSection,
			wantErr: true,
		},
		{
			name:   "[Fail] update error",
			noteID: noteID.String(),
			sections: []note.Section{
				{ID: existingSection.ID.String(), Content: "updated"},
			},
			secRow:  existingSection,
			rowErr:  errors.New("update err"),
			wantErr: true,
		},
		{
			name:   "[Fail] create error",
			noteID: noteID.String(),
			sections: []note.Section{
				{FieldID: newSection.FieldID.String(), Content: "created"},
			},
			rowErr:  errors.New("create err"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := mockdb.NewNoteDBTX(nil, tt.rowErr, nil).WithSectionRow(tt.secRow)
			repo := &NoteRepository{queries: generated.New(mock)}
			err := repo.ReplaceSections(context.Background(), tt.noteID, tt.sections)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
			} else if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			_ = now
		})
	}
}
