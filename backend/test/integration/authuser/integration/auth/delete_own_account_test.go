package auth_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	td "github.com/Leviosa-care/leviosa/backend/test/helpers"

	"github.com/Leviosa-care/leviosa/backend/internal/common/auth/session"
	"github.com/Leviosa-care/leviosa/backend/internal/common/contracts/identity"
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
		td.ClearSessionsRedis(t, ctx, redisClient)

		// Create guest user (below minimum role)
		guestUser := td.NewTestUser(t, "guest@example.com", "Guest", "User")
		guestUser.State = domain.Active
		guestUser.Role = identity.GuestStr
		guestUserEncx, err := domain.ProcessUserEncx(ctx, crypto, guestUser)
		require.NoError(t, err)
		err = td.InsertUserEncx(t, ctx, guestUserEncx, testPool, crypto)
		require.NoError(t, err)

		// Create session for guest user
		sessionInfo := &session.SessionInfo{
			ID:     uuid.New(),
			UserID: guestUser.ID,
			Role:   identity.Guest, // Below Standard minimum role
			State:  session.SessionActive,
		}
		accessToken := td.CreateSessionWithEncryption(t, ctx, sessionInfo, redisClient, crypto)

		// Act - make request with guest user token
		req := td.NewDeleteOwnAccountRequest(t, ctx, testServerURL, accessToken)
		resp, err := client.Do(req)

		// Assert
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)

		// Verify user still exists
		existingUserEncx, err := td.GetUserEnxByID(t, ctx, guestUser.ID, testPool, crypto)
		require.NoError(t, err)
		existingUser, err := domain.DecryptUserEncx(ctx, crypto, existingUserEncx)
		require.NoError(t, err)
		assert.NotNil(t, existingUser)
	})

	t.Run("should successfully delete own account as standard user", func(t *testing.T) {
		// Clean state
		td.ClearUsersTable(t, ctx, testPool)
		td.ClearSessionsRedis(t, ctx, redisClient)
		td.ClearOTPKeys(t, ctx, redisClient)

		// Create standard user
		user := td.NewTestUser(t, "self-delete@example.com", "Self", "Delete")
		user.State = domain.Active
		userEncx, err := domain.ProcessUserEncx(ctx, crypto, user)
		require.NoError(t, err)
		err = td.InsertUserEncx(t, ctx, userEncx, testPool, crypto)
		require.NoError(t, err)

		// Create session for user
		sessionInfo := &session.SessionInfo{
			ID:     uuid.New(),
			UserID: user.ID,
			Role:   identity.Standard,
			State:  session.SessionActive,
		}
		accessToken := td.CreateSessionWithEncryption(t, ctx, sessionInfo, redisClient, crypto)

		// Create additional data to be cleaned up
		// Create OTP
		td.CreateOTP(t, ctx, user.Email, redisClient, crypto)

		// Create additional session
		additionalSessionInfo := &session.SessionInfo{
			ID:     uuid.New(),
			UserID: user.ID,
			Role:   identity.Standard,
			State:  session.SessionActive,
		}
		td.CreateSessionWithEncryption(t, ctx, additionalSessionInfo, redisClient, crypto)

		// Act - delete own account
		req := td.NewDeleteOwnAccountRequest(t, ctx, testServerURL, accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Verify response message
		message := td.ParseDeleteUserResponse(t, resp)
		assert.Equal(t, "Account deleted successfully", message)

		// Verify user is deleted from database
		deletedUserEncx, err := td.GetUserEnxByID(t, ctx, user.ID, testPool, crypto)
		assert.Error(t, err)
		assert.Equal(t, *deletedUserEncx, domain.UserEncx{})

		// Verify all sessions are cleared
		session1 := td.GetSessionByID(t, ctx, sessionInfo.ID, redisClient)
		assert.Nil(t, session1)
		session2 := td.GetSessionByID(t, ctx, additionalSessionInfo.ID, redisClient)
		assert.Nil(t, session2)

		// Verify OTP is cleared
		_, err = td.GetOTPByEmail(t, ctx, user.Email, redisClient, crypto)
		assert.Error(t, err)
	})

	t.Run("should successfully delete own account as premium user", func(t *testing.T) {
		// Clean state
		td.ClearUsersTable(t, ctx, testPool)
		td.ClearSessionsRedis(t, ctx, redisClient)

		// Create premium user
		user := td.NewTestUser(t, "premium@example.com", "Premium", "User")
		user.State = domain.Active
		user.Role = identity.PremiumStr

		userEncx, err := domain.ProcessUserEncx(ctx, crypto, user)
		require.NoError(t, err)

		err = td.InsertUserEncx(t, ctx, userEncx, testPool, crypto)
		require.NoError(t, err)

		// Create session for premium user
		sessionInfo := &session.SessionInfo{
			ID:     uuid.New(),
			UserID: user.ID,
			Role:   identity.Premium,
			State:  session.SessionActive,
		}
		accessToken := td.CreateSessionWithEncryption(t, ctx, sessionInfo, redisClient, crypto)

		// Act - delete own account
		req := td.NewDeleteOwnAccountRequest(t, ctx, testServerURL, accessToken)
		resp, err := client.Do(req)

		// Assert
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Verify user is deleted
		deletedUserEncx, err := td.GetUserEnxByID(t, ctx, user.ID, testPool, crypto)
		assert.Error(t, err, "Should get 'no row found' type of error")
		assert.Equal(t, *deletedUserEncx, domain.UserEncx{})
	})

	t.Run("should successfully delete own account with Stripe customer", func(t *testing.T) {
		// Clean state
		td.ClearUsersTable(t, ctx, testPool)
		td.ClearSessionsRedis(t, ctx, redisClient)

		// Create user with Stripe customer ID
		user := td.NewTestUser(t, "stripe-self@example.com", "Stripe", "Self")
		user.State = domain.Active
		user.StripeCustomerID = "cus_self_test456"

		userEncx, err := domain.ProcessUserEncx(ctx, crypto, user)
		require.NoError(t, err)

		err = td.InsertUserEncx(t, ctx, userEncx, testPool, crypto)
		require.NoError(t, err)

		// Create session for user
		sessionInfo := &session.SessionInfo{
			ID:     uuid.New(),
			UserID: user.ID,
			Role:   identity.Standard,
			State:  session.SessionActive,
		}
		accessToken := td.CreateSessionWithEncryption(t, ctx, sessionInfo, redisClient, crypto)

		// Act - delete own account with Stripe customer
		req := td.NewDeleteOwnAccountRequest(t, ctx, testServerURL, accessToken)
		resp, err := client.Do(req)

		// Assert
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Verify user is deleted (Stripe deletion handled internally)
		deletedUserEncx, err := td.GetUserEnxByID(t, ctx, user.ID, testPool, crypto)
		assert.Error(t, err)

		assert.Equal(t, *deletedUserEncx, domain.UserEncx{})
	})

	t.Run("should handle user not found error", func(t *testing.T) {
		// Clean state
		td.ClearUsersTable(t, ctx, testPool)
		td.ClearSessionsRedis(t, ctx, redisClient)

		// Create a session for a user that doesn't exist in database
		nonExistentUserID := uuid.New()
		sessionInfo := &session.SessionInfo{
			ID:     uuid.New(),
			UserID: nonExistentUserID,
			Role:   identity.Standard,
			State:  session.SessionActive,
		}
		accessToken := td.CreateSessionWithEncryption(t, ctx, sessionInfo, redisClient, crypto)

		// Act - try to delete account for non-existent user
		req := td.NewDeleteOwnAccountRequest(t, ctx, testServerURL, accessToken)
		resp, err := client.Do(req)

		// Assert
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}
