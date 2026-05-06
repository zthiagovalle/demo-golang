package repositories_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var (
	testCtx      = context.Background()
	testPool     *pgxpool.Pool
	testDataPath string
)

func TestMain(m *testing.M) {
	pgContainer, err := postgres.Run(testCtx,
		"postgres:17-alpine",
		postgres.WithDatabase("demo"),
		postgres.WithUsername("demo"),
		postgres.WithPassword("demo"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").WithOccurrence(2),
		),
	)
	if err != nil {
		log.Fatalf("postgres container: %v", err)
	}
	defer func() { _ = pgContainer.Terminate(testCtx) }()

	connStr, err := pgContainer.ConnectionString(testCtx, "sslmode=disable")
	if err != nil {
		log.Fatalf("connection string: %v", err)
	}

	migrationsPath := mustAbs("../../../migrations")
	if err := runMigrations(connStr, migrationsPath); err != nil {
		log.Fatalf("migrations: %v", err)
	}

	testPool, err = pgxpool.New(testCtx, connStr)
	if err != nil {
		log.Fatalf("pgx pool: %v", err)
	}
	defer testPool.Close()

	testDataPath = mustAbs("../../../development-environment/database/tests-dataset")

	os.Exit(m.Run())
}

func runMigrations(connStr, migrationsPath string) error {
	mig, err := migrate.New(fmt.Sprintf("file://%s", migrationsPath), connStr)
	if err != nil {
		return err
	}
	defer mig.Close()
	if err := mig.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}

func mustAbs(rel string) string {
	abs, err := filepath.Abs(rel)
	if err != nil {
		log.Fatalf("abs path: %v", err)
	}
	return abs
}

func loadDatasets(t *testing.T, files ...string) {
	t.Helper()
	for _, file := range files {
		path := filepath.Join(testDataPath, file)
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("read dataset %s: %v", file, err)
		}
		if _, err := testPool.Exec(testCtx, string(data)); err != nil {
			t.Fatalf("exec dataset %s: %v", file, err)
		}
	}
}
