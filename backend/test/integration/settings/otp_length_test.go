package helpers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	tu "github.com/Leviosa-care/leviosa/backend/internal/common/testutils"
	"github.com/Leviosa-care/leviosa/backend/internal/settings/domain"
	httpEndpoints "github.com/Leviosa-care/leviosa/backend/internal/settings/interface/http"
	th "github.com/Leviosa-care/leviosa/backend/test/helpers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestGetOTPLength TEST_PATH=test/integration/settings/otp_length_test.go

func TestGetOTPLength(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("should return 404 when OTP length not set", func(t *testing.T) {
		th.ClearSettingTestData(t, ctx, testPool, s3Client)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup admin user and create authenticated request
		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		req := th.NewGetOTPLengthRequest(t, ctx, testServerURL)

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
		assert.Contains(t, respBody.Error, "OTP length")
	})

	t.Run("should successfully retrieve OTP length (admin endpoint)", func(t *testing.T) {
		th.ClearSettingTestData(t, ctx, testPool, s3Client)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup: Insert OTP length directly into database
		length := 6
		th.InsertOTPLength(t, ctx, length, testPool)

		// Setup admin user and create authenticated request
		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		req := th.NewGetOTPLengthRequest(t, ctx, testServerURL)

		// Add authentication to the request
		req.Header = tu.CreateAuthHeader(accessToken)
		req.AddCookie(tu.CreateAuthCookie(accessToken))

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

// make test-func TEST_NAME=TestSetOTPLength TEST_PATH=test/integration/settings/otp_length_test.go

func TestSetOTPLength(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("should successfully set OTP length", func(t *testing.T) {
		th.ClearSettingTestData(t, ctx, testPool, s3Client)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Create a test channel for RabbitMQ verification
		// COMMENTED OUT: RabbitMQ verification disabled
		// testCh := th.GetRabbitMQChannel(t, testMQConn)
		// defer testCh.Close()

		// Purge queues to ensure clean state
		// th.PurgeSettingsQueues(t, testCh)

		length := 8

		// Setup admin user and create authenticated request
		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		request := domain.SetOTPLengthRequest{Length: length}
		req := th.NewSetOTPLengthRequest(t, ctx, testServerURL, request)

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

		var respBody domain.SetOTPLengthResponse
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		require.NoError(t, err)
		assert.True(t, respBody.Success)

		// Verify data was persisted directly in database
		retrievedLength, err := th.GetOTPLength(t, ctx, testPool)
		require.NoError(t, err)
		assert.Equal(t, length, retrievedLength)

		// COMMENTED OUT: Verify RabbitMQ message was published (verification disabled)
		// COMMENTED OUT: RabbitMQ verification disabled
		// th.VerifySettingsUpdateMessage(t, testCh, settings.OTPLength, 8)
	})

	t.Run("should return 400 for length less than 4 digits", func(t *testing.T) {
		invalidLengths := []int{0, -1, 1, 2, 3}

		for _, length := range invalidLengths {
			t.Run(fmt.Sprintf("invalid length: %d", length), func(t *testing.T) {
				th.ClearSettingTestData(t, ctx, testPool, s3Client)
				defer tu.ClearAuthData(t, ctx, authCtx)

				// Setup admin user and create authenticated request
				accessToken := tu.SetupAdminUser(t, ctx, authCtx)

				request := domain.SetOTPLengthRequest{Length: length}
				req := th.NewSetOTPLengthRequest(t, ctx, testServerURL, request)

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
				assert.Contains(t, respBody.Error, "length_min")
			})
		}
	})

	t.Run("should return 400 for length greater than 10 digits", func(t *testing.T) {
		invalidLengths := []int{11, 12, 15, 20}

		for _, length := range invalidLengths {
			t.Run(fmt.Sprintf("invalid length: %d", length), func(t *testing.T) {
				th.ClearSettingTestData(t, ctx, testPool, s3Client)
				defer tu.ClearAuthData(t, ctx, authCtx)

				// Setup admin user and create authenticated request
				accessToken := tu.SetupAdminUser(t, ctx, authCtx)

				request := domain.SetOTPLengthRequest{Length: length}
				req := th.NewSetOTPLengthRequest(t, ctx, testServerURL, request)

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
				assert.Contains(t, respBody.Error, "length_max")
			})
		}
	})

	t.Run("should successfully accept valid length ranges", func(t *testing.T) {
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

				request := domain.SetOTPLengthRequest{Length: test.length}
				req := th.NewSetOTPLengthRequest(t, ctx, testServerURL, request)

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

				var respBody domain.SetOTPLengthResponse
				err = json.NewDecoder(resp.Body).Decode(&respBody)
				require.NoError(t, err)
				assert.True(t, respBody.Success)

				// Verify the length was stored correctly directly in database
				length, err := th.GetOTPLength(t, ctx, testPool)
				require.NoError(t, err)
				assert.Equal(t, test.length, length)

				// COMMENTED OUT: Verify RabbitMQ message was published (verification disabled)
				// COMMENTED OUT: RabbitMQ verification disabled
				// th.VerifySettingsUpdateMessage(t, testCh, settings.OTPLength, length)
			})
		}
	})

	t.Run("should return 415 for incorrect content type", func(t *testing.T) {
		th.ClearSettingTestData(t, ctx, testPool, s3Client)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Setup admin user and create authenticated request
		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		request := domain.SetOTPLengthRequest{Length: 6}
		req := th.NewSetOTPLengthRequest(t, ctx, testServerURL, request)
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

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, testServerURL+httpEndpoints.AdminSetOTPLengthEndpoint,
			strings.NewReader(`{"length": 6, "unknown_field": "value"}`))
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

	t.Run("should successfully update existing OTP length", func(t *testing.T) {
		th.ClearSettingTestData(t, ctx, testPool, s3Client)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Create a test channel for RabbitMQ verification
		// COMMENTED OUT: RabbitMQ verification disabled
		// testCh := th.GetRabbitMQChannel(t, testMQConn)
		// defer testCh.Close()

		// Purge queues to ensure clean state
		// th.PurgeSettingsQueues(t, testCh)

		// Set initial length
		oldLength := 4
		th.InsertOTPLength(t, ctx, oldLength, testPool)

		// Update to new length
		newLength := 8

		// Setup admin user and create authenticated request
		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		request := domain.SetOTPLengthRequest{Length: newLength}
		req := th.NewSetOTPLengthRequest(t, ctx, testServerURL, request)

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

		// Verify updated length directly in database
		retrievedLength, err := th.GetOTPLength(t, ctx, testPool)
		require.NoError(t, err)
		assert.Equal(t, newLength, retrievedLength)

		// COMMENTED OUT: Verify RabbitMQ message was published (verification disabled)
		// COMMENTED OUT: RabbitMQ verification disabled
		// th.VerifySettingsUpdateMessage(t, testCh, settings.OTPLength, newLength)
	})
}
