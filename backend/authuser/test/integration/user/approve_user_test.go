package user_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/Leviosa-care/authuser/internal/domain"
	td "github.com/Leviosa-care/authuser/test/helpers"

	"github.com/Leviosa-care/core/contracts/identity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TEST=TestApproveUser make test-integration-user-test

func TestApproveUser(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("should successfully approve pending user", func(t *testing.T) {
		// Clean state
		td.ClearUsersTable(t, ctx, testPool)

		accessToken := setupAdminUser(t, ctx)

		// Create a pending user
		pendingUser := td.NewTestUser("pending@example.com", "John", "Doe")
		pendingUser.State = domain.Pending
		td.InsertUserWithEncryption(t, ctx, pendingUser, testPool, crypto)

		// Prepare approval request
		request := domain.ApproveUserRequest{
			UserID: pendingUser.ID,
			Role:   identity.StandardStr,
		}

		// Act
		req := td.NewApproveUserRequest(t, ctx, testServerURL, request, accessToken)
		resp, err := client.Do(req)

		// Assert HTTP response
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Parse response
		response := td.ParseApproveUserResponse(t, resp)
		assert.Equal(t, "User approved successfully", response["message"])

		// Verify user state in database
		user := td.GetUserByIDFromDB(t, ctx, pendingUser.ID, testPool, crypto)
		assert.Equal(t, domain.Active, user.State)
		assert.Equal(t, identity.StandardStr, user.Role)
		// assert.Nil(t, user.RoleEncrypted, "RoleEncrypted should be nil after approval")
	})

	t.Run("should return error when user not found", func(t *testing.T) {
		// Clean state
		td.ClearUsersTable(t, ctx, testPool)

		accessToken := setupAdminUser(t, ctx)

		// Prepare approval request with non-existent user ID
		request := domain.ApproveUserRequest{
			UserID: uuid.New(),
			Role:   identity.StandardStr,
		}

		// Act
		req := td.NewApproveUserRequest(t, ctx, testServerURL, request, accessToken)
		resp, err := client.Do(req)

		// Assert
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("should return error when user is not in pending state", func(t *testing.T) {
		// Clean state
		td.ClearUsersTable(t, ctx, testPool)

		accessToken := setupAdminUser(t, ctx)
		// Create an active user (not pending)
		activeUser := td.NewTestUser("active@example.com", "Jane", "Smith")
		activeUser.State = domain.Active
		td.InsertUserWithEncryption(t, ctx, activeUser, testPool, crypto)

		// Prepare approval request
		request := domain.ApproveUserRequest{
			UserID: activeUser.ID,
			Role:   identity.StandardStr,
		}

		// Act
		req := td.NewApproveUserRequest(t, ctx, testServerURL, request, accessToken)
		resp, err := client.Do(req)

		// Assert
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusConflict, resp.StatusCode) // User not in pending state
	})

	t.Run("should return error when user is unverified", func(t *testing.T) {
		// Clean state
		td.ClearUsersTable(t, ctx, testPool)

		accessToken := setupAdminUser(t, ctx)

		// Create an unverified user
		unverifiedUser := td.NewTestUser("unverified@example.com", "Bob", "Wilson")
		unverifiedUser.State = domain.Unverified
		td.InsertUserWithEncryption(t, ctx, unverifiedUser, testPool, crypto)

		// Prepare approval request
		request := domain.ApproveUserRequest{
			UserID: unverifiedUser.ID,
			Role:   identity.StandardStr,
		}

		// Act
		req := td.NewApproveUserRequest(t, ctx, testServerURL, request, accessToken)
		resp, err := client.Do(req)

		// Assert
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusConflict, resp.StatusCode) // User not in pending state

		// Verify user state unchanged in database
		user := td.GetUserByIDFromDB(t, ctx, unverifiedUser.ID, testPool, crypto)
		assert.Equal(t, domain.Unverified, user.State)
	})

	t.Run("should return error with invalid role", func(t *testing.T) {
		// Clean state
		td.ClearUsersTable(t, ctx, testPool)

		accessToken := setupAdminUser(t, ctx)

		// Create a pending user
		pendingUser := td.NewTestUser("pending2@example.com", "Alice", "Johnson")
		pendingUser.State = domain.Pending
		td.InsertUserWithEncryption(t, ctx, pendingUser, testPool, crypto)

		// Prepare approval request with invalid role
		request := domain.ApproveUserRequest{
			UserID: pendingUser.ID,
			Role:   "invalid_role",
		}

		// Act
		req := td.NewApproveUserRequest(t, ctx, testServerURL, request, accessToken)
		resp, err := client.Do(req)

		// Assert
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		// Verify user state unchanged in database
		user := td.GetUserByIDFromDB(t, ctx, pendingUser.ID, testPool, crypto)
		assert.Equal(t, domain.Pending, user.State)
	})

	t.Run("should return error with malformed JSON", func(t *testing.T) {
		// Clean state
		td.ClearUsersTable(t, ctx, testPool)

		// Act
		accessToken := setupAdminUser(t, ctx)
		req := td.NewMalformedApproveUserRequest(t, ctx, testServerURL, accessToken)
		resp, err := client.Do(req)

		// Assert
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should return error with empty user ID", func(t *testing.T) {
		// Clean state
		td.ClearUsersTable(t, ctx, testPool)

		accessToken := setupAdminUser(t, ctx)

		// Prepare approval request with zero UUID
		request := domain.ApproveUserRequest{
			UserID: uuid.UUID{}, // Zero UUID
			Role:   identity.StandardStr,
		}

		// Act
		req := td.NewApproveUserRequest(t, ctx, testServerURL, request, accessToken)
		resp, err := client.Do(req)

		// Assert
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusNotFound, resp.StatusCode) // Empty UUID will not be found
	})

	t.Run("should approve user with different valid roles", func(t *testing.T) {
		// Clean state
		td.ClearUsersTable(t, ctx, testPool)

		accessToken := setupAdminUser(t, ctx)

		// Test different roles
		testCases := []struct {
			name string
			role string
		}{
			{name: "standard role", role: identity.StandardStr},
			{name: "premium role", role: identity.PremiumStr},
			{name: "partner role", role: identity.PartnerStr},
		}

		for i, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				time.Sleep(500 * time.Millisecond)
				// Create a pending user
				pendingUser := &domain.User{
					ID:    uuid.New(),
					State: domain.Pending,
					// Email:      email,
					Email:      fmt.Sprintf("pending_%s_%d@example.com", tc.role, i),
					FirstName:  fmt.Sprintf("firstname_%d", i),
					LastName:   fmt.Sprintf("lastname_%d", i),
					Password:   fmt.Sprintf("qPDAR0.4Z8{vpCO]_%d", i),
					Telephone:  fmt.Sprintf("061234567%d", i),
					Role:       tc.role,
					CreatedAt:  time.Now(),
					LoggedInAt: time.Now(),
				}

				td.InsertUserWithEncryption(t, ctx, pendingUser, testPool, crypto)

				// Prepare approval request
				request := domain.ApproveUserRequest{
					UserID: pendingUser.ID,
					Role:   tc.role,
				}

				// Act
				req := td.NewApproveUserRequest(t, ctx, testServerURL, request, accessToken)
				resp, err := client.Do(req)

				// Assert
				require.NoError(t, err)
				defer resp.Body.Close()
				assert.Equal(t, http.StatusOK, resp.StatusCode)

				// Verify user state in database
				user := td.GetUserByIDFromDB(t, ctx, pendingUser.ID, testPool, crypto)
				assert.Equal(t, domain.Active, user.State)
				assert.Equal(t, tc.role, user.Role)
			})
		}
	})
}
