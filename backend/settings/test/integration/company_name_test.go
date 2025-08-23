package testdata

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/Leviosa-care/core/contracts/settings"
	"github.com/Leviosa-care/settings/internal/domain"
	td "github.com/Leviosa-care/settings/test/testdata"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetCompanyName(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("should return 404 when company name not set", func(t *testing.T) {
		td.ClearSettingsTable(t, ctx, testPool)

		req := td.NewGetCompanyNameRequest(t, ctx, testServerURL)
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
		td.ClearSettingsTable(t, ctx, testPool)

		// Setup: Insert company name directly into database
		td.InsertCompanyName(t, ctx, "Test Company Inc", testPool)

		// Test: Get the company name via HTTP
		req := td.NewGetCompanyNameRequest(t, ctx, testServerURL)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var respBody domain.GetCompanyNameResponse
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		require.NoError(t, err)
		assert.Equal(t, "Test Company Inc", respBody.Name)
	})
}

func TestSetCompanyName(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("should successfully set company name", func(t *testing.T) {
		td.ClearSettingsTable(t, ctx, testPool)

		// Create a test channel for RabbitMQ verification
		testCh := td.GetRabbitMQChannel(t, testMQConn)
		defer testCh.Close()

		// Purge queues to ensure clean state
		td.PurgeSettingsQueues(t, testCh)

		request := domain.SetCompanyNameRequest{Name: "New Company Name"}
		req := td.NewSetCompanyNameRequest(t, ctx, testServerURL, request)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var respBody domain.SetCompanyNameResponse
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		require.NoError(t, err)
		assert.True(t, respBody.Success)

		// Verify data was persisted directly in database
		name, err := td.GetCompanyNameFromDB(t, ctx, testPool)
		require.NoError(t, err)
		assert.Equal(t, "New Company Name", name)

		// Verify RabbitMQ message was published
		td.VerifySettingsUpdateMessage(t, testCh, settings.CompanyName, "New Company Name")
	})

	t.Run("should return 400 for empty company name", func(t *testing.T) {
		td.ClearSettingsTable(t, ctx, testPool)

		request := domain.SetCompanyNameRequest{Name: ""}
		req := td.NewSetCompanyNameRequest(t, ctx, testServerURL, request)

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
		td.ClearSettingsTable(t, ctx, testPool)

		request := domain.SetCompanyNameRequest{Name: "   "}
		req := td.NewSetCompanyNameRequest(t, ctx, testServerURL, request)

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
		td.ClearSettingsTable(t, ctx, testPool)

		longName := string(make([]byte, 256))
		for i := range longName {
			longName = longName[:i] + "A" + longName[i+1:]
		}

		request := domain.SetCompanyNameRequest{Name: longName}
		req := td.NewSetCompanyNameRequest(t, ctx, testServerURL, request)

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
		td.ClearSettingsTable(t, ctx, testPool)

		request := domain.SetCompanyNameRequest{Name: "Test Company"}
		req := td.NewSetCompanyNameRequest(t, ctx, testServerURL, request)
		req.Header.Set("Content-Type", "text/plain")

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
		td.ClearSettingsTable(t, ctx, testPool)

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, testServerURL+"/admin/settings/name",
			strings.NewReader(`{"name": "test", "invalid_field": "value"}`))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

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
		td.ClearSettingsTable(t, ctx, testPool)

		// Set initial name
		request1 := domain.SetCompanyNameRequest{Name: "Initial Company"}
		req1 := td.NewSetCompanyNameRequest(t, ctx, testServerURL, request1)
		resp1, err := client.Do(req1)
		require.NoError(t, err)
		defer resp1.Body.Close()
		require.Equal(t, http.StatusOK, resp1.StatusCode)

		// Update to new name
		request2 := domain.SetCompanyNameRequest{Name: "Updated Company"}
		req2 := td.NewSetCompanyNameRequest(t, ctx, testServerURL, request2)
		resp2, err := client.Do(req2)
		require.NoError(t, err)
		defer resp2.Body.Close()
		require.Equal(t, http.StatusOK, resp2.StatusCode)

		// Verify updated name directly in database
		name, err := td.GetCompanyNameFromDB(t, ctx, testPool)
		require.NoError(t, err)
		assert.Equal(t, "Updated Company", name)
	})
}

