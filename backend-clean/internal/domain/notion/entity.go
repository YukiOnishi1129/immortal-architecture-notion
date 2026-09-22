// Package notion holds domain models for the Notion integration.
package notion

import "time"

// Sync represents the link between a note and its Notion page.
// PageID is nil until the note is synced to Notion for the first time.
type Sync struct {
	NoteID   string
	PageID   *string
	PageURL  *string
	SyncedAt *time.Time
}

// ParentPage is the Notion page a note is created under.
// It is configured per template.
type ParentPage struct {
	PageID string
}
