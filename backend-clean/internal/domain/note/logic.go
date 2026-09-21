package note

import (
	domainerr "immortal-architecture-cqrs/backend/internal/domain/errors"
	"immortal-architecture-cqrs/backend/internal/domain/template"
	"strings"
)

// Validate checks if status is valid.
func (s NoteStatus) Validate() error {
	if s != StatusDraft && s != StatusPublish {
		return domainerr.ErrInvalidStatus
	}
	return nil
}

// CanChangeStatus validates status transition.
func CanChangeStatus(from, to NoteStatus) error {
	if from == StatusDraft && to == StatusPublish {
		return nil
	}
	if from == StatusPublish && to == StatusDraft {
		return nil
	}
	if from == to {
		return nil
	}
	return domainerr.ErrInvalidStatusChange
}

// ValidateSections checks that sections match template fields and required fields are filled.
func ValidateSections(tplFields []template.Field, sections []Section) error {
	if len(sections) == 0 {
		return domainerr.ErrSectionsMissing
	}
	lookup := make(map[string]template.Field)
	for _, f := range tplFields {
		lookup[f.ID] = f
	}
	seen := make(map[string]bool)
	for _, s := range sections {
		f, ok := lookup[s.FieldID]
		if !ok {
			return domainerr.ErrSectionsMissing
		}
		if seen[s.FieldID] {
			return domainerr.ErrSectionsMissing
		}
		seen[s.FieldID] = true
		if f.IsRequired && s.Content == "" {
			return domainerr.ErrRequiredFieldEmpty
		}
	}
	// ensure all template fields are covered
	if len(seen) != len(lookup) {
		return domainerr.ErrSectionsMissing
	}
	return nil
}

// ValidateNoteForCreate validates a note creation attempt against template and required fields.
func ValidateNoteForCreate(title string, tpl template.Template, sections []Section) error {
	if strings.TrimSpace(title) == "" {
		return domainerr.ErrTitleRequired
	}
	if tpl.ID == "" {
		return domainerr.ErrTemplateOwnerRequired // template missing indicator
	}
	if tpl.OwnerID == "" {
		return domainerr.ErrOwnerRequired
	}
	if err := template.ValidateTemplate(tpl); err != nil {
		return err
	}
	if err := ValidateSections(tpl.Fields, sections); err != nil {
		return err
	}
	return nil
}

// ValidateNoteOwnership ensures only owner can mutate a note.
func ValidateNoteOwnership(noteOwnerID, actorID string) error {
	if strings.TrimSpace(noteOwnerID) == "" || strings.TrimSpace(actorID) == "" {
		return domainerr.ErrOwnerRequired
	}
	if noteOwnerID != actorID {
		return domainerr.ErrUnauthorized
	}
	return nil
}
