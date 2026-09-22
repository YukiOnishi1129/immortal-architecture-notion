//go:build e2e

// Package e2e runs the HTTP API against a real PostgreSQL database.
//
// Only Notion is faked: everything else (Echo, the controllers, the use cases,
// sqlc, the read model) is the production wiring, so these tests cover the
// paths that unit tests with mocked repositories cannot.
//
// Run with: make test-e2e
package e2e

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	containerName = "mini-notion-e2e-db"
	dbUser        = "e2e"
	dbPassword    = "e2e"
	dbName        = "mini_notion_e2e"
	hostPort      = "55433"

	// setupTimeout bounds starting the container and migrating it.
	setupTimeout = 2 * time.Minute
	// readyTimeout bounds waiting for the database to accept connections.
	readyTimeout = 60 * time.Second
)

var testPool *pgxpool.Pool

func TestMain(m *testing.M) {
	// Bring the database up under a deadline. A hung Docker daemon would
	// otherwise block the whole run with no output, which is painful in CI.
	// os.Exit skips deferred calls, so the context is released explicitly
	// before the tests start rather than with defer.
	setupCtx, cancelSetup := context.WithTimeout(context.Background(), setupTimeout)

	dsn, stop, err := startPostgres(setupCtx)
	if err != nil {
		cancelSetup()
		fmt.Fprintf(os.Stderr, "e2e: cannot start postgres: %v\n", err)
		os.Exit(1)
	}

	pool, err := pgxpool.New(setupCtx, dsn)
	if err != nil {
		cancelSetup()
		stop()
		fmt.Fprintf(os.Stderr, "e2e: cannot connect: %v\n", err)
		os.Exit(1)
	}
	testPool = pool

	if err := applyMigrations(setupCtx, pool); err != nil {
		cancelSetup()
		pool.Close()
		stop()
		fmt.Fprintf(os.Stderr, "e2e: migrations failed: %v\n", err)
		os.Exit(1)
	}

	// The deadline covers setup only: the tests must not inherit it, and the
	// teardown below still has to run after a slow suite.
	cancelSetup()

	code := m.Run()

	pool.Close()
	stop()
	os.Exit(code)
}

// startPostgres runs a throwaway PostgreSQL container and waits for it to
// accept connections. The container is removed when the tests finish.
func startPostgres(ctx context.Context) (string, func(), error) {
	// Remove a container left behind by an interrupted run.
	_ = exec.Command("docker", "rm", "-f", containerName).Run() //nolint:errcheck // best effort

	run := exec.CommandContext(ctx, "docker", "run", "-d", "--rm",
		"--name", containerName,
		"-e", "POSTGRES_USER="+dbUser,
		"-e", "POSTGRES_PASSWORD="+dbPassword,
		"-e", "POSTGRES_DB="+dbName,
		"-p", hostPort+":5432",
		"postgres:16-alpine",
	)
	if out, err := run.CombinedOutput(); err != nil {
		return "", func() {}, fmt.Errorf("docker run: %v: %s", err, out)
	}

	stop := func() {
		_ = exec.Command("docker", "rm", "-f", containerName).Run() //nolint:errcheck // best effort
	}

	dsn := fmt.Sprintf("postgres://%s:%s@localhost:%s/%s?sslmode=disable",
		dbUser, dbPassword, hostPort, dbName)

	if err := waitForPostgres(ctx, dsn); err != nil {
		stop()
		return "", func() {}, err
	}
	return dsn, stop, nil
}

func waitForPostgres(ctx context.Context, dsn string) error {
	deadline := time.Now().Add(readyTimeout)
	for time.Now().Before(deadline) {
		// Stop early when the caller's deadline has already passed, instead of
		// retrying against a context that can no longer succeed.
		if err := ctx.Err(); err != nil {
			return fmt.Errorf("waiting for postgres: %w", err)
		}

		pool, err := pgxpool.New(ctx, dsn)
		if err == nil {
			pingErr := pool.Ping(ctx)
			pool.Close()
			if pingErr == nil {
				return nil
			}
		}
		time.Sleep(500 * time.Millisecond)
	}
	return fmt.Errorf("postgres did not become ready within %s", readyTimeout)
}

// applyMigrations runs the .up.sql files in order, so the schema under test is
// the same one the application uses.
func applyMigrations(ctx context.Context, pool *pgxpool.Pool) error {
	dir, err := filepath.Abs(filepath.Join("..", "..", "migrations"))
	if err != nil {
		return err
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	files := make([]string, 0, len(entries))
	for _, e := range entries {
		if strings.HasSuffix(e.Name(), ".up.sql") {
			files = append(files, e.Name())
		}
	}
	sort.Strings(files)

	for _, name := range files {
		sqlBytes, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			return err
		}
		if _, err := pool.Exec(ctx, string(sqlBytes)); err != nil {
			return fmt.Errorf("%s: %w", name, err)
		}
	}
	return nil
}
