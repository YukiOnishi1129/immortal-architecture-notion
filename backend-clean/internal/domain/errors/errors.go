// Package errors defines domain-level error values.
package errors

import "errors"

var (
	// ErrNotFound indicates resource not found.
	ErrNotFound = errors.New("not found")
	// ErrTemplateInUse indicates template is referenced by notes.
	ErrTemplateInUse = errors.New("template is used by notes")
	// ErrUnauthorized indicates authorization failure.
	ErrUnauthorized = errors.New("unauthorized")
	// ErrInvalidStatus indicates invalid status value.
	ErrInvalidStatus = errors.New("invalid status")
	// ErrInvalidStatusChange indicates invalid status transition.
	ErrInvalidStatusChange = errors.New("invalid status change")
	// ErrInvalidTemplateField indicates invalid template field definition.
	ErrInvalidTemplateField = errors.New("invalid template field")
	// ErrTemplateNameRequired indicates template name missing.
	ErrTemplateNameRequired = errors.New("template name is required")
	// ErrTemplateOwnerRequired indicates template owner missing.
	ErrTemplateOwnerRequired = errors.New("template owner is required")
	// ErrFieldRequired indicates at least one field is required.
	ErrFieldRequired = errors.New("template requires at least one field")
	// ErrFieldOrderInvalid indicates invalid field order.
	ErrFieldOrderInvalid = errors.New("field order must be greater than zero and unique")
	// ErrFieldLabelRequired indicates field label missing.
	ErrFieldLabelRequired = errors.New("field label is required")
	// ErrSectionsMissing indicates sections don't match template.
	ErrSectionsMissing = errors.New("sections do not match template fields")
	// ErrRequiredFieldEmpty indicates required field content missing.
	ErrRequiredFieldEmpty = errors.New("required field content is empty")
	// ErrProviderRequired indicates provider missing.
	ErrProviderRequired = errors.New("provider is required")
	// ErrProviderAccountRequired indicates provider account id missing.
	ErrProviderAccountRequired = errors.New("provider account id is required")
	// ErrTitleRequired indicates title missing.
	ErrTitleRequired = errors.New("title is required")
	// ErrOwnerRequired indicates owner missing.
	ErrOwnerRequired = errors.New("owner is required")
)
