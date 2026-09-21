package note

import "time"

// NoteStatus represents publication status.
//
//nolint:revive
type NoteStatus string

// Status constants.
const (
	StatusDraft   NoteStatus = "Draft"
	StatusPublish NoteStatus = "Publish"
)

// Note aggregate root.
type Note struct {
	ID         string
	Title      string
	TemplateID string
	OwnerID    string
	Status     NoteStatus
	Sections   []Section
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// Section represents note content for a field.
type Section struct {
	ID      string
	NoteID  string
	FieldID string
	Content string
}
