package auth_test

import (
	"context"
	"net/http"
	"strings"
	"testing"
	"time"

	aggregatorHandler "github.com/Leviosa-care/authuser/internal/adapters/http/auth"
	sessionRepository "github.com/Leviosa-care/authuser/internal/adapters/redis/session"
	"github.com/Leviosa-care/authuser/internal/domain"

	td "github.com/Leviosa-care/authuser/test/helpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TEST=TestValidatePasswordResetOTP make test-integration-auth-test

func TestValidatePasswordResetOTP(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	// No RabbitMQ setup needed since we're injecting OTPs directly

	t.Run("should successfully validate password reset OTP and return reset token", func(t *testing.T) {
		// Clean state and insert existing user
		td.ClearAllTestData(t, ctx, testPool, testClient)

		existingEmail := "reset@example.com"
		td.InsertTestUser(t, ctx, existingEmail, "Reset", "User", testPool, crypto)

		// Directly create and insert OTP for testing (no HTTP request needed)
		otp := td.NewValidOTP(existingEmail)
		emailHash := crypto.HashBasic(ctx, []byte(existingEmail))
		otp.EmailHash = emailHash

		// Encrypt OTP before storing
		err := crypto.ProcessStruct(ctx, otp)
		require.NoError(t, err)

		// Store OTP in Redis with TTL
		ttl := 10 * time.Minute
		td.InsertOTP(t, ctx, otp, testClient, ttl)

		// Step 2: Validate password reset OTP
		validateRequest := domain.ValidatePasswordResetOTPRequest{
			Email: existingEmail,
			Code:  otp.Code,
		}
		req := td.NewValidatePasswordResetOTPRequest(t, ctx, testServerURL, validateRequest)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert HTTP response
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		message, status, expiresAt := td.ParseValidatePasswordResetOTPResponse(t, resp)
		assert.Equal(t, "Password reset OTP validated successfully", message)
		assert.Equal(t, "validated", status)
		assert.NotEmpty(t, expiresAt)

		// Verify password reset token cookie was set
		resetTokenCookie := td.GetPasswordResetTokenCookie(t, resp)
		assert.NotEmpty(t, resetTokenCookie.Value, "Password reset token should not be empty")
		assert.Equal(t, "/auth/password/reset/confirm", resetTokenCookie.Path, "Cookie path should be restricted to confirm endpoint")
		assert.True(t, resetTokenCookie.HttpOnly, "Cookie should be HttpOnly")
		assert.True(t, resetTokenCookie.Secure, "Cookie should be Secure")
		assert.Equal(t, http.SameSiteStrictMode, resetTokenCookie.SameSite, "Cookie should use SameSite=Strict")

		// Verify OTP was consumed (no longer exists in Redis)
		exists := td.CheckOTPExists(t, ctx, emailHash, testClient)
		assert.False(t, exists, "OTP should be consumed and no longer exist in Redis")

		// Verify reset token was stored in Redis using repository helper
		tokenHash := crypto.HashBasic(ctx, []byte(resetTokenCookie.Value))
		resetSessionKey := sessionRepository.FormatResetSessionKey(tokenHash)
		storedEmail, err := testClient.Get(ctx, resetSessionKey).Result()
		require.NoError(t, err, "Reset session should exist in Redis")
		assert.Equal(t, existingEmail, storedEmail, "Stored email should match")

		// Verify reset token has proper TTL (should be around 15 minutes)
		ttl, err = testClient.TTL(ctx, resetSessionKey).Result()
		require.NoError(t, err)
		assert.True(t, ttl > 13*time.Minute, "Reset token TTL should be around 15 minutes")
		assert.True(t, ttl <= 15*time.Minute, "Reset token TTL should not exceed 15 minutes")
	})

	t.Run("should return not found when OTP does not exist", func(t *testing.T) {
		// Clean state and insert existing user
		td.ClearAllTestData(t, ctx, testPool, testClient)

		existingEmail := "noopt@example.com"
		td.InsertTestUser(t, ctx, existingEmail, "No", "OTP", testPool, crypto)

		// Try to validate OTP without requesting one first
		validateRequest := domain.ValidatePasswordResetOTPRequest{
			Email: existingEmail,
			Code:  "123456",
		}
		req := td.NewValidatePasswordResetOTPRequest(t, ctx, testServerURL, validateRequest)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert not found response
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		errorMsg, statusCode := td.ParseErrorResponse(t, resp)
		assert.Contains(t, strings.ToLower(errorMsg), "not found")
		assert.Equal(t, http.StatusNotFound, statusCode)

		// Verify no reset token cookie was set
		for _, cookie := range resp.Cookies() {
			assert.NotEqual(t, aggregatorHandler.PasswordResetTokenCookieName, cookie.Name, "No reset token cookie should be set")
		}
	})

	t.Run("should return unauthorized for wrong OTP code", func(t *testing.T) {
		// Clean state and insert existing user
		td.ClearAllTestData(t, ctx, testPool, testClient)

		existingEmail := "wrongcode@example.com"
		td.InsertTestUser(t, ctx, existingEmail, "Wrong", "Code", testPool, crypto)

		// Directly create and insert OTP for testing
		otp := td.NewValidOTP(existingEmail)
		emailHash := crypto.HashBasic(ctx, []byte(existingEmail))
		otp.EmailHash = emailHash

		// Encrypt OTP before storing
		err := crypto.ProcessStruct(ctx, otp)
		require.NoError(t, err)

		// Store OTP in Redis with TTL
		ttl := 10 * time.Minute
		td.InsertOTP(t, ctx, otp, testClient, ttl)

		// Step 2: Try to validate with wrong OTP code
		validateRequest := domain.ValidatePasswordResetOTPRequest{
			Email: existingEmail,
			Code:  "000000", // Wrong code
		}
		req := td.NewValidatePasswordResetOTPRequest(t, ctx, testServerURL, validateRequest)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert unauthorized response
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		errorMsg, statusCode := td.ParseErrorResponse(t, resp)
		assert.Contains(t, strings.ToLower(errorMsg), "mismatch")
		assert.Equal(t, http.StatusUnauthorized, statusCode)

		// Verify no reset token cookie was set
		for _, cookie := range resp.Cookies() {
			assert.NotEqual(t, aggregatorHandler.PasswordResetTokenCookieName, cookie.Name, "No reset token cookie should be set")
		}

		// Verify OTP still exists (failed attempt)
		emailHash = crypto.HashBasic(ctx, []byte(existingEmail))
		exists := td.CheckOTPExists(t, ctx, emailHash, testClient)
		assert.True(t, exists, "OTP should still exist after failed validation")
	})

	t.Run("should return bad request for invalid email format", func(t *testing.T) {
		// Clean state
		td.ClearAllTestData(t, ctx, testPool, testClient)

		// Test with invalid email
		validateRequest := domain.ValidatePasswordResetOTPRequest{
			Email: "invalid-email",
			Code:  "123456",
		}
		req := td.NewValidatePasswordResetOTPRequest(t, ctx, testServerURL, validateRequest)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert bad request response
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		errorMsg, statusCode := td.ParseErrorResponse(t, resp)
		assert.Contains(t, strings.ToLower(errorMsg), "email")
		assert.Equal(t, http.StatusBadRequest, statusCode)
	})

	t.Run("should return bad request for invalid OTP code format", func(t *testing.T) {
		// Clean state
		td.ClearAllTestData(t, ctx, testPool, testClient)

		// Test with invalid OTP code (non-numeric)
		validateRequest := domain.ValidatePasswordResetOTPRequest{
			Email: "test@example.com",
			Code:  "ABC123", // Non-numeric code
		}
		req := td.NewValidatePasswordResetOTPRequest(t, ctx, testServerURL, validateRequest)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert bad request response
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		errorMsg, statusCode := td.ParseErrorResponse(t, resp)
		assert.Contains(t, strings.ToLower(errorMsg), "code")
		assert.Equal(t, http.StatusBadRequest, statusCode)
	})

	t.Run("should return gone for expired OTP", func(t *testing.T) {
		// Clean state and insert existing user
		td.ClearAllTestData(t, ctx, testPool, testClient)

		existingEmail := "expired@example.com"
		td.InsertTestUser(t, ctx, existingEmail, "Expired", "OTP", testPool, crypto)

		// Directly create and insert an expired OTP for testing
		otp := td.NewExpiredOTP(existingEmail)
		emailHash := crypto.HashBasic(ctx, []byte(existingEmail))
		otp.EmailHash = emailHash

		// Encrypt OTP before storing
		err := crypto.ProcessStruct(ctx, otp)
		require.NoError(t, err)

		// Store already-expired OTP in Redis with short TTL to simulate expiration
		ttl := 1 * time.Millisecond
		td.InsertOTP(t, ctx, otp, testClient, ttl)

		// Wait for expiration
		time.Sleep(10 * time.Millisecond)

		// Step 2: Try to validate expired OTP
		validateRequest := domain.ValidatePasswordResetOTPRequest{
			Email: existingEmail,
			Code:  otp.Code,
		}
		req := td.NewValidatePasswordResetOTPRequest(t, ctx, testServerURL, validateRequest)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert gone response (OTP expired)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode) // Expired OTP results in not found
		errorMsg, statusCode := td.ParseErrorResponse(t, resp)
		assert.Contains(t, strings.ToLower(errorMsg), "not found")
		assert.Equal(t, http.StatusNotFound, statusCode)
	})

	t.Run("should handle concurrent OTP validation attempts", func(t *testing.T) {
		// Clean state and insert existing user
		td.ClearAllTestData(t, ctx, testPool, testClient)

		existingEmail := "concurrent@example.com"
		td.InsertTestUser(t, ctx, existingEmail, "Concurrent", "User", testPool, crypto)

		// Directly create and insert OTP for testing
		otp := td.NewValidOTP(existingEmail)
		emailHash := crypto.HashBasic(ctx, []byte(existingEmail))
		otp.EmailHash = emailHash

		// Encrypt OTP before storing
		err := crypto.ProcessStruct(ctx, otp)
		require.NoError(t, err)

		// Store OTP in Redis with TTL
		ttl := 10 * time.Minute
		td.InsertOTP(t, ctx, otp, testClient, ttl)

		// Step 2: Make concurrent validation requests
		validateRequest := domain.ValidatePasswordResetOTPRequest{
			Email: existingEmail,
			Code:  otp.Code,
		}

		// Create two concurrent requests
		req1 := td.NewValidatePasswordResetOTPRequest(t, ctx, testServerURL, validateRequest)
		req2 := td.NewValidatePasswordResetOTPRequest(t, ctx, testServerURL, validateRequest)

		// Execute concurrently
		results := make(chan *http.Response, 2)
		go func() {
			resp, _ := client.Do(req1)
			results <- resp
		}()
		go func() {
			resp, _ := client.Do(req2)
			results <- resp
		}()

		// Collect responses
		resp1 := <-results
		resp2 := <-results
		defer resp1.Body.Close()
		defer resp2.Body.Close()

		// One should succeed, one should fail (OTP already consumed)
		successCount := 0
		conflictCount := 0

		if resp1.StatusCode == http.StatusOK {
			successCount++
		} else if resp1.StatusCode == http.StatusConflict {
			conflictCount++
		}

		if resp2.StatusCode == http.StatusOK {
			successCount++
		} else if resp2.StatusCode == http.StatusConflict {
			conflictCount++
		}

		// Exactly one should succeed, one should be conflicted (already consumed)
		assert.Equal(t, 1, successCount, "Exactly one request should succeed")
		assert.Equal(t, 1, conflictCount, "Exactly one request should be conflicted")

		// Verify OTP is consumed
		exists := td.CheckOTPExists(t, ctx, emailHash, testClient)
		assert.False(t, exists, "OTP should be consumed")
	})

	t.Run("should handle malformed JSON request", func(t *testing.T) {
		// Clean state
		td.ClearAllTestData(t, ctx, testPool, testClient)

		// Create request with malformed JSON
		req, err := http.NewRequestWithContext(
			ctx,
			http.MethodPost,
			testServerURL+"/auth/password/reset/validate",
			strings.NewReader(`{"email": "test@example.com", "code": }`),
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

	t.Run("should return bad request for empty fields", func(t *testing.T) {
		// Clean state
		td.ClearAllTestData(t, ctx, testPool, testClient)

		// Test with empty email and code
		validateRequest := domain.ValidatePasswordResetOTPRequest{
			Email: "",
			Code:  "",
		}
		req := td.NewValidatePasswordResetOTPRequest(t, ctx, testServerURL, validateRequest)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert bad request response
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		errorMsg, statusCode := td.ParseErrorResponse(t, resp)
		assert.Contains(t, strings.ToLower(errorMsg), "email")
		assert.Equal(t, http.StatusBadRequest, statusCode)
	})
}
