package user_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	th "github.com/Leviosa-care/leviosa/backend/test/helpers"

	userEndpoints "github.com/Leviosa-care/leviosa/backend/internal/authuser/interface/user"

	"github.com/Leviosa-care/leviosa/backend/internal/common/auth/session"
	"github.com/Leviosa-care/leviosa/backend/internal/common/contracts/identity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TEST=TestChangePassword make test-integration-user-test

func TestChangePassword(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	email := "test@example.com"
	firstname := "John"
	lastname := "Doe"

	t.Run("should successfully change password with valid credentials", func(t *testing.T) {
		// Clean state
		th.ClearUsersTable(t, ctx, testPool)
		th.ClearSessionsRedis(t, ctx, redisClient)

		// Create test user with known passwords
		oldPassword := "oldPassword123!"
		newPassword := "newPassword456!"

		user := th.NewTestUser(t, email, firstname, lastname)
		user.Password = oldPassword
		userEncx, err := domain.ProcessUserEncx(ctx, crypto, user)

		err = th.InsertUserEncx(t, ctx, userEncx, testPool, crypto)
		require.NoError(t, err)

		// Create active session for the user
		sessionInfo := &session.SessionInfo{
			ID:     uuid.New(),
			UserID: user.ID,
			Role:   identity.Standard,
			State:  session.SessionActive,
		}
		accessToken := th.CreateSessionWithEncryption(t, ctx, sessionInfo, redisClient, crypto)

		// Make change password request
		request := domain.ChangePasswordRequest{
			OldPassword: oldPassword,
			NewPassword: newPassword,
		}
		req := th.NewChangePasswordRequest(t, ctx, testServerURL, request, accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert HTTP response
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		message := th.ParseChangePasswordResponse(t, resp)
		assert.Equal(t, "Password changed successfully", message)

		// Verify password was changed in database
		updatedUser, err := th.GetUserEnxByID(t, ctx, user.ID, testPool, crypto)
		require.NoError(t, err)

		// Verify old password no longer works by attempting to verify it
		match, err := crypto.CompareSecureHashAndValue(ctx, oldPassword, updatedUser.PasswordHashSecure)
		assert.False(t, match)
		assert.NoError(t, err)

		// Verify new password works
		match, err = crypto.CompareSecureHashAndValue(ctx, newPassword, updatedUser.PasswordHashSecure)
		assert.True(t, match)
		assert.NoError(t, err)
	})

	t.Run("should fail with incorrect old password", func(t *testing.T) {
		// Clean state
		th.ClearUsersTable(t, ctx, testPool)
		th.ClearSessionsRedis(t, ctx, redisClient)

		// Create test user with known password
		actualPassword := "actualPassword123!"

		user := th.NewTestUser(t, email, firstname, lastname)
		user.Password = actualPassword
		userEncx, err := domain.ProcessUserEncx(ctx, crypto, user)

		err = th.InsertUserEncx(t, ctx, userEncx, testPool, crypto)
		require.NoError(t, err)

		// Create active session for the user
		sessionInfo := &session.SessionInfo{
			ID:     uuid.New(),
			UserID: user.ID,
			Role:   identity.Standard,
			State:  session.SessionActive,
		}
		accessToken := th.CreateSessionWithEncryption(t, ctx, sessionInfo, redisClient, crypto)

		// Make change password request with wrong old password
		request := domain.ChangePasswordRequest{
			OldPassword: "wrongOldPassword",
			NewPassword: "newPassword456!",
		}
		req := th.NewChangePasswordRequest(t, ctx, testServerURL, request, accessToken)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert HTTP response
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		retrievedUser, err := th.GetUserEnxByID(t, ctx, user.ID, testPool, crypto)
		require.NoError(t, err)

		// Verify password was NOT changed in database
		match, err := crypto.CompareSecureHashAndValue(ctx, actualPassword, retrievedUser.PasswordHashSecure)
		assert.True(t, match)
		assert.NoError(t, err)
	})

	t.Run("should fail with same old and new password", func(t *testing.T) {
		// Clean state
		th.ClearUsersTable(t, ctx, testPool)
		th.ClearSessionsRedis(t, ctx, redisClient)

		// Create test user with known password
		password := "samePassword123!"
		user := th.NewTestUser(t, email, firstname, lastname)
		user.Password = password
		userEncx, err := domain.ProcessUserEncx(ctx, crypto, user)

		err = th.InsertUserEncx(t, ctx, userEncx, testPool, crypto)
		require.NoError(t, err)

		// Create active session for the user
		sessionInfo := &session.SessionInfo{
			ID:     uuid.New(),
			UserID: user.ID,
			Role:   identity.Standard,
			State:  session.SessionActive,
		}
		accessToken := th.CreateSessionWithEncryption(t, ctx, sessionInfo, redisClient, crypto)

		// Make change password request with same old and new password
		request := domain.ChangePasswordRequest{
			OldPassword: password,
			NewPassword: password,
		}
		req := th.NewChangePasswordRequest(t, ctx, testServerURL, request, accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert HTTP response
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should fail without authentication", func(t *testing.T) {
		// Make change password request without authentication cookie
		request := domain.ChangePasswordRequest{
			OldPassword: "oldPassword123!",
			NewPassword: "newPassword456!",
		}
		req := th.NewChangePasswordRequestWithoutAuth(t, ctx, testServerURL, request)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert HTTP response
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("should fail with invalid JSON", func(t *testing.T) {
		// Clean state
		th.ClearUsersTable(t, ctx, testPool)
		th.ClearSessionsRedis(t, ctx, redisClient)

		// Create test user
		password := "password123!"

		user := th.NewTestUser(t, email, firstname, lastname)
		user.Password = password
		userEncx, err := domain.ProcessUserEncx(ctx, crypto, user)

		err = th.InsertUserEncx(t, ctx, userEncx, testPool, crypto)
		require.NoError(t, err)

		// Create active session for the user
		sessionInfo := &session.SessionInfo{
			ID:     uuid.New(),
			UserID: user.ID,
			Role:   identity.Standard,
			State:  session.SessionActive,
		}
		accessToken := th.CreateSessionWithEncryption(t, ctx, sessionInfo, redisClient, crypto)

		// Make request with invalid JSON
		req := th.NewInvalidJSONRequest(t, ctx, testServerURL, http.MethodPatch, userEndpoints.ChangePasswordEndpoint)
		th.AddAuthCookie(req, accessToken)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert HTTP response
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should fail with missing old password in request", func(t *testing.T) {
		// Clean state
		th.ClearUsersTable(t, ctx, testPool)
		th.ClearSessionsRedis(t, ctx, redisClient)

		// Create test user
		oldPassword := "oldPassword123!"
		newPassword := "newPassword456!"

		user := th.NewTestUser(t, email, firstname, lastname)
		user.Password = oldPassword
		userEncx, err := domain.ProcessUserEncx(ctx, crypto, user)

		err = th.InsertUserEncx(t, ctx, userEncx, testPool, crypto)
		require.NoError(t, err)

		// Create active session for the user
		sessionInfo := &session.SessionInfo{
			ID:     uuid.New(),
			UserID: user.ID,
			Role:   identity.Standard,
			State:  session.SessionActive,
		}
		accessToken := th.CreateSessionWithEncryption(t, ctx, sessionInfo, redisClient, crypto)

		// Make change password request with missing old password
		request := domain.ChangePasswordRequest{
			OldPassword: "", // Empty old password
			NewPassword: newPassword,
		}
		req := th.NewChangePasswordRequest(t, ctx, testServerURL, request, accessToken)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert HTTP response
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should fail with missing new password", func(t *testing.T) {
		// Clean state
		th.ClearUsersTable(t, ctx, testPool)
		th.ClearSessionsRedis(t, ctx, redisClient)

		// Create test user
		oldPassword := "oldPassword123!"

		user := th.NewTestUser(t, email, firstname, lastname)
		user.Password = oldPassword
		userEncx, err := domain.ProcessUserEncx(ctx, crypto, user)

		err = th.InsertUserEncx(t, ctx, userEncx, testPool, crypto)
		require.NoError(t, err)

		// Create active session for the user
		sessionInfo := &session.SessionInfo{
			ID:     uuid.New(),
			UserID: user.ID,
			Role:   identity.Standard,
			State:  session.SessionActive,
		}
		accessToken := th.CreateSessionWithEncryption(t, ctx, sessionInfo, redisClient, crypto)

		// Make change password request with missing new password
		request := domain.ChangePasswordRequest{
			OldPassword: "password123!",
			NewPassword: "", // Empty new password
		}
		req := th.NewChangePasswordRequest(t, ctx, testServerURL, request, accessToken)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert HTTP response
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should fail for non-existent user", func(t *testing.T) {
		// Clean state
		th.ClearUsersTable(t, ctx, testPool)
		th.ClearSessionsRedis(t, ctx, redisClient)

		// Create active session for non-existent user
		nonExistentUserID := uuid.New()
		sessionInfo := &session.SessionInfo{
			ID:     uuid.New(),
			UserID: nonExistentUserID,
			Role:   identity.Standard,
			State:  session.SessionActive,
		}
		accessToken := th.CreateSessionWithEncryption(t, ctx, sessionInfo, redisClient, crypto)

		// Make change password request
		request := domain.ChangePasswordRequest{
			OldPassword: "oldPassword123!",
			NewPassword: "newPassword456!",
		}
		req := th.NewChangePasswordRequest(t, ctx, testServerURL, request, accessToken)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert HTTP response
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("should work with visitor role but upgraded to standard", func(t *testing.T) {
		// Clean state
		th.ClearUsersTable(t, ctx, testPool)
		th.ClearSessionsRedis(t, ctx, redisClient)

		// Create test user with known password (standard role required for password change)
		oldPassword := "oldPassword123!"
		newPassword := "newPassword456!"

		user := th.NewTestUser(t, email, firstname, lastname)
		user.Password = oldPassword
		userEncx, err := domain.ProcessUserEncx(ctx, crypto, user)

		err = th.InsertUserEncx(t, ctx, userEncx, testPool, crypto)
		require.NoError(t, err)

		// Create session with standard role (required by middleware)
		sessionInfo := &session.SessionInfo{
			ID:     uuid.New(),
			UserID: user.ID,
			Role:   identity.Standard,
			State:  session.SessionActive,
		}
		accessToken := th.CreateSessionWithEncryption(t, ctx, sessionInfo, redisClient, crypto)

		// Make change password request
		request := domain.ChangePasswordRequest{
			OldPassword: oldPassword,
			NewPassword: newPassword,
		}
		req := th.NewChangePasswordRequest(t, ctx, testServerURL, request, accessToken)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert HTTP response
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		message := th.ParseChangePasswordResponse(t, resp)
		assert.Equal(t, "Password changed successfully", message)

		updatedUser, err := th.GetUserEnxByID(t, ctx, user.ID, testPool, crypto)
		require.NoError(t, err)

		// Verify old password no longer works by attempting to verify it
		match, err := crypto.CompareSecureHashAndValue(ctx, oldPassword, updatedUser.PasswordHashSecure)
		assert.False(t, match)
		assert.NoError(t, err)

		// Verify new password works
		match, err = crypto.CompareSecureHashAndValue(ctx, newPassword, updatedUser.PasswordHashSecure)
		assert.True(t, match)
		assert.NoError(t, err)
	})
}
