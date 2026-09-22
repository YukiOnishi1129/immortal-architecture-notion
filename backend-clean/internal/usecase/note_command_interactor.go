package usecase

import (
	"context"
	"strings"

	domainerr "immortal-architecture-notion/backend/internal/domain/errors"
	"immortal-architecture-notion/backend/internal/domain/note"
	"immortal-architecture-notion/backend/internal/domain/service"
	"immortal-architecture-notion/backend/internal/port"
)

// NoteCommandInteractor handles note command (write) use cases.
type NoteCommandInteractor struct {
	notes         port.NoteRepository
	readModelRepo port.NoteReadModelRepository
	templates     port.TemplateRepository
	tx            port.TxManager
	output        port.NoteCommandOutputPort

	// notion is nil when the integration is not configured.
	// Notes stay fully usable in that case; only syncing is skipped.
	notion port.NotionClient
}

var _ port.NoteCommandInputPort = (*NoteCommandInteractor)(nil)

// NewNoteCommandInteractor creates NoteCommandInteractor.
func NewNoteCommandInteractor(
	notes port.NoteRepository,
	readModelRepo port.NoteReadModelRepository,
	templates port.TemplateRepository,
	tx port.TxManager,
	output port.NoteCommandOutputPort,
	notionClient port.NotionClient,
) *NoteCommandInteractor {
	return &NoteCommandInteractor{
		notes:         notes,
		readModelRepo: readModelRepo,
		templates:     templates,
		tx:            tx,
		output:        output,
		notion:        notionClient,
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

	// The Notion page must reflect what is about to be saved, so the edited
	// content is assembled before the call.
	edited := u.editedNote(current, input)
	synced, err := u.syncEdit(ctx, edited)
	if err != nil {
		return err
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

		if synced != nil {
			if err := u.notes.SaveNotionPage(txCtx, input.ID,
				synced.pageID, synced.pageURL, synced.syncedAt); err != nil {
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

	// Notion is called first, outside any transaction. If it fails the
	// database is untouched and the note keeps its current status, which is
	// the agreed behavior: a note is never published without its page.
	synced, err := u.syncStatusChange(ctx, current, input.Status)
	if err != nil {
		return err
	}

	err = u.tx.WithinTransaction(ctx, func(txCtx context.Context) error {
		if _, err := u.notes.UpdateStatus(txCtx, input.ID, input.Status); err != nil {
			return err
		}
		if synced != nil {
			if err := u.notes.SaveNotionPage(txCtx, input.ID,
				synced.pageID, synced.pageURL, synced.syncedAt); err != nil {
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
		// The page exists on Notion but its id was never stored, so clean it up.
		u.cleanUpOrphanPage(ctx, synced, err)
		return err
	}

	n, err := u.notes.Get(ctx, input.ID)
	if err != nil {
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
		NotionPageURL:  wm.Note.NotionPageURL,
	}
}
