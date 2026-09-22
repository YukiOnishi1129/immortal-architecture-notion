//go:build e2e

package e2e

import (
	"net/http"
	"testing"
)

// updateTemplate sends a template update. Fields are passed through as given,
// so a test can decide whether to keep their ids.
func (a *app) updateTemplate(templateID, ownerID string, body map[string]any) response {
	a.t.Helper()
	return a.do(http.MethodPut, "/api/templates/"+templateID+"?ownerId="+ownerID, body)
}

// fieldsOf turns a template's fields back into an update payload, keeping ids.
func fieldsOf(tpl templateResponse) []map[string]any {
	out := make([]map[string]any, 0, len(tpl.Fields))
	for _, f := range tpl.Fields {
		// Echo back exactly what the API returned, so a test only changes what
		// it means to change.
		out = append(out, map[string]any{
			"id": f.ID, "label": f.Label, "order": f.Order, "isRequired": f.IsRequired,
		})
	}
	return out
}

// readTemplate fetches a template through the API.
func (a *app) readTemplate(templateID, ownerID string) templateResponse {
	a.t.Helper()
	res := a.do(http.MethodGet, "/api/templates/"+templateID+"?ownerId="+ownerID, nil).
		expect(a.t, http.StatusOK)
	var out templateResponse
	res.decode(a.t, &out)
	return out
}

// QA-14 a template already used by a note can still be given a parent page.
// Without this there is no way to switch existing templates on, since every
// template worth linking is one that already has notes.
func TestE2E_InUseTemplateCanBeLinkedToNotion(t *testing.T) {
	a := newApp(t)
	owner := seedAccount(t)
	tpl := a.createTemplate(owner, unlinkedTemplate)
	a.createNote(owner, tpl, draftNote) // the template is now in use

	res := a.updateTemplate(tpl.ID, owner, map[string]any{
		"id":                  tpl.ID,
		"name":                tpl.Name,
		"fields":              fieldsOf(tpl),
		"notionParentPageUrl": parentPageURL,
	})
	if res.code != http.StatusOK {
		t.Fatalf("status = %d, want 200: %s", res.code, res.body)
	}

	got := a.readTemplate(tpl.ID, owner)
	if got.NotionParentPageURL == nil {
		t.Fatalf("notionParentPageUrl is nil, want the parent page to be stored")
	}
	// The id is normalized, so the stored URL is the canonical form.
	if *got.NotionParentPageURL != "https://www.notion.so/"+parentPageID {
		t.Fatalf("stored url = %q, want the normalized parent page", *got.NotionParentPageURL)
	}
}

// QA-37 fields a note already wrote to cannot be renamed or removed, because
// the note stores its content per field.
func TestE2E_InUseTemplateFieldsAreProtected(t *testing.T) {
	tests := []struct {
		name     string
		mutate   func(fields []map[string]any) []map[string]any
		wantCode int
	}{
		{
			name: "[Fail] renaming a used field",
			mutate: func(f []map[string]any) []map[string]any {
				f[0]["label"] = "書き換えた"
				return f
			},
			wantCode: http.StatusBadRequest,
		},
		{
			name: "[Fail] removing a used field",
			mutate: func(f []map[string]any) []map[string]any {
				return f[:0]
			},
			wantCode: http.StatusBadRequest,
		},
		{
			name: "[Fail] dropping the ids looks like a delete and re-create",
			mutate: func(f []map[string]any) []map[string]any {
				for _, field := range f {
					delete(field, "id")
				}
				return f
			},
			wantCode: http.StatusBadRequest,
		},
		{
			name: "[Success] adding a field leaves existing notes alone",
			mutate: func(f []map[string]any) []map[string]any {
				return append(f, map[string]any{
					"label": "補足", "order": len(f) + 1, "isRequired": false,
				})
			},
			wantCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := newApp(t)
			owner := seedAccount(t)
			tpl := a.createTemplate(owner, linkedTemplate)
			note := a.createNote(owner, tpl, draftNote)

			res := a.updateTemplate(tpl.ID, owner, map[string]any{
				"id":     tpl.ID,
				"name":   tpl.Name,
				"fields": tt.mutate(fieldsOf(tpl)),
			})
			if res.code != tt.wantCode {
				t.Fatalf("status = %d, want %d: %s", res.code, tt.wantCode, res.body)
			}

			// Whatever happened, the note must still read back intact.
			got := a.getNote(note.ID, owner)
			if len(got.Sections) == 0 {
				t.Fatalf("the note lost its sections")
			}
			if got.Sections[0].Content != draftNote.Body {
				t.Fatalf("section content = %q, want %q", got.Sections[0].Content, draftNote.Body)
			}
		})
	}
}

// A template no note uses yet can be edited freely, which is the behaviour
// the protection above must not break.
func TestE2E_UnusedTemplateFieldsStayEditable(t *testing.T) {
	a := newApp(t)
	owner := seedAccount(t)
	tpl := a.createTemplate(owner, linkedTemplate)

	fields := fieldsOf(tpl)
	fields[0]["label"] = "書き換えた"

	a.updateTemplate(tpl.ID, owner, map[string]any{
		"id": tpl.ID, "name": tpl.Name, "fields": fields,
	}).expect(t, http.StatusOK)

	got := a.readTemplate(tpl.ID, owner)
	if got.Fields[0].Label != "書き換えた" {
		t.Fatalf("label = %q, want it renamed", got.Fields[0].Label)
	}
}

// Reordering fields must work even though (template_id, "order") is unique:
// moving a field onto a position another one still holds would clash unless
// the update is done in two passes.
func TestE2E_FieldsCanBeReordered(t *testing.T) {
	a := newApp(t)
	owner := seedAccount(t)

	res := a.do(http.MethodPost, "/api/templates", map[string]any{
		"name": "Report", "ownerId": owner,
		"fields": []map[string]any{
			{"label": "A", "order": 1, "isRequired": true},
			{"label": "B", "order": 2, "isRequired": false},
			{"label": "C", "order": 3, "isRequired": false},
		},
	}).expect(t, http.StatusOK)
	var tpl templateResponse
	res.decode(t, &tpl)

	// Reverse the order, keeping every id.
	fields := fieldsOf(tpl)
	reversed := make([]map[string]any, 0, len(fields))
	for i := len(fields) - 1; i >= 0; i-- {
		f := fields[i]
		f["order"] = len(fields) - i
		reversed = append(reversed, f)
	}

	a.updateTemplate(tpl.ID, owner, map[string]any{
		"id": tpl.ID, "name": tpl.Name, "fields": reversed,
	}).expect(t, http.StatusOK)

	got := a.readTemplate(tpl.ID, owner)
	labels := make([]string, 0, len(got.Fields))
	for _, f := range got.Fields {
		labels = append(labels, f.Label)
	}
	if len(labels) != 3 || labels[0] != "C" || labels[1] != "B" || labels[2] != "A" {
		t.Fatalf("order = %v, want [C B A]", labels)
	}
}

// Reordering and adding a field in the same request also has to avoid the
// unique constraint, since the new row needs a position too.
func TestE2E_FieldsCanBeReorderedWhileAdding(t *testing.T) {
	a := newApp(t)
	owner := seedAccount(t)

	res := a.do(http.MethodPost, "/api/templates", map[string]any{
		"name": "Report", "ownerId": owner,
		"fields": []map[string]any{
			{"label": "A", "order": 1, "isRequired": true},
			{"label": "B", "order": 2, "isRequired": false},
		},
	}).expect(t, http.StatusOK)
	var tpl templateResponse
	res.decode(t, &tpl)

	fields := fieldsOf(tpl)
	swapped := []map[string]any{fields[1], fields[0]}
	swapped[0]["order"] = 1
	swapped[1]["order"] = 2
	swapped = append(swapped, map[string]any{
		"label": "C", "order": 3, "isRequired": false,
	})

	a.updateTemplate(tpl.ID, owner, map[string]any{
		"id": tpl.ID, "name": tpl.Name, "fields": swapped,
	}).expect(t, http.StatusOK)

	got := a.readTemplate(tpl.ID, owner)
	labels := make([]string, 0, len(got.Fields))
	for _, f := range got.Fields {
		labels = append(labels, f.Label)
	}
	if len(labels) != 3 || labels[0] != "B" || labels[1] != "A" || labels[2] != "C" {
		t.Fatalf("order = %v, want [B A C]", labels)
	}
}

// Orders are only required to be positive and unique, so they can be sparse.
// The parking positions used while reordering must clear the highest order in
// play, not just the number of fields, or a requested order lands on one.
func TestE2E_SparseOrdersCanBeReordered(t *testing.T) {
	a := newApp(t)
	owner := seedAccount(t)

	res := a.do(http.MethodPost, "/api/templates", map[string]any{
		"name": "Report", "ownerId": owner,
		"fields": []map[string]any{
			{"label": "A", "order": 1, "isRequired": true},
			{"label": "B", "order": 7, "isRequired": false},
		},
	}).expect(t, http.StatusOK)
	var tpl templateResponse
	res.decode(t, &tpl)

	// Put a note on it so the in-place path is taken.
	a.do(http.MethodPost, "/api/notes", map[string]any{
		"title": "n", "templateId": tpl.ID, "ownerId": owner,
		"sections": []map[string]any{
			{"fieldId": tpl.Fields[0].ID, "content": "x"},
			{"fieldId": tpl.Fields[1].ID, "content": "y"},
		},
	}).expect(t, http.StatusOK)

	// Move B to 6 and add a field at 7. With a count-based parking offset the
	// new row would collide with a parked one.
	fields := fieldsOf(tpl)
	fields[1]["order"] = 6
	fields = append(fields, map[string]any{
		"label": "C", "order": 7, "isRequired": false,
	})

	a.updateTemplate(tpl.ID, owner, map[string]any{
		"id": tpl.ID, "name": tpl.Name, "fields": fields,
	}).expect(t, http.StatusOK)

	got := a.readTemplate(tpl.ID, owner)
	if len(got.Fields) != 3 {
		t.Fatalf("fields = %d, want 3", len(got.Fields))
	}
	if got.Fields[1].Order != 6 || got.Fields[2].Order != 7 {
		t.Fatalf("orders = %d,%d want 6,7", got.Fields[1].Order, got.Fields[2].Order)
	}
}
