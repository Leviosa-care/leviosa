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

// make test-func TEST_NAME=TestGetAllPartnersByProducts TEST_PATH=test/integration/authuser/partner/get_all_partners_by_products_test.go

func TestGetAllPartnersByProducts(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("should successfully get all partners by multiple products", func(t *testing.T) {
		// Clean state
		td.ClearPartnersTable(t, ctx, testPool)
		td.ClearSessionsRedis(t, ctx, redisClient)

		// Create 3 different products
		product1 := uuid.New()
		product2 := uuid.New()
		product3 := uuid.New()

		// Create partner 1 with products 1 and 2
		user1 := td.NewTestUser(t, "partner1@example.com", "Partner", "One")
		user1.State = domain.Active
		user1.Role = identity.PartnerStr
		userEncx1, err := domain.ProcessUserEncx(ctx, crypto, user1)
		require.NoError(t, err)
		err = td.InsertUserEncx(t, ctx, userEncx1, testPool)
		require.NoError(t, err)

		partner1 := td.NewTestPartner(t, user1.ID)
		partner1.ProductIDs = []uuid.UUID{product1, product2}
		partnerEncx1, err := domain.ProcessPartnerEncx(ctx, crypto, partner1)
		require.NoError(t, err)
		td.InsertPartnerEncx(t, ctx, partnerEncx1, testPool)
		require.NoError(t, err)
		time.Sleep(10 * time.Millisecond)

		// Create partner 2 with products 2 and 3
		user2 := td.NewTestUser(t, "partner2@example.com", "Partner", "Two")
		user2.State = domain.Active
		user2.Role = identity.PartnerStr
		userEncx2, err := domain.ProcessUserEncx(ctx, crypto, user2)
		require.NoError(t, err)
		err = td.InsertUserEncx(t, ctx, userEncx2, testPool)
		require.NoError(t, err)

		partner2 := td.NewTestPartner(t, user2.ID)
		partner2.ProductIDs = []uuid.UUID{product2, product3}
		partnerEncx2, err := domain.ProcessPartnerEncx(ctx, crypto, partner2)
		require.NoError(t, err)
		td.InsertPartnerEncx(t, ctx, partnerEncx2, testPool)
		require.NoError(t, err)
		time.Sleep(10 * time.Millisecond)

		// Create partner 3 with only product 1
		user3 := td.NewTestUser(t, "partner3@example.com", "Partner", "Three")
		user3.State = domain.Active
		user3.Role = identity.PartnerStr
		userEncx3, err := domain.ProcessUserEncx(ctx, crypto, user3)
		require.NoError(t, err)
		err = td.InsertUserEncx(t, ctx, userEncx3, testPool)
		require.NoError(t, err)

		partner3 := td.NewTestPartner(t, user3.ID)
		partner3.ProductIDs = []uuid.UUID{product1}
		partnerEncx3, err := domain.ProcessPartnerEncx(ctx, crypto, partner3)
		require.NoError(t, err)
		td.InsertPartnerEncx(t, ctx, partnerEncx3, testPool)
		require.NoError(t, err)

		// Act - request partners for products 1 and 3 (should return partners 1, 2, and 3)
		req := td.NewGetAllPartnersByProductsRequest(t, ctx, testServerURL, []uuid.UUID{product1, product3})
		resp, err := client.Do(req)

		// Assert HTTP response
		assert.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Parse response body
		partners := td.ParsePartnersListResponse(t, resp)
		require.Len(t, partners, 3, "Expected 3 partners matching products 1 or 3")

		// Verify partners are ordered by created_at DESC (newest first)
		for i := 0; i < len(partners)-1; i++ {
			assert.True(t, partners[i].CreatedAt.After(partners[i+1].CreatedAt) || partners[i].CreatedAt.Equal(partners[i+1].CreatedAt),
				"Partners should be ordered by created_at DESC")
		}

		// Verify all returned partners have at least one of the requested products
		partnerIDs := map[uuid.UUID]bool{partner1.ID: true, partner2.ID: true, partner3.ID: true}
		for _, responsePartner := range partners {
			assert.True(t, partnerIDs[responsePartner.ID], "Returned partner should be one of the expected partners")

			// Verify partner has at least one of the requested products
			hasProduct := false
			for _, prodID := range responsePartner.ProductIDs {
				if prodID == product1 || prodID == product3 {
					hasProduct = true
					break
				}
			}
			assert.True(t, hasProduct, "Partner should have at least one of the requested products")
		}
	})

	t.Run("should return partners matching any of the provided products", func(t *testing.T) {
		// Clean state
		td.ClearPartnersTable(t, ctx, testPool)

		// Create products
		productA := uuid.New()
		productB := uuid.New()
		productC := uuid.New()

		// Partner with only product A
		user1 := td.NewTestUser(t, "partnerA@example.com", "Partner", "A")
		user1.State = domain.Active
		user1.Role = identity.PartnerStr
		userEncx1, err := domain.ProcessUserEncx(ctx, crypto, user1)
		require.NoError(t, err)
		err = td.InsertUserEncx(t, ctx, userEncx1, testPool)
		require.NoError(t, err)

		partner1 := td.NewTestPartner(t, user1.ID)
		partner1.ProductIDs = []uuid.UUID{productA}
		partnerEncx1, err := domain.ProcessPartnerEncx(ctx, crypto, partner1)
		require.NoError(t, err)
		td.InsertPartnerEncx(t, ctx, partnerEncx1, testPool)
		require.NoError(t, err)

		// Partner with only product B
		user2 := td.NewTestUser(t, "partnerB@example.com", "Partner", "B")
		user2.State = domain.Active
		user2.Role = identity.PartnerStr
		userEncx2, err := domain.ProcessUserEncx(ctx, crypto, user2)
		require.NoError(t, err)
		err = td.InsertUserEncx(t, ctx, userEncx2, testPool)
		require.NoError(t, err)

		partner2 := td.NewTestPartner(t, user2.ID)
		partner2.ProductIDs = []uuid.UUID{productB}
		partnerEncx2, err := domain.ProcessPartnerEncx(ctx, crypto, partner2)
		require.NoError(t, err)
		td.InsertPartnerEncx(t, ctx, partnerEncx2, testPool)
		require.NoError(t, err)

		// Partner with only product C (should NOT be returned)
		user3 := td.NewTestUser(t, "partnerC@example.com", "Partner", "C")
		user3.State = domain.Active
		user3.Role = identity.PartnerStr
		userEncx3, err := domain.ProcessUserEncx(ctx, crypto, user3)
		require.NoError(t, err)
		err = td.InsertUserEncx(t, ctx, userEncx3, testPool)
		require.NoError(t, err)

		partner3 := td.NewTestPartner(t, user3.ID)
		partner3.ProductIDs = []uuid.UUID{productC}
		partnerEncx3, err := domain.ProcessPartnerEncx(ctx, crypto, partner3)
		require.NoError(t, err)
		td.InsertPartnerEncx(t, ctx, partnerEncx3, testPool)
		require.NoError(t, err)

		// Act - request partners for products A and B only
		req := td.NewGetAllPartnersByProductsRequest(t, ctx, testServerURL, []uuid.UUID{productA, productB})
		resp, err := client.Do(req)

		// Assert
		assert.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		partners := td.ParsePartnersListResponse(t, resp)
		require.Len(t, partners, 2, "Expected 2 partners with products A or B")

		// Verify correct partners returned
		returnedIDs := make(map[uuid.UUID]bool)
		for _, p := range partners {
			returnedIDs[p.ID] = true
		}
		assert.True(t, returnedIDs[partner1.ID], "Partner 1 should be returned")
		assert.True(t, returnedIDs[partner2.ID], "Partner 2 should be returned")
		assert.False(t, returnedIDs[partner3.ID], "Partner 3 should NOT be returned")
	})

	t.Run("should return empty array when no partners match any product", func(t *testing.T) {
		// Clean state
		td.ClearPartnersTable(t, ctx, testPool)

		// Create partners with different products
		otherProduct := uuid.New()
		user := td.NewTestUser(t, "partner@example.com", "Partner", "Test")
		user.State = domain.Active
		user.Role = identity.PartnerStr
		userEncx, err := domain.ProcessUserEncx(ctx, crypto, user)
		require.NoError(t, err)
		err = td.InsertUserEncx(t, ctx, userEncx, testPool)
		require.NoError(t, err)

		partner := td.NewTestPartner(t, user.ID)
		partner.ProductIDs = []uuid.UUID{otherProduct}
		partnerEncx, err := domain.ProcessPartnerEncx(ctx, crypto, partner)
		require.NoError(t, err)
		td.InsertPartnerEncx(t, ctx, partnerEncx, testPool)
		require.NoError(t, err)

		// Act - request partners for products that don't exist
		nonExistentProducts := []uuid.UUID{uuid.New(), uuid.New()}
		req := td.NewGetAllPartnersByProductsRequest(t, ctx, testServerURL, nonExistentProducts)
		resp, err := client.Do(req)

		// Assert
		assert.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		partners := td.ParsePartnersListResponse(t, resp)
		assert.Empty(t, partners, "Expected empty partners array")
		assert.NotNil(t, partners, "Partners array should not be nil")
	})

	t.Run("should return 400 when no product IDs provided", func(t *testing.T) {
		// Clean state
		td.ClearPartnersTable(t, ctx, testPool)

		// Act - request with no product IDs
		req := td.NewGetAllPartnersByProductsRequest(t, ctx, testServerURL, []uuid.UUID{})
		resp, err := client.Do(req)

		// Assert
		assert.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should return 400 when product ID has invalid format", func(t *testing.T) {
		// Clean state
		td.ClearPartnersTable(t, ctx, testPool)

		// Act - request with invalid product ID
		req := td.NewGetAllPartnersByProductsRequestWithStrings(t, ctx, testServerURL, []string{"invalid-uuid", uuid.New().String()})
		resp, err := client.Do(req)

		// Assert - should return 400 bad request
		assert.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
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
		req := td.NewGetAllPartnersByProductsRequest(t, ctx, testServerURL, []uuid.UUID{productID})
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
		req := td.NewGetAllPartnersByProductsRequest(t, ctx, testServerURL, []uuid.UUID{productID})
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
		targetProduct1 := uuid.New()
		targetProduct2 := uuid.New()
		product1 := uuid.New()
		product2 := uuid.New()

		// Create partner with multiple products including target products
		user1 := td.NewTestUser(t, "multi1@example.com", "Multi", "One")
		user1.State = domain.Active
		user1.Role = identity.PartnerStr
		userEncx1, err := domain.ProcessUserEncx(ctx, crypto, user1)
		require.NoError(t, err)
		err = td.InsertUserEncx(t, ctx, userEncx1, testPool)
		require.NoError(t, err)

		partner1 := td.NewTestPartner(t, user1.ID)
		partner1.ProductIDs = []uuid.UUID{product1, targetProduct1, product2}
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
		partner2.ProductIDs = []uuid.UUID{targetProduct2}
		partnerEncx2, err := domain.ProcessPartnerEncx(ctx, crypto, partner2)
		require.NoError(t, err)
		td.InsertPartnerEncx(t, ctx, partnerEncx2, testPool)
		require.NoError(t, err)

		// Create partner without target products
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

		// Act - request partners for target products
		req := td.NewGetAllPartnersByProductsRequest(t, ctx, testServerURL, []uuid.UUID{targetProduct1, targetProduct2})
		resp, err := client.Do(req)

		// Assert HTTP response
		assert.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Parse response body - should return only partners with target products
		partners := td.ParsePartnersListResponse(t, resp)
		require.Len(t, partners, 2, "Expected 2 partners with the target products")

		// Verify both partners have at least one target product
		partnerIDs := []uuid.UUID{partner1.ID, partner2.ID}
		for _, responsePartner := range partners {
			assert.Contains(t, partnerIDs, responsePartner.ID, "Returned partner should be one of the expected partners")

			hasTargetProduct := false
			for _, prodID := range responsePartner.ProductIDs {
				if prodID == targetProduct1 || prodID == targetProduct2 {
					hasTargetProduct = true
					break
				}
			}
			assert.True(t, hasTargetProduct, "Partner should have at least one target product")
		}

		// Verify partner3 is NOT in the results
		for _, responsePartner := range partners {
			assert.NotEqual(t, partner3.ID, responsePartner.ID, "Partner without target products should not be returned")
		}
	})
}
