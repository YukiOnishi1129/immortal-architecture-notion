package port

import (
	"context"

	"immortal-architecture-cqrs/backend/internal/domain/note"
)

// NoteQueryInputPort defines query (read) use case inputs for notes.
type NoteQueryInputPort interface {
	List(ctx context.Context, filters note.Filters) error
	Get(ctx context.Context, id string) error
}

// NoteQueryOutputPort defines query (read) presenters for notes.
type NoteQueryOutputPort interface {
	PresentNoteList(ctx context.Context, notes []note.ReadModel) error
	PresentNote(ctx context.Context, note *note.ReadModel) error
}

// NoteReadModelRepository abstracts read model persistence for notes.
// It handles both reading (List/Get) and synchronization (Upsert/Delete)
// called from the command side.
type NoteReadModelRepository interface {
	List(ctx context.Context, filters note.Filters) ([]note.ReadModel, error)
	Get(ctx context.Context, id string) (*note.ReadModel, error)
	Upsert(ctx context.Context, model note.ReadModel) error
	Delete(ctx context.Context, id string) error
}
