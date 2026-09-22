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

	// NotionParentPageID is where notes of this template are created in Notion.
	// Empty means the template is not linked to Notion.
	NotionParentPageID string
}

// Field represents a template field definition.
type Field struct {
	ID         string
	Label      string
	Order      int
	IsRequired bool
}
