package helpers

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/common/contracts/settings"
	tu "github.com/Leviosa-care/leviosa/backend/internal/common/testutils"
	httpEndpoints "github.com/Leviosa-care/settings/internal/adapters/http"
	"github.com/Leviosa-care/settings/internal/domain"
	th "github.com/Leviosa-care/settings/test/helpers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TEST=TestGetCompanyName make test-integration-test

func TestGetCompanyName(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("should return 404 when company name not set", func(t *testing.T) {
		th.ClearSettingsTable(t, ctx, testPool)

		req := th.NewGetCompanyNameRequest(t, ctx, testServerURL)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		var respBody struct {
			Error string `json:"error"`
		}
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		require.NoError(t, err)
		assert.Contains(t, respBody.Error, "company name")
	})

	t.Run("should successfully retrieve company name", func(t *testing.T) {
		th.ClearSettingsTable(t, ctx, testPool)

		// Setup: Insert company name directly into database
		name := "Test Company Inc"
		th.InsertCompanyName(t, ctx, name, testPool)

		// Test: Get the company name via HTTP
		req := th.NewGetCompanyNameRequest(t, ctx, testServerURL)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var respBody domain.GetCompanyNameResponse
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		require.NoError(t, err)
		assert.Equal(t, name, respBody.Name)
	})
}

// TEST=TestSetCompanyName make test-integration-test

func TestSetCompanyName(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("should successfully set company name", func(t *testing.T) {
		th.ClearSettingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Create a test channel for RabbitMQ verification
		testCh := th.GetRabbitMQChannel(t, testMQConn)
		defer testCh.Close()

		// Purge queues to ensure clean state
		th.PurgeSettingsQueues(t, testCh)

		// Setup admin user and create authenticated request
		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		name := "New Company Name"
		request := domain.SetCompanyNameRequest{Name: name}
		req := th.NewSetCompanyNameRequest(t, ctx, testServerURL, request)

		// Add authentication to the request
		authHeader := tu.CreateAuthHeader(accessToken)
		for key, values := range authHeader {
			for _, value := range values {
				req.Header.Add(key, value)
			}
		}
		req.AddCookie(tu.CreateAuthCookie(accessToken))

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var respBody domain.SetCompanyNameResponse
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		require.NoError(t, err)
		assert.True(t, respBody.Success)

		// Verify data was persisted directly in database
		retrievedName, err := th.GetCompanyNameFromDB(t, ctx, testPool)
		require.NoError(t, err)
		assert.Equal(t, name, retrievedName)

		// Verify RabbitMQ message was published
		th.VerifySettingsUpdateMessage(t, testCh, settings.CompanyName, name)
	})

	t.Run("should return 400 for empty company name", func(t *testing.T) {
		th.ClearSettingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup admin user and create authenticated request
		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		request := domain.SetCompanyNameRequest{Name: ""}
		req := th.NewSetCompanyNameRequest(t, ctx, testServerURL, request)

		// Add authentication to the request
		authHeader := tu.CreateAuthHeader(accessToken)
		for key, values := range authHeader {
			for _, value := range values {
				req.Header.Add(key, value)
			}
		}
		req.AddCookie(tu.CreateAuthCookie(accessToken))

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var respBody struct {
			Error string `json:"error"`
		}
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		require.NoError(t, err)
		assert.Contains(t, respBody.Error, "name_required")
	})

	t.Run("should return 400 for whitespace-only company name", func(t *testing.T) {
		th.ClearSettingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)
		request := domain.SetCompanyNameRequest{Name: "   "}
		req := th.NewSetCompanyNameRequest(t, ctx, testServerURL, request)

		// Add authentication to the request
		authHeader := tu.CreateAuthHeader(accessToken)
		for key, values := range authHeader {
			for _, value := range values {
				req.Header.Add(key, value)
			}
		}
		req.AddCookie(tu.CreateAuthCookie(accessToken))

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var respBody struct {
			Error string `json:"error"`
		}
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		require.NoError(t, err)
		assert.Contains(t, respBody.Error, "name_empty")
	})

	t.Run("should return 400 for company name exceeding 255 characters", func(t *testing.T) {
		th.ClearSettingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		longName := string(make([]byte, 256))
		for i := range longName {
			longName = longName[:i] + "A" + longName[i+1:]
		}

		request := domain.SetCompanyNameRequest{Name: longName}
		req := th.NewSetCompanyNameRequest(t, ctx, testServerURL, request)

		// Add authentication to the request
		authHeader := tu.CreateAuthHeader(accessToken)
		for key, values := range authHeader {
			for _, value := range values {
				req.Header.Add(key, value)
			}
		}
		req.AddCookie(tu.CreateAuthCookie(accessToken))

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var respBody struct {
			Error string `json:"error"`
		}
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		require.NoError(t, err)
		assert.Contains(t, respBody.Error, "name_length")
	})

	t.Run("should return 415 for incorrect content type", func(t *testing.T) {
		th.ClearSettingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)
		request := domain.SetCompanyNameRequest{Name: "Test Company"}
		req := th.NewSetCompanyNameRequest(t, ctx, testServerURL, request)
		req.Header.Set("Content-Type", "text/plain")

		// Add authentication to the request
		authHeader := tu.CreateAuthHeader(accessToken)
		for key, values := range authHeader {
			for _, value := range values {
				req.Header.Add(key, value)
			}
		}
		req.AddCookie(tu.CreateAuthCookie(accessToken))

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnsupportedMediaType, resp.StatusCode)

		var respBody struct {
			Error string `json:"error"`
		}
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		require.NoError(t, err)
		assert.Contains(t, respBody.Error, "unsupported media type")
	})

	t.Run("should return 400 for invalid JSON", func(t *testing.T) {
		th.ClearSettingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, testServerURL+httpEndpoints.SetCompanyNameEndpoint,
			strings.NewReader(`{"name": "test", "invalid_field": "value"}`))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		// Add authentication to the request
		authHeader := tu.CreateAuthHeader(accessToken)
		for key, values := range authHeader {
			for _, value := range values {
				req.Header.Add(key, value)
			}
		}
		req.AddCookie(tu.CreateAuthCookie(accessToken))

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var respBody struct {
			Error string `json:"error"`
		}
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		require.NoError(t, err)
		assert.Contains(t, respBody.Error, "invalid request body")
	})

	t.Run("should successfully update existing company name", func(t *testing.T) {
		th.ClearSettingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Create a test channel for RabbitMQ verification
		testCh := th.GetRabbitMQChannel(t, testMQConn)
		defer testCh.Close()

		// Purge queues to ensure clean state
		th.PurgeSettingsQueues(t, testCh)

		// Set initial name
		oldName := "Initial Company"
		th.InsertCompanyName(t, ctx, oldName, testPool)

		// Setup admin user and update to new name
		accessToken := tu.SetupAdminUser(t, ctx, authCtx)
		newName := "Updated Company"
		request := domain.SetCompanyNameRequest{Name: newName}
		req := th.NewSetCompanyNameRequest(t, ctx, testServerURL, request)

		// Add authentication to the request
		authHeader := tu.CreateAuthHeader(accessToken)
		for key, values := range authHeader {
			for _, value := range values {
				req.Header.Add(key, value)
			}
		}
		req.AddCookie(tu.CreateAuthCookie(accessToken))

		resp2, err := client.Do(req)
		require.NoError(t, err)
		defer resp2.Body.Close()
		require.Equal(t, http.StatusOK, resp2.StatusCode)

		// Verify updated name directly in database
		name, err := th.GetCompanyNameFromDB(t, ctx, testPool)
		require.NoError(t, err)
		assert.Equal(t, newName, name)

		// Verify RabbitMQ message was published
		th.VerifySettingsUpdateMessage(t, testCh, settings.CompanyName, newName)
	})
}
