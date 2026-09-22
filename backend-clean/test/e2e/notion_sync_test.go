//go:build e2e

package e2e

import (
	"errors"
	"net/http"
	"testing"
)

// The parent page every linked template points at. The id inside the URL is
// what the application stores after normalizing it.
const (
	parentPageURL = "https://www.notion.so/workspace/Parent-1429989fe8ac4effbc8f57f56486db54"
	parentPageID  = "1429989fe8ac4effbc8f57f56486db54"
)

// linkedTemplate is wired to the parent page above; unlinkedTemplate is not,
// so publishing its notes must be refused.
var (
	linkedTemplate = templateInput{
		Name: "Daily Report", FieldLabel: "Summary", ParentPageURL: parentPageURL,
	}
	unlinkedTemplate = templateInput{
		Name: "Daily Report", FieldLabel: "Summary",
	}
)

var draftNote = noteInput{Title: "Sprint 12 report", Body: "shipped the importer"}

func (a *app) publish(noteID, ownerID string) response {
	a.t.Helper()
	return a.do(http.MethodPost, "/api/notes/"+noteID+"/publish?ownerId="+ownerID, nil)
}

func (a *app) unpublish(noteID, ownerID string) response {
	a.t.Helper()
	return a.do(http.MethodPost, "/api/notes/"+noteID+"/unpublish?ownerId="+ownerID, nil)
}

// TestE2E_Publish covers publishing a draft under the conditions that change
// the outcome: whether Notion is configured, reachable, and given a parent page.
func TestE2E_Publish(t *testing.T) {
	tests := []struct {
		name string

		// input
		template      templateInput
		note          noteInput
		notionEnabled bool
		notionDown    bool

		// expected output
		wantCode    int
		wantStatus  string
		wantLinked  bool // notionPageUrl is set
		wantPages   int  // pages in Notion
		wantNoCalls bool // Notion was never called
	}{
		{
			name:          "[Success] QA-01 publishing creates a notion page",
			template:      linkedTemplate,
			note:          draftNote,
			notionEnabled: true,
			wantCode:      http.StatusOK,
			wantStatus:    "Publish",
			wantLinked:    true,
			wantPages:     1,
		},
		{
			name:          "[Fail] QA-20 notion outage leaves the note a draft",
			template:      linkedTemplate,
			note:          draftNote,
			notionEnabled: true,
			notionDown:    true,
			wantCode:      http.StatusInternalServerError,
			wantStatus:    "Draft",
			wantLinked:    false,
			wantPages:     0,
		},
		{
			name:          "[Fail] QA-30b no parent page on the template",
			template:      unlinkedTemplate,
			note:          draftNote,
			notionEnabled: true,
			wantCode:      http.StatusBadRequest,
			wantStatus:    "Draft",
			wantLinked:    false,
			wantPages:     0,
			wantNoCalls:   true,
		},
		{
			name:          "[Success] QA-22 publishes normally while notion is not configured",
			template:      linkedTemplate,
			note:          draftNote,
			notionEnabled: false,
			wantCode:      http.StatusOK,
			wantStatus:    "Publish",
			wantLinked:    false,
			wantPages:     0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := newApp(t)
			if !tt.notionEnabled {
				a = newAppWithoutNotion(t)
			}
			owner := seedAccount(t)
			tpl := a.createTemplate(owner, tt.template)
			n := a.createNote(owner, tpl, tt.note)
			if tt.notionDown {
				a.notion.failWith(errors.New("notion is down"))
			}

			res := a.publish(n.ID, owner)
			if res.code != tt.wantCode {
				t.Fatalf("status = %d, want %d: %s", res.code, tt.wantCode, res.body)
			}

			got := a.getNote(n.ID, owner)
			if got.Status != tt.wantStatus {
				t.Fatalf("note status = %q, want %q", got.Status, tt.wantStatus)
			}
			if linked := got.NotionPageURL != nil; linked != tt.wantLinked {
				t.Fatalf("notionPageUrl set = %v, want %v", linked, tt.wantLinked)
			}
			if count := a.notion.pageCount(); count != tt.wantPages {
				t.Fatalf("notion pages = %d, want %d", count, tt.wantPages)
			}
			if tt.wantNoCalls {
				if calls := a.notion.callNames(); len(calls) != 0 {
					t.Fatalf("notion calls = %v, want none", calls)
				}
			}
		})
	}
}

// TestE2E_PublishedPageContent checks what actually lands on the Notion page,
// which the table above only counts.
//
//	in:  template "Daily Report" with field "Summary", parent page 1429989f...
//	     note     "Sprint 12 report" / "shipped the importer"
//	out: one page under that parent, titled "Sprint 12 report",
//	     holding one "Summary" section with the note body
func TestE2E_PublishedPageContent(t *testing.T) {
	a := newApp(t)
	owner := seedAccount(t)
	tpl := a.createTemplate(owner, linkedTemplate)
	n := a.createNote(owner, tpl, draftNote)

	a.publish(n.ID, owner).expect(t, http.StatusOK)

	got := a.getNote(n.ID, owner)
	page := a.notion.page(pageIDFromURL(*got.NotionPageURL))
	if page == nil {
		t.Fatalf("no notion page behind %q", *got.NotionPageURL)
	}

	if page.parentID != parentPageID {
		t.Fatalf("parent page = %q, want %q", page.parentID, parentPageID)
	}
	if page.title != draftNote.Title {
		t.Fatalf("title = %q, want %q", page.title, draftNote.Title)
	}
	if len(page.sections) != 1 {
		t.Fatalf("sections = %d, want 1", len(page.sections))
	}
	if page.sections[0].Label != linkedTemplate.FieldLabel {
		t.Fatalf("section label = %q, want %q", page.sections[0].Label, linkedTemplate.FieldLabel)
	}
	if page.sections[0].Content != draftNote.Body {
		t.Fatalf("section content = %q, want %q", page.sections[0].Content, draftNote.Body)
	}
}

// TestE2E_Edit covers editing a note. Only a published note is mirrored to
// Notion; a draft has no page yet.
func TestE2E_Edit(t *testing.T) {
	edited := noteInput{Title: "Sprint 12 retro", Body: "importer slipped"}

	tests := []struct {
		name string

		// input
		publishFirst bool
		edit         noteInput

		// expected output
		wantNotionCalls bool
		wantPageTitle   string // checked only when the page is expected to change
		wantPageBody    string
	}{
		{
			name:            "[Success] QA-02 editing a published note updates its page",
			publishFirst:    true,
			edit:            edited,
			wantNotionCalls: true,
			wantPageTitle:   edited.Title,
			wantPageBody:    edited.Body,
		},
		{
			name:            "[Success] editing a draft never reaches notion",
			publishFirst:    false,
			edit:            noteInput{Title: "Sprint 12 draft", Body: "still writing"},
			wantNotionCalls: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := newApp(t)
			owner := seedAccount(t)
			tpl := a.createTemplate(owner, linkedTemplate)
			n := a.createNote(owner, tpl, draftNote)

			sectionID := n.Sections[0].ID
			var pageID string
			if tt.publishFirst {
				a.publish(n.ID, owner).expect(t, http.StatusOK)
				published := a.getNote(n.ID, owner)
				pageID = pageIDFromURL(*published.NotionPageURL)
				sectionID = published.Sections[0].ID
				a.notion.resetCalls()
			}

			a.updateNote(n.ID, owner, tt.edit, sectionID).expect(t, http.StatusOK)

			calls := a.notion.callNames()
			if got := len(calls) > 0; got != tt.wantNotionCalls {
				t.Fatalf("notion called = %v (%v), want %v", got, calls, tt.wantNotionCalls)
			}
			if !tt.wantNotionCalls {
				return
			}

			// The edit must land on the existing page, not a new one.
			if count := a.notion.pageCount(); count != 1 {
				t.Fatalf("notion pages = %d, want 1 (no duplicate)", count)
			}
			page := a.notion.page(pageID)
			if page.title != tt.wantPageTitle {
				t.Fatalf("page title = %q, want %q", page.title, tt.wantPageTitle)
			}
			if page.sections[0].Content != tt.wantPageBody {
				t.Fatalf("page content = %q, want %q", page.sections[0].Content, tt.wantPageBody)
			}
		})
	}
}

// TestE2E_UnpublishAndRepublish is a sequence rather than a single action, so
// it stays outside the tables above.
//
//	QA-03 unpublishing trashes the page but keeps it
//	QA-04 republishing restores that same page, so the URL never changes
func TestE2E_UnpublishAndRepublish(t *testing.T) {
	a := newApp(t)
	owner := seedAccount(t)
	tpl := a.createTemplate(owner, linkedTemplate)
	n := a.createNote(owner, tpl, draftNote)
	a.publish(n.ID, owner).expect(t, http.StatusOK)

	firstURL := *a.getNote(n.ID, owner).NotionPageURL
	pageID := pageIDFromURL(firstURL)

	// --- unpublish: trashed, but the page and its id survive ---
	a.unpublish(n.ID, owner).expect(t, http.StatusOK)
	if page := a.notion.page(pageID); page == nil || !page.trashed {
		t.Fatalf("page was not trashed: %+v", page)
	}
	if status := a.getNote(n.ID, owner).Status; status != "Draft" {
		t.Fatalf("status = %q, want Draft", status)
	}

	// --- republish: the same page comes back ---
	a.publish(n.ID, owner).expect(t, http.StatusOK)
	if page := a.notion.page(pageID); page == nil || page.trashed {
		t.Fatalf("page was not restored: %+v", page)
	}
	if count := a.notion.pageCount(); count != 1 {
		t.Fatalf("notion pages = %d, want 1 (republish must reuse the page)", count)
	}
	if secondURL := *a.getNote(n.ID, owner).NotionPageURL; secondURL != firstURL {
		t.Fatalf("url changed: %q -> %q", firstURL, secondURL)
	}
}

// TestE2E_RetryAfterNotionRecovers completes QA-20: the user only has to press
// the same button again once Notion is back.
func TestE2E_RetryAfterNotionRecovers(t *testing.T) {
	a := newApp(t)
	owner := seedAccount(t)
	tpl := a.createTemplate(owner, linkedTemplate)
	n := a.createNote(owner, tpl, draftNote)

	a.notion.failWith(errors.New("notion is down"))
	a.publish(n.ID, owner)
	if status := a.getNote(n.ID, owner).Status; status != "Draft" {
		t.Fatalf("status = %q, want Draft while notion is down", status)
	}

	a.notion.recover()
	a.publish(n.ID, owner).expect(t, http.StatusOK)

	got := a.getNote(n.ID, owner)
	if got.Status != "Publish" {
		t.Fatalf("status = %q after retry, want Publish", got.Status)
	}
	if got.NotionPageURL == nil {
		t.Fatalf("notionPageUrl is nil after a successful retry")
	}
}

// TestE2E_DeleteKeepsNotionPage covers QA-12: the Notion page may hold edits
// made in Notion, so deleting the note must not remove it.
func TestE2E_DeleteKeepsNotionPage(t *testing.T) {
	a := newApp(t)
	owner := seedAccount(t)
	tpl := a.createTemplate(owner, linkedTemplate)
	n := a.createNote(owner, tpl, draftNote)
	a.publish(n.ID, owner).expect(t, http.StatusOK)
	a.notion.resetCalls()

	a.do(http.MethodDelete, "/api/notes/"+n.ID+"?ownerId="+owner, nil).expect(t, http.StatusOK)

	if count := a.notion.pageCount(); count != 1 {
		t.Fatalf("notion pages = %d, want the page to remain", count)
	}
	if calls := a.notion.callNames(); len(calls) != 0 {
		t.Fatalf("notion calls = %v, want none on delete", calls)
	}
}

func (a *app) syncToNotion(noteID, ownerID string) response {
	a.t.Helper()
	return a.do(http.MethodPost, "/api/notes/"+noteID+"/notion-sync?ownerId="+ownerID, nil)
}

// QA-23 notes published before the integration existed have no Notion page,
// and no status change will reach them. The manual sync creates one without
// taking the note out of the published list.
func TestE2E_ManualSyncLinksAlreadyPublishedNote(t *testing.T) {
	// Publish with the integration off, which is the state those notes are in.
	a := newAppWithoutNotion(t)
	owner := seedAccount(t)
	tpl := a.createTemplate(owner, linkedTemplate)
	n := a.createNote(owner, tpl, draftNote)
	a.publish(n.ID, owner).expect(t, http.StatusOK)

	published := a.getNote(n.ID, owner)
	if published.Status != "Publish" || published.NotionPageURL != nil {
		t.Fatalf("setup failed: status=%q url=%v", published.Status, published.NotionPageURL)
	}

	// The integration is now configured.
	a.withNotion()
	a.syncToNotion(n.ID, owner).expect(t, http.StatusOK)

	got := a.getNote(n.ID, owner)
	if got.Status != "Publish" {
		t.Fatalf("status = %q, want it to stay Publish", got.Status)
	}
	if got.NotionPageURL == nil {
		t.Fatalf("notionPageUrl is nil, want the note linked")
	}
	if count := a.notion.pageCount(); count != 1 {
		t.Fatalf("notion pages = %d, want 1", count)
	}
	page := a.notion.page(pageIDFromURL(*got.NotionPageURL))
	if page.title != draftNote.Title {
		t.Fatalf("notion title = %q, want %q", page.title, draftNote.Title)
	}
}

// The manual sync is only for published notes that have no page yet.
func TestE2E_ManualSyncRejected(t *testing.T) {
	tests := []struct {
		name         string
		publishFirst bool
		syncTwice    bool
	}{
		{name: "[Fail] a draft has nothing to show yet"},
		{name: "[Fail] already linked", publishFirst: true, syncTwice: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := newApp(t)
			owner := seedAccount(t)
			tpl := a.createTemplate(owner, linkedTemplate)
			n := a.createNote(owner, tpl, draftNote)
			if tt.publishFirst {
				a.publish(n.ID, owner).expect(t, http.StatusOK)
			}

			res := a.syncToNotion(n.ID, owner)
			if res.code != http.StatusBadRequest {
				t.Fatalf("status = %d, want 400: %s", res.code, res.body)
			}
			if tt.syncTwice && a.notion.pageCount() != 1 {
				t.Fatalf("notion pages = %d, want the original one only", a.notion.pageCount())
			}
		})
	}
}

// pageIDFromURL extracts the id from a page URL the fake returns.
func pageIDFromURL(url string) string {
	const prefix = "https://www.notion.so/"
	if len(url) <= len(prefix) {
		return ""
	}
	return url[len(prefix):]
}

// An update the application rejects must not reach Notion. The validation
// happens after the page would be written, so an invalid request could
// otherwise leave the page holding content the note never stored.
func TestE2E_RejectedEditDoesNotTouchNotion(t *testing.T) {
	a := newApp(t)
	owner := seedAccount(t)
	tpl := a.createTemplate(owner, linkedTemplate)
	n := a.createNote(owner, tpl, draftNote)
	a.publish(n.ID, owner).expect(t, http.StatusOK)

	published := a.getNote(n.ID, owner)
	pageID := pageIDFromURL(*published.NotionPageURL)
	a.notion.resetCalls()

	// A section id that does not belong to this note.
	res := a.updateNote(n.ID, owner, noteInput{Title: "Broken edit", Body: "nope"},
		"11111111-1111-1111-1111-111111111111")
	if res.code == http.StatusOK {
		t.Fatalf("the update was accepted, want a rejection: %s", res.body)
	}

	// The stored note must be untouched...
	got := a.getNote(n.ID, owner)
	if got.Title != draftNote.Title {
		t.Fatalf("title = %q, want %q (the edit was rejected)", got.Title, draftNote.Title)
	}

	// ...and so must the Notion page.
	page := a.notion.page(pageID)
	if page.title != draftNote.Title {
		t.Fatalf("notion title = %q, want %q: a rejected edit reached Notion",
			page.title, draftNote.Title)
	}
	if calls := a.notion.callNames(); len(calls) != 0 {
		t.Fatalf("notion calls = %v, want none for a rejected edit", calls)
	}
}
