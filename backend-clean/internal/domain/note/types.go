// Package note holds note domain models.
package note

// Filters for listing notes.
type Filters struct {
	Status     *NoteStatus
	TemplateID *string
	OwnerID    *string
	Query      *string
}

// SectionWithField represents a section with template field metadata.
type SectionWithField struct {
	Section    Section
	FieldLabel string
	FieldOrder int
	IsRequired bool
}

// WithMeta represents a note with template metadata.
type WithMeta struct {
	Note           Note
	TemplateName   string
	OwnerFirstName string
	OwnerLastName  string
	OwnerThumbnail *string
	Sections       []SectionWithField
}
