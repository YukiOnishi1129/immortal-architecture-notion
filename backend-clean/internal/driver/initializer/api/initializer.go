// Package initializer wires dependencies for the API server.
package initializer

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	httpcontroller "immortal-architecture-notion/backend/internal/adapter/http/controller"
	openapi "immortal-architecture-notion/backend/internal/adapter/http/generated/openapi"
	"immortal-architecture-notion/backend/internal/driver/config"
	driverdb "immortal-architecture-notion/backend/internal/driver/db"
	"immortal-architecture-notion/backend/internal/driver/factory"
	httpfactory "immortal-architecture-notion/backend/internal/driver/factory/http"
)

// BuildServer composes all dependencies and returns an Echo server, config, and cleanup function.
func BuildServer(ctx context.Context) (*echo.Echo, *config.Config, func(), error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, nil, func() {}, err
	}

	pool, err := driverdb.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		return nil, nil, func() {}, err
	}
	cleanup := func() {
		pool.Close()
	}

	txMgr := driverdb.NewTxManager(pool)

	accountRepoFactory := factory.NewAccountRepoFactory(pool)
	templateRepoFactory := factory.NewTemplateRepoFactory(pool)
	noteRepoFactory := factory.NewNoteRepoFactory(pool)
	noteReadModelRepoFactory := factory.NewNoteReadModelRepoFactory(pool)
	txFactory := factory.NewTxFactory(txMgr)

	accountOutputFactory := httpfactory.NewAccountOutputFactory()
	templateOutputFactory := httpfactory.NewTemplateOutputFactory()
	noteCommandOutputFactory := httpfactory.NewNoteCommandOutputFactory()
	noteQueryOutputFactory := httpfactory.NewNoteQueryOutputFactory()

	accountInputFactory := factory.NewAccountInputFactory()
	templateInputFactory := factory.NewTemplateInputFactory()
	noteCommandInputFactory := factory.NewNoteCommandInputFactory()
	noteQueryInputFactory := factory.NewNoteQueryInputFactory()

	e := echo.New()

	// Allow frontend (localhost:3000) to call the API during development.
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: cfg.AllowedOrigins,
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.PATCH, echo.OPTIONS},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			echo.HeaderAuthorization,
		},
	}))

	ac := httpcontroller.NewAccountController(accountInputFactory, accountOutputFactory, accountRepoFactory)
	nc := httpcontroller.NewNoteController(
		noteCommandInputFactory, noteCommandOutputFactory,
		noteQueryInputFactory, noteQueryOutputFactory,
		noteRepoFactory, noteReadModelRepoFactory,
		templateRepoFactory, txFactory,
	)
	tc := httpcontroller.NewTemplateController(templateInputFactory, templateOutputFactory, templateRepoFactory, txFactory)
	server := httpcontroller.NewServer(ac, nc, tc)
	openapi.RegisterHandlers(e, server)

	return e, cfg, cleanup, nil
}
