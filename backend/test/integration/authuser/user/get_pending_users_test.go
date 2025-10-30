package user_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	tu "github.com/Leviosa-care/leviosa/backend/internal/common/testutils"
	td "github.com/Leviosa-care/leviosa/backend/test/helpers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestGetPendingUsers TEST_PATH=test/integration/authuser/user/get_pending_users_test.go

func TestGetPendingUsers(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("should return empty array when no pending users exist", func(t *testing.T) {
		// Clean state
		td.ClearUsersTable(t, ctx, testPool)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

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

		const userCount = 3
		const pendingUserCount = userCount - 1

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		baseUsers := [userCount]basePendingUser{
			{
				email:     "firstpending@example.com",
				firstname: "John",
				lastname:  "DOE",
				state:     domain.Pending,
			},
			{
				email:     "secondpending@example.com",
				firstname: "Jane",
				lastname:  "SMITH",
				state:     domain.Pending,
			},
			{
				email:     "active@example.com",
				firstname: "Active",
				lastname:  "User",
				state:     domain.Active,
			},
		}

		// Create and insert test users with different states
		for _, baseUser := range baseUsers {
			user := td.NewTestUser(t, baseUser.email, baseUser.firstname, baseUser.lastname)
			user.State = baseUser.state

			userEncx, err := domain.ProcessUserEncx(ctx, crypto, user)
			require.NoError(t, err)

			err = td.InsertUserEncx(t, ctx, userEncx, testPool)
			require.NoError(t, err)
		}

		// Act
		req := td.NewGetPendingUsersRequest(t, ctx, testServerURL, accessToken)
		resp, err := client.Do(req)

		// Assert
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		pendingUsers := td.ParseGetPendingUsersResponse(t, resp)
		assert.Len(t, pendingUsers, pendingUserCount, "Should return only pending users") // the -1 is to remove the active user

		// Verify user data (should be decrypted in response)
		emails := make(map[string]struct{}, len(pendingUsers))
		states := make(map[domain.UserState]int, len(pendingUsers))
		for _, pendingUser := range pendingUsers {
			emails[pendingUser.Email] = struct{}{}
			states[pendingUser.State]++
		}

		for _, bu := range baseUsers {
			if bu.state == domain.Pending {
				_, exists := emails[bu.email]
				assert.Truef(t, exists, "expected email %q to be in users", bu.email)
			}
		}

		assert.Equal(t, states[domain.Pending], pendingUserCount) // 2 pending users and 1 active
	})

	t.Run("should return users ordered by creation date descending", func(t *testing.T) {
		// Clean state
		td.ClearUsersTable(t, ctx, testPool)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		const userCount = 3
		const pendingUserCount = userCount - 1

		lastname := "User"

		basePendingUsers := [userCount]basePendingUser{
			{email: "firstpendinguser@example.com", firstname: "First"},
			{email: "secondpendinguser@example.com", firstname: "Second"},
			{email: "thirdpendinguser@example.com", firstname: "Third"},
		}

		// Create and insert test users with different states
		for _, basePendingUser := range basePendingUsers {
			var user *domain.User
			user = td.NewTestUser(t, basePendingUser.email, basePendingUser.firstname, lastname)
			user.State = domain.Pending
			userEncx, err := domain.ProcessUserEncx(ctx, crypto, user)
			require.NoError(t, err)
			err = td.InsertUserEncx(t, ctx, userEncx, testPool)
			require.NoError(t, err)
			time.Sleep(10 * time.Millisecond)
		}

		// Act
		req := td.NewGetPendingUsersRequest(t, ctx, testServerURL, accessToken)
		resp, err := client.Do(req)

		// Assert
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		users := td.ParseGetPendingUsersResponse(t, resp)
		require.Len(t, users, userCount)

		// Should be ordered by creation date descending (newest first)
		for i, bu := range basePendingUsers {
			assert.Equal(t, bu.email, users[userCount-1-i].Email)
		}
	})

	t.Run("should properly decrypt and return all user fields", func(t *testing.T) {
		// Clean state
		td.ClearUsersTable(t, ctx, testPool)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		email := "complete@example.com"
		firstname := "John"
		lastname := "Doe"
		state := domain.Pending
		telephone := "1234567890"

		picture := "profile.jpg"
		gender := domain.GenderMan.String()
		address1 := "123 Main St"
		city := "New York"
		postalCode := "10001"

		layout := "2006-01-02"
		date := "1990-01-01"
		expectedBirthDate, err := time.ParseInLocation(layout, date, time.Local)
		require.NoError(t, err)

		// Insert test user with all fields populated
		testUser := td.NewTestUser(t, email, firstname, lastname)
		testUser.State = state
		testUser.Telephone = telephone
		testUser.Picture = picture
		testUser.BirthDate = expectedBirthDate
		testUser.Gender = gender
		testUser.Address1 = address1
		testUser.City = city
		testUser.PostalCode = postalCode

		testUserEncx, err := domain.ProcessUserEncx(ctx, crypto, testUser)
		require.NoError(t, err)

		err = td.InsertUserEncx(t, ctx, testUserEncx, testPool)
		require.NoError(t, err)

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
		assert.Equal(t, email, user.Email)
		assert.Equal(t, firstname, user.FirstName)
		assert.Equal(t, lastname, user.LastName)
		assert.Equal(t, telephone, user.Telephone)
		assert.Equal(t, picture, user.Picture)
		assert.Equal(t, expectedBirthDate, user.BirthDate)
		assert.Equal(t, gender, user.Gender)
		assert.Equal(t, address1, user.Address1)
		assert.Equal(t, city, user.City)
		assert.Equal(t, postalCode, user.PostalCode)
		assert.Equal(t, domain.Pending, user.State)
	})

	t.Run("should handle users with optional fields as empty", func(t *testing.T) {
		// Clean state
		td.ClearUsersTable(t, ctx, testPool)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		email := "minimal@example.com"
		firstname := "Min"
		lastname := "User"
		state := domain.Pending

		// Insert minimal user (only required fields)
		minimalUser := td.NewTestUser(t, email, firstname, lastname)
		minimalUser.State = state
		// Leave optional fields empty

		minimalUserEncx, err := domain.ProcessUserEncx(ctx, crypto, minimalUser)
		require.NoError(t, err)

		err = td.InsertUserEncx(t, ctx, minimalUserEncx, testPool)
		require.NoError(t, err)

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
		assert.Equal(t, email, user.Email)
		assert.Equal(t, firstname, user.FirstName)
		assert.Equal(t, lastname, user.LastName)
		assert.Empty(t, user.Picture)
		assert.Empty(t, user.Gender)
		assert.Empty(t, user.Address1)
		assert.Empty(t, user.City)
		assert.Empty(t, user.PostalCode)
	})

	t.Run("should handle concurrent requests properly", func(t *testing.T) {
		// Clean state
		td.ClearUsersTable(t, ctx, testPool)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		const userCount = 5

		// Insert multiple pending users
		for i := range userCount {
			user := td.NewTestUser(t, fmt.Sprintf("concurrent%d@example.com", i), "User", fmt.Sprintf("%d", i))
			user.State = domain.Pending
			userEncx, err := domain.ProcessUserEncx(ctx, crypto, user)
			require.NoError(t, err)
			err = td.InsertUserEncx(t, ctx, userEncx, testPool)
			require.NoError(t, err)
		}

		const concurrentRequestCount = 3

		// Make concurrent requests
		responses := make(chan *http.Response, concurrentRequestCount)
		for range concurrentRequestCount {
			go func() {
				req := td.NewGetPendingUsersRequest(t, ctx, testServerURL, accessToken)
				resp, err := client.Do(req)
				require.NoError(t, err)
				responses <- resp
			}()
		}

		// Collect and verify all responses
		successCount := 0
		for range concurrentRequestCount {
			resp := <-responses
			defer resp.Body.Close()

			if resp.StatusCode == http.StatusOK {
				successCount++
				users := td.ParseGetPendingUsersResponse(t, resp)
				assert.Len(t, users, userCount, "Each concurrent request should return all pending users")
			}
		}

		assert.Equal(t, concurrentRequestCount, successCount, "All concurrent requests should succeed")
	})

	t.Run("should exclude users with unverified state", func(t *testing.T) {
		// Clean state
		td.ClearUsersTable(t, ctx, testPool)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		email := "pending@example.com"

		// Insert users with different states including unverified
		pendingUser := td.NewTestUser(t, email, "Pending", "User")
		pendingUser.State = domain.Pending
		pendingUserEncx, err := domain.ProcessUserEncx(ctx, crypto, pendingUser)
		require.NoError(t, err)
		err = td.InsertUserEncx(t, ctx, pendingUserEncx, testPool)
		require.NoError(t, err)

		unverifiedUser := td.NewTestUser(t, "unverified@example.com", "Unverified", "User")
		unverifiedUser.State = domain.Unverified
		unverifiedUserEncx, err := domain.ProcessUserEncx(ctx, crypto, unverifiedUser)
		require.NoError(t, err)
		err = td.InsertUserEncx(t, ctx, unverifiedUserEncx, testPool)
		require.NoError(t, err)

		// Act
		req := td.NewGetPendingUsersRequest(t, ctx, testServerURL, accessToken)
		resp, err := client.Do(req)

		// Assert
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		users := td.ParseGetPendingUsersResponse(t, resp)
		require.Len(t, users, 1)
		assert.Equal(t, email, users[0].Email)
		assert.Equal(t, domain.Pending, users[0].State)
	})

	t.Run("should handle large number of pending users efficiently", func(t *testing.T) {
		// Clean state
		td.ClearUsersTable(t, ctx, testPool)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		// Insert many pending users
		const userCount = 50

		for i := range userCount {
			user := td.NewTestUser(t, fmt.Sprintf("bulk%d@example.com", i), "Bulk", fmt.Sprintf("User%d", i))
			user.State = domain.Pending
			userEncx, err := domain.ProcessUserEncx(ctx, crypto, user)
			require.NoError(t, err)
			err = td.InsertUserEncx(t, ctx, userEncx, testPool)
			require.NoError(t, err)
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

type basePendingUser struct {
	email     string
	firstname string
	lastname  string
	state     domain.UserState
}
