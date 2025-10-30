package auth_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	td "github.com/Leviosa-care/leviosa/backend/test/helpers"

	ck "github.com/Leviosa-care/leviosa/backend/internal/common/auth/cookies"
	"github.com/Leviosa-care/leviosa/backend/internal/common/auth/session"
	"github.com/Leviosa-care/leviosa/backend/internal/common/contracts/identity"
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
		testUser := td.NewTestUser(t, "deletetest@example.com", "Delete", "User")
		testUser.State = domain.Active

		testUserEncx, err := domain.ProcessUserEncx(ctx, crypto, testUser)
		require.NoError(t, err)

		err = td.InsertUserEncx(t, ctx, testUserEncx, testPool, crypto)
		require.NoError(t, err)

		// Act - make request without authentication
		req := td.NewDeleteUserByAdminRequestWithoutAuth(t, ctx, testServerURL, testUser.ID)
		resp, err := client.Do(req)

		// Assert
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

		// Verify user still exists
		existingUserEncx, err := td.GetUserEnxByID(t, ctx, testUser.ID, testPool, crypto)
		require.NoError(t, err)

		existingUser, err := td.GetUserEnxByID(t, ctx, existingUserEncx.ID, testPool, crypto)
		require.NoError(t, err)

		assert.NotNil(t, existingUser)
	})

	t.Run("should return 403 with non-admin user", func(t *testing.T) {
		// Clean state
		td.ClearUsersTable(t, ctx, testPool)
		td.ClearSessionsRedis(t, ctx, redisClient)

		// Create standard user (non-admin)
		standardUser := td.NewTestUser(t, "standard@example.com", "Standard", "User")
		standardUser.State = domain.Active

		standardUserEncx, err := domain.ProcessUserEncx(ctx, crypto, standardUser)
		require.NoError(t, err)

		err = td.InsertUserEncx(t, ctx, standardUserEncx, testPool, crypto)
		require.NoError(t, err)

		// Create session for standard user
		sessionInfo := &session.SessionInfo{
			ID:     uuid.New(),
			UserID: standardUser.ID,
			Role:   identity.Standard,
			State:  session.SessionActive,
		}
		accessToken := td.CreateSessionWithEncryption(t, ctx, sessionInfo, redisClient, crypto)

		// Create user to delete
		targetUser := td.NewTestUser(t, "target@example.com", "Target", "User")
		targetUser.State = domain.Active

		targetUserEncx, err := domain.ProcessUserEncx(ctx, crypto, targetUser)
		require.NoError(t, err)

		err = td.InsertUserEncx(t, ctx, targetUserEncx, testPool, crypto)
		require.NoError(t, err)

		// Act - make request with standard user token
		req := td.NewDeleteUserByAdminRequest(t, ctx, testServerURL, targetUser.ID, accessToken)
		resp, err := client.Do(req)

		// Assert
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)

		// Verify target user still exists
		existingUserEncx, err := td.GetUserEnxByID(t, ctx, targetUser.ID, testPool, crypto)
		require.NoError(t, err)

		existingUser, err := domain.DecryptUserEncx(ctx, crypto, existingUserEncx)
		require.NoError(t, err)

		assert.NotNil(t, existingUser)
	})

	t.Run("should return 400 for invalid user ID format", func(t *testing.T) {
		// Clean state
		td.ClearUsersTable(t, ctx, testPool)
		td.ClearSessionsRedis(t, ctx, redisClient)

		// Create admin user
		adminUser := td.NewTestUser(t, "admin@example.com", "Admin", "User")
		adminUser.State = domain.Active

		adminUserEncx, err := domain.ProcessUserEncx(ctx, crypto, adminUser)
		require.NoError(t, err)

		err = td.InsertUserEncx(t, ctx, adminUserEncx, testPool, crypto)
		require.NoError(t, err)

		// Create session for admin user
		sessionInfo := &session.SessionInfo{
			ID:     uuid.New(),
			UserID: adminUser.ID,
			Role:   identity.Administrator,
			State:  session.SessionActive,
		}
		accessToken := td.CreateSessionWithEncryption(t, ctx, sessionInfo, redisClient, crypto)

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
		td.ClearSessionsRedis(t, ctx, redisClient)

		// Create admin user
		adminUser := td.NewTestUser(t, "admin@example.com", "Admin", "User")
		adminUser.State = domain.Active

		adminUserEncx, err := domain.ProcessUserEncx(ctx, crypto, adminUser)
		require.NoError(t, err)

		err = td.InsertUserEncx(t, ctx, adminUserEncx, testPool, crypto)
		require.NoError(t, err)

		// Create session for admin user
		sessionInfo := &session.SessionInfo{
			ID:     uuid.New(),
			UserID: adminUser.ID,
			Role:   identity.Administrator,
			State:  session.SessionActive,
		}
		accessToken := td.CreateSessionWithEncryption(t, ctx, sessionInfo, redisClient, crypto)

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
		td.ClearSessionsRedis(t, ctx, redisClient)
		td.ClearOTPKeys(t, ctx, redisClient)

		// Create admin user
		adminUser := td.NewTestUser(t, "admin@example.com", "Admin", "User")
		adminUser.State = domain.Active

		adminUserEncx, err := domain.ProcessUserEncx(ctx, crypto, adminUser)
		require.NoError(t, err)

		err = td.InsertUserEncx(t, ctx, adminUserEncx, testPool, crypto)
		require.NoError(t, err)

		// Create session for admin user
		adminSessionInfo := &session.SessionInfo{
			ID:     uuid.New(),
			UserID: adminUser.ID,
			Role:   identity.Administrator,
			State:  session.SessionActive,
		}
		accessToken := td.CreateSessionWithEncryption(t, ctx, adminSessionInfo, redisClient, crypto)

		// Create user to delete with sessions and OTPs
		targetUser := td.NewTestUser(t, "target@example.com", "Target", "User")
		targetUser.State = domain.Active

		targetUserEncx, err := domain.ProcessUserEncx(ctx, crypto, targetUser)
		require.NoError(t, err)

		err = td.InsertUserEncx(t, ctx, targetUserEncx, testPool, crypto)
		require.NoError(t, err)

		// Create sessions for target user
		targetSessionInfo := &session.SessionInfo{
			ID:     uuid.New(),
			UserID: targetUser.ID,
			Role:   identity.Standard,
			State:  session.SessionActive,
		}
		td.CreateSessionWithEncryption(t, ctx, targetSessionInfo, redisClient, crypto)

		// Create OTP for target user
		td.CreateOTP(t, ctx, targetUser.Email, redisClient, crypto)

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
		deletedUserEncx, err := td.GetUserEnxByID(t, ctx, targetUser.ID, testPool, crypto)
		require.Error(t, err)
		assert.Equal(t, *deletedUserEncx, domain.UserEncx{})

		// Verify sessions are cleared
		targetSession := td.GetSessionByID(t, ctx, targetSessionInfo.ID, redisClient)
		assert.Nil(t, targetSession)

		// Verify OTP is cleared
		_, err = td.GetOTPEncxByEmail(t, ctx, targetUser.Email, redisClient, crypto)
		assert.Error(t, err)

		// Verify admin user still exists
		existingAdminEncx, err := td.GetUserEnxByID(t, ctx, adminUser.ID, testPool, crypto)
		require.NoError(t, err)

		existingAdmin, err := domain.DecryptUserEncx(ctx, crypto, existingAdminEncx)
		require.NoError(t, err)

		assert.NotNil(t, existingAdmin)
	})

	t.Run("should delete user with Stripe customer", func(t *testing.T) {
		// Clean state
		td.ClearUsersTable(t, ctx, testPool)
		td.ClearSessionsRedis(t, ctx, redisClient)

		// Create admin user
		adminUser := td.NewTestUser(t, "admin@example.com", "Admin", "User")
		adminUser.State = domain.Active

		adminUserEncx, err := domain.ProcessUserEncx(ctx, crypto, adminUser)
		require.NoError(t, err)

		err = td.InsertUserEncx(t, ctx, adminUserEncx, testPool, crypto)
		require.NoError(t, err)

		// Create session for admin user
		sessionInfo := &session.SessionInfo{
			ID:     uuid.New(),
			UserID: adminUser.ID,
			Role:   identity.Administrator,
			State:  session.SessionActive,
		}
		accessToken := td.CreateSessionWithEncryption(t, ctx, sessionInfo, redisClient, crypto)

		// Create user with Stripe customer ID
		targetUser := td.NewTestUser(t, "stripe-user@example.com", "Stripe", "User")
		targetUser.State = domain.Active
		targetUser.StripeCustomerID = "cus_test123" // Mock Stripe customer ID

		targetUserEncx, err := domain.ProcessUserEncx(ctx, crypto, targetUser)
		require.NoError(t, err)

		err = td.InsertUserEncx(t, ctx, targetUserEncx, testPool, crypto)
		require.NoError(t, err)

		// Act - delete user with Stripe customer
		req := td.NewDeleteUserByAdminRequest(t, ctx, testServerURL, targetUser.ID, accessToken)
		resp, err := client.Do(req)

		// Assert
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Verify user is deleted
		deletedUserEncx, err := td.GetUserEnxByID(t, ctx, targetUser.ID, testPool, crypto)
		assert.Error(t, err)
		assert.Equal(t, *deletedUserEncx, domain.UserEncx{})
	})
}
