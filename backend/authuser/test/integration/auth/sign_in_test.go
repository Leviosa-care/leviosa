package auth_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/Leviosa-care/authuser/internal/domain"
	td "github.com/Leviosa-care/authuser/test/helpers"

	ck "github.com/Leviosa-care/core/auth/cookies"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TEST=TestSignIn make test-integration-auth-test

func TestSignIn(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	validEmail := "signin-test@example.com"
	_ = validEmail
	validPassword := "K9$!qf>2]Ez~:Kb6C(D3RqP8"
	_ = validPassword

	t.Run("should successfully sign in with valid credentials", func(t *testing.T) {
		// Clear all test data
		td.ClearAllTestData(t, ctx, testPool, testClient)

		// Create an active user with the test password
		user, err := td.NewTestUserWithEncryption(validEmail, "John", "Doe", crypto)
		require.NoError(t, err)
		user.State = domain.Active
		user.Password = validPassword

		// Process encryption/hashing again to hash the password
		err = crypto.ProcessStruct(ctx, user)
		require.NoError(t, err)

		// Insert the user into database
		td.InsertUser(t, ctx, user, testPool)

		// Create sign-in request
		request := domain.SignInRequest{
			Email:    validEmail,
			Password: validPassword,
		}

		// Make HTTP request
		req := td.NewSignInRequest(t, ctx, testServerURL, request)
		resp, err := client.Do(req)

		// Assert HTTP response
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		// Parse response JSON
		var response struct {
			Message string `json:"message"`
			Status  string `json:"status"`
		}
		td.ParseJSONResponse(t, resp, &response)

		// Validate response structure
		assert.Equal(t, "user logged in successfully", response.Message)
		assert.Equal(t, "created", response.Status)

		// Verify access token cookie is set
		cookies := resp.Cookies()
		var accessTokenCookie *http.Cookie
		for _, cookie := range cookies {
			if cookie.Name == ck.AccessTokenCookieName {
				accessTokenCookie = cookie
				break
			}
		}
		require.NotNil(t, accessTokenCookie, "Access token cookie should be set")

		// Verify cookie attributes for access token
		assert.True(t, accessTokenCookie.HttpOnly, "Access token cookie should be HttpOnly")
		assert.True(t, accessTokenCookie.Secure, "Access token cookie should be Secure")
		assert.Equal(t, http.SameSiteStrictMode, accessTokenCookie.SameSite, "Access token cookie should have SameSite=Strict")
		assert.Equal(t, "/", accessTokenCookie.Path, "Access token cookie path should be /")
		assert.NotEmpty(t, accessTokenCookie.Value, "Access token cookie should have a value")

		// Verify refresh token cookie is set
		var refreshTokenCookie *http.Cookie
		for _, cookie := range cookies {
			if cookie.Name == ck.RefreshTokenCookieName {
				refreshTokenCookie = cookie
				break
			}
		}
		require.NotNil(t, refreshTokenCookie, "Refresh token cookie should be set")

		// Verify cookie attributes for refresh token
		assert.True(t, refreshTokenCookie.HttpOnly, "Refresh token cookie should be HttpOnly")
		assert.True(t, refreshTokenCookie.Secure, "Refresh token cookie should be Secure")
		assert.Equal(t, http.SameSiteStrictMode, refreshTokenCookie.SameSite, "Refresh token cookie should have SameSite=Strict")
		assert.Equal(t, ck.RefreshEndpoint, refreshTokenCookie.Path, "Refresh token cookie path should be /")
		assert.NotEmpty(t, refreshTokenCookie.Value, "Refresh token cookie should have a value")
	})

	t.Run("should fail with invalid email", func(t *testing.T) {
		// Clear all test data
		td.ClearAllTestData(t, ctx, testPool, testClient)

		// Create sign-in request with non-existent email
		request := domain.SignInRequest{
			Email:    "nonexistent@example.com",
			Password: validPassword,
		}

		// Make HTTP request
		req := td.NewSignInRequest(t, ctx, testServerURL, request)
		resp, err := client.Do(req)

		// Assert HTTP response
		require.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("should fail with invalid password", func(t *testing.T) {
		// Clear all test data
		td.ClearAllTestData(t, ctx, testPool, testClient)

		// Create an active user with the test password
		user, err := td.NewTestUserWithEncryption(validEmail, "John", "Doe", crypto)
		require.NoError(t, err)
		user.State = domain.Active
		user.Password = validPassword

		// Process encryption/hashing again to hash the password
		err = crypto.ProcessStruct(ctx, user)
		require.NoError(t, err)

		// Insert the user into database
		td.InsertUser(t, ctx, user, testPool)

		// Create sign-in request with wrong password
		request := domain.SignInRequest{
			Email:    validEmail,
			Password: "WrongPassword123!",
		}

		// Make HTTP request
		req := td.NewSignInRequest(t, ctx, testServerURL, request)
		resp, err := client.Do(req)

		// Assert HTTP response
		require.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("should fail with inactive user account", func(t *testing.T) {
		// Clear all test data
		td.ClearAllTestData(t, ctx, testPool, testClient)

		// Create a pending user (not yet activated)
		user, err := td.NewTestUserWithEncryption(validEmail, "John", "Doe", crypto)
		require.NoError(t, err)
		user.State = domain.Pending
		user.Password = validPassword

		// Process encryption/hashing
		err = crypto.ProcessStruct(ctx, user)
		require.NoError(t, err)

		// Insert the user into database
		td.InsertUser(t, ctx, user, testPool)

		// Create sign-in request
		request := domain.SignInRequest{
			Email:    validEmail,
			Password: validPassword,
		}

		// Make HTTP request
		req := td.NewSignInRequest(t, ctx, testServerURL, request)
		resp, err := client.Do(req)

		// Assert HTTP response
		require.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("should fail with unverified user account", func(t *testing.T) {
		// Clear all test data
		td.ClearAllTestData(t, ctx, testPool, testClient)

		// Create an unverified user
		user, err := td.NewTestUserWithEncryption(validEmail, "John", "Doe", crypto)
		require.NoError(t, err)
		user.State = domain.Unverified
		user.Password = validPassword

		// Process encryption/hashing
		err = crypto.ProcessStruct(ctx, user)
		require.NoError(t, err)

		// Insert the user into database
		td.InsertUser(t, ctx, user, testPool)

		// Create sign-in request
		request := domain.SignInRequest{
			Email:    validEmail,
			Password: validPassword,
		}

		// Make HTTP request
		req := td.NewSignInRequest(t, ctx, testServerURL, request)
		resp, err := client.Do(req)

		// Assert HTTP response
		require.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("should fail with invalid request body", func(t *testing.T) {
		// Make HTTP request with invalid JSON
		req := td.NewInvalidJSONRequest(t, ctx, testServerURL, "POST", "/auth/login")
		resp, err := client.Do(req)

		// Assert HTTP response
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should fail with missing email field", func(t *testing.T) {
		// Create sign-in request missing email
		request := domain.SignInRequest{
			Password: validPassword,
		}

		// Make HTTP request
		req := td.NewSignInRequest(t, ctx, testServerURL, request)
		resp, err := client.Do(req)

		// Assert HTTP response
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should fail with missing password field", func(t *testing.T) {
		// Create sign-in request missing password
		request := domain.SignInRequest{
			Email: validEmail,
		}

		// Make HTTP request
		req := td.NewSignInRequest(t, ctx, testServerURL, request)
		resp, err := client.Do(req)

		// Assert HTTP response
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}
