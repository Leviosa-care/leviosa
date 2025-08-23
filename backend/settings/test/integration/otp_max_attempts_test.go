package testdata

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
	td "github.com/Leviosa-care/settings/test/testdata"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetOTPMaxAttempts(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("should return 404 when OTP max attempts not set", func(t *testing.T) {
		td.ClearAllTestData(t, ctx, testPool)

		req := td.NewGetOTPMaxAttemptsRequest(t, ctx, testServerURL)
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
		td.ClearAllTestData(t, ctx, testPool)

		// Setup: Insert OTP max attempts directly into database
		td.InsertOTPMaxAttempts(t, ctx, 5, testPool)

		// Test: Get the OTP max attempts
		req := td.NewGetOTPMaxAttemptsRequest(t, ctx, testServerURL)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var respBody domain.GetOTPMaxAttemptsResponse
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		require.NoError(t, err)
		assert.Equal(t, 5, respBody.MaxAttempts)
	})
}

func TestSetOTPMaxAttempts(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("should successfully set OTP max attempts", func(t *testing.T) {
		td.ClearAllTestData(t, ctx, testPool)
		
		// Create a test channel for RabbitMQ verification
		testCh := td.GetRabbitMQChannel(t, testMQConn)
		defer testCh.Close()
		
		// Purge queues to ensure clean state
		td.PurgeSettingsQueues(t, testCh)

		request := domain.SetOTPMaxAttemptsRequest{MaxAttempts: 3}
		req := td.NewSetOTPMaxAttemptsRequest(t, ctx, testServerURL, request)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var respBody domain.SetOTPMaxAttemptsResponse
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		require.NoError(t, err)
		assert.True(t, respBody.Success)

		// Verify data was persisted directly in database
		maxAttempts, err := td.GetOTPMaxAttemptsFromDB(t, ctx, testPool)
		require.NoError(t, err)
		assert.Equal(t, 3, maxAttempts)

		// Verify RabbitMQ message was published (note: integer values are published as integers)
		td.VerifySettingsUpdateMessage(t, testCh, settings.OTPMaxAttempts, 3)
	})

	t.Run("should return 400 for max attempts less than 1", func(t *testing.T) {
		td.ClearAllTestData(t, ctx, testPool)

		invalidAttempts := []int{0, -1, -5}

		for _, attempts := range invalidAttempts {
			t.Run(fmt.Sprintf("invalid attempts: %d", attempts), func(t *testing.T) {
				request := domain.SetOTPMaxAttemptsRequest{MaxAttempts: attempts}
				req := td.NewSetOTPMaxAttemptsRequest(t, ctx, testServerURL, request)

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
		td.ClearAllTestData(t, ctx, testPool)

		invalidAttempts := []int{11, 15, 20, 100}

		for _, attempts := range invalidAttempts {
			t.Run(fmt.Sprintf("invalid attempts: %d", attempts), func(t *testing.T) {
				request := domain.SetOTPMaxAttemptsRequest{MaxAttempts: attempts}
				req := td.NewSetOTPMaxAttemptsRequest(t, ctx, testServerURL, request)

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
		td.ClearAllTestData(t, ctx, testPool)

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
				td.ClearAllTestData(t, ctx, testPool)

				request := domain.SetOTPMaxAttemptsRequest{MaxAttempts: test.attempts}
				req := td.NewSetOTPMaxAttemptsRequest(t, ctx, testServerURL, request)

				resp, err := client.Do(req)
				require.NoError(t, err)
				defer resp.Body.Close()

				assert.Equal(t, http.StatusOK, resp.StatusCode)

				var respBody domain.SetOTPMaxAttemptsResponse
				err = json.NewDecoder(resp.Body).Decode(&respBody)
				require.NoError(t, err)
				assert.True(t, respBody.Success)

				// Verify the max attempts was stored correctly directly in database
				maxAttempts, err := td.GetOTPMaxAttemptsFromDB(t, ctx, testPool)
				require.NoError(t, err)
				assert.Equal(t, test.attempts, maxAttempts)
			})
		}
	})

	t.Run("should return 415 for incorrect content type", func(t *testing.T) {
		td.ClearAllTestData(t, ctx, testPool)

		request := domain.SetOTPMaxAttemptsRequest{MaxAttempts: 3}
		req := td.NewSetOTPMaxAttemptsRequest(t, ctx, testServerURL, request)
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
		td.ClearAllTestData(t, ctx, testPool)

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, testServerURL+"/admin/settings/otp/max-attempts", 
			strings.NewReader(`{"max_attempts": 3, "unknown_field": "value"}`))
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

	t.Run("should successfully update existing OTP max attempts", func(t *testing.T) {
		td.ClearAllTestData(t, ctx, testPool)

		// Set initial max attempts
		request1 := domain.SetOTPMaxAttemptsRequest{MaxAttempts: 3}
		req1 := td.NewSetOTPMaxAttemptsRequest(t, ctx, testServerURL, request1)
		resp1, err := client.Do(req1)
		require.NoError(t, err)
		defer resp1.Body.Close()
		require.Equal(t, http.StatusOK, resp1.StatusCode)

		// Update to new max attempts
		request2 := domain.SetOTPMaxAttemptsRequest{MaxAttempts: 7}
		req2 := td.NewSetOTPMaxAttemptsRequest(t, ctx, testServerURL, request2)
		resp2, err := client.Do(req2)
		require.NoError(t, err)
		defer resp2.Body.Close()
		require.Equal(t, http.StatusOK, resp2.StatusCode)

		// Verify updated max attempts directly in database
		maxAttempts, err := td.GetOTPMaxAttemptsFromDB(t, ctx, testPool)
		require.NoError(t, err)
		assert.Equal(t, 7, maxAttempts)
	})

	t.Run("should handle security considerations properly", func(t *testing.T) {
		td.ClearAllTestData(t, ctx, testPool)

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
				td.ClearAllTestData(t, ctx, testPool)

				request := domain.SetOTPMaxAttemptsRequest{MaxAttempts: test.attempts}
				req := td.NewSetOTPMaxAttemptsRequest(t, ctx, testServerURL, request)

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