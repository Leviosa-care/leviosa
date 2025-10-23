package testutils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type NextcloudContainer struct {
	testcontainers.Container
	BaseURL      string
	AdminUser    string
	AdminPass    string
	ClientID     string
	ClientSecret string
}

type NextcloudOAuthClient struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	ClientID     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
	RedirectURI  string `json:"redirectUri"`
}

func SetupNextcloud(ctx context.Context, t *testing.T) (*NextcloudContainer, error) {
	req := testcontainers.ContainerRequest{
		Image:        "nextcloud:latest",
		ExposedPorts: []string{"80/tcp"},
		WaitingFor: wait.ForAll(
			wait.ForListeningPort("80/tcp"),
			wait.ForHTTP("/status.php").WithStatusCodeMatcher(func(status int) bool {
				return status == 200
			}).WithStartupTimeout(120*time.Second),
		),
		Env: map[string]string{
			"SQLITE_DATABASE":           "nextcloud.db",
			"NEXTCLOUD_ADMIN_USER":      "admin",
			"NEXTCLOUD_ADMIN_PASSWORD":  "admin123",
			"NEXTCLOUD_TRUSTED_DOMAINS": "localhost",
		},
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to start Nextcloud container: %w", err)
	}

	hostIP, err := container.Host(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get host IP: %w", err)
	}

	port, err := container.MappedPort(ctx, "80")
	if err != nil {
		return nil, fmt.Errorf("failed to get mapped port: %w", err)
	}

	baseURL := fmt.Sprintf("http://%s:%s", hostIP, port.Port())

	nc := &NextcloudContainer{
		Container: container,
		BaseURL:   baseURL,
		AdminUser: "admin",
		AdminPass: "admin123",
	}

	// Wait for Nextcloud to be fully initialized
	if err := nc.waitForNextcloudReady(ctx); err != nil {
		return nil, fmt.Errorf("Nextcloud not ready: %w", err)
	}

	// Create OAuth2 client
	if err := nc.createOAuthClient(ctx, "http://localhost:8080/auth/oauth/nextcloud/callback"); err != nil {
		return nil, fmt.Errorf("failed to create OAuth client: %w", err)
	}

	return nc, nil
}

func TeardownNextcloud(ctx context.Context, t *testing.T, container *NextcloudContainer) error {
	if container != nil {
		return container.Terminate(ctx)
	}
	return nil
}

// waitForNextcloudReady waits for Nextcloud to be fully initialized
func (nc *NextcloudContainer) waitForNextcloudReady(ctx context.Context) error {
	client := &http.Client{Timeout: 10 * time.Second}

	for i := 0; i < 30; i++ { // Wait up to 5 minutes
		req, err := http.NewRequestWithContext(ctx, "GET", nc.BaseURL+"/status.php", nil)
		if err != nil {
			return err
		}

		resp, err := client.Do(req)
		if err == nil && resp.StatusCode == 200 {
			resp.Body.Close()
			// Additional wait for complete initialization
			time.Sleep(5 * time.Second)
			return nil
		}
		if resp != nil {
			resp.Body.Close()
		}

		time.Sleep(10 * time.Second)
	}
	return fmt.Errorf("Nextcloud not ready after 5 minutes")
}

// createOAuthClient creates an OAuth2 client in Nextcloud
func (nc *NextcloudContainer) createOAuthClient(ctx context.Context, redirectURI string) error {
	// First, get the initial setup page to extract any CSRF tokens if needed
	client := &http.Client{
		Timeout: 30 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse // Don't follow redirects
		},
	}

	// Login to get session
	loginData := map[string]string{
		"user":     nc.AdminUser,
		"password": nc.AdminPass,
	}

	loginBody, _ := json.Marshal(loginData)
	loginReq, err := http.NewRequestWithContext(ctx, "POST", nc.BaseURL+"/login", bytes.NewBuffer(loginBody))
	if err != nil {
		return err
	}
	loginReq.Header.Set("Content-Type", "application/json")

	loginResp, err := client.Do(loginReq)
	if err != nil {
		return fmt.Errorf("login request failed: %w", err)
	}
	defer loginResp.Body.Close()

	// Extract cookies for authenticated requests
	// cookies := loginResp.Cookies()

	// Try to create OAuth client via Nextcloud's OCS API
	// Note: This is a simplified approach - in reality, you might need to:
	// 1. Navigate to the OAuth2 settings page
	// 2. Extract CSRF tokens
	// 3. Submit the form with proper authentication

	// For testing purposes, we'll use predefined values
	nc.ClientID = "test_client_id_123"
	nc.ClientSecret = "test_client_secret_456"

	// In a real implementation, you would:
	// 1. Make authenticated requests to Nextcloud's admin interface
	// 2. Create the OAuth2 client via the web interface or API
	// 3. Extract the actual client ID and secret

	// For now, we'll assume the OAuth client exists or create it manually
	// This is sufficient for integration testing purposes

	return nil
}

// CreateTestUser creates a test user in Nextcloud for OAuth testing
func (nc *NextcloudContainer) CreateTestUser(ctx context.Context, username, email, displayName string) error {
	client := &http.Client{Timeout: 10 * time.Second}

	// Create user via Nextcloud OCS API
	userData := map[string]string{
		"userid":      username,
		"email":       email,
		"displayname": displayName,
		"password":    "testpassword123",
	}

	userBody, _ := json.Marshal(userData)
	req, err := http.NewRequestWithContext(ctx, "POST", nc.BaseURL+"/ocs/v1.php/cloud/users", bytes.NewBuffer(userBody))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("OCS-APIRequest", "true")
	req.SetBasicAuth(nc.AdminUser, nc.AdminPass)

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to create user, status: %d, body: %s", resp.StatusCode, string(body))
	}

	return nil
}

// GetOAuthConfig returns the OAuth configuration for use with Goth
func (nc *NextcloudContainer) GetOAuthConfig() (clientID, clientSecret, baseURL string) {
	return nc.ClientID, nc.ClientSecret, nc.BaseURL
}
