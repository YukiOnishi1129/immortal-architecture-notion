package initializer

import (
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"

	httpcontroller "immortal-architecture-notion/backend/internal/adapter/http/controller"
	openapi "immortal-architecture-notion/backend/internal/adapter/http/generated/openapi"
	"immortal-architecture-notion/backend/internal/driver/factory"
	httpfactory "immortal-architecture-notion/backend/internal/driver/factory/http"
)

// Smoke test to ensure initializer wiring does not panic and returns non-nil server.
func TestNewServer_Wiring(t *testing.T) {
	// use nil pool since factories are functional closures; server wiring should not panic
	var pool *pgxpool.Pool
	ac := httpcontroller.NewAccountController(
		factory.NewAccountInputFactory(),
		httpfactory.NewAccountOutputFactory(),
		factory.NewAccountRepoFactory(pool),
	)
	tc := httpcontroller.NewTemplateController(
		factory.NewTemplateInputFactory(),
		httpfactory.NewTemplateOutputFactory(),
		factory.NewTemplateRepoFactory(pool),
		factory.NewTxFactory(nil),
	)
	nc := httpcontroller.NewNoteController(
		factory.NewNoteCommandInputFactory(),
		httpfactory.NewNoteCommandOutputFactory(),
		factory.NewNoteQueryInputFactory(),
		httpfactory.NewNoteQueryOutputFactory(),
		factory.NewNoteRepoFactory(pool),
		factory.NewNoteReadModelRepoFactory(pool),
		factory.NewTemplateRepoFactory(pool),
		factory.NewTxFactory(nil),
	)

	srv := httpcontroller.NewServer(ac, nc, tc)
	if srv == nil {
		t.Fatalf("server is nil")
	}

	e := echo.New()
	// ensure register does not panic
	openapi.RegisterHandlersWithBaseURL(e, srv, "")
}
