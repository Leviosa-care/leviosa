package auth_test

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	authEndpoints "github.com/Leviosa-care/leviosa/backend/internal/authuser/interface/auth"
	ck "github.com/Leviosa-care/leviosa/backend/internal/common/auth/cookies"
	"github.com/Leviosa-care/leviosa/backend/internal/common/auth/session"
	"github.com/Leviosa-care/leviosa/backend/internal/common/contracts/identity"
	td "github.com/Leviosa-care/leviosa/backend/test/helpers"

	"github.com/google/uuid"
	"github.com/hengadev/encx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestValidateOTPCreatePendingUser TEST_PATH=test/integration/authuser/auth/validate_otp_create_pending_user_test.go

func TestValidateOTPCreatePendingUser(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("should successfully validate OTP and create pending user", func(t *testing.T) {
		// Clean state
		td.ClearAuthTestData(t, ctx, testPool, redisClient)

		email := "newuser@example.com"

		// First, create an OTP for this email
		otp := td.NewTestOTP(email)

		otpEncx, err := domain.ProcessOTPEncx(ctx, crypto, otp)
		require.NoError(t, err)

		err = td.InsertOTPEncx(t, ctx, otpEncx, redisClient, 10*time.Minute)
		require.NoError(t, err)

		// Test validation with correct code
		request := domain.ValidateOTPRequest{
			Email: email,
			Code:  otp.Code,
		}
		req := td.NewValidateOTPCreatePendingUserRequest(t, ctx, testServerURL, request)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert HTTP response
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
		message, status := td.ParseValidateOTPCreatePendingUserResponse(t, resp)
		assert.Equal(t, "Pending user created successfully", message)
		assert.Equal(t, "created", status)

		// Verify access token cookie is set using raw cookie parsing
		cookies := resp.Cookies()
		var accessTokenCookie *http.Cookie
		for _, cookie := range cookies {
			if cookie.Name == ck.AccessTokenCookieName {
				accessTokenCookie = cookie
				break
			}
		}
		assert.NotNil(t, accessTokenCookie, "Access token cookie should be set")

		// Verify cookie attributes
		assert.True(t, accessTokenCookie.HttpOnly, "Cookie should be HttpOnly")
		assert.True(t, accessTokenCookie.Secure, "Cookie should be Secure")
		assert.Equal(t, http.SameSiteStrictMode, accessTokenCookie.SameSite, "Cookie should have SameSite=Strict")
		assert.Equal(t, "/", accessTokenCookie.Path, "Cookie path should be /")
		assert.NotEmpty(t, accessTokenCookie.Value, "Cookie should have a value (token)")

		// Verify refresh token cookie is set
		var refreshTokenCookie *http.Cookie
		for _, cookie := range cookies {
			if cookie.Name == ck.RefreshTokenCookieName {
				refreshTokenCookie = cookie
				break
			}
		}
		assert.NotNil(t, refreshTokenCookie, "Refresh token cookie should be set")

		// Verify refresh cookie attributes
		assert.True(t, refreshTokenCookie.HttpOnly, "Refresh cookie should be HttpOnly")
		assert.True(t, refreshTokenCookie.Secure, "Refresh cookie should be Secure")
		assert.Equal(t, http.SameSiteStrictMode, refreshTokenCookie.SameSite, "Refresh cookie should have SameSite=Strict")
		assert.Equal(t, ck.RefreshEndpoint, refreshTokenCookie.Path, "Refresh cookie path should be /auth/refresh")
		assert.NotEmpty(t, refreshTokenCookie.Value, "Refresh cookie should have a value (token)")

		// Verify session exists in Redis using raw Redis queries
		accessTokenBytes, err := encx.SerializeValue(accessTokenCookie.Value)
		assert.NoError(t, err)
		accessTokenHash := crypto.HashBasic(ctx, accessTokenBytes)
		accessTokenKey := session.FormatAccessTokenKey(accessTokenHash)
		accessSessionIDStr, err := redisClient.Get(ctx, accessTokenKey).Result()
		assert.NoError(t, err, "Token mapping should exist in Redis")
		assert.NoError(t, uuid.Validate(accessSessionIDStr))

		// Verify refresh token exists in Redis
		refreshTokenBytes, err := encx.SerializeValue(refreshTokenCookie.Value)
		assert.NoError(t, err)
		refreshTokenHash := crypto.HashBasic(ctx, refreshTokenBytes)
		refreshTokenKey := session.FormatRefreshTokenKey(refreshTokenHash)
		refreshSessionIDStr, err := redisClient.Get(ctx, refreshTokenKey).Result()
		assert.NoError(t, err, "Refresh token mapping should exist in Redis")

		// Verify that access and refresh tokens store the same session ID value in Redis
		assert.Equal(t, accessSessionIDStr, refreshSessionIDStr, "Refresh token should map to same session")

		// Verify session data exists using raw Redis queries
		sessionKey := session.FormatSessionKey(accessSessionIDStr)
		// sessionBytes, err := redisClient.Get(ctx, sessionKey).Result()
		sessionBytes, err := redisClient.Get(ctx, sessionKey).Bytes()
		assert.NoError(t, err, "Session data should exist in Redis")
		assert.NotEmpty(t, sessionBytes, "Session data should not be empty")

		// Verify session data integrity by decrypting
		var sessionEncx session.SessionEncx
		err = json.Unmarshal(sessionBytes, &sessionEncx)
		assert.NoError(t, err, "Failed to unmarshal SessionEncx")

		// Decrypt the SessionEncx to get the original session
		s, err := session.DecryptSessionEncx(ctx, crypto, &sessionEncx)
		assert.NoError(t, err, "Failed to decrypt session")

		assert.Equal(t, session.SessionPending, s.State, "Session should be pending")
		assert.NotEmpty(t, s.UserID, "Session should have user ID")
		assert.NotEmpty(t, sessionEncx.AccessTokenHash, "Session should have access token hash")
		assert.NotEmpty(t, sessionEncx.RefreshTokenHash, "Session should have refresh token hash")
		assert.Equal(t, s.Role, identity.Visitor, "Session should have a user with role 'visitor'")

		// Verify pending session has shorter expiration time (30 minutes instead of 24 hours)
		expectedExpiry := s.CreatedAt.Add(session.PendingSessionDuration)
		timeDiff := time.Duration(math.Abs(float64(s.ExpiresAt.Sub(expectedExpiry))))
		assert.LessOrEqual(t, timeDiff, 5*time.Second, "Pending session should expire in ~30 minutes, not 24 hours")

		// Verify OTP was consumed (should be deleted)
		emailBytes, err := encx.SerializeValue(email)
		assert.NoError(t, err)
		emailHash := crypto.HashBasic(ctx, emailBytes)
		exists := td.CheckOTPExists(t, ctx, emailHash, redisClient)
		assert.False(t, exists, "OTP should be consumed and deleted")

		// Verify user was created using raw SQL
		userExists := td.CheckUserExistsByEmailHashSQL(t, ctx, emailHash, testPool)
		assert.True(t, userExists, "User should exist in database")

		// Verify user has correct state and data using raw SQL
		userState := td.GetUserStateSQL(t, ctx, emailHash, testPool)
		assert.Equal(t, domain.Unverified, userState)

		// Verify encrypted fields are properly populated
		hasEncryptedFields := td.CheckUserHasEncryptedFieldsSQL(t, ctx, emailHash, testPool)
		assert.True(t, hasEncryptedFields, "User should have encrypted fields populated")
	})

	t.Run("should return unauthorized for invalid OTP code", func(t *testing.T) {
		// Clean state
		td.ClearAuthTestData(t, ctx, testPool, redisClient)

		email := "invalidcode@example.com"

		// Create an OTP for this email
		otp := td.NewTestOTP(email)
		otpEncx, err := domain.ProcessOTPEncx(ctx, crypto, otp)
		require.NoError(t, err)

		err = td.InsertOTPEncx(t, ctx, otpEncx, redisClient, 10*time.Minute)
		require.NoError(t, err)

		// Test validation with wrong code
		request := domain.ValidateOTPRequest{
			Email: email,
			Code:  "999999", // Wrong code
		}
		req := td.NewValidateOTPCreatePendingUserRequest(t, ctx, testServerURL, request)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert unauthorized response
		errorMsg, statusCode := td.ParseErrorResponse(t, resp)
		assert.Contains(t, errorMsg, "mismatch")
		assert.Equal(t, http.StatusUnauthorized, statusCode)

		// Verify no access cookie is set
		cookies := resp.Cookies()
		var accessTokenCookie *http.Cookie
		for _, cookie := range cookies {
			if cookie.Name == ck.AccessTokenCookieName {
				accessTokenCookie = cookie
				break
			}
		}
		assert.Nil(t, accessTokenCookie, "No access token cookie should be set on validation failure")

		// Verify no refresh token cookie is set
		var refreshTokenCookie *http.Cookie
		for _, cookie := range cookies {
			if cookie.Name == ck.RefreshTokenCookieName {
				refreshTokenCookie = cookie
				break
			}
		}
		assert.Nil(t, refreshTokenCookie, "No refresh token cookie should be set on validation failure")

		// Verify no sessions exist in Redis using raw Redis queries
		sessionKeys, err := redisClient.Keys(ctx, fmt.Sprintf("%s*", session.SessionKeyPrefix)).Result()
		require.NoError(t, err)
		assert.Empty(t, sessionKeys, "No session keys should exist in Redis after failure")

		tokenKeys, err := redisClient.Keys(ctx, fmt.Sprintf("%s*", session.AccessTokenKeyPrefix)).Result()
		require.NoError(t, err)
		assert.Empty(t, tokenKeys, "No token keys should exist in Redis after failure")

		refreshTokenKeys, err := redisClient.Keys(ctx, fmt.Sprintf("%s*", session.RefreshTokenKeyPrefix)).Result()
		require.NoError(t, err)
		assert.Empty(t, refreshTokenKeys, "No refresh token keys should exist in Redis after failure")

		// Verify OTP still exists (should be incremented but not deleted)
		emailBytes, err := encx.SerializeValue(email)
		require.NoError(t, err)
		emailHash := crypto.HashBasic(ctx, emailBytes)
		exists := td.CheckOTPExists(t, ctx, emailHash, redisClient)
		assert.True(t, exists, "OTP should still exist after failed attempt")

		// Verify no user was created using raw SQL
		userExists := td.CheckUserExistsByEmailHashSQL(t, ctx, emailHash, testPool)
		assert.False(t, userExists, "No user should be created with invalid OTP")
	})

	t.Run("should return not found when OTP does not exist", func(t *testing.T) {
		// Clean state
		td.ClearAuthTestData(t, ctx, testPool, redisClient)

		email := "nonotp@example.com"

		// Test validation without creating OTP first
		request := domain.ValidateOTPRequest{
			Email: email,
			Code:  "123456",
		}
		req := td.NewValidateOTPCreatePendingUserRequest(t, ctx, testServerURL, request)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert not found response
		errorMsg, statusCode := td.ParseErrorResponse(t, resp)
		assert.Contains(t, errorMsg, "not found")
		assert.Equal(t, http.StatusNotFound, statusCode)

		// Verify no access cookie is set
		cookies := resp.Cookies()
		var accessTokenCookie *http.Cookie
		for _, cookie := range cookies {
			if cookie.Name == ck.AccessTokenCookieName {
				accessTokenCookie = cookie
				break
			}
		}
		assert.Nil(t, accessTokenCookie, "No access token cookie should be set when OTP not found")

		// Verify no refresh token cookie is set
		var refreshTokenCookie *http.Cookie
		for _, cookie := range cookies {
			if cookie.Name == ck.RefreshTokenCookieName {
				refreshTokenCookie = cookie
				break
			}
		}
		assert.Nil(t, refreshTokenCookie, "No refresh token cookie should be set when OTP not found")

		// Verify no sessions exist in Redis using raw Redis queries
		sessionKeys, err := redisClient.Keys(ctx, fmt.Sprintf("%s*", session.SessionKeyPrefix)).Result()
		require.NoError(t, err)
		assert.Empty(t, sessionKeys, "No session keys should exist in Redis when OTP not found")

		tokenKeys, err := redisClient.Keys(ctx, fmt.Sprintf("%s*", session.AccessTokenKeyPrefix)).Result()
		require.NoError(t, err)
		assert.Empty(t, tokenKeys, "No token keys should exist in Redis when OTP not found")

		refreshTokenKeys, err := redisClient.Keys(ctx, fmt.Sprintf("%s*", session.RefreshTokenKeyPrefix)).Result()
		require.NoError(t, err)
		assert.Empty(t, refreshTokenKeys, "No refresh token keys should exist in Redis when OTP not found")

		// Verify no user was created using raw SQL
		emailBytes, err := encx.SerializeValue(email)
		require.NoError(t, err)
		emailHash := crypto.HashBasic(ctx, emailBytes)
		userExists := td.CheckUserExistsByEmailHashSQL(t, ctx, emailHash, testPool)
		assert.False(t, userExists, "No user should be created without valid OTP")
	})

	t.Run("should return gone when OTP is expired", func(t *testing.T) {
		// Clean state
		td.ClearAuthTestData(t, ctx, testPool, redisClient)

		email := "expired@example.com"

		// Create an OTP with past expiration time but longer Redis TTL
		otp := td.NewTestOTP(email)
		otp.ExpiresAt = time.Now().Add(-1 * time.Minute) // Logically expired

		otpEncx, err := domain.ProcessOTPEncx(ctx, crypto, otp)
		require.NoError(t, err)

		err = td.InsertOTPEncx(t, ctx, otpEncx, redisClient, 10*time.Minute)
		require.NoError(t, err)

		// No sleep needed - OTP exists but is logically expired

		// Test validation with expired OTP
		request := domain.ValidateOTPRequest{
			Email: email,
			Code:  otp.Code,
		}
		req := td.NewValidateOTPCreatePendingUserRequest(t, ctx, testServerURL, request)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert gone response (expired)
		errorMsg, statusCode := td.ParseErrorResponse(t, resp)
		assert.Contains(t, errorMsg, "expired")
		assert.Equal(t, http.StatusGone, statusCode)

		// Verify no access cookie is set
		cookies := resp.Cookies()
		var accessTokenCookie *http.Cookie
		for _, cookie := range cookies {
			if cookie.Name == ck.AccessTokenCookieName {
				accessTokenCookie = cookie
				break
			}
		}
		assert.Nil(t, accessTokenCookie, "No access token cookie should be set when OTP is expired")

		// Verify no refresh token cookie is set
		var refreshTokenCookie *http.Cookie
		for _, cookie := range cookies {
			if cookie.Name == ck.RefreshTokenCookieName {
				refreshTokenCookie = cookie
				break
			}
		}
		assert.Nil(t, refreshTokenCookie, "No refresh token cookie should be set when OTP is expired")

		// Verify no sessions exist in Redis using raw Redis queries
		sessionKeys, err := redisClient.Keys(ctx, fmt.Sprintf("%s*", session.SessionKeyPrefix)).Result()
		require.NoError(t, err)
		assert.Empty(t, sessionKeys, "No session keys should exist in Redis when OTP is expired")

		tokenKeys, err := redisClient.Keys(ctx, fmt.Sprintf("%s*", session.AccessTokenKeyPrefix)).Result()
		require.NoError(t, err)
		assert.Empty(t, tokenKeys, "No token keys should exist in Redis when OTP is expired")

		refreshTokenKeys, err := redisClient.Keys(ctx, fmt.Sprintf("%s*", session.RefreshTokenKeyPrefix)).Result()
		require.NoError(t, err)
		assert.Empty(t, refreshTokenKeys, "No refresh token keys should exist in Redis when OTP is expired")

		// Verify no user was created using raw SQL
		emailBytes, err := encx.SerializeValue(email)
		require.NoError(t, err)
		emailHash := crypto.HashBasic(ctx, emailBytes)
		userExists := td.CheckUserExistsByEmailHashSQL(t, ctx, emailHash, testPool)
		assert.False(t, userExists, "No user should be created with expired OTP")
	})

	t.Run("should return too many requests when max attempts exceeded", func(t *testing.T) {
		// Clean state
		td.ClearAuthTestData(t, ctx, testPool, redisClient)

		email := "maxattempts@example.com"

		// Create an OTP with max attempts already reached
		otp := td.NewTestOTP(email)
		otp.Attempts = 5 // Assuming max attempts is 5

		otpEncx, err := domain.ProcessOTPEncx(ctx, crypto, otp)
		require.NoError(t, err)

		err = td.InsertOTPEncx(t, ctx, otpEncx, redisClient, 10*time.Minute)
		require.NoError(t, err)

		// Test validation (should be rate limited)
		request := domain.ValidateOTPRequest{
			Email: email,
			Code:  otp.Code,
		}
		req := td.NewValidateOTPCreatePendingUserRequest(t, ctx, testServerURL, request)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert rate limit response
		errorMsg, statusCode := td.ParseErrorResponse(t, resp)
		assert.Contains(t, errorMsg, "attempts")
		assert.Equal(t, http.StatusTooManyRequests, statusCode)

		// Verify no access cookie is set
		cookies := resp.Cookies()
		var accessTokenCookie *http.Cookie
		for _, cookie := range cookies {
			if cookie.Name == ck.AccessTokenCookieName {
				accessTokenCookie = cookie
				break
			}
		}
		assert.Nil(t, accessTokenCookie, "No access token cookie should be set when rate limited")

		// Verify no refresh token cookie is set
		var refreshTokenCookie *http.Cookie
		for _, cookie := range cookies {
			if cookie.Name == ck.RefreshTokenCookieName {
				refreshTokenCookie = cookie
				break
			}
		}
		assert.Nil(t, refreshTokenCookie, "No refresh token cookie should be set when rate limited")

		// Verify no sessions exist in Redis
		sessionKeys, err := redisClient.Keys(ctx, fmt.Sprintf("%s*", session.SessionKeyPrefix)).Result()
		require.NoError(t, err)
		assert.Empty(t, sessionKeys, "No session keys should exist in Redis when rate limited")

		tokenKeys, err := redisClient.Keys(ctx, fmt.Sprintf("%s*", session.AccessTokenKeyPrefix)).Result()
		require.NoError(t, err)
		assert.Empty(t, tokenKeys, "No token keys should exist in Redis when rate limited")

		refreshTokenKeys, err := redisClient.Keys(ctx, fmt.Sprintf("%s*", session.RefreshTokenKeyPrefix)).Result()
		require.NoError(t, err)
		assert.Empty(t, refreshTokenKeys, "No refresh token keys should exist in Redis when rate limited")

		// Verify no user was created using raw SQL
		emailBytes, err := encx.SerializeValue(email)
		require.NoError(t, err)
		emailHash := crypto.HashBasic(ctx, emailBytes)
		userExists := td.CheckUserExistsByEmailHashSQL(t, ctx, emailHash, testPool)
		assert.False(t, userExists, "No user should be created when rate limited")
	})

	t.Run("should return bad request for invalid email format", func(t *testing.T) {
		// Clean state
		td.ClearAuthTestData(t, ctx, testPool, redisClient)

		// Test with invalid email format
		request := domain.ValidateOTPRequest{
			Email: "invalid-email-format",
			Code:  "123456",
		}
		req := td.NewValidateOTPCreatePendingUserRequest(t, ctx, testServerURL, request)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert bad request response
		errorMsg, statusCode := td.ParseErrorResponse(t, resp)
		assert.Contains(t, errorMsg, "invalid")
		assert.Equal(t, http.StatusBadRequest, statusCode)
	})

	t.Run("should return bad request for invalid OTP format", func(t *testing.T) {
		// Clean state
		td.ClearAuthTestData(t, ctx, testPool, redisClient)

		// Test with invalid OTP code format
		request := domain.ValidateOTPRequest{
			Email: "valid@example.com",
			Code:  "12345", // Too short
		}
		req := td.NewValidateOTPCreatePendingUserRequest(t, ctx, testServerURL, request)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert bad request response
		errorMsg, statusCode := td.ParseErrorResponse(t, resp)
		assert.Contains(t, errorMsg, "validation failed")
		assert.Equal(t, http.StatusBadRequest, statusCode)
	})

	t.Run("should succeed when user already exists (race condition handling)", func(t *testing.T) {
		// Clean state
		td.ClearAuthTestData(t, ctx, testPool, redisClient)

		email := "existing@example.com"

		// Pre-create the user (simulating race condition)
		user := td.NewTestUser(t, email, "Existing", "User")
		userEncx, err := domain.ProcessUserEncx(ctx, crypto, user)
		require.NoError(t, err)
		err = td.InsertUserEncx(t, ctx, userEncx, testPool)
		require.NoError(t, err)

		// TODO: I need a session with that user ID in it, that is how I get the session

		// Create valid OTP for the same email
		otp := td.NewTestOTP(email)

		otpEncx, err := domain.ProcessOTPEncx(ctx, crypto, otp)
		require.NoError(t, err)

		err = td.InsertOTPEncx(t, ctx, otpEncx, redisClient, 10*time.Minute)
		require.NoError(t, err)

		// Test validation (should succeed despite existing user)
		request := domain.ValidateOTPRequest{
			Email: email,
			Code:  otp.Code,
		}
		req := td.NewValidateOTPCreatePendingUserRequest(t, ctx, testServerURL, request)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert success response (conflict is handled gracefully)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
		message, status := td.ParseValidateOTPCreatePendingUserResponse(t, resp)
		assert.Equal(t, "Pending user created successfully", message)
		assert.Equal(t, "created", status)

		// Verify access cookie is set even when user already exists
		cookies := resp.Cookies()
		var accessTokenCookie *http.Cookie
		for _, cookie := range cookies {
			if cookie.Name == ck.AccessTokenCookieName {
				accessTokenCookie = cookie
				break
			}
		}
		require.NotNil(t, accessTokenCookie, "Access token cookie should be set even when user exists")
		assert.NotEmpty(t, accessTokenCookie.Value, "Cookie should have a value (token)")

		// Verify refresh token cookie is set even when user already exists
		var refreshTokenCookie *http.Cookie
		for _, cookie := range cookies {
			if cookie.Name == ck.RefreshTokenCookieName {
				refreshTokenCookie = cookie
				break
			}
		}
		require.NotNil(t, refreshTokenCookie, "Refresh token cookie should be set even when user exists")
		assert.NotEmpty(t, refreshTokenCookie.Value, "Refresh cookie should have a value (token hash)")

		// Verify session exists in Redis using raw Redis queries
		accessToken := accessTokenCookie.Value
		accessTokenBytes, err := encx.SerializeValue(accessToken)
		require.NoError(t, err)
		accessTokenHash := crypto.HashBasic(ctx, accessTokenBytes)
		accessTokenKey := session.AccessTokenKeyPrefix + accessTokenHash
		accessSessionIDStr, err := redisClient.Get(ctx, accessTokenKey).Result()
		require.NoError(t, err, "Token mapping should exist in Redis even when user exists")

		// Verify refresh token exists in Redis even when user already exists
		refreshToken := refreshTokenCookie.Value
		refreshTokenBytes, err := encx.SerializeValue(refreshToken)
		require.NoError(t, err)
		refreshTokenHash := crypto.HashBasic(ctx, refreshTokenBytes)
		refreshTokenKey := session.RefreshTokenKeyPrefix + refreshTokenHash
		refreshSessionIDStr, err := redisClient.Get(ctx, refreshTokenKey).Result()
		require.NoError(t, err, "Refresh token mapping should exist in Redis even when user exists")
		assert.Equal(t, accessSessionIDStr, refreshSessionIDStr, "Refresh token should map to same session even when user exists")

		sessionKey := session.SessionKeyPrefix + accessSessionIDStr
		sessionDataBytes, err := redisClient.Get(ctx, sessionKey).Result()
		_ = sessionDataBytes
		require.NoError(t, err, "Session data should exist in Redis even when user exists")

		// Verify OTP was consumed
		emailBytes, err := encx.SerializeValue(email)
		require.NoError(t, err)
		emailHash := crypto.HashBasic(ctx, emailBytes)
		exists := td.CheckOTPExists(t, ctx, emailHash, redisClient)
		assert.False(t, exists, "OTP should be consumed even when user exists")

		// Verify user still exists using raw SQL
		userExists := td.CheckUserExistsByEmailHashSQL(t, ctx, emailHash, testPool)
		assert.True(t, userExists, "User should still exist")
	})

	t.Run("should handle concurrent validation attempts", func(t *testing.T) {
		// Clean state
		td.ClearAuthTestData(t, ctx, testPool, redisClient)

		email := "concurrent@example.com"

		// Create OTP
		otp := td.NewTestOTP(email)

		otpEncx, err := domain.ProcessOTPEncx(ctx, crypto, otp)
		require.NoError(t, err)

		err = td.InsertOTPEncx(t, ctx, otpEncx, redisClient, 10*time.Minute)
		require.NoError(t, err)

		// Launch concurrent validation attempts
		numAttempts := 3
		results := make(chan *http.Response, numAttempts)
		errors := make(chan error, numAttempts)

		for range numAttempts {
			go func() {
				request := domain.ValidateOTPRequest{
					Email: email,
					Code:  otp.Code,
				}
				req := td.NewValidateOTPCreatePendingUserRequest(t, ctx, testServerURL, request)
				resp, err := client.Do(req)
				if err != nil {
					t.Errorf("HTTP request failed: %v", err)
					errors <- err
					return
				}
				results <- resp
			}()
		}

		// Collect results
		successCount := 0
		requestsProcessed := 0

		// Process responses and errors
		for requestsProcessed < numAttempts {
			select {
			case resp := <-results:
				defer resp.Body.Close()
				requestsProcessed++

				t.Logf("Concurrent request %d response status: %d", requestsProcessed, resp.StatusCode)

				if resp.StatusCode == http.StatusCreated {
					successCount++
					t.Logf("Request %d succeeded", requestsProcessed)
				} else {
					errorMsg, statusCode := td.ParseErrorResponse(t, resp)
					t.Logf("Request %d failed with status %d: %s", requestsProcessed, statusCode, errorMsg)

					// Expect concurrent requests to get 409 Conflict (already consumed)
					if resp.StatusCode == http.StatusConflict {
						assert.Contains(t, errorMsg, "already consumed", "Concurrent requests should get 'already consumed' error")
					}
				}

			case err := <-errors:
				requestsProcessed++
				t.Errorf("Request %d had HTTP error: %v", requestsProcessed, err)
			}
		}

		// Debug: Log final success count
		t.Logf("Total successful requests: %d (expected: 1)", successCount)

		// Only one request should succeed due to OTP consumption
		assert.Equal(t, 1, successCount, "Only one concurrent request should succeed")

		// Verify user was created using raw SQL
		emailBytes, err := encx.SerializeValue(email)
		require.NoError(t, err)
		emailHash := crypto.HashBasic(ctx, emailBytes)
		userExists := td.CheckUserExistsByEmailHashSQL(t, ctx, emailHash, testPool)
		assert.True(t, userExists, "User should be created by successful request")

		// Verify OTP was consumed
		exists := td.CheckOTPExists(t, ctx, emailHash, redisClient)
		assert.False(t, exists, "OTP should be consumed")

		// Debug: Check how many sessions were created
		sessionKeys, err := redisClient.Keys(ctx, fmt.Sprintf("%s*", session.SessionKeyPrefix)).Result()
		require.NoError(t, err)
		t.Logf("Number of sessions created: %d (expected: 1)", len(sessionKeys))
		assert.Equal(t, 1, len(sessionKeys), "Only one session should be created")

		tokenKeys, err := redisClient.Keys(ctx, fmt.Sprintf("%s*", session.AccessTokenKeyPrefix)).Result()
		require.NoError(t, err)
		t.Logf("Number of token created: %d (expected: 1)", len(tokenKeys))
		assert.Equal(t, 1, len(tokenKeys), "Only one token should be created")
	})

	t.Run("should handle malformed JSON request", func(t *testing.T) {
		// Clean state
		td.ClearAuthTestData(t, ctx, testPool, redisClient)

		// Create malformed JSON request
		malformedJSON := `{"email": "test@example.com", "invalid_field": true}`
		req, err := http.NewRequestWithContext(
			ctx,
			http.MethodPost,
			testServerURL+authEndpoints.ValidateOTPCreatePendingEndpoint,
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
	})

	t.Run("should verify pending session characteristics and completion workflow", func(t *testing.T) {
		// Clean state
		td.ClearAuthTestData(t, ctx, testPool, redisClient)

		email := "completiontest@example.com"

		// Create and validate OTP to get a pending session
		otp := td.NewTestOTP(email)

		otpEncx, err := domain.ProcessOTPEncx(ctx, crypto, otp)
		require.NoError(t, err)

		err = td.InsertOTPEncx(t, ctx, otpEncx, redisClient, 10*time.Minute)
		require.NoError(t, err)

		// Validate OTP and create pending session
		request := domain.ValidateOTPRequest{
			Email: email,
			Code:  otp.Code,
		}
		req := td.NewValidateOTPCreatePendingUserRequest(t, ctx, testServerURL, request)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		// Extract access cookie
		cookies := resp.Cookies()
		var accessTokenCookie *http.Cookie
		for _, cookie := range cookies {
			if cookie.Name == ck.AccessTokenCookieName {
				accessTokenCookie = cookie
				break
			}
		}
		require.NotNil(t, accessTokenCookie, "Access token cookie should be set")

		// Extract refresh token cookie
		var refreshTokenCookie *http.Cookie
		for _, cookie := range cookies {
			if cookie.Name == ck.RefreshTokenCookieName {
				refreshTokenCookie = cookie
				break
			}
		}
		require.NotNil(t, refreshTokenCookie, "Refresh token cookie should be set")

		// Verify session in Redis
		accessToken := accessTokenCookie.Value
		accessTokenBytes, err := encx.SerializeValue(accessToken)
		require.NoError(t, err)
		accessTokenHash := crypto.HashBasic(ctx, accessTokenBytes)
		accessTokenKey := session.AccessTokenKeyPrefix + accessTokenHash
		accessSessionIDStr, err := redisClient.Get(ctx, accessTokenKey).Result()
		require.NoError(t, err)

		// Verify refresh token exists in Redis
		refreshToken := refreshTokenCookie.Value
		refreshTokenBytes, err := encx.SerializeValue(refreshToken)
		require.NoError(t, err)
		refreshTokenHash := crypto.HashBasic(ctx, refreshTokenBytes)
		refreshTokenKey := session.RefreshTokenKeyPrefix + refreshTokenHash
		refreshSessionIDStr, err := redisClient.Get(ctx, refreshTokenKey).Result()
		require.NoError(t, err, "Refresh token mapping should exist in Redis")
		assert.Equal(t, accessSessionIDStr, refreshSessionIDStr, "Refresh token should map to same session")

		sessionKey := session.SessionKeyPrefix + accessSessionIDStr
		// sessionBytes, err := redisClient.Get(ctx, sessionKey).Result()
		sessionBytes, err := redisClient.Get(ctx, sessionKey).Bytes()
		require.NoError(t, err)

		// Decode and verify session characteristics
		var sessionEncx session.SessionEncx
		err = json.Unmarshal(sessionBytes, &sessionEncx)
		require.NoError(t, err, "Failed to unmarshal SessionEncx")

		// Decrypt the SessionEncx to get the original session
		s, err := session.DecryptSessionEncx(ctx, crypto, &sessionEncx)
		require.NoError(t, err, "Failed to decrypt session")

		// sessionData := []byte(sessionBytes)
		// session := td.DecodeSessionWithDecryption(t, sessionData, crypto)

		// Verify it's a pending session with correct expiration
		assert.Equal(t, session.SessionPending, s.State, "Session should be pending")
		assert.Nil(t, s.CompletedAt, "Session should not have completion timestamp yet")

		// Verify session expires in ~30 minutes (pending duration), not 24 hours
		expectedExpiry := s.CreatedAt.Add(session.PendingSessionDuration)
		timeDiff := time.Duration(math.Abs(float64(s.ExpiresAt.Sub(expectedExpiry))))
		assert.LessOrEqual(t, timeDiff, 5*time.Second, "Pending session should have 30-minute expiration")

		// Verify session does NOT expire in 24 hours (active session duration)
		activeExpiry := s.CreatedAt.Add(session.ActiveSessionDuration)
		activeTimeDiff := time.Duration(math.Abs(float64(s.ExpiresAt.Sub(activeExpiry))))
		assert.Greater(t, activeTimeDiff, 20*time.Hour, "Pending session should NOT have 24-hour expiration")

		t.Logf("Session created at: %s", s.CreatedAt.Format(time.RFC3339))
		t.Logf("Session expires at: %s", s.ExpiresAt.Format(time.RFC3339))
		t.Logf("Expected pending expiry: %s", expectedExpiry.Format(time.RFC3339))
		t.Logf("Time difference: %s", timeDiff.String())
		t.Logf("Pending session duration: %s", session.PendingSessionDuration.String())
		t.Logf("Active session duration: %s", session.ActiveSessionDuration.String())
	})
}
