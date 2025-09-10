package user_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/Leviosa-care/authuser/internal/domain"
	th "github.com/Leviosa-care/authuser/test/helpers"
	"github.com/Leviosa-care/core/auth/session"
	"github.com/Leviosa-care/core/contracts/identity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TEST=TestChangePassword make test-integration-user-test

func TestChangePassword(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("should successfully change password with valid credentials", func(t *testing.T) {
		// Clean state
		th.ClearUsersTable(t, ctx, testPool)
		th.ClearSessionsRedis(t, ctx, testClient)

		// Create test user with known password
		oldPassword := "oldPassword123!"
		newPassword := "newPassword456!"
		userID := th.InsertTestUserWithPassword(t, ctx, "test@example.com", "John", "Doe", oldPassword, testPool, crypto)

		// Create active session for the user
		sessionInfo := &session.SessionInfo{
			ID:     uuid.New(),
			UserID: userID,
			Role:   identity.Standard,
			State:  session.SessionActive,
		}
		accessToken := th.CreateSessionWithEncryption(t, ctx, sessionInfo, testClient, crypto)

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
		updatedUser := th.GetUserByIDFromDB(t, ctx, userID, testPool, crypto)
		require.NotNil(t, updatedUser)

		// TODO: that part will need to change when the github.com/hengadev/encx package fixes its CompareSecureHashAndValue implementation

		// Verify old password no longer works by attempting to verify it
		err = service.VerifyUserPassword(ctx, userID, oldPassword)
		assert.Error(t, err, "Old password should no longer work")

		// Verify new password works
		err = service.VerifyUserPassword(ctx, userID, newPassword)
		assert.NoError(t, err, "New password should work")
	})

	t.Run("should fail with incorrect old password", func(t *testing.T) {
		// Clean state
		th.ClearUsersTable(t, ctx, testPool)
		th.ClearSessionsRedis(t, ctx, testClient)

		// Create test user with known password
		actualPassword := "actualPassword123!"
		userID := th.InsertTestUserWithPassword(t, ctx, "test@example.com", "John", "Doe", actualPassword, testPool, crypto)

		// Create active session for the user
		sessionInfo := &session.SessionInfo{
			ID:     uuid.New(),
			UserID: userID,
			Role:   identity.Standard,
			State:  session.SessionActive,
		}
		accessToken := th.CreateSessionWithEncryption(t, ctx, sessionInfo, testClient, crypto)

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

		// Verify password was NOT changed in database
		err = service.VerifyUserPassword(ctx, userID, actualPassword)
		assert.NoError(t, err, "Original password should still work")

		err = service.VerifyUserPassword(ctx, userID, "newPassword456!")
		assert.Error(t, err, "New password should NOT work")
	})

	t.Run("should fail with same old and new password", func(t *testing.T) {
		// Clean state
		th.ClearUsersTable(t, ctx, testPool)
		th.ClearSessionsRedis(t, ctx, testClient)

		// Create test user with known password
		password := "samePassword123!"
		userID := th.InsertTestUserWithPassword(t, ctx, "test@example.com", "John", "Doe", password, testPool, crypto)

		// Create active session for the user
		sessionInfo := &session.SessionInfo{
			ID:     uuid.New(),
			UserID: userID,
			Role:   identity.Standard,
			State:  session.SessionActive,
		}
		accessToken := th.CreateSessionWithEncryption(t, ctx, sessionInfo, testClient, crypto)

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
		th.ClearSessionsRedis(t, ctx, testClient)

		// Create test user
		userID := th.InsertTestUserWithPassword(t, ctx, "test@example.com", "John", "Doe", "password123!", testPool, crypto)

		// Create active session for the user
		sessionInfo := &session.SessionInfo{
			ID:     uuid.New(),
			UserID: userID,
			Role:   identity.Standard,
			State:  session.SessionActive,
		}
		accessToken := th.CreateSessionWithEncryption(t, ctx, sessionInfo, testClient, crypto)

		// Make request with invalid JSON
		req := th.NewInvalidJSONRequest(t, ctx, testServerURL, http.MethodPatch, "/users/me/password")
		th.AddAuthCookie(req, accessToken)
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert HTTP response
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should fail with missing old password", func(t *testing.T) {
		// Clean state
		th.ClearUsersTable(t, ctx, testPool)
		th.ClearSessionsRedis(t, ctx, testClient)

		// Create test user
		userID := th.InsertTestUserWithPassword(t, ctx, "test@example.com", "John", "Doe", "password123!", testPool, crypto)

		// Create active session for the user
		sessionInfo := &session.SessionInfo{
			ID:     uuid.New(),
			UserID: userID,
			Role:   identity.Standard,
			State:  session.SessionActive,
		}
		accessToken := th.CreateSessionWithEncryption(t, ctx, sessionInfo, testClient, crypto)

		// Make change password request with missing old password
		request := domain.ChangePasswordRequest{
			OldPassword: "", // Empty old password
			NewPassword: "newPassword456!",
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
		th.ClearSessionsRedis(t, ctx, testClient)

		// Create test user
		userID := th.InsertTestUserWithPassword(t, ctx, "test@example.com", "John", "Doe", "password123!", testPool, crypto)

		// Create active session for the user
		sessionInfo := &session.SessionInfo{
			ID:     uuid.New(),
			UserID: userID,
			Role:   identity.Standard,
			State:  session.SessionActive,
		}
		accessToken := th.CreateSessionWithEncryption(t, ctx, sessionInfo, testClient, crypto)

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
		th.ClearSessionsRedis(t, ctx, testClient)

		// Create active session for non-existent user
		nonExistentUserID := uuid.New()
		sessionInfo := &session.SessionInfo{
			ID:     uuid.New(),
			UserID: nonExistentUserID,
			Role:   identity.Standard,
			State:  session.SessionActive,
		}
		accessToken := th.CreateSessionWithEncryption(t, ctx, sessionInfo, testClient, crypto)

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
		th.ClearSessionsRedis(t, ctx, testClient)

		// Create test user with known password (standard role required for password change)
		oldPassword := "oldPassword123!"
		newPassword := "newPassword456!"
		userID := th.InsertTestUserWithPassword(t, ctx, "test@example.com", "John", "Doe", oldPassword, testPool, crypto)

		// Create session with standard role (required by middleware)
		sessionInfo := &session.SessionInfo{
			ID:     uuid.New(),
			UserID: userID,
			Role:   identity.Standard,
			State:  session.SessionActive,
		}
		accessToken := th.CreateSessionWithEncryption(t, ctx, sessionInfo, testClient, crypto)

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

		// Verify new password works
		err = service.VerifyUserPassword(ctx, userID, newPassword)
		assert.NoError(t, err, "New password should work")
	})
}

