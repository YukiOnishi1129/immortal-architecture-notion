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
	"immortal-architecture-cqrs/backend/internal/domain/template"
)

func TestTemplateRepository_Create(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	row := &generated.Template{
		ID:        pgtype.UUID{Bytes: [16]byte{1}, Valid: true},
		Name:      "tpl",
		OwnerID:   pgtype.UUID{Bytes: [16]byte{2}, Valid: true},
		UpdatedAt: pgtype.Timestamptz{Time: now, Valid: true},
	}
	tests := []struct {
		name    string
		tpl     template.Template
		row     *generated.Template
		rowErr  error
		wantErr bool
	}{
		{name: "[Success] create template", tpl: template.Template{Name: "tpl", OwnerID: row.OwnerID.String()}, row: row},
		{name: "[Fail] invalid owner uuid", tpl: template.Template{Name: "tpl", OwnerID: "bad-uuid"}, wantErr: true},
		{name: "[Fail] query error", tpl: template.Template{Name: "tpl", OwnerID: row.OwnerID.String()}, rowErr: errors.New("db error"), wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := mockdb.NewTemplateDBTX(row, nil, tt.rowErr, nil)
			repo := &TemplateRepository{queries: generated.New(mock)}
			got, err := repo.Create(context.Background(), tt.tpl)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.Name != tt.tpl.Name {
				t.Fatalf("name = %s, want %s", got.Name, tt.tpl.Name)
			}
		})
	}
}

func TestTemplateRepository_Get(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	detail := &generated.GetTemplateByIDRow{
		ID:             pgtype.UUID{Bytes: [16]byte{1}, Valid: true},
		Name:           "tpl",
		OwnerID:        pgtype.UUID{Bytes: [16]byte{2}, Valid: true},
		UpdatedAt:      pgtype.Timestamptz{Time: now, Valid: true},
		OwnerFirstName: "Taro",
		OwnerLastName:  "Yamada",
		OwnerThumbnail: pgtype.Text{String: "thumb", Valid: true},
		IsUsed:         false,
	}
	tests := []struct {
		name    string
		id      string
		row     *generated.GetTemplateByIDRow
		rowErr  error
		wantErr error
	}{
		{name: "[Success] get template", id: detail.ID.String(), row: detail},
		{name: "[Fail] invalid uuid", id: "bad-uuid", row: detail, wantErr: errors.New("invalid")},
		{name: "[Fail] not found", id: detail.ID.String(), rowErr: pgx.ErrNoRows, wantErr: domainerr.ErrNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := mockdb.NewTemplateDBTX(nil, tt.row, tt.rowErr, nil)
			repo := &TemplateRepository{queries: generated.New(mock)}
			got, err := repo.Get(context.Background(), tt.id)
			if tt.wantErr == nil {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if got.Template.Name != tt.row.Name {
					t.Fatalf("name = %s, want %s", got.Template.Name, tt.row.Name)
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

func TestTemplateRepository_Update(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	tplRow := &generated.Template{
		ID:        pgtype.UUID{Bytes: [16]byte{1}, Valid: true},
		Name:      "tpl2",
		OwnerID:   pgtype.UUID{Bytes: [16]byte{2}, Valid: true},
		UpdatedAt: pgtype.Timestamptz{Time: now, Valid: true},
	}
	tests := []struct {
		name    string
		tpl     template.Template
		row     *generated.Template
		rowErr  error
		wantErr error
	}{
		{name: "[Success] update template", tpl: template.Template{ID: tplRow.ID.String(), Name: "tpl2"}, row: tplRow},
		{name: "[Fail] invalid uuid", tpl: template.Template{ID: "bad-uuid", Name: "tpl2"}, wantErr: errors.New("invalid")},
		{name: "[Fail] not found", tpl: template.Template{ID: tplRow.ID.String(), Name: "tpl2"}, rowErr: pgx.ErrNoRows, wantErr: domainerr.ErrNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := mockdb.NewTemplateDBTX(tt.row, nil, tt.rowErr, nil)
			repo := &TemplateRepository{queries: generated.New(mock)}
			got, err := repo.Update(context.Background(), tt.tpl)
			if tt.wantErr == nil {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if got.Name != tt.tpl.Name {
					t.Fatalf("name = %s, want %s", got.Name, tt.tpl.Name)
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

func TestTemplateRepository_Delete(t *testing.T) {
	baseID := pgtype.UUID{Bytes: [16]byte{1}, Valid: true}
	tests := []struct {
		name    string
		id      string
		execErr error
		wantErr bool
	}{
		{name: "[Success] delete template", id: baseID.String()},
		{name: "[Fail] invalid uuid", id: "bad-uuid", wantErr: true},
		{name: "[Fail] exec error", id: baseID.String(), execErr: errors.New("db error"), wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := mockdb.NewTemplateDBTX(nil, nil, nil, tt.execErr)
			repo := &TemplateRepository{queries: generated.New(mock)}
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

func TestTemplateRepository_List(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	// list returns emptyRows; we assert success and error path
	tests := []struct {
		name     string
		queryErr error
		wantErr  bool
	}{
		{name: "[Success] list templates"},
		{name: "[Fail] query error", queryErr: errors.New("db error"), wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := mockdb.NewTemplateDBTX(nil, nil, nil, nil)
			mock.QueryErr = tt.queryErr
			repo := &TemplateRepository{queries: generated.New(mock)}
			_, err := repo.List(context.Background(), template.Filters{})
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
			} else if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			_ = now // silence unused if removed
		})
	}
}

func TestTemplateRepository_ReplaceFields(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	tplID := pgtype.UUID{Bytes: [16]byte{1}, Valid: true}
	fieldRow := &generated.Field{
		ID:         pgtype.UUID{Bytes: [16]byte{2}, Valid: true},
		TemplateID: tplID,
		Label:      "lbl",
		Order:      1,
		IsRequired: true,
	}
	tests := []struct {
		name    string
		tplID   string
		field   template.Field
		rowErr  error
		execErr error
		wantErr bool
	}{
		{name: "[Success] replace fields", tplID: tplID.String(), field: template.Field{Label: "lbl", Order: 1, IsRequired: true}},
		{name: "[Fail] invalid tpl uuid", tplID: "bad-uuid", field: template.Field{Label: "lbl"}, wantErr: true},
		{name: "[Fail] delete error", tplID: tplID.String(), field: template.Field{Label: "lbl"}, execErr: errors.New("del error"), wantErr: true},
		{name: "[Fail] create error", tplID: tplID.String(), field: template.Field{Label: "lbl"}, rowErr: errors.New("create error"), wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := mockdb.NewTemplateDBTX(nil, nil, tt.rowErr, tt.execErr)
			mock.FieldRow = fieldRow
			repo := &TemplateRepository{queries: generated.New(mock)}
			err := repo.ReplaceFields(context.Background(), tt.tplID, []template.Field{tt.field})
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
