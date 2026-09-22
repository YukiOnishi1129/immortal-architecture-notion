package usecase

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	domainerr "immortal-architecture-notion/backend/internal/domain/errors"
	"immortal-architecture-notion/backend/internal/domain/note"
	"immortal-architecture-notion/backend/internal/domain/notion"
	"immortal-architecture-notion/backend/internal/port"
)

// syncResult carries what should be written to the notes table after a
// successful Notion call. A nil result means nothing needs to be stored.
type syncResult struct {
	pageID   *string
	pageURL  *string
	syncedAt *time.Time
}

// notionEnabled reports whether the integration is configured.
// When it is not, every note operation still works; syncing is skipped.
func (u *NoteCommandInteractor) notionEnabled() bool {
	return u.notion != nil
}

// syncStatusChange performs the Notion side of a status change.
//
// It runs OUTSIDE the transaction and BEFORE the database is touched, so a
// failure leaves the database untouched and the caller can simply return the
// error. Holding a transaction open across an external HTTP call would keep
// row locks for the duration of the request and can exhaust the pool.
func (u *NoteCommandInteractor) syncStatusChange(
	ctx context.Context, current *note.WithMeta, to note.NoteStatus,
) (*syncResult, error) {
	if !u.notionEnabled() {
		return nil, nil
	}

	sync := syncFromNote(current.Note)
	switch notion.ActionForStatusChange(to, sync) {
	case notion.ActionCreate:
		parentID, err := u.parentPageID(ctx, current.Note.TemplateID)
		if err != nil {
			return nil, err
		}
		page, err := u.notion.CreatePage(ctx, parentID, current.Note.Title, toNotionSections(current))
		if err != nil {
			return nil, wrapSyncErr(err)
		}
		return syncedNow(page), nil

	case notion.ActionRestore:
		page, err := u.notion.Restore(ctx, sync.PageIDOrEmpty())
		if err != nil {
			return nil, wrapSyncErr(err)
		}
		return syncedNow(page), nil

	case notion.ActionTrash:
		if err := u.notion.Trash(ctx, sync.PageIDOrEmpty()); err != nil {
			return nil, wrapSyncErr(err)
		}
		// The page id is kept so the same page can be restored later.
		return nil, nil

	case notion.ActionNone, notion.ActionUpdate:
		return nil, nil
	}
	return nil, nil
}

// syncEdit mirrors a content edit to Notion. Only published, already synced
// notes have a page to update.
func (u *NoteCommandInteractor) syncEdit(ctx context.Context, current *note.WithMeta) (*syncResult, error) {
	if !u.notionEnabled() {
		return nil, nil
	}

	sync := syncFromNote(current.Note)
	if notion.ActionForEdit(current.Note.Status, sync) != notion.ActionUpdate {
		return nil, nil
	}

	page, err := u.notion.UpdatePage(ctx, sync.PageIDOrEmpty(), current.Note.Title, toNotionSections(current))
	if err != nil {
		return nil, wrapSyncErr(err)
	}
	return syncedNow(page), nil
}

// cleanUpOrphanPage trashes a page that was created but never recorded.
//
// Without this the page id is lost, so the page could never be found again
// from the application. The cleanup itself can fail, in which case there is
// nothing left to do but leave a trace for an operator.
func (u *NoteCommandInteractor) cleanUpOrphanPage(ctx context.Context, res *syncResult, cause error) {
	if res == nil || res.pageID == nil || !u.notionEnabled() {
		return
	}
	if err := u.notion.Trash(ctx, *res.pageID); err != nil {
		slog.ErrorContext(ctx, "orphan notion page left behind",
			"page_id", *res.pageID, "cleanup_error", err, "cause", cause)
	}
}

// parentPageID reads the Notion parent page configured on the note's template.
func (u *NoteCommandInteractor) parentPageID(ctx context.Context, templateID string) (string, error) {
	tpl, err := u.templates.Get(ctx, templateID)
	if err != nil {
		return "", err
	}
	parentID := tpl.Template.NotionParentPageID
	// Guarded here as well as in the UI, because the API can be called directly.
	if err := notion.ValidateParentPageID(parentID); err != nil {
		return "", err
	}
	return parentID, nil
}

// editedNote returns what the note will look like after the update, so the
// content sent to Notion matches what is about to be written to the database.
// Sections left out of the request keep their current content.
func (u *NoteCommandInteractor) editedNote(
	current *note.WithMeta, input port.NoteUpdateInput,
) *note.WithMeta {
	if !u.notionEnabled() {
		return current
	}

	edited := *current
	edited.Note.Title = input.Title

	if input.Sections == nil {
		return &edited
	}

	updates := make(map[string]string, len(input.Sections))
	for _, s := range input.Sections {
		updates[s.SectionID] = s.Content
	}

	sections := make([]note.SectionWithField, len(current.Sections))
	copy(sections, current.Sections)
	for i := range sections {
		if content, ok := updates[sections[i].Section.ID]; ok {
			sections[i].Section.Content = content
		}
	}
	edited.Sections = sections
	return &edited
}

func syncFromNote(n note.Note) *notion.Sync {
	return &notion.Sync{
		NoteID:   n.ID,
		PageID:   n.NotionPageID,
		PageURL:  n.NotionPageURL,
		SyncedAt: n.NotionSyncedAt,
	}
}

func syncedNow(page *port.NotionPage) *syncResult {
	now := time.Now().UTC()
	id, url := page.PageID, page.URL
	return &syncResult{pageID: &id, pageURL: &url, syncedAt: &now}
}

// toNotionSections converts note content into the blocks sent to Notion,
// keeping the field order shown in the app.
func toNotionSections(wm *note.WithMeta) []port.NotionSection {
	sections := make([]port.NotionSection, 0, len(wm.Sections))
	for _, s := range wm.Sections {
		sections = append(sections, port.NotionSection{
			Label:   s.FieldLabel,
			Content: s.Section.Content,
		})
	}
	return sections
}

// wrapSyncErr keeps the underlying cause but lets callers match on a single
// sentinel, so the HTTP layer can turn every sync failure into one message.
func wrapSyncErr(err error) error {
	return fmt.Errorf("%w: %w", domainerr.ErrNotionSyncFailed, err)
}
