package storage

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
)

// RunTestsWithTempDB sets up a container with a temporary
// PostgreSQL database for testing.
//
// Returns the exit code from testRunner or 1 on setup failure.
func RunTestsWithTempDB(testRunner func() int) int { //nolint:funlen
	dockerPool, err := dockertest.NewPool("")
	if err != nil {
		slog.Error("Could not construct dockerPool", "error", err)
		return 1
	}

	err = dockerPool.Client.Ping()
	if err != nil {
		slog.Error("Could not connect to Docker", "error", err)
		return 1
	}

	// pull an image, create a container based on it and run it
	resource, err := dockerPool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "latest",
		Env: []string{
			"POSTGRES_DB=go_test",
			"POSTGRES_PASSWORD=secret",
			"listen_addresses = '*'",
		},
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		slog.Error("Could not start resource", "error", err)
		return 1
	}

	hostAndPort, err := getHostPort(resource, "5432/tcp")
	if err != nil {
		slog.Error("Could not get host and port", "error", err)
		return 1
	}

	databaseUrl := fmt.Sprintf(
		"postgres://postgres:secret@%s/go_test?sslmode=disable",
		hostAndPort,
	)

	slog.Info("Connecting to database on " + databaseUrl)

	// exponential backoff-retry, because the application
	// in the container might not be ready to accept connections yet
	dockerPool.MaxWait = 10 * time.Second //nolint:mnd
	err = dockerPool.Retry(func() error {
		_, err = Init(Config{
			DatabaseUrl: databaseUrl,
		})
		if err != nil {
			return err
		}

		return pool.Ping(context.Background())
	})

	if err != nil {
		slog.Error("Could not connect to database", "error", err)
		return 1
	}

	defer func() {
		if err = dockerPool.Purge(resource); err != nil {
			slog.Error("Could not purge resource", "error", err)
		}
	}()

	return testRunner()
}

func getHostPort(resource *dockertest.Resource, id string) (string, error) {
	dockerURL := os.Getenv("DOCKER_HOST")
	if dockerURL == "" {
		return resource.GetHostPort(id), nil
	}

	u, err := url.Parse(dockerURL)
	if err != nil {
		return "", err
	}

	return u.Hostname() + ":" + resource.GetPort(id), nil
}

// RunTestsWithMigratedDB applies all up migrations,
// runs the provided test runner,
// and then rolls back all migrations.
//
// Returns the exit code from testRunner or 1 on setup failure.
func RunTestsWithMigratedDB(testRunner func() int) int {
	m, err := migrator()
	if err != nil {
		slog.Error(err.Error())
		return 1
	}

	err = m.Up()
	if err != nil {
		slog.Error(err.Error())
		return 1
	}

	defer func(m *migrate.Migrate) {
		err = m.Down()
		if err != nil { // coverage-ignore
			slog.Error(err.Error())
		}

		err, err2 := m.Close()
		if err != nil { // coverage-ignore
			slog.Error(err.Error())
		}

		if err2 != nil { // coverage-ignore
			slog.Error(err2.Error())
		}
	}(m)

	return testRunner()
}

// TestWithMigratedDB applies all up migrations,
// then registers a cleanup function
// to roll them back after the test finishes.
//
// This helper is intended for individual tests.
func TestWithMigratedDB(t *testing.T) {
	m, err := migrator()
	if err != nil {
		t.Fatal(err)
	}

	err = m.Up()
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		err = m.Down()
		if err != nil { // coverage-ignore
			t.Fatal(err)
		}

		err, err2 := m.Close()
		if err != nil { // coverage-ignore
			t.Fatal(err)
		}

		if err2 != nil { // coverage-ignore
			t.Fatal(err2)
		}
	})
}

// TestWithTransaction sets up a rolled-back transaction for a test.
// Fails the test if setup fails.
func TestWithTransaction(t *testing.T) {
	tx, err := activePool().Begin(t.Context())
	if err != nil {
		t.Fatal(err)
	}

	// Create a savepoint to prevent any function
	// from performing a real commit.
	// This ensures changes are rolled back after the test.
	nestedTx, err := tx.Begin(t.Context())
	if err != nil {
		t.Fatal(err)
	}

	prevConn := Conn

	// Mock the Conn function
	// to return the nested transaction.
	Conn = func() Connection {
		return nestedTx
	}

	t.Cleanup(func() {
		// Use Background context,
		// because t.Context() is already closed on Cleanup
		ctx := context.Background() //nolint:usetesting
		if err = tx.Rollback(ctx); err != nil &&
			!errors.Is(err, pgx.ErrTxClosed) {
			slog.Error(err.Error())
		}

		Conn = prevConn
	})
}

func migrator() (*migrate.Migrate, error) {
	db := stdlib.OpenDBFromPool(activePool())

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return nil, err
	}

	// Compute path relative to this file (resolves to `.../storage/migrations`)
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		return nil, errors.New("unable to determine the current file path")
	}

	migrationsPath := "file://" + filepath.Join(
		filepath.Dir(currentFile),
		"migrations",
	)

	return migrate.NewWithDatabaseInstance(
		migrationsPath,
		activePool().Config().ConnConfig.Database,
		driver,
	)
}
