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

// Action is what needs to happen on the Notion side for a status change.
type Action string

// Actions derived from a note's status transition.
const (
	// ActionNone means Notion is left untouched.
	ActionNone Action = "none"
	// ActionCreate creates a new page under the template's parent page.
	ActionCreate Action = "create"
	// ActionRestore brings a previously trashed page back.
	ActionRestore Action = "restore"
	// ActionTrash moves the page to the Notion trash, keeping its id.
	ActionTrash Action = "trash"
	// ActionUpdate rewrites the title and content of an existing page.
	ActionUpdate Action = "update"
)
