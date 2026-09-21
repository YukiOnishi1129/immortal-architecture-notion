// Package template holds template domain models.
package template

import "time"

// Template represents a note template aggregate.
type Template struct {
	ID        string
	Name      string
	OwnerID   string
	Fields    []Field
	UpdatedAt time.Time
}

// Field represents a template field definition.
type Field struct {
	ID         string
	Label      string
	Order      int
	IsRequired bool
}
