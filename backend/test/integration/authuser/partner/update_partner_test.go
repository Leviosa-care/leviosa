package partner_test

import (
	"context"
	"net/http"
	"strings"
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

// make test-func TEST_NAME=TestUpdatePartner TEST_PATH=test/integration/authuser/partner/update_partner_test.go

func TestUpdatePartner(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("should successfully update partner bio", func(t *testing.T) {
		// Clean state
		td.ClearPartnersTable(t, ctx, testPool)
		td.ClearSessionsRedis(t, ctx, redisClient)

		// Use admin for simplicity
		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		// Create partner user
		partnerUser := td.NewTestUser(t, "partner1@example.com", "Partner", "One")
		partnerUser.State = domain.Active
		partnerUser.Role = identity.PartnerStr
		partnerUserEncx, err := domain.ProcessUserEncx(ctx, crypto, partnerUser)
		require.NoError(t, err)
		err = td.InsertUserEncx(t, ctx, partnerUserEncx, testPool)
		require.NoError(t, err)

		// Create partner profile
		testPartner := td.NewTestPartner(t, partnerUser.ID)
		testPartner.Bio = "Original bio"
		testPartner.Experience = "Original experience"
		partnerEncx, err := domain.ProcessPartnerEncx(ctx, crypto, testPartner)
		require.NoError(t, err)
		err = td.InsertPartnerEncx(t, ctx, partnerEncx, testPool)
		require.NoError(t, err)

		// Act - update partner bio
		updatedBio := "Updated bio description"
		updateRequest := domain.UpdatePartnerRequest{
			Bio: &updatedBio,
		}

		req := td.NewUpdatePartnerRequest(t, ctx, testServerURL, testPartner.ID, updateRequest, accessToken)
		resp, err := client.Do(req)

		// Assert HTTP response
		assert.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Parse response
		partnerResp := td.ParsePartnerResponse(t, resp)
		assert.Equal(t, updatedBio, partnerResp.Bio)
		assert.Equal(t, "Original experience", partnerResp.Experience)

		// Verify database changes
		retrievedPartnerEncx, err := td.GetPartnerEncxByID(t, ctx, testPartner.ID, testPool)
		assert.NoError(t, err)
		retrievedPartner, err := domain.DecryptPartnerEncx(ctx, crypto, retrievedPartnerEncx)
		assert.NoError(t, err)

		assert.Equal(t, updatedBio, retrievedPartner.Bio)
		assert.Equal(t, "Original experience", retrievedPartner.Experience)
	})

	t.Run("should successfully update partner experience", func(t *testing.T) {
		// Clean state
		td.ClearPartnersTable(t, ctx, testPool)
		td.ClearSessionsRedis(t, ctx, redisClient)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		// Create partner user
		partnerUser := td.NewTestUser(t, "partner2@example.com", "Partner", "Two")
		partnerUser.State = domain.Active
		partnerUser.Role = identity.PartnerStr
		partnerUserEncx, err := domain.ProcessUserEncx(ctx, crypto, partnerUser)
		require.NoError(t, err)
		err = td.InsertUserEncx(t, ctx, partnerUserEncx, testPool)
		require.NoError(t, err)

		// Create partner profile
		testPartner := td.NewTestPartner(t, partnerUser.ID)
		testPartner.Bio = "Original bio"
		testPartner.Experience = "Original experience"
		partnerEncx, err := domain.ProcessPartnerEncx(ctx, crypto, testPartner)
		require.NoError(t, err)
		err = td.InsertPartnerEncx(t, ctx, partnerEncx, testPool)
		require.NoError(t, err)

		// Act - update only experience
		updatedExperience := "Updated experience details"
		updateRequest := domain.UpdatePartnerRequest{
			Experience: &updatedExperience,
		}

		req := td.NewUpdatePartnerRequest(t, ctx, testServerURL, testPartner.ID, updateRequest, accessToken)
		resp, err := client.Do(req)

		// Assert
		assert.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Verify database changes
		retrievedPartnerEncx, err := td.GetPartnerEncxByID(t, ctx, testPartner.ID, testPool)
		assert.NoError(t, err)
		retrievedPartner, err := domain.DecryptPartnerEncx(ctx, crypto, retrievedPartnerEncx)
		assert.NoError(t, err)

		assert.Equal(t, "Original bio", retrievedPartner.Bio, "Bio should remain unchanged")
		assert.Equal(t, updatedExperience, retrievedPartner.Experience)
	})

	t.Run("should successfully update multiple fields at once", func(t *testing.T) {
		// Clean state
		td.ClearPartnersTable(t, ctx, testPool)
		td.ClearSessionsRedis(t, ctx, redisClient)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		// Create partner user
		partnerUser := td.NewTestUser(t, "partner3@example.com", "Partner", "Three")
		partnerUser.State = domain.Active
		partnerUser.Role = identity.PartnerStr
		partnerUserEncx, err := domain.ProcessUserEncx(ctx, crypto, partnerUser)
		require.NoError(t, err)
		err = td.InsertUserEncx(t, ctx, partnerUserEncx, testPool)
		require.NoError(t, err)

		// Create partner profile
		testPartner := td.NewTestPartner(t, partnerUser.ID)
		testPartner.Bio = "Original bio"
		testPartner.Experience = "Original experience"
		partnerEncx, err := domain.ProcessPartnerEncx(ctx, crypto, testPartner)
		require.NoError(t, err)
		err = td.InsertPartnerEncx(t, ctx, partnerEncx, testPool)
		require.NoError(t, err)

		// Act - update both bio and experience
		updatedBio := "Updated bio description"
		updatedExperience := "Updated experience details"
		updateRequest := domain.UpdatePartnerRequest{
			Bio:        &updatedBio,
			Experience: &updatedExperience,
		}

		req := td.NewUpdatePartnerRequest(t, ctx, testServerURL, testPartner.ID, updateRequest, accessToken)
		resp, err := client.Do(req)

		// Assert
		assert.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Verify database changes
		retrievedPartnerEncx, err := td.GetPartnerEncxByID(t, ctx, testPartner.ID, testPool)
		assert.NoError(t, err)
		retrievedPartner, err := domain.DecryptPartnerEncx(ctx, crypto, retrievedPartnerEncx)
		assert.NoError(t, err)

		assert.Equal(t, updatedBio, retrievedPartner.Bio)
		assert.Equal(t, updatedExperience, retrievedPartner.Experience)
	})

	t.Run("should allow admin to update any partner profile", func(t *testing.T) {
		// Clean state
		td.ClearPartnersTable(t, ctx, testPool)
		td.ClearSessionsRedis(t, ctx, redisClient)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		// Create partner user
		partnerUser := td.NewTestUser(t, "partner4@example.com", "John", "Partner")
		partnerUser.State = domain.Active
		partnerUser.Role = identity.PartnerStr
		partnerUserEncx, err := domain.ProcessUserEncx(ctx, crypto, partnerUser)
		require.NoError(t, err)
		err = td.InsertUserEncx(t, ctx, partnerUserEncx, testPool)
		require.NoError(t, err)

		// Create partner profile
		testPartner := td.NewTestPartner(t, partnerUser.ID)
		testPartner.Bio = "Original bio"
		partnerEncx, err := domain.ProcessPartnerEncx(ctx, crypto, testPartner)
		require.NoError(t, err)
		err = td.InsertPartnerEncx(t, ctx, partnerEncx, testPool)
		require.NoError(t, err)

		// Act - admin updates partner profile
		updatedBio := "Admin updated bio"
		updateRequest := domain.UpdatePartnerRequest{
			Bio: &updatedBio,
		}

		req := td.NewUpdatePartnerRequest(t, ctx, testServerURL, testPartner.ID, updateRequest, accessToken)
		resp, err := client.Do(req)

		// Assert
		assert.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Verify database changes
		retrievedPartnerEncx, err := td.GetPartnerEncxByID(t, ctx, testPartner.ID, testPool)
		assert.NoError(t, err)
		retrievedPartner, err := domain.DecryptPartnerEncx(ctx, crypto, retrievedPartnerEncx)
		assert.NoError(t, err)

		assert.Equal(t, updatedBio, retrievedPartner.Bio)
	})

	t.Run("should return 404 when partner not found", func(t *testing.T) {
		// Clean state
		td.ClearPartnersTable(t, ctx, testPool)
		td.ClearSessionsRedis(t, ctx, redisClient)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		// Act - try to update non-existent partner
		nonExistentID := uuid.New()
		updatedBio := "Updated bio"
		updateRequest := domain.UpdatePartnerRequest{
			Bio: &updatedBio,
		}

		req := td.NewUpdatePartnerRequest(t, ctx, testServerURL, nonExistentID, updateRequest, accessToken)
		resp, err := client.Do(req)

		// Assert
		assert.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("should return 400 for bio exceeding max length", func(t *testing.T) {
		// Clean state
		td.ClearPartnersTable(t, ctx, testPool)
		td.ClearSessionsRedis(t, ctx, redisClient)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		// Create partner user
		partnerUser := td.NewTestUser(t, "partner5@example.com", "Partner", "Five")
		partnerUser.State = domain.Active
		partnerUser.Role = identity.PartnerStr
		partnerUserEncx, err := domain.ProcessUserEncx(ctx, crypto, partnerUser)
		require.NoError(t, err)
		err = td.InsertUserEncx(t, ctx, partnerUserEncx, testPool)
		require.NoError(t, err)

		// Create partner profile
		testPartner := td.NewTestPartner(t, partnerUser.ID)
		partnerEncx, err := domain.ProcessPartnerEncx(ctx, crypto, testPartner)
		require.NoError(t, err)
		err = td.InsertPartnerEncx(t, ctx, partnerEncx, testPool)
		require.NoError(t, err)

		// Act - bio exceeds 1000 characters
		invalidBio := strings.Repeat("a", 1001)
		updateRequest := domain.UpdatePartnerRequest{
			Bio: &invalidBio,
		}

		req := td.NewUpdatePartnerRequest(t, ctx, testServerURL, testPartner.ID, updateRequest, accessToken)
		resp, err := client.Do(req)

		// Assert
		assert.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should return 400 for experience exceeding max length", func(t *testing.T) {
		// Clean state
		td.ClearPartnersTable(t, ctx, testPool)
		td.ClearSessionsRedis(t, ctx, redisClient)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		// Create partner user
		partnerUser := td.NewTestUser(t, "partner6@example.com", "Partner", "Six")
		partnerUser.State = domain.Active
		partnerUser.Role = identity.PartnerStr
		partnerUserEncx, err := domain.ProcessUserEncx(ctx, crypto, partnerUser)
		require.NoError(t, err)
		err = td.InsertUserEncx(t, ctx, partnerUserEncx, testPool)
		require.NoError(t, err)

		// Create partner profile
		testPartner := td.NewTestPartner(t, partnerUser.ID)
		partnerEncx, err := domain.ProcessPartnerEncx(ctx, crypto, testPartner)
		require.NoError(t, err)
		err = td.InsertPartnerEncx(t, ctx, partnerEncx, testPool)
		require.NoError(t, err)

		// Act - experience exceeds 2000 characters
		invalidExperience := strings.Repeat("a", 2001)
		updateRequest := domain.UpdatePartnerRequest{
			Experience: &invalidExperience,
		}

		req := td.NewUpdatePartnerRequest(t, ctx, testServerURL, testPartner.ID, updateRequest, accessToken)
		resp, err := client.Do(req)

		// Assert
		assert.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should return 401 without authentication", func(t *testing.T) {
		// Clean state
		td.ClearPartnersTable(t, ctx, testPool)

		// Create partner
		partnerUser := td.NewTestUser(t, "partner7@example.com", "Partner", "Seven")
		partnerUser.State = domain.Active
		partnerUserEncx, err := domain.ProcessUserEncx(ctx, crypto, partnerUser)
		require.NoError(t, err)
		err = td.InsertUserEncx(t, ctx, partnerUserEncx, testPool)
		require.NoError(t, err)

		testPartner := td.NewTestPartner(t, partnerUser.ID)
		partnerEncx, err := domain.ProcessPartnerEncx(ctx, crypto, testPartner)
		require.NoError(t, err)
		err = td.InsertPartnerEncx(t, ctx, partnerEncx, testPool)
		require.NoError(t, err)

		// Act - no authentication
		updatedBio := "Updated bio"
		updateRequest := domain.UpdatePartnerRequest{
			Bio: &updatedBio,
		}

		req := td.NewUpdatePartnerRequest(t, ctx, testServerURL, testPartner.ID, updateRequest, "")

		resp, err := client.Do(req)

		// Assert
		assert.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("should return 400 for invalid partner ID format", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, redisClient)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		// Act - invalid UUID format
		// Manually construct request with invalid UUID
		req, err := http.NewRequestWithContext(
			ctx,
			http.MethodPut,
			testServerURL+"/partners/invalid-uuid",
			strings.NewReader(`{"bio":"Updated bio"}`),
		)
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

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

	t.Run("should handle encrypted partner data correctly", func(t *testing.T) {
		// Clean state
		td.ClearPartnersTable(t, ctx, testPool)
		td.ClearSessionsRedis(t, ctx, redisClient)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		// Create partner user with full encrypted data
		partnerUser := td.NewTestUser(t, "partner8@example.com", "Partner", "Eight")
		partnerUser.State = domain.Active
		partnerUser.Role = identity.PartnerStr
		userEncx, err := domain.ProcessUserEncx(ctx, crypto, partnerUser)
		require.NoError(t, err)
		err = td.InsertUserEncx(t, ctx, userEncx, testPool)
		require.NoError(t, err)

		// Partner with full data including Stripe info
		partner := td.NewTestPartner(t, partnerUser.ID)
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

		// Act - update partner
		updatedBio := "Updated bio with new information"
		updateRequest := domain.UpdatePartnerRequest{
			Bio: &updatedBio,
		}

		req := td.NewUpdatePartnerRequest(t, ctx, testServerURL, partner.ID, updateRequest, accessToken)
		resp, err := client.Do(req)

		// Assert
		assert.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Verify partner was updated correctly
		retrievedPartnerEncx, err := td.GetPartnerEncxByID(t, ctx, partner.ID, testPool)
		assert.NoError(t, err)
		retrievedPartner, err := domain.DecryptPartnerEncx(ctx, crypto, retrievedPartnerEncx)
		assert.NoError(t, err)

		assert.Equal(t, updatedBio, retrievedPartner.Bio)
		assert.Equal(t, "10+ years in healthcare", retrievedPartner.Experience, "Experience should remain unchanged")
		assert.Equal(t, "acct_1234567890", retrievedPartner.StripeConnectedAccountID, "Stripe data should remain unchanged")
		assert.Equal(t, 2, len(retrievedPartner.CategoryIDs), "CategoryIDs should remain unchanged")
		assert.Equal(t, 3, len(retrievedPartner.ProductIDs), "ProductIDs should remain unchanged")
	})
}
