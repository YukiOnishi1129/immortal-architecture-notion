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

func TestAccountController_CreateOrGet(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		createErr  error
		wantStatus int
		wantBody   string
	}{
		{
			name:       "[Success] create or get",
			body:       `{"email":"user@example.com","name":"Taro","provider":"google","providerAccountId":"pid"}`,
			wantStatus: http.StatusOK,
		},
		{
			name:       "[Fail] usecase error",
			body:       `{"email":"user@example.com","name":"Taro","provider":"google","providerAccountId":"pid"}`,
			createErr:  domainerr.ErrNotFound,
			wantStatus: http.StatusNotFound,
			wantBody:   domainerr.ErrNotFound.Error(),
		},
		{
			name:       "[Fail] bind error",
			body:       `not-json`,
			wantStatus: http.StatusBadRequest,
			wantBody:   "invalid body",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := presenter.NewAccountPresenter()
			input := &ctrlmock.AccountInputStub{CreateErr: tt.createErr}
			ctrl := NewAccountController(
				func(repo port.AccountRepository, output port.AccountOutputPort) port.AccountInputPort {
					input.Output = output
					return input
				},
				func() *presenter.AccountPresenter { return p },
				func() port.AccountRepository { return nil },
			)

			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/api/accounts/auth", bytes.NewBufferString(tt.body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			_ = ctrl.CreateOrGet(c)
			assertStatusBody(t, rec, tt.wantStatus, tt.wantBody)
			if tt.wantStatus == http.StatusOK && (p.Response() == nil || p.Response().Email != "user@example.com") {
				t.Fatalf("presenter response not set: %+v", p.Response())
			}
		})
	}
}

func TestAccountController_GetCurrent(t *testing.T) {
	tests := []struct {
		name       string
		headerID   string
		getErr     error
		wantStatus int
		wantBody   string
	}{
		{name: "[Success] with header", headerID: "acc-1", wantStatus: http.StatusOK},
		{name: "[Fail] missing header", headerID: "", wantStatus: http.StatusForbidden, wantBody: domainerr.ErrUnauthorized.Error()},
		{name: "[Fail] usecase not found", headerID: "acc-1", getErr: domainerr.ErrNotFound, wantStatus: http.StatusNotFound, wantBody: domainerr.ErrNotFound.Error()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := presenter.NewAccountPresenter()
			input := &ctrlmock.AccountInputStub{GetErr: tt.getErr}
			ctrl := NewAccountController(
				func(repo port.AccountRepository, output port.AccountOutputPort) port.AccountInputPort {
					input.Output = output
					return input
				},
				func() *presenter.AccountPresenter { return p },
				func() port.AccountRepository { return nil },
			)
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/api/accounts/me", nil)
			if tt.headerID != "" {
				req.Header.Set("X-Account-ID", tt.headerID)
			}
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			_ = ctrl.GetCurrent(c)
			assertStatusBody(t, rec, tt.wantStatus, tt.wantBody)
		})
	}
}

func TestAccountController_GetByID(t *testing.T) {
	tests := []struct {
		name       string
		getErr     error
		wantStatus int
		wantBody   string
	}{
		{name: "[Success] get by id", wantStatus: http.StatusOK},
		{name: "[Fail] not found", getErr: domainerr.ErrNotFound, wantStatus: http.StatusNotFound, wantBody: domainerr.ErrNotFound.Error()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := presenter.NewAccountPresenter()
			input := &ctrlmock.AccountInputStub{GetErr: tt.getErr}
			ctrl := NewAccountController(
				func(repo port.AccountRepository, output port.AccountOutputPort) port.AccountInputPort {
					input.Output = output
					return input
				},
				func() *presenter.AccountPresenter { return p },
				func() port.AccountRepository { return nil },
			)

			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/api/accounts/acc-1", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			_ = ctrl.GetByID(c, "acc-1")
			assertStatusBody(t, rec, tt.wantStatus, tt.wantBody)
		})
	}
}

func TestAccountController_GetAccountByEmail(t *testing.T) {
	tests := []struct {
		name       string
		email      string
		getErr     error
		wantStatus int
		wantBody   string
	}{
		{name: "[Success] get by email", email: "user@example.com", wantStatus: http.StatusOK},
		{name: "[Fail] not found", email: "missing@example.com", getErr: domainerr.ErrNotFound, wantStatus: http.StatusNotFound, wantBody: domainerr.ErrNotFound.Error()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := presenter.NewAccountPresenter()
			input := &ctrlmock.AccountInputStub{GetErr: tt.getErr}
			ctrl := NewAccountController(
				func(repo port.AccountRepository, output port.AccountOutputPort) port.AccountInputPort {
					input.Output = output
					return input
				},
				func() *presenter.AccountPresenter { return p },
				func() port.AccountRepository { return nil },
			)

			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/api/accounts/by-email?email="+tt.email, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// Parse query params
			_ = ctrl.GetAccountByEmail(c, openapi.AccountsGetAccountByEmailParams{Email: tt.email})
			assertStatusBody(t, rec, tt.wantStatus, tt.wantBody)
			if tt.wantStatus == http.StatusOK && (p.Response() == nil || p.Response().Email != tt.email) {
				t.Fatalf("presenter response not set or email mismatch: %+v", p.Response())
			}
		})
	}
}
