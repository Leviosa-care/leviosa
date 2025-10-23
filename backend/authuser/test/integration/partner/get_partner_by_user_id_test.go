package partner_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/Leviosa-care/authuser/internal/domain"
	td "github.com/Leviosa-care/authuser/test/helpers"
	"github.com/Leviosa-care/core/contracts/identity"
	tu "github.com/Leviosa-care/core/testutils"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TEST=TestGetPartnerByUserID make test-integration-partner-test

func TestGetPartnerByUserID(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("should successfully get own partner profile", func(t *testing.T) {
		// Clean state
		td.ClearPartnersTable(t, ctx, testPool)
		td.ClearSessionsRedis(t, ctx, redisClient)

		// Create partner user
		partnerUser := td.NewTestUser(t, "partner@example.com", "John", "Partner")
		partnerUser.State = domain.Active
		partnerUser.Role = identity.PartnerStr
		partnerUserEncx, err := domain.ProcessUserEncx(ctx, crypto, partnerUser)
		require.NoError(t, err)
		err = td.InsertUserEncx(t, ctx, partnerUserEncx, testPool, crypto)
		require.NoError(t, err)

		// Create partner profile
		testPartner := &domain.Partner{
			ID:             uuid.New(),
			UserID:         partnerUser.ID,
			Bio:            "My bio",
			Experience:     "My experience",
			Certifications: []string{"My Cert"},
			CategoryIDs:    []uuid.UUID{},
			ProductIDs:     []uuid.UUID{},
			IsVerified:     true,
		}
		td.InsertPartner(t, ctx, testPartner, testPool, crypto)

		// Create session for partner user
		accessToken := tu.SetupSessionForUser(t, ctx, authCtx, partnerUser.ID, identity.Partner)

		// Act - get own partner profile
		req := td.NewGetPartnerByUserIDRequest(t, ctx, testServerURL, partnerUser.ID)
		req.Header.Set("Cookie", td.ToCookieString(accessToken))

		resp, err := client.Do(req)

		// Assert
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// TODO: Parse response and verify partner data matches
	})

	t.Run("should successfully get any partner profile with admin role", func(t *testing.T) {
		// Clean state
		td.ClearPartnersTable(t, ctx, testPool)
		td.ClearSessionsRedis(t, ctx, redisClient)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		// Create partner user
		partnerUser := td.NewTestUser(t, "partner@example.com", "John", "Partner")
		partnerUser.State = domain.Active
		partnerUser.Role = identity.PartnerStr
		partnerUserEncx, err := domain.ProcessUserEncx(ctx, crypto, partnerUser)
		require.NoError(t, err)
		err = td.InsertUserEncx(t, ctx, partnerUserEncx, testPool, crypto)
		require.NoError(t, err)

		testPartner := &domain.Partner{
			ID:             uuid.New(),
			UserID:         partnerUser.ID,
			Bio:            "Test bio",
			Experience:     "Test experience",
			Certifications: []string{"Test Cert"},
			CategoryIDs:    []uuid.UUID{},
			ProductIDs:     []uuid.UUID{},
			IsVerified:     true,
		}
		td.InsertPartner(t, ctx, testPartner, testPool, crypto)

		// Act - admin gets partner profile
		req := td.NewGetPartnerByUserIDRequest(t, ctx, testServerURL, partnerUser.ID)
		req.Header.Set("Cookie", td.ToCookieString(accessToken))

		resp, err := client.Do(req)

		// Assert
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("should return 403 when partner tries to access another partner's profile", func(t *testing.T) {
		// Clean state
		td.ClearPartnersTable(t, ctx, testPool)
		td.ClearSessionsRedis(t, ctx, redisClient)

		// Create first partner (the one making the request)
		partner1 := td.NewTestUser(t, "partner1@example.com", "Partner", "One")
		partner1.State = domain.Active
		partner1.Role = identity.PartnerStr
		partner1Encx, err := domain.ProcessUserEncx(ctx, crypto, partner1)
		require.NoError(t, err)
		err = td.InsertUserEncx(t, ctx, partner1Encx, testPool, crypto)
		require.NoError(t, err)

		testPartner1 := &domain.Partner{
			ID:      uuid.New(),
			UserID:  partner1.ID,
			Bio:     "Partner 1 bio",
			IsVerified: true,
		}
		td.InsertPartner(t, ctx, testPartner1, testPool, crypto)

		accessToken := tu.SetupSessionForUser(t, ctx, authCtx, partner1.ID, identity.Partner)

		// Create second partner (target)
		partner2 := td.NewTestUser(t, "partner2@example.com", "Partner", "Two")
		partner2.State = domain.Active
		partner2.Role = identity.PartnerStr
		partner2Encx, err := domain.ProcessUserEncx(ctx, crypto, partner2)
		require.NoError(t, err)
		err = td.InsertUserEncx(t, ctx, partner2Encx, testPool, crypto)
		require.NoError(t, err)

		testPartner2 := &domain.Partner{
			ID:      uuid.New(),
			UserID:  partner2.ID,
			Bio:     "Partner 2 bio",
			IsVerified: true,
		}
		td.InsertPartner(t, ctx, testPartner2, testPool, crypto)

		// Act - partner1 tries to access partner2's profile
		req := td.NewGetPartnerByUserIDRequest(t, ctx, testServerURL, partner2.ID)
		req.Header.Set("Cookie", td.ToCookieString(accessToken))

		resp, err := client.Do(req)

		// Assert
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})

	t.Run("should return 404 when partner profile doesn't exist", func(t *testing.T) {
		// Clean state
		td.ClearPartnersTable(t, ctx, testPool)
		td.ClearSessionsRedis(t, ctx, redisClient)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		// Create user without partner profile
		userWithoutPartner := td.NewTestUser(t, "user@example.com", "Regular", "User")
		userWithoutPartner.State = domain.Active
		userEncx, err := domain.ProcessUserEncx(ctx, crypto, userWithoutPartner)
		require.NoError(t, err)
		err = td.InsertUserEncx(t, ctx, userEncx, testPool, crypto)
		require.NoError(t, err)

		// Act
		req := td.NewGetPartnerByUserIDRequest(t, ctx, testServerURL, userWithoutPartner.ID)
		req.Header.Set("Cookie", td.ToCookieString(accessToken))

		resp, err := client.Do(req)

		// Assert
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("should return 401 without authentication", func(t *testing.T) {
		// Clean state
		td.ClearPartnersTable(t, ctx, testPool)

		// Create partner
		partnerUser := td.NewTestUser(t, "partner@example.com", "Partner", "User")
		partnerUser.State = domain.Active
		partnerUserEncx, err := domain.ProcessUserEncx(ctx, crypto, partnerUser)
		require.NoError(t, err)
		err = td.InsertUserEncx(t, ctx, partnerUserEncx, testPool, crypto)
		require.NoError(t, err)

		testPartner := &domain.Partner{
			ID:      uuid.New(),
			UserID:  partnerUser.ID,
			Bio:     "Test bio",
			IsVerified: true,
		}
		td.InsertPartner(t, ctx, testPartner, testPool, crypto)

		// Act
		req := td.NewGetPartnerByUserIDRequest(t, ctx, testServerURL, partnerUser.ID)
		// No session cookie

		resp, err := client.Do(req)

		// Assert
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}
