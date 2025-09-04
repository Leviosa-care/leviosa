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

		req := th.NewGetOTPDurationRequest(t, ctx, testServerURL)
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

		// Setup: Insert OTP duration directly into database
		duration := 300
		th.InsertOTPDuration(t, ctx, duration, testPool) // 5 minutes

		// Test: Get the OTP duration via HTTP
		req := th.NewGetOTPDurationRequest(t, ctx, testServerURL)
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

		// Create a test channel for RabbitMQ verification
		testCh := th.GetRabbitMQChannel(t, testMQConn)
		defer testCh.Close()

		// Purge queues to ensure clean state
		th.PurgeSettingsQueues(t, testCh)

		duration := 600
		request := domain.SetOTPDurationRequest{Duration: duration} // 10 minutes
		req := th.NewSetOTPDurationRequest(t, ctx, testServerURL, request)

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
		th.ClearSettingsTable(t, ctx, testPool)

		invalidDurations := []int{0, -1, 30, 59}

		for _, duration := range invalidDurations {
			t.Run(fmt.Sprintf("invalid duration: %d", duration), func(t *testing.T) {
				request := domain.SetOTPDurationRequest{Duration: duration}
				req := th.NewSetOTPDurationRequest(t, ctx, testServerURL, request)

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
		th.ClearSettingsTable(t, ctx, testPool)

		invalidDurations := []int{3601, 7200, 10000}

		for _, duration := range invalidDurations {
			t.Run(fmt.Sprintf("invalid duration: %d", duration), func(t *testing.T) {
				request := domain.SetOTPDurationRequest{Duration: duration}
				req := th.NewSetOTPDurationRequest(t, ctx, testServerURL, request)

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
		th.ClearSettingsTable(t, ctx, testPool)

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

				// Create a test channel for RabbitMQ verification
				testCh := th.GetRabbitMQChannel(t, testMQConn)
				defer testCh.Close()

				// Purge queues to ensure clean state
				th.PurgeSettingsQueues(t, testCh)

				request := domain.SetOTPDurationRequest{Duration: test.duration}
				req := th.NewSetOTPDurationRequest(t, ctx, testServerURL, request)

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

		request := domain.SetOTPDurationRequest{Duration: 300}
		req := th.NewSetOTPDurationRequest(t, ctx, testServerURL, request)
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
		th.ClearSettingsTable(t, ctx, testPool)

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, testServerURL+"/admin/settings/otp/duration",
			strings.NewReader(`{"duration": 300, "unknown_field": "value"}`))
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

	t.Run("should successfully update existing OTP duration", func(t *testing.T) {
		th.ClearSettingsTable(t, ctx, testPool)

		// Set initial duration
		request1 := domain.SetOTPDurationRequest{Duration: 300}
		req1 := th.NewSetOTPDurationRequest(t, ctx, testServerURL, request1)
		resp1, err := client.Do(req1)
		require.NoError(t, err)
		defer resp1.Body.Close()
		require.Equal(t, http.StatusOK, resp1.StatusCode)

		// Update to new duration
		newDuration := 900
		request2 := domain.SetOTPDurationRequest{Duration: newDuration}
		req2 := th.NewSetOTPDurationRequest(t, ctx, testServerURL, request2)
		resp2, err := client.Do(req2)
		require.NoError(t, err)
		defer resp2.Body.Close()
		require.Equal(t, http.StatusOK, resp2.StatusCode)

		// Verify updated duration directly in database
		retrievedDuration, err := th.GetOTPDurationFromDB(t, ctx, testPool)
		require.NoError(t, err)
		assert.Equal(t, newDuration, retrievedDuration)
	})
}
