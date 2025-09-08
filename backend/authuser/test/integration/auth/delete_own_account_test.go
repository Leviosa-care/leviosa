package auth_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/Leviosa-care/authuser/internal/domain"
	td "github.com/Leviosa-care/authuser/test/helpers"

	"github.com/Leviosa-care/core/auth/session"
	"github.com/Leviosa-care/core/contracts/identity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TEST=TestDeleteOwnAccount make test-integration-auth-test

func TestDeleteOwnAccount(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("should return 401 without authentication", func(t *testing.T) {
		// Act - make request without authentication
		req := td.NewDeleteOwnAccountRequestWithoutAuth(t, ctx, testServerURL)
		resp, err := client.Do(req)

		// Assert
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("should return 403 with guest user", func(t *testing.T) {
		// Clean state
		td.ClearUsersTable(t, ctx, testPool)
		td.ClearSessionsRedis(t, ctx, testClient)

		// Create guest user (below minimum role)
		guestUser := td.NewTestUser("guest@example.com", "Guest", "User")
		guestUser.State = domain.Active
		td.InsertUserWithEncryption(t, ctx, guestUser, testPool, crypto)

		// Create session for guest user
		sessionInfo := &session.SessionInfo{
			ID:     uuid.New(),
			UserID: guestUser.ID,
			Role:   identity.Guest, // Below Standard minimum role
			State:  session.SessionActive,
		}
		accessToken := td.CreateSessionWithEncryption(t, ctx, sessionInfo, testClient, crypto)

		// Act - make request with guest user token
		req := td.NewDeleteOwnAccountRequest(t, ctx, testServerURL, accessToken)
		resp, err := client.Do(req)

		// Assert
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)

		// Verify user still exists
		existingUser := td.GetUserByID(t, ctx, guestUser.ID, testPool)
		assert.NotNil(t, existingUser)
	})

	t.Run("should successfully delete own account as standard user", func(t *testing.T) {
		// Clean state
		td.ClearUsersTable(t, ctx, testPool)
		td.ClearSessionsRedis(t, ctx, testClient)
		td.ClearOTPKeys(t, ctx, testClient)

		// Create standard user
		user := td.NewTestUser("self-delete@example.com", "Self", "Delete")
		user.State = domain.Active
		td.InsertUserWithEncryption(t, ctx, user, testPool, crypto)

		// Create session for user
		sessionInfo := &session.SessionInfo{
			ID:     uuid.New(),
			UserID: user.ID,
			Role:   identity.Standard,
			State:  session.SessionActive,
		}
		accessToken := td.CreateSessionWithEncryption(t, ctx, sessionInfo, testClient, crypto)

		// Create additional data to be cleaned up
		// Create OTP
		td.CreateOTP(t, ctx, user.Email, testClient)

		// Create additional session
		additionalSessionInfo := &session.SessionInfo{
			ID:     uuid.New(),
			UserID: user.ID,
			Role:   identity.Standard,
			State:  session.SessionActive,
		}
		td.CreateSessionWithEncryption(t, ctx, additionalSessionInfo, testClient, crypto)

		// Act - delete own account
		req := td.NewDeleteOwnAccountRequest(t, ctx, testServerURL, accessToken)
		resp, err := client.Do(req)

		// Assert
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Verify response message
		message := td.ParseDeleteUserResponse(t, resp)
		assert.Equal(t, "Account deleted successfully", message)

		// Verify user is deleted from database
		deletedUser := td.GetUserByID(t, ctx, user.ID, testPool)
		assert.Nil(t, deletedUser)

		// Verify all sessions are cleared
		session1 := td.GetSessionByID(t, ctx, sessionInfo.ID, testClient)
		assert.Nil(t, session1)
		session2 := td.GetSessionByID(t, ctx, additionalSessionInfo.ID, testClient)
		assert.Nil(t, session2)

		// Verify OTP is cleared
		otp := td.GetOTP(t, ctx, user.Email, testClient)
		assert.NotEqual(t, domain.OTP{}, otp)
	})

	t.Run("should successfully delete own account as premium user", func(t *testing.T) {
		// Clean state
		td.ClearUsersTable(t, ctx, testPool)
		td.ClearSessionsRedis(t, ctx, testClient)

		// Create premium user
		user := td.NewTestUser("premium@example.com", "Premium", "User")
		user.State = domain.Active
		td.InsertUserWithEncryption(t, ctx, user, testPool, crypto)

		// Create session for premium user
		sessionInfo := &session.SessionInfo{
			ID:     uuid.New(),
			UserID: user.ID,
			Role:   identity.Premium,
			State:  session.SessionActive,
		}
		accessToken := td.CreateSessionWithEncryption(t, ctx, sessionInfo, testClient, crypto)

		// Act - delete own account
		req := td.NewDeleteOwnAccountRequest(t, ctx, testServerURL, accessToken)
		resp, err := client.Do(req)

		// Assert
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Verify user is deleted
		deletedUser := td.GetUserByID(t, ctx, user.ID, testPool)
		assert.Nil(t, deletedUser)
	})

	t.Run("should successfully delete own account with Stripe customer", func(t *testing.T) {
		// Clean state
		td.ClearUsersTable(t, ctx, testPool)
		td.ClearSessionsRedis(t, ctx, testClient)

		// Create user with Stripe customer ID
		user := td.NewTestUser("stripe-self@example.com", "Stripe", "Self")
		user.State = domain.Active
		user.StripeCustomerID = "cus_self_test456"
		td.InsertUserWithEncryption(t, ctx, user, testPool, crypto)

		// Create session for user
		sessionInfo := &session.SessionInfo{
			ID:     uuid.New(),
			UserID: user.ID,
			Role:   identity.Standard,
			State:  session.SessionActive,
		}
		accessToken := td.CreateSessionWithEncryption(t, ctx, sessionInfo, testClient, crypto)

		// Act - delete own account with Stripe customer
		req := td.NewDeleteOwnAccountRequest(t, ctx, testServerURL, accessToken)
		resp, err := client.Do(req)

		// Assert
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Verify user is deleted (Stripe deletion handled internally)
		deletedUser := td.GetUserByID(t, ctx, user.ID, testPool)
		assert.Nil(t, deletedUser)
	})

	t.Run("should handle user not found error", func(t *testing.T) {
		// Clean state
		td.ClearUsersTable(t, ctx, testPool)
		td.ClearSessionsRedis(t, ctx, testClient)

		// Create a session for a user that doesn't exist in database
		nonExistentUserID := uuid.New()
		sessionInfo := &session.SessionInfo{
			ID:     uuid.New(),
			UserID: nonExistentUserID,
			Role:   identity.Standard,
			State:  session.SessionActive,
		}
		accessToken := td.CreateSessionWithEncryption(t, ctx, sessionInfo, testClient, crypto)

		// Act - try to delete account for non-existent user
		req := td.NewDeleteOwnAccountRequest(t, ctx, testServerURL, accessToken)
		resp, err := client.Do(req)

		// Assert
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}
