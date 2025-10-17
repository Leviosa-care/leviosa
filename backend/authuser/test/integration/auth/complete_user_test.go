package auth_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/Leviosa-care/authuser/internal/domain"
	td "github.com/Leviosa-care/authuser/test/helpers"

	authEndpoints "github.com/Leviosa-care/authuser/internal/adapters/http/auth"

	ck "github.com/Leviosa-care/core/auth/cookies"
	"github.com/Leviosa-care/core/auth/session"
	"github.com/Leviosa-care/core/errs"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newPendingUser(email string) *domain.User {
	return &domain.User{
		ID:    uuid.New(),
		Email: email,
		State: domain.Unverified,
	}
}

// TEST=TestCompleteUser make test-integration-auth-test

func TestCompleteUser(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	validEmail := "test@example.com"

	// Create valid complete user request
	validRequest := domain.CompleteUserRequest{
		Password:  "qPDAR0.4Z8{vpCO]",
		FirstName: "John",
		LastName:  "Doe",
		BirthDate: time.Now().AddDate(-25, 0, 0),
		Gender: domain.GenderInput{
			Gender: domain.GenderMan,
		},
		Telephone:  "0612345678",
		PostalCode: "75001",
		City:       "Paris",
		Address1:   "123 Rue de Rivoli",
		Address2:   "Apt 4B",
	}

	t.Run("should successfully complete user registration with pending session", func(t *testing.T) {
		// Clean state
		td.ClearAllTestData(t, ctx, testPool, redisClient)

		pendingUser := newPendingUser(validEmail)
		pendingUserEncx, err := domain.ProcessUserEncx(ctx, crypto, pendingUser)
		require.NoError(t, err)
		err = td.InsertUserEncx(t, ctx, pendingUserEncx, testPool, crypto)
		require.NoError(t, err)

		// Create a pending session for this specific user
		// pendingSession := td.CreateTestPendingSessionWithUserIDAndCrypto(t, pendingUser.ID, crypto)
		pendingSession, err := td.NewTestSession(t, crypto)
		pendingSession.UserID = pendingUser.ID
		pendingSession.State = session.SessionPending

		pendingSessionEncx, err := session.ProcessSessionEncx(ctx, crypto, pendingSession)
		require.NoError(t, err)

		td.InsertSessionEncx(t, ctx, redisClient, pendingSessionEncx, time.Hour)

		// Make HTTP request with session access token
		req := td.NewCompleteUserRequest(t, ctx, testServerURL, validRequest, pendingSession.AccessToken)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert HTTP response
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		message := td.ParseCompleteUserResponse(t, resp)
		assert.Equal(t, "User registration completed successfully", message)

		// Verify user was created in database
		userEncx, err := td.GetUserEnxByID(t, ctx, pendingUser.ID, testPool, crypto)
		require.NoError(t, err)

		user, err := domain.DecryptUserEncx(ctx, crypto, userEncx)
		require.NoError(t, err)

		verifyCompletedUserFields(t, pendingSession.UserID, &validRequest, user)

		// Verify session was removed after completion
		sessionValueExists, err := redisClient.Exists(ctx, session.SessionKeyPrefix+"*").Result()
		assert.Equal(t, int64(0), sessionValueExists, "All sessions should be removed after completion")
		userSessionValueExists, err := redisClient.Exists(ctx, session.UserSessionIndexKeyPrefix+"*").Result()
		assert.Equal(t, int64(0), userSessionValueExists, "All user sessions should be removed after completion")
		accessTokenValueExists, err := redisClient.Exists(ctx, session.AccessTokenKeyPrefix+"*").Result()
		assert.Equal(t, int64(0), accessTokenValueExists, "All access token sessions should be removed after completion")
		refreshTokenValueExists, err := redisClient.Exists(ctx, session.RefreshTokenKeyPrefix+"*").Result()
		assert.Equal(t, int64(0), refreshTokenValueExists, "All refresh token sessions should be removed after completion")
	})

	t.Run("should return 401 when session cookie is missing", func(t *testing.T) {
		// Clean state
		td.ClearAllTestData(t, ctx, testPool, redisClient)

		// Make HTTP request without session token
		req := td.NewCompleteUserRequest(t, ctx, testServerURL, validRequest, "")
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert HTTP response
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		errorMsg, statusCode := td.ParseErrorResponse(t, resp)
		assert.Equal(t, http.StatusUnauthorized, statusCode)
		assert.Contains(t, errorMsg, errs.ErrUnauthorized.Error())
	})

	t.Run("should return 401 when session cookie is empty", func(t *testing.T) {
		// Clean state
		td.ClearAllTestData(t, ctx, testPool, redisClient)

		// Make HTTP request with empty session token
		req := td.NewCompleteUserRequest(t, ctx, testServerURL, validRequest, "")
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert HTTP response
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		errorMsg, statusCode := td.ParseErrorResponse(t, resp)
		assert.Equal(t, http.StatusUnauthorized, statusCode)
		assert.Contains(t, errorMsg, errs.ErrUnauthorized.Error())
	})

	t.Run("should return 401 when session does not exist", func(t *testing.T) {
		// Clean state
		td.ClearAllTestData(t, ctx, testPool, redisClient)

		// Make HTTP request with non-existent session randomToken
		randomToken, err := session.GenerateToken()
		require.NoError(t, err)
		req := td.NewCompleteUserRequest(t, ctx, testServerURL, validRequest, randomToken)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert HTTP response
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("should return 401 when session is expired", func(t *testing.T) {
		// Clean state
		td.ClearAllTestData(t, ctx, testPool, redisClient)

		// Create an expired session
		// expiredSession := td.CreateTestPendingSessionWithCrypto(t, crypto)
		expiredSession, err := td.NewTestSession(t, crypto)
		require.NoError(t, err)
		expiredSession.State = session.SessionPending
		expiredSession.ExpiresAt = time.Now().Add(-1 * time.Hour) // Expired 1 hour ago

		// Re-encrypt with updated expiry
		sessionEncx, err := session.ProcessSessionEncx(ctx, crypto, expiredSession)
		require.NoError(t, err)

		td.InsertSessionEncx(t, ctx, redisClient, sessionEncx, time.Hour)

		// Make HTTP request with expired session token
		req := td.NewCompleteUserRequest(t, ctx, testServerURL, validRequest, expiredSession.AccessToken)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert HTTP response
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("should return 409 when session is already active", func(t *testing.T) {
		// Clean state
		td.ClearAllTestData(t, ctx, testPool, redisClient)

		// Create an active session
		// activeSession := td.CreateTestSessionWithCrypto(t, crypto)
		activeSession, err := td.NewTestSession(t, crypto)
		activeSession.State = session.SessionActive
		require.NoError(t, err)

		sessionEncx, err := session.ProcessSessionEncx(ctx, crypto, activeSession)
		require.NoError(t, err)
		td.InsertSessionEncx(t, ctx, redisClient, sessionEncx, time.Hour)

		// Make HTTP request with active session token
		req := td.NewCompleteUserRequest(t, ctx, testServerURL, validRequest, activeSession.AccessToken)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert HTTP response
		assert.Equal(t, http.StatusConflict, resp.StatusCode)
	})

	t.Run("should return 400 for invalid JSON request body", func(t *testing.T) {
		// Clean state
		td.ClearAllTestData(t, ctx, testPool, redisClient)

		// Create a pending session
		// pendingSession := td.CreateTestPendingSessionWithCrypto(t, crypto)
		pendingSession, err := td.NewTestSession(t, crypto)
		require.NoError(t, err)
		pendingSession.State = session.SessionPending

		pendingSessionEncx, err := session.ProcessSessionEncx(ctx, crypto, pendingSession)
		require.NoError(t, err)

		td.InsertSessionEncx(t, ctx, redisClient, pendingSessionEncx, time.Hour)

		// Make HTTP request with invalid JSON (manually crafted)
		req, err := http.NewRequestWithContext(
			ctx,
			http.MethodPost,
			testServerURL+authEndpoints.CompleteUserEndpoint,
			nil, // No body - will cause JSON decode error
		)
		require.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")
		cookie := &http.Cookie{
			Name:  ck.AccessTokenCookieName,
			Value: pendingSession.AccessToken,
		}
		req.AddCookie(cookie)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert HTTP response
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		errorMsg, statusCode := td.ParseErrorResponse(t, resp)
		assert.Equal(t, http.StatusBadRequest, statusCode)
		assert.Contains(t, errorMsg, "invalid request body")
	})

	t.Run("should return 400 for invalid password", func(t *testing.T) {
		// Clean state
		td.ClearAllTestData(t, ctx, testPool, redisClient)

		// Create a pending session
		// pendingSession := td.CreateTestPendingSessionWithCrypto(t, crypto)
		pendingSession, err := td.NewTestSession(t, crypto)
		require.NoError(t, err)
		pendingSession.State = session.SessionPending

		pendingSessionEncx, err := session.ProcessSessionEncx(ctx, crypto, pendingSession)
		require.NoError(t, err)

		td.InsertSessionEncx(t, ctx, redisClient, pendingSessionEncx, time.Hour)

		request := validRequest
		request.Password = "weak"

		// Make HTTP request
		req := td.NewCompleteUserRequest(t, ctx, testServerURL, request, pendingSession.AccessToken)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert HTTP response
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should return 400 for invalid birth date (future)", func(t *testing.T) {
		// Clean state
		td.ClearAllTestData(t, ctx, testPool, redisClient)

		// Create a pending session
		// pendingSession := td.CreateTestPendingSessionWithCrypto(t, crypto)
		pendingSession, err := td.NewTestSession(t, crypto)
		require.NoError(t, err)
		pendingSession.State = session.SessionPending

		pendingSessionEncx, err := session.ProcessSessionEncx(ctx, crypto, pendingSession)
		require.NoError(t, err)

		td.InsertSessionEncx(t, ctx, redisClient, pendingSessionEncx, time.Hour)

		request := validRequest
		request.BirthDate = time.Now().AddDate(1, 0, 0) // Future date

		// Make HTTP request
		req := td.NewCompleteUserRequest(t, ctx, testServerURL, request, pendingSession.AccessToken)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert HTTP response
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should return 400 for invalid birth date (too young)", func(t *testing.T) {
		// Clean state
		td.ClearAllTestData(t, ctx, testPool, redisClient)

		// Create a pending session
		// pendingSession := td.CreateTestPendingSessionWithCrypto(t, crypto)
		pendingSession, err := td.NewTestSession(t, crypto)
		require.NoError(t, err)
		pendingSession.State = session.SessionPending

		pendingSessionEncx, err := session.ProcessSessionEncx(ctx, crypto, pendingSession)
		require.NoError(t, err)

		td.InsertSessionEncx(t, ctx, redisClient, pendingSessionEncx, time.Hour)

		request := validRequest
		request.BirthDate = time.Now().AddDate(-10, 0, 0) // 10 years old (under 13)

		// Make HTTP request
		req := td.NewCompleteUserRequest(t, ctx, testServerURL, request, pendingSession.AccessToken)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert HTTP response
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should return 400 for invalid phone number", func(t *testing.T) {
		// Clean state
		td.ClearAllTestData(t, ctx, testPool, redisClient)

		// Create a pending session
		// pendingSession := td.CreateTestPendingSessionWithCrypto(t, crypto)
		pendingSession, err := td.NewTestSession(t, crypto)
		require.NoError(t, err)
		pendingSession.State = session.SessionPending

		pendingSessionEncx, err := session.ProcessSessionEncx(ctx, crypto, pendingSession)
		require.NoError(t, err)

		td.InsertSessionEncx(t, ctx, redisClient, pendingSessionEncx, time.Hour)

		request := validRequest
		request.Telephone = "invalid-phone" // Invalid phone format

		// Make HTTP request
		req := td.NewCompleteUserRequest(t, ctx, testServerURL, request, pendingSession.AccessToken)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert HTTP response
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should return 400 for invalid gender", func(t *testing.T) {
		// Clean state
		td.ClearAllTestData(t, ctx, testPool, redisClient)

		// Create a pending session
		// pendingSession := td.CreateTestPendingSessionWithCrypto(t, crypto)
		pendingSession, err := td.NewTestSession(t, crypto)
		pendingSession.State = session.SessionPending
		require.NoError(t, err)

		pendingSessionEncx, err := session.ProcessSessionEncx(ctx, crypto, pendingSession)
		require.NoError(t, err)

		td.InsertSessionEncx(t, ctx, redisClient, pendingSessionEncx, time.Hour)

		request := validRequest
		request.Gender = domain.GenderInput{
			Gender: "invalid_gender", // Invalid gender value
		}

		// Make HTTP request
		req := td.NewCompleteUserRequest(t, ctx, testServerURL, request, pendingSession.AccessToken)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert HTTP response
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should successfully complete user registration with custom gender", func(t *testing.T) {
		// Clean state
		td.ClearAllTestData(t, ctx, testPool, redisClient)

		pendingUser := newPendingUser(validEmail)
		pendingUserEncx, err := domain.ProcessUserEncx(ctx, crypto, pendingUser)
		require.NoError(t, err)
		err = td.InsertUserEncx(t, ctx, pendingUserEncx, testPool, crypto)
		require.NoError(t, err)

		// Create a pending session
		// pendingSession := td.CreateTestPendingSessionWithUserIDAndCrypto(t, pendingUser.ID, crypto)
		pendingSession, err := td.NewTestSession(t, crypto)
		pendingSession.State = session.SessionPending
		pendingSession.UserID = pendingUser.ID
		require.NoError(t, err)

		pendingSessionEncx, err := session.ProcessSessionEncx(ctx, crypto, pendingSession)
		require.NoError(t, err)

		td.InsertSessionEncx(t, ctx, redisClient, pendingSessionEncx, time.Hour)

		request := validRequest
		request.Gender = domain.GenderInput{
			Gender:       domain.GenderCustom,
			CustomGender: "Genderfluid",
		}

		// Make HTTP request
		req := td.NewCompleteUserRequest(t, ctx, testServerURL, request, pendingSession.AccessToken)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert HTTP response
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		message := td.ParseCompleteUserResponse(t, resp)
		assert.Equal(t, "User registration completed successfully", message)

		// Verify user was created in database
		userEncx, err := td.GetUserEnxByID(t, ctx, pendingUser.ID, testPool, crypto)
		require.NoError(t, err)
		user, err := domain.DecryptUserEncx(ctx, crypto, userEncx)
		require.NoError(t, err)
		verifyCompletedUserFields(t, pendingSession.UserID, &request, user)

		// Verify session was removed after completion
		sessionValueExists, err := redisClient.Exists(ctx, session.SessionKeyPrefix+"*").Result()
		assert.Equal(t, int64(0), sessionValueExists, "All sessions should be removed after completion")
		userSessionValueExists, err := redisClient.Exists(ctx, session.UserSessionIndexKeyPrefix+"*").Result()
		assert.Equal(t, int64(0), userSessionValueExists, "All user sessions should be removed after completion")
		accessTokenValueExists, err := redisClient.Exists(ctx, session.AccessTokenKeyPrefix+"*").Result()
		assert.Equal(t, int64(0), accessTokenValueExists, "All access token sessions should be removed after completion")
		refreshTokenValueExists, err := redisClient.Exists(ctx, session.RefreshTokenKeyPrefix+"*").Result()
		assert.Equal(t, int64(0), refreshTokenValueExists, "All refresh token sessions should be removed after completion")
	})

	t.Run("should return 400 for custom gender without customGender field", func(t *testing.T) {
		// Clean state
		td.ClearAllTestData(t, ctx, testPool, redisClient)

		// Create a pending session
		// pendingSession := td.CreateTestPendingSessionWithCrypto(t, crypto)
		pendingSession, err := td.NewTestSession(t, crypto)
		pendingSession.State = session.SessionPending
		require.NoError(t, err)

		pendingSessionEncx, err := session.ProcessSessionEncx(ctx, crypto, pendingSession)
		require.NoError(t, err)

		td.InsertSessionEncx(t, ctx, redisClient, pendingSessionEncx, time.Hour)

		request := validRequest
		request.Gender = domain.GenderInput{
			Gender: domain.GenderCustom,
		}

		// Make HTTP request
		req := td.NewCompleteUserRequest(t, ctx, testServerURL, request, pendingSession.AccessToken)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert HTTP response
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should return 400 for missing required fields", func(t *testing.T) {
		// Clean state
		td.ClearAllTestData(t, ctx, testPool, redisClient)

		// Create a pending session
		// pendingSession := td.CreateTestPendingSessionWithCrypto(t, crypto)
		pendingSession, err := td.NewTestSession(t, crypto)
		pendingSession.State = session.SessionPending
		require.NoError(t, err)

		pendingSessionEncx, err := session.ProcessSessionEncx(ctx, crypto, pendingSession)
		require.NoError(t, err)

		td.InsertSessionEncx(t, ctx, redisClient, pendingSessionEncx, time.Hour)

		// Create complete user request with missing required fields
		request := domain.CompleteUserRequest{
			Password: validRequest.Password,
			// Missing all other required fields
		}

		// Make HTTP request
		req := td.NewCompleteUserRequest(t, ctx, testServerURL, request, pendingSession.AccessToken)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert HTTP response
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should successfully complete user registration with minimal valid data", func(t *testing.T) {
		// Clean state
		td.ClearAllTestData(t, ctx, testPool, redisClient)

		// Create an unverified user first
		pendingUser := newPendingUser("test@example.com")
		pendingUserEncx, err := domain.ProcessUserEncx(ctx, crypto, pendingUser)
		require.NoError(t, err)
		err = td.InsertUserEncx(t, ctx, pendingUserEncx, testPool, crypto)
		require.NoError(t, err)

		// Create a pending session
		// pendingSession := td.CreateTestPendingSessionWithUserIDAndCrypto(t, pendingUser.ID, crypto)
		pendingSession, err := td.NewTestSession(t, crypto)
		pendingSession.State = session.SessionPending
		pendingSession.UserID = pendingUser.ID
		require.NoError(t, err)

		pendingSessionEncx, err := session.ProcessSessionEncx(ctx, crypto, pendingSession)
		require.NoError(t, err)

		td.InsertSessionEncx(t, ctx, redisClient, pendingSessionEncx, time.Hour)

		// Create complete user request with minimal valid data (no Address2)
		request := domain.CompleteUserRequest{
			Password:  "qPDAR0.4Z8{vpCO]",
			FirstName: "Jo", // Minimum 2 characters
			LastName:  "Do", // Minimum 2 characters
			BirthDate: time.Now().AddDate(-25, 0, 0),
			Gender: domain.GenderInput{
				Gender: domain.GenderPreferNotToSay,
			},
			Telephone:  "0612345678",
			PostalCode: "75001", // Minimum 3 characters
			City:       "Paris", // Minimum 2 characters
			Address1:   "1 St",  // Minimum 5 characters but this should fail
		}

		// Make HTTP request
		req := td.NewCompleteUserRequest(t, ctx, testServerURL, request, pendingSession.AccessToken)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Address1 with 4 characters should fail (minimum is 5)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		// Fix the address and try again
		request.Address1 = "123 Rue de Rivoli" // 6 characters, should pass
		req = td.NewCompleteUserRequest(t, ctx, testServerURL, request, pendingSession.AccessToken)
		resp, err = client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert HTTP response
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		message := td.ParseCompleteUserResponse(t, resp)
		assert.Equal(t, "User registration completed successfully", message)

		userEncx, err := td.GetUserEnxByID(t, ctx, pendingUser.ID, testPool, crypto)
		require.NoError(t, err)
		user, err := domain.DecryptUserEncx(ctx, crypto, userEncx)
		require.NoError(t, err)
		verifyCompletedUserFields(t, pendingSession.UserID, &request, user)

		// Verify session was removed after completion
		sessionValueExists, err := redisClient.Exists(ctx, session.SessionKeyPrefix+"*").Result()
		assert.Equal(t, int64(0), sessionValueExists, "All sessions should be removed after completion")
		userSessionValueExists, err := redisClient.Exists(ctx, session.UserSessionIndexKeyPrefix+"*").Result()
		assert.Equal(t, int64(0), userSessionValueExists, "All user sessions should be removed after completion")
		accessTokenValueExists, err := redisClient.Exists(ctx, session.AccessTokenKeyPrefix+"*").Result()
		assert.Equal(t, int64(0), accessTokenValueExists, "All access token sessions should be removed after completion")
		refreshTokenValueExists, err := redisClient.Exists(ctx, session.RefreshTokenKeyPrefix+"*").Result()
		assert.Equal(t, int64(0), refreshTokenValueExists, "All refresh token sessions should be removed after completion")
	})
}

func verifyCompletedUserFields(t *testing.T, userID uuid.UUID, expectedUser *domain.CompleteUserRequest, actualUser *domain.User) {
	t.Helper()
	assert.Equal(t, userID, actualUser.ID)
	assert.Equal(t, expectedUser.FirstName, actualUser.FirstName)
	assert.Equal(t, expectedUser.LastName, actualUser.LastName)
	assert.NotEmpty(t, actualUser.BirthDate)
	if expectedUser.Gender.Gender == domain.GenderCustom {
		assert.Equal(t, expectedUser.Gender.CustomGender, actualUser.Gender)
	} else {
		assert.Equal(t, expectedUser.Gender.Gender, domain.Gender(actualUser.Gender))
	}
	assert.Equal(t, expectedUser.Telephone, actualUser.Telephone)
	assert.Equal(t, expectedUser.PostalCode, actualUser.PostalCode)
	assert.Equal(t, expectedUser.City, actualUser.City)
	assert.Equal(t, expectedUser.Address1, actualUser.Address1)
	assert.Equal(t, expectedUser.Address2, actualUser.Address2)
	assert.NotEmpty(t, actualUser.StripeCustomerID)
}
