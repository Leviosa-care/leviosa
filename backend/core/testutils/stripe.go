// testutils/stripe_mock.go
package testutils

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// StripeMockContainer holds the testcontainer.Container and the URL to the mock server.
type StripeMockContainer struct {
	testcontainers.Container
	URL string
}

// SetupStripeMock starts a stripe-mock Docker container.
// It will expose the mock Stripe API on a mapped port.
func SetupStripeMock(ctx context.Context, t *testing.T) (*StripeMockContainer, error) {
	req := testcontainers.ContainerRequest{
		Image:        "stripe/stripe-mock:latest", // Or a specific version like 'v0.124.0'
		ExposedPorts: []string{"12111/tcp"},       // Default stripe-mock port
		WaitingFor:   wait.ForListeningPort("12111/tcp").WithStartupTimeout(2 * time.Minute),
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		if t != nil {
			t.Logf("Failed to start stripe-mock container: %v", err)
		}
		return nil, fmt.Errorf("failed to start stripe-mock container: %w", err)
	}

	hostIP, err := container.Host(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get host IP for stripe-mock: %w", err)
	}
	port, err := container.MappedPort(ctx, "12111")
	if err != nil {
		return nil, fmt.Errorf("failed to get mapped port for stripe-mock: %w", err)
	}

	mockURL := fmt.Sprintf("http://%s:%s", hostIP, port.Port())
	if t != nil {
		t.Logf("Stripe-mock started at %s", mockURL)
	}

	// You might want to add a health check/ping to the mock server here if it has one.
	// For stripe-mock, waiting for the port is usually sufficient.

	return &StripeMockContainer{Container: container, URL: mockURL}, nil
}

// TeardownStripeMock terminates the stripe-mock Docker container.
func TeardownStripeMock(ctx context.Context, t *testing.T, container *StripeMockContainer) {
	if container == nil {
		return
	}
	if err := container.Terminate(ctx); err != nil {
		if t != nil {
			t.Logf("Failed to terminate stripe-mock container: %v", err)
		}
	}
}
