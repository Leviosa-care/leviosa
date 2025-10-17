package user_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/Leviosa-care/authuser/internal/domain"
	td "github.com/Leviosa-care/authuser/test/helpers"
	ck "github.com/Leviosa-care/core/auth/cookies"
	tu "github.com/Leviosa-care/core/testutils"

	"github.com/Leviosa-care/core/contracts/identity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TEST=TestUpdateUserRole make test-integration-user-test

func TestUpdateUserRole(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("should successfully update user role from standard to premium", func(t *testing.T) {
		// Clean state
		td.ClearUsersTable(t, ctx, testPool)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		// Create an active user with standard role
		user := td.NewTestUser(t, "user@example.com", "John", "Doe")
		user.State = domain.Active
		user.Role = identity.StandardStr
		userEncx, err := domain.ProcessUserEncx(ctx, crypto, user)
		require.NoError(t, err)
		err = td.InsertUserEncx(t, ctx, userEncx, testPool, crypto)
		require.NoError(t, err)

		// Act - Update role to premium
		req := td.NewUpdateUserRoleRequest(t, ctx, testServerURL, user.ID, identity.PremiumStr, accessToken)
		resp, err := client.Do(req)

		// Assert HTTP response
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Parse response
		response := td.ParseUpdateUserRoleResponse(t, resp)
		assert.Equal(t, "User role updated successfully", response["message"])

		// Verify role updated in database
		updatedUserEncx, err := td.GetUserEnxByID(t, ctx, user.ID, testPool, crypto)
		require.NoError(t, err)

		updatedUser, err := domain.DecryptUserEncx(ctx, crypto, updatedUserEncx)
		require.NoError(t, err)

		assert.Equal(t, domain.Active, updatedUser.State)
		assert.Equal(t, identity.PremiumStr, updatedUser.Role)
		assert.Equal(t, user.Email, updatedUser.Email) // Other fields unchanged
		assert.Equal(t, user.FirstName, updatedUser.FirstName)
	})

	t.Run("should successfully update user role from premium to partner", func(t *testing.T) {
		// Clean state
		td.ClearUsersTable(t, ctx, testPool)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		// Create an active user with premium role
		user := td.NewTestUser(t, "premium@example.com", "Jane", "Smith")
		user.State = domain.Active
		user.Role = identity.PremiumStr
		userEncx, err := domain.ProcessUserEncx(ctx, crypto, user)
		require.NoError(t, err)
		err = td.InsertUserEncx(t, ctx, userEncx, testPool, crypto)
		require.NoError(t, err)

		// Act - Update role to partner
		req := td.NewUpdateUserRoleRequest(t, ctx, testServerURL, user.ID, identity.PartnerStr, accessToken)
		resp, err := client.Do(req)

		// Assert HTTP response
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Parse response
		response := td.ParseUpdateUserRoleResponse(t, resp)
		assert.Equal(t, "User role updated successfully", response["message"])

		// Verify role updated in database
		updatedUserEncx, err := td.GetUserEnxByID(t, ctx, user.ID, testPool, crypto)
		require.NoError(t, err)

		updatedUser, err := domain.DecryptUserEncx(ctx, crypto, updatedUserEncx)
		require.NoError(t, err)

		assert.Equal(t, identity.PartnerStr, updatedUser.Role)
	})

	t.Run("should return error when user not found", func(t *testing.T) {
		// Clean state
		td.ClearUsersTable(t, ctx, testPool)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		// Use non-existent user ID
		nonExistentID := uuid.New()

		// Act
		req := td.NewUpdateUserRoleRequest(t, ctx, testServerURL, nonExistentID, identity.StandardStr, accessToken)
		resp, err := client.Do(req)

		// Assert
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("should return error with invalid user ID format", func(t *testing.T) {
		// Clean state
		td.ClearUsersTable(t, ctx, testPool)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		// Create request with invalid UUID path parameter
		req, err := http.NewRequestWithContext(
			ctx,
			http.MethodPatch,
			fmt.Sprintf("%s/admin/users/invalid-uuid/role", testServerURL),
			nil,
		)
		require.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")
		if accessToken != "" {
			req.AddCookie(&http.Cookie{Name: ck.AccessTokenCookieName, Value: accessToken})
		}

		// Act
		resp, err := client.Do(req)

		// Assert
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should return error with invalid role", func(t *testing.T) {
		// Clean state
		td.ClearUsersTable(t, ctx, testPool)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		// Create an active user
		user := td.NewTestUser(t, "user2@example.com", "Bob", "Wilson")
		user.State = domain.Active
		user.Role = identity.StandardStr
		userEncx, err := domain.ProcessUserEncx(ctx, crypto, user)
		require.NoError(t, err)
		err = td.InsertUserEncx(t, ctx, userEncx, testPool, crypto)
		require.NoError(t, err)

		// Act - Try to update with invalid role
		req := td.NewUpdateUserRoleRequest(t, ctx, testServerURL, user.ID, "invalid_role", accessToken)
		resp, err := client.Do(req)

		// Assert
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		// Verify role unchanged in database
		unchangedUserEncx, err := td.GetUserEnxByID(t, ctx, user.ID, testPool, crypto)
		require.NoError(t, err)

		unchangedUser, err := domain.DecryptUserEncx(ctx, crypto, unchangedUserEncx)
		require.NoError(t, err)

		assert.Equal(t, identity.StandardStr, unchangedUser.Role)
	})

	t.Run("should return error with malformed JSON", func(t *testing.T) {
		// Clean state
		td.ClearUsersTable(t, ctx, testPool)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		// Create an active user
		user := td.NewTestUser(t, "user3@example.com", "Alice", "Johnson")
		user.State = domain.Active

		userEncx, err := domain.ProcessUserEncx(ctx, crypto, user)
		require.NoError(t, err)

		err = td.InsertUserEncx(t, ctx, userEncx, testPool, crypto)
		require.NoError(t, err)

		// Act
		req := td.NewMalformedUpdateUserRoleRequest(t, ctx, testServerURL, user.ID, accessToken)
		resp, err := client.Do(req)

		// Assert
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should return error with empty role", func(t *testing.T) {
		// Clean state
		td.ClearUsersTable(t, ctx, testPool)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		// Create an active user
		user := td.NewTestUser(t, "user@example.com", "Charlie", "Brown")
		user.State = domain.Active
		user.Role = identity.StandardStr

		userEncx, err := domain.ProcessUserEncx(ctx, crypto, user)
		require.NoError(t, err)

		err = td.InsertUserEncx(t, ctx, userEncx, testPool, crypto)
		require.NoError(t, err)

		// Act - Try to update with empty role
		req := td.NewUpdateUserRoleRequest(t, ctx, testServerURL, user.ID, "", accessToken)
		resp, err := client.Do(req)

		// Assert
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		// Verify role unchanged in database
		unchangedUserEncx, err := td.GetUserEnxByID(t, ctx, user.ID, testPool, crypto)
		require.NoError(t, err)

		unchangedUser, err := domain.DecryptUserEncx(ctx, crypto, unchangedUserEncx)
		require.NoError(t, err)

		assert.Equal(t, identity.StandardStr, unchangedUser.Role)
	})

	t.Run("should return error when missing user ID in path", func(t *testing.T) {
		// Clean state
		td.ClearUsersTable(t, ctx, testPool)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		// Create request without user ID in path
		req, err := http.NewRequestWithContext(
			ctx,
			http.MethodPatch,
			fmt.Sprintf("%s/admin/users//role", testServerURL), // Empty user ID
			nil,
		)
		require.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")
		if accessToken != "" {
			req.AddCookie(&http.Cookie{Name: ck.AccessTokenCookieName, Value: accessToken})
		}

		// Act
		resp, err := client.Do(req)

		// Assert
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode) // Route not found
	})

	t.Run("should work with different valid roles", func(t *testing.T) {
		// Clean state
		td.ClearUsersTable(t, ctx, testPool)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		// Test different role updates
		testCases := []struct {
			name     string
			fromRole string
			toRole   string
		}{
			{name: "guest to standard", fromRole: identity.GuestStr, toRole: identity.StandardStr},
			{name: "standard to premium", fromRole: identity.StandardStr, toRole: identity.PremiumStr},
			{name: "premium to partner", fromRole: identity.PremiumStr, toRole: identity.PartnerStr},
			{name: "partner to administrator", fromRole: identity.PartnerStr, toRole: identity.AdministratorStr},
		}

		for i, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// Create a user with initial role
				user := &domain.User{
					ID:         uuid.New(),
					State:      domain.Active,
					Email:      fmt.Sprintf("testuser%d@example.com", i),
					FirstName:  fmt.Sprintf("First%d", i),
					LastName:   fmt.Sprintf("Last%d", i),
					Password:   fmt.Sprintf("password%d123", i),
					Telephone:  fmt.Sprintf("123456%04d", i),
					Role:       tc.fromRole,
					CreatedAt:  time.Now(),
					LoggedInAt: time.Now(),
				}
				userEncx, err := domain.ProcessUserEncx(ctx, crypto, user)
				require.NoError(t, err)

				err = td.InsertUserEncx(t, ctx, userEncx, testPool, crypto)
				require.NoError(t, err)

				// Act - Update role
				req := td.NewUpdateUserRoleRequest(t, ctx, testServerURL, user.ID, tc.toRole, accessToken)
				resp, err := client.Do(req)

				// Assert
				require.NoError(t, err)
				defer resp.Body.Close()
				assert.Equal(t, http.StatusOK, resp.StatusCode)

				// Verify role updated in database
				updatedUserEncx, err := td.GetUserEnxByID(t, ctx, user.ID, testPool, crypto)
				require.NoError(t, err)

				updatedUser, err := domain.DecryptUserEncx(ctx, crypto, updatedUserEncx)
				require.NoError(t, err)

				assert.Equal(t, tc.toRole, updatedUser.Role)
				assert.Equal(t, domain.Active, updatedUser.State) // State unchanged
			})
		}
	})

	t.Run("should handle role update on pending users", func(t *testing.T) {
		// Clean state
		td.ClearUsersTable(t, ctx, testPool)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		// Create a pending user
		pendingUser := td.NewTestUser(t, "pending@example.com", "Pending", "User")
		pendingUser.State = domain.Pending
		pendingUser.Role = identity.GuestStr

		pendingUserEncx, err := domain.ProcessUserEncx(ctx, crypto, pendingUser)
		require.NoError(t, err)

		err = td.InsertUserEncx(t, ctx, pendingUserEncx, testPool, crypto)
		require.NoError(t, err)

		// Act - Update role of pending user
		req := td.NewUpdateUserRoleRequest(t, ctx, testServerURL, pendingUser.ID, identity.StandardStr, accessToken)
		resp, err := client.Do(req)

		// Assert HTTP response
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Parse response
		response := td.ParseUpdateUserRoleResponse(t, resp)
		assert.Equal(t, "User role updated successfully", response["message"])

		// Verify role updated but state remains pending
		updatedUserEncx, err := td.GetUserEnxByID(t, ctx, pendingUser.ID, testPool, crypto)
		require.NoError(t, err)

		updatedUser, err := domain.DecryptUserEncx(ctx, crypto, updatedUserEncx)
		require.NoError(t, err)

		assert.Equal(t, identity.StandardStr, updatedUser.Role)
		assert.Equal(t, domain.Pending, updatedUser.State) // State unchanged
	})

	t.Run("should handle role downgrade", func(t *testing.T) {
		// Clean state
		td.ClearUsersTable(t, ctx, testPool)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		// Create a user with administrator role
		adminUser := td.NewTestUser(t, "admin@example.com", "Admin", "User")
		adminUser.State = domain.Active
		adminUser.Role = identity.AdministratorStr

		adminUserEncx, err := domain.ProcessUserEncx(ctx, crypto, adminUser)
		require.NoError(t, err)

		err = td.InsertUserEncx(t, ctx, adminUserEncx, testPool, crypto)
		require.NoError(t, err)

		// Act - Downgrade from administrator to standard
		req := td.NewUpdateUserRoleRequest(t, ctx, testServerURL, adminUser.ID, identity.StandardStr, accessToken)
		resp, err := client.Do(req)

		// Assert HTTP response
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Verify role downgraded in database
		updatedUserEncx, err := td.GetUserEnxByID(t, ctx, adminUser.ID, testPool, crypto)
		require.NoError(t, err)

		updatedUser, err := domain.DecryptUserEncx(ctx, crypto, updatedUserEncx)
		require.NoError(t, err)

		assert.Equal(t, identity.StandardStr, updatedUser.Role)
		assert.Equal(t, domain.Active, updatedUser.State)
	})
}
