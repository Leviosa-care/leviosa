package testutils

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib" // Your Postgres driver
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type PostgresContainer struct {
	testcontainers.Container
	ConnectionString string
}

func SetupPostgres(ctx context.Context, t *testing.T) (*PostgresContainer, error) {
	req := testcontainers.ContainerRequest{
		Image:        "postgres:17.5-alpine3.21",
		ExposedPorts: []string{"5432/tcp"},
		WaitingFor:   wait.ForListeningPort("5432/tcp"),
		Env: map[string]string{
			"POSTGRES_DB":       "testdb",
			"POSTGRES_USER":     "testuser",
			"POSTGRES_PASSWORD": "testpassword",
		},
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to start container: %w", err)
	}

	hostIP, err := container.Host(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get host IP: %w", err)
	}
	port, err := container.MappedPort(ctx, "5432")
	if err != nil {
		return nil, fmt.Errorf("failed to get mapped port: %w", err)
	}

	connStr := fmt.Sprintf("host=%s port=%s user=testuser password=testpassword dbname=testdb sslmode=disable",
		hostIP, port.Port())

	// Ping DB to ensure connection is ready
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open db connection: %w", err)
	}
	if err = db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping db: %w", err)
	}
	db.Close() // Close initial ping connection

	return &PostgresContainer{Container: container, ConnectionString: connStr}, nil
}

func TeardownPostgres(ctx context.Context, t *testing.T, container *PostgresContainer) {
	if err := container.Terminate(ctx); err != nil {
		t.Fatalf("failed to terminate container: %v", err)
	}
}
