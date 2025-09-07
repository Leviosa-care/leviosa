package auth_test

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/Leviosa-care/authuser/internal/domain"
	td "github.com/Leviosa-care/authuser/test/helpers"
	"github.com/google/uuid"

	ck "github.com/Leviosa-care/core/auth/cookies"
	sx "github.com/Leviosa-care/core/auth/session"
	"github.com/Leviosa-care/core/contracts/identity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TEST=TestValidateOTPCreatePendingUser make test-integration-auth-test

func TestValidateOTPCreatePendingUser(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("should successfully validate OTP and create pending user", func(t *testing.T) {
		// Clean state
		td.ClearAllTestData(t, ctx, testPool, testClient)

		email := "newuser@example.com"

		// First, create an OTP for this email
		otp := td.NewValidOTP(email)
		err := crypto.ProcessStruct(ctx, otp)
		require.NoError(t, err)
		td.InsertOTP(t, ctx, otp, testClient, 10*time.Minute)

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
		require.NotNil(t, accessTokenCookie, "Access token cookie should be set")

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
		require.NotNil(t, refreshTokenCookie, "Refresh token cookie should be set")

		// Verify refresh cookie attributes
		assert.True(t, refreshTokenCookie.HttpOnly, "Refresh cookie should be HttpOnly")
		assert.True(t, refreshTokenCookie.Secure, "Refresh cookie should be Secure")
		assert.Equal(t, http.SameSiteStrictMode, refreshTokenCookie.SameSite, "Refresh cookie should have SameSite=Strict")
		assert.Equal(t, ck.RefreshEndpoint, refreshTokenCookie.Path, "Refresh cookie path should be /auth/refresh")
		assert.NotEmpty(t, refreshTokenCookie.Value, "Refresh cookie should have a value (token)")

		// Verify session exists in Redis using raw Redis queries
		accessTokenHash := crypto.HashBasic(ctx, []byte(accessTokenCookie.Value))
		accessTokenKey := sx.FormatAccessTokenKey(accessTokenHash)
		accessSessionIDStr, err := testClient.Get(ctx, accessTokenKey).Result()
		require.NoError(t, err, "Token mapping should exist in Redis")
		require.NoError(t, uuid.Validate(accessSessionIDStr))

		// Verify refresh token exists in Redis
		refreshTokenHash := crypto.HashBasic(ctx, []byte(refreshTokenCookie.Value))
		refreshTokenKey := sx.FormatRefreshTokenKey(refreshTokenHash)
		refreshSessionIDStr, err := testClient.Get(ctx, refreshTokenKey).Result()
		require.NoError(t, err, "Refresh token mapping should exist in Redis")

		// Verify that access and refresh tokens store the same session ID value in Redis
		assert.Equal(t, accessSessionIDStr, refreshSessionIDStr, "Refresh token should map to same session")

		// Verify session data exists using raw Redis queries
		sessionKey := sx.FormatSessionKey(accessSessionIDStr)
		sessionDataBytes, err := testClient.Get(ctx, sessionKey).Result()
		require.NoError(t, err, "Session data should exist in Redis")
		assert.NotEmpty(t, sessionDataBytes, "Session data should not be empty")

		// Verify session data integrity by decrypting
		sessionData := []byte(sessionDataBytes)
		session := td.DecodeSessionWithDecryption(t, sessionData, crypto)
		assert.Equal(t, sx.SessionPending, session.State, "Session should be pending")
		assert.NotEmpty(t, session.UserID, "Session should have user ID")
		assert.NotEmpty(t, session.AccessTokenHash, "Session should have access token hash")
		assert.NotEmpty(t, session.RefreshTokenHash, "Session should have refresh token hash")
		assert.Equal(t, session.Role, identity.Visitor, "Session should have a user with role 'visitor'")

		// Verify pending session has shorter expiration time (30 minutes instead of 24 hours)
		expectedExpiry := session.CreatedAt.Add(sx.PendingSessionDuration)
		timeDiff := time.Duration(math.Abs(float64(session.ExpiresAt.Sub(expectedExpiry))))
		assert.LessOrEqual(t, timeDiff, 5*time.Second, "Pending session should expire in ~30 minutes, not 24 hours")

		// Verify OTP was consumed (should be deleted)
		emailHash := crypto.HashBasic(ctx, []byte(email))
		exists := td.CheckOTPExists(t, ctx, emailHash, testClient)
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
		td.ClearAllTestData(t, ctx, testPool, testClient)

		email := "invalidcode@example.com"

		// Create an OTP for this email
		otp := td.NewValidOTP(email)
		err := crypto.ProcessStruct(ctx, otp)
		require.NoError(t, err)
		td.InsertOTP(t, ctx, otp, testClient, 10*time.Minute)

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
		sessionKeys, err := testClient.Keys(ctx, fmt.Sprintf("%s*", sx.SessionKeyPrefix)).Result()
		require.NoError(t, err)
		assert.Empty(t, sessionKeys, "No session keys should exist in Redis after failure")

		tokenKeys, err := testClient.Keys(ctx, fmt.Sprintf("%s*", sx.AccessTokenKeyPrefix)).Result()
		require.NoError(t, err)
		assert.Empty(t, tokenKeys, "No token keys should exist in Redis after failure")

		refreshTokenKeys, err := testClient.Keys(ctx, fmt.Sprintf("%s*", sx.RefreshTokenKeyPrefix)).Result()
		require.NoError(t, err)
		assert.Empty(t, refreshTokenKeys, "No refresh token keys should exist in Redis after failure")

		// Verify OTP still exists (should be incremented but not deleted)
		emailHash := crypto.HashBasic(ctx, []byte(email))
		exists := td.CheckOTPExists(t, ctx, emailHash, testClient)
		assert.True(t, exists, "OTP should still exist after failed attempt")

		// Verify no user was created using raw SQL
		userExists := td.CheckUserExistsByEmailHashSQL(t, ctx, emailHash, testPool)
		assert.False(t, userExists, "No user should be created with invalid OTP")
	})

	t.Run("should return not found when OTP does not exist", func(t *testing.T) {
		// Clean state
		td.ClearAllTestData(t, ctx, testPool, testClient)

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
		sessionKeys, err := testClient.Keys(ctx, fmt.Sprintf("%s*", sx.SessionKeyPrefix)).Result()
		require.NoError(t, err)
		assert.Empty(t, sessionKeys, "No session keys should exist in Redis when OTP not found")

		tokenKeys, err := testClient.Keys(ctx, fmt.Sprintf("%s*", sx.AccessTokenKeyPrefix)).Result()
		require.NoError(t, err)
		assert.Empty(t, tokenKeys, "No token keys should exist in Redis when OTP not found")

		refreshTokenKeys, err := testClient.Keys(ctx, fmt.Sprintf("%s*", sx.RefreshTokenKeyPrefix)).Result()
		require.NoError(t, err)
		assert.Empty(t, refreshTokenKeys, "No refresh token keys should exist in Redis when OTP not found")

		// Verify no user was created using raw SQL
		emailHash := crypto.HashBasic(ctx, []byte(email))
		userExists := td.CheckUserExistsByEmailHashSQL(t, ctx, emailHash, testPool)
		assert.False(t, userExists, "No user should be created without valid OTP")
	})

	t.Run("should return gone when OTP is expired", func(t *testing.T) {
		// Clean state
		td.ClearAllTestData(t, ctx, testPool, testClient)

		email := "expired@example.com"

		// Create an OTP with past expiration time but longer Redis TTL
		otp := td.NewValidOTP(email)
		err := crypto.ProcessStruct(ctx, otp)
		require.NoError(t, err)
		otp.ExpiresAt = time.Now().Add(-1 * time.Minute)      // Logically expired
		td.InsertOTP(t, ctx, otp, testClient, 10*time.Minute) // But Redis keeps it

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
		sessionKeys, err := testClient.Keys(ctx, fmt.Sprintf("%s*", sx.SessionKeyPrefix)).Result()
		require.NoError(t, err)
		assert.Empty(t, sessionKeys, "No session keys should exist in Redis when OTP is expired")

		tokenKeys, err := testClient.Keys(ctx, fmt.Sprintf("%s*", sx.AccessTokenKeyPrefix)).Result()
		require.NoError(t, err)
		assert.Empty(t, tokenKeys, "No token keys should exist in Redis when OTP is expired")

		refreshTokenKeys, err := testClient.Keys(ctx, fmt.Sprintf("%s*", sx.RefreshTokenKeyPrefix)).Result()
		require.NoError(t, err)
		assert.Empty(t, refreshTokenKeys, "No refresh token keys should exist in Redis when OTP is expired")

		// Verify no user was created using raw SQL
		emailHash := crypto.HashBasic(ctx, []byte(email))
		userExists := td.CheckUserExistsByEmailHashSQL(t, ctx, emailHash, testPool)
		assert.False(t, userExists, "No user should be created with expired OTP")
	})

	t.Run("should return too many requests when max attempts exceeded", func(t *testing.T) {
		// Clean state
		td.ClearAllTestData(t, ctx, testPool, testClient)

		email := "maxattempts@example.com"

		// Create an OTP with max attempts already reached
		otp := td.NewValidOTP(email)
		err := crypto.ProcessStruct(ctx, otp)
		require.NoError(t, err)
		otp.Attempts = 5 // Assuming max attempts is 5
		td.InsertOTP(t, ctx, otp, testClient, 10*time.Minute)

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
		sessionKeys, err := testClient.Keys(ctx, fmt.Sprintf("%s*", sx.SessionKeyPrefix)).Result()
		require.NoError(t, err)
		assert.Empty(t, sessionKeys, "No session keys should exist in Redis when rate limited")

		tokenKeys, err := testClient.Keys(ctx, fmt.Sprintf("%s*", sx.AccessTokenKeyPrefix)).Result()
		require.NoError(t, err)
		assert.Empty(t, tokenKeys, "No token keys should exist in Redis when rate limited")

		refreshTokenKeys, err := testClient.Keys(ctx, fmt.Sprintf("%s*", sx.RefreshTokenKeyPrefix)).Result()
		require.NoError(t, err)
		assert.Empty(t, refreshTokenKeys, "No refresh token keys should exist in Redis when rate limited")

		// Verify no user was created using raw SQL
		emailHash := crypto.HashBasic(ctx, []byte(email))
		userExists := td.CheckUserExistsByEmailHashSQL(t, ctx, emailHash, testPool)
		assert.False(t, userExists, "No user should be created when rate limited")
	})

	t.Run("should return bad request for invalid email format", func(t *testing.T) {
		// Clean state
		td.ClearAllTestData(t, ctx, testPool, testClient)

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
		td.ClearAllTestData(t, ctx, testPool, testClient)

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
		td.ClearAllTestData(t, ctx, testPool, testClient)

		email := "existing@example.com"

		// Pre-create the user (simulating race condition)
		td.InsertTestUser(t, ctx, email, "Existing", "User", testPool, crypto)

		// TODO: I need a session with that user ID in it, that is how I get the session

		// Create valid OTP for the same email
		otp := td.NewValidOTP(email)
		crypto.ProcessStruct(ctx, otp)
		td.InsertOTP(t, ctx, otp, testClient, 10*time.Minute)

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
		accessTokenHash := crypto.HashBasic(ctx, []byte(accessToken))
		accessTokenKey := sx.AccessTokenKeyPrefix + accessTokenHash
		accessSessionIDStr, err := testClient.Get(ctx, accessTokenKey).Result()
		require.NoError(t, err, "Token mapping should exist in Redis even when user exists")

		// Verify refresh token exists in Redis even when user already exists
		refreshToken := refreshTokenCookie.Value
		refreshTokenHash := crypto.HashBasic(ctx, []byte(refreshToken))
		refreshTokenKey := sx.RefreshTokenKeyPrefix + refreshTokenHash
		refreshSessionIDStr, err := testClient.Get(ctx, refreshTokenKey).Result()
		require.NoError(t, err, "Refresh token mapping should exist in Redis even when user exists")
		assert.Equal(t, accessSessionIDStr, refreshSessionIDStr, "Refresh token should map to same session even when user exists")

		sessionKey := sx.SessionKeyPrefix + accessSessionIDStr
		sessionDataBytes, err := testClient.Get(ctx, sessionKey).Result()
		_ = sessionDataBytes
		require.NoError(t, err, "Session data should exist in Redis even when user exists")

		// Verify OTP was consumed
		emailHash := crypto.HashBasic(ctx, []byte(email))
		exists := td.CheckOTPExists(t, ctx, emailHash, testClient)
		assert.False(t, exists, "OTP should be consumed even when user exists")

		// Verify user still exists using raw SQL
		userExists := td.CheckUserExistsByEmailHashSQL(t, ctx, emailHash, testPool)
		assert.True(t, userExists, "User should still exist")
	})

	t.Run("should handle concurrent validation attempts", func(t *testing.T) {
		// Clean state
		td.ClearAllTestData(t, ctx, testPool, testClient)

		email := "concurrent@example.com"

		// Create OTP
		otp := td.NewValidOTP(email)
		crypto.ProcessStruct(ctx, otp)
		td.InsertOTP(t, ctx, otp, testClient, 10*time.Minute)

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
		emailHash := crypto.HashBasic(ctx, []byte(email))
		userExists := td.CheckUserExistsByEmailHashSQL(t, ctx, emailHash, testPool)
		assert.True(t, userExists, "User should be created by successful request")

		// Verify OTP was consumed
		exists := td.CheckOTPExists(t, ctx, emailHash, testClient)
		assert.False(t, exists, "OTP should be consumed")

		// Debug: Check how many sessions were created
		sessionKeys, err := testClient.Keys(ctx, fmt.Sprintf("%s*", sx.SessionKeyPrefix)).Result()
		require.NoError(t, err)
		t.Logf("Number of sessions created: %d (expected: 1)", len(sessionKeys))
		assert.Equal(t, 1, len(sessionKeys), "Only one session should be created")

		tokenKeys, err := testClient.Keys(ctx, fmt.Sprintf("%s*", sx.AccessTokenKeyPrefix)).Result()
		require.NoError(t, err)
		t.Logf("Number of token created: %d (expected: 1)", len(tokenKeys))
		assert.Equal(t, 1, len(tokenKeys), "Only one token should be created")
	})

	t.Run("should handle malformed JSON request", func(t *testing.T) {
		// Clean state
		td.ClearAllTestData(t, ctx, testPool, testClient)

		// Create malformed JSON request
		malformedJSON := `{"email": "test@example.com", "invalid_field": true}`
		req, err := http.NewRequestWithContext(
			ctx,
			http.MethodPost,
			testServerURL+"/auth/otp",
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
		td.ClearAllTestData(t, ctx, testPool, testClient)

		email := "completiontest@example.com"

		// Create and validate OTP to get a pending session
		otp := td.NewValidOTP(email)
		err := crypto.ProcessStruct(ctx, otp)
		require.NoError(t, err)
		td.InsertOTP(t, ctx, otp, testClient, 10*time.Minute)

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
		accessTokenHash := crypto.HashBasic(ctx, []byte(accessToken))
		accessTokenKey := sx.AccessTokenKeyPrefix + accessTokenHash
		accessSessionIDStr, err := testClient.Get(ctx, accessTokenKey).Result()
		require.NoError(t, err)

		// Verify refresh token exists in Redis
		refreshToken := refreshTokenCookie.Value
		refreshTokenHash := crypto.HashBasic(ctx, []byte(refreshToken))
		refreshTokenKey := sx.RefreshTokenKeyPrefix + refreshTokenHash
		refreshSessionIDStr, err := testClient.Get(ctx, refreshTokenKey).Result()
		require.NoError(t, err, "Refresh token mapping should exist in Redis")
		assert.Equal(t, accessSessionIDStr, refreshSessionIDStr, "Refresh token should map to same session")

		sessionKey := sx.SessionKeyPrefix + accessSessionIDStr
		sessionDataBytes, err := testClient.Get(ctx, sessionKey).Result()
		require.NoError(t, err)

		// Decode and verify session characteristics
		sessionData := []byte(sessionDataBytes)
		session := td.DecodeSessionWithDecryption(t, sessionData, crypto)

		// Verify it's a pending session with correct expiration
		assert.Equal(t, sx.SessionPending, session.State, "Session should be pending")
		assert.Nil(t, session.CompletedAt, "Session should not have completion timestamp yet")

		// Verify session expires in ~30 minutes (pending duration), not 24 hours
		expectedExpiry := session.CreatedAt.Add(sx.PendingSessionDuration)
		timeDiff := time.Duration(math.Abs(float64(session.ExpiresAt.Sub(expectedExpiry))))
		assert.LessOrEqual(t, timeDiff, 5*time.Second, "Pending session should have 30-minute expiration")

		// Verify session does NOT expire in 24 hours (active session duration)
		activeExpiry := session.CreatedAt.Add(sx.ActiveSessionDuration)
		activeTimeDiff := time.Duration(math.Abs(float64(session.ExpiresAt.Sub(activeExpiry))))
		assert.Greater(t, activeTimeDiff, 20*time.Hour, "Pending session should NOT have 24-hour expiration")

		t.Logf("Session created at: %s", session.CreatedAt.Format(time.RFC3339))
		t.Logf("Session expires at: %s", session.ExpiresAt.Format(time.RFC3339))
		t.Logf("Expected pending expiry: %s", expectedExpiry.Format(time.RFC3339))
		t.Logf("Time difference: %s", timeDiff.String())
		t.Logf("Pending session duration: %s", sx.PendingSessionDuration.String())
		t.Logf("Active session duration: %s", sx.ActiveSessionDuration.String())
	})
}
