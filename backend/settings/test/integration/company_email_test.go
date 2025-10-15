package helpers

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/Leviosa-care/core/contracts/settings"
	"github.com/Leviosa-care/core/errs"
	tu "github.com/Leviosa-care/core/testutils"
	"github.com/Leviosa-care/core/validation"
	httpEndpoints "github.com/Leviosa-care/settings/internal/adapters/http"
	"github.com/Leviosa-care/settings/internal/domain"
	th "github.com/Leviosa-care/settings/test/helpers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TEST=TestGetCompanyEmail make test-integration-test

func TestGetCompanyEmail(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("should return 404 when company email not set", func(t *testing.T) {
		th.ClearSettingsTable(t, ctx, testPool)

		req := th.NewGetCompanyEmailRequest(t, ctx, testServerURL)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		var respBody struct {
			Error string `json:"error"`
		}
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		require.NoError(t, err)
		assert.Contains(t, respBody.Error, "company email")
	})

	t.Run("should successfully retrieve company email", func(t *testing.T) {
		th.ClearSettingsTable(t, ctx, testPool)

		// Setup: Insert company email directly into database
		email := "support@testcompany.com"
		th.InsertCompanyEmail(t, ctx, email, testPool)

		// Test: Get the company email via HTTP
		req := th.NewGetCompanyEmailRequest(t, ctx, testServerURL)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var respBody domain.GetCompanyEmailResponse
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		require.NoError(t, err)
		assert.Equal(t, email, respBody.Email)
	})
}

// TEST=TestSetCompanyEmail make test-integration-test

func TestSetCompanyEmail(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("should successfully set company email", func(t *testing.T) {
		th.ClearSettingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Create a test channel for RabbitMQ verification
		testCh := th.GetRabbitMQChannel(t, testMQConn)
		defer testCh.Close()

		// Purge queues to ensure clean state
		th.PurgeSettingsQueues(t, testCh)

		// Setup admin user and create authenticated request
		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		request := domain.SetCompanyEmailRequest{Email: "contact@newcompany.com"}
		req := th.NewSetCompanyEmailRequest(t, ctx, testServerURL, request)

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

		var respBody domain.SetCompanyEmailResponse
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		require.NoError(t, err)
		assert.True(t, respBody.Success)

		// Verify data was persisted directly in database
		email, err := th.GetCompanyEmailFromDB(t, ctx, testPool)
		require.NoError(t, err)
		assert.Equal(t, "contact@newcompany.com", email)

		// Verify RabbitMQ message was published
		th.VerifySettingsUpdateMessage(t, testCh, settings.CompanyEmail, "contact@newcompany.com")
	})

	t.Run("should return 400 for empty email", func(t *testing.T) {
		th.ClearSettingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)
		request := domain.SetCompanyEmailRequest{Email: ""}
		req := th.NewSetCompanyEmailRequest(t, ctx, testServerURL, request)

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
		assert.Contains(t, respBody.Error, "email_required")
	})

	t.Run("should return 400 for invalid email format", func(t *testing.T) {
		invalidEmails := []string{
			"notanemail",
			"@domain.com",
			"user@",
			"user@domain",
			"user space@domain.com",
			"user..double.dot@domain.com",
		}

		for _, email := range invalidEmails {
			t.Run("invalid email: "+email, func(t *testing.T) {
				th.ClearSettingsTable(t, ctx, testPool)
				defer tu.ClearAuthData(t, ctx, authCtx)

				accessToken := tu.SetupAdminUser(t, ctx, authCtx)
				request := domain.SetCompanyEmailRequest{Email: email}
				req := th.NewSetCompanyEmailRequest(t, ctx, testServerURL, request)

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
				assert.Contains(t, respBody.Error, "email_format")
			})
		}
	})

	t.Run("should return 400 for email exceeding 255 characters", func(t *testing.T) {
		th.ClearSettingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		// Create a long email address
		longLocalPart := strings.Repeat("a", validation.EmailMaxLength)
		longEmail := longLocalPart + "@domain.com" // Total > 255 chars

		request := domain.SetCompanyEmailRequest{Email: longEmail}
		req := th.NewSetCompanyEmailRequest(t, ctx, testServerURL, request)

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
		assert.Contains(t, respBody.Error, errs.ErrInvalidValue.Error())
	})

	t.Run("should successfully accept valid email formats", func(t *testing.T) {
		validEmails := []string{
			"simple@domain.com",
			"user.name@domain.com",
			"user+tag@domain.com",
			"123@domain.com",
			"test@sub.domain.com",
			"user@domain-name.com",
		}

		for _, email := range validEmails {
			t.Run("valid email: "+email, func(t *testing.T) {
				th.ClearSettingsTable(t, ctx, testPool)
				defer tu.ClearAuthData(t, ctx, authCtx)

				// Create a test channel for RabbitMQ verification
				testCh := th.GetRabbitMQChannel(t, testMQConn)
				defer testCh.Close()

				// Purge queues to ensure clean state
				th.PurgeSettingsQueues(t, testCh)

				// Setup admin user and create authenticated request
				accessToken := tu.SetupAdminUser(t, ctx, authCtx)

				request := domain.SetCompanyEmailRequest{Email: email}
				req := th.NewSetCompanyEmailRequest(t, ctx, testServerURL, request)

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

				var respBody domain.SetCompanyEmailResponse
				err = json.NewDecoder(resp.Body).Decode(&respBody)
				require.NoError(t, err)
				assert.True(t, respBody.Success)

				// Verify RabbitMQ message was published
				th.VerifySettingsUpdateMessage(t, testCh, settings.CompanyEmail, email)
			})
		}
	})

	t.Run("should return 415 for incorrect content type", func(t *testing.T) {
		th.ClearSettingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)
		request := domain.SetCompanyEmailRequest{Email: "test@company.com"}
		req := th.NewSetCompanyEmailRequest(t, ctx, testServerURL, request)

		req.Header.Set("Content-Type", "text/plain")
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
		th.ClearSettingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, testServerURL+httpEndpoints.SetCompanyEmailEndpoint,
			strings.NewReader(`{"email": "test@company.com", "unknown_field": "value"}`))
		require.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")
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

	t.Run("should successfully update existing company email", func(t *testing.T) {
		th.ClearSettingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Create a test channel for RabbitMQ verification
		testCh := th.GetRabbitMQChannel(t, testMQConn)
		defer testCh.Close()

		// Purge queues to ensure clean state
		th.PurgeSettingsQueues(t, testCh)

		// Set initial email
		oldEmail := "old@company.com"
		th.InsertCompanyEmail(t, ctx, oldEmail, testPool)

		// Setup admin user and update to new email
		accessToken := tu.SetupAdminUser(t, ctx, authCtx)
		newEmail := "new@company.com"
		request := domain.SetCompanyEmailRequest{Email: newEmail}
		req := th.NewSetCompanyEmailRequest(t, ctx, testServerURL, request)
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

		// Verify updated email directly in database
		email, err := th.GetCompanyEmailFromDB(t, ctx, testPool)
		require.NoError(t, err)
		assert.Equal(t, "new@company.com", email)

		// Verify RabbitMQ message was published
		th.VerifySettingsUpdateMessage(t, testCh, settings.CompanyEmail, newEmail)
	})
}
