package controller

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"

	ctrlmock "immortal-architecture-cqrs/backend/internal/adapter/http/controller/mock"
	openapi "immortal-architecture-cqrs/backend/internal/adapter/http/generated/openapi"
	"immortal-architecture-cqrs/backend/internal/adapter/http/presenter"
	domainerr "immortal-architecture-cqrs/backend/internal/domain/errors"
	"immortal-architecture-cqrs/backend/internal/port"
)

func TestTemplateController_Create(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		wantStatus int
	}{
		{
			name:       "[Success] create template",
			body:       `{"name":"Template","ownerId":"00000000-0000-0000-0000-000000000002","fields":[{"label":"Title","order":1,"isRequired":true}]}`,
			wantStatus: http.StatusOK,
		},
		{
			name:       "[Fail] bind error",
			body:       `not-json`,
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := presenter.NewTemplatePresenter()
			input := &ctrlmock.TemplateInputStub{}
			ctrl := NewTemplateController(
				func(repo port.TemplateRepository, tx port.TxManager, output port.TemplateOutputPort) port.TemplateInputPort {
					input.Output = output
					return input
				},
				func() *presenter.TemplatePresenter { return p },
				func() port.TemplateRepository { return nil },
				func() port.TxManager { return nil },
			)

			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/api/templates", bytes.NewBufferString(tt.body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			_ = ctrl.Create(c)
			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
		})
	}
}

func TestTemplateController_List(t *testing.T) {
	tests := []struct {
		name       string
		inErr      error
		wantStatus int
	}{
		{name: "[Success] list templates", wantStatus: http.StatusOK},
		{name: "[Fail] repo error", inErr: domainerr.ErrNotFound, wantStatus: http.StatusNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			p := presenter.NewTemplatePresenter()
			input := &ctrlmock.TemplateInputStub{Err: tt.inErr}
			ctrl := NewTemplateController(
				func(repo port.TemplateRepository, tx port.TxManager, output port.TemplateOutputPort) port.TemplateInputPort {
					input.Output = output
					return input
				},
				func() *presenter.TemplatePresenter { return p },
				func() port.TemplateRepository { return nil },
				func() port.TxManager { return nil },
			)
			req := httptest.NewRequest(http.MethodGet, "/api/templates", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			_ = ctrl.List(c, openapi.TemplatesListTemplatesParams{})
			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
		})
	}
}

func TestTemplateController_Get(t *testing.T) {
	tests := []struct {
		name       string
		inErr      error
		wantStatus int
	}{
		{name: "[Success] get template", wantStatus: http.StatusOK},
		{name: "[Fail] not found", inErr: domainerr.ErrNotFound, wantStatus: http.StatusNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			p := presenter.NewTemplatePresenter()
			input := &ctrlmock.TemplateInputStub{Err: tt.inErr}
			ctrl := NewTemplateController(
				func(repo port.TemplateRepository, tx port.TxManager, output port.TemplateOutputPort) port.TemplateInputPort {
					input.Output = output
					return input
				},
				func() *presenter.TemplatePresenter { return p },
				func() port.TemplateRepository { return nil },
				func() port.TxManager { return nil },
			)
			req := httptest.NewRequest(http.MethodGet, "/api/templates/t1", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			_ = ctrl.GetByID(c, "t1")
			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
		})
	}
}

func TestTemplateController_Update(t *testing.T) {
	e := echo.New()
	p := presenter.NewTemplatePresenter()
	input := &ctrlmock.TemplateInputStub{}
	ctrl := NewTemplateController(
		func(repo port.TemplateRepository, tx port.TxManager, output port.TemplateOutputPort) port.TemplateInputPort {
			input.Output = output
			return input
		},
		func() *presenter.TemplatePresenter { return p },
		func() port.TemplateRepository { return nil },
		func() port.TxManager { return nil },
	)

	tests := []struct {
		name       string
		body       string
		params     openapi.TemplatesUpdateTemplateParams
		inErr      error
		wantStatus int
	}{
		{
			name:       "[Success] update template",
			body:       `{"name":"updated","fields":[{"id":"f1","label":"Title","order":1,"isRequired":true}]}`,
			params:     openapi.TemplatesUpdateTemplateParams{OwnerId: "owner"},
			wantStatus: http.StatusOK,
		},
		{
			name:       "[Fail] missing owner",
			body:       `{"name":"updated","fields":[{"id":"f1","label":"Title","order":1,"isRequired":true}]}`,
			params:     openapi.TemplatesUpdateTemplateParams{OwnerId: ""},
			wantStatus: http.StatusForbidden,
		},
		{
			name:       "[Fail] usecase error",
			body:       `{"name":"updated","fields":[{"id":"f1","label":"Title","order":1,"isRequired":true}]}`,
			params:     openapi.TemplatesUpdateTemplateParams{OwnerId: "owner"},
			inErr:      domainerr.ErrNotFound,
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input.Err = tt.inErr
			req := httptest.NewRequest(http.MethodPut, "/api/templates/t1", bytes.NewBufferString(tt.body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			_ = ctrl.Update(c, "t1", tt.params)
			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
		})
	}
}

func TestTemplateController_Delete(t *testing.T) {
	tests := []struct {
		name       string
		ownerID    string
		inErr      error
		wantStatus int
	}{
		{name: "[Success] delete template", ownerID: "owner", wantStatus: http.StatusOK},
		{name: "[Fail] owner missing", ownerID: "", wantStatus: http.StatusForbidden},
		{name: "[Fail] not found", ownerID: "owner", inErr: domainerr.ErrNotFound, wantStatus: http.StatusNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			p := presenter.NewTemplatePresenter()
			input := &ctrlmock.TemplateInputStub{Err: tt.inErr}
			ctrl := NewTemplateController(
				func(repo port.TemplateRepository, tx port.TxManager, output port.TemplateOutputPort) port.TemplateInputPort {
					input.Output = output
					return input
				},
				func() *presenter.TemplatePresenter { return p },
				func() port.TemplateRepository { return nil },
				func() port.TxManager { return nil },
			)

			req := httptest.NewRequest(http.MethodDelete, "/api/templates/t1", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			_ = ctrl.Delete(c, "t1", openapi.TemplatesDeleteTemplateParams{OwnerId: tt.ownerID})
			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
		})
	}
}
