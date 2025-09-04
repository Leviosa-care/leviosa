package helpers

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/Leviosa-care/core/contracts/settings"
	"github.com/Leviosa-care/settings/internal/domain"
	th "github.com/Leviosa-care/settings/test/helpers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TEST=TestGetCompanyInstagram make test-integration-test

func TestGetCompanyInstagram(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("should return 404 when company instagram not set", func(t *testing.T) {
		th.ClearAllTestData(t, ctx, testPool)

		req := th.NewGetCompanyInstagramRequest(t, ctx, testServerURL)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		var respBody struct {
			Error string `json:"error"`
		}
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		require.NoError(t, err)
		assert.Contains(t, respBody.Error, "company instagram")
	})

	t.Run("should successfully retrieve company instagram", func(t *testing.T) {
		th.ClearAllTestData(t, ctx, testPool)

		// Setup: Insert company instagram directly into database
		insta := "https://instagram.com/testcompany"
		th.InsertCompanyInstagram(t, ctx, insta, testPool)

		// Test: Get the company instagram
		req := th.NewGetCompanyInstagramRequest(t, ctx, testServerURL)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var respBody domain.GetCompanyInstagramResponse
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		require.NoError(t, err)
		assert.Equal(t, insta, respBody.Instagram)
	})
}

// TEST=TestSetCompanyInstagram make test-integration-test

func TestSetCompanyInstagram(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("should successfully set company instagram", func(t *testing.T) {
		th.ClearAllTestData(t, ctx, testPool)

		// Create a test channel for RabbitMQ verification
		testCh := th.GetRabbitMQChannel(t, testMQConn)
		defer testCh.Close()

		// Purge queues to ensure clean state
		th.PurgeSettingsQueues(t, testCh)

		insta := "https://instagram.com/mycompany"
		request := domain.SetCompanyInstagramRequest{Instagram: insta}
		req := th.NewSetCompanyInstagramRequest(t, ctx, testServerURL, request)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var respBody domain.SetCompanyInstagramResponse
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		require.NoError(t, err)
		assert.True(t, respBody.Success)

		// Verify data was persisted directly in database
		instagram, err := th.GetCompanyInstagramFromDB(t, ctx, testPool)
		require.NoError(t, err)
		assert.Equal(t, insta, instagram)

		// Verify RabbitMQ message was published
		th.VerifySettingsUpdateMessage(t, testCh, settings.CompanyInstagram, insta)
	})

	t.Run("should return 400 for empty instagram link", func(t *testing.T) {
		th.ClearAllTestData(t, ctx, testPool)

		request := domain.SetCompanyInstagramRequest{Instagram: ""}
		req := th.NewSetCompanyInstagramRequest(t, ctx, testServerURL, request)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var respBody struct {
			Error string `json:"error"`
		}
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		require.NoError(t, err)
		assert.Contains(t, respBody.Error, "instagram_required")
	})

	t.Run("should return 400 for invalid URL format", func(t *testing.T) {
		th.ClearAllTestData(t, ctx, testPool)

		invalidUrls := []string{
			"not-a-url",
			"just-text",
			"instagram.com/user",       // missing protocol
			"ftp://instagram.com/user", // wrong protocol but valid URL
			"http://",                  // incomplete URL
			"https://",                 // incomplete URL
		}

		for _, url := range invalidUrls {
			t.Run("invalid URL: "+url, func(t *testing.T) {
				request := domain.SetCompanyInstagramRequest{Instagram: url}
				req := th.NewSetCompanyInstagramRequest(t, ctx, testServerURL, request)

				resp, err := client.Do(req)
				require.NoError(t, err)
				defer resp.Body.Close()

				assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

				var respBody struct {
					Error string `json:"error"`
				}
				err = json.NewDecoder(resp.Body).Decode(&respBody)
				require.NoError(t, err)
				assert.Contains(t, respBody.Error, "instagram_format")
			})
		}
	})

	t.Run("should return 400 for instagram link exceeding 255 characters", func(t *testing.T) {
		th.ClearAllTestData(t, ctx, testPool)

		// Create a long URL
		longPath := strings.Repeat("a", 240)
		longUrl := "https://instagram.com/" + longPath // Total > 255 chars

		request := domain.SetCompanyInstagramRequest{Instagram: longUrl}
		req := th.NewSetCompanyInstagramRequest(t, ctx, testServerURL, request)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var respBody struct {
			Error string `json:"error"`
		}
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		require.NoError(t, err)
		assert.Contains(t, respBody.Error, "instagram_length")
	})

	t.Run("should successfully accept various valid URL formats", func(t *testing.T) {
		th.ClearAllTestData(t, ctx, testPool)

		validUrls := []string{
			"https://instagram.com/company",
			"https://www.instagram.com/company",
			"http://instagram.com/company",
			"https://instagram.com/company_name",
			"https://instagram.com/company.official",
			"https://instagram.com/company123",
			"https://instagram.com/p/ABC123/",        // specific post
			"https://instagram.com/company/?hl=en",   // with query params
			"https://m.instagram.com/company",        // mobile version
			"https://business.instagram.com/company", // business version
		}

		for _, url := range validUrls {
			t.Run("valid URL: "+url, func(t *testing.T) {
				th.ClearAllTestData(t, ctx, testPool)

				// Create a test channel for RabbitMQ verification
				testCh := th.GetRabbitMQChannel(t, testMQConn)
				defer testCh.Close()

				// Purge queues to ensure clean state
				th.PurgeSettingsQueues(t, testCh)

				request := domain.SetCompanyInstagramRequest{Instagram: url}
				req := th.NewSetCompanyInstagramRequest(t, ctx, testServerURL, request)

				resp, err := client.Do(req)
				require.NoError(t, err)
				defer resp.Body.Close()

				assert.Equal(t, http.StatusOK, resp.StatusCode)

				var respBody domain.SetCompanyInstagramResponse
				err = json.NewDecoder(resp.Body).Decode(&respBody)
				require.NoError(t, err)
				assert.True(t, respBody.Success)

				// Verify RabbitMQ message was published
				th.VerifySettingsUpdateMessage(t, testCh, settings.CompanyInstagram, url)
			})
		}
	})

	t.Run("should accept non-Instagram URLs (business requirement)", func(t *testing.T) {
		th.ClearAllTestData(t, ctx, testPool)

		// Test that the field accepts any valid URL, not just Instagram
		nonInstagramUrls := []string{
			"https://twitter.com/company",
			"https://facebook.com/company",
			"https://linkedin.com/company/company-name",
			"https://company.com/social",
		}

		for _, url := range nonInstagramUrls {
			t.Run("non-Instagram URL: "+url, func(t *testing.T) {
				th.ClearAllTestData(t, ctx, testPool)

				// Create a test channel for RabbitMQ verification
				testCh := th.GetRabbitMQChannel(t, testMQConn)
				defer testCh.Close()

				// Purge queues to ensure clean state
				th.PurgeSettingsQueues(t, testCh)

				request := domain.SetCompanyInstagramRequest{Instagram: url}
				req := th.NewSetCompanyInstagramRequest(t, ctx, testServerURL, request)

				resp, err := client.Do(req)
				require.NoError(t, err)
				defer resp.Body.Close()

				assert.Equal(t, http.StatusOK, resp.StatusCode)

				var respBody domain.SetCompanyInstagramResponse
				err = json.NewDecoder(resp.Body).Decode(&respBody)
				require.NoError(t, err)
				assert.True(t, respBody.Success)

				// Verify RabbitMQ message was published
				th.VerifySettingsUpdateMessage(t, testCh, settings.CompanyInstagram, url)
			})
		}
	})

	t.Run("should return 415 for incorrect content type", func(t *testing.T) {
		th.ClearAllTestData(t, ctx, testPool)

		request := domain.SetCompanyInstagramRequest{Instagram: "https://instagram.com/test"}
		req := th.NewSetCompanyInstagramRequest(t, ctx, testServerURL, request)
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

	t.Run("should return 400 for unknown JSON fields", func(t *testing.T) {
		th.ClearAllTestData(t, ctx, testPool)

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, testServerURL+"/admin/settings/instagram",
			strings.NewReader(`{"instagram": "https://instagram.com/test", "unknown_field": "value"}`))
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

	t.Run("should successfully update existing company instagram", func(t *testing.T) {
		th.ClearAllTestData(t, ctx, testPool)

		// Create a test channel for RabbitMQ verification
		testCh := th.GetRabbitMQChannel(t, testMQConn)
		defer testCh.Close()

		// Purge queues to ensure clean state
		th.PurgeSettingsQueues(t, testCh)

		// Set initial instagram

		oldInsta := "https://instagram.com/oldcompany"
		th.InsertCompanyInstagram(t, ctx, oldInsta, testPool)

		// Update to new instagram
		newInsta := "https://instagram.com/newcompany"
		request2 := domain.SetCompanyInstagramRequest{Instagram: newInsta}
		req2 := th.NewSetCompanyInstagramRequest(t, ctx, testServerURL, request2)
		resp2, err := client.Do(req2)
		require.NoError(t, err)
		defer resp2.Body.Close()
		require.Equal(t, http.StatusOK, resp2.StatusCode)

		// Verify updated instagram directly in database
		instagram, err := th.GetCompanyInstagramFromDB(t, ctx, testPool)
		require.NoError(t, err)
		assert.Equal(t, newInsta, instagram)

		// Verify RabbitMQ message was published
		th.VerifySettingsUpdateMessage(t, testCh, settings.CompanyInstagram, newInsta)
	})
}
