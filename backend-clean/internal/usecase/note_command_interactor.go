package usecase

import (
	"context"
	"strings"

	domainerr "immortal-architecture-cqrs/backend/internal/domain/errors"
	"immortal-architecture-cqrs/backend/internal/domain/note"
	"immortal-architecture-cqrs/backend/internal/domain/service"
	"immortal-architecture-cqrs/backend/internal/port"
)

// NoteCommandInteractor handles note command (write) use cases.
type NoteCommandInteractor struct {
	notes         port.NoteRepository
	readModelRepo port.NoteReadModelRepository
	templates     port.TemplateRepository
	tx            port.TxManager
	output        port.NoteCommandOutputPort
}

var _ port.NoteCommandInputPort = (*NoteCommandInteractor)(nil)

// NewNoteCommandInteractor creates NoteCommandInteractor.
func NewNoteCommandInteractor(
	notes port.NoteRepository,
	readModelRepo port.NoteReadModelRepository,
	templates port.TemplateRepository,
	tx port.TxManager,
	output port.NoteCommandOutputPort,
) *NoteCommandInteractor {
	return &NoteCommandInteractor{
		notes:         notes,
		readModelRepo: readModelRepo,
		templates:     templates,
		tx:            tx,
		output:        output,
	}
}

// Create creates a note and synchronizes the read model.
func (u *NoteCommandInteractor) Create(ctx context.Context, input port.NoteCreateInput) error {
	if input.OwnerID == "" {
		return domainerr.ErrOwnerRequired
	}

	tpl, err := u.templates.Get(ctx, input.TemplateID)
	if err != nil {
		return err
	}

	sections, err := buildSections("", input.Sections)
	if err != nil {
		return err
	}
	if err := note.ValidateNoteForCreate(input.Title, tpl.Template, sections); err != nil {
		return err
	}

	var noteID string
	err = u.tx.WithinTransaction(ctx, func(txCtx context.Context) error {
		newNote := note.Note{
			Title:      input.Title,
			TemplateID: tpl.Template.ID,
			OwnerID:    input.OwnerID,
			Status:     note.StatusDraft,
			Sections:   sections,
		}
		nn, err := u.notes.Create(txCtx, newNote)
		if err != nil {
			return err
		}
		noteID = nn.ID
		sectionsWithID, err := buildSections(noteID, input.Sections)
		if err != nil {
			return err
		}
		if err := note.ValidateSections(tpl.Template.Fields, sectionsWithID); err != nil {
			return err
		}
		if err := u.notes.ReplaceSections(txCtx, noteID, sectionsWithID); err != nil {
			return err
		}

		// Synchronize read model
		created, err := u.notes.Get(txCtx, noteID)
		if err != nil {
			return err
		}
		return u.readModelRepo.Upsert(txCtx, toReadModel(created))
	})
	if err != nil {
		return err
	}
	n, err := u.notes.Get(ctx, noteID)
	if err != nil {
		return err
	}
	return u.output.PresentNote(ctx, n)
}

// Update updates a note and synchronizes the read model.
func (u *NoteCommandInteractor) Update(ctx context.Context, input port.NoteUpdateInput) error {
	current, err := u.notes.Get(ctx, input.ID)
	if err != nil {
		return err
	}
	if err := note.ValidateNoteOwnership(current.Note.OwnerID, input.OwnerID); err != nil {
		return err
	}
	if strings.TrimSpace(input.Title) == "" {
		return domainerr.ErrTitleRequired
	}

	err = u.tx.WithinTransaction(ctx, func(txCtx context.Context) error {
		_, err := u.notes.Update(txCtx, note.Note{
			ID:    input.ID,
			Title: input.Title,
		})
		if err != nil {
			return err
		}
		if input.Sections != nil {
			tpl, err := u.templates.Get(ctx, current.Note.TemplateID)
			if err != nil {
				return err
			}
			sections, err := buildSectionsForUpdate(current.Sections, tpl.Template.Fields, input.Sections, current.Note.ID)
			if err != nil {
				return err
			}
			if err := note.ValidateSections(tpl.Template.Fields, sections); err != nil {
				return err
			}
			if err := u.notes.ReplaceSections(txCtx, input.ID, sections); err != nil {
				return err
			}
		}

		// Synchronize read model
		updated, err := u.notes.Get(txCtx, input.ID)
		if err != nil {
			return err
		}
		return u.readModelRepo.Upsert(txCtx, toReadModel(updated))
	})
	if err != nil {
		return err
	}
	n, err := u.notes.Get(ctx, input.ID)
	if err != nil {
		return err
	}
	return u.output.PresentNote(ctx, n)
}

// ChangeStatus changes note status and synchronizes the read model.
func (u *NoteCommandInteractor) ChangeStatus(ctx context.Context, input port.NoteStatusChangeInput) error {
	current, err := u.notes.Get(ctx, input.ID)
	if err != nil {
		return err
	}
	if err := note.ValidateNoteOwnership(current.Note.OwnerID, input.OwnerID); err != nil {
		return err
	}
	if err := input.Status.Validate(); err != nil {
		return err
	}
	if input.Status == note.StatusPublish {
		if err := service.CanPublish(current.Note, input.OwnerID); err != nil {
			return err
		}
	} else {
		if err := service.CanUnpublish(current.Note, input.OwnerID); err != nil {
			return err
		}
	}
	if err := note.CanChangeStatus(current.Note.Status, input.Status); err != nil {
		return err
	}

	if _, err := u.notes.UpdateStatus(ctx, input.ID, input.Status); err != nil {
		return err
	}

	// Synchronize read model
	n, err := u.notes.Get(ctx, input.ID)
	if err != nil {
		return err
	}
	if err := u.readModelRepo.Upsert(ctx, toReadModel(n)); err != nil {
		return err
	}
	return u.output.PresentNote(ctx, n)
}

// Delete deletes a note and removes the read model.
func (u *NoteCommandInteractor) Delete(ctx context.Context, id, ownerID string) error {
	current, err := u.notes.Get(ctx, id)
	if err != nil {
		return err
	}
	if err := note.ValidateNoteOwnership(current.Note.OwnerID, ownerID); err != nil {
		return err
	}

	err = u.tx.WithinTransaction(ctx, func(txCtx context.Context) error {
		if err := u.notes.Delete(txCtx, id); err != nil {
			return err
		}
		return u.readModelRepo.Delete(txCtx, id)
	})
	if err != nil {
		return err
	}
	return u.output.PresentNoteDeleted(ctx)
}

// toReadModel converts a WithMeta to a ReadModel for synchronization.
func toReadModel(wm *note.WithMeta) note.ReadModel {
	sections := make([]note.SectionReadModel, 0, len(wm.Sections))
	for _, s := range wm.Sections {
		sections = append(sections, note.SectionReadModel{
			ID:         s.Section.ID,
			FieldID:    s.Section.FieldID,
			FieldLabel: s.FieldLabel,
			FieldOrder: s.FieldOrder,
			IsRequired: s.IsRequired,
			Content:    s.Section.Content,
		})
	}
	return note.ReadModel{
		ID:             wm.Note.ID,
		Title:          wm.Note.Title,
		Status:         wm.Note.Status,
		TemplateID:     wm.Note.TemplateID,
		TemplateName:   wm.TemplateName,
		OwnerID:        wm.Note.OwnerID,
		OwnerFirstName: wm.OwnerFirstName,
		OwnerLastName:  wm.OwnerLastName,
		OwnerThumbnail: wm.OwnerThumbnail,
		Sections:       sections,
		CreatedAt:      wm.Note.CreatedAt,
		UpdatedAt:      wm.Note.UpdatedAt,
	}
}
