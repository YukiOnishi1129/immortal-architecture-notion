package usecase

import (
	"context"

	"immortal-architecture-cqrs/backend/internal/domain/note"
	"immortal-architecture-cqrs/backend/internal/port"
)

// NoteQueryInteractor handles note query (read) use cases.
type NoteQueryInteractor struct {
	readModelRepo port.NoteReadModelRepository
	output        port.NoteQueryOutputPort
}

var _ port.NoteQueryInputPort = (*NoteQueryInteractor)(nil)

// NewNoteQueryInteractor creates NoteQueryInteractor.
func NewNoteQueryInteractor(readModelRepo port.NoteReadModelRepository, output port.NoteQueryOutputPort) *NoteQueryInteractor {
	return &NoteQueryInteractor{
		readModelRepo: readModelRepo,
		output:        output,
	}
}

// List returns notes from the read model.
func (u *NoteQueryInteractor) List(ctx context.Context, filters note.Filters) error {
	notes, err := u.readModelRepo.List(ctx, filters)
	if err != nil {
		return err
	}
	return u.output.PresentNoteList(ctx, notes)
}

// Get returns a single note from the read model.
func (u *NoteQueryInteractor) Get(ctx context.Context, id string) error {
	n, err := u.readModelRepo.Get(ctx, id)
	if err != nil {
		return err
	}
	return u.output.PresentNote(ctx, n)
}
