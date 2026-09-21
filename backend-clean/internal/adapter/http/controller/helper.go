package controller

import (
	"errors"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"

	openapi "immortal-architecture-cqrs/backend/internal/adapter/http/generated/openapi"
	"immortal-architecture-cqrs/backend/internal/domain/account"
	domainerr "immortal-architecture-cqrs/backend/internal/domain/errors"
)

func handleError(ctx echo.Context, err error) error {
	switch {
	case errors.Is(err, domainerr.ErrNotFound):
		return ctx.JSON(http.StatusNotFound, openapi.ModelsNotFoundError{Code: openapi.ModelsNotFoundErrorCodeNOTFOUND, Message: err.Error()})
	case errors.Is(err, domainerr.ErrUnauthorized):
		return ctx.JSON(http.StatusForbidden, openapi.ModelsForbiddenError{Code: openapi.ModelsForbiddenErrorCodeFORBIDDEN, Message: err.Error()})
	case errors.Is(err, account.ErrInvalidEmail), errors.Is(err, account.ErrInvalidName):
		return ctx.JSON(http.StatusBadRequest, openapi.ModelsBadRequestError{Code: openapi.ModelsBadRequestErrorCodeBADREQUEST, Message: err.Error()})
	case errors.Is(err, domainerr.ErrInvalidStatus) || errors.Is(err, domainerr.ErrInvalidStatusChange) || errors.Is(err, domainerr.ErrInvalidTemplateField):
		return ctx.JSON(http.StatusBadRequest, openapi.ModelsBadRequestError{Code: openapi.ModelsBadRequestErrorCodeBADREQUEST, Message: err.Error()})
	default:
		return ctx.JSON(http.StatusInternalServerError, openapi.ModelsErrorResponse{Code: "INTERNAL_ERROR", Message: err.Error()})
	}
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
