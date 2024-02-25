package testhelper

import (
	"context"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"path/filepath"
	"time"
)

type TestDatabase struct {
	*postgres.PostgresContainer
	ConnectionString string
}

// CreateTestDatabase creates a database container that can be used  for testing
func CreateTestDatabase(ctx context.Context) (*TestDatabase, error) {
	pgContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:15.3-alpine"),
		postgres.WithInitScripts(filepath.Join("..", "db", "schema.sql")),
		postgres.WithDatabase("postgres"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		return nil, err
	}
	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return nil, err
	}
	return &TestDatabase{
		PostgresContainer: pgContainer,
		ConnectionString:  connStr,
	}, nil
}
