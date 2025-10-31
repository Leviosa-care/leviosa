package auth_test

import (
	"context"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	authEndpoints "github.com/Leviosa-care/leviosa/backend/internal/authuser/interface/auth"
	td "github.com/Leviosa-care/leviosa/backend/test/helpers"
	"github.com/hengadev/encx"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestRequestPasswordReset TEST_PATH=test/integration/authuser/auth/request_password_reset_test.go:w

func TestRequestPasswordReset(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	// Setup RabbitMQ for OTP message verification
	ch, err := testMQConn.Channel()
	require.NoError(t, err)
	defer ch.Close()

	// Setup OTP queue and start consuming
	td.SetupOTPQueue(t, ch)
	msgs := td.ConsumeOTPMessages(t, ch)

	t.Run("should successfully request password reset for registered email", func(t *testing.T) {
		// Clean state and insert existing user
		td.ClearAllTestData(t, ctx, testPool, redisClient)
		td.PurgeOTPQueue(t, ch)

		existingEmail := "user@example.com"
		user := td.NewTestUser(t, existingEmail, "John", "DOE")
		userEncx, err := domain.ProcessUserEncx(ctx, crypto, user)
		require.NoError(t, err)
		err = td.InsertUserEncx(t, ctx, userEncx, testPool)
		require.NoError(t, err)

		// Test password reset request
		request := domain.RequestPasswordResetRequest{
			Email: existingEmail,
		}
		req := td.NewRequestPasswordResetRequest(t, ctx, testServerURL, request)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert HTTP response
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		message, status := td.ParseRequestPasswordResetResponse(t, resp)
		assert.Equal(t, "Password reset email sent successfully", message)
		assert.Equal(t, "sent", status)

		// Verify OTP was created in Redis
		emailBytes, err := encx.SerializeValue(request.Email)
		require.NoError(t, err)
		emailHash := crypto.HashBasic(ctx, emailBytes)
		exists := td.CheckOTPExists(t, ctx, emailHash, redisClient)
		assert.True(t, exists, "OTP should exist in Redis")

		// Verify OTP has proper TTL
		ttl := td.GetOTPTTL(t, ctx, emailHash, redisClient)
		assert.True(t, ttl > 8*time.Minute, "OTP TTL should be around 10 minutes")
		assert.True(t, ttl <= 10*time.Minute, "OTP TTL should not exceed 10 minutes")

		// Verify RabbitMQ message was sent
		delivery := td.WaitForOTPMessage(t, msgs, 2*time.Second)

		// Get OTP from Redis to verify the code
		otpEncx, err := td.GetOTPEncxByEmailHash(t, ctx, emailHash, redisClient)
		assert.NoError(t, err)

		otp, err := domain.DecryptOTPEncx(ctx, crypto, otpEncx)
		require.NoError(t, err)

		// Verify message content
		td.VerifyOTPMessage(t, delivery, request.Email, otp.Code)
	})

	t.Run("should return not found when email is not registered", func(t *testing.T) {
		// Clean state
		td.ClearAllTestData(t, ctx, testPool, redisClient)
		td.PurgeOTPQueue(t, ch)

		// Test with non-existent email
		request := domain.RequestPasswordResetRequest{
			Email: "nonexistent@example.com",
		}
		req := td.NewRequestPasswordResetRequest(t, ctx, testServerURL, request)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert not found response
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		errorMsg, statusCode := td.ParseErrorResponse(t, resp)
		assert.Contains(t, errorMsg, "email is not registered")
		assert.Equal(t, http.StatusNotFound, statusCode)

		// Verify no OTP was created
		emailBytes, err := encx.SerializeValue(request.Email)
		require.NoError(t, err)
		emailHash := crypto.HashBasic(ctx, emailBytes)
		exists := td.CheckOTPExists(t, ctx, emailHash, redisClient)
		assert.False(t, exists, "OTP should not exist in Redis")

		// Verify no RabbitMQ message was sent
		select {
		case <-msgs:
			t.Fatal("No message should be sent for non-existent email")
		case <-time.After(1 * time.Second):
			// Expected: no message should be received
		}
	})

	t.Run("should return bad request for invalid email format", func(t *testing.T) {
		// Clean state
		td.ClearAllTestData(t, ctx, testPool, redisClient)
		td.PurgeOTPQueue(t, ch)

		// Test with invalid email
		request := domain.RequestPasswordResetRequest{
			Email: "invalid-email",
		}
		req := td.NewRequestPasswordResetRequest(t, ctx, testServerURL, request)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert bad request response
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		errorMsg, statusCode := td.ParseErrorResponse(t, resp)
		assert.Contains(t, strings.ToLower(errorMsg), "email")
		assert.Equal(t, http.StatusBadRequest, statusCode)

		// Verify no OTP was created
		emailBytes, err := encx.SerializeValue(request.Email)
		require.NoError(t, err)
		emailHash := crypto.HashBasic(ctx, emailBytes)
		exists := td.CheckOTPExists(t, ctx, emailHash, redisClient)
		assert.False(t, exists, "OTP should not exist in Redis")

		// Verify no RabbitMQ message was sent
		select {
		case <-msgs:
			t.Fatal("No message should be sent for invalid email")
		case <-time.After(1 * time.Second):
			// Expected: no message should be received
		}
	})

	t.Run("should return bad request for empty email", func(t *testing.T) {
		// Clean state
		td.ClearAllTestData(t, ctx, testPool, redisClient)
		td.PurgeOTPQueue(t, ch)

		// Test with empty email
		request := domain.RequestPasswordResetRequest{
			Email: "",
		}
		req := td.NewRequestPasswordResetRequest(t, ctx, testServerURL, request)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert bad request response
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		errorMsg, statusCode := td.ParseErrorResponse(t, resp)
		assert.Contains(t, strings.ToLower(errorMsg), "email")
		assert.Equal(t, http.StatusBadRequest, statusCode)

		// Verify no OTP was create
		emailBytes, err := encx.SerializeValue(request.Email)
		require.NoError(t, err)
		emailHash := crypto.HashBasic(ctx, emailBytes)
		exists := td.CheckOTPExists(t, ctx, emailHash, redisClient)
		assert.False(t, exists, "OTP should not exist in Redis")

		// Verify no RabbitMQ message was sent
		select {
		case <-msgs:
			t.Fatal("No message should be sent for empty email")
		case <-time.After(1 * time.Second):
			// Expected: no message should be received
		}
	})

	t.Run("should handle OTP rate limiting properly", func(t *testing.T) {
		// Clean state and insert existing user
		td.ClearAllTestData(t, ctx, testPool, redisClient)
		td.PurgeOTPQueue(t, ch)

		existingEmail := "ratelimit@example.com"
		// td.InsertTestUser(t, ctx, existingEmail, "Rate", "Limited", testPool, crypto)
		user := td.NewTestUser(t, existingEmail, "Rate", "Limited")
		userEncx, err := domain.ProcessUserEncx(ctx, crypto, user)
		require.NoError(t, err)
		err = td.InsertUserEncx(t, ctx, userEncx, testPool)
		require.NoError(t, err)

		// Send multiple password reset requests rapidly
		request := domain.RequestPasswordResetRequest{
			Email: existingEmail,
		}

		// First request should succeed
		req1 := td.NewRequestPasswordResetRequest(t, ctx, testServerURL, request)
		resp1, err := client.Do(req1)
		require.NoError(t, err)
		defer resp1.Body.Close()

		assert.Equal(t, http.StatusOK, resp1.StatusCode)

		// Immediate second request might hit rate limit (depending on OTP service configuration)
		req2 := td.NewRequestPasswordResetRequest(t, ctx, testServerURL, request)
		resp2, err := client.Do(req2)
		require.NoError(t, err)
		defer resp2.Body.Close()

		// Either succeeds (if no rate limiting) or returns 429
		if resp2.StatusCode == http.StatusTooManyRequests {
			errorMsg, statusCode := td.ParseErrorResponse(t, resp2)
			assert.Contains(t, strings.ToLower(errorMsg), "rate")
			assert.Equal(t, http.StatusTooManyRequests, statusCode)
		} else {
			// If rate limiting is not strict, it might still succeed
			assert.Equal(t, http.StatusOK, resp2.StatusCode)
		}
	})

	t.Run("should handle malformed JSON request", func(t *testing.T) {
		// Clean state
		td.ClearAllTestData(t, ctx, testPool, redisClient)

		// Create request with malformed JSON
		req, err := http.NewRequestWithContext(
			ctx,
			http.MethodPost,
			testServerURL+authEndpoints.RequestPasswordResetEndpoint,
			strings.NewReader(`{"email": "test@example.com", "invalid": }`),
		)
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert bad request response
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		errorMsg, statusCode := td.ParseErrorResponse(t, resp)
		assert.Contains(t, strings.ToLower(errorMsg), "request body")
		assert.Equal(t, http.StatusBadRequest, statusCode)
	})

	t.Run("should handle missing content-type header", func(t *testing.T) {
		// Clean state and insert existing user
		td.ClearAllTestData(t, ctx, testPool, redisClient)

		existingEmail := "contenttype@example.com"
		// td.InsertTestUser(t, ctx, existingEmail, "Content", "Type", testPool, crypto)

		user := td.NewTestUser(t, existingEmail, "Content", "Type")
		userEncx, err := domain.ProcessUserEncx(ctx, crypto, user)
		require.NoError(t, err)
		err = td.InsertUserEncx(t, ctx, userEncx, testPool)
		require.NoError(t, err)

		// Create request without Content-Type header
		req, err := http.NewRequestWithContext(
			ctx,
			http.MethodPost,
			testServerURL+authEndpoints.RequestPasswordResetEndpoint,
			strings.NewReader(`{"email": "`+existingEmail+`"}`),
		)
		require.NoError(t, err)
		// Intentionally not setting Content-Type header

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Should still work as Go's decoder can handle JSON without explicit Content-Type
		// but response might vary based on implementation
		if resp.StatusCode == http.StatusOK {
			message, status := td.ParseRequestPasswordResetResponse(t, resp)
			assert.Equal(t, "Password reset email sent successfully", message)
			assert.Equal(t, "sent", status)
		} else {
			// If the server is strict about Content-Type, it might return 400
			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		}
	})

	t.Run("edge case: very long email address", func(t *testing.T) {
		// Clean state
		td.ClearAllTestData(t, ctx, testPool, redisClient)
		td.PurgeOTPQueue(t, ch)

		// Create a very long but technically valid email
		longEmail := strings.Repeat("a", 240) + "@example.com"

		// Test with very long email
		request := domain.RequestPasswordResetRequest{
			Email: longEmail,
		}
		req := td.NewRequestPasswordResetRequest(t, ctx, testServerURL, request)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Should return not found since user doesn't exist
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		errorMsg, statusCode := td.ParseErrorResponse(t, resp)
		assert.Contains(t, errorMsg, "email is not registered")
		assert.Equal(t, http.StatusNotFound, statusCode)
	})
}
