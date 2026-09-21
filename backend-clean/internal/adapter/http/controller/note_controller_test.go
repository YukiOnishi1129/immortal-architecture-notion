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
	"immortal-architecture-cqrs/backend/internal/domain/note"
	"immortal-architecture-cqrs/backend/internal/port"
)

// newTestController creates a NoteController with command and query stubs.
func newTestController(commandStub *ctrlmock.NoteCommandInputStub, queryStub *ctrlmock.NoteQueryInputStub) *NoteController {
	cp := presenter.NewNoteCommandPresenter()
	qp := presenter.NewNoteQueryPresenter()
	return NewNoteController(
		func(noteRepo port.NoteRepository, readModelRepo port.NoteReadModelRepository, tplRepo port.TemplateRepository, tx port.TxManager, output port.NoteCommandOutputPort) port.NoteCommandInputPort {
			commandStub.Output = output
			return commandStub
		},
		func() *presenter.NoteCommandPresenter { return cp },
		func(readModelRepo port.NoteReadModelRepository, output port.NoteQueryOutputPort) port.NoteQueryInputPort {
			queryStub.Output = output
			return queryStub
		},
		func() *presenter.NoteQueryPresenter { return qp },
		func() port.NoteRepository { return nil },
		func() port.NoteReadModelRepository { return nil },
		func() port.TemplateRepository { return nil },
		func() port.TxManager { return nil },
	)
}

func TestNoteController_Create(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		wantStatus int
		wantBody   string
	}{
		{
			name:       "[Success] create note",
			body:       `{"title":"Hello","templateId":"00000000-0000-0000-0000-000000000001","ownerId":"00000000-0000-0000-0000-000000000002","sections":[{"fieldId":"f1","content":"c1"}]}`,
			wantStatus: http.StatusOK,
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
			ctrl := newTestController(&ctrlmock.NoteCommandInputStub{}, &ctrlmock.NoteQueryInputStub{})

			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/api/notes", bytes.NewBufferString(tt.body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			_ = ctrl.Create(c)
			assertStatusBody(t, rec, tt.wantStatus, tt.wantBody)
		})
	}
}

func TestNoteController_List(t *testing.T) {
	tests := []struct {
		name       string
		filters    openapi.NotesListNotesParams
		inErr      error
		wantStatus int
		wantBody   string
	}{
		{name: "[Success] list notes", filters: openapi.NotesListNotesParams{}, wantStatus: http.StatusOK},
		{name: "[Fail] repo error", filters: openapi.NotesListNotesParams{}, inErr: domainerr.ErrNotFound, wantStatus: http.StatusNotFound, wantBody: domainerr.ErrNotFound.Error()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			queryStub := &ctrlmock.NoteQueryInputStub{
				ReadModels: []note.ReadModel{{ID: "n1"}},
				Err:        tt.inErr,
			}
			ctrl := newTestController(&ctrlmock.NoteCommandInputStub{}, queryStub)

			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/api/notes", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			_ = ctrl.List(c, tt.filters)
			assertStatusBody(t, rec, tt.wantStatus, tt.wantBody)
		})
	}
}

func TestNoteController_Get(t *testing.T) {
	tests := []struct {
		name       string
		inErr      error
		wantStatus int
		wantBody   string
	}{
		{name: "[Success] get note", wantStatus: http.StatusOK},
		{name: "[Fail] not found", inErr: domainerr.ErrNotFound, wantStatus: http.StatusNotFound, wantBody: domainerr.ErrNotFound.Error()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			queryStub := &ctrlmock.NoteQueryInputStub{Err: tt.inErr}
			ctrl := newTestController(&ctrlmock.NoteCommandInputStub{}, queryStub)

			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/api/notes/n1", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			_ = ctrl.GetByID(c, "n1")
			assertStatusBody(t, rec, tt.wantStatus, tt.wantBody)
		})
	}
}

func TestNoteController_Update(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		params     openapi.NotesUpdateNoteParams
		inErr      error
		wantStatus int
		wantBody   string
	}{
		{name: "[Success] update note", body: `{"title":"New","sections":[{"id":"sec1","content":"c"}]}`, params: openapi.NotesUpdateNoteParams{OwnerId: "owner"}, wantStatus: http.StatusOK},
		{name: "[Fail] missing owner", body: `{"title":"New","sections":[{"id":"sec1","content":"c"}]}`, params: openapi.NotesUpdateNoteParams{OwnerId: ""}, wantStatus: http.StatusForbidden, wantBody: domainerr.ErrUnauthorized.Error()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			commandStub := &ctrlmock.NoteCommandInputStub{Err: tt.inErr}
			ctrl := newTestController(commandStub, &ctrlmock.NoteQueryInputStub{})

			e := echo.New()
			req := httptest.NewRequest(http.MethodPut, "/api/notes/n1", bytes.NewBufferString(tt.body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			_ = ctrl.Update(c, "n1", tt.params)
			assertStatusBody(t, rec, tt.wantStatus, tt.wantBody)
		})
	}
}

func TestNoteController_Publish(t *testing.T) {
	tests := []struct {
		name       string
		ownerID    string
		inErr      error
		wantStatus int
		wantBody   string
	}{
		{name: "[Success] publish note", ownerID: "owner", wantStatus: http.StatusOK},
		{name: "[Fail] publish missing owner", ownerID: "", wantStatus: http.StatusInternalServerError, wantBody: domainerr.ErrOwnerRequired.Error()},
		{name: "[Fail] publish forbidden", ownerID: "owner", inErr: domainerr.ErrUnauthorized, wantStatus: http.StatusForbidden, wantBody: domainerr.ErrUnauthorized.Error()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			commandStub := &ctrlmock.NoteCommandInputStub{Err: tt.inErr}
			ctrl := newTestController(commandStub, &ctrlmock.NoteQueryInputStub{})

			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/api/notes/n1/publish", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			_ = ctrl.Publish(c, "n1", openapi.NotesPublishNoteParams{OwnerId: tt.ownerID})
			assertStatusBody(t, rec, tt.wantStatus, tt.wantBody)
		})
	}
}

func TestNoteController_Unpublish(t *testing.T) {
	tests := []struct {
		name       string
		ownerID    string
		inErr      error
		wantStatus int
		wantBody   string
	}{
		{name: "[Success] unpublish note", ownerID: "owner", wantStatus: http.StatusOK},
		{name: "[Fail] unpublish missing owner", ownerID: "", wantStatus: http.StatusInternalServerError, wantBody: domainerr.ErrOwnerRequired.Error()},
		{name: "[Fail] unpublish forbidden", ownerID: "owner", inErr: domainerr.ErrUnauthorized, wantStatus: http.StatusForbidden, wantBody: domainerr.ErrUnauthorized.Error()},
		{name: "[Fail] unpublish not found", ownerID: "owner", inErr: domainerr.ErrNotFound, wantStatus: http.StatusNotFound, wantBody: domainerr.ErrNotFound.Error()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			commandStub := &ctrlmock.NoteCommandInputStub{Err: tt.inErr}
			ctrl := newTestController(commandStub, &ctrlmock.NoteQueryInputStub{})

			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/api/notes/n1/unpublish", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			_ = ctrl.Unpublish(c, "n1", openapi.NotesUnpublishNoteParams{OwnerId: tt.ownerID})
			assertStatusBody(t, rec, tt.wantStatus, tt.wantBody)
		})
	}
}

func TestNoteController_Delete(t *testing.T) {
	tests := []struct {
		name       string
		ownerID    string
		inErr      error
		wantStatus int
		wantBody   string
	}{
		{name: "[Success] delete", ownerID: "owner", wantStatus: http.StatusOK},
		{name: "[Fail] owner missing", ownerID: "", wantStatus: http.StatusForbidden, wantBody: domainerr.ErrUnauthorized.Error()},
		{name: "[Fail] not found", ownerID: "owner", inErr: domainerr.ErrNotFound, wantStatus: http.StatusNotFound, wantBody: domainerr.ErrNotFound.Error()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			commandStub := &ctrlmock.NoteCommandInputStub{Err: tt.inErr}
			ctrl := newTestController(commandStub, &ctrlmock.NoteQueryInputStub{})

			e := echo.New()
			req := httptest.NewRequest(http.MethodDelete, "/api/notes/n1", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			_ = ctrl.Delete(c, "n1", openapi.NotesDeleteNoteParams{OwnerId: tt.ownerID})
			assertStatusBody(t, rec, tt.wantStatus, tt.wantBody)
		})
	}
}
