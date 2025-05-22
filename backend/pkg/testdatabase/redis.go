package testdatabase

import (
	"context"
	"fmt"

	"github.com/hengadev/leviosa/pkg/redisutil"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type TestDatabase struct {
	container testcontainers.Container
	DB        *redis.Client
}

func NewTestDatabase(ctx context.Context) (*TestDatabase, error) {
	// setup db container
	container, dbConn, err := createContainer(ctx)
	if err != nil {
		return nil, fmt.Errorf("createContainer: %w", err)
	}
	return &TestDatabase{
		container: container,
		DB:        dbConn,
	}, nil
}

func (tdb *TestDatabase) TearDown() {
	tdb.DB.Close()
	// remove test container
	_ = tdb.container.Terminate(context.Background())
}

func createContainer(ctx context.Context) (testcontainers.Container, *redis.Client, error) {
	req := testcontainers.ContainerRequest{
		Name:         fmt.Sprintf("redis-test-%s", uuid.NewString()),
		Image:        "redis:7.4",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForLog("database system is ready to accept connections"),
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return container, nil, fmt.Errorf("failed to start container: %v", err)
	}

	db, err := redisutil.Connect(ctx, redisutil.WithAddr("localhost:6379"))
	if err != nil {
		return container, db, fmt.Errorf("failed to establish database connection: %v", err)
	}

	return container, db, nil
}
