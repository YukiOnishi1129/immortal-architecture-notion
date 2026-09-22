package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"

	domainerr "immortal-architecture-notion/backend/internal/domain/errors"
	"immortal-architecture-notion/backend/internal/domain/template"
	"immortal-architecture-notion/backend/internal/port"
	uc "immortal-architecture-notion/backend/internal/usecase"
	mockusecase "immortal-architecture-notion/backend/internal/usecase/mock"
)

func TestTemplateInteractor_Create(t *testing.T) {
	tests := []struct {
		name       string
		input      port.TemplateCreateInput
		created    *template.Template
		withFields *template.WithUsage
		createErr  error
		wantError  error
	}{
		{
			name: "[Success] create with fields",
			input: port.TemplateCreateInput{
				Name:    "Template",
				OwnerID: "owner-1",
				Fields: []template.Field{
					{ID: "f1", Label: "Title", Order: 1, IsRequired: true},
				},
			},
			created: &template.Template{ID: "tpl-1", Name: "Template", OwnerID: "owner-1"},
			withFields: &template.WithUsage{
				Template: template.Template{
					ID:      "tpl-1",
					Name:    "Template",
					OwnerID: "owner-1",
					Fields:  []template.Field{{ID: "f1", Label: "Title", Order: 1, IsRequired: true}},
				},
			},
		},
		{
			name: "[Fail] validation error",
			input: port.TemplateCreateInput{
				Name:    "",
				OwnerID: "owner-1",
				Fields:  []template.Field{{ID: "f1", Label: "Title", Order: 1}},
			},
			wantError: domainerr.ErrTemplateNameRequired,
		},
		{
			name: "[Fail] repo create error",
			input: port.TemplateCreateInput{
				Name:    "Template",
				OwnerID: "owner-1",
				Fields:  []template.Field{{ID: "f1", Label: "Title", Order: 1}},
			},
			createErr: errors.New("repo error"),
			wantError: errors.New("repo error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mockusecase.NewMockTemplateRepository(ctrl)
			tx := mockusecase.NewMockTxManager(ctrl)
			out := mockusecase.NewMockTemplateOutputPort(ctrl)

			// set expectations based on test case data
			if tt.created != nil || tt.createErr != nil {
				tx.EXPECT().WithinTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
					func(_ context.Context, fn func(context.Context) error) error {
						return fn(context.Background())
					},
				)
			}

			if tt.created != nil || tt.createErr != nil {
				repo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(tt.created, tt.createErr)
			}
			if tt.created != nil && tt.createErr == nil {
				repo.EXPECT().ReplaceFields(gomock.Any(), tt.created.ID, gomock.Any()).Return(nil)
				repo.EXPECT().Get(gomock.Any(), tt.created.ID).Return(tt.withFields, nil)
				out.EXPECT().PresentTemplate(gomock.Any(), tt.withFields).Return(nil)
			}

			interactor := uc.NewTemplateInteractor(repo, tx, out)
			err := interactor.Create(context.Background(), tt.input)

			if tt.wantError == nil && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.wantError != nil {
				if tt.wantError.Error() != err.Error() {
					t.Fatalf("want %v, got %v", tt.wantError, err)
				}
			}
		})
	}
}

func TestTemplateInteractor_List(t *testing.T) {
	tests := []struct {
		name      string
		filters   template.Filters
		result    []template.WithUsage
		repoErr   error
		wantError error
	}{
		{
			name:    "[Success] list templates",
			filters: template.Filters{OwnerID: strPtr("owner-1")},
			result: []template.WithUsage{
				{Template: template.Template{ID: "tpl-1", Name: "tpl"}},
			},
		},
		{
			name:      "[Fail] repo error",
			filters:   template.Filters{},
			repoErr:   errors.New("repo error"),
			wantError: errors.New("repo error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mockusecase.NewMockTemplateRepository(ctrl)
			tx := mockusecase.NewMockTxManager(ctrl)
			out := mockusecase.NewMockTemplateOutputPort(ctrl)

			repo.EXPECT().List(gomock.Any(), tt.filters).Return(tt.result, tt.repoErr)
			if tt.repoErr == nil {
				out.EXPECT().PresentTemplateList(gomock.Any(), tt.result).Return(nil)
			}

			interactor := uc.NewTemplateInteractor(repo, tx, out)
			err := interactor.List(context.Background(), tt.filters)

			if tt.wantError == nil && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.wantError != nil && tt.wantError.Error() != err.Error() {
				t.Fatalf("want %v, got %v", tt.wantError, err)
			}
		})
	}
}

func TestTemplateInteractor_Get(t *testing.T) {
	tests := []struct {
		name      string
		id        string
		result    *template.WithUsage
		repoErr   error
		wantError error
	}{
		{
			name:   "[Success] get template",
			id:     "tpl-1",
			result: &template.WithUsage{Template: template.Template{ID: "tpl-1", Name: "tpl"}},
		},
		{
			name:      "[Fail] not found",
			id:        "missing",
			repoErr:   domainerr.ErrNotFound,
			wantError: domainerr.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mockusecase.NewMockTemplateRepository(ctrl)
			tx := mockusecase.NewMockTxManager(ctrl)
			out := mockusecase.NewMockTemplateOutputPort(ctrl)

			repo.EXPECT().Get(gomock.Any(), tt.id).Return(tt.result, tt.repoErr)
			if tt.repoErr == nil {
				out.EXPECT().PresentTemplate(gomock.Any(), tt.result).Return(nil)
			}

			interactor := uc.NewTemplateInteractor(repo, tx, out)
			err := interactor.Get(context.Background(), tt.id)

			if tt.wantError == nil && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.wantError != nil && !errors.Is(err, tt.wantError) {
				t.Fatalf("want %v, got %v", tt.wantError, err)
			}
		})
	}
}

func TestTemplateInteractor_Update(t *testing.T) {
	tests := []struct {
		name        string
		input       port.TemplateUpdateInput
		current     *template.WithUsage
		getErr      error
		updateErr   error
		replaceErr  error
		wantError   error
		expectTxRun bool
	}{
		{
			name: "[Success] update name and fields",
			input: port.TemplateUpdateInput{
				ID:      "tpl-1",
				Name:    "updated",
				OwnerID: "owner-1",
				Fields: []template.Field{
					{ID: "f1", Label: "Title", Order: 1, IsRequired: true},
				},
			},
			current: &template.WithUsage{
				Template: template.Template{ID: "tpl-1", Name: "old", OwnerID: "owner-1"},
			},
			expectTxRun: true,
		},
		{
			name: "[Fail] owner required",
			input: port.TemplateUpdateInput{
				ID:      "tpl-1",
				Name:    "updated",
				OwnerID: "",
			},
			current:   &template.WithUsage{Template: template.Template{ID: "tpl-1", OwnerID: "owner-1"}},
			wantError: domainerr.ErrTemplateOwnerRequired,
		},
		{
			name: "[Fail] validate fields",
			input: port.TemplateUpdateInput{
				ID:      "tpl-1",
				Name:    "updated",
				OwnerID: "owner-1",
				Fields:  []template.Field{},
			},
			current:   &template.WithUsage{Template: template.Template{ID: "tpl-1", OwnerID: "owner-1"}},
			wantError: domainerr.ErrFieldRequired,
		},
		{
			name: "[Fail] repo get error",
			input: port.TemplateUpdateInput{
				ID:      "tpl-1",
				OwnerID: "owner-1",
			},
			getErr:    errors.New("get err"),
			wantError: errors.New("get err"),
		},
		{
			name: "[Fail] update error",
			input: port.TemplateUpdateInput{
				ID:      "tpl-1",
				Name:    "updated",
				OwnerID: "owner-1",
			},
			current:     &template.WithUsage{Template: template.Template{ID: "tpl-1", OwnerID: "owner-1"}},
			updateErr:   errors.New("update err"),
			wantError:   errors.New("update err"),
			expectTxRun: true,
		},
		{
			name: "[Fail] replace fields error",
			input: port.TemplateUpdateInput{
				ID:      "tpl-1",
				Name:    "updated",
				OwnerID: "owner-1",
				Fields: []template.Field{
					{ID: "f1", Label: "Title", Order: 1, IsRequired: true},
				},
			},
			current:     &template.WithUsage{Template: template.Template{ID: "tpl-1", OwnerID: "owner-1"}},
			updateErr:   nil,
			replaceErr:  errors.New("replace err"),
			wantError:   errors.New("replace err"),
			expectTxRun: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mockusecase.NewMockTemplateRepository(ctrl)
			tx := mockusecase.NewMockTxManager(ctrl)
			out := mockusecase.NewMockTemplateOutputPort(ctrl)

			repo.EXPECT().Get(gomock.Any(), tt.input.ID).Return(tt.current, tt.getErr)
			if tt.getErr == nil && tt.expectTxRun {
				tx.EXPECT().WithinTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
					func(_ context.Context, fn func(context.Context) error) error {
						return fn(context.Background())
					},
				)
			}
			if tt.getErr == nil && tt.expectTxRun {
				repo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(&tt.current.Template, tt.updateErr)
				if tt.updateErr == nil && tt.input.Fields != nil {
					repo.EXPECT().ReplaceFields(gomock.Any(), tt.input.ID, tt.input.Fields).Return(tt.replaceErr)
				}
			}
			if tt.getErr == nil && tt.expectTxRun && tt.updateErr == nil && tt.replaceErr == nil {
				repo.EXPECT().Get(gomock.Any(), tt.input.ID).Return(tt.current, nil)
				out.EXPECT().PresentTemplate(gomock.Any(), tt.current).Return(nil)
			}

			interactor := uc.NewTemplateInteractor(repo, tx, out)
			err := interactor.Update(context.Background(), tt.input)

			if tt.wantError == nil && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.wantError != nil && (err == nil || tt.wantError.Error() != err.Error()) {
				t.Fatalf("want %v, got %v", tt.wantError, err)
			}
		})
	}
}

func TestTemplateInteractor_Delete(t *testing.T) {
	tests := []struct {
		name      string
		id        string
		ownerID   string
		current   *template.WithUsage
		getErr    error
		deleteErr error
		wantError error
		expectDel bool
	}{
		{
			name:      "[Success] delete",
			id:        "tpl-1",
			ownerID:   "owner-1",
			current:   &template.WithUsage{Template: template.Template{ID: "tpl-1", OwnerID: "owner-1"}, IsUsed: false},
			expectDel: true,
		},
		{
			name:      "[Fail] get error",
			id:        "tpl-1",
			ownerID:   "owner-1",
			getErr:    errors.New("get err"),
			wantError: errors.New("get err"),
		},
		{
			name:      "[Fail] owner required",
			id:        "tpl-1",
			ownerID:   "",
			current:   &template.WithUsage{Template: template.Template{ID: "tpl-1", OwnerID: "owner-1"}},
			wantError: domainerr.ErrTemplateOwnerRequired,
		},
		{
			name:      "[Fail] in use",
			id:        "tpl-1",
			ownerID:   "owner-1",
			current:   &template.WithUsage{Template: template.Template{ID: "tpl-1", OwnerID: "owner-1"}, IsUsed: true},
			wantError: domainerr.ErrTemplateInUse,
		},
		{
			name:      "[Fail] delete error",
			id:        "tpl-1",
			ownerID:   "owner-1",
			current:   &template.WithUsage{Template: template.Template{ID: "tpl-1", OwnerID: "owner-1"}},
			deleteErr: errors.New("delete err"),
			wantError: errors.New("delete err"),
			expectDel: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mockusecase.NewMockTemplateRepository(ctrl)
			tx := mockusecase.NewMockTxManager(ctrl)
			out := mockusecase.NewMockTemplateOutputPort(ctrl)

			repo.EXPECT().Get(gomock.Any(), tt.id).Return(tt.current, tt.getErr)
			if tt.getErr == nil && tt.expectDel {
				repo.EXPECT().Delete(gomock.Any(), tt.id).Return(tt.deleteErr)
			}
			if tt.getErr == nil && tt.wantError == nil && tt.deleteErr == nil {
				out.EXPECT().PresentTemplateDeleted(gomock.Any()).Return(nil)
			}

			interactor := uc.NewTemplateInteractor(repo, tx, out)
			err := interactor.Delete(context.Background(), tt.id, tt.ownerID)

			if tt.wantError == nil && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.wantError != nil && (err == nil || tt.wantError.Error() != err.Error()) {
				t.Fatalf("want %v, got %v", tt.wantError, err)
			}
		})
	}
}

// TestTemplateInteractor_Update_NotionParentPage pins down how the optional
// notionParentPageUrl field is applied, since the existing Update cases accept
// any template and never assert the value that is written.
func TestTemplateInteractor_Update_NotionParentPage(t *testing.T) {
	const storedID = "1429989fe8ac4effbc8f57f56486db54"

	strPtr := func(s string) *string { return &s }

	tests := []struct {
		name    string
		sentURL *string
		want    string
	}{
		{
			name:    "[Keep] field omitted keeps the stored page",
			sentURL: nil,
			want:    storedID,
		},
		{
			name:    "[Clear] empty string clears the link",
			sentURL: strPtr(""),
			want:    "",
		},
		{
			name:    "[Replace] a new url overwrites the stored page",
			sentURL: strPtr("https://www.notion.so/team/Daily-aaaaaaaabbbbccccddddeeeeeeeeeeee"),
			want:    "aaaaaaaabbbbccccddddeeeeeeeeeeee",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mockusecase.NewMockTemplateRepository(ctrl)
			tx := mockusecase.NewMockTxManager(ctrl)
			out := mockusecase.NewMockTemplateOutputPort(ctrl)

			current := &template.WithUsage{
				Template: template.Template{
					ID:                 "tpl-1",
					Name:               "old",
					OwnerID:            "owner-1",
					NotionParentPageID: storedID,
				},
			}

			repo.EXPECT().Get(gomock.Any(), "tpl-1").Return(current, nil).Times(2)
			tx.EXPECT().WithinTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
				func(ctx context.Context, fn func(context.Context) error) error { return fn(ctx) },
			)

			var written string
			repo.EXPECT().Update(gomock.Any(), gomock.Any()).DoAndReturn(
				func(_ context.Context, tpl template.Template) (*template.Template, error) {
					written = tpl.NotionParentPageID
					return &current.Template, nil
				},
			)
			out.EXPECT().PresentTemplate(gomock.Any(), current).Return(nil)

			interactor := uc.NewTemplateInteractor(repo, tx, out)
			err := interactor.Update(context.Background(), port.TemplateUpdateInput{
				ID:                  "tpl-1",
				Name:                "updated",
				OwnerID:             "owner-1",
				NotionParentPageURL: tt.sentURL,
			})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if written != tt.want {
				t.Fatalf("NotionParentPageID = %q, want %q", written, tt.want)
			}
		})
	}
}

// TestTemplateInteractor_Update_InvalidNotionURL checks that a malformed url is
// rejected before anything is written.
func TestTemplateInteractor_Update_InvalidNotionURL(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mockusecase.NewMockTemplateRepository(ctrl)
	tx := mockusecase.NewMockTxManager(ctrl)
	out := mockusecase.NewMockTemplateOutputPort(ctrl)

	current := &template.WithUsage{
		Template: template.Template{ID: "tpl-1", Name: "old", OwnerID: "owner-1"},
	}
	repo.EXPECT().Get(gomock.Any(), "tpl-1").Return(current, nil)

	bad := "https://www.notion.so/not-a-page"
	interactor := uc.NewTemplateInteractor(repo, tx, out)
	err := interactor.Update(context.Background(), port.TemplateUpdateInput{
		ID:                  "tpl-1",
		Name:                "updated",
		OwnerID:             "owner-1",
		NotionParentPageURL: &bad,
	})
	if !errors.Is(err, domainerr.ErrInvalidNotionParentURL) {
		t.Fatalf("want %v, got %v", domainerr.ErrInvalidNotionParentURL, err)
	}
}
