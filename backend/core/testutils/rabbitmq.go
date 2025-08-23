package testutils

import (
	"context"
	"fmt"
	"testing"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type RabbitMQContainer struct {
	testcontainers.Container
	ConnectionString string
	ManagementURL    string
}

func SetupRabbitMQ(ctx context.Context, t *testing.T) (*RabbitMQContainer, error) {
	req := testcontainers.ContainerRequest{
		Image:        "rabbitmq:3-management-alpine",
		ExposedPorts: []string{"5672/tcp", "15672/tcp"}, // AMQP and Management UI ports
		WaitingFor: wait.ForAll(
			wait.ForListeningPort("5672/tcp"),
			wait.ForListeningPort("15672/tcp"),
			wait.ForLog("Server startup complete"),
		).WithDeadline(60 * time.Second),
		Env: map[string]string{
			"RABBITMQ_DEFAULT_USER": "testuser",
			"RABBITMQ_DEFAULT_PASS": "testpassword",
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

	// Get AMQP port
	amqpPort, err := container.MappedPort(ctx, "5672")
	if err != nil {
		return nil, fmt.Errorf("failed to get AMQP mapped port: %w", err)
	}

	// Get Management UI port
	mgmtPort, err := container.MappedPort(ctx, "15672")
	if err != nil {
		return nil, fmt.Errorf("failed to get management mapped port: %w", err)
	}

	connStr := fmt.Sprintf("amqp://testuser:testpassword@%s:%s/",
		hostIP, amqpPort.Port())

	mgmtURL := fmt.Sprintf("http://testuser:testpassword@%s:%s",
		hostIP, mgmtPort.Port())

	// Test connection to ensure RabbitMQ is ready
	conn, err := amqp.Dial(connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}
	conn.Close() // Close initial test connection

	return &RabbitMQContainer{
		Container:        container,
		ConnectionString: connStr,
		ManagementURL:    mgmtURL,
	}, nil
}

func TeardownRabbitMQ(ctx context.Context, t *testing.T, container *RabbitMQContainer) {
	if err := container.Terminate(ctx); err != nil {
		t.Fatalf("failed to terminate container: %v", err)
	}
}

// Helper function to create a connection for tests
func (r *RabbitMQContainer) NewConnection() (*amqp.Connection, error) {
	return amqp.Dial(r.ConnectionString)
}

// Helper function to create a channel for tests
func (r *RabbitMQContainer) NewChannel() (*amqp.Channel, *amqp.Connection, error) {
	conn, err := r.NewConnection()
	if err != nil {
		return nil, nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, nil, err
	}

	return ch, conn, nil
}
