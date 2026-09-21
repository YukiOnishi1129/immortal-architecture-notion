// Package factory provides constructors for driver-level wiring.
package factory

import (
	"immortal-architecture-notion/backend/internal/port"
	"immortal-architecture-notion/backend/internal/usecase"
)

// NewAccountInputFactory returns a factory for AccountInteractor.
func NewAccountInputFactory() func(repo port.AccountRepository, output port.AccountOutputPort) port.AccountInputPort {
	return func(repo port.AccountRepository, output port.AccountOutputPort) port.AccountInputPort {
		return usecase.NewAccountInteractor(repo, output)
	}
}

// NewTemplateInputFactory returns a factory for TemplateInteractor.
func NewTemplateInputFactory() func(repo port.TemplateRepository, tx port.TxManager, output port.TemplateOutputPort) port.TemplateInputPort {
	return func(repo port.TemplateRepository, tx port.TxManager, output port.TemplateOutputPort) port.TemplateInputPort {
		return usecase.NewTemplateInteractor(repo, tx, output)
	}
}

// NewNoteCommandInputFactory returns a factory for NoteCommandInteractor.
func NewNoteCommandInputFactory() func(noteRepo port.NoteRepository, readModelRepo port.NoteReadModelRepository, tplRepo port.TemplateRepository, tx port.TxManager, output port.NoteCommandOutputPort) port.NoteCommandInputPort {
	return func(noteRepo port.NoteRepository, readModelRepo port.NoteReadModelRepository, tplRepo port.TemplateRepository, tx port.TxManager, output port.NoteCommandOutputPort) port.NoteCommandInputPort {
		return usecase.NewNoteCommandInteractor(noteRepo, readModelRepo, tplRepo, tx, output)
	}
}

// NewNoteQueryInputFactory returns a factory for NoteQueryInteractor.
func NewNoteQueryInputFactory() func(readModelRepo port.NoteReadModelRepository, output port.NoteQueryOutputPort) port.NoteQueryInputPort {
	return func(readModelRepo port.NoteReadModelRepository, output port.NoteQueryOutputPort) port.NoteQueryInputPort {
		return usecase.NewNoteQueryInteractor(readModelRepo, output)
	}
}
