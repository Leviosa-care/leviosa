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

// TEST=TestGetOTPDuration make test-integration-test

func TestGetOTPDuration(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("should return 404 when OTP duration not set", func(t *testing.T) {
		th.ClearSettingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup admin user and create authenticated request
		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		req := th.NewGetOTPDurationRequest(t, ctx, testServerURL)

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
		assert.Contains(t, respBody.Error, "OTP duration")
	})

	t.Run("should successfully retrieve OTP duration (admin endpoint)", func(t *testing.T) {
		th.ClearSettingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup: Insert OTP duration directly into database
		duration := 300
		th.InsertOTPDuration(t, ctx, duration, testPool) // 5 minutes

		// Setup admin user and create authenticated request
		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		req := th.NewGetOTPDurationRequest(t, ctx, testServerURL)

		// Add authentication to the request
		req.Header = tu.CreateAuthHeader(accessToken)
		req.AddCookie(tu.CreateAuthCookie(accessToken))

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var respBody domain.GetOTPDurationResponse
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		require.NoError(t, err)
		assert.Equal(t, duration, respBody.Duration)
	})
}

// TEST=TestSetOTPDuration make test-integration-test

func TestSetOTPDuration(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("should successfully set OTP duration", func(t *testing.T) {
		th.ClearSettingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Create a test channel for RabbitMQ verification
		testCh := th.GetRabbitMQChannel(t, testMQConn)
		defer testCh.Close()

		// Purge queues to ensure clean state
		th.PurgeSettingsQueues(t, testCh)

		// Setup admin user and create authenticated request
		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		duration := 600
		request := domain.SetOTPDurationRequest{Duration: duration} // 10 minutes
		req := th.NewSetOTPDurationRequest(t, ctx, testServerURL, request)

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

		var respBody domain.SetOTPDurationResponse
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		require.NoError(t, err)
		assert.True(t, respBody.Success)

		// Verify data was persisted directly in database
		retrievedDuration, err := th.GetOTPDurationFromDB(t, ctx, testPool)
		require.NoError(t, err)
		assert.Equal(t, duration, retrievedDuration)

		// Verify RabbitMQ message was published
		th.VerifySettingsUpdateMessage(t, testCh, settings.OTPDuration, duration)
	})

	t.Run("should return 400 for duration less than 60 seconds", func(t *testing.T) {

		invalidDurations := []int{0, -1, 30, 59}

		for _, duration := range invalidDurations {
			t.Run(fmt.Sprintf("invalid duration: %d", duration), func(t *testing.T) {
				th.ClearSettingsTable(t, ctx, testPool)
				defer tu.ClearAuthData(t, ctx, authCtx)

				// Setup admin user and create authenticated request
				accessToken := tu.SetupAdminUser(t, ctx, authCtx)

				request := domain.SetOTPDurationRequest{Duration: duration}
				req := th.NewSetOTPDurationRequest(t, ctx, testServerURL, request)

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
				assert.Contains(t, respBody.Error, "duration_min")
			})
		}
	})

	t.Run("should return 400 for duration greater than 3600 seconds", func(t *testing.T) {
		invalidDurations := []int{3601, 7200, 10000}

		for _, duration := range invalidDurations {
			t.Run(fmt.Sprintf("invalid duration: %d", duration), func(t *testing.T) {
				th.ClearSettingsTable(t, ctx, testPool)
				defer tu.ClearAuthData(t, ctx, authCtx)

				// Setup admin user and create authenticated request
				accessToken := tu.SetupAdminUser(t, ctx, authCtx)

				request := domain.SetOTPDurationRequest{Duration: duration}
				req := th.NewSetOTPDurationRequest(t, ctx, testServerURL, request)

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
				assert.Contains(t, respBody.Error, "duration_max")
			})
		}
	})

	t.Run("should successfully accept valid duration ranges", func(t *testing.T) {
		validDurations := []struct {
			duration int
			name     string
		}{
			{60, "1 minute (minimum)"},
			{120, "2 minutes"},
			{300, "5 minutes"},
			{600, "10 minutes"},
			{900, "15 minutes"},
			{1800, "30 minutes"},
			{3600, "1 hour (maximum)"},
		}

		for _, test := range validDurations {
			t.Run(test.name, func(t *testing.T) {
				th.ClearSettingsTable(t, ctx, testPool)
				defer tu.ClearAuthData(t, ctx, authCtx)

				// Create a test channel for RabbitMQ verification
				testCh := th.GetRabbitMQChannel(t, testMQConn)
				defer testCh.Close()

				// Purge queues to ensure clean state
				th.PurgeSettingsQueues(t, testCh)

				// Setup admin user and create authenticated request
				accessToken := tu.SetupAdminUser(t, ctx, authCtx)

				request := domain.SetOTPDurationRequest{Duration: test.duration}
				req := th.NewSetOTPDurationRequest(t, ctx, testServerURL, request)

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

				var respBody domain.SetOTPDurationResponse
				err = json.NewDecoder(resp.Body).Decode(&respBody)
				require.NoError(t, err)
				assert.True(t, respBody.Success)

				// Verify the duration was stored correctly directly in database
				retrievedDuration, err := th.GetOTPDurationFromDB(t, ctx, testPool)
				require.NoError(t, err)
				assert.Equal(t, test.duration, retrievedDuration)

				// Verify RabbitMQ message was published
				th.VerifySettingsUpdateMessage(t, testCh, settings.OTPDuration, test.duration)
			})
		}
	})

	t.Run("should return 415 for incorrect content type", func(t *testing.T) {
		th.ClearSettingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup admin user and create authenticated request
		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		request := domain.SetOTPDurationRequest{Duration: 300}
		req := th.NewSetOTPDurationRequest(t, ctx, testServerURL, request)
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
		th.ClearSettingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup admin user and create authenticated request
		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, testServerURL+httpEndpoints.AdminSetOTPDurationEndpoint,
			strings.NewReader(`{"duration": 300, "unknown_field": "value"}`))
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

	t.Run("should successfully update existing OTP duration", func(t *testing.T) {
		th.ClearSettingsTable(t, ctx, testPool)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup: Insert OTP duration directly into database
		duration := 300
		th.InsertOTPDuration(t, ctx, duration, testPool) // 5 minutes

		// Setup admin user and create authenticated request
		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		// Update to new duration
		newDuration := 900
		request := domain.SetOTPDurationRequest{Duration: newDuration}
		req := th.NewSetOTPDurationRequest(t, ctx, testServerURL, request)

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

		// Verify updated duration directly in database
		retrievedDuration, err := th.GetOTPDurationFromDB(t, ctx, testPool)
		require.NoError(t, err)
		assert.Equal(t, newDuration, retrievedDuration)
	})
}
