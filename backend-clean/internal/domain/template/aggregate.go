// Package template holds template domain models.
package template

// ReplaceFields enforces updates through the aggregate root.
func (t *Template) ReplaceFields(fields []Field) error {
	fields, err := NormalizeAndValidate(fields)
	if err != nil {
		return err
	}
	t.Fields = fields
	return nil
}
