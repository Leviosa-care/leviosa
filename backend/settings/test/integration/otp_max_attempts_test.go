package helpers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/Leviosa-care/core/contracts/settings"
	tu "github.com/Leviosa-care/core/testutils"
	httpEndpoints "github.com/Leviosa-care/settings/internal/adapters/http"
	"github.com/Leviosa-care/settings/internal/domain"
	th "github.com/Leviosa-care/settings/test/helpers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TEST=TestGetOTPMaxAttempts make test-integration-test

func TestGetOTPMaxAttempts(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("should return 404 when OTP max attempts not set", func(t *testing.T) {
		th.ClearAllTestData(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup admin user and create authenticated request
		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		req := th.NewGetOTPMaxAttemptsRequest(t, ctx, testServerURL)

		// Add authentication to the request
		req.Header = tu.CreateAuthHeader(accessToken)
		req.AddCookie(tu.CreateAuthCookie(accessToken))

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		var respBody struct {
			Error string `json:"error"`
		}
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		require.NoError(t, err)
		assert.Contains(t, respBody.Error, "OTP max attempts")
	})

	t.Run("should successfully retrieve OTP max attempts (admin endpoint)", func(t *testing.T) {
		th.ClearAllTestData(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup: Insert OTP max attempts directly into database
		attempts := 5
		th.InsertOTPMaxAttempts(t, ctx, attempts, testPool)

		// Setup admin user and create authenticated request
		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		// Test: Get the OTP max attempts
		req := th.NewGetOTPMaxAttemptsRequest(t, ctx, testServerURL)

		// Add authentication to the request
		req.Header = tu.CreateAuthHeader(accessToken)
		req.AddCookie(tu.CreateAuthCookie(accessToken))

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var respBody domain.GetOTPMaxAttemptsResponse
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		require.NoError(t, err)
		assert.Equal(t, attempts, respBody.MaxAttempts)
	})
}

// TEST=TestSetOTPMaxAttempts make test-integration-test

func TestSetOTPMaxAttempts(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("should successfully set OTP max attempts", func(t *testing.T) {
		th.ClearAllTestData(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Create a test channel for RabbitMQ verification
		testCh := th.GetRabbitMQChannel(t, testMQConn)
		defer testCh.Close()

		// Purge queues to ensure clean state
		th.PurgeSettingsQueues(t, testCh)

		attempts := 3

		// Setup admin user and create authenticated request
		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		request := domain.SetOTPMaxAttemptsRequest{MaxAttempts: attempts}
		req := th.NewSetOTPMaxAttemptsRequest(t, ctx, testServerURL, request)

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

		var respBody domain.SetOTPMaxAttemptsResponse
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		require.NoError(t, err)
		assert.True(t, respBody.Success)

		// Verify data was persisted directly in database
		maxAttempts, err := th.GetOTPMaxAttemptsFromDB(t, ctx, testPool)
		require.NoError(t, err)
		assert.Equal(t, attempts, maxAttempts)

		// Verify RabbitMQ message was published (note: integer values are published as integers)
		th.VerifySettingsUpdateMessage(t, testCh, settings.OTPMaxAttempts, attempts)
	})

	t.Run("should return 400 for max attempts less than 1", func(t *testing.T) {
		invalidAttempts := []int{0, -1, -5}

		for _, attempts := range invalidAttempts {
			t.Run(fmt.Sprintf("invalid attempts: %d", attempts), func(t *testing.T) {
				th.ClearAllTestData(t, ctx, testPool)
				defer tu.ClearAuthData(t, ctx, authCtx)

				// Setup admin user and create authenticated request
				accessToken := tu.SetupAdminUser(t, ctx, authCtx)

				request := domain.SetOTPMaxAttemptsRequest{MaxAttempts: attempts}
				req := th.NewSetOTPMaxAttemptsRequest(t, ctx, testServerURL, request)

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
				assert.Contains(t, respBody.Error, "max_attempts_min")
			})
		}
	})

	t.Run("should return 400 for max attempts greater than 10", func(t *testing.T) {
		invalidAttempts := []int{11, 15, 20, 100}

		for _, attempts := range invalidAttempts {
			t.Run(fmt.Sprintf("invalid attempts: %d", attempts), func(t *testing.T) {
				th.ClearAllTestData(t, ctx, testPool)
				defer tu.ClearAuthData(t, ctx, authCtx)

				// Setup admin user and create authenticated request
				accessToken := tu.SetupAdminUser(t, ctx, authCtx)

				request := domain.SetOTPMaxAttemptsRequest{MaxAttempts: attempts}
				req := th.NewSetOTPMaxAttemptsRequest(t, ctx, testServerURL, request)

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
				assert.Contains(t, respBody.Error, "max_attempts_max")
			})
		}
	})

	t.Run("should successfully accept valid max attempts ranges", func(t *testing.T) {
		validAttempts := []struct {
			attempts int
			name     string
		}{
			{1, "1 attempt (minimum, very restrictive)"},
			{2, "2 attempts (restrictive)"},
			{3, "3 attempts (standard)"},
			{4, "4 attempts"},
			{5, "5 attempts (generous)"},
			{6, "6 attempts"},
			{7, "7 attempts"},
			{8, "8 attempts"},
			{9, "9 attempts"},
			{10, "10 attempts (maximum, very generous)"},
		}

		for _, test := range validAttempts {
			t.Run(test.name, func(t *testing.T) {
				th.ClearAllTestData(t, ctx, testPool)
				defer tu.ClearAuthData(t, ctx, authCtx)

				// Create a test channel for RabbitMQ verification
				testCh := th.GetRabbitMQChannel(t, testMQConn)
				defer testCh.Close()

				// Purge queues to ensure clean state
				th.PurgeSettingsQueues(t, testCh)

				// Setup admin user and create authenticated request
				accessToken := tu.SetupAdminUser(t, ctx, authCtx)

				request := domain.SetOTPMaxAttemptsRequest{MaxAttempts: test.attempts}
				req := th.NewSetOTPMaxAttemptsRequest(t, ctx, testServerURL, request)

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

				var respBody domain.SetOTPMaxAttemptsResponse
				err = json.NewDecoder(resp.Body).Decode(&respBody)
				require.NoError(t, err)
				assert.True(t, respBody.Success)

				// Verify the max attempts was stored correctly directly in database
				maxAttempts, err := th.GetOTPMaxAttemptsFromDB(t, ctx, testPool)
				require.NoError(t, err)
				assert.Equal(t, test.attempts, maxAttempts)

				// Verify RabbitMQ message was published (note: integer values are published as integers)
				th.VerifySettingsUpdateMessage(t, testCh, settings.OTPMaxAttempts, test.attempts)
			})
		}
	})

	t.Run("should return 415 for incorrect content type", func(t *testing.T) {
		th.ClearAllTestData(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup admin user and create authenticated request
		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		request := domain.SetOTPMaxAttemptsRequest{MaxAttempts: 3}
		req := th.NewSetOTPMaxAttemptsRequest(t, ctx, testServerURL, request)
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
		th.ClearAllTestData(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup admin user and create authenticated request
		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, testServerURL+httpEndpoints.AdminSetOTPMaxAttemptsEndpoint,
			strings.NewReader(`{"max_attempts": 3, "unknown_field": "value"}`))
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

	t.Run("should successfully update existing OTP max attempts", func(t *testing.T) {
		th.ClearAllTestData(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Create a test channel for RabbitMQ verification
		testCh := th.GetRabbitMQChannel(t, testMQConn)
		defer testCh.Close()

		// Purge queues to ensure clean state
		th.PurgeSettingsQueues(t, testCh)

		// Set initial max attempts
		th.InsertOTPMaxAttempts(t, ctx, 3, testPool)

		// Update to new max attempts
		newAttempts := 7

		// Setup admin user and create authenticated request
		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		request := domain.SetOTPMaxAttemptsRequest{MaxAttempts: newAttempts}
		req := th.NewSetOTPMaxAttemptsRequest(t, ctx, testServerURL, request)

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

		// Verify updated max attempts directly in database
		maxAttempts, err := th.GetOTPMaxAttemptsFromDB(t, ctx, testPool)
		require.NoError(t, err)
		assert.Equal(t, 7, maxAttempts)

		// Verify RabbitMQ message was published (note: integer values are published as integers)
		th.VerifySettingsUpdateMessage(t, testCh, settings.OTPMaxAttempts, newAttempts)
	})

	t.Run("should handle security considerations properly", func(t *testing.T) {
		// Test edge cases for security
		securityTests := []struct {
			attempts    int
			description string
		}{
			{1, "very restrictive (lockout after first failure)"},
			{3, "balanced security (standard recommendation)"},
			{10, "maximum allowed (still secure but user-friendly)"},
		}

		for _, test := range securityTests {
			t.Run(fmt.Sprintf("security test: %s", test.description), func(t *testing.T) {
				th.ClearAllTestData(t, ctx, testPool)
				defer tu.ClearAuthData(t, ctx, authCtx)

				// Setup admin user and create authenticated request
				accessToken := tu.SetupAdminUser(t, ctx, authCtx)

				request := domain.SetOTPMaxAttemptsRequest{MaxAttempts: test.attempts}
				req := th.NewSetOTPMaxAttemptsRequest(t, ctx, testServerURL, request)

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

				var respBody domain.SetOTPMaxAttemptsResponse
				err = json.NewDecoder(resp.Body).Decode(&respBody)
				require.NoError(t, err)
				assert.True(t, respBody.Success)
			})
		}
	})
}
