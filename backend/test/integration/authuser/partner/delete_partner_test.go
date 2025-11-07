package partner_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	ck "github.com/Leviosa-care/leviosa/backend/internal/common/auth/cookies"
	"github.com/Leviosa-care/leviosa/backend/internal/common/contracts/identity"
	tu "github.com/Leviosa-care/leviosa/backend/internal/common/testutils"
	td "github.com/Leviosa-care/leviosa/backend/test/helpers"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestDeletePartner TEST_PATH=test/integration/authuser/partner/delete_partner_test.go

func TestDeletePartner(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("should successfully delete partner but preserve user account", func(t *testing.T) {
		// Clean state
		td.ClearPartnersTable(t, ctx, testPool)
		td.ClearSessionsRedis(t, ctx, redisClient)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		// Create test user and partner
		testUser := td.NewTestUser(t, "partner@example.com", "John", "Partner")
		testUser.State = domain.Active
		testUser.Role = identity.PartnerStr
		testUserEncx, err := domain.ProcessUserEncx(ctx, crypto, testUser)
		require.NoError(t, err)
		err = td.InsertUserEncx(t, ctx, testUserEncx, testPool)
		require.NoError(t, err)

		testPartner := td.NewTestPartner(t, testUser.ID)
		partnerEncx, err := domain.ProcessPartnerEncx(ctx, crypto, testPartner)
		require.NoError(t, err)
		err = td.InsertPartnerEncx(t, ctx, partnerEncx, testPool)
		require.NoError(t, err)

		// Act
		req := td.NewDeletePartnerRequest(t, ctx, testServerURL, testPartner.ID, accessToken)

		resp, err := client.Do(req)

		// Assert HTTP response
		assert.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)

		// Verify partner was deleted
		_, err = td.GetPartnerEncxByUserID(t, ctx, testUser.ID, testPool)
		assert.Error(t, err, "Partner should be deleted")

		// Verify user still exists
		userEncx, err := td.GetUserEnxByID(t, ctx, testUser.ID, testPool)
		assert.NoError(t, err, "User should still exist after partner deletion")
		user, err := domain.DecryptUserEncx(ctx, crypto, userEncx)
		assert.NoError(t, err)
		assert.Equal(t, testUser.ID, user.ID)
		assert.Equal(t, testUser.Email, user.Email)
	})

	t.Run("should return 404 when partner not found", func(t *testing.T) {
		// Clean state
		td.ClearPartnersTable(t, ctx, testPool)
		td.ClearSessionsRedis(t, ctx, redisClient)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		// Act - try to delete non-existent partner
		nonExistentID := uuid.New()
		req := td.NewDeletePartnerRequest(t, ctx, testServerURL, nonExistentID, accessToken)

		resp, err := client.Do(req)

		// Assert
		assert.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("should return 403 for non-admin user", func(t *testing.T) {
		// Clean state
		td.ClearPartnersTable(t, ctx, testPool)
		td.ClearSessionsRedis(t, ctx, redisClient)

		// Setup standard user session
		accessToken := tu.SetupStandardUser(t, ctx, authCtx)

		// Create partner to delete
		partnerUser := td.NewTestUser(t, "partner@example.com", "Partner", "User")
		partnerUser.State = domain.Active
		partnerUser.Role = identity.PartnerStr
		partnerUserEncx, err := domain.ProcessUserEncx(ctx, crypto, partnerUser)
		require.NoError(t, err)
		err = td.InsertUserEncx(t, ctx, partnerUserEncx, testPool)
		require.NoError(t, err)

		testPartner := td.NewTestPartner(t, partnerUser.ID)
		partnerEncx, err := domain.ProcessPartnerEncx(ctx, crypto, testPartner)
		require.NoError(t, err)
		err = td.InsertPartnerEncx(t, ctx, partnerEncx, testPool)
		require.NoError(t, err)

		// Act - standard user tries to delete partner
		req := td.NewDeletePartnerRequest(t, ctx, testServerURL, testPartner.ID, accessToken)

		resp, err := client.Do(req)

		// Assert
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)

		// Verify partner still exists
		exists, err := td.CheckPartnerExistsByUserID(t, ctx, partnerUser.ID, testPool)
		assert.NoError(t, err)
		assert.True(t, exists, "Partner should still exist after failed delete")
	})

	t.Run("should return 401 without authentication", func(t *testing.T) {
		// Clean state
		td.ClearPartnersTable(t, ctx, testPool)

		// Create partner to delete
		testUser := td.NewTestUser(t, "partner@example.com", "Partner", "User")
		testUser.State = domain.Active
		testUser.Role = identity.PartnerStr
		testUserEncx, err := domain.ProcessUserEncx(ctx, crypto, testUser)
		require.NoError(t, err)
		err = td.InsertUserEncx(t, ctx, testUserEncx, testPool)
		require.NoError(t, err)

		testPartner := td.NewTestPartner(t, testUser.ID)
		partnerEncx, err := domain.ProcessPartnerEncx(ctx, crypto, testPartner)
		require.NoError(t, err)
		err = td.InsertPartnerEncx(t, ctx, partnerEncx, testPool)
		require.NoError(t, err)

		// Act - no session cookie
		req := td.NewDeletePartnerRequest(t, ctx, testServerURL, testPartner.ID, "")

		resp, err := client.Do(req)

		// Assert
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

		// Verify partner still exists
		exists, err := td.CheckPartnerExistsByUserID(t, ctx, testUser.ID, testPool)
		assert.NoError(t, err)
		assert.True(t, exists, "Partner should still exist after unauthorized delete attempt")
	})

	t.Run("should return 400 for invalid UUID format", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, redisClient)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		// Act - invalid UUID in path
		req, err := http.NewRequestWithContext(
			ctx,
			http.MethodDelete,
			testServerURL+"/admin/partners/invalid-uuid",
			nil,
		)
		require.NoError(t, err)

		// Add access token cookie properly
		cookie := &http.Cookie{
			Name:  ck.AccessTokenCookieName,
			Value: accessToken,
		}
		req.AddCookie(cookie)

		resp, err := client.Do(req)

		// Assert
		assert.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should delete partner when multiple partners exist", func(t *testing.T) {
		// Clean state
		td.ClearPartnersTable(t, ctx, testPool)
		td.ClearSessionsRedis(t, ctx, redisClient)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		// Create 3 partners
		partners := make([]*domain.Partner, 0, 3)
		for i := 0; i < 3; i++ {
			user := td.NewTestUser(t,
				fmt.Sprintf("partner%d@example.com", i),
				"Partner",
				fmt.Sprintf("User%d", i))
			user.State = domain.Active
			user.Role = identity.PartnerStr
			userEncx, err := domain.ProcessUserEncx(ctx, crypto, user)
			require.NoError(t, err)
			err = td.InsertUserEncx(t, ctx, userEncx, testPool)
			require.NoError(t, err)

			partner := td.NewTestPartner(t, user.ID)
			partnerEncx, err := domain.ProcessPartnerEncx(ctx, crypto, partner)
			require.NoError(t, err)
			err = td.InsertPartnerEncx(t, ctx, partnerEncx, testPool)
			require.NoError(t, err)

			partners = append(partners, partner)
		}

		// Verify initial count
		initialCount, err := td.CountPartners(t, ctx, testPool)
		require.NoError(t, err)
		assert.Equal(t, 3, initialCount)

		// Act - delete the second partner
		req := td.NewDeletePartnerRequest(t, ctx, testServerURL, partners[1].ID, accessToken)

		resp, err := client.Do(req)

		// Assert
		assert.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)

		// Verify count decreased by 1
		finalCount, err := td.CountPartners(t, ctx, testPool)
		assert.NoError(t, err)
		assert.Equal(t, 2, finalCount)

		// Verify deleted partner doesn't exist
		_, err = td.GetPartnerEncxByID(t, ctx, partners[1].ID, testPool)
		assert.Error(t, err)

		// Verify other partners still exist
		_, err = td.GetPartnerEncxByID(t, ctx, partners[0].ID, testPool)
		assert.NoError(t, err)
		_, err = td.GetPartnerEncxByID(t, ctx, partners[2].ID, testPool)
		assert.NoError(t, err)
	})

	t.Run("should handle partner with encrypted data correctly", func(t *testing.T) {
		// Clean state
		td.ClearPartnersTable(t, ctx, testPool)
		td.ClearSessionsRedis(t, ctx, redisClient)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		// Create user and partner with full encrypted data
		testUser := td.NewTestUser(t, "partner@example.com", "John", "Doe")
		testUser.State = domain.Active
		testUser.Role = identity.PartnerStr
		userEncx, err := domain.ProcessUserEncx(ctx, crypto, testUser)
		require.NoError(t, err)
		err = td.InsertUserEncx(t, ctx, userEncx, testPool)
		require.NoError(t, err)

		// Partner with full data including Stripe info
		partner := td.NewTestPartner(t, testUser.ID)
		partner.Bio = "Experienced professional with multiple certifications"
		partner.Experience = "10+ years in healthcare"
		partner.StripeConnectedAccountID = "acct_1234567890"
		partner.StripeAccountStatus = domain.StripeAccountStatusActive
		partner.StripeOnboardingComplete = true
		partner.CategoryIDs = []uuid.UUID{uuid.New(), uuid.New()}
		partner.ProductIDs = []uuid.UUID{uuid.New(), uuid.New(), uuid.New()}

		partnerEncx, err := domain.ProcessPartnerEncx(ctx, crypto, partner)
		require.NoError(t, err)
		err = td.InsertPartnerEncx(t, ctx, partnerEncx, testPool)
		require.NoError(t, err)

		// Act - delete partner
		req := td.NewDeletePartnerRequest(t, ctx, testServerURL, partner.ID, accessToken)

		resp, err := client.Do(req)

		// Assert
		assert.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)

		// Verify partner was deleted
		_, err = td.GetPartnerEncxByID(t, ctx, partner.ID, testPool)
		assert.Error(t, err, "Partner with encrypted data should be deleted")

		// Verify user still exists with all data intact
		retrievedUserEncx, err := td.GetUserEnxByID(t, ctx, testUser.ID, testPool)
		assert.NoError(t, err)
		retrievedUser, err := domain.DecryptUserEncx(ctx, crypto, retrievedUserEncx)
		assert.NoError(t, err)
		assert.Equal(t, testUser.Email, retrievedUser.Email)
		assert.Equal(t, testUser.FirstName, retrievedUser.FirstName)
		assert.Equal(t, testUser.LastName, retrievedUser.LastName)
	})
}
