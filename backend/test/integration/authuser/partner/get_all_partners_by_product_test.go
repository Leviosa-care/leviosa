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

// make test-func TEST_NAME=TestGetAllPartnersByProduct TEST_PATH=test/integration/authuser/partner/get_all_partners_by_product_test.go

func TestGetAllPartnersByProduct(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("should successfully get all partners by product", func(t *testing.T) {
		// Clean state
		td.ClearPartnersTable(t, ctx, testPool)
		td.ClearSessionsRedis(t, ctx, redisClient)

		// Create a shared product ID that multiple partners will use
		sharedProductID := uuid.New()
		otherProductID := uuid.New()

		// Create 3 partners with the shared product
		partnersWithProduct := make([]*domain.Partner, 0, 3)
		for i := 0; i < 3; i++ {
			user := td.NewTestUser(t,
				"partner"+string(rune(i))+"@example.com",
				"Partner",
				string(rune('A'+i)))
			user.State = domain.Active
			user.Role = identity.PartnerStr
			testUserEncx, err := domain.ProcessUserEncx(ctx, crypto, user)
			require.NoError(t, err)
			err = td.InsertUserEncx(t, ctx, testUserEncx, testPool)
			require.NoError(t, err)

			partner := td.NewTestPartner(t, user.ID)
			// Assign the shared product to this partner
			partner.ProductIDs = []uuid.UUID{sharedProductID}
			partnerEncx, err := domain.ProcessPartnerEncx(ctx, crypto, partner)
			require.NoError(t, err)
			td.InsertPartnerEncx(t, ctx, partnerEncx, testPool)
			require.NoError(t, err)
			partnersWithProduct = append(partnersWithProduct, partner)

			// Small delay to ensure different timestamps
			time.Sleep(10 * time.Millisecond)
		}

		// Create 2 partners with a different product (should NOT be returned)
		for i := 0; i < 2; i++ {
			user := td.NewTestUser(t,
				"other"+string(rune(i))+"@example.com",
				"Other",
				string(rune('X'+i)))
			user.State = domain.Active
			user.Role = identity.PartnerStr
			testUserEncx, err := domain.ProcessUserEncx(ctx, crypto, user)
			require.NoError(t, err)
			err = td.InsertUserEncx(t, ctx, testUserEncx, testPool)
			require.NoError(t, err)

			partner := td.NewTestPartner(t, user.ID)
			// Assign a different product
			partner.ProductIDs = []uuid.UUID{otherProductID}
			partnerEncx, err := domain.ProcessPartnerEncx(ctx, crypto, partner)
			require.NoError(t, err)
			td.InsertPartnerEncx(t, ctx, partnerEncx, testPool)
			require.NoError(t, err)

			time.Sleep(10 * time.Millisecond)
		}

		// Act - request partners for the shared product
		req := td.NewGetAllPartnersByProductRequest(t, ctx, testServerURL, sharedProductID)
		resp, err := client.Do(req)

		// Assert HTTP response
		assert.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Parse response body
		partners := td.ParsePartnersListResponse(t, resp)
		require.Len(t, partners, 3, "Expected 3 partners in response for the shared product")

		// Verify partners are ordered by created_at DESC (newest first)
		for i := 0; i < len(partners)-1; i++ {
			assert.True(t, partners[i].CreatedAt.After(partners[i+1].CreatedAt) || partners[i].CreatedAt.Equal(partners[i+1].CreatedAt),
				"Partners should be ordered by created_at DESC")
		}

		// Verify each partner has the correct product
		for _, responsePartner := range partners {
			// Get encrypted partner from database by ID
			partnerEncx, err := td.GetPartnerEncxByID(t, ctx, responsePartner.ID, testPool)
			require.NoError(t, err, "Failed to get partner from database")

			// Decrypt partner
			dbPartner, err := domain.DecryptPartnerEncx(ctx, crypto, partnerEncx)
			require.NoError(t, err, "Failed to decrypt partner")

			// Verify the shared product is present in the partner's products
			assert.Contains(t, dbPartner.ProductIDs, sharedProductID, "Partner should have the shared product")
			assert.Equal(t, dbPartner.ID, responsePartner.ID, "ID mismatch")
			assert.Equal(t, dbPartner.Bio, responsePartner.Bio, "Bio mismatch")
			assert.Equal(t, dbPartner.Experience, responsePartner.Experience, "Experience mismatch")
		}
	})

	t.Run("should return empty array when no partners exist for product", func(t *testing.T) {
		// Clean state
		td.ClearPartnersTable(t, ctx, testPool)
		td.ClearSessionsRedis(t, ctx, redisClient)

		// Create partners with different products
		otherProductID := uuid.New()
		for i := 0; i < 2; i++ {
			user := td.NewTestUser(t,
				"partner"+string(rune(i))+"@example.com",
				"Partner",
				string(rune('A'+i)))
			user.State = domain.Active
			user.Role = identity.PartnerStr
			testUserEncx, err := domain.ProcessUserEncx(ctx, crypto, user)
			require.NoError(t, err)
			err = td.InsertUserEncx(t, ctx, testUserEncx, testPool)
			require.NoError(t, err)

			partner := td.NewTestPartner(t, user.ID)
			partner.ProductIDs = []uuid.UUID{otherProductID}
			partnerEncx, err := domain.ProcessPartnerEncx(ctx, crypto, partner)
			require.NoError(t, err)
			td.InsertPartnerEncx(t, ctx, partnerEncx, testPool)
			require.NoError(t, err)
		}

		// Act - request partners for a product that no partner has
		nonExistentProductID := uuid.New()
		req := td.NewGetAllPartnersByProductRequest(t, ctx, testServerURL, nonExistentProductID)
		resp, err := client.Do(req)

		// Assert HTTP response
		assert.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Parse and verify empty array structure
		partners := td.ParsePartnersListResponse(t, resp)
		assert.Empty(t, partners, "Expected empty partners array")
		assert.NotNil(t, partners, "Partners array should not be nil, should be empty array")
	})

	t.Run("should return 400 when product ID is invalid", func(t *testing.T) {
		// Clean state
		td.ClearPartnersTable(t, ctx, testPool)

		// Act - request with invalid product ID
		req := td.NewGetAllPartnersByProductRequestWithInvalidID(t, ctx, testServerURL, "invalid-uuid")
		resp, err := client.Do(req)

		// Assert - should return 400 bad request
		assert.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should return 404 when product ID is empty", func(t *testing.T) {
		// Clean state
		td.ClearPartnersTable(t, ctx, testPool)

		// Act - request with empty product ID (URL: /partners/products/)
		req := td.NewGetAllPartnersByProductRequestWithInvalidID(t, ctx, testServerURL, "")
		resp, err := client.Do(req)

		// Assert - should return 404 because route doesn't match (Go's http.ServeMux behavior)
		assert.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("should return 500 when partner DEK is corrupted", func(t *testing.T) {
		// Clean state
		td.ClearPartnersTable(t, ctx, testPool)

		// Create product and partner
		productID := uuid.New()
		user := td.NewTestUser(t, "partner@example.com", "John", "Partner")
		user.State = domain.Active
		user.Role = identity.PartnerStr
		userEncx, err := domain.ProcessUserEncx(ctx, crypto, user)
		require.NoError(t, err)
		err = td.InsertUserEncx(t, ctx, userEncx, testPool)
		require.NoError(t, err)

		partner := td.NewTestPartner(t, user.ID)
		partner.ProductIDs = []uuid.UUID{productID}
		partnerEncx, err := domain.ProcessPartnerEncx(ctx, crypto, partner)
		require.NoError(t, err)
		err = td.InsertPartnerEncx(t, ctx, partnerEncx, testPool)
		require.NoError(t, err)

		// Corrupt the DEK to simulate decryption failure
		td.CorruptPartnerDEK(t, ctx, partner.ID, testPool)

		// Act
		req := td.NewGetAllPartnersByProductRequest(t, ctx, testServerURL, productID)
		resp, err := client.Do(req)

		// Assert - should return 500 due to decryption failure
		assert.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})

	t.Run("should return 500 when key version is invalid", func(t *testing.T) {
		// Clean state
		td.ClearPartnersTable(t, ctx, testPool)

		// Create product and partner
		productID := uuid.New()
		user := td.NewTestUser(t, "partner2@example.com", "Jane", "Partner")
		user.State = domain.Active
		user.Role = identity.PartnerStr
		userEncx, err := domain.ProcessUserEncx(ctx, crypto, user)
		require.NoError(t, err)
		err = td.InsertUserEncx(t, ctx, userEncx, testPool)
		require.NoError(t, err)

		partner := td.NewTestPartner(t, user.ID)
		partner.ProductIDs = []uuid.UUID{productID}
		partnerEncx, err := domain.ProcessPartnerEncx(ctx, crypto, partner)
		require.NoError(t, err)
		err = td.InsertPartnerEncx(t, ctx, partnerEncx, testPool)
		require.NoError(t, err)

		// Set an invalid key version to simulate decryption failure
		td.SetInvalidKeyVersion(t, ctx, partner.ID, testPool, 99999)

		// Act
		req := td.NewGetAllPartnersByProductRequest(t, ctx, testServerURL, productID)
		resp, err := client.Do(req)

		// Assert - should return 500 due to decryption failure
		assert.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})

	t.Run("should handle partners with multiple products correctly", func(t *testing.T) {
		// Clean state
		td.ClearPartnersTable(t, ctx, testPool)
		td.ClearSessionsRedis(t, ctx, redisClient)

		// Create target product and other products
		targetProductID := uuid.New()
		product1 := uuid.New()
		product2 := uuid.New()

		// Create partner with multiple products including the target
		user1 := td.NewTestUser(t, "multi1@example.com", "Multi", "One")
		user1.State = domain.Active
		user1.Role = identity.PartnerStr
		userEncx1, err := domain.ProcessUserEncx(ctx, crypto, user1)
		require.NoError(t, err)
		err = td.InsertUserEncx(t, ctx, userEncx1, testPool)
		require.NoError(t, err)

		partner1 := td.NewTestPartner(t, user1.ID)
		partner1.ProductIDs = []uuid.UUID{product1, targetProductID, product2}
		partnerEncx1, err := domain.ProcessPartnerEncx(ctx, crypto, partner1)
		require.NoError(t, err)
		td.InsertPartnerEncx(t, ctx, partnerEncx1, testPool)
		require.NoError(t, err)

		// Create partner with only target product
		user2 := td.NewTestUser(t, "multi2@example.com", "Multi", "Two")
		user2.State = domain.Active
		user2.Role = identity.PartnerStr
		userEncx2, err := domain.ProcessUserEncx(ctx, crypto, user2)
		require.NoError(t, err)
		err = td.InsertUserEncx(t, ctx, userEncx2, testPool)
		require.NoError(t, err)

		partner2 := td.NewTestPartner(t, user2.ID)
		partner2.ProductIDs = []uuid.UUID{targetProductID}
		partnerEncx2, err := domain.ProcessPartnerEncx(ctx, crypto, partner2)
		require.NoError(t, err)
		td.InsertPartnerEncx(t, ctx, partnerEncx2, testPool)
		require.NoError(t, err)

		// Create partner without target product
		user3 := td.NewTestUser(t, "multi3@example.com", "Multi", "Three")
		user3.State = domain.Active
		user3.Role = identity.PartnerStr
		userEncx3, err := domain.ProcessUserEncx(ctx, crypto, user3)
		require.NoError(t, err)
		err = td.InsertUserEncx(t, ctx, userEncx3, testPool)
		require.NoError(t, err)

		partner3 := td.NewTestPartner(t, user3.ID)
		partner3.ProductIDs = []uuid.UUID{product1, product2}
		partnerEncx3, err := domain.ProcessPartnerEncx(ctx, crypto, partner3)
		require.NoError(t, err)
		td.InsertPartnerEncx(t, ctx, partnerEncx3, testPool)
		require.NoError(t, err)

		// Act - request partners for target product
		req := td.NewGetAllPartnersByProductRequest(t, ctx, testServerURL, targetProductID)
		resp, err := client.Do(req)

		// Assert HTTP response
		assert.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Parse response body - should return only partners with target product
		partners := td.ParsePartnersListResponse(t, resp)
		require.Len(t, partners, 2, "Expected 2 partners with the target product")

		// Verify both partners have the target product
		partnerIDs := []uuid.UUID{partner1.ID, partner2.ID}
		for _, responsePartner := range partners {
			assert.Contains(t, partnerIDs, responsePartner.ID, "Returned partner should be one of the expected partners")
			assert.Contains(t, responsePartner.ProductIDs, targetProductID, "Partner should have target product")
		}

		// Verify partner3 is NOT in the results
		for _, responsePartner := range partners {
			assert.NotEqual(t, partner3.ID, responsePartner.ID, "Partner without target product should not be returned")
		}
	})
}
