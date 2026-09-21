package service

import (
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"

	sqldb "immortal-architecture-bad-api/backend/internal/db/sqlc"
	openapi "immortal-architecture-bad-api/backend/internal/generated/openapi"
)

func newEchoContext() echo.Context {
	e := echo.New()
	req := httptest.NewRequest("POST", "/", nil)
	rec := httptest.NewRecorder()
	return e.NewContext(req, rec)
}

func TestTemplateService_CreateTemplate_Validation(t *testing.T) {
	t.Parallel()

	// TemplateService.CreateTemplate が本来担保したい責務:
	// 1. ownerID の形式チェック（唯一ユニットテストで確認できる）
	// 2. fields を sqldb.Field に詰め替え、order を自動補完
	// 3. トランザクション管理（Begin/Commit/Rollback）と CreateField の挿入結果検証
	// 4. ownerID/fields に紐づく DB 制約エラーをハンドリング
	// しかし Service が pgxpool/sqlc と直結しているため、2 以降は DB か巨大モックを用意しない限り再現できない。

	tests := []struct {
		name        string
		ownerID     string
		fields      []sqldb.Field
		expectedErr error
	}{
		{
			name:        "invalid owner id",
			ownerID:     "not-a-uuid",
			expectedErr: ErrInvalidAccountID,
		},
		// {
		// 	name:    "success with auto order",
		// 	ownerID: "11111111-1111-1111-1111-111111111111",
		// 	fields: []sqldb.Field{
		// 		{Label: "Title", Order: 0, IsRequired: true},
		// 		{Label: "Body", Order: 0, IsRequired: false},
		// 	},
		// 	expectedErr: nil, // ← DB とトランザクションがないと確認できない
		// },
		// {
		// 	name:    "create field fails -> transaction rolls back",
		// 	ownerID: "11111111-1111-1111-1111-111111111111",
		// 	fields:  []sqldb.Field{{Label: "Title", Order: 1}},
		// 	expectedErr: errors.New("tx rollback should run"), // ← 実際は sqlc/pgx のエラーを検証したい
		// },
		// {
		// 	name:    "owner violates FK",
		// 	ownerID: "22222222-2222-2222-2222-222222222222",
		// 	fields:  []sqldb.Field{{Label: "Title", Order: 1}},
		// 	expectedErr: errors.New("foreign key violation"), // ← 実際は DB が返す固有エラーを想定
		// },
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx := newEchoContext()
			svc := &TemplateService{}

			_, err := svc.CreateTemplate(ctx, tt.ownerID, "Template", tt.fields)
			if tt.expectedErr == nil {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				return
			}
			if err == nil || !errors.Is(err, tt.expectedErr) {
				t.Fatalf("expected %v, got %v", tt.expectedErr, err)
			}
		})
	}
}

func TestTemplateService_UpdateTemplate_Validation(t *testing.T) {
	t.Parallel()

	// TemplateService.UpdateTemplate が担うべき責務:
	// 1. templateID の形式チェック。
	// 2. Name の更新。
	// 3. fields が nil/空/複数の場合のバリデーションと `syncTemplateFields` の Insert/Update/Delete。
	// 4. テンプレートの使用状況に応じた ErrTemplateInUse の返却。
	// 5. すべてトランザクション内で完結させる。
	// 現行コードでは 2-5 が pgx/sqlc に依存しており、ユニットテストでは templateID チェックしか走らせられない。

	tests := []struct {
		name        string
		templateID  string
		fields      []openapi.ModelsUpdateFieldRequest
		expectedErr error
	}{
		{
			name:        "invalid template id",
			templateID:  "bad-id",
			fields:      nil,
			expectedErr: ErrInvalidTemplateID,
		},
		// {
		// 	name:       "success update with fields",
		// 	templateID: "11111111-1111-1111-1111-111111111111",
		// 	fields: []openapi.ModelsUpdateFieldRequest{
		// 		{Id: nil, Label: "Title", Order: 1, IsRequired: true},
		// 	},
		// 	expectedErr: nil,
		// },
		// {
		// 	name: "template in use",
		// 	templateID: "11111111-1111-1111-1111-111111111111",
		// 	fields: []openapi.ModelsUpdateFieldRequest{
		// 		{Id: nil, Label: "Title", Order: 1, IsRequired: true},
		// 	},
		// 	expectedErr: ErrTemplateInUse,
		// },
		// {
		// 	name:       "empty fields should fail before sync",
		// 	templateID: "11111111-1111-1111-1111-111111111111",
		// 	fields:     []openapi.ModelsUpdateFieldRequest{},
		// 	expectedErr: errors.New("at least one field is required"),
		// },
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx := newEchoContext()
			svc := &TemplateService{}

			_, err := svc.UpdateTemplate(ctx, tt.templateID, "Template", tt.fields)
			if tt.expectedErr == nil {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				return
			}
			if err == nil || !errors.Is(err, tt.expectedErr) {
				t.Fatalf("expected %v, got %v", tt.expectedErr, err)
			}
		})
	}
}

func TestTemplateService_DeleteTemplate_Validation(t *testing.T) {
	t.Parallel()

	// TemplateService.DeleteTemplate が担うべき責務:
	// 1. templateID の形式チェック。
	// 2. テンプレートが存在するか確認し、なければ ErrTemplateNotFound。
	// 3. テンプレートを参照しているノートがあれば ErrTemplateInUse。
	// 4. トランザクション中に DeleteTemplate, CheckTemplateInUse の失敗を正しく扱う。
	// 2-4 は DB/トランザクション依存のため、現状の構造ではユニットテストから触れない。

	tests := []struct {
		name        string
		templateID  string
		expectedErr error
	}{
		{
			name:        "invalid template id",
			templateID:  "bad-id",
			expectedErr: ErrInvalidTemplateID,
		},
		// {
		// 	name:        "success deletion",
		// 	templateID:  "11111111-1111-1111-1111-111111111111",
		// 	expectedErr: nil,
		// },
		// {
		// 	name:        "template in use blocks deletion",
		// 	templateID:  "11111111-1111-1111-1111-111111111111",
		// 	expectedErr: ErrTemplateInUse,
		// },
		// {
		// 	name:        "template not found",
		// 	templateID:  "22222222-2222-2222-2222-222222222222",
		// 	expectedErr: ErrTemplateNotFound,
		// },
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx := newEchoContext()
			svc := &TemplateService{}

			err := svc.DeleteTemplate(ctx, tt.templateID)
			if tt.expectedErr == nil {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				return
			}
			if err == nil || !errors.Is(err, tt.expectedErr) {
				t.Fatalf("expected %v, got %v", tt.expectedErr, err)
			}
		})
	}
}
