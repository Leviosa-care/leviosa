package auth_test

import (
	"context"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/Leviosa-care/authuser/internal/domain"
	td "github.com/Leviosa-care/authuser/test/helpers"

	"github.com/hengadev/encx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TEST=TestCheckEmailSendOTP make test-integration-auth-test

func TestCheckEmailSendOTP(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	// Setup RabbitMQ for OTP message verification
	ch, err := testMQConn.Channel()
	require.NoError(t, err)
	defer ch.Close()

	// Setup OTP queue and start consuming
	td.SetupOTPQueue(t, ch)
	msgs := td.ConsumeOTPMessages(t, ch)

	t.Run("should successfully request email verification for available email", func(t *testing.T) {
		// Clean state
		td.ClearAllTestData(t, ctx, testPool, redisClient)
		td.PurgeOTPQueue(t, ch)

		// Test with valid, available email
		request := domain.CheckEmailAvailabilityRequest{
			Email: "newuser@example.com",
		}
		req := td.NewCheckEmailSendOTPRequest(t, ctx, testServerURL, request)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert HTTP response
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		message, status := td.ParseCheckEmailSendOTPResponse(t, resp)
		assert.Equal(t, "Verification email sent successfully", message)
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
		otpEncx, err := td.GetOTPEncxByEmailHash(t, ctx, emailHash, redisClient, crypto)
		assert.NoError(t, err)
		otp, err := domain.DecryptOTPEncx(ctx, crypto, otpEncx)
		assert.NoError(t, err)

		// Verify message content
		td.VerifyOTPMessage(t, delivery, request.Email, otp.Code)
	})

	t.Run("should return conflict when email is already registered", func(t *testing.T) {
		// Clean state and insert existing user
		td.ClearAllTestData(t, ctx, testPool, redisClient)
		td.PurgeOTPQueue(t, ch)

		existingEmail := "existing@example.com"
		user := td.NewTestUser(t, existingEmail, "John", "Doe")
		userEncx, err := domain.ProcessUserEncx(ctx, crypto, user)
		require.NoError(t, err)
		err = td.InsertUserEncx(t, ctx, userEncx, testPool, crypto)
		require.NoError(t, err)

		// Test with existing email
		request := domain.CheckEmailAvailabilityRequest{
			Email: existingEmail,
		}
		req := td.NewCheckEmailSendOTPRequest(t, ctx, testServerURL, request)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert conflict response
		errorMsg, statusCode := td.ParseErrorResponse(t, resp)
		assert.Contains(t, errorMsg, "email is already registered")
		assert.Equal(t, http.StatusConflict, statusCode)

		// Verify no OTP was created
		emailBytes, err := encx.SerializeValue(request.Email)
		require.NoError(t, err)
		emailHash := crypto.HashBasic(ctx, emailBytes)
		exists := td.CheckOTPExists(t, ctx, emailHash, redisClient)
		assert.False(t, exists, "No OTP should be created for existing email")

		// Verify no RabbitMQ message was sent
		td.VerifyNoOTPMessage(t, msgs, 500*time.Millisecond)
	})

	t.Run("should return rate limit when OTP already exists", func(t *testing.T) {
		// Clean state
		td.ClearAllTestData(t, ctx, testPool, redisClient)
		td.PurgeOTPQueue(t, ch)

		email := "ratelimit@example.com"

		// Insert existing OTP
		existingOTP := td.NewTestOTP(email)
		existingOTPEncx, err := domain.ProcessOTPEncx(ctx, crypto, existingOTP)
		require.NoError(t, err)
		err = td.InsertOTPEncx(t, ctx, existingOTPEncx, redisClient, 10*time.Minute)
		require.NoError(t, err)

		// Test with same email (should hit rate limit)
		request := domain.CheckEmailAvailabilityRequest{
			Email: email,
		}
		req := td.NewCheckEmailSendOTPRequest(t, ctx, testServerURL, request)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert rate limit response
		assert.Equal(t, http.StatusTooManyRequests, resp.StatusCode)
		errorMsg, statusCode := td.ParseErrorResponse(t, resp)
		assert.Contains(t, errorMsg, "already active")
		assert.Equal(t, http.StatusTooManyRequests, statusCode)

		// Verify no new RabbitMQ message was sent
		td.VerifyNoOTPMessage(t, msgs, 500*time.Millisecond)
	})

	t.Run("should return bad request for invalid email", func(t *testing.T) {
		// Clean state
		td.ClearAllTestData(t, ctx, testPool, redisClient)
		td.PurgeOTPQueue(t, ch)

		// Test with invalid email
		request := domain.CheckEmailAvailabilityRequest{
			Email: "invalid-email-format",
		}
		req := td.NewCheckEmailSendOTPRequest(t, ctx, testServerURL, request)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert bad request response
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		errorMsg, statusCode := td.ParseErrorResponse(t, resp)
		assert.Contains(t, errorMsg, "invalid")
		assert.Equal(t, http.StatusBadRequest, statusCode)

		// Verify no OTP was created
		exists := td.CheckOTPExists(t, ctx, "invalid-hash", redisClient)
		assert.False(t, exists, "No OTP should be created for invalid email")

		// Verify no RabbitMQ message was sent
		td.VerifyNoOTPMessage(t, msgs, 500*time.Millisecond)
	})

	t.Run("should return bad request for empty email", func(t *testing.T) {
		// Clean state
		td.ClearAllTestData(t, ctx, testPool, redisClient)
		td.PurgeOTPQueue(t, ch)

		// Test with empty email
		request := domain.CheckEmailAvailabilityRequest{
			Email: "",
		}
		req := td.NewCheckEmailSendOTPRequest(t, ctx, testServerURL, request)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert bad request response
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		errorMsg, statusCode := td.ParseErrorResponse(t, resp)
		assert.Contains(t, errorMsg, "required")
		assert.Equal(t, http.StatusBadRequest, statusCode)

		// Verify no RabbitMQ message was sent
		td.VerifyNoOTPMessage(t, msgs, 500*time.Millisecond)
	})

	t.Run("should handle malformed JSON request", func(t *testing.T) {
		// Clean state
		td.ClearAllTestData(t, ctx, testPool, redisClient)
		td.PurgeOTPQueue(t, ch)

		// Create malformed JSON request
		malformedJSON := `{"email": "test@example.com", "invalid_field": true}`
		req, err := http.NewRequestWithContext(
			ctx,
			http.MethodPost,
			testServerURL+"/auth/email",
			strings.NewReader(malformedJSON),
		)
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert bad request response for unknown fields
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		errorMsg, statusCode := td.ParseErrorResponse(t, resp)
		assert.Contains(t, errorMsg, "unknown field")
		assert.Equal(t, http.StatusBadRequest, statusCode)

		// Verify no RabbitMQ message was sent
		td.VerifyNoOTPMessage(t, msgs, 500*time.Millisecond)
	})

	t.Run("should successfully handle concurrent requests for different emails", func(t *testing.T) {
		// Clean state
		td.ClearAllTestData(t, ctx, testPool, redisClient)
		td.PurgeOTPQueue(t, ch)

		// Test concurrent requests
		emails := []string{
			"concurrent1@example.com",
			"concurrent2@example.com",
			"concurrent3@example.com",
		}

		// Send concurrent requests
		responses := make(chan *http.Response, len(emails))
		for _, email := range emails {
			go func(e string) {
				request := domain.CheckEmailAvailabilityRequest{Email: e}
				req := td.NewCheckEmailSendOTPRequest(t, ctx, testServerURL, request)
				resp, err := client.Do(req)
				require.NoError(t, err)
				responses <- resp
			}(email)
		}

		// Collect and verify all responses
		successCount := 0
		for range len(emails) {
			resp := <-responses
			defer resp.Body.Close()

			if resp.StatusCode == http.StatusOK {
				successCount++
			}
		}

		// All requests should succeed
		assert.Equal(t, len(emails), successCount, "All concurrent requests should succeed")

		// Verify all OTPs were created
		for _, email := range emails {
			emailBytes, err := encx.SerializeValue(email)
			require.NoError(t, err)
			emailHash := crypto.HashBasic(ctx, emailBytes)
			exists := td.CheckOTPExists(t, ctx, emailHash, redisClient)
			assert.True(t, exists, "OTP should exist for email: %s", email)
		}

		// Verify all RabbitMQ messages were sent (allow some time for async processing)
		messagesReceived := 0
		timeout := time.After(3 * time.Second)
		for messagesReceived < len(emails) {
			select {
			case delivery := <-msgs:
				messagesReceived++
				delivery.Ack(false)
			case <-timeout:
				t.Fatalf("Timeout waiting for messages. Expected %d, received %d", len(emails), messagesReceived)
			}
		}
		assert.Equal(t, len(emails), messagesReceived, "All RabbitMQ messages should be received")
	})
}
