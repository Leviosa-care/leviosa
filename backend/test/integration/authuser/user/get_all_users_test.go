package user_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	td "github.com/Leviosa-care/leviosa/backend/test/helpers"
	"github.com/Leviosa-care/leviosa/backend/internal/common/contracts/identity"
	tu "github.com/Leviosa-care/leviosa/backend/internal/common/testutils"

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

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

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

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		baseUsers := [3]baseUser{
			{email: "pending@example.com", firstname: "John", lastname: "Pending", state: domain.Pending},
			{email: "active@example.com", firstname: "Jane", lastname: "Active", state: domain.Active},
			{email: "unverified@example.com", firstname: "Bob", lastname: "Unverified", state: domain.Unverified},
		}

		// Insert test users with different states
		for _, baseUser := range baseUsers {
			user := td.NewTestUser(t, baseUser.email, baseUser.firstname, baseUser.lastname)
			user.State = baseUser.state

			userEncx, err := domain.ProcessUserEncx(ctx, crypto, user)
			require.NoError(t, err)

			err = td.InsertUserEncx(t, ctx, userEncx, testPool, crypto)
			require.NoError(t, err)
		}

		// Act
		req := td.NewGetAllUsersRequest(t, ctx, testServerURL, accessToken)
		resp, err := client.Do(req)

		// Assert
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		users := td.ParseGetAllUsersResponse(t, resp)
		assert.Len(t, users, len(baseUsers)+1, "Should return all users regardless of state (3 test users + 1 admin)")

		// Verify user data (should be decrypted in response)
		emails := make(map[string]struct{}, len(users))
		roles := make(map[string]struct{}, len(users))
		states := make(map[domain.UserState]struct{}, len(users))
		for _, u := range users {
			emails[u.Email] = struct{}{}
			roles[u.Role] = struct{}{}
			states[u.State] = struct{}{}
		}

		// Verify every base user email is present
		for _, bu := range baseUsers {
			_, exists := emails[bu.email]
			assert.Truef(t, exists, "expected email %s to be in users", bu.email)
		}
		// Verify at least one administrator role is present
		_, hasAdmin := roles[identity.AdministratorStr]
		assert.True(t, hasAdmin, "expected at least one Administrator user")

		// Verify all required states are present
		expectedStates := []domain.UserState{domain.Pending, domain.Active, domain.Unverified}
		for _, s := range expectedStates {
			_, exists := states[s]
			assert.Truef(t, exists, "expected state %s to be present among users", s)
		}
	})

	t.Run("should return users ordered by creation date descending", func(t *testing.T) {
		// Clean state
		td.ClearUsersTable(t, ctx, testPool)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		const userCount = 3

		baseUsers := [userCount]baseUser{
			{email: "first@example.com", firstname: "First", lastname: "User", state: domain.Active},
			{email: "second@example.com", firstname: "Second", lastname: "User", state: domain.Pending},
			{email: "third@example.com", firstname: "Third", lastname: "User", state: domain.Unverified},
		}

		// Insert test users with different states
		for _, baseUser := range baseUsers {
			user := td.NewTestUser(t, baseUser.email, baseUser.firstname, baseUser.lastname)
			user.State = baseUser.state

			userEncx, err := domain.ProcessUserEncx(ctx, crypto, user)
			require.NoError(t, err)

			err = td.InsertUserEncx(t, ctx, userEncx, testPool, crypto)
			require.NoError(t, err)

			time.Sleep(10 * time.Millisecond)
		}

		// Act
		req := td.NewGetAllUsersRequest(t, ctx, testServerURL, accessToken)
		resp, err := client.Do(req)

		// Assert
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		users := td.ParseGetAllUsersResponse(t, resp)
		require.Len(t, users, userCount+1, "Should return all users regardless of state (3 test users + 1 admin")

		// Should be ordered by creation date descending (newest first)
		// Admin user is created in setupAdminUser (first), so will be last
		// Test users created in order: first, second, third (newest)
		for i, bu := range baseUsers {
			assert.Equal(t, bu.email, users[userCount-1-i].Email)
		}
		assert.Equal(t, identity.AdministratorStr, users[userCount].Role)
	})

	t.Run("should properly decrypt and return all user fields", func(t *testing.T) {
		// Clean state
		td.ClearUsersTable(t, ctx, testPool)
		layout := "2006-01-02"
		date := "1995-05-12"
		expectedBirthDate, err := time.ParseInLocation(layout, date, time.Local)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		email := "complete@example.com"
		firstname := "John"
		lastname := "DOE"
		state := domain.Active
		telephone := "1234567890"
		picture := "profile.jpg"
		gender := domain.GenderMan.String()
		address1 := "123 Main St"
		city := "New York"
		postalcode := "10001"

		// Insert test user with all fields populated
		testUser := td.NewTestUser(t, email, firstname, lastname)
		testUser.State = state
		testUser.Telephone = telephone
		testUser.Picture = picture
		testUser.BirthDate = expectedBirthDate
		testUser.Gender = gender
		testUser.Address1 = address1
		testUser.City = city
		testUser.PostalCode = postalcode

		testUserEncx, err := domain.ProcessUserEncx(ctx, crypto, testUser)
		require.NoError(t, err)

		err = td.InsertUserEncx(t, ctx, testUserEncx, testPool, crypto)
		require.NoError(t, err)

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
		assert.Equal(t, email, user.Email)
		assert.Equal(t, firstname, user.FirstName)
		assert.Equal(t, lastname, user.LastName)
		assert.Equal(t, telephone, user.Telephone)
		assert.Equal(t, picture, user.Picture)
		assert.Equal(t, expectedBirthDate, user.BirthDate)
		assert.Equal(t, gender, user.Gender)
		assert.Equal(t, address1, user.Address1)
		assert.Equal(t, city, user.City)
		assert.Equal(t, postalcode, user.PostalCode)
		assert.Equal(t, state, user.State)
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

		minimalUserEncx, err := domain.ProcessUserEncx(ctx, crypto, minimalUser)
		require.NoError(t, err)
		// Leave optional fields empty

		td.InsertUserEncx(t, ctx, minimalUserEncx, testPool, crypto)
		require.NoError(t, err)

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
		assert.Equal(t, email, user.Email)
		assert.Equal(t, firstname, user.FirstName)
		assert.Equal(t, lastname, user.LastName)
		assert.Empty(t, user.Picture)
		assert.Empty(t, user.Gender)
		assert.Empty(t, user.Address1)
		assert.Empty(t, user.City)
		assert.Empty(t, user.PostalCode)
	})

	t.Run("should include all user states in mixed scenario", func(t *testing.T) {
		// Clean state
		td.ClearUsersTable(t, ctx, testPool)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		email := "mixed%d@example.com"
		firstname := "User"
		lastname := "%d"

		// Create states with all possible states
		states := [5]*domain.User{
			{State: domain.Pending},
			{State: domain.Active},
			{State: domain.Unverified},
			{State: domain.Pending}, // Another pending user
			{State: domain.Active},  // Another active user
		}

		// Insert users
		for i, user := range states {
			testUser := td.NewTestUser(t, fmt.Sprintf(email, i), firstname, fmt.Sprintf(lastname, i))
			testUser.State = user.State
			testUserEncx, err := domain.ProcessUserEncx(ctx, crypto, testUser)
			require.NoError(t, err)
			err = td.InsertUserEncx(t, ctx, testUserEncx, testPool, crypto)
			require.NoError(t, err)
		}

		// Act
		req := td.NewGetAllUsersRequest(t, ctx, testServerURL, accessToken)
		resp, err := client.Do(req)

		// Assert
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		responseUsers := td.ParseGetAllUsersResponse(t, resp)
		require.Len(t, responseUsers, len(states)+1, "Should return all users regardless of state (5 test users + 1 admin")

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

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		const userCount = 5

		states := [3]domain.UserState{domain.Pending, domain.Active, domain.Unverified}

		// Insert multiple users with different states
		for i := range userCount {
			user := td.NewTestUser(t, fmt.Sprintf("concurrent%d@example.com", i), "User", fmt.Sprintf("%d", i))
			user.State = states[i%len(states)]

			userEncx, err := domain.ProcessUserEncx(ctx, crypto, user)
			require.NoError(t, err)

			err = td.InsertUserEncx(t, ctx, userEncx, testPool, crypto)
			require.NoError(t, err)
		}

		const concurrentRequestCount = 3

		// Make concurrent requests
		responses := make(chan *http.Response, concurrentRequestCount)
		for range concurrentRequestCount {
			go func() {
				req := td.NewGetAllUsersRequest(t, ctx, testServerURL, accessToken)

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
				users := td.ParseGetAllUsersResponse(t, resp)
				assert.Len(t, users, userCount+1, "Each concurrent request should return all users (5 test + 1 admin)")
			}
		}

		assert.Equal(t, concurrentRequestCount, successCount, "All concurrent requests should succeed")
	})

	t.Run("should handle large number of users efficiently", func(t *testing.T) {
		// Clean state
		td.ClearUsersTable(t, ctx, testPool)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		// Insert many users with various states
		const userCount = 100
		states := [3]domain.UserState{domain.Pending, domain.Active, domain.Unverified}

		for i := range userCount {
			user := td.NewTestUser(t, fmt.Sprintf("bulk%d@example.com", i), "Bulk", fmt.Sprintf("User%d", i))
			user.State = states[i%len(states)] // Cycle through states

			userEncx, err := domain.ProcessUserEncx(ctx, crypto, user)
			require.NoError(t, err)

			err = td.InsertUserEncx(t, ctx, userEncx, testPool, crypto)
			require.NoError(t, err)
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
		testUser := td.NewTestUser(t, "test@example.com", "Test", "User")
		testUser.State = domain.Active
		testUserEncx, err := domain.ProcessUserEncx(ctx, crypto, testUser)
		require.NoError(t, err)

		err = td.InsertUserEncx(t, ctx, testUserEncx, testPool, crypto)
		require.NoError(t, err)

		// Act without admin auth (this should fail if auth middleware is properly configured)
		req := td.NewGetAllUsersRequestWithoutAuth(t, ctx, testServerURL)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, resp.StatusCode, http.StatusUnauthorized)
	})
}

type baseUser struct {
	email     string
	firstname string
	lastname  string
	state     domain.UserState
}
