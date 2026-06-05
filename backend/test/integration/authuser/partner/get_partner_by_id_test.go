package partner_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/contracts/identity"
	td "github.com/Leviosa-care/leviosa/backend/test/helpers"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestGetPartnerByID TEST_PATH=test/integration/authuser/partner/get_partner_by_id_test.go

func TestGetPartnerByID(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("should return enriched public partner by ID", func(t *testing.T) {
		// Clean state
		td.ClearPartnersTable(t, ctx, testPool)

		// Create test user and associated partner
		user := td.NewTestUser(t, "partner@example.com", "Alice", "Martin")
		user.State = domain.Active
		user.Role = identity.PartnerStr
		user.Picture = "https://example.com/alice.jpg"
		userEncx, err := domain.ProcessUserEncx(ctx, crypto, user)
		require.NoError(t, err)
		err = td.InsertUserEncx(t, ctx, userEncx, testPool)
		require.NoError(t, err)

		partner := td.NewTestPartner(t, user.ID)
		partner.Occupation = "Kinésithérapeute du sport"
		partner.Quote = "Le mouvement est la vie"
		partner.Tags = []string{"sport", "rééducation"}
		partner.Bio = "Bio détaillée"
		partner.Experience = "10 ans d'expérience"
		partner.StripeAccountStatus = domain.StripeAccountStatusActive
		partner.CategoryIDs = []uuid.UUID{uuid.New()}
		partnerEncx, err := domain.ProcessPartnerEncx(ctx, crypto, partner)
		require.NoError(t, err)
		err = td.InsertPartnerEncx(t, ctx, partnerEncx, testPool)
		require.NoError(t, err)

		// Act
		req := td.NewGetPartnerByIDRequest(t, ctx, testServerURL, partner.ID)
		resp, err := client.Do(req)

		// Assert HTTP response
		assert.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Parse response body as PublicPartnerResponse
		var responsePartner domain.PublicPartnerResponse
		err = json.NewDecoder(resp.Body).Decode(&responsePartner)
		require.NoError(t, err)
		require.NotNil(t, responsePartner.ID)

		// Verify enriched fields
		assert.Equal(t, partner.ID, responsePartner.ID)
		assert.Equal(t, "Alice", responsePartner.FirstName)
		assert.Equal(t, "Martin", responsePartner.LastName)
		assert.Equal(t, "Kinésithérapeute du sport", responsePartner.Occupation)
		assert.Equal(t, "Le mouvement est la vie", responsePartner.Quote)
		assert.Equal(t, []string{"sport", "rééducation"}, responsePartner.Tags)
		assert.Equal(t, "Bio détaillée", responsePartner.Bio)
		assert.Equal(t, "10 ans d'expérience", responsePartner.Experience)
		assert.Equal(t, partner.CategoryIDs, responsePartner.CategoryIDs)
	})

	t.Run("should return 404 when partner not found", func(t *testing.T) {
		// Clean state
		td.ClearPartnersTable(t, ctx, testPool)

		// Act
		nonExistentID := uuid.New()
		req := td.NewGetPartnerByIDRequest(t, ctx, testServerURL, nonExistentID)

		resp, err := client.Do(req)

		// Assert
		assert.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("should return 404 for disabled partner", func(t *testing.T) {
		// Clean state
		td.ClearPartnersTable(t, ctx, testPool)

		user := td.NewTestUser(t, "disabled@example.com", "Bob", "Durand")
		user.State = domain.Active
		user.Role = identity.PartnerStr
		userEncx, err := domain.ProcessUserEncx(ctx, crypto, user)
		require.NoError(t, err)
		err = td.InsertUserEncx(t, ctx, userEncx, testPool)
		require.NoError(t, err)

		partner := td.NewTestPartner(t, user.ID)
		partner.StripeAccountStatus = domain.StripeAccountStatusDisabled
		partnerEncx, err := domain.ProcessPartnerEncx(ctx, crypto, partner)
		require.NoError(t, err)
		err = td.InsertPartnerEncx(t, ctx, partnerEncx, testPool)
		require.NoError(t, err)

		// Act
		req := td.NewGetPartnerByIDRequest(t, ctx, testServerURL, partner.ID)
		resp, err := client.Do(req)

		assert.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("should return 404 for inactive user partner", func(t *testing.T) {
		// Clean state
		td.ClearPartnersTable(t, ctx, testPool)

		user := td.NewTestUser(t, "pending@example.com", "Charlie", "Lambert")
		user.State = domain.Pending
		user.Role = identity.PartnerStr
		userEncx, err := domain.ProcessUserEncx(ctx, crypto, user)
		require.NoError(t, err)
		err = td.InsertUserEncx(t, ctx, userEncx, testPool)
		require.NoError(t, err)

		partner := td.NewTestPartner(t, user.ID)
		partner.StripeAccountStatus = domain.StripeAccountStatusActive
		partnerEncx, err := domain.ProcessPartnerEncx(ctx, crypto, partner)
		require.NoError(t, err)
		err = td.InsertPartnerEncx(t, ctx, partnerEncx, testPool)
		require.NoError(t, err)

		// Act
		req := td.NewGetPartnerByIDRequest(t, ctx, testServerURL, partner.ID)
		resp, err := client.Do(req)

		assert.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("should return 400 for invalid UUID format", func(t *testing.T) {
		// Clean state
		td.ClearSessionsRedis(t, ctx, redisClient)

		// Act
		req, err := http.NewRequestWithContext(
			ctx,
			http.MethodGet,
			testServerURL+"/partners/invalid-uuid",
			nil,
		)
		require.NoError(t, err)

		resp, err := client.Do(req)

		// Assert
		assert.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}
