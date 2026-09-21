// Package controller exposes intentionally bloated HTTP handlers.
package controller

import (
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"

	sqldb "immortal-architecture-bad-api/backend/internal/db/sqlc"
	openapi "immortal-architecture-bad-api/backend/internal/generated/openapi"
)

func respondError(ctx echo.Context, status int, message string) error {
	return ctx.JSON(status, openapi.ModelsErrorResponse{
		Code:    http.StatusText(status),
		Message: message,
	})
}

func accountToResponse(acc *sqldb.Account) openapi.ModelsAccountResponse {
	return openapi.ModelsAccountResponse{
		Id:          uuidToString(acc.ID),
		Email:       acc.Email,
		FirstName:   acc.FirstName,
		LastName:    acc.LastName,
		FullName:    strings.TrimSpace(acc.FirstName + " " + acc.LastName),
		Thumbnail:   textToString(acc.Thumbnail),
		LastLoginAt: timestamptzToTime(acc.LastLoginAt),
		CreatedAt:   timestamptzToTime(acc.CreatedAt),
		UpdatedAt:   timestamptzToTime(acc.UpdatedAt),
	}
}

func uuidToString(id pgtype.UUID) string {
	if !id.Valid {
		return ""
	}
	u, err := uuid.FromBytes(id.Bytes[:])
	if err != nil {
		return ""
	}
	return u.String()
}

func timestamptzToTime(ts pgtype.Timestamptz) time.Time {
	if !ts.Valid {
		return time.Time{}
	}
	return ts.Time.UTC()
}

func textToString(val pgtype.Text) *string {
	if !val.Valid {
		return nil
	}
	s := strings.TrimSpace(val.String)
	if s == "" {
		return nil
	}
	return &s
}
