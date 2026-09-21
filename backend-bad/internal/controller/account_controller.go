package controller

import (
	"errors"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"

	openapi "immortal-architecture-bad-api/backend/internal/generated/openapi"
	"immortal-architecture-bad-api/backend/internal/service"
)

// AccountsCreateOrGetAccount handles OAuth upsert requests.
func (c *Controller) AccountsCreateOrGetAccount(ctx echo.Context) error {
	var body openapi.AccountsCreateOrGetAccountJSONRequestBody
	if err := ctx.Bind(&body); err != nil {
		return respondError(ctx, http.StatusBadRequest, "invalid payload")
	}

	account, err := c.accountService.CreateOrGetAccount(ctx, service.CreateOrGetAccountInput{
		Email:             body.Email,
		FullName:          body.Name,
		Provider:          body.Provider,
		ProviderAccountID: body.ProviderAccountId,
		Thumbnail:         body.Thumbnail,
	})
	if err != nil {
		return respondError(ctx, http.StatusInternalServerError, "failed to upsert account")
	}
	return ctx.JSON(http.StatusOK, accountToResponse(account))
}

// AccountsGetCurrentAccount resolves account via header/query.
func (c *Controller) AccountsGetCurrentAccount(ctx echo.Context) error {
	accountID := ctx.Request().Header.Get("X-Account-ID")
	if accountID == "" {
		accountID = ctx.QueryParam("accountId")
	}
	if strings.TrimSpace(accountID) == "" {
		return respondError(ctx, http.StatusUnauthorized, "missing X-Account-ID header")
	}

	account, err := c.accountService.GetAccountByID(ctx, accountID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrAccountNotFound):
			return respondError(ctx, http.StatusNotFound, "account not found")
		case errors.Is(err, service.ErrInvalidAccountID):
			return respondError(ctx, http.StatusBadRequest, "invalid account id")
		default:
			return respondError(ctx, http.StatusInternalServerError, "failed to fetch account")
		}
	}
	return ctx.JSON(http.StatusOK, accountToResponse(account))
}

// AccountsGetAccountById returns account by ID.
// revive:disable-next-line:var-naming // Method name fixed by generated interface.
func (c *Controller) AccountsGetAccountById(ctx echo.Context, accountID string) error {
	if strings.TrimSpace(accountID) == "" {
		return respondError(ctx, http.StatusBadRequest, "account id is required")
	}
	account, err := c.accountService.GetAccountByID(ctx, accountID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrAccountNotFound):
			return respondError(ctx, http.StatusNotFound, "account not found")
		case errors.Is(err, service.ErrInvalidAccountID):
			return respondError(ctx, http.StatusBadRequest, "invalid account id")
		default:
			return respondError(ctx, http.StatusInternalServerError, "failed to fetch account")
		}
	}
	return ctx.JSON(http.StatusOK, accountToResponse(account))
}
