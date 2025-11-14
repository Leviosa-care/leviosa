package partner_test

import (
	"context"
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

// make test-func TEST_NAME=TestVerifyPartner TEST_PATH=test/integration/authuser/partner/verify_partner_test.go

func TestVerifyPartner(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("should successfully verify unverified partner", func(t *testing.T) {
		// Clean state
		td.ClearPartnersTable(t, ctx, testPool)
		td.ClearSessionsRedis(t, ctx, redisClient)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		// Create a test user with partner role (unverified)
		testUser := td.NewTestUser(t, "partner1@example.com", "John", "Partner")
		testUser.State = domain.Pending // Partner awaiting verification
		testUserEncx, err := domain.ProcessUserEncx(ctx, crypto, testUser)
		require.NoError(t, err)
		err = td.InsertUserEncx(t, ctx, testUserEncx, testPool)
		require.NoError(t, err)

		// Create an unverified partner profile
		testPartner := &domain.Partner{
			ID:          uuid.New(),
			UserID:      testUser.ID,
			Bio:         "Experienced healthcare professional",
			Experience:  "10 years in home care",
			CategoryIDs: []uuid.UUID{},
			ProductIDs:  []uuid.UUID{},
		}
		partnerEncx, err := domain.ProcessPartnerEncx(ctx, crypto, testPartner)
		require.NoError(t, err)
		err = td.InsertPartnerEncx(t, ctx, partnerEncx, testPool)
		require.NoError(t, err)

		// Act - verify the partner
		req := td.NewVerifyPartnerRequest(t, ctx, testServerURL, testPartner.ID, accessToken)

		resp, err := client.Do(req)

		// Assert HTTP response
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Verify user state and role were updated
		userEncx, err := td.GetUserEnxByID(t, ctx, testUser.ID, testPool)
		require.NoError(t, err)
		user, err := domain.DecryptUserEncx(ctx, crypto, userEncx)
		require.NoError(t, err)
		assert.Equal(t, domain.Active, user.State, "User state should be active")
		assert.Equal(t, identity.PartnerStr, user.Role, "User role should be partner")
	})

	t.Run("should return error when partner not found", func(t *testing.T) {
		// Clean state
		td.ClearPartnersTable(t, ctx, testPool)
		td.ClearSessionsRedis(t, ctx, redisClient)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		// Act - try to verify non-existent partner
		nonExistentID := uuid.New()
		req := td.NewVerifyPartnerRequest(t, ctx, testServerURL, nonExistentID, accessToken)

		resp, err := client.Do(req)

		// Assert
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("should return error when partner is already verified", func(t *testing.T) {
		// Clean state
		td.ClearPartnersTable(t, ctx, testPool)
		td.ClearSessionsRedis(t, ctx, redisClient)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		// Create a test user
		testUser := td.NewTestUser(t, "partner2@example.com", "Jane", "Verified")
		testUser.State = domain.Active
		testUser.Role = identity.PartnerStr
		testUserEncx, err := domain.ProcessUserEncx(ctx, crypto, testUser)
		require.NoError(t, err)
		err = td.InsertUserEncx(t, ctx, testUserEncx, testPool)
		require.NoError(t, err)

		// Create an already verified partner (user state is Active and Stripe status is active)
		testPartner := &domain.Partner{
			ID:                       uuid.New(),
			UserID:                   testUser.ID,
			Bio:                      "Already verified",
			Experience:               "5 years",
			CategoryIDs:              []uuid.UUID{},
			ProductIDs:               []uuid.UUID{},
			StripeAccountStatus:      domain.StripeAccountStatusActive,
			StripeOnboardingComplete: true,
		}
		partnerEncx, err := domain.ProcessPartnerEncx(ctx, crypto, testPartner)
		require.NoError(t, err)
		err = td.InsertPartnerEncx(t, ctx, partnerEncx, testPool)
		require.NoError(t, err)

		// Act - try to verify already verified partner
		req := td.NewVerifyPartnerRequest(t, ctx, testServerURL, testPartner.ID, accessToken)

		resp, err := client.Do(req)

		// Assert
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusConflict, resp.StatusCode)
	})

	t.Run("should require admin role", func(t *testing.T) {
		// Clean state
		td.ClearPartnersTable(t, ctx, testPool)
		td.ClearSessionsRedis(t, ctx, redisClient)

		// Create standard user session (not admin)
		accessToken := tu.SetupStandardUser(t, ctx, authCtx)

		// Create a partner to verify
		partnerUser := td.NewTestUser(t, "partner3@example.com", "Partner", "User")
		partnerUser.State = domain.Pending
		partnerUserEncx, err := domain.ProcessUserEncx(ctx, crypto, partnerUser)
		require.NoError(t, err)
		err = td.InsertUserEncx(t, ctx, partnerUserEncx, testPool)
		require.NoError(t, err)

		testPartner := &domain.Partner{
			ID:          uuid.New(),
			UserID:      partnerUser.ID,
			Bio:         "Test bio",
			Experience:  "Test experience",
			CategoryIDs: []uuid.UUID{},
			ProductIDs:  []uuid.UUID{},
		}
		partnerEncx, err := domain.ProcessPartnerEncx(ctx, crypto, testPartner)
		require.NoError(t, err)
		err = td.InsertPartnerEncx(t, ctx, partnerEncx, testPool)
		require.NoError(t, err)

		// Act - try to verify with non-admin role
		req := td.NewVerifyPartnerRequest(t, ctx, testServerURL, testPartner.ID, accessToken)

		resp, err := client.Do(req)

		// Assert
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})

	t.Run("should require authentication", func(t *testing.T) {
		// Clean state
		td.ClearPartnersTable(t, ctx, testPool)

		// Create a partner to verify
		testUser := td.NewTestUser(t, "partner4@example.com", "Partner", "User")
		testUser.State = domain.Pending
		testUserEncx, err := domain.ProcessUserEncx(ctx, crypto, testUser)
		require.NoError(t, err)
		err = td.InsertUserEncx(t, ctx, testUserEncx, testPool)
		require.NoError(t, err)

		testPartner := &domain.Partner{
			ID:          uuid.New(),
			UserID:      testUser.ID,
			Bio:         "Test bio",
			Experience:  "Test experience",
			CategoryIDs: []uuid.UUID{},
			ProductIDs:  []uuid.UUID{},
		}
		partnerEncx, err := domain.ProcessPartnerEncx(ctx, crypto, testPartner)
		require.NoError(t, err)
		err = td.InsertPartnerEncx(t, ctx, partnerEncx, testPool)
		require.NoError(t, err)

		// Act - try to verify without authentication
		req := td.NewVerifyPartnerRequest(t, ctx, testServerURL, testPartner.ID, "")

		resp, err := client.Do(req)

		// Assert
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("should return error for invalid partner ID format", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, redisClient)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		// Act - try with invalid UUID format
		req, err := http.NewRequestWithContext(
			ctx,
			http.MethodPost,
			testServerURL+"/admin/partners/invalid-uuid/verify",
			nil,
		)
		require.NoError(t, err)

		// Add session cookie
		if accessToken != "" {
			cookie := &http.Cookie{
				Name:  ck.AccessTokenCookieName,
				Value: accessToken,
			}
			req.AddCookie(cookie)
		}

		resp, err := client.Do(req)

		// Assert
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}
