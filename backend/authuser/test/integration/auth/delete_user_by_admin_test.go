package auth_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/Leviosa-care/authuser/internal/domain"
	td "github.com/Leviosa-care/authuser/test/helpers"

	ck "github.com/Leviosa-care/core/auth/cookies"
	"github.com/Leviosa-care/core/auth/session"
	"github.com/Leviosa-care/core/contracts/identity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TEST=TestDeleteUserByAdmin make test-integration-auth-test

func TestDeleteUserByAdmin(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("should return 401 without authentication", func(t *testing.T) {
		// Clean state
		td.ClearUsersTable(t, ctx, testPool)

		// Create test user to delete
		testUser := td.NewTestUser("deletetest@example.com", "Delete", "User")
		testUser.State = domain.Active
		td.InsertUserWithEncryption(t, ctx, testUser, testPool, crypto)

		// Act - make request without authentication
		req := td.NewDeleteUserByAdminRequestWithoutAuth(t, ctx, testServerURL, testUser.ID)
		resp, err := client.Do(req)

		// Assert
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

		// Verify user still exists
		existingUser := td.GetUserByID(t, ctx, testUser.ID, testPool)
		assert.NotNil(t, existingUser)
	})

	t.Run("should return 403 with non-admin user", func(t *testing.T) {
		// Clean state
		td.ClearUsersTable(t, ctx, testPool)
		td.ClearSessionsRedis(t, ctx, testClient)

		// Create standard user (non-admin)
		standardUser := td.NewTestUser("standard@example.com", "Standard", "User")
		standardUser.State = domain.Active
		td.InsertUserWithEncryption(t, ctx, standardUser, testPool, crypto)

		// Create session for standard user
		sessionInfo := &session.SessionInfo{
			ID:     uuid.New(),
			UserID: standardUser.ID,
			Role:   identity.Standard,
			State:  session.SessionActive,
		}
		accessToken := td.CreateSessionWithEncryption(t, ctx, sessionInfo, testClient, crypto)

		// Create user to delete
		targetUser := td.NewTestUser("target@example.com", "Target", "User")
		targetUser.State = domain.Active
		td.InsertUserWithEncryption(t, ctx, targetUser, testPool, crypto)

		// Act - make request with standard user token
		req := td.NewDeleteUserByAdminRequest(t, ctx, testServerURL, targetUser.ID, accessToken)
		resp, err := client.Do(req)

		// Assert
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)

		// Verify target user still exists
		existingUser := td.GetUserByID(t, ctx, targetUser.ID, testPool)
		assert.NotNil(t, existingUser)
	})

	t.Run("should return 400 for invalid user ID format", func(t *testing.T) {
		// Clean state
		td.ClearUsersTable(t, ctx, testPool)
		td.ClearSessionsRedis(t, ctx, testClient)

		// Create admin user
		adminUser := td.NewTestUser("admin@example.com", "Admin", "User")
		adminUser.State = domain.Active
		td.InsertUserWithEncryption(t, ctx, adminUser, testPool, crypto)

		// Create session for admin user
		sessionInfo := &session.SessionInfo{
			ID:     uuid.New(),
			UserID: adminUser.ID,
			Role:   identity.Administrator,
			State:  session.SessionActive,
		}
		accessToken := td.CreateSessionWithEncryption(t, ctx, sessionInfo, testClient, crypto)

		// Act - make request with invalid user ID
		req, err := http.NewRequestWithContext(
			ctx,
			http.MethodDelete,
			testServerURL+"/admin/auth/users/invalid-uuid",
			nil,
		)
		require.NoError(t, err)
		req.AddCookie(&http.Cookie{Name: ck.AccessTokenCookieName, Value: accessToken})
		resp, err := client.Do(req)

		// Assert
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should return 404 for non-existent user", func(t *testing.T) {
		// Clean state
		td.ClearUsersTable(t, ctx, testPool)
		td.ClearSessionsRedis(t, ctx, testClient)

		// Create admin user
		adminUser := td.NewTestUser("admin@example.com", "Admin", "User")
		adminUser.State = domain.Active
		td.InsertUserWithEncryption(t, ctx, adminUser, testPool, crypto)

		// Create session for admin user
		sessionInfo := &session.SessionInfo{
			ID:     uuid.New(),
			UserID: adminUser.ID,
			Role:   identity.Administrator,
			State:  session.SessionActive,
		}
		accessToken := td.CreateSessionWithEncryption(t, ctx, sessionInfo, testClient, crypto)

		// Act - make request with non-existent user ID
		nonExistentID := uuid.New()
		req := td.NewDeleteUserByAdminRequest(t, ctx, testServerURL, nonExistentID, accessToken)
		resp, err := client.Do(req)

		// Assert
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("should successfully delete existing user", func(t *testing.T) {
		// Clean state
		td.ClearUsersTable(t, ctx, testPool)
		td.ClearSessionsRedis(t, ctx, testClient)
		td.ClearOTPKeys(t, ctx, testClient)

		// Create admin user
		adminUser := td.NewTestUser("admin@example.com", "Admin", "User")
		adminUser.State = domain.Active
		td.InsertUserWithEncryption(t, ctx, adminUser, testPool, crypto)

		// Create session for admin user
		adminSessionInfo := &session.SessionInfo{
			ID:     uuid.New(),
			UserID: adminUser.ID,
			Role:   identity.Administrator,
			State:  session.SessionActive,
		}
		accessToken := td.CreateSessionWithEncryption(t, ctx, adminSessionInfo, testClient, crypto)

		// Create user to delete with sessions and OTPs
		targetUser := td.NewTestUser("target@example.com", "Target", "User")
		targetUser.State = domain.Active
		td.InsertUserWithEncryption(t, ctx, targetUser, testPool, crypto)

		// Create sessions for target user
		targetSessionInfo := &session.SessionInfo{
			ID:     uuid.New(),
			UserID: targetUser.ID,
			Role:   identity.Standard,
			State:  session.SessionActive,
		}
		td.CreateSessionWithEncryption(t, ctx, targetSessionInfo, testClient, crypto)

		// Create OTP for target user
		td.CreateOTP(t, ctx, targetUser.Email, testClient)

		// Act - delete user
		req := td.NewDeleteUserByAdminRequest(t, ctx, testServerURL, targetUser.ID, accessToken)
		resp, err := client.Do(req)

		// Assert
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Verify response message
		message := td.ParseDeleteUserResponse(t, resp)
		assert.Equal(t, "User deleted successfully", message)

		// Verify user is deleted from database
		deletedUser := td.GetUserByID(t, ctx, targetUser.ID, testPool)
		assert.Nil(t, deletedUser)

		// Verify sessions are cleared
		targetSession := td.GetSessionByID(t, ctx, targetSessionInfo.ID, testClient)
		assert.Nil(t, targetSession)

		// Verify OTP is cleared
		otp := td.GetOTP(t, ctx, targetUser.Email, testClient)
		assert.NotEqual(t, domain.OTP{}, otp)

		// Verify admin user still exists
		existingAdmin := td.GetUserByID(t, ctx, adminUser.ID, testPool)
		assert.NotNil(t, existingAdmin)
	})

	t.Run("should delete user with Stripe customer", func(t *testing.T) {
		// Clean state
		td.ClearUsersTable(t, ctx, testPool)
		td.ClearSessionsRedis(t, ctx, testClient)

		// Create admin user
		adminUser := td.NewTestUser("admin@example.com", "Admin", "User")
		adminUser.State = domain.Active
		td.InsertUserWithEncryption(t, ctx, adminUser, testPool, crypto)

		// Create session for admin user
		sessionInfo := &session.SessionInfo{
			ID:     uuid.New(),
			UserID: adminUser.ID,
			Role:   identity.Administrator,
			State:  session.SessionActive,
		}
		accessToken := td.CreateSessionWithEncryption(t, ctx, sessionInfo, testClient, crypto)

		// Create user with Stripe customer ID
		targetUser := td.NewTestUser("stripe-user@example.com", "Stripe", "User")
		targetUser.State = domain.Active
		targetUser.StripeCustomerID = "cus_test123" // Mock Stripe customer ID
		td.InsertUserWithEncryption(t, ctx, targetUser, testPool, crypto)

		// Act - delete user with Stripe customer
		req := td.NewDeleteUserByAdminRequest(t, ctx, testServerURL, targetUser.ID, accessToken)
		resp, err := client.Do(req)

		// Assert
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Verify user is deleted
		deletedUser := td.GetUserByID(t, ctx, targetUser.ID, testPool)
		assert.Nil(t, deletedUser)
	})
}
