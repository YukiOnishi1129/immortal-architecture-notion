package service

import (
	"github.com/labstack/echo/v4"

	sqldb "immortal-architecture-bad-api/backend/internal/db/sqlc"
	openapi "immortal-architecture-bad-api/backend/internal/generated/openapi"
)

// AccountServicer represents the AccountService contract.
type AccountServicer interface {
	CreateOrGetAccount(echo.Context, CreateOrGetAccountInput) (*sqldb.Account, error)
	GetAccountByID(echo.Context, string) (*sqldb.Account, error)
}

// TemplateServicer represents the TemplateService contract.
type TemplateServicer interface {
	ListTemplates(echo.Context, TemplateFilters) ([]*openapi.ModelsTemplateResponse, error)
	CreateTemplate(echo.Context, string, string, []sqldb.Field) (*openapi.ModelsTemplateResponse, error)
	DeleteTemplate(echo.Context, string) error
	GetTemplate(echo.Context, string) (*openapi.ModelsTemplateResponse, error)
	UpdateTemplate(echo.Context, string, string, []openapi.ModelsUpdateFieldRequest) (*openapi.ModelsTemplateResponse, error)
}

// NoteServicer represents the NoteService contract.
type NoteServicer interface {
	ListNotes(echo.Context, NoteFilters) ([]*openapi.ModelsNoteResponse, error)
	CreateNote(echo.Context, string, string, string, map[string]string) (*openapi.ModelsNoteResponse, error)
	DeleteNote(echo.Context, string) error
	GetNote(echo.Context, string) (*openapi.ModelsNoteResponse, error)
	UpdateNote(echo.Context, string, string, map[string]string) (*openapi.ModelsNoteResponse, error)
	ChangeStatus(echo.Context, string, string) (*openapi.ModelsNoteResponse, error)
}
