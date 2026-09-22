package notion

import "immortal-architecture-notion/backend/internal/domain/note"

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

// ActionForStatusChange decides what to do on Notion when a note's status
// changes. The rules come from the design decision table:
//
//	Draft -> Publish (first time)  create
//	Draft -> Publish (again)       restore, because Notion has no hard delete
//	Publish -> Draft               trash, keeping the page id
//
// Restoring instead of re-creating keeps the page id and URL stable.
func ActionForStatusChange(to note.NoteStatus, sync *Sync) Action {
	if to == note.StatusPublish {
		if sync.IsFirstSync() {
			return ActionCreate
		}
		return ActionRestore
	}
	if sync.IsFirstSync() {
		// Never synced, so there is no page to trash.
		return ActionNone
	}
	return ActionTrash
}

// ActionForEdit decides what to do on Notion when a note's content is edited.
// Only published notes are mirrored; drafts are not on Notion yet.
func ActionForEdit(status note.NoteStatus, sync *Sync) Action {
	if status != note.StatusPublish || sync.IsFirstSync() {
		return ActionNone
	}
	return ActionUpdate
}
