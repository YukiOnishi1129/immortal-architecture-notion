// Package note holds note domain models.
package note

import domainerr "immortal-architecture-cqrs/backend/internal/domain/errors"

// Aggregate helper methods to enforce child through parent.

// ReplaceSections updates sections via aggregate root.
func (n *Note) ReplaceSections(sections []Section) error {
	// Ensure sections belong to this note
	for i := range sections {
		if sections[i].NoteID == "" {
			sections[i].NoteID = n.ID
		}
		if sections[i].NoteID != n.ID {
			return domainerr.ErrUnauthorized
		}
	}
	n.Sections = sections
	return nil
}
