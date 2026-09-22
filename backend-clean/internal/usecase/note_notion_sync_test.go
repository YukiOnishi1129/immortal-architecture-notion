package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"

	domainerr "immortal-architecture-notion/backend/internal/domain/errors"
	"immortal-architecture-notion/backend/internal/domain/note"
	"immortal-architecture-notion/backend/internal/domain/template"
	"immortal-architecture-notion/backend/internal/port"
	uc "immortal-architecture-notion/backend/internal/usecase"
	mockusecase "immortal-architecture-notion/backend/internal/usecase/mock"
)

const (
	testParentPageID = "1429989fe8ac4effbc8f57f56486db54"
	testNotionPageID = "aaaaaaaabbbbccccddddeeeeeeeeeeee"
	testNotionURL    = "https://www.notion.so/aaaaaaaabbbbccccddddeeeeeeeeeeee"
)

func strptr(s string) *string { return &s }

// strptrMatcher matches a *string by its value, since the pointer identity
// differs between the test and the code under test.
type strptrMatcher string

func (m strptrMatcher) Matches(x any) bool {
	got, ok := x.(*string)
	return ok && got != nil && *got == string(m)
}

func (m strptrMatcher) String() string { return "points to " + string(m) }

type notionDeps struct {
	notes      *mockusecase.MockNoteRepository
	readModels *mockusecase.MockNoteReadModelRepository
	templates  *mockusecase.MockTemplateRepository
	tx         *mockusecase.MockTxManager
	output     *mockusecase.MockNoteCommandOutputPort
	notion     *mockusecase.MockNotionClient
}

func newNotionDeps(ctrl *gomock.Controller) *notionDeps {
	return &notionDeps{
		notes:      mockusecase.NewMockNoteRepository(ctrl),
		readModels: mockusecase.NewMockNoteReadModelRepository(ctrl),
		templates:  mockusecase.NewMockTemplateRepository(ctrl),
		tx:         mockusecase.NewMockTxManager(ctrl),
		output:     mockusecase.NewMockNoteCommandOutputPort(ctrl),
		notion:     mockusecase.NewMockNotionClient(ctrl),
	}
}

func (d *notionDeps) interactor() port.NoteCommandInputPort {
	return uc.NewNoteCommandInteractor(d.notes, d.readModels, d.templates, d.tx, d.output, d.notion)
}

// expectTx makes the transaction run its body inline.
func (d *notionDeps) expectTx() {
	d.tx.EXPECT().WithinTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, fn func(context.Context) error) error { return fn(ctx) },
	)
}

func syncNoteWithMeta(status note.NoteStatus, pageID *string) *note.WithMeta {
	return &note.WithMeta{
		Note: note.Note{
			ID: "note-1", Title: "Title", TemplateID: "tpl-1", OwnerID: "owner-1",
			Status: status, NotionPageID: pageID,
		},
		Sections: []note.SectionWithField{
			{Section: note.Section{ID: "sec-1", NoteID: "note-1", FieldID: "f1", Content: "body"},
				FieldLabel: "Field", FieldOrder: 1, IsRequired: true},
		},
	}
}

func templateWithParent(parentID string) *template.WithUsage {
	return &template.WithUsage{
		Template: template.Template{ID: "tpl-1", OwnerID: "owner-1", NotionParentPageID: parentID},
	}
}

// QA-01: publishing a note for the first time creates a Notion page and
// stores its id.
func TestChangeStatus_FirstPublishCreatesPage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	d := newNotionDeps(ctrl)
	current := syncNoteWithMeta(note.StatusDraft, nil)

	d.notes.EXPECT().Get(gomock.Any(), "note-1").Return(current, nil)
	d.templates.EXPECT().Get(gomock.Any(), "tpl-1").Return(templateWithParent(testParentPageID), nil)
	d.notion.EXPECT().
		CreatePage(gomock.Any(), testParentPageID, "Title", []port.NotionSection{{Label: "Field", Content: "body"}}).
		Return(&port.NotionPage{PageID: testNotionPageID, URL: testNotionURL}, nil)

	d.expectTx()
	d.notes.EXPECT().UpdateStatus(gomock.Any(), "note-1", note.StatusPublish).Return(&current.Note, nil)
	d.notes.EXPECT().
		SaveNotionPage(gomock.Any(), "note-1", strptrMatcher(testNotionPageID), strptrMatcher(testNotionURL), gomock.Any()).
		Return(nil)
	d.notes.EXPECT().Get(gomock.Any(), "note-1").Return(current, nil).Times(2)
	d.readModels.EXPECT().Upsert(gomock.Any(), gomock.Any()).Return(nil)
	d.output.EXPECT().PresentNote(gomock.Any(), current).Return(nil)

	err := d.interactor().ChangeStatus(context.Background(), port.NoteStatusChangeInput{
		ID: "note-1", OwnerID: "owner-1", Status: note.StatusPublish,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// QA-04: publishing again restores the existing page instead of creating
// a second one, so the URL stays the same.
func TestChangeStatus_RepublishRestoresPage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	d := newNotionDeps(ctrl)
	current := syncNoteWithMeta(note.StatusDraft, strptr(testNotionPageID))

	d.notes.EXPECT().Get(gomock.Any(), "note-1").Return(current, nil)
	d.notion.EXPECT().Restore(gomock.Any(), testNotionPageID).
		Return(&port.NotionPage{PageID: testNotionPageID, URL: testNotionURL}, nil)

	d.expectTx()
	d.notes.EXPECT().UpdateStatus(gomock.Any(), "note-1", note.StatusPublish).Return(&current.Note, nil)
	d.notes.EXPECT().SaveNotionPage(gomock.Any(), "note-1", gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	d.notes.EXPECT().Get(gomock.Any(), "note-1").Return(current, nil).Times(2)
	d.readModels.EXPECT().Upsert(gomock.Any(), gomock.Any()).Return(nil)
	d.output.EXPECT().PresentNote(gomock.Any(), current).Return(nil)

	if err := d.interactor().ChangeStatus(context.Background(), port.NoteStatusChangeInput{
		ID: "note-1", OwnerID: "owner-1", Status: note.StatusPublish,
	}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// QA-03: unpublishing moves the page to the trash but keeps the page id,
// so the same page can be restored later.
func TestChangeStatus_UnpublishTrashesPage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	d := newNotionDeps(ctrl)
	current := syncNoteWithMeta(note.StatusPublish, strptr(testNotionPageID))

	d.notes.EXPECT().Get(gomock.Any(), "note-1").Return(current, nil)
	d.notion.EXPECT().Trash(gomock.Any(), testNotionPageID).Return(nil)

	d.expectTx()
	d.notes.EXPECT().UpdateStatus(gomock.Any(), "note-1", note.StatusDraft).Return(&current.Note, nil)
	// The page id must NOT be cleared, so SaveNotionPage is not called.
	d.notes.EXPECT().Get(gomock.Any(), "note-1").Return(current, nil).Times(2)
	d.readModels.EXPECT().Upsert(gomock.Any(), gomock.Any()).Return(nil)
	d.output.EXPECT().PresentNote(gomock.Any(), current).Return(nil)

	if err := d.interactor().ChangeStatus(context.Background(), port.NoteStatusChangeInput{
		ID: "note-1", OwnerID: "owner-1", Status: note.StatusDraft,
	}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// QA-20: the core of the design. If Notion fails, nothing is written and the
// note keeps its current status.
func TestChangeStatus_NotionFailureBlocksPublish(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	d := newNotionDeps(ctrl)
	current := syncNoteWithMeta(note.StatusDraft, nil)

	d.notes.EXPECT().Get(gomock.Any(), "note-1").Return(current, nil)
	d.templates.EXPECT().Get(gomock.Any(), "tpl-1").Return(templateWithParent(testParentPageID), nil)
	d.notion.EXPECT().CreatePage(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, errors.New("notion is down"))

	// No transaction, no status update, no read model write.
	err := d.interactor().ChangeStatus(context.Background(), port.NoteStatusChangeInput{
		ID: "note-1", OwnerID: "owner-1", Status: note.StatusPublish,
	})
	if !errors.Is(err, domainerr.ErrNotionSyncFailed) {
		t.Fatalf("want ErrNotionSyncFailed, got %v", err)
	}
}

// QA-30b: publishing is refused when the template has no parent page,
// even when the API is called directly.
func TestChangeStatus_ParentPageNotSet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	d := newNotionDeps(ctrl)
	current := syncNoteWithMeta(note.StatusDraft, nil)

	d.notes.EXPECT().Get(gomock.Any(), "note-1").Return(current, nil)
	d.templates.EXPECT().Get(gomock.Any(), "tpl-1").Return(templateWithParent(""), nil)

	err := d.interactor().ChangeStatus(context.Background(), port.NoteStatusChangeInput{
		ID: "note-1", OwnerID: "owner-1", Status: note.StatusPublish,
	})
	if !errors.Is(err, domainerr.ErrNotionParentNotSet) {
		t.Fatalf("want ErrNotionParentNotSet, got %v", err)
	}
}

// If the database write fails after the page was created, the page would be
// orphaned because its id was never stored, so it is trashed again.
func TestChangeStatus_DBFailureCleansUpOrphanPage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	d := newNotionDeps(ctrl)
	current := syncNoteWithMeta(note.StatusDraft, nil)
	dbErr := errors.New("db is down")

	d.notes.EXPECT().Get(gomock.Any(), "note-1").Return(current, nil)
	d.templates.EXPECT().Get(gomock.Any(), "tpl-1").Return(templateWithParent(testParentPageID), nil)
	d.notion.EXPECT().CreatePage(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(&port.NotionPage{PageID: testNotionPageID, URL: testNotionURL}, nil)

	d.expectTx()
	d.notes.EXPECT().UpdateStatus(gomock.Any(), "note-1", note.StatusPublish).Return(nil, dbErr)
	d.notion.EXPECT().Trash(gomock.Any(), testNotionPageID).Return(nil)

	if err := d.interactor().ChangeStatus(context.Background(), port.NoteStatusChangeInput{
		ID: "note-1", OwnerID: "owner-1", Status: note.StatusPublish,
	}); !errors.Is(err, dbErr) {
		t.Fatalf("want db error, got %v", err)
	}
}

// QA-22: with no API key the client is nil, and notes still work.
func TestChangeStatus_NotionDisabledStillPublishes(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	d := newNotionDeps(ctrl)
	current := syncNoteWithMeta(note.StatusDraft, nil)

	d.notes.EXPECT().Get(gomock.Any(), "note-1").Return(current, nil)
	d.expectTx()
	d.notes.EXPECT().UpdateStatus(gomock.Any(), "note-1", note.StatusPublish).Return(&current.Note, nil)
	d.notes.EXPECT().Get(gomock.Any(), "note-1").Return(current, nil).Times(2)
	d.readModels.EXPECT().Upsert(gomock.Any(), gomock.Any()).Return(nil)
	d.output.EXPECT().PresentNote(gomock.Any(), current).Return(nil)

	interactor := uc.NewNoteCommandInteractor(d.notes, d.readModels, d.templates, d.tx, d.output, nil)
	if err := interactor.ChangeStatus(context.Background(), port.NoteStatusChangeInput{
		ID: "note-1", OwnerID: "owner-1", Status: note.StatusPublish,
	}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// QA-02: editing a published note updates the existing page rather than
// creating a new one.
func TestUpdate_PublishedNoteUpdatesPage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	d := newNotionDeps(ctrl)
	current := syncNoteWithMeta(note.StatusPublish, strptr(testNotionPageID))

	d.notes.EXPECT().Get(gomock.Any(), "note-1").Return(current, nil)
	// The new content must be sent, not the stored one.
	d.notion.EXPECT().
		UpdatePage(gomock.Any(), testNotionPageID, "New title", []port.NotionSection{{Label: "Field", Content: "new body"}}).
		Return(&port.NotionPage{PageID: testNotionPageID, URL: testNotionURL}, nil)

	d.expectTx()
	d.notes.EXPECT().Update(gomock.Any(), gomock.Any()).Return(&current.Note, nil)
	d.templates.EXPECT().Get(gomock.Any(), "tpl-1").Return(&template.WithUsage{
		Template: template.Template{
			ID: "tpl-1", OwnerID: "owner-1",
			Fields: []template.Field{{ID: "f1", Label: "Field", Order: 1, IsRequired: true}},
		},
	}, nil)
	d.notes.EXPECT().ReplaceSections(gomock.Any(), "note-1", gomock.Any()).Return(nil)
	d.notes.EXPECT().SaveNotionPage(gomock.Any(), "note-1", gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	d.notes.EXPECT().Get(gomock.Any(), "note-1").Return(current, nil).Times(2)
	d.readModels.EXPECT().Upsert(gomock.Any(), gomock.Any()).Return(nil)
	d.output.EXPECT().PresentNote(gomock.Any(), current).Return(nil)

	if err := d.interactor().Update(context.Background(), port.NoteUpdateInput{
		ID: "note-1", OwnerID: "owner-1", Title: "New title",
		Sections: []port.SectionUpdateInput{{SectionID: "sec-1", Content: "new body"}},
	}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// Editing a draft must not touch Notion at all.
func TestUpdate_DraftDoesNotTouchNotion(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	d := newNotionDeps(ctrl)
	current := syncNoteWithMeta(note.StatusDraft, nil)

	d.notes.EXPECT().Get(gomock.Any(), "note-1").Return(current, nil)
	d.expectTx()
	d.notes.EXPECT().Update(gomock.Any(), gomock.Any()).Return(&current.Note, nil)
	d.notes.EXPECT().Get(gomock.Any(), "note-1").Return(current, nil).Times(2)
	d.readModels.EXPECT().Upsert(gomock.Any(), gomock.Any()).Return(nil)
	d.output.EXPECT().PresentNote(gomock.Any(), current).Return(nil)
	// No notion call is expected; gomock fails the test if one happens.

	if err := d.interactor().Update(context.Background(), port.NoteUpdateInput{
		ID: "note-1", OwnerID: "owner-1", Title: "New title",
	}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// QA-12: deleting a note leaves the Notion page alone.
func TestDelete_LeavesNotionPage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	d := newNotionDeps(ctrl)
	current := syncNoteWithMeta(note.StatusPublish, strptr(testNotionPageID))

	d.notes.EXPECT().Get(gomock.Any(), "note-1").Return(current, nil)
	d.expectTx()
	d.notes.EXPECT().Delete(gomock.Any(), "note-1").Return(nil)
	d.readModels.EXPECT().Delete(gomock.Any(), "note-1").Return(nil)
	d.output.EXPECT().PresentNoteDeleted(gomock.Any()).Return(nil)
	// No Trash call: the page stays in Notion.

	if err := d.interactor().Delete(context.Background(), "note-1", "owner-1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
