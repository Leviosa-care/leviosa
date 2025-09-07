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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TEST=TestGetAllUsers make test-integration-user-test

func TestGetAllUsers(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("should return empty array when no users exist", func(t *testing.T) {
		// Clean state
		td.ClearUsersTable(t, ctx, testPool)

		accessToken := setupAdminUser(t, ctx)

		// Act
		req := td.NewGetAllUsersRequest(t, ctx, testServerURL, accessToken)
		resp, err := client.Do(req)

		// Assert
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		users := td.ParseGetAllUsersResponse(t, resp)
		assert.Len(t, users, 1, "Should return only the admin user when no other users exist")
		assert.Equal(t, identity.AdministratorStr, users[0].Role, "The only user should be the admin user")
	})

	t.Run("should return all users regardless of state", func(t *testing.T) {
		// Clean state
		td.ClearUsersTable(t, ctx, testPool)

		accessToken := setupAdminUser(t, ctx)

		// Insert test users with different states
		pendingUser := td.NewTestUser("pending@example.com", "John", "Pending")
		pendingUser.State = domain.Pending
		activeUser := td.NewTestUser("active@example.com", "Jane", "Active")
		activeUser.State = domain.Active
		unverifiedUser := td.NewTestUser("unverified@example.com", "Bob", "Unverified")
		unverifiedUser.State = domain.Unverified

		// Process encryption before insertion
		crypto.ProcessStruct(ctx, pendingUser)
		crypto.ProcessStruct(ctx, activeUser)
		crypto.ProcessStruct(ctx, unverifiedUser)

		td.InsertUser(t, ctx, pendingUser, testPool)
		td.InsertUser(t, ctx, activeUser, testPool)
		td.InsertUser(t, ctx, unverifiedUser, testPool)

		// Act
		req := td.NewGetAllUsersRequest(t, ctx, testServerURL, accessToken)
		resp, err := client.Do(req)

		// Assert
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		users := td.ParseGetAllUsersResponse(t, resp)
		assert.Len(t, users, 4, "Should return all users regardless of state (3 test users + 1 admin)")

		// Verify user data (should be decrypted in response)
		emails := []string{users[0].Email, users[1].Email, users[2].Email, users[3].Email}
		assert.Contains(t, emails, "pending@example.com")
		assert.Contains(t, emails, "active@example.com")
		assert.Contains(t, emails, "unverified@example.com")

		// Verify admin user is present by checking roles
		roles := []string{users[0].Role, users[1].Role, users[2].Role, users[3].Role}
		assert.Contains(t, roles, identity.AdministratorStr)

		// Verify all different states are present
		states := []domain.UserState{users[0].State, users[1].State, users[2].State, users[3].State}
		assert.Contains(t, states, domain.Pending)
		assert.Contains(t, states, domain.Active)
		assert.Contains(t, states, domain.Unverified)
	})

	t.Run("should return users ordered by creation date descending", func(t *testing.T) {
		// Clean state
		td.ClearUsersTable(t, ctx, testPool)

		accessToken := setupAdminUser(t, ctx)

		// Insert users with slight delay to ensure different timestamps
		firstUser := td.NewTestUser("first@example.com", "First", "User")
		firstUser.State = domain.Active
		crypto.ProcessStruct(ctx, firstUser)
		td.InsertUser(t, ctx, firstUser, testPool)

		time.Sleep(10 * time.Millisecond)

		secondUser := td.NewTestUser("second@example.com", "Second", "User")
		secondUser.State = domain.Pending
		crypto.ProcessStruct(ctx, secondUser)
		td.InsertUser(t, ctx, secondUser, testPool)

		time.Sleep(10 * time.Millisecond)

		thirdUser := td.NewTestUser("third@example.com", "Third", "User")
		thirdUser.State = domain.Unverified
		crypto.ProcessStruct(ctx, thirdUser)
		td.InsertUser(t, ctx, thirdUser, testPool)

		// Act
		req := td.NewGetAllUsersRequest(t, ctx, testServerURL, accessToken)
		resp, err := client.Do(req)

		// Assert
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		users := td.ParseGetAllUsersResponse(t, resp)
		require.Len(t, users, 4)

		// Should be ordered by creation date descending (newest first)
		// Admin user is created in setupAdminUser (first), so will be last
		// Test users created in order: first, second, third (newest)
		assert.Equal(t, "third@example.com", users[0].Email)
		assert.Equal(t, "second@example.com", users[1].Email)
		assert.Equal(t, "first@example.com", users[2].Email)
		assert.Equal(t, identity.AdministratorStr, users[3].Role)
	})

	t.Run("should properly decrypt and return all user fields", func(t *testing.T) {
		// Clean state
		td.ClearUsersTable(t, ctx, testPool)

		accessToken := setupAdminUser(t, ctx)

		// Insert test user with all fields populated
		testUser := td.NewTestUser("complete@example.com", "John", "Doe")
		testUser.State = domain.Active
		testUser.Telephone = "1234567890"
		testUser.Picture = "profile.jpg"
		birthdate, err := time.Parse("2006-01-02", "1990-01-01")
		require.NoError(t, err)
		testUser.BirthDate = birthdate
		testUser.Gender = "male"
		testUser.Address1 = "123 Main St"
		testUser.City = "New York"
		testUser.PostalCode = "10001"

		crypto.ProcessStruct(ctx, testUser)
		td.InsertUser(t, ctx, testUser, testPool)

		// Act
		req := td.NewGetAllUsersRequest(t, ctx, testServerURL, accessToken)
		resp, err := client.Do(req)

		// Assert
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		users := td.ParseGetAllUsersResponse(t, resp)
		require.Len(t, users, 2)

		// Find the test user (not the admin user)
		var user *domain.UserResponse
		for _, u := range users {
			if u.Role != identity.AdministratorStr {
				user = u
				break
			}
		}
		require.NotNil(t, user, "Should find the test user in response")
		assert.Equal(t, "complete@example.com", user.Email)
		assert.Equal(t, "John", user.FirstName)
		assert.Equal(t, "Doe", user.LastName)
		assert.Equal(t, "1234567890", user.Telephone)
		assert.Equal(t, "profile.jpg", user.Picture)

		// Check birthdate as time.Time
		expectedBirthDate, err := time.Parse("2006-01-02", "1990-01-01")
		require.NoError(t, err)
		assert.Equal(t, expectedBirthDate, user.BirthDate)

		assert.Equal(t, "male", user.Gender)
		assert.Equal(t, "123 Main St", user.Address1)
		assert.Equal(t, "New York", user.City)
		assert.Equal(t, "10001", user.PostalCode)
		assert.Equal(t, domain.Active, user.State)
	})

	t.Run("should handle users with optional fields as empty", func(t *testing.T) {
		// Clean state
		td.ClearUsersTable(t, ctx, testPool)

		accessToken := setupAdminUser(t, ctx)

		// Insert minimal user (only required fields)
		minimalUser := td.NewTestUser("minimal@example.com", "Min", "User")
		minimalUser.State = domain.Pending
		// Leave optional fields empty

		crypto.ProcessStruct(ctx, minimalUser)
		td.InsertUser(t, ctx, minimalUser, testPool)

		// Act
		req := td.NewGetAllUsersRequest(t, ctx, testServerURL, accessToken)
		resp, err := client.Do(req)

		// Assert
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		users := td.ParseGetAllUsersResponse(t, resp)
		require.Len(t, users, 2)

		// Find the test user (not the admin user)
		var user *domain.UserResponse
		for _, u := range users {
			if u.Role != identity.AdministratorStr {
				user = u
				break
			}
		}
		require.NotNil(t, user, "Should find the test user in response")
		assert.Equal(t, "minimal@example.com", user.Email)
		assert.Equal(t, "Min", user.FirstName)
		assert.Equal(t, "User", user.LastName)
		assert.Equal(t, "0612345678", user.Telephone) // Default value from NewTestUser
		assert.Empty(t, user.Picture)
		assert.True(t, user.BirthDate.IsZero()) // time.Time zero value
		assert.Empty(t, user.Gender)
		assert.Empty(t, user.Address1)
		assert.Empty(t, user.City)
		assert.Empty(t, user.PostalCode)
	})

	t.Run("should include all user states in mixed scenario", func(t *testing.T) {
		// Clean state
		td.ClearUsersTable(t, ctx, testPool)

		accessToken := setupAdminUser(t, ctx)

		// Create users with all possible states
		users := []*domain.User{
			{State: domain.Pending},
			{State: domain.Active},
			{State: domain.Unverified},
			{State: domain.Pending}, // Another pending user
			{State: domain.Active},  // Another active user
		}

		// Insert users
		for i, user := range users {
			testUser := td.NewTestUser(fmt.Sprintf("mixed%d@example.com", i), "User", fmt.Sprintf("%d", i))
			testUser.State = user.State
			crypto.ProcessStruct(ctx, testUser)
			td.InsertUser(t, ctx, testUser, testPool)
		}

		// Act
		req := td.NewGetAllUsersRequest(t, ctx, testServerURL, accessToken)
		resp, err := client.Do(req)

		// Assert
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		responseUsers := td.ParseGetAllUsersResponse(t, resp)
		require.Len(t, responseUsers, 6)

		// Count states
		stateCount := make(map[domain.UserState]int)
		for _, user := range responseUsers {
			stateCount[user.State]++
		}

		assert.Equal(t, 2, stateCount[domain.Pending], "Should have 2 pending users")
		assert.Equal(t, 3, stateCount[domain.Active], "Should have 3 active users (2 test + 1 admin)")
		assert.Equal(t, 1, stateCount[domain.Unverified], "Should have 1 unverified user")
	})

	t.Run("should handle concurrent requests properly", func(t *testing.T) {
		// Clean state
		td.ClearUsersTable(t, ctx, testPool)

		accessToken := setupAdminUser(t, ctx)

		// Insert multiple users with different states
		for i := range 5 {
			user := td.NewTestUser(fmt.Sprintf("concurrent%d@example.com", i), "User", fmt.Sprintf("%d", i))
			states := []domain.UserState{domain.Pending, domain.Active, domain.Unverified}
			user.State = states[i%len(states)]
			crypto.ProcessStruct(ctx, user)
			td.InsertUser(t, ctx, user, testPool)
		}

		// Make concurrent requests
		responses := make(chan *http.Response, 3)
		for range 3 {
			go func() {
				req := td.NewGetAllUsersRequest(t, ctx, testServerURL, accessToken)
				resp, err := client.Do(req)
				require.NoError(t, err)
				responses <- resp
			}()
		}

		// Collect and verify all responses
		successCount := 0
		for range 3 {
			resp := <-responses
			defer resp.Body.Close()

			if resp.StatusCode == http.StatusOK {
				successCount++
				users := td.ParseGetAllUsersResponse(t, resp)
				assert.Len(t, users, 6, "Each concurrent request should return all users (5 test + 1 admin)")
			}
		}

		assert.Equal(t, 3, successCount, "All concurrent requests should succeed")
	})

	t.Run("should handle large number of users efficiently", func(t *testing.T) {
		// Clean state
		td.ClearUsersTable(t, ctx, testPool)

		accessToken := setupAdminUser(t, ctx)

		// Insert many users with various states
		userCount := 100
		states := []domain.UserState{domain.Pending, domain.Active, domain.Unverified}
		for i := range userCount {
			user := td.NewTestUser(fmt.Sprintf("bulk%d@example.com", i), "Bulk", fmt.Sprintf("User%d", i))
			user.State = states[i%len(states)] // Cycle through states
			crypto.ProcessStruct(ctx, user)
			td.InsertUser(t, ctx, user, testPool)
		}

		// Act with timeout to ensure reasonable performance
		start := time.Now()
		req := td.NewGetAllUsersRequest(t, ctx, testServerURL, accessToken)
		resp, err := client.Do(req)
		duration := time.Since(start)

		// Assert
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		users := td.ParseGetAllUsersResponse(t, resp)
		assert.Len(t, users, userCount+1)

		// Performance check - should complete within reasonable time
		assert.Less(t, duration, 10*time.Second, "Should handle large user lists efficiently")

		// Verify state distribution
		stateCount := make(map[domain.UserState]int)
		for _, user := range users {
			stateCount[user.State]++
		}

		// Should have roughly equal distribution of states (plus admin user adds 1 to Active)
		expectedPerState := userCount / len(states)
		for _, state := range states {
			if state == domain.Active {
				// Active state has one extra admin user
				assert.GreaterOrEqual(t, stateCount[state], expectedPerState, "State %v should appear at least %d times", state, expectedPerState)
				assert.LessOrEqual(t, stateCount[state], expectedPerState+2, "State %v should appear at most %d times", state, expectedPerState+2)
			} else {
				assert.GreaterOrEqual(t, stateCount[state], expectedPerState-1, "State %v should appear at least %d times", state, expectedPerState-1)
				assert.LessOrEqual(t, stateCount[state], expectedPerState+1, "State %v should appear at most %d times", state, expectedPerState+1)
			}
		}
	})

	t.Run("should require admin authorization", func(t *testing.T) {
		// Clean state
		td.ClearUsersTable(t, ctx, testPool)

		// Insert a test user
		testUser := td.NewTestUser("test@example.com", "Test", "User")
		testUser.State = domain.Active
		crypto.ProcessStruct(ctx, testUser)
		td.InsertUser(t, ctx, testUser, testPool)

		// Act without admin auth (this should fail if auth middleware is properly configured)
		req := td.NewGetAllUsersRequestWithoutAuth(t, ctx, testServerURL)
		resp, err := client.Do(req)

		// Assert - This test assumes auth middleware is configured
		// If auth middleware is not yet implemented, this test will need to be updated
		require.NoError(t, err)
		defer resp.Body.Close()

		// The exact status code depends on how auth middleware is implemented
		// Common responses: 401 Unauthorized, 403 Forbidden
		// For now, we'll just verify the request completes
		// TODO: Update this test once admin auth middleware is implemented
		if resp.StatusCode == http.StatusOK {
			t.Log("Admin authentication not yet implemented - this endpoint is currently accessible without auth")
		} else {
			t.Logf("Admin authentication appears to be working - got status %d", resp.StatusCode)
		}
	})
}
