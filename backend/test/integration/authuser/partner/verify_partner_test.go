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

// TEST=TestVerifyPartner make test-integration-partner-test

func TestVerifyPartner(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("should successfully verify unverified partner", func(t *testing.T) {
		// Clean state
		td.ClearPartnersTable(t, ctx, testPool)
		td.ClearSessionsRedis(t, ctx, redisClient)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		// Create a test user with partner role (unverified)
		testUser := td.NewTestUser(t, "partner@example.com", "John", "Partner")
		testUser.State = domain.Pending // Partner awaiting verification
		testUserEncx, err := domain.ProcessUserEncx(ctx, crypto, testUser)
		require.NoError(t, err)
		err = td.InsertUserEncx(t, ctx, testUserEncx, testPool, crypto)
		require.NoError(t, err)

		// Create an unverified partner profile
		testPartner := &domain.Partner{
			ID:             uuid.New(),
			UserID:         testUser.ID,
			Bio:            "Experienced healthcare professional",
			Experience:     "10 years in home care",
			// Certifications: []string{"CPR Certified", "First Aid"},
			CategoryIDs:    []uuid.UUID{},
			ProductIDs:     []uuid.UUID{},
			IsVerified:     false,
			VerifiedAt:     nil,
			VerifiedByUserID: nil,
		}
		td.InsertPartner(t, ctx, testPartner, testPool, crypto)

		// Act - verify the partner
		req := td.NewVerifyPartnerRequest(t, ctx, testServerURL, testPartner.ID)
		req.Header.Set("Cookie", td.ToCookieString(accessToken))

		resp, err := client.Do(req)

		// Assert HTTP response
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Verify partner verification status in database
		partnerEncx, err := td.GetPartnerByUserID(t, ctx, testUser.ID, testPool)
		require.NoError(t, err)
		partner, err := domain.DecryptPartnerEncx(ctx, crypto, partnerEncx)
		require.NoError(t, err)
		assert.True(t, partner.IsVerified, "Partner should be verified")
		assert.NotNil(t, partner.VerifiedAt, "VerifiedAt should be set")
		assert.NotNil(t, partner.VerifiedByUserID, "VerifiedByUserID should be set")

		// Verify user state and role were updated
		userEncx, err := td.GetUserEnxByID(t, ctx, testUser.ID, testPool, crypto)
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
		req := td.NewVerifyPartnerRequest(t, ctx, testServerURL, nonExistentID)
		req.Header.Set("Cookie", td.ToCookieString(accessToken))

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
		testUser := td.NewTestUser(t, "verified-partner@example.com", "Jane", "Verified")
		testUser.State = domain.Active
		testUser.Role = identity.PartnerStr
		testUserEncx, err := domain.ProcessUserEncx(ctx, crypto, testUser)
		require.NoError(t, err)
		err = td.InsertUserEncx(t, ctx, testUserEncx, testPool, crypto)
		require.NoError(t, err)

		// Create an already verified partner
		verifiedAt := time.Now()
		verifierID := uuid.New()
		testPartner := &domain.Partner{
			ID:               uuid.New(),
			UserID:           testUser.ID,
			Bio:              "Already verified",
			Experience:       "5 years",
			// Certifications:   []string{"Certification 1"},
			CategoryIDs:      []uuid.UUID{},
			ProductIDs:       []uuid.UUID{},
			IsVerified:       true,
			VerifiedAt:       &verifiedAt,
			VerifiedByUserID: &verifierID,
		}
		td.InsertPartner(t, ctx, testPartner, testPool, crypto)

		// Act - try to verify already verified partner
		req := td.NewVerifyPartnerRequest(t, ctx, testServerURL, testPartner.ID)
		req.Header.Set("Cookie", td.ToCookieString(accessToken))

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
		standardUser := td.NewTestUser(t, "standard@example.com", "Standard", "User")
		standardUser.State = domain.Active
		standardUserEncx, err := domain.ProcessUserEncx(ctx, crypto, standardUser)
		require.NoError(t, err)
		err = td.InsertUserEncx(t, ctx, standardUserEncx, testPool, crypto)
		require.NoError(t, err)

		accessToken := tu.SetupSessionForUser(t, ctx, authCtx, standardUser.ID, identity.Standard)

		// Create a partner to verify
		partnerUser := td.NewTestUser(t, "partner@example.com", "Partner", "User")
		partnerUser.State = domain.Pending
		partnerUserEncx, err := domain.ProcessUserEncx(ctx, crypto, partnerUser)
		require.NoError(t, err)
		err = td.InsertUserEncx(t, ctx, partnerUserEncx, testPool, crypto)
		require.NoError(t, err)

		testPartner := &domain.Partner{
			ID:             uuid.New(),
			UserID:         partnerUser.ID,
			Bio:            "Test bio",
			Experience:     "Test experience",
			// Certifications: []string{},
			CategoryIDs:    []uuid.UUID{},
			ProductIDs:     []uuid.UUID{},
			IsVerified:     false,
		}
		td.InsertPartner(t, ctx, testPartner, testPool, crypto)

		// Act - try to verify with non-admin role
		req := td.NewVerifyPartnerRequest(t, ctx, testServerURL, testPartner.ID)
		req.Header.Set("Cookie", td.ToCookieString(accessToken))

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
		testUser := td.NewTestUser(t, "partner@example.com", "Partner", "User")
		testUser.State = domain.Pending
		testUserEncx, err := domain.ProcessUserEncx(ctx, crypto, testUser)
		require.NoError(t, err)
		err = td.InsertUserEncx(t, ctx, testUserEncx, testPool, crypto)
		require.NoError(t, err)

		testPartner := &domain.Partner{
			ID:             uuid.New(),
			UserID:         testUser.ID,
			Bio:            "Test bio",
			Experience:     "Test experience",
			// Certifications: []string{},
			CategoryIDs:    []uuid.UUID{},
			ProductIDs:     []uuid.UUID{},
			IsVerified:     false,
		}
		td.InsertPartner(t, ctx, testPartner, testPool, crypto)

		// Act - try to verify without authentication
		req := td.NewVerifyPartnerRequest(t, ctx, testServerURL, testPartner.ID)
		// No session cookie

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
		req.Header.Set("Cookie", td.ToCookieString(accessToken))

		resp, err := client.Do(req)

		// Assert
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}
