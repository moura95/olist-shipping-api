package repository_test

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"github/moura95/olist-shipping-api/internal/repository"
)

var (
	testDB        *sql.DB
	testQueries   *repository.Queries
	testContainer *postgres.PostgresContainer
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	var err error
	testContainer, err = postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:15-alpine"),
		postgres.WithDatabase("test_db"),
		postgres.WithUsername("test_user"),
		postgres.WithPassword("test_password"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30*time.Second)),
	)
	if err != nil {
		log.Fatalf("Failed to start PostgreSQL container: %v", err)
	}

	connStr, err := testContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		log.Fatalf("Failed to get connection string: %v", err)
	}

	testDB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err := testDB.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	if err := runMigrations(testDB); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	testQueries = repository.New(testDB)

	code := m.Run()

	if err := testContainer.Terminate(ctx); err != nil {
		log.Printf("Failed to terminate container: %v", err)
	}

	os.Exit(code)
}

func runMigrations(db *sql.DB) error {
	migrationsPath := "../../db/migrations"

	files, err := ioutil.ReadDir(migrationsPath)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %v", err)
	}

	var upFiles []string
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".sql" && filepath.Base(file.Name())[len(filepath.Base(file.Name()))-6:] == "up.sql" {
			upFiles = append(upFiles, file.Name())
		}
	}

	for _, fileName := range upFiles {
		content, err := ioutil.ReadFile(filepath.Join(migrationsPath, fileName))
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %v", fileName, err)
		}

		if _, err := db.Exec(string(content)); err != nil {
			return fmt.Errorf("failed to execute migration %s: %v", fileName, err)
		}
	}

	return nil
}

func cleanupTestData(t *testing.T) {
	ctx := context.Background()

	tables := []string{
		"packages",
	}

	for _, table := range tables {
		_, err := testDB.ExecContext(ctx, fmt.Sprintf("DELETE FROM %s", table))
		if err != nil {
			t.Logf("Failed to cleanup table %s: %v", table, err)
		}
	}
}
