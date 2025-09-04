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

// TEST=TestGetOTPLength make test-integration-test

func TestGetOTPLength(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("should return 404 when OTP length not set", func(t *testing.T) {
		th.ClearAllTestData(t, ctx, testPool)

		req := th.NewGetOTPLengthRequest(t, ctx, testServerURL)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		var respBody struct {
			Error string `json:"error"`
		}
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		require.NoError(t, err)
		assert.Contains(t, respBody.Error, "OTP length")
	})

	t.Run("should successfully retrieve OTP length (admin endpoint)", func(t *testing.T) {
		th.ClearAllTestData(t, ctx, testPool)

		// Setup: Insert OTP length directly into database
		length := 6
		th.InsertOTPLength(t, ctx, length, testPool)

		// Test: Get the OTP length
		req := th.NewGetOTPLengthRequest(t, ctx, testServerURL)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var respBody domain.GetOTPLengthResponse
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		require.NoError(t, err)
		assert.Equal(t, length, respBody.Length)
	})
}

// TEST=TestSetOTPLength make test-integration-test

func TestSetOTPLength(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("should successfully set OTP length", func(t *testing.T) {
		th.ClearAllTestData(t, ctx, testPool)

		// Create a test channel for RabbitMQ verification
		testCh := th.GetRabbitMQChannel(t, testMQConn)
		defer testCh.Close()

		// Purge queues to ensure clean state
		th.PurgeSettingsQueues(t, testCh)

		length := 8
		request := domain.SetOTPLengthRequest{Length: length}
		req := th.NewSetOTPLengthRequest(t, ctx, testServerURL, request)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var respBody domain.SetOTPLengthResponse
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		require.NoError(t, err)
		assert.True(t, respBody.Success)

		// Verify data was persisted directly in database
		retrievedLength, err := th.GetOTPLengthFromDB(t, ctx, testPool)
		require.NoError(t, err)
		assert.Equal(t, length, retrievedLength)

		// Verify RabbitMQ message was published
		th.VerifySettingsUpdateMessage(t, testCh, settings.OTPLength, 8)
	})

	t.Run("should return 400 for length less than 4 digits", func(t *testing.T) {
		th.ClearAllTestData(t, ctx, testPool)

		invalidLengths := []int{0, -1, 1, 2, 3}

		for _, length := range invalidLengths {
			t.Run(fmt.Sprintf("invalid length: %d", length), func(t *testing.T) {
				request := domain.SetOTPLengthRequest{Length: length}
				req := th.NewSetOTPLengthRequest(t, ctx, testServerURL, request)

				resp, err := client.Do(req)
				require.NoError(t, err)
				defer resp.Body.Close()

				assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

				var respBody struct {
					Error string `json:"error"`
				}
				err = json.NewDecoder(resp.Body).Decode(&respBody)
				require.NoError(t, err)
				assert.Contains(t, respBody.Error, "length_min")
			})
		}
	})

	t.Run("should return 400 for length greater than 10 digits", func(t *testing.T) {
		th.ClearAllTestData(t, ctx, testPool)

		invalidLengths := []int{11, 12, 15, 20}

		for _, length := range invalidLengths {
			t.Run(fmt.Sprintf("invalid length: %d", length), func(t *testing.T) {
				request := domain.SetOTPLengthRequest{Length: length}
				req := th.NewSetOTPLengthRequest(t, ctx, testServerURL, request)

				resp, err := client.Do(req)
				require.NoError(t, err)
				defer resp.Body.Close()

				assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

				var respBody struct {
					Error string `json:"error"`
				}
				err = json.NewDecoder(resp.Body).Decode(&respBody)
				require.NoError(t, err)
				assert.Contains(t, respBody.Error, "length_max")
			})
		}
	})

	t.Run("should successfully accept valid length ranges", func(t *testing.T) {
		th.ClearAllTestData(t, ctx, testPool)

		validLengths := []struct {
			length int
			name   string
		}{
			{4, "4 digits (minimum)"},
			{5, "5 digits"},
			{6, "6 digits (common)"},
			{7, "7 digits"},
			{8, "8 digits (secure)"},
			{9, "9 digits"},
			{10, "10 digits (maximum)"},
		}

		for _, test := range validLengths {
			t.Run(test.name, func(t *testing.T) {
				th.ClearAllTestData(t, ctx, testPool)

				// Create a test channel for RabbitMQ verification
				testCh := th.GetRabbitMQChannel(t, testMQConn)
				defer testCh.Close()

				// Purge queues to ensure clean state
				th.PurgeSettingsQueues(t, testCh)

				request := domain.SetOTPLengthRequest{Length: test.length}
				req := th.NewSetOTPLengthRequest(t, ctx, testServerURL, request)

				resp, err := client.Do(req)
				require.NoError(t, err)
				defer resp.Body.Close()

				assert.Equal(t, http.StatusOK, resp.StatusCode)

				var respBody domain.SetOTPLengthResponse
				err = json.NewDecoder(resp.Body).Decode(&respBody)
				require.NoError(t, err)
				assert.True(t, respBody.Success)

				// Verify the length was stored correctly directly in database
				length, err := th.GetOTPLengthFromDB(t, ctx, testPool)
				require.NoError(t, err)
				assert.Equal(t, test.length, length)

				// Verify RabbitMQ message was published
				th.VerifySettingsUpdateMessage(t, testCh, settings.OTPLength, length)
			})
		}
	})

	t.Run("should return 415 for incorrect content type", func(t *testing.T) {
		th.ClearAllTestData(t, ctx, testPool)

		request := domain.SetOTPLengthRequest{Length: 6}
		req := th.NewSetOTPLengthRequest(t, ctx, testServerURL, request)
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

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, testServerURL+"/admin/settings/otp/length",
			strings.NewReader(`{"length": 6, "unknown_field": "value"}`))
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

	t.Run("should successfully update existing OTP length", func(t *testing.T) {
		th.ClearAllTestData(t, ctx, testPool)

		// Create a test channel for RabbitMQ verification
		testCh := th.GetRabbitMQChannel(t, testMQConn)
		defer testCh.Close()

		// Purge queues to ensure clean state
		th.PurgeSettingsQueues(t, testCh)

		// Set initial length
		oldLength := 4
		th.InsertOTPLength(t, ctx, oldLength, testPool)

		// Update to new length
		newLength := 8
		request2 := domain.SetOTPLengthRequest{Length: newLength}
		req2 := th.NewSetOTPLengthRequest(t, ctx, testServerURL, request2)
		resp2, err := client.Do(req2)
		require.NoError(t, err)
		defer resp2.Body.Close()
		require.Equal(t, http.StatusOK, resp2.StatusCode)

		// Verify updated length directly in database
		retrievedLength, err := th.GetOTPLengthFromDB(t, ctx, testPool)
		require.NoError(t, err)
		assert.Equal(t, newLength, retrievedLength)

		// Verify RabbitMQ message was published
		th.VerifySettingsUpdateMessage(t, testCh, settings.OTPLength, newLength)
	})
}
