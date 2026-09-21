package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"

	sqldb "immortal-architecture-bad-api/backend/internal/db/sqlc"
	openapi "immortal-architecture-bad-api/backend/internal/generated/openapi"
	"immortal-architecture-bad-api/backend/internal/service"
	mockservice "immortal-architecture-bad-api/backend/internal/service/mocks"
)

func TestTemplatesListTemplates(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		queryParam     string
		setup          func(*mockservice.MockTemplateServicer)
		wantStatus     int
		wantLen        int
		wantErrMessage string
	}{
		{
			name:       "success",
			queryParam: "foo",
			setup: func(m *mockservice.MockTemplateServicer) {
				m.EXPECT().ListTemplates(gomock.Any(), gomock.Any()).DoAndReturn(func(_ echo.Context, filters service.TemplateFilters) ([]*openapi.ModelsTemplateResponse, error) {
					if filters.Query == nil || *filters.Query != "foo" {
						return nil, errors.New("unexpected query")
					}
					return []*openapi.ModelsTemplateResponse{{Id: "tpl-1", Name: "Template"}}, nil
				})
			},
			wantStatus: http.StatusOK,
			wantLen:    1,
		},
		{
			name: "service error",
			setup: func(m *mockservice.MockTemplateServicer) {
				m.EXPECT().ListTemplates(gomock.Any(), gomock.Any()).Return(nil, errors.New("boom"))
			},
			wantStatus:     http.StatusInternalServerError,
			wantErrMessage: "failed to list templates",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			mockSvc := mockservice.NewMockTemplateServicer(mockCtrl)
			if tt.setup != nil {
				tt.setup(mockSvc)
			}

			e := echo.New()
			url := "/api/templates"
			if tt.queryParam != "" {
				url += "?q=" + tt.queryParam
			}
			req := httptest.NewRequest(http.MethodGet, url, nil)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)

			var params openapi.TemplatesListTemplatesParams
			if tt.queryParam != "" {
				params.Q = &tt.queryParam
			}

			ctrl := &Controller{templateService: mockSvc}
			if err := ctrl.TemplatesListTemplates(ctx, params); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if rec.Code != tt.wantStatus {
				t.Fatalf("expected status %d, got %d", tt.wantStatus, rec.Code)
			}

			if tt.wantStatus == http.StatusOK {
				var body []openapi.ModelsTemplateResponse
				if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
					t.Fatalf("failed to parse response: %v", err)
				}
				if len(body) != tt.wantLen {
					t.Fatalf("expected %d templates, got %d", tt.wantLen, len(body))
				}
			} else {
				var body openapi.ModelsErrorResponse
				if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
					t.Fatalf("failed to parse error response: %v", err)
				}
				if body.Message != tt.wantErrMessage {
					t.Fatalf("unexpected error message: %s", body.Message)
				}
			}
		})
	}
}

func TestTemplatesCreateTemplate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		body           string
		setup          func(*mockservice.MockTemplateServicer)
		wantStatus     int
		wantErrMessage string
	}{
		{
			name: "success",
			body: `{
                "name": "Template",
                "ownerId": "11111111-1111-1111-1111-111111111111",
                "fields": [{"label":"Title","order":1,"isRequired":true}]
            }`,
			setup: func(m *mockservice.MockTemplateServicer) {
				m.EXPECT().CreateTemplate(gomock.Any(), "11111111-1111-1111-1111-111111111111", "Template", gomock.Any()).DoAndReturn(
					func(_ echo.Context, ownerID string, name string, fields []sqldb.Field) (*openapi.ModelsTemplateResponse, error) {
						if ownerID == "" || name != "Template" || len(fields) != 1 {
							return nil, errors.New("unexpected inputs")
						}
						return &openapi.ModelsTemplateResponse{Id: "tpl-1", Name: name}, nil
					},
				)
			},
			wantStatus: http.StatusCreated,
		},
		{
			name:           "invalid payload",
			body:           `{`,
			wantStatus:     http.StatusBadRequest,
			wantErrMessage: "invalid payload",
		},
		{
			name: "service error",
			body: `{
                "name": "Template",
                "ownerId": "11111111-1111-1111-1111-111111111111",
                "fields": [{"label":"Title","order":1,"isRequired":true}]
            }`,
			setup: func(m *mockservice.MockTemplateServicer) {
				m.EXPECT().CreateTemplate(gomock.Any(), "11111111-1111-1111-1111-111111111111", "Template", gomock.Any()).Return(nil, service.ErrInvalidAccountID)
			},
			wantStatus:     http.StatusBadRequest,
			wantErrMessage: "invalid owner id",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			mockSvc := mockservice.NewMockTemplateServicer(mockCtrl)
			if tt.setup != nil {
				tt.setup(mockSvc)
			}

			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/api/templates", bytes.NewBufferString(tt.body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)

			ctrl := &Controller{templateService: mockSvc}
			if err := ctrl.TemplatesCreateTemplate(ctx); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if rec.Code != tt.wantStatus {
				t.Fatalf("expected status %d, got %d", tt.wantStatus, rec.Code)
			}
			if tt.wantStatus != http.StatusCreated {
				var body openapi.ModelsErrorResponse
				if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
					t.Fatalf("failed to parse error response: %v", err)
				}
				if body.Message != tt.wantErrMessage {
					t.Fatalf("unexpected error message: %s", body.Message)
				}
			}
		})
	}
}

func TestTemplatesDeleteTemplate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		setup          func(*mockservice.MockTemplateServicer)
		wantStatus     int
		wantErrMessage string
	}{
		{
			name: "success",
			setup: func(m *mockservice.MockTemplateServicer) {
				m.EXPECT().DeleteTemplate(gomock.Any(), "tpl-1").Return(nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "not found",
			setup: func(m *mockservice.MockTemplateServicer) {
				m.EXPECT().DeleteTemplate(gomock.Any(), "tpl-1").Return(service.ErrTemplateNotFound)
			},
			wantStatus:     http.StatusNotFound,
			wantErrMessage: "template not found",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			mockSvc := mockservice.NewMockTemplateServicer(mockCtrl)
			if tt.setup != nil {
				tt.setup(mockSvc)
			}

			e := echo.New()
			req := httptest.NewRequest(http.MethodDelete, "/api/templates/tpl-1", nil)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)
			ctx.SetParamNames("templateId")
			ctx.SetParamValues("tpl-1")

			ctrl := &Controller{templateService: mockSvc}
			if err := ctrl.TemplatesDeleteTemplate(ctx, "tpl-1"); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if rec.Code != tt.wantStatus {
				t.Fatalf("expected status %d, got %d", tt.wantStatus, rec.Code)
			}
			if tt.wantStatus != http.StatusOK {
				var body openapi.ModelsErrorResponse
				if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
					t.Fatalf("failed to parse error response: %v", err)
				}
				if body.Message != tt.wantErrMessage {
					t.Fatalf("unexpected error message: %s", body.Message)
				}
			}
		})
	}
}

func TestTemplatesGetTemplateById(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		setup          func(*mockservice.MockTemplateServicer)
		wantStatus     int
		wantErrMessage string
	}{
		{
			name: "success",
			setup: func(m *mockservice.MockTemplateServicer) {
				m.EXPECT().GetTemplate(gomock.Any(), "tpl-1").Return(&openapi.ModelsTemplateResponse{Id: "tpl-1", Name: "Template"}, nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "not found",
			setup: func(m *mockservice.MockTemplateServicer) {
				m.EXPECT().GetTemplate(gomock.Any(), "tpl-1").Return(nil, service.ErrTemplateNotFound)
			},
			wantStatus:     http.StatusNotFound,
			wantErrMessage: "template not found",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			mockSvc := mockservice.NewMockTemplateServicer(mockCtrl)
			if tt.setup != nil {
				tt.setup(mockSvc)
			}

			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/api/templates/tpl-1", nil)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)
			ctx.SetParamNames("templateId")
			ctx.SetParamValues("tpl-1")

			ctrl := &Controller{templateService: mockSvc}
			if err := ctrl.TemplatesGetTemplateById(ctx, "tpl-1"); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if rec.Code != tt.wantStatus {
				t.Fatalf("expected status %d, got %d", tt.wantStatus, rec.Code)
			}
			if tt.wantStatus != http.StatusOK {
				var body openapi.ModelsErrorResponse
				if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
					t.Fatalf("failed to parse error response: %v", err)
				}
				if body.Message != tt.wantErrMessage {
					t.Fatalf("unexpected error message: %s", body.Message)
				}
			}
		})
	}
}

func TestTemplatesUpdateTemplate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		body           string
		setup          func(*mockservice.MockTemplateServicer)
		wantStatus     int
		wantErrMessage string
	}{
		{
			name: "success",
			body: `{
                "name": "Template",
                "fields": [{"id":"11111111-1111-1111-1111-111111111111","label":"Title","order":1,"isRequired":true}]
            }`,
			setup: func(m *mockservice.MockTemplateServicer) {
				m.EXPECT().UpdateTemplate(gomock.Any(), "tpl-1", "Template", gomock.Any()).DoAndReturn(
					func(_ echo.Context, id, name string, fields []openapi.ModelsUpdateFieldRequest) (*openapi.ModelsTemplateResponse, error) {
						if id != "tpl-1" || name != "Template" || len(fields) != 1 {
							return nil, errors.New("unexpected inputs")
						}
						return &openapi.ModelsTemplateResponse{Id: id, Name: name}, nil
					},
				)
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "service conflict",
			body: `{
                "name": "Template",
                "fields": [{"id":"11111111-1111-1111-1111-111111111111","label":"Title","order":1,"isRequired":true}]
            }`,
			setup: func(m *mockservice.MockTemplateServicer) {
				m.EXPECT().UpdateTemplate(gomock.Any(), "tpl-1", "Template", gomock.Any()).Return(nil, service.ErrTemplateInUse)
			},
			wantStatus:     http.StatusConflict,
			wantErrMessage: "template is used by existing notes",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			mockSvc := mockservice.NewMockTemplateServicer(mockCtrl)
			if tt.setup != nil {
				tt.setup(mockSvc)
			}

			e := echo.New()
			req := httptest.NewRequest(http.MethodPut, "/api/templates/tpl-1", bytes.NewBufferString(tt.body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)
			ctx.SetParamNames("templateId")
			ctx.SetParamValues("tpl-1")

			ctrl := &Controller{templateService: mockSvc}
			if err := ctrl.TemplatesUpdateTemplate(ctx, "tpl-1"); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if rec.Code != tt.wantStatus {
				t.Fatalf("expected status %d, got %d", tt.wantStatus, rec.Code)
			}
			if tt.wantStatus != http.StatusOK {
				var body openapi.ModelsErrorResponse
				if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
					t.Fatalf("failed to parse error response: %v", err)
				}
				if body.Message != tt.wantErrMessage {
					t.Fatalf("unexpected error message: %s", body.Message)
				}
			}
		})
	}
}
