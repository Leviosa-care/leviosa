package helpers

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"
	"time"

	// "github.com/Leviosa-care/leviosa/backend/internal/common/contracts/settings"
	tu "github.com/Leviosa-care/leviosa/backend/internal/common/testutils"
	"github.com/Leviosa-care/leviosa/backend/internal/settings/domain"
	httpEndpoints "github.com/Leviosa-care/leviosa/backend/internal/settings/interface/http"
	th "github.com/Leviosa-care/leviosa/backend/test/helpers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestGetCompanyInstagram TEST_PATH=test/integration/settings/company_instagram_test.go

func TestGetCompanyInstagram(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("should return 404 when company instagram not set", func(t *testing.T) {
		th.ClearSettingTestData(t, ctx, testPool, s3Client)

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
		th.ClearSettingTestData(t, ctx, testPool, s3Client)

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

// make test-func TEST_NAME=TestSetCompanyInstagram TEST_PATH=test/integration/settings/company_instagram_test.go

func TestSetCompanyInstagram(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("should successfully set company instagram", func(t *testing.T) {
		th.ClearSettingTestData(t, ctx, testPool, s3Client)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Create a test channel for RabbitMQ verification
		// COMMENTED OUT: RabbitMQ verification disabled
		// testCh := th.GetRabbitMQChannel(t, testMQConn)
		// defer testCh.Close()

		// Purge queues to ensure clean state
		// th.PurgeSettingsQueues(t, testCh)

		insta := "https://instagram.com/mycompany"

		// Setup admin user and create authenticated request
		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		request := domain.SetCompanyInstagramRequest{Instagram: insta}
		req := th.NewSetCompanyInstagramRequest(t, ctx, testServerURL, request)

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

		var respBody domain.SetCompanyInstagramResponse
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		require.NoError(t, err)
		assert.True(t, respBody.Success)

		// Verify data was persisted directly in database
		instagram, err := th.GetCompanyInstagram(t, ctx, testPool)
		require.NoError(t, err)
		assert.Equal(t, insta, instagram)

		// COMMENTED OUT: Verify RabbitMQ message was published (verification disabled)
		// COMMENTED OUT: RabbitMQ verification disabled
		// th.VerifySettingsUpdateMessage(t, testCh, settings.CompanyInstagram, insta)
	})

	t.Run("should return 400 for empty instagram link", func(t *testing.T) {
		th.ClearSettingTestData(t, ctx, testPool, s3Client)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup admin user and create authenticated request
		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		request := domain.SetCompanyInstagramRequest{Instagram: ""}
		req := th.NewSetCompanyInstagramRequest(t, ctx, testServerURL, request)

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
		assert.Contains(t, respBody.Error, "instagram_required")
	})

	t.Run("should return 400 for invalid URL format", func(t *testing.T) {
		th.ClearSettingTestData(t, ctx, testPool, s3Client)

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
				defer tu.ClearAuthData(t, ctx, authCtx)
				// Setup admin user and create authenticated request
				accessToken := tu.SetupAdminUser(t, ctx, authCtx)

				request := domain.SetCompanyInstagramRequest{Instagram: url}
				req := th.NewSetCompanyInstagramRequest(t, ctx, testServerURL, request)

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
				assert.Contains(t, respBody.Error, "instagram_format")
			})
		}
	})

	t.Run("should return 400 for instagram link exceeding 255 characters", func(t *testing.T) {
		th.ClearSettingTestData(t, ctx, testPool, s3Client)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Create a long URL
		longPath := strings.Repeat("a", 240)
		longUrl := "https://instagram.com/" + longPath // Total > 255 chars

		// Setup admin user and create authenticated request
		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		request := domain.SetCompanyInstagramRequest{Instagram: longUrl}
		req := th.NewSetCompanyInstagramRequest(t, ctx, testServerURL, request)

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
		assert.Contains(t, respBody.Error, "instagram_length")
	})

	t.Run("should successfully accept various valid URL formats", func(t *testing.T) {
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
				th.ClearSettingTestData(t, ctx, testPool, s3Client)
				defer tu.ClearAuthData(t, ctx, authCtx)

				// Setup admin user and create authenticated request
				accessToken := tu.SetupAdminUser(t, ctx, authCtx)

				// Create a test channel for RabbitMQ verification
				// COMMENTED OUT: RabbitMQ verification disabled
				// testCh := th.GetRabbitMQChannel(t, testMQConn)
				// defer testCh.Close()

				// Purge queues to ensure clean state
				// th.PurgeSettingsQueues(t, testCh)

				request := domain.SetCompanyInstagramRequest{Instagram: url}
				req := th.NewSetCompanyInstagramRequest(t, ctx, testServerURL, request)

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

				var respBody domain.SetCompanyInstagramResponse
				err = json.NewDecoder(resp.Body).Decode(&respBody)
				require.NoError(t, err)
				assert.True(t, respBody.Success)

				// COMMENTED OUT: Verify RabbitMQ message was published (verification disabled)
				// COMMENTED OUT: RabbitMQ verification disabled
				// th.VerifySettingsUpdateMessage(t, testCh, settings.CompanyInstagram, url)
			})
		}
	})

	t.Run("should accept non-Instagram URLs (business requirement)", func(t *testing.T) {
		// Test that the field accepts any valid URL, not just Instagram
		nonInstagramUrls := []string{
			"https://twitter.com/company",
			"https://facebook.com/company",
			"https://linkedin.com/company/company-name",
			"https://company.com/social",
		}

		for _, url := range nonInstagramUrls {
			t.Run("non-Instagram URL: "+url, func(t *testing.T) {
				th.ClearSettingTestData(t, ctx, testPool, s3Client)
				defer tu.ClearAuthData(t, ctx, authCtx)

				// Create a test channel for RabbitMQ verification
				// COMMENTED OUT: RabbitMQ verification disabled
				// testCh := th.GetRabbitMQChannel(t, testMQConn)
				// defer testCh.Close()

				// Purge queues to ensure clean state
				// th.PurgeSettingsQueues(t, testCh)

				// Setup admin user and create authenticated request
				accessToken := tu.SetupAdminUser(t, ctx, authCtx)

				request := domain.SetCompanyInstagramRequest{Instagram: url}
				req := th.NewSetCompanyInstagramRequest(t, ctx, testServerURL, request)

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

				var respBody domain.SetCompanyInstagramResponse
				err = json.NewDecoder(resp.Body).Decode(&respBody)
				require.NoError(t, err)
				assert.True(t, respBody.Success)

				// COMMENTED OUT: Verify RabbitMQ message was published (verification disabled)
				// COMMENTED OUT: RabbitMQ verification disabled
				// th.VerifySettingsUpdateMessage(t, testCh, settings.CompanyInstagram, url)
			})
		}
	})

	t.Run("should return 415 for incorrect content type", func(t *testing.T) {
		th.ClearSettingTestData(t, ctx, testPool, s3Client)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup admin user and create authenticated request
		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		request := domain.SetCompanyInstagramRequest{Instagram: "https://instagram.com/test"}
		req := th.NewSetCompanyInstagramRequest(t, ctx, testServerURL, request)
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

	t.Run("should return 400 for unknown JSON fields", func(t *testing.T) {
		th.ClearSettingTestData(t, ctx, testPool, s3Client)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup admin user and create authenticated request
		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, testServerURL+httpEndpoints.SetCompanyInstagramEndpoint,
			strings.NewReader(`{"instagram": "https://instagram.com/test", "unknown_field": "value"}`))
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

	t.Run("should successfully update existing company instagram", func(t *testing.T) {
		th.ClearSettingTestData(t, ctx, testPool, s3Client)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Create a test channel for RabbitMQ verification
		// COMMENTED OUT: RabbitMQ verification disabled
		// testCh := th.GetRabbitMQChannel(t, testMQConn)
		// defer testCh.Close()

		// Purge queues to ensure clean state
		// th.PurgeSettingsQueues(t, testCh)

		// Set initial instagram

		oldInsta := "https://instagram.com/oldcompany"
		th.InsertCompanyInstagram(t, ctx, oldInsta, testPool)

		// Setup admin user and create authenticated request
		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		// Update to new instagram
		newInsta := "https://instagram.com/newcompany"

		request := domain.SetCompanyInstagramRequest{Instagram: newInsta}
		req := th.NewSetCompanyInstagramRequest(t, ctx, testServerURL, request)

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
		require.Equal(t, http.StatusOK, resp.StatusCode)

		// Verify updated instagram directly in database
		instagram, err := th.GetCompanyInstagram(t, ctx, testPool)
		require.NoError(t, err)
		assert.Equal(t, newInsta, instagram)

		// COMMENTED OUT: Verify RabbitMQ message was published (verification disabled)
		// COMMENTED OUT: RabbitMQ verification disabled
		// th.VerifySettingsUpdateMessage(t, testCh, settings.CompanyInstagram, newInsta)
	})
}
