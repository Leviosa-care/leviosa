package partner_test

import (
	"context"
	"net/http"
	"strings"
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

// TEST=TestUpdatePartner make test-integration-partner-test

func TestUpdatePartner(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("should successfully update own partner profile", func(t *testing.T) {
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
			Bio:            "Original bio",
			Experience:     "Original experience",
			// Certifications: []string{"Original Cert"},
			CategoryIDs:    []uuid.UUID{},
			ProductIDs:     []uuid.UUID{},
			IsVerified:     true,
		}
		td.InsertPartner(t, ctx, testPartner, testPool, crypto)

		accessToken := tu.SetupSessionForUser(t, ctx, authCtx, partnerUser.ID, identity.Partner)

		// Act - update partner profile
		updatedBio := "Updated bio description"
		updatedExperience := "Updated experience details"
		updatedCerts := []string{"Cert 1", "Cert 2", "Cert 3"}

		updateRequest := domain.UpdatePartnerRequest{
			Bio:            &updatedBio,
			Experience:     &updatedExperience,
			// Certifications: &updatedCerts,
		}

		req := td.NewUpdatePartnerRequest(t, ctx, testServerURL, testPartner.ID, updateRequest)
		req.Header.Set("Cookie", td.ToCookieString(accessToken))

		resp, err := client.Do(req)

		// Assert
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Verify database changes
		partnerEncx, err := td.GetPartnerByUserID(t, ctx, partnerUser.ID, testPool)
		require.NoError(t, err)
		partner, err := domain.DecryptPartnerEncx(ctx, crypto, partnerEncx)
		require.NoError(t, err)

		assert.Equal(t, updatedBio, partner.Bio)
		assert.Equal(t, updatedExperience, partner.Experience)
		assert.Equal(t, updatedCerts, partner.Certifications)
	})

	t.Run("should successfully update any partner profile with admin role", func(t *testing.T) {
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
			Bio:            "Original bio",
			Experience:     "Original experience",
			// Certifications: []string{"Original Cert"},
			CategoryIDs:    []uuid.UUID{},
			ProductIDs:     []uuid.UUID{},
			IsVerified:     true,
		}
		td.InsertPartner(t, ctx, testPartner, testPool, crypto)

		// Act - admin updates partner profile
		updatedBio := "Admin updated bio"
		updateRequest := domain.UpdatePartnerRequest{
			Bio: &updatedBio,
		}

		req := td.NewUpdatePartnerRequest(t, ctx, testServerURL, testPartner.ID, updateRequest)
		req.Header.Set("Cookie", td.ToCookieString(accessToken))

		resp, err := client.Do(req)

		// Assert
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Verify database changes
		partnerEncx, err := td.GetPartnerByUserID(t, ctx, partnerUser.ID, testPool)
		require.NoError(t, err)
		partner, err := domain.DecryptPartnerEncx(ctx, crypto, partnerEncx)
		require.NoError(t, err)

		assert.Equal(t, updatedBio, partner.Bio)
		assert.Equal(t, "Original experience", partner.Experience, "Experience should remain unchanged")
	})

	t.Run("should support partial updates (only provided fields)", func(t *testing.T) {
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

		testPartner := &domain.Partner{
			ID:             uuid.New(),
			UserID:         partnerUser.ID,
			Bio:            "Original bio",
			Experience:     "Original experience",
			// Certifications: []string{"Cert 1", "Cert 2"},
			CategoryIDs:    []uuid.UUID{},
			ProductIDs:     []uuid.UUID{},
			IsVerified:     true,
		}
		td.InsertPartner(t, ctx, testPartner, testPool, crypto)

		accessToken := tu.SetupSessionForUser(t, ctx, authCtx, partnerUser.ID, identity.Partner)

		// Act - update only bio
		updatedBio := "Only bio updated"
		updateRequest := domain.UpdatePartnerRequest{
			Bio: &updatedBio,
		}

		req := td.NewUpdatePartnerRequest(t, ctx, testServerURL, testPartner.ID, updateRequest)
		req.Header.Set("Cookie", td.ToCookieString(accessToken))

		resp, err := client.Do(req)

		// Assert
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Verify only bio changed
		partnerEncx, err := td.GetPartnerByUserID(t, ctx, partnerUser.ID, testPool)
		require.NoError(t, err)
		partner, err := domain.DecryptPartnerEncx(ctx, crypto, partnerEncx)
		require.NoError(t, err)

		assert.Equal(t, updatedBio, partner.Bio, "Bio should be updated")
		assert.Equal(t, "Original experience", partner.Experience, "Experience should remain unchanged")
		assert.Equal(t, []string{"Cert 1", "Cert 2"}, partner.Certifications, "Certifications should remain unchanged")
	})

	t.Run("should return 403 when partner tries to update another partner's profile", func(t *testing.T) {
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
			ID:         uuid.New(),
			UserID:     partner1.ID,
			Bio:        "Partner 1 bio",
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
			ID:         uuid.New(),
			UserID:     partner2.ID,
			Bio:        "Partner 2 bio",
			IsVerified: true,
		}
		td.InsertPartner(t, ctx, testPartner2, testPool, crypto)

		// Act - partner1 tries to update partner2's profile
		updatedBio := "Malicious update"
		updateRequest := domain.UpdatePartnerRequest{
			Bio: &updatedBio,
		}

		req := td.NewUpdatePartnerRequest(t, ctx, testServerURL, testPartner2.ID, updateRequest)
		req.Header.Set("Cookie", td.ToCookieString(accessToken))

		resp, err := client.Do(req)

		// Assert
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)

		// Verify partner2's bio was not changed
		partnerEncx, err := td.GetPartnerByUserID(t, ctx, partner2.ID, testPool)
		require.NoError(t, err)
		partner, err := domain.DecryptPartnerEncx(ctx, crypto, partnerEncx)
		require.NoError(t, err)
		assert.Equal(t, "Partner 2 bio", partner.Bio, "Bio should remain unchanged")
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

		req := td.NewUpdatePartnerRequest(t, ctx, testServerURL, nonExistentID, updateRequest)
		req.Header.Set("Cookie", td.ToCookieString(accessToken))

		resp, err := client.Do(req)

		// Assert
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("should return 400 for bio exceeding max length", func(t *testing.T) {
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

		testPartner := &domain.Partner{
			ID:         uuid.New(),
			UserID:     partnerUser.ID,
			Bio:        "Original bio",
			IsVerified: true,
		}
		td.InsertPartner(t, ctx, testPartner, testPool, crypto)

		accessToken := tu.SetupSessionForUser(t, ctx, authCtx, partnerUser.ID, identity.Partner)

		// Act - bio exceeds 1000 characters
		invalidBio := strings.Repeat("a", 1001)
		updateRequest := domain.UpdatePartnerRequest{
			Bio: &invalidBio,
		}

		req := td.NewUpdatePartnerRequest(t, ctx, testServerURL, testPartner.ID, updateRequest)
		req.Header.Set("Cookie", td.ToCookieString(accessToken))

		resp, err := client.Do(req)

		// Assert
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should return 400 for experience exceeding max length", func(t *testing.T) {
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

		testPartner := &domain.Partner{
			ID:         uuid.New(),
			UserID:     partnerUser.ID,
			Experience: "Original experience",
			IsVerified: true,
		}
		td.InsertPartner(t, ctx, testPartner, testPool, crypto)

		accessToken := tu.SetupSessionForUser(t, ctx, authCtx, partnerUser.ID, identity.Partner)

		// Act - experience exceeds 2000 characters
		invalidExperience := strings.Repeat("a", 2001)
		updateRequest := domain.UpdatePartnerRequest{
			Experience: &invalidExperience,
		}

		req := td.NewUpdatePartnerRequest(t, ctx, testServerURL, testPartner.ID, updateRequest)
		req.Header.Set("Cookie", td.ToCookieString(accessToken))

		resp, err := client.Do(req)

		// Assert
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should return 400 for too many certifications", func(t *testing.T) {
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

		testPartner := &domain.Partner{
			ID:             uuid.New(),
			UserID:         partnerUser.ID,
			// Certifications: []string{"Cert 1"},
			IsVerified:     true,
		}
		td.InsertPartner(t, ctx, testPartner, testPool, crypto)

		accessToken := tu.SetupSessionForUser(t, ctx, authCtx, partnerUser.ID, identity.Partner)

		// Act - more than 20 certifications
		invalidCerts := make([]string, 21)
		for i := 0; i < 21; i++ {
			invalidCerts[i] = "Certification " + string(rune('A'+i))
		}

		updateRequest := domain.UpdatePartnerRequest{
			// Certifications: &invalidCerts,
		}

		req := td.NewUpdatePartnerRequest(t, ctx, testServerURL, testPartner.ID, updateRequest)
		req.Header.Set("Cookie", td.ToCookieString(accessToken))

		resp, err := client.Do(req)

		// Assert
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should return 400 for empty certification in array", func(t *testing.T) {
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

		testPartner := &domain.Partner{
			ID:             uuid.New(),
			UserID:         partnerUser.ID,
			// Certifications: []string{"Valid Cert"},
			IsVerified:     true,
		}
		td.InsertPartner(t, ctx, testPartner, testPool, crypto)

		accessToken := tu.SetupSessionForUser(t, ctx, authCtx, partnerUser.ID, identity.Partner)

		// Act - certifications array contains empty string
		invalidCerts := []string{"Valid Cert", "", "Another Cert"}
		updateRequest := domain.UpdatePartnerRequest{
			// Certifications: &invalidCerts,
		}

		req := td.NewUpdatePartnerRequest(t, ctx, testServerURL, testPartner.ID, updateRequest)
		req.Header.Set("Cookie", td.ToCookieString(accessToken))

		resp, err := client.Do(req)

		// Assert
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
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
			ID:         uuid.New(),
			UserID:     partnerUser.ID,
			Bio:        "Test bio",
			IsVerified: true,
		}
		td.InsertPartner(t, ctx, testPartner, testPool, crypto)

		// Act - no authentication
		updatedBio := "Updated bio"
		updateRequest := domain.UpdatePartnerRequest{
			Bio: &updatedBio,
		}

		req := td.NewUpdatePartnerRequest(t, ctx, testServerURL, testPartner.ID, updateRequest)
		// No session cookie

		resp, err := client.Do(req)

		// Assert
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("should return 400 for invalid partner ID format", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, redisClient)

		accessToken := tu.SetupAdminUser(t, ctx, authCtx)

		// Act - invalid UUID format
		updatedBio := "Updated bio"
		updateRequest := domain.UpdatePartnerRequest{
			Bio: &updatedBio,
		}

		// Manually construct request with invalid UUID
		req, err := http.NewRequestWithContext(
			ctx,
			http.MethodPut,
			testServerURL+"/partners/invalid-uuid",
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
