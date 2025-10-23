package partner_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	td "github.com/Leviosa-care/leviosa/backend/test/helpers"
	"github.com/Leviosa-care/leviosa/backend/internal/common/contracts/identity"
	tu "github.com/Leviosa-care/leviosa/backend/internal/common/testutils"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TEST=TestGetAllPartners make test-integration-partner-test

func TestGetAllPartners(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("should successfully get all partners with admin role", func(t *testing.T) {
		// Clean state
		td.ClearPartnersTable(t, ctx, testPool)
		td.ClearSessionsRedis(t, ctx, redisClient)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		// Create multiple test partners
		for i := 0; i < 3; i++ {
			testUser := td.NewTestUser(t,
				"partner"+string(rune(i))+"@example.com",
				"Partner",
				string(rune('A'+i)))
			testUser.State = domain.Active
			testUser.Role = identity.PartnerStr
			testUserEncx, err := domain.ProcessUserEncx(ctx, crypto, testUser)
			require.NoError(t, err)
			err = td.InsertUserEncx(t, ctx, testUserEncx, testPool, crypto)
			require.NoError(t, err)

			testPartner := &domain.Partner{
				ID:             uuid.New(),
				UserID:         testUser.ID,
				Bio:            "Bio for partner " + string(rune('A'+i)),
				Experience:     "Experience " + string(rune('A'+i)),
				Certifications: []string{"Cert " + string(rune('A'+i))},
				CategoryIDs:    []uuid.UUID{},
				ProductIDs:     []uuid.UUID{},
				IsVerified:     i%2 == 0, // Alternate verified status
			}
			td.InsertPartner(t, ctx, testPartner, testPool, crypto)
		}

		// Act
		req := td.NewGetAllPartnersRequest(t, ctx, testServerURL)
		req.Header.Set("Cookie", td.ToCookieString(accessToken))

		resp, err := client.Do(req)

		// Assert
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// TODO: Parse response and verify all partners are returned with user info
	})

	t.Run("should return empty array when no partners exist", func(t *testing.T) {
		// Clean state
		td.ClearPartnersTable(t, ctx, testPool)
		td.ClearSessionsRedis(t, ctx, redisClient)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		// Act
		req := td.NewGetAllPartnersRequest(t, ctx, testServerURL)
		req.Header.Set("Cookie", td.ToCookieString(accessToken))

		resp, err := client.Do(req)

		// Assert
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// TODO: Parse response and verify empty partners array
	})

	t.Run("should return 403 for non-admin user", func(t *testing.T) {
		// Clean state
		td.ClearPartnersTable(t, ctx, testPool)
		td.ClearSessionsRedis(t, ctx, redisClient)

		// Create standard user
		standardUser := td.NewTestUser(t, "standard@example.com", "Standard", "User")
		standardUser.State = domain.Active
		standardUserEncx, err := domain.ProcessUserEncx(ctx, crypto, standardUser)
		require.NoError(t, err)
		err = td.InsertUserEncx(t, ctx, standardUserEncx, testPool, crypto)
		require.NoError(t, err)

		accessToken := tu.SetupSessionForUser(t, ctx, authCtx, standardUser.ID, identity.Standard)

		// Act
		req := td.NewGetAllPartnersRequest(t, ctx, testServerURL)
		req.Header.Set("Cookie", td.ToCookieString(accessToken))

		resp, err := client.Do(req)

		// Assert
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})

	t.Run("should return 401 without authentication", func(t *testing.T) {
		// Act
		req := td.NewGetAllPartnersRequest(t, ctx, testServerURL)
		// No session cookie

		resp, err := client.Do(req)

		// Assert
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("should include user information for each partner", func(t *testing.T) {
		// Clean state
		td.ClearPartnersTable(t, ctx, testPool)
		td.ClearSessionsRedis(t, ctx, redisClient)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		// Create test partner
		testUser := td.NewTestUser(t, "partner@example.com", "John", "Partner")
		testUser.State = domain.Active
		testUser.Role = identity.PartnerStr
		testUserEncx, err := domain.ProcessUserEncx(ctx, crypto, testUser)
		require.NoError(t, err)
		err = td.InsertUserEncx(t, ctx, testUserEncx, testPool, crypto)
		require.NoError(t, err)

		testPartner := &domain.Partner{
			ID:             uuid.New(),
			UserID:         testUser.ID,
			Bio:            "Test bio",
			Experience:     "Test experience",
			Certifications: []string{"Cert1"},
			CategoryIDs:    []uuid.UUID{},
			ProductIDs:     []uuid.UUID{},
			IsVerified:     true,
		}
		td.InsertPartner(t, ctx, testPartner, testPool, crypto)

		// Act
		req := td.NewGetAllPartnersRequest(t, ctx, testServerURL)
		req.Header.Set("Cookie", td.ToCookieString(accessToken))

		resp, err := client.Do(req)

		// Assert
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// TODO: Parse response and verify partner includes complete user information
	})
}
