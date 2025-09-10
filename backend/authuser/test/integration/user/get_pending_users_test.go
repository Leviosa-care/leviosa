package user_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/Leviosa-care/authuser/internal/domain"
	td "github.com/Leviosa-care/authuser/test/helpers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TEST=TestGetPendingUsers make test-integration-user-test

func TestGetPendingUsers(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("should return empty array when no pending users exist", func(t *testing.T) {
		// Clean state
		td.ClearUsersTable(t, ctx, testPool)

		accessToken := setupAdminUser(t, ctx)

		// Act
		req := td.NewGetPendingUsersRequest(t, ctx, testServerURL, accessToken)
		resp, err := client.Do(req)

		// Assert
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		users := td.ParseGetPendingUsersResponse(t, resp)
		assert.Empty(t, users, "Should return empty array when no pending users")
	})

	t.Run("should return only pending users", func(t *testing.T) {
		// Clean state
		td.ClearUsersTable(t, ctx, testPool)

		accessToken := setupAdminUser(t, ctx)

		// Insert test users with different states
		pendingUser1 := td.NewTestUser("pending1@example.com", "John", "Doe")
		pendingUser1.State = domain.Pending
		pendingUser2 := td.NewTestUser("pending2@example.com", "Jane", "Smith")
		pendingUser2.State = domain.Pending
		activeUser := td.NewTestUser("active@example.com", "Active", "User")
		activeUser.State = domain.Active

		// Process encryption before insertion
		crypto.ProcessStruct(ctx, pendingUser1)
		crypto.ProcessStruct(ctx, pendingUser2)
		crypto.ProcessStruct(ctx, activeUser)

		td.InsertUser(t, ctx, pendingUser1, testPool)
		td.InsertUser(t, ctx, pendingUser2, testPool)
		td.InsertUser(t, ctx, activeUser, testPool)

		// Act
		req := td.NewGetPendingUsersRequest(t, ctx, testServerURL, accessToken)
		resp, err := client.Do(req)

		// Assert
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		users := td.ParseGetPendingUsersResponse(t, resp)
		assert.Len(t, users, 2, "Should return only pending users")

		// Verify user data (should be decrypted in response)
		emails := []string{users[0].Email, users[1].Email}
		assert.Contains(t, emails, "pending1@example.com")
		assert.Contains(t, emails, "pending2@example.com")

		// Verify all users have pending state
		for _, user := range users {
			assert.Equal(t, domain.Pending, user.State)
		}
	})

	t.Run("should return users ordered by creation date descending", func(t *testing.T) {
		// Clean state
		td.ClearUsersTable(t, ctx, testPool)

		accessToken := setupAdminUser(t, ctx)

		// Insert users with slight delay to ensure different timestamps
		firstUser := td.NewTestUser("first@example.com", "First", "User")
		firstUser.State = domain.Pending
		crypto.ProcessStruct(ctx, firstUser)
		td.InsertUser(t, ctx, firstUser, testPool)

		time.Sleep(10 * time.Millisecond)

		secondUser := td.NewTestUser("second@example.com", "Second", "User")
		secondUser.State = domain.Pending
		crypto.ProcessStruct(ctx, secondUser)
		td.InsertUser(t, ctx, secondUser, testPool)

		time.Sleep(10 * time.Millisecond)

		thirdUser := td.NewTestUser("third@example.com", "Third", "User")
		thirdUser.State = domain.Pending
		crypto.ProcessStruct(ctx, thirdUser)
		td.InsertUser(t, ctx, thirdUser, testPool)

		// Act
		req := td.NewGetPendingUsersRequest(t, ctx, testServerURL, accessToken)
		resp, err := client.Do(req)

		// Assert
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		users := td.ParseGetPendingUsersResponse(t, resp)
		require.Len(t, users, 3)

		// Should be ordered by creation date descending (newest first)
		assert.Equal(t, "third@example.com", users[0].Email)
		assert.Equal(t, "second@example.com", users[1].Email)
		assert.Equal(t, "first@example.com", users[2].Email)
	})

	t.Run("should properly decrypt and return all user fields", func(t *testing.T) {
		// Clean state
		td.ClearUsersTable(t, ctx, testPool)

		accessToken := setupAdminUser(t, ctx)

		// Insert test user with all fields populated
		testUser := td.NewTestUser("complete@example.com", "John", "Doe")
		testUser.State = domain.Pending
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
		req := td.NewGetPendingUsersRequest(t, ctx, testServerURL, accessToken)
		resp, err := client.Do(req)

		// Assert
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		users := td.ParseGetPendingUsersResponse(t, resp)
		require.Len(t, users, 1)

		user := users[0]
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
		assert.Equal(t, domain.Pending, user.State)
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
		req := td.NewGetPendingUsersRequest(t, ctx, testServerURL, accessToken)
		resp, err := client.Do(req)

		// Assert
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		users := td.ParseGetPendingUsersResponse(t, resp)
		require.Len(t, users, 1)

		user := users[0]
		assert.Equal(t, "minimal@example.com", user.Email)
		assert.Equal(t, "Min", user.FirstName)
		assert.Equal(t, "User", user.LastName)
		assert.Equal(t, "0123456789", user.Telephone) // Default value from NewTestUser
		assert.Empty(t, user.Picture)
		assert.True(t, user.BirthDate.IsZero()) // time.Time zero value
		assert.Empty(t, user.Gender)
		assert.Empty(t, user.Address1)
		assert.Empty(t, user.City)
		assert.Empty(t, user.PostalCode)
	})

	t.Run("should handle concurrent requests properly", func(t *testing.T) {
		// Clean state
		td.ClearUsersTable(t, ctx, testPool)

		accessToken := setupAdminUser(t, ctx)

		// Insert multiple pending users
		for i := range 5 {
			user := td.NewTestUser(fmt.Sprintf("concurrent%d@example.com", i), "User", fmt.Sprintf("%d", i))
			user.State = domain.Pending
			crypto.ProcessStruct(ctx, user)
			td.InsertUser(t, ctx, user, testPool)
		}

		// Make concurrent requests
		responses := make(chan *http.Response, 3)
		for range 3 {
			go func() {
				req := td.NewGetPendingUsersRequest(t, ctx, testServerURL, accessToken)
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
				users := td.ParseGetPendingUsersResponse(t, resp)
				assert.Len(t, users, 5, "Each concurrent request should return all pending users")
			}
		}

		assert.Equal(t, 3, successCount, "All concurrent requests should succeed")
	})

	t.Run("should exclude users with unverified state", func(t *testing.T) {
		// Clean state
		td.ClearUsersTable(t, ctx, testPool)

		accessToken := setupAdminUser(t, ctx)

		// Insert users with different states including unverified
		pendingUser := td.NewTestUser("pending@example.com", "Pending", "User")
		pendingUser.State = domain.Pending
		unverifiedUser := td.NewTestUser("unverified@example.com", "Unverified", "User")
		unverifiedUser.State = domain.Unverified

		crypto.ProcessStruct(ctx, pendingUser)
		crypto.ProcessStruct(ctx, unverifiedUser)

		td.InsertUser(t, ctx, pendingUser, testPool)
		td.InsertUser(t, ctx, unverifiedUser, testPool)

		// Act
		req := td.NewGetPendingUsersRequest(t, ctx, testServerURL, accessToken)
		resp, err := client.Do(req)

		// Assert
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		users := td.ParseGetPendingUsersResponse(t, resp)
		require.Len(t, users, 1)
		assert.Equal(t, "pending@example.com", users[0].Email)
		assert.Equal(t, domain.Pending, users[0].State)
	})

	t.Run("should handle large number of pending users efficiently", func(t *testing.T) {
		// Clean state
		td.ClearUsersTable(t, ctx, testPool)

		accessToken := setupAdminUser(t, ctx)

		// Insert many pending users
		userCount := 50
		for i := range userCount {
			user := td.NewTestUser(fmt.Sprintf("bulk%d@example.com", i), "Bulk", fmt.Sprintf("User%d", i))
			user.State = domain.Pending
			crypto.ProcessStruct(ctx, user)
			td.InsertUser(t, ctx, user, testPool)
		}

		// Act with timeout to ensure reasonable performance
		start := time.Now()
		req := td.NewGetPendingUsersRequest(t, ctx, testServerURL, accessToken)
		resp, err := client.Do(req)
		duration := time.Since(start)

		// Assert
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		users := td.ParseGetPendingUsersResponse(t, resp)
		assert.Len(t, users, userCount)

		// Performance check - should complete within reasonable time
		assert.Less(t, duration, 5*time.Second, "Should handle large user lists efficiently")
	})
}
