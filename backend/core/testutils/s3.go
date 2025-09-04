// testutils/localstack.go
package testutils

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// LocalstackContainer holds the testcontainer.Container and endpoint URLs.
type LocalstackContainer struct {
	testcontainers.Container
	S3Endpoint string
	// You can add other service endpoints if needed, e.g., SQSEndpoint, LambdaEndpoint
}

// Enhanced SetupLocalstack with better debugging and wait strategies
func SetupLocalstack(ctx context.Context, t *testing.T) (*LocalstackContainer, error) {
	// Add timeout to the context to prevent indefinite hanging
	setupCtx, cancel := context.WithTimeout(ctx, 3*time.Minute)
	defer cancel()

	if t != nil {
		t.Log("Starting Localstack container setup...")
	} else {
		log.Println("Starting Localstack container setup...")
	}

	req := testcontainers.ContainerRequest{
		Image:        "localstack/localstack:4.6", // Use a specific stable version instead of latest
		ExposedPorts: []string{"4566/tcp"},
		Env: map[string]string{
			"SERVICES":              "s3",
			"AWS_DEFAULT_REGION":    "us-east-1",
			"AWS_ACCESS_KEY_ID":     "test",
			"AWS_SECRET_ACCESS_KEY": "test",
			"DEBUG":                 "1",     // Enable debug logging
			"LS_LOG":                "trace", // More verbose logging
			"PERSISTENCE":           "0",     // Disable persistence for faster startup
		},
		// Use multiple wait strategies as fallback
		WaitingFor: wait.ForAll(
			wait.ForLog("Ready.").WithStartupTimeout(2 * time.Minute),
			// wait.ForHTTP("/health").OnPort("4566").WithStartupTimeout(2*time.Minute),
		).WithDeadline(2 * time.Minute),
	}

	if t != nil {
		t.Log("Creating container...")
	} else {
		log.Println("Creating container...")
	}

	container, err := testcontainers.GenericContainer(setupCtx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		if t != nil {
			t.Logf("Failed to start Localstack container: %v", err)
		} else {
			log.Printf("Failed to start Localstack container: %v", err)
		}
		return nil, fmt.Errorf("failed to start Localstack container: %w", err)
	}

	if t != nil {
		t.Log("Container started, getting connection details...")
	} else {
		log.Println("Container started, getting connection details...")
	}

	hostIP, err := container.Host(setupCtx)
	if err != nil {
		return nil, fmt.Errorf("failed to get host IP for Localstack: %w", err)
	}

	port, err := container.MappedPort(setupCtx, "4566")
	if err != nil {
		return nil, fmt.Errorf("failed to get mapped port for Localstack: %w", err)
	}

	s3Endpoint := fmt.Sprintf("http://%s:%s", hostIP, port.Port())

	if t != nil {
		t.Logf("Localstack S3 endpoint: %s", s3Endpoint)
	} else {
		log.Printf("Localstack S3 endpoint: %s", s3Endpoint)
	}

	return &LocalstackContainer{Container: container, S3Endpoint: s3Endpoint}, nil
}

// TeardownLocalstack terminates the Localstack Docker container.
func TeardownLocalstack(ctx context.Context, t *testing.T, container *LocalstackContainer) {
	if container == nil {
		return
	}
	if err := container.Terminate(ctx); err != nil {
		if t != nil {
			t.Logf("Failed to terminate Localstack container: %v", err)
		}
	}
}

func NewLocalstackS3Config(s3Endpoint string) aws.Config {
	// Load default config first
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("us-east-1"), // Or the region you set for Localstack
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("test", "test", "")),
	)
	if err != nil {
		panic(fmt.Sprintf("failed to load AWS config for Localstack: %v", err))
	}

	return cfg
}
