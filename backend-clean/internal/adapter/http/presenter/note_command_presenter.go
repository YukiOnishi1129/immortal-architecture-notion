package presenter

import (
	"context"

	openapi "immortal-architecture-cqrs/backend/internal/adapter/http/generated/openapi"
	"immortal-architecture-cqrs/backend/internal/domain/note"
	"immortal-architecture-cqrs/backend/internal/port"
)

// NoteCommandPresenter converts note domain models to OpenAPI responses for command operations.
type NoteCommandPresenter struct {
	note      *openapi.ModelsNoteResponse
	deletedOK bool
}

var _ port.NoteCommandOutputPort = (*NoteCommandPresenter)(nil)

// NewNoteCommandPresenter creates a new NoteCommandPresenter.
func NewNoteCommandPresenter() *NoteCommandPresenter {
	return &NoteCommandPresenter{}
}

// PresentNote stores single note response.
func (p *NoteCommandPresenter) PresentNote(_ context.Context, n *note.WithMeta) error {
	resp := toNoteResponse(*n)
	p.note = &resp
	return nil
}

// PresentNoteDeleted marks delete success.
func (p *NoteCommandPresenter) PresentNoteDeleted(_ context.Context) error {
	p.deletedOK = true
	return nil
}

// Note returns the last note response.
func (p *NoteCommandPresenter) Note() *openapi.ModelsNoteResponse {
	return p.note
}

// DeleteResponse returns deletion success response.
func (p *NoteCommandPresenter) DeleteResponse() openapi.ModelsSuccessResponse {
	return openapi.ModelsSuccessResponse{Success: p.deletedOK}
}
