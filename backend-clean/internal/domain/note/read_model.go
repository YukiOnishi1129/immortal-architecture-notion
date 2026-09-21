package note

import "time"

// ReadModel is the denormalized read model for notes.
// It contains all data needed to display a note without JOINs.
type ReadModel struct {
	ID             string
	Title          string
	Status         NoteStatus
	TemplateID     string
	TemplateName   string
	OwnerID        string
	OwnerFirstName string
	OwnerLastName  string
	OwnerThumbnail *string
	Sections       []SectionReadModel
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// SectionReadModel is the denormalized read model for sections.
type SectionReadModel struct {
	ID         string `json:"id"`
	FieldID    string `json:"field_id"`
	FieldLabel string `json:"field_label"`
	FieldOrder int    `json:"field_order"`
	IsRequired bool   `json:"is_required"`
	Content    string `json:"content"`
}
