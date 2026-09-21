package template

import (
	"strings"

	domainerr "immortal-architecture-cqrs/backend/internal/domain/errors"
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
