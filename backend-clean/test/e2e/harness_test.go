//go:build e2e

package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"

	httpcontroller "immortal-architecture-notion/backend/internal/adapter/http/controller"
	openapi "immortal-architecture-notion/backend/internal/adapter/http/generated/openapi"
	driverdb "immortal-architecture-notion/backend/internal/driver/db"
	"immortal-architecture-notion/backend/internal/driver/factory"
	httpfactory "immortal-architecture-notion/backend/internal/driver/factory/http"
)

// app is the API under test, wired exactly as production does it except for
// the Notion client.
type app struct {
	server *echo.Echo
	notion *fakeNotion
	t      *testing.T

	// Kept so a test can move between the two wirings without losing data.
	notionEnabledServer  *echo.Echo
	notionDisabledServer *echo.Echo
}

func newApp(t *testing.T) *app {
	t.Helper()
	resetDatabase(t)

	notion := newFakeNotion()
	txMgr := driverdb.NewTxManager(testPool)

	ac := httpcontroller.NewAccountController(
		factory.NewAccountInputFactory(),
		httpfactory.NewAccountOutputFactory(),
		factory.NewAccountRepoFactory(testPool),
	)
	tc := httpcontroller.NewTemplateController(
		factory.NewTemplateInputFactory(),
		httpfactory.NewTemplateOutputFactory(),
		factory.NewTemplateRepoFactory(testPool),
		factory.NewTxFactory(txMgr),
	)
	nc := httpcontroller.NewNoteController(
		factory.NewNoteCommandInputFactory(notion),
		httpfactory.NewNoteCommandOutputFactory(),
		factory.NewNoteQueryInputFactory(),
		httpfactory.NewNoteQueryOutputFactory(),
		factory.NewNoteRepoFactory(testPool),
		factory.NewNoteReadModelRepoFactory(testPool),
		factory.NewTemplateRepoFactory(testPool),
		factory.NewTxFactory(txMgr),
	)

	e := echo.New()
	openapi.RegisterHandlers(e, httpcontroller.NewServer(ac, nc, tc))
	return &app{server: e, notionEnabledServer: e, notion: notion, t: t}
}

// newAppWithoutNotion wires the API with the integration switched off,
// which is what happens when NOTION_API_KEY is not set.
func newAppWithoutNotion(t *testing.T) *app {
	t.Helper()
	a := newApp(t)
	nc := httpcontroller.NewNoteController(
		factory.NewNoteCommandInputFactory(nil),
		httpfactory.NewNoteCommandOutputFactory(),
		factory.NewNoteQueryInputFactory(),
		httpfactory.NewNoteQueryOutputFactory(),
		factory.NewNoteRepoFactory(testPool),
		factory.NewNoteReadModelRepoFactory(testPool),
		factory.NewTemplateRepoFactory(testPool),
		factory.NewTxFactory(driverdb.NewTxManager(testPool)),
	)
	tc := httpcontroller.NewTemplateController(
		factory.NewTemplateInputFactory(),
		httpfactory.NewTemplateOutputFactory(),
		factory.NewTemplateRepoFactory(testPool),
		factory.NewTxFactory(driverdb.NewTxManager(testPool)),
	)
	ac := httpcontroller.NewAccountController(
		factory.NewAccountInputFactory(),
		httpfactory.NewAccountOutputFactory(),
		factory.NewAccountRepoFactory(testPool),
	)
	e := echo.New()
	openapi.RegisterHandlers(e, httpcontroller.NewServer(ac, nc, tc))
	// The fake Notion client is kept, so a test can switch back to the wired
	// server and act on data created while the integration was off.
	a.notionDisabledServer = e
	a.server = e
	return a
}

// withNotion switches back to the server that has the Notion client wired,
// leaving the data created while it was off in place.
func (a *app) withNotion() *app {
	a.t.Helper()
	a.server = a.notionEnabledServer
	return a
}

// resetDatabase clears every table so each test starts from a known state.
func resetDatabase(t *testing.T) {
	t.Helper()
	_, err := testPool.Exec(context.Background(),
		`TRUNCATE note_read_models, sections, notes, fields, templates, accounts CASCADE`)
	if err != nil {
		t.Fatalf("reset database: %v", err)
	}
}

type response struct {
	code int
	body []byte
}

// do sends a request through the real router.
func (a *app) do(method, path string, body any) response {
	a.t.Helper()

	var reader *bytes.Reader
	if body != nil {
		raw, err := json.Marshal(body)
		if err != nil {
			a.t.Fatalf("encode body: %v", err)
		}
		reader = bytes.NewReader(raw)
	} else {
		reader = bytes.NewReader(nil)
	}

	req := httptest.NewRequest(method, path, reader)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	a.server.ServeHTTP(rec, req)
	return response{code: rec.Code, body: rec.Body.Bytes()}
}

func (r response) decode(t *testing.T, out any) {
	t.Helper()
	if err := json.Unmarshal(r.body, out); err != nil {
		t.Fatalf("decode response (%d): %v: %s", r.code, err, r.body)
	}
}

// expect fails the test unless the status matches. The code is spelled out at
// every call site so the expected outcome is readable there.
//
//nolint:unparam // only 200 is asserted today; other codes are checked inline
func (r response) expect(t *testing.T, code int) response {
	t.Helper()
	if r.code != code {
		t.Fatalf("status = %d, want %d: %s", r.code, code, r.body)
	}
	return r
}

// seedAccount inserts an owner directly, since sign-in is out of scope here.
func seedAccount(t *testing.T) string {
	t.Helper()
	var id string
	err := testPool.QueryRow(context.Background(),
		`INSERT INTO accounts (email, first_name, last_name, provider, provider_account_id)
		 VALUES ('e2e@example.com', 'E2E', 'User', 'google', 'e2e-1') RETURNING id::text`).Scan(&id)
	if err != nil {
		t.Fatalf("seed account: %v", err)
	}
	return id
}

type templateResponse struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	IsUsed bool   `json:"isUsed"`
	Fields []struct {
		ID         string `json:"id"`
		Label      string `json:"label"`
		Order      int    `json:"order"`
		IsRequired bool   `json:"isRequired"`
	} `json:"fields"`
	NotionParentPageURL *string `json:"notionParentPageUrl"`
}

type noteResponse struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Status   string `json:"status"`
	Sections []struct {
		ID      string `json:"id"`
		Content string `json:"content"`
	} `json:"sections"`
	NotionPageURL *string `json:"notionPageUrl"`
}

// templateInput is the data a test sets up a template with.
// FieldLabel becomes the heading of the corresponding Notion block.
type templateInput struct {
	Name          string
	FieldLabel    string
	ParentPageURL string // empty means the template is not linked to Notion
}

// noteInput is the data a test sets up a note with.
// Title becomes the Notion page title, Body its only paragraph.
type noteInput struct {
	Title string
	Body  string
}

// createTemplate creates a template and returns the API response.
func (a *app) createTemplate(ownerID string, in templateInput) templateResponse {
	a.t.Helper()
	body := map[string]any{
		"name":    in.Name,
		"ownerId": ownerID,
		"fields": []map[string]any{
			{"label": in.FieldLabel, "order": 1, "isRequired": true},
		},
	}
	if in.ParentPageURL != "" {
		body["notionParentPageUrl"] = in.ParentPageURL
	}
	res := a.do(http.MethodPost, "/api/templates", body).expect(a.t, http.StatusOK)

	var out templateResponse
	res.decode(a.t, &out)
	return out
}

// createNote creates a draft note with one section and returns the response.
func (a *app) createNote(ownerID string, tpl templateResponse, in noteInput) noteResponse {
	a.t.Helper()
	res := a.do(http.MethodPost, "/api/notes", map[string]any{
		"title":      in.Title,
		"templateId": tpl.ID,
		"ownerId":    ownerID,
		"sections": []map[string]any{
			{"fieldId": tpl.Fields[0].ID, "content": in.Body},
		},
	}).expect(a.t, http.StatusOK)

	var out noteResponse
	res.decode(a.t, &out)
	return out
}

// updateNote sends an edit. Every section must be included, because an update
// replaces them all.
func (a *app) updateNote(noteID, ownerID string, in noteInput, sectionID string) response {
	a.t.Helper()
	return a.do(http.MethodPut, "/api/notes/"+noteID+"?ownerId="+ownerID, map[string]any{
		"title": in.Title,
		"sections": []map[string]any{
			{"id": sectionID, "content": in.Body},
		},
	})
}

// getNote reads a note back through the query side (the read model).
func (a *app) getNote(id, ownerID string) noteResponse {
	a.t.Helper()
	res := a.do(http.MethodGet, "/api/notes/"+id+"?ownerId="+ownerID, nil).expect(a.t, http.StatusOK)
	var out noteResponse
	res.decode(a.t, &out)
	return out
}
