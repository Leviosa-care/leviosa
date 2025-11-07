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

// make test-func TEST_NAME=TestGetAllPartnersByCategory TEST_PATH=test/integration/authuser/partner/get_all_partners_by_category_test.go

func TestGetAllPartnersByCategory(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("should successfully get all partners by category", func(t *testing.T) {
		// Clean state
		td.ClearPartnersTable(t, ctx, testPool)
		td.ClearSessionsRedis(t, ctx, redisClient)

		// Create a shared category ID that multiple partners will use
		sharedCategoryID := uuid.New()
		otherCategoryID := uuid.New()

		// Create 3 partners with the shared category
		partnersWithCategory := make([]*domain.Partner, 0, 3)
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
			// Assign the shared category to this partner
			partner.CategoryIDs = []uuid.UUID{sharedCategoryID}
			partnerEncx, err := domain.ProcessPartnerEncx(ctx, crypto, partner)
			require.NoError(t, err)
			td.InsertPartnerEncx(t, ctx, partnerEncx, testPool)
			require.NoError(t, err)
			partnersWithCategory = append(partnersWithCategory, partner)

			// Small delay to ensure different timestamps
			time.Sleep(10 * time.Millisecond)
		}

		// Create 2 partners with a different category (should NOT be returned)
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
			// Assign a different category
			partner.CategoryIDs = []uuid.UUID{otherCategoryID}
			partnerEncx, err := domain.ProcessPartnerEncx(ctx, crypto, partner)
			require.NoError(t, err)
			td.InsertPartnerEncx(t, ctx, partnerEncx, testPool)
			require.NoError(t, err)

			time.Sleep(10 * time.Millisecond)
		}

		// Act - request partners for the shared category
		req := td.NewGetAllPartnersByCategoryRequest(t, ctx, testServerURL, sharedCategoryID)
		resp, err := client.Do(req)

		// Assert HTTP response
		assert.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Parse response body
		partners := td.ParsePartnersListResponse(t, resp)
		require.Len(t, partners, 3, "Expected 3 partners in response for the shared category")

		// Verify partners are ordered by created_at DESC (newest first)
		for i := 0; i < len(partners)-1; i++ {
			assert.True(t, partners[i].CreatedAt.After(partners[i+1].CreatedAt) || partners[i].CreatedAt.Equal(partners[i+1].CreatedAt),
				"Partners should be ordered by created_at DESC")
		}

		// Verify each partner has the correct category
		for _, responsePartner := range partners {
			// Get encrypted partner from database by ID
			partnerEncx, err := td.GetPartnerEncxByID(t, ctx, responsePartner.ID, testPool)
			require.NoError(t, err, "Failed to get partner from database")

			// Decrypt partner
			dbPartner, err := domain.DecryptPartnerEncx(ctx, crypto, partnerEncx)
			require.NoError(t, err, "Failed to decrypt partner")

			// Verify the shared category is present in the partner's categories
			assert.Contains(t, dbPartner.CategoryIDs, sharedCategoryID, "Partner should have the shared category")
			assert.Equal(t, dbPartner.ID, responsePartner.ID, "ID mismatch")
			assert.Equal(t, dbPartner.Bio, responsePartner.Bio, "Bio mismatch")
			assert.Equal(t, dbPartner.Experience, responsePartner.Experience, "Experience mismatch")
		}
	})

	t.Run("should return empty array when no partners exist for category", func(t *testing.T) {
		// Clean state
		td.ClearPartnersTable(t, ctx, testPool)
		td.ClearSessionsRedis(t, ctx, redisClient)

		// Create partners with different categories
		otherCategoryID := uuid.New()
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
			partner.CategoryIDs = []uuid.UUID{otherCategoryID}
			partnerEncx, err := domain.ProcessPartnerEncx(ctx, crypto, partner)
			require.NoError(t, err)
			td.InsertPartnerEncx(t, ctx, partnerEncx, testPool)
			require.NoError(t, err)
		}

		// Act - request partners for a category that no partner has
		nonExistentCategoryID := uuid.New()
		req := td.NewGetAllPartnersByCategoryRequest(t, ctx, testServerURL, nonExistentCategoryID)
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

	t.Run("should return 400 when category ID is invalid", func(t *testing.T) {
		// Clean state
		td.ClearPartnersTable(t, ctx, testPool)

		// Act - request with invalid category ID
		req := td.NewGetAllPartnersByCategoryRequestWithInvalidID(t, ctx, testServerURL, "invalid-uuid")
		resp, err := client.Do(req)

		// Assert - should return 400 bad request
		assert.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should return 404 when category ID is empty", func(t *testing.T) {
		// Clean state
		td.ClearPartnersTable(t, ctx, testPool)

		// Act - request with empty category ID (URL: /partners/categories/)
		req := td.NewGetAllPartnersByCategoryRequestWithInvalidID(t, ctx, testServerURL, "")
		resp, err := client.Do(req)

		// Assert - should return 404 because route doesn't match (Go's http.ServeMux behavior)
		assert.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("should return 500 when partner DEK is corrupted", func(t *testing.T) {
		// Clean state
		td.ClearPartnersTable(t, ctx, testPool)

		// Create category and partner
		categoryID := uuid.New()
		user := td.NewTestUser(t, "partner@example.com", "John", "Partner")
		user.State = domain.Active
		user.Role = identity.PartnerStr
		userEncx, err := domain.ProcessUserEncx(ctx, crypto, user)
		require.NoError(t, err)
		err = td.InsertUserEncx(t, ctx, userEncx, testPool)
		require.NoError(t, err)

		partner := td.NewTestPartner(t, user.ID)
		partner.CategoryIDs = []uuid.UUID{categoryID}
		partnerEncx, err := domain.ProcessPartnerEncx(ctx, crypto, partner)
		require.NoError(t, err)
		err = td.InsertPartnerEncx(t, ctx, partnerEncx, testPool)
		require.NoError(t, err)

		// Corrupt the DEK to simulate decryption failure
		td.CorruptPartnerDEK(t, ctx, partner.ID, testPool)

		// Act
		req := td.NewGetAllPartnersByCategoryRequest(t, ctx, testServerURL, categoryID)
		resp, err := client.Do(req)

		// Assert - should return 500 due to decryption failure
		assert.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})

	t.Run("should return 500 when key version is invalid", func(t *testing.T) {
		// Clean state
		td.ClearPartnersTable(t, ctx, testPool)

		// Create category and partner
		categoryID := uuid.New()
		user := td.NewTestUser(t, "partner2@example.com", "Jane", "Partner")
		user.State = domain.Active
		user.Role = identity.PartnerStr
		userEncx, err := domain.ProcessUserEncx(ctx, crypto, user)
		require.NoError(t, err)
		err = td.InsertUserEncx(t, ctx, userEncx, testPool)
		require.NoError(t, err)

		partner := td.NewTestPartner(t, user.ID)
		partner.CategoryIDs = []uuid.UUID{categoryID}
		partnerEncx, err := domain.ProcessPartnerEncx(ctx, crypto, partner)
		require.NoError(t, err)
		err = td.InsertPartnerEncx(t, ctx, partnerEncx, testPool)
		require.NoError(t, err)

		// Set an invalid key version to simulate decryption failure
		td.SetInvalidKeyVersion(t, ctx, partner.ID, testPool, 99999)

		// Act
		req := td.NewGetAllPartnersByCategoryRequest(t, ctx, testServerURL, categoryID)
		resp, err := client.Do(req)

		// Assert - should return 500 due to decryption failure
		assert.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})

	t.Run("should handle partners with multiple categories correctly", func(t *testing.T) {
		// Clean state
		td.ClearPartnersTable(t, ctx, testPool)
		td.ClearSessionsRedis(t, ctx, redisClient)

		// Create target category and other categories
		targetCategoryID := uuid.New()
		category1 := uuid.New()
		category2 := uuid.New()

		// Create partner with multiple categories including the target
		user1 := td.NewTestUser(t, "multi1@example.com", "Multi", "One")
		user1.State = domain.Active
		user1.Role = identity.PartnerStr
		userEncx1, err := domain.ProcessUserEncx(ctx, crypto, user1)
		require.NoError(t, err)
		err = td.InsertUserEncx(t, ctx, userEncx1, testPool)
		require.NoError(t, err)

		partner1 := td.NewTestPartner(t, user1.ID)
		partner1.CategoryIDs = []uuid.UUID{category1, targetCategoryID, category2}
		partnerEncx1, err := domain.ProcessPartnerEncx(ctx, crypto, partner1)
		require.NoError(t, err)
		td.InsertPartnerEncx(t, ctx, partnerEncx1, testPool)
		require.NoError(t, err)

		// Create partner with only target category
		user2 := td.NewTestUser(t, "multi2@example.com", "Multi", "Two")
		user2.State = domain.Active
		user2.Role = identity.PartnerStr
		userEncx2, err := domain.ProcessUserEncx(ctx, crypto, user2)
		require.NoError(t, err)
		err = td.InsertUserEncx(t, ctx, userEncx2, testPool)
		require.NoError(t, err)

		partner2 := td.NewTestPartner(t, user2.ID)
		partner2.CategoryIDs = []uuid.UUID{targetCategoryID}
		partnerEncx2, err := domain.ProcessPartnerEncx(ctx, crypto, partner2)
		require.NoError(t, err)
		td.InsertPartnerEncx(t, ctx, partnerEncx2, testPool)
		require.NoError(t, err)

		// Create partner without target category
		user3 := td.NewTestUser(t, "multi3@example.com", "Multi", "Three")
		user3.State = domain.Active
		user3.Role = identity.PartnerStr
		userEncx3, err := domain.ProcessUserEncx(ctx, crypto, user3)
		require.NoError(t, err)
		err = td.InsertUserEncx(t, ctx, userEncx3, testPool)
		require.NoError(t, err)

		partner3 := td.NewTestPartner(t, user3.ID)
		partner3.CategoryIDs = []uuid.UUID{category1, category2}
		partnerEncx3, err := domain.ProcessPartnerEncx(ctx, crypto, partner3)
		require.NoError(t, err)
		td.InsertPartnerEncx(t, ctx, partnerEncx3, testPool)
		require.NoError(t, err)

		// Act - request partners for target category
		req := td.NewGetAllPartnersByCategoryRequest(t, ctx, testServerURL, targetCategoryID)
		resp, err := client.Do(req)

		// Assert HTTP response
		assert.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Parse response body - should return only partners with target category
		partners := td.ParsePartnersListResponse(t, resp)
		require.Len(t, partners, 2, "Expected 2 partners with the target category")

		// Verify both partners have the target category
		partnerIDs := []uuid.UUID{partner1.ID, partner2.ID}
		for _, responsePartner := range partners {
			assert.Contains(t, partnerIDs, responsePartner.ID, "Returned partner should be one of the expected partners")
			assert.Contains(t, responsePartner.CategoryIDs, targetCategoryID, "Partner should have target category")
		}

		// Verify partner3 is NOT in the results
		for _, responsePartner := range partners {
			assert.NotEqual(t, partner3.ID, responsePartner.ID, "Partner without target category should not be returned")
		}
	})
}
