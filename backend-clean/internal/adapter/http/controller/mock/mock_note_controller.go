package mock

import (
	"context"

	"immortal-architecture-notion/backend/internal/domain/note"
	"immortal-architecture-notion/backend/internal/port"
)

// NoteCommandInputStub is a stub for the CQRS command input port.
type NoteCommandInputStub struct {
	Err    error
	Output port.NoteCommandOutputPort
}

func (s *NoteCommandInputStub) Create(ctx context.Context, input port.NoteCreateInput) error {
	if s.Output != nil && s.Err == nil {
		_ = s.Output.PresentNote(ctx, &note.WithMeta{Note: note.Note{ID: "note-1", OwnerID: input.OwnerID}})
	}
	return s.Err
}

func (s *NoteCommandInputStub) Update(ctx context.Context, input port.NoteUpdateInput) error {
	if s.Output != nil && s.Err == nil {
		_ = s.Output.PresentNote(ctx, &note.WithMeta{Note: note.Note{ID: input.ID, OwnerID: input.OwnerID}})
	}
	return s.Err
}

func (s *NoteCommandInputStub) ChangeStatus(ctx context.Context, input port.NoteStatusChangeInput) error {
	if s.Output != nil && s.Err == nil {
		_ = s.Output.PresentNote(ctx, &note.WithMeta{Note: note.Note{ID: input.ID, OwnerID: input.OwnerID, Status: input.Status}})
	}
	return s.Err
}

func (s *NoteCommandInputStub) Delete(ctx context.Context, id, ownerID string) error {
	if s.Output != nil && s.Err == nil {
		_ = s.Output.PresentNoteDeleted(ctx)
	}
	return s.Err
}

// NoteQueryInputStub is a stub for the CQRS query input port.
type NoteQueryInputStub struct {
	Err        error
	Output     port.NoteQueryOutputPort
	ReadModels []note.ReadModel
	ReadModel  *note.ReadModel
}

func (s *NoteQueryInputStub) List(ctx context.Context, filters note.Filters) error {
	if s.Output != nil && s.Err == nil {
		_ = s.Output.PresentNoteList(ctx, s.ReadModels)
	}
	return s.Err
}

func (s *NoteQueryInputStub) Get(ctx context.Context, id string) error {
	if s.Output != nil && s.Err == nil {
		resp := s.ReadModel
		if resp == nil {
			resp = &note.ReadModel{ID: id}
		}
		_ = s.Output.PresentNote(ctx, resp)
	}
	return s.Err
}
