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

// noteWithMeta builds a note aggregate for test fixtures.
func noteWithMeta(id, ownerID string, status note.NoteStatus) *note.WithMeta {
	return &note.WithMeta{
		Note: note.Note{
			ID:         id,
			Title:      "Title",
			TemplateID: "tpl-1",
			OwnerID:    ownerID,
			Status:     status,
		},
		TemplateName: "Template",
		Sections: []note.SectionWithField{
			{
				Section:    note.Section{ID: "sec-1", NoteID: id, FieldID: "f1", Content: "content"},
				FieldLabel: "Field",
				FieldOrder: 1,
				IsRequired: true,
			},
		},
	}
}

// templateWithUsage builds a template aggregate for test fixtures.
func templateWithUsage() *template.WithUsage {
	return &template.WithUsage{
		Template: template.Template{
			ID:      "tpl-1",
			Name:    "Template",
			OwnerID: "owner-1",
			Fields: []template.Field{
				{ID: "f1", Label: "Field", Order: 1, IsRequired: true},
			},
		},
	}
}

// runInTx makes the TxManager mock execute the callback inline.
func runInTx(tx *mockusecase.MockTxManager) {
	tx.EXPECT().WithinTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
		func(_ context.Context, fn func(context.Context) error) error {
			return fn(context.Background())
		},
	)
}

func TestNoteCommandInteractor_Create(t *testing.T) {
	validInput := port.NoteCreateInput{
		Title:      "Title",
		TemplateID: "tpl-1",
		OwnerID:    "owner-1",
		Sections:   []port.SectionInput{{FieldID: "f1", Content: "content"}},
	}

	tests := []struct {
		name      string
		input     port.NoteCreateInput
		tpl       *template.WithUsage
		tplErr    error
		createErr error
		wantError error
		expectTx  bool
	}{
		{
			name:     "[Success] create a draft note",
			input:    validInput,
			tpl:      templateWithUsage(),
			expectTx: true,
		},
		{
			name: "[Fail] owner is required",
			input: port.NoteCreateInput{
				Title:      "Title",
				TemplateID: "tpl-1",
				Sections:   []port.SectionInput{{FieldID: "f1", Content: "content"}},
			},
			wantError: domainerr.ErrOwnerRequired,
		},
		{
			name:      "[Fail] template not found",
			input:     validInput,
			tplErr:    domainerr.ErrNotFound,
			wantError: domainerr.ErrNotFound,
		},
		{
			name: "[Fail] sections are missing",
			input: port.NoteCreateInput{
				Title:      "Title",
				TemplateID: "tpl-1",
				OwnerID:    "owner-1",
			},
			tpl:       templateWithUsage(),
			wantError: domainerr.ErrSectionsMissing,
		},
		{
			name: "[Fail] title is required",
			input: port.NoteCreateInput{
				Title:      "",
				TemplateID: "tpl-1",
				OwnerID:    "owner-1",
				Sections:   []port.SectionInput{{FieldID: "f1", Content: "content"}},
			},
			tpl:       templateWithUsage(),
			wantError: domainerr.ErrTitleRequired,
		},
		{
			name:      "[Fail] repository create error",
			input:     validInput,
			tpl:       templateWithUsage(),
			createErr: errors.New("db error"),
			wantError: errors.New("db error"),
			expectTx:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			notes := mockusecase.NewMockNoteRepository(ctrl)
			readModels := mockusecase.NewMockNoteReadModelRepository(ctrl)
			templates := mockusecase.NewMockTemplateRepository(ctrl)
			tx := mockusecase.NewMockTxManager(ctrl)
			output := mockusecase.NewMockNoteCommandOutputPort(ctrl)

			if tt.input.OwnerID != "" {
				templates.EXPECT().Get(gomock.Any(), tt.input.TemplateID).Return(tt.tpl, tt.tplErr)
			}

			if tt.expectTx {
				created := noteWithMeta("note-1", "owner-1", note.StatusDraft)
				runInTx(tx)
				notes.EXPECT().Create(gomock.Any(), gomock.Any()).
					Return(&created.Note, tt.createErr)

				if tt.createErr == nil {
					notes.EXPECT().ReplaceSections(gomock.Any(), "note-1", gomock.Any()).Return(nil)
					// once inside the transaction for the read model, once after it
					notes.EXPECT().Get(gomock.Any(), "note-1").Return(created, nil).Times(2)
					readModels.EXPECT().Upsert(gomock.Any(), gomock.Any()).Return(nil)
					output.EXPECT().PresentNote(gomock.Any(), created).Return(nil)
				}
			}

			interactor := uc.NewNoteCommandInteractor(notes, readModels, templates, tx, output)
			err := interactor.Create(context.Background(), tt.input)

			assertError(t, err, tt.wantError)
		})
	}
}

func TestNoteCommandInteractor_Update(t *testing.T) {
	tests := []struct {
		name      string
		input     port.NoteUpdateInput
		current   *note.WithMeta
		getErr    error
		updateErr error
		wantError error
		expectTx  bool
	}{
		{
			name: "[Success] update title only",
			input: port.NoteUpdateInput{
				ID: "note-1", Title: "New Title", OwnerID: "owner-1",
			},
			current:  noteWithMeta("note-1", "owner-1", note.StatusDraft),
			expectTx: true,
		},
		{
			name: "[Fail] not owner",
			input: port.NoteUpdateInput{
				ID: "note-1", Title: "New Title", OwnerID: "other",
			},
			current:   noteWithMeta("note-1", "owner-1", note.StatusDraft),
			wantError: domainerr.ErrUnauthorized,
		},
		{
			name: "[Fail] title is required",
			input: port.NoteUpdateInput{
				ID: "note-1", Title: "   ", OwnerID: "owner-1",
			},
			current:   noteWithMeta("note-1", "owner-1", note.StatusDraft),
			wantError: domainerr.ErrTitleRequired,
		},
		{
			name: "[Fail] note not found",
			input: port.NoteUpdateInput{
				ID: "note-1", Title: "New Title", OwnerID: "owner-1",
			},
			getErr:    domainerr.ErrNotFound,
			wantError: domainerr.ErrNotFound,
		},
		{
			name: "[Fail] repository update error",
			input: port.NoteUpdateInput{
				ID: "note-1", Title: "New Title", OwnerID: "owner-1",
			},
			current:   noteWithMeta("note-1", "owner-1", note.StatusDraft),
			updateErr: errors.New("db error"),
			wantError: errors.New("db error"),
			expectTx:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			notes := mockusecase.NewMockNoteRepository(ctrl)
			readModels := mockusecase.NewMockNoteReadModelRepository(ctrl)
			templates := mockusecase.NewMockTemplateRepository(ctrl)
			tx := mockusecase.NewMockTxManager(ctrl)
			output := mockusecase.NewMockNoteCommandOutputPort(ctrl)

			notes.EXPECT().Get(gomock.Any(), tt.input.ID).Return(tt.current, tt.getErr)

			if tt.expectTx {
				runInTx(tx)
				notes.EXPECT().Update(gomock.Any(), gomock.Any()).
					Return(&tt.current.Note, tt.updateErr)

				if tt.updateErr == nil {
					// once inside the transaction for the read model, once after it
					notes.EXPECT().Get(gomock.Any(), tt.input.ID).Return(tt.current, nil).Times(2)
					readModels.EXPECT().Upsert(gomock.Any(), gomock.Any()).Return(nil)
					output.EXPECT().PresentNote(gomock.Any(), tt.current).Return(nil)
				}
			}

			interactor := uc.NewNoteCommandInteractor(notes, readModels, templates, tx, output)
			err := interactor.Update(context.Background(), tt.input)

			assertError(t, err, tt.wantError)
		})
	}
}

func TestNoteCommandInteractor_ChangeStatus(t *testing.T) {
	tests := []struct {
		name       string
		input      port.NoteStatusChangeInput
		current    *note.WithMeta
		getErr     error
		updateErr  error
		upsertErr  error
		wantError  error
		expectCall bool // UpdateStatus まで到達するか
	}{
		{
			name: "[Success] draft to publish",
			input: port.NoteStatusChangeInput{
				ID: "note-1", OwnerID: "owner-1", Status: note.StatusPublish,
			},
			current:    noteWithMeta("note-1", "owner-1", note.StatusDraft),
			expectCall: true,
		},
		{
			name: "[Success] publish to draft",
			input: port.NoteStatusChangeInput{
				ID: "note-1", OwnerID: "owner-1", Status: note.StatusDraft,
			},
			current:    noteWithMeta("note-1", "owner-1", note.StatusPublish),
			expectCall: true,
		},
		{
			name: "[Fail] not owner",
			input: port.NoteStatusChangeInput{
				ID: "note-1", OwnerID: "other", Status: note.StatusPublish,
			},
			current:   noteWithMeta("note-1", "owner-1", note.StatusDraft),
			wantError: domainerr.ErrUnauthorized,
		},
		{
			name: "[Fail] invalid status value",
			input: port.NoteStatusChangeInput{
				ID: "note-1", OwnerID: "owner-1", Status: note.NoteStatus("Unknown"),
			},
			current:   noteWithMeta("note-1", "owner-1", note.StatusDraft),
			wantError: domainerr.ErrInvalidStatus,
		},
		{
			// 現状の仕様: CanChangeStatus は from == to を許容する（冪等）
			name: "[Success] publishing an already published note is allowed",
			input: port.NoteStatusChangeInput{
				ID: "note-1", OwnerID: "owner-1", Status: note.StatusPublish,
			},
			current:    noteWithMeta("note-1", "owner-1", note.StatusPublish),
			expectCall: true,
		},
		{
			name: "[Fail] note not found",
			input: port.NoteStatusChangeInput{
				ID: "note-1", OwnerID: "owner-1", Status: note.StatusPublish,
			},
			getErr:    domainerr.ErrNotFound,
			wantError: domainerr.ErrNotFound,
		},
		{
			name: "[Fail] update status error",
			input: port.NoteStatusChangeInput{
				ID: "note-1", OwnerID: "owner-1", Status: note.StatusPublish,
			},
			current:    noteWithMeta("note-1", "owner-1", note.StatusDraft),
			updateErr:  errors.New("db error"),
			wantError:  errors.New("db error"),
			expectCall: true,
		},
		{
			name: "[Fail] read model upsert error",
			input: port.NoteStatusChangeInput{
				ID: "note-1", OwnerID: "owner-1", Status: note.StatusPublish,
			},
			current:    noteWithMeta("note-1", "owner-1", note.StatusDraft),
			upsertErr:  errors.New("upsert error"),
			wantError:  errors.New("upsert error"),
			expectCall: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			notes := mockusecase.NewMockNoteRepository(ctrl)
			readModels := mockusecase.NewMockNoteReadModelRepository(ctrl)
			templates := mockusecase.NewMockTemplateRepository(ctrl)
			tx := mockusecase.NewMockTxManager(ctrl)
			output := mockusecase.NewMockNoteCommandOutputPort(ctrl)

			notes.EXPECT().Get(gomock.Any(), tt.input.ID).Return(tt.current, tt.getErr)

			if tt.expectCall {
				notes.EXPECT().UpdateStatus(gomock.Any(), tt.input.ID, tt.input.Status).
					Return(&tt.current.Note, tt.updateErr)

				if tt.updateErr == nil {
					// ChangeStatus reloads the note after updating.
					notes.EXPECT().Get(gomock.Any(), tt.input.ID).Return(tt.current, nil)
					readModels.EXPECT().Upsert(gomock.Any(), gomock.Any()).Return(tt.upsertErr)

					if tt.upsertErr == nil {
						output.EXPECT().PresentNote(gomock.Any(), tt.current).Return(nil)
					}
				}
			}

			interactor := uc.NewNoteCommandInteractor(notes, readModels, templates, tx, output)
			err := interactor.ChangeStatus(context.Background(), tt.input)

			assertError(t, err, tt.wantError)
		})
	}
}

func TestNoteCommandInteractor_Delete(t *testing.T) {
	tests := []struct {
		name      string
		noteID    string
		ownerID   string
		current   *note.WithMeta
		getErr    error
		deleteErr error
		rmErr     error
		wantError error
		expectTx  bool
	}{
		{
			name:     "[Success] delete own note",
			noteID:   "note-1",
			ownerID:  "owner-1",
			current:  noteWithMeta("note-1", "owner-1", note.StatusDraft),
			expectTx: true,
		},
		{
			name:      "[Fail] not owner",
			noteID:    "note-1",
			ownerID:   "other",
			current:   noteWithMeta("note-1", "owner-1", note.StatusDraft),
			wantError: domainerr.ErrUnauthorized,
		},
		{
			name:      "[Fail] note not found",
			noteID:    "note-1",
			ownerID:   "owner-1",
			getErr:    domainerr.ErrNotFound,
			wantError: domainerr.ErrNotFound,
		},
		{
			name:      "[Fail] delete error",
			noteID:    "note-1",
			ownerID:   "owner-1",
			current:   noteWithMeta("note-1", "owner-1", note.StatusDraft),
			deleteErr: errors.New("db error"),
			wantError: errors.New("db error"),
			expectTx:  true,
		},
		{
			name:      "[Fail] read model delete error",
			noteID:    "note-1",
			ownerID:   "owner-1",
			current:   noteWithMeta("note-1", "owner-1", note.StatusDraft),
			rmErr:     errors.New("rm error"),
			wantError: errors.New("rm error"),
			expectTx:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			notes := mockusecase.NewMockNoteRepository(ctrl)
			readModels := mockusecase.NewMockNoteReadModelRepository(ctrl)
			templates := mockusecase.NewMockTemplateRepository(ctrl)
			tx := mockusecase.NewMockTxManager(ctrl)
			output := mockusecase.NewMockNoteCommandOutputPort(ctrl)

			notes.EXPECT().Get(gomock.Any(), tt.noteID).Return(tt.current, tt.getErr)

			if tt.expectTx {
				runInTx(tx)
				notes.EXPECT().Delete(gomock.Any(), tt.noteID).Return(tt.deleteErr)

				if tt.deleteErr == nil {
					readModels.EXPECT().Delete(gomock.Any(), tt.noteID).Return(tt.rmErr)

					if tt.rmErr == nil {
						output.EXPECT().PresentNoteDeleted(gomock.Any()).Return(nil)
					}
				}
			}

			interactor := uc.NewNoteCommandInteractor(notes, readModels, templates, tx, output)
			err := interactor.Delete(context.Background(), tt.noteID, tt.ownerID)

			assertError(t, err, tt.wantError)
		})
	}
}

// assertError compares an actual error against the expected one.
// Sentinel errors are matched with errors.Is, others by message.
func assertError(t *testing.T, got, want error) {
	t.Helper()

	if want == nil {
		if got != nil {
			t.Fatalf("unexpected error: %v", got)
		}
		return
	}

	if got == nil {
		t.Fatalf("expected error %v, got nil", want)
	}

	if errors.Is(got, want) {
		return
	}
	if got.Error() != want.Error() {
		t.Fatalf("expected error %v, got %v", want, got)
	}
}
