// Package http provides factory functions for HTTP adapters.
package http

import httppresenter "immortal-architecture-notion/backend/internal/adapter/http/presenter"

// NewAccountOutputFactory returns a factory for HTTP AccountPresenter.
func NewAccountOutputFactory() func() *httppresenter.AccountPresenter {
	return func() *httppresenter.AccountPresenter {
		return httppresenter.NewAccountPresenter()
	}
}

// NewTemplateOutputFactory returns a factory for HTTP TemplatePresenter.
func NewTemplateOutputFactory() func() *httppresenter.TemplatePresenter {
	return func() *httppresenter.TemplatePresenter {
		return httppresenter.NewTemplatePresenter()
	}
}

// NewNoteCommandOutputFactory returns a factory for HTTP NoteCommandPresenter.
func NewNoteCommandOutputFactory() func() *httppresenter.NoteCommandPresenter {
	return func() *httppresenter.NoteCommandPresenter {
		return httppresenter.NewNoteCommandPresenter()
	}
}

// NewNoteQueryOutputFactory returns a factory for HTTP NoteQueryPresenter.
func NewNoteQueryOutputFactory() func() *httppresenter.NoteQueryPresenter {
	return func() *httppresenter.NoteQueryPresenter {
		return httppresenter.NewNoteQueryPresenter()
	}
}
