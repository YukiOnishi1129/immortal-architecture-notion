// Package port defines application ports (interfaces).
package port

import (
	"context"
	"time"

	"immortal-architecture-notion/backend/internal/domain/note"
)

// NoteRepository abstracts note persistence.
type NoteRepository interface {
	List(ctx context.Context, filters note.Filters) ([]note.WithMeta, error)
	Get(ctx context.Context, id string) (*note.WithMeta, error)
	Create(ctx context.Context, n note.Note) (*note.Note, error)
	Update(ctx context.Context, n note.Note) (*note.Note, error)
	UpdateStatus(ctx context.Context, id string, status note.NoteStatus) (*note.Note, error)
	Delete(ctx context.Context, id string) error
	ReplaceSections(ctx context.Context, noteID string, sections []note.Section) error

	// SaveNotionPage stores the Notion page a note is linked to.
	// Passing nil for pageID clears the link.
	SaveNotionPage(ctx context.Context, noteID string, pageID, pageURL *string, syncedAt *time.Time) error
}

// NoteCreateInput is input for creating notes.
type NoteCreateInput struct {
	Title      string
	TemplateID string
	OwnerID    string
	Sections   []SectionInput
}

// SectionInput is input for creating sections.
type SectionInput struct {
	FieldID string
	Content string
}

// NoteUpdateInput is input for updating notes.
type NoteUpdateInput struct {
	ID       string
	Title    string
	OwnerID  string
	Sections []SectionUpdateInput
}

// SectionUpdateInput is input for updating sections.
type SectionUpdateInput struct {
	SectionID string
	Content   string
}

// NoteStatusChangeInput is input for status changes.
type NoteStatusChangeInput struct {
	ID      string
	OwnerID string
	Status  note.NoteStatus
}
