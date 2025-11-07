package partner_test

import (
	"context"
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

	t.Run("should successfully get partner by ID", func(t *testing.T) {
		// Clean state
		td.ClearPartnersTable(t, ctx, testPool)

		// Create test user and associated partner
		user := td.NewTestUser(t, "partner@example.com", "John", "Partner")
		user.State = domain.Active
		user.Role = identity.PartnerStr
		userEncx, err := domain.ProcessUserEncx(ctx, crypto, user)
		require.NoError(t, err)
		err = td.InsertUserEncx(t, ctx, userEncx, testPool)
		require.NoError(t, err)

		partner := td.NewTestPartner(t, user.ID)
		partnerEncx, err := domain.ProcessPartnerEncx(ctx, crypto, partner)
		err = td.InsertPartnerEncx(t, ctx, partnerEncx, testPool)
		require.NoError(t, err)

		// Act
		req := td.NewGetPartnerByIDRequest(t, ctx, testServerURL, partner.ID)
		resp, err := client.Do(req)

		// Assert HTTP response
		assert.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Parse response body
		responsePartner := td.ParsePartnerResponse(t, resp)
		require.NotNil(t, responsePartner, "Response partner should not be nil")

		// Get encrypted partner from database
		dbPartnerEncx, err := td.GetPartnerEncxByID(t, ctx, partner.ID, testPool)
		require.NoError(t, err, "Failed to get partner from database")

		// Decrypt partner
		dbPartner, err := domain.DecryptPartnerEncx(ctx, crypto, dbPartnerEncx)
		require.NoError(t, err, "Failed to decrypt partner")

		// Verify all fields match database
		assert.Equal(t, dbPartner.ID, responsePartner.ID, "ID mismatch")
		assert.Equal(t, dbPartner.Bio, responsePartner.Bio, "Bio mismatch")
		assert.Equal(t, dbPartner.Experience, responsePartner.Experience, "Experience mismatch")
		assert.Equal(t, dbPartner.CategoryIDs, responsePartner.CategoryIDs, "CategoryIDs mismatch")
		assert.Equal(t, dbPartner.ProductIDs, responsePartner.ProductIDs, "ProductIDs mismatch")
		assert.WithinDuration(t, dbPartner.CreatedAt, responsePartner.CreatedAt, time.Second, "CreatedAt mismatch")
		assert.WithinDuration(t, dbPartner.UpdatedAt, responsePartner.UpdatedAt, time.Second, "UpdatedAt mismatch")
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

	t.Run("should return 500 when partner DEK is corrupted", func(t *testing.T) {
		// Clean state
		td.ClearPartnersTable(t, ctx, testPool)

		// Create test user and partner
		user := td.NewTestUser(t, "partner@example.com", "John", "Partner")
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

		// Corrupt the DEK to simulate decryption failure
		td.CorruptPartnerDEK(t, ctx, partner.ID, testPool)

		// Act
		req := td.NewGetPartnerByIDRequest(t, ctx, testServerURL, partner.ID)
		resp, err := client.Do(req)

		// Assert - should return 500 due to decryption failure
		assert.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})

	t.Run("should return 500 when key version is invalid", func(t *testing.T) {
		// Clean state
		td.ClearPartnersTable(t, ctx, testPool)

		// Create test user and partner
		user := td.NewTestUser(t, "partner2@example.com", "Jane", "Partner")
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

		// Set an invalid key version to simulate decryption failure
		td.SetInvalidKeyVersion(t, ctx, partner.ID, testPool, 99999)

		// Act
		req := td.NewGetPartnerByIDRequest(t, ctx, testServerURL, partner.ID)
		resp, err := client.Do(req)

		// Assert - should return 500 due to decryption failure
		assert.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}
