package presenter

import (
	"context"

	openapi "immortal-architecture-notion/backend/internal/adapter/http/generated/openapi"
	"immortal-architecture-notion/backend/internal/domain/note"
	"immortal-architecture-notion/backend/internal/port"
)

// NoteQueryPresenter converts note read models to OpenAPI responses for query operations.
type NoteQueryPresenter struct {
	note  *openapi.ModelsNoteResponse
	notes []openapi.ModelsNoteResponse
}

var _ port.NoteQueryOutputPort = (*NoteQueryPresenter)(nil)

// NewNoteQueryPresenter creates a new NoteQueryPresenter.
func NewNoteQueryPresenter() *NoteQueryPresenter {
	return &NoteQueryPresenter{}
}

// PresentNoteList stores note list response from read models.
func (p *NoteQueryPresenter) PresentNoteList(_ context.Context, notes []note.ReadModel) error {
	res := make([]openapi.ModelsNoteResponse, 0, len(notes))
	for _, n := range notes {
		res = append(res, readModelToResponse(n))
	}
	p.notes = res
	return nil
}

// PresentNote stores single note response from read model.
func (p *NoteQueryPresenter) PresentNote(_ context.Context, n *note.ReadModel) error {
	resp := readModelToResponse(*n)
	p.note = &resp
	return nil
}

// Note returns the last note response.
func (p *NoteQueryPresenter) Note() *openapi.ModelsNoteResponse {
	return p.note
}

// Notes returns the note list response.
func (p *NoteQueryPresenter) Notes() []openapi.ModelsNoteResponse {
	return p.notes
}

func readModelToResponse(n note.ReadModel) openapi.ModelsNoteResponse {
	sections := make([]openapi.ModelsSection, 0, len(n.Sections))
	for _, s := range n.Sections {
		sections = append(sections, openapi.ModelsSection{
			Id:         s.ID,
			FieldId:    s.FieldID,
			FieldLabel: s.FieldLabel,
			Content:    s.Content,
			IsRequired: s.IsRequired,
		})
	}
	return openapi.ModelsNoteResponse{
		Id:           n.ID,
		Title:        n.Title,
		TemplateId:   n.TemplateID,
		TemplateName: n.TemplateName,
		OwnerId:      n.OwnerID,
		Owner: openapi.ModelsAccountSummary{
			Id:        n.OwnerID,
			FirstName: n.OwnerFirstName,
			LastName:  n.OwnerLastName,
			Thumbnail: n.OwnerThumbnail,
		},
		Status:    openapi.ModelsNoteStatus(n.Status),
		Sections:  sections,
		CreatedAt: n.CreatedAt,
		UpdatedAt: n.UpdatedAt,
	}
}
