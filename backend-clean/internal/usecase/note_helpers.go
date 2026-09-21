package usecase

import (
	domainerr "immortal-architecture-cqrs/backend/internal/domain/errors"
	"immortal-architecture-cqrs/backend/internal/domain/note"
	"immortal-architecture-cqrs/backend/internal/domain/template"
	"immortal-architecture-cqrs/backend/internal/port"
)

func buildSections(noteID string, inputs []port.SectionInput) ([]note.Section, error) {
	if len(inputs) == 0 {
		return nil, domainerr.ErrSectionsMissing
	}

	sections := make([]note.Section, 0, len(inputs))
	for _, s := range inputs {
		sections = append(sections, note.Section{
			FieldID: s.FieldID,
			NoteID:  noteID,
			Content: s.Content,
		})
	}
	return sections, nil
}

// buildSectionsForUpdate maps update inputs to sections using existing sections' field IDs.
func buildSectionsForUpdate(existing []note.SectionWithField, templateFields []template.Field, inputs []port.SectionUpdateInput, noteID string) ([]note.Section, error) {
	fieldBySection := make(map[string]string, len(existing))
	for _, s := range existing {
		fieldBySection[s.Section.ID] = s.Section.FieldID
	}
	sections := make([]note.Section, 0, len(inputs))
	for _, in := range inputs {
		fieldID, ok := fieldBySection[in.SectionID]
		if !ok {
			return nil, domainerr.ErrSectionsMissing
		}
		sections = append(sections, note.Section{
			ID:      in.SectionID,
			FieldID: fieldID,
			NoteID:  noteID,
			Content: in.Content,
		})
	}
	if err := note.ValidateSections(templateFields, sections); err != nil {
		return nil, err
	}
	return sections, nil
}
