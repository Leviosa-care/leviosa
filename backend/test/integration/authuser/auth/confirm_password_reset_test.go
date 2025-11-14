package auth_test

import (
	"context"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	sessionRepository "github.com/Leviosa-care/leviosa/backend/internal/authuser/infrastructure/redis/session"
	aggregatorHandler "github.com/Leviosa-care/leviosa/backend/internal/authuser/interface/auth"
	session "github.com/Leviosa-care/leviosa/backend/internal/common/auth/session"
	td "github.com/Leviosa-care/leviosa/backend/test/helpers"

	"github.com/hengadev/encx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestConfirmPasswordReset TEST_PATH=test/integration/authuser/auth/confirm_password_reset_test.go

func TestConfirmPasswordReset(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	// No RabbitMQ setup needed since we're directly injecting data

	t.Run("should successfully confirm password reset with valid token", func(t *testing.T) {
		// Clean state and insert existing user
		td.ClearAuthTestData(t, ctx, testPool, redisClient)

		existingEmail := "resetconfirm@example.com"
		oldPassword := td.GenerateStrongPassword(t)
		newPassword := td.GenerateStrongPassword(t)

		// Insert user with old password
		user := td.NewTestUser(t, existingEmail, "Reset", "Confirm")
		user.Password = oldPassword
		user.State = domain.Active
		userEncx, err := domain.ProcessUserEncx(ctx, crypto, user)
		require.NoError(t, err)
		err = td.InsertUserEncx(t, ctx, userEncx, testPool)
		require.NoError(t, err)

		// Generate a reset token and store reset session directly
		resetToken, err := session.GenerateToken()
		require.NoError(t, err)
		resetTokenBytes, err := encx.SerializeValue(resetToken)
		require.NoError(t, err)
		tokenHash := crypto.HashBasic(ctx, resetTokenBytes)

		// Store reset session in Redis directly
		resetSessionKey := sessionRepository.FormatResetSessionKey(tokenHash)
		ttl := 15 * time.Minute
		err = redisClient.Set(ctx, resetSessionKey, existingEmail, ttl).Err()
		require.NoError(t, err)

		// Create reset cookie for the request
		resetTokenCookie := &http.Cookie{
			Name:  aggregatorHandler.PasswordResetTokenCookieName,
			Value: resetToken,
		}

		// Confirm password reset
		confirmRequest := domain.ConfirmPasswordResetRequest{
			Token:       resetTokenCookie.Value,
			NewPassword: newPassword,
		}
		req := td.NewConfirmPasswordResetRequestWithCookie(t, ctx, testServerURL, confirmRequest, resetTokenCookie)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert HTTP response
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		message, status := td.ParseConfirmPasswordResetResponse(t, resp)
		assert.Equal(t, "Password reset completed successfully", message)
		assert.Equal(t, "completed", status)

		// Verify reset token cookie was cleared
		foundResetCookie := false
		for _, cookie := range resp.Cookies() {
			if cookie.Name == aggregatorHandler.PasswordResetTokenCookieName && cookie.MaxAge == -1 {
				foundResetCookie = true
				break
			}
		}
		assert.True(t, foundResetCookie, "Reset token cookie should be cleared")

		// Verify reset session was consumed (no longer exists in Redis)
		resetSessionTokenBytes, err := encx.SerializeValue(resetTokenCookie.Value)
		tokenHash = crypto.HashBasic(ctx, resetSessionTokenBytes)
		resetSessionKey = sessionRepository.FormatResetSessionKey(tokenHash)
		_, err = redisClient.Get(ctx, resetSessionKey).Result()
		assert.Error(t, err, "Reset session should be consumed and no longer exist")

		// Verify password was actually changed by attempting sign in with new password
		signInRequest := domain.SignInRequest{
			Email:    existingEmail,
			Password: newPassword,
		}
		signInReq := td.NewSignInRequest(t, ctx, testServerURL, signInRequest)
		signInResp, err := client.Do(signInReq)
		require.NoError(t, err)
		defer signInResp.Body.Close()
		assert.Equal(t, http.StatusCreated, signInResp.StatusCode, "Should be able to sign in with new password")

		// Verify cannot sign in with old password
		oldSignInRequest := domain.SignInRequest{
			Email:    existingEmail,
			Password: oldPassword,
		}
		oldSignInReq := td.NewSignInRequest(t, ctx, testServerURL, oldSignInRequest)
		oldSignInResp, err := client.Do(oldSignInReq)
		require.NoError(t, err)
		defer oldSignInResp.Body.Close()
		assert.Equal(t, http.StatusUnauthorized, oldSignInResp.StatusCode, "Should not be able to sign in with old password")
	})

	t.Run("should return not found for invalid reset token", func(t *testing.T) {
		// Clean state
		td.ClearAuthTestData(t, ctx, testPool, redisClient)

		newPassword := td.GenerateStrongPassword(t)

		// Try to confirm with invalid token
		confirmRequest := domain.ConfirmPasswordResetRequest{
			Token:       "invalid-token-123",
			NewPassword: newPassword,
		}
		req := td.NewConfirmPasswordResetRequest(t, ctx, testServerURL, confirmRequest)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert not found response
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		errorMsg, statusCode := td.ParseErrorResponse(t, resp)
		assert.Contains(t, strings.ToLower(errorMsg), "not found")
		assert.Equal(t, http.StatusNotFound, statusCode)
	})

	t.Run("should return not found for expired reset token", func(t *testing.T) {
		// Clean state and insert existing user
		td.ClearAuthTestData(t, ctx, testPool, redisClient)

		existingEmail := "expired@example.com"

		user := td.NewTestUser(t, existingEmail, "Reset", "Confirm")
		userEncx, err := domain.ProcessUserEncx(ctx, crypto, user)
		require.NoError(t, err)
		err = td.InsertUserEncx(t, ctx, userEncx, testPool)
		require.NoError(t, err)

		// Generate a reset token and store with very short TTL to simulate expiration
		resetToken, err := session.GenerateToken()
		require.NoError(t, err)
		tokenBytes, err := encx.SerializeValue(resetToken)
		require.NoError(t, err)
		tokenHash := crypto.HashBasic(ctx, tokenBytes)
		emailBytes, err := encx.SerializeValue(existingEmail)
		require.NoError(t, err)
		emailHash := crypto.HashBasic(ctx, emailBytes)

		// Store reset session in Redis with short TTL
		resetSessionKey := sessionRepository.FormatResetSessionKey(tokenHash)
		err = redisClient.Set(ctx, resetSessionKey, emailHash, 1*time.Millisecond).Err()
		require.NoError(t, err)

		// Wait for expiration
		time.Sleep(10 * time.Millisecond)

		// Create reset cookie for the request
		resetTokenCookie := &http.Cookie{
			Name:  aggregatorHandler.PasswordResetTokenCookieName,
			Value: resetToken,
		}

		newPassword := td.GenerateStrongPassword(t)

		// Step 3: Try to confirm with expired token
		confirmRequest := domain.ConfirmPasswordResetRequest{
			Token:       resetTokenCookie.Value,
			NewPassword: newPassword,
		}
		req := td.NewConfirmPasswordResetRequestWithCookie(t, ctx, testServerURL, confirmRequest, resetTokenCookie)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert not found response (expired token results in not found)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		errorMsg, statusCode := td.ParseErrorResponse(t, resp)
		assert.Contains(t, strings.ToLower(errorMsg), "not found")
		assert.Equal(t, http.StatusNotFound, statusCode)
	})

	t.Run("should return bad request for invalid password format", func(t *testing.T) {
		// Clean state and insert existing user
		td.ClearAuthTestData(t, ctx, testPool, redisClient)

		existingEmail := "weakpass@example.com"

		user := td.NewTestUser(t, existingEmail, "Reset", "Confirm")
		userEncx, err := domain.ProcessUserEncx(ctx, crypto, user)
		require.NoError(t, err)
		err = td.InsertUserEncx(t, ctx, userEncx, testPool)
		require.NoError(t, err)

		// Generate a reset token and store reset session directly
		resetToken, err := session.GenerateToken()
		require.NoError(t, err)
		tokenBytes, err := encx.SerializeValue(resetToken)
		require.NoError(t, err)
		tokenHash := crypto.HashBasic(ctx, tokenBytes)
		emailBytes, err := encx.SerializeValue(existingEmail)
		require.NoError(t, err)
		emailHash := crypto.HashBasic(ctx, emailBytes)

		// Store reset session in Redis directly
		resetSessionKey := sessionRepository.FormatResetSessionKey(tokenHash)
		ttl := 15 * time.Minute
		err = redisClient.Set(ctx, resetSessionKey, emailHash, ttl).Err()
		require.NoError(t, err)

		// Create reset cookie for the request
		resetTokenCookie := &http.Cookie{
			Name:  aggregatorHandler.PasswordResetTokenCookieName,
			Value: resetToken,
		}

		// Step 3: Try to confirm with weak password
		confirmRequest := domain.ConfirmPasswordResetRequest{
			Token:       resetTokenCookie.Value,
			NewPassword: "weak", // Invalid password
		}
		req := td.NewConfirmPasswordResetRequestWithCookie(t, ctx, testServerURL, confirmRequest, resetTokenCookie)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert bad request response
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		errorMsg, statusCode := td.ParseErrorResponse(t, resp)
		assert.Contains(t, strings.ToLower(errorMsg), "password")
		assert.Equal(t, http.StatusBadRequest, statusCode)

		// Verify reset token was not consumed (still exists)
		tokenBytes, err = encx.SerializeValue(resetTokenCookie.Value)
		tokenHash = crypto.HashBasic(ctx, tokenBytes)
		resetSessionKey = sessionRepository.FormatResetSessionKey(tokenHash)
		_, err = redisClient.Get(ctx, resetSessionKey).Result()
		assert.NoError(t, err, "Reset token should still exist after failed validation")
	})

	t.Run("should handle token in request body instead of cookie", func(t *testing.T) {
		// Clean state and insert existing user
		td.ClearAuthTestData(t, ctx, testPool, redisClient)

		existingEmail := "tokenbody@example.com"
		newPassword := td.GenerateStrongPassword(t)

		user := td.NewTestUser(t, existingEmail, "Token", "Body")
		userEncx, err := domain.ProcessUserEncx(ctx, crypto, user)
		require.NoError(t, err)
		err = td.InsertUserEncx(t, ctx, userEncx, testPool)
		require.NoError(t, err)

		// Generate a reset token and store reset session directly
		resetToken, err := session.GenerateToken()
		require.NoError(t, err)
		tokenBytes, err := encx.SerializeValue(resetToken)
		require.NoError(t, err)
		tokenHash := crypto.HashBasic(ctx, tokenBytes)

		// Store reset session in Redis directly
		resetSessionKey := sessionRepository.FormatResetSessionKey(tokenHash)
		ttl := 15 * time.Minute
		err = redisClient.Set(ctx, resetSessionKey, existingEmail, ttl).Err()
		require.NoError(t, err)

		// Confirm password reset with token in request body (no cookie)
		confirmRequest := domain.ConfirmPasswordResetRequest{
			Token:       resetToken,
			NewPassword: newPassword,
		}
		req := td.NewConfirmPasswordResetRequest(t, ctx, testServerURL, confirmRequest) // Don't add cookie
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Should still work with token in request body
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		message, status := td.ParseConfirmPasswordResetResponse(t, resp)
		assert.Equal(t, "Password reset completed successfully", message)
		assert.Equal(t, "completed", status)
	})

	t.Run("should return bad request for missing token", func(t *testing.T) {
		// Clean state
		td.ClearAuthTestData(t, ctx, testPool, redisClient)

		newPassword := td.GenerateStrongPassword(t)

		// Try to confirm without token
		confirmRequest := domain.ConfirmPasswordResetRequest{
			Token:       "", // Missing token
			NewPassword: newPassword,
		}
		req := td.NewConfirmPasswordResetRequest(t, ctx, testServerURL, confirmRequest)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert bad request response
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		errorMsg, statusCode := td.ParseErrorResponse(t, resp)
		assert.Contains(t, strings.ToLower(errorMsg), "token")
		assert.Equal(t, http.StatusBadRequest, statusCode)
	})

	t.Run("should handle concurrent confirmation attempts", func(t *testing.T) {
		// Clean state and insert existing user
		td.ClearAuthTestData(t, ctx, testPool, redisClient)

		existingEmail := "concurrent@example.com"
		newPassword := td.GenerateStrongPassword(t)

		user := td.NewTestUser(t, existingEmail, "Concurrent", "User")
		userEncx, err := domain.ProcessUserEncx(ctx, crypto, user)
		require.NoError(t, err)
		err = td.InsertUserEncx(t, ctx, userEncx, testPool)
		require.NoError(t, err)

		// Generate a reset token and store reset session directly
		resetToken, err := session.GenerateToken()
		require.NoError(t, err)
		tokenBytes, err := encx.SerializeValue(resetToken)
		require.NoError(t, err)
		tokenHash := crypto.HashBasic(ctx, tokenBytes)

		// Store reset session in Redis directly
		resetSessionKey := sessionRepository.FormatResetSessionKey(tokenHash)
		ttl := 15 * time.Minute
		err = redisClient.Set(ctx, resetSessionKey, existingEmail, ttl).Err()
		require.NoError(t, err)

		// Create reset cookie for the request
		resetTokenCookie := &http.Cookie{
			Name:  aggregatorHandler.PasswordResetTokenCookieName,
			Value: resetToken,
		}

		// Step 3: Make concurrent confirmation requests
		confirmRequest := domain.ConfirmPasswordResetRequest{
			Token:       resetTokenCookie.Value,
			NewPassword: newPassword,
		}

		req1 := td.NewConfirmPasswordResetRequestWithCookie(t, ctx, testServerURL, confirmRequest, resetTokenCookie)
		req2 := td.NewConfirmPasswordResetRequestWithCookie(t, ctx, testServerURL, confirmRequest, resetTokenCookie)

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

		// One should succeed, one should fail (token already consumed)
		successCount := 0
		notFoundCount := 0

		if resp1.StatusCode == http.StatusOK {
			successCount++
		} else if resp1.StatusCode == http.StatusNotFound {
			notFoundCount++
		}

		if resp2.StatusCode == http.StatusOK {
			successCount++
		} else if resp2.StatusCode == http.StatusNotFound {
			notFoundCount++
		}

		// Exactly one should succeed, one should fail
		assert.Equal(t, 1, successCount, "Exactly one request should succeed")
		assert.Equal(t, 1, notFoundCount, "Exactly one request should fail with not found")
	})

	t.Run("should handle malformed JSON request", func(t *testing.T) {
		// Clean state
		td.ClearAuthTestData(t, ctx, testPool, redisClient)

		// Create request with malformed JSON
		req, err := http.NewRequestWithContext(
			ctx,
			http.MethodPost,
			testServerURL+aggregatorHandler.PasswordResetEndpoint,
			strings.NewReader(`{"token": "test-token", "new_password": }`),
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
}
