package template

import (
	"strings"

	domainerr "immortal-architecture-notion/backend/internal/domain/errors"
)

// NormalizeAndValidate sets missing order and validates fields.
func NormalizeAndValidate(fields []Field) ([]Field, error) {
	for i := range fields {
		if fields[i].Order == 0 {
			fields[i].Order = i + 1
		}
	}
	if err := validateFields(fields); err != nil {
		return nil, err
	}
	return fields, nil
}

func validateFields(fields []Field) error {
	seen := make(map[int]bool)
	if len(fields) == 0 {
		return domainerr.ErrFieldRequired
	}
	for _, f := range fields {
		if f.Label == "" {
			return domainerr.ErrFieldLabelRequired
		}
		order := f.Order
		if order <= 0 {
			return domainerr.ErrFieldOrderInvalid
		}
		if seen[order] {
			return domainerr.ErrFieldOrderInvalid
		}
		seen[order] = true
	}
	return nil
}

// ValidateTemplate ensures template has required attributes and valid fields.
func ValidateTemplate(t Template) error {
	if t.Name == "" {
		return domainerr.ErrTemplateNameRequired
	}
	if t.OwnerID == "" {
		return domainerr.ErrTemplateOwnerRequired
	}
	if _, err := NormalizeAndValidate(t.Fields); err != nil {
		return err
	}
	return nil
}

// CanDeleteTemplate returns error if template is in use.
func CanDeleteTemplate(isUsed bool) error {
	if isUsed {
		return domainerr.ErrTemplateInUse
	}
	return nil
}

// ValidateTemplateOwnership ensures only owner can mutate a template.
func ValidateTemplateOwnership(templateOwnerID, actorID string) error {
	if strings.TrimSpace(templateOwnerID) == "" || strings.TrimSpace(actorID) == "" {
		return domainerr.ErrTemplateOwnerRequired
	}
	if templateOwnerID != actorID {
		return domainerr.ErrUnauthorized
	}
	return nil
}

// ValidateFieldsChange rejects edits that would break notes already written
// from this template.
//
// A note stores its content per field id, so renaming a field changes what
// those notes show, and removing one destroys that content. Fields no note
// has written to yet are free to change, as are new fields.
//
// Only the fields are protected. The template name and its Notion parent page
// can always be edited, since no note content depends on them.
func ValidateFieldsChange(current, incoming []Field, usedFieldIDs []string) error {
	used := make(map[string]bool, len(usedFieldIDs))
	for _, id := range usedFieldIDs {
		used[id] = true
	}

	byID := make(map[string]Field, len(current))
	for _, f := range current {
		byID[f.ID] = f
	}

	seen := make(map[string]bool, len(incoming))
	for _, f := range incoming {
		if f.ID == "" {
			// A new field: existing notes are unaffected by it.
			continue
		}
		seen[f.ID] = true

		existing, ok := byID[f.ID]
		if !ok || !used[f.ID] {
			// Unknown ids are stored as new fields, and a field no note has
			// written to can still be changed freely.
			continue
		}
		if existing.Label != f.Label || existing.IsRequired != f.IsRequired {
			return domainerr.ErrTemplateInUse
		}
	}

	// Anything left out would be deleted.
	for id := range byID {
		if !seen[id] && used[id] {
			return domainerr.ErrTemplateInUse
		}
	}
	return nil
}
