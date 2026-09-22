package controller

import (
	"errors"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"

	openapi "immortal-architecture-notion/backend/internal/adapter/http/generated/openapi"
	"immortal-architecture-notion/backend/internal/domain/account"
	domainerr "immortal-architecture-notion/backend/internal/domain/errors"
)

func handleError(ctx echo.Context, err error) error {
	switch {
	case errors.Is(err, domainerr.ErrNotFound):
		return ctx.JSON(http.StatusNotFound, openapi.ModelsNotFoundError{Code: openapi.ModelsNotFoundErrorCodeNOTFOUND, Message: err.Error()})
	case errors.Is(err, domainerr.ErrUnauthorized):
		return ctx.JSON(http.StatusForbidden, openapi.ModelsForbiddenError{Code: openapi.ModelsForbiddenErrorCodeFORBIDDEN, Message: err.Error()})
	case errors.Is(err, account.ErrInvalidEmail), errors.Is(err, account.ErrInvalidName):
		return ctx.JSON(http.StatusBadRequest, openapi.ModelsBadRequestError{Code: openapi.ModelsBadRequestErrorCodeBADREQUEST, Message: err.Error()})
	case isBadRequest(err):
		return ctx.JSON(http.StatusBadRequest, openapi.ModelsBadRequestError{Code: openapi.ModelsBadRequestErrorCodeBADREQUEST, Message: err.Error()})
	case errors.Is(err, domainerr.ErrNotionSyncFailed):
		// The note was left untouched, so retrying the same action is safe.
		// The underlying cause is not exposed; it is only useful in the logs.
		return ctx.JSON(http.StatusInternalServerError, openapi.ModelsErrorResponse{
			Code:    "NOTION_SYNC_FAILED",
			Message: "failed to sync with Notion, please try again",
		})
	default:
		return ctx.JSON(http.StatusInternalServerError, openapi.ModelsErrorResponse{Code: "INTERNAL_ERROR", Message: err.Error()})
	}
}

// isBadRequest reports whether the error is the caller's fault rather than a
// server failure. Everything listed here is rejected before anything is
// written, so the client can fix the request and retry.
func isBadRequest(err error) bool {
	for _, target := range []error{
		domainerr.ErrInvalidStatus,
		domainerr.ErrInvalidStatusChange,
		domainerr.ErrInvalidTemplateField,
		domainerr.ErrFieldRequired,
		domainerr.ErrTemplateInUse,
		domainerr.ErrTemplateNameRequired,
		domainerr.ErrNotionParentNotSet,
		domainerr.ErrInvalidNotionParentURL,
	} {
		if errors.Is(err, target) {
			return true
		}
	}
	return false
}

func currentAccountID(ctx echo.Context) (string, error) {
	id := ctx.Request().Header.Get("X-Account-ID")
	if strings.TrimSpace(id) == "" {
		return "", domainerr.ErrUnauthorized
	}
	return id, nil
}

func valueOrEmpty(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
