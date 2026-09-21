package port

import (
	"context"

	"immortal-architecture-cqrs/backend/internal/domain/note"
)

// NoteCommandInputPort defines command (write) use case inputs for notes.
type NoteCommandInputPort interface {
	Create(ctx context.Context, input NoteCreateInput) error
	Update(ctx context.Context, input NoteUpdateInput) error
	ChangeStatus(ctx context.Context, input NoteStatusChangeInput) error
	Delete(ctx context.Context, id, ownerID string) error
}

// NoteCommandOutputPort defines command (write) presenters for notes.
type NoteCommandOutputPort interface {
	PresentNote(ctx context.Context, note *note.WithMeta) error
	PresentNoteDeleted(ctx context.Context) error
}
