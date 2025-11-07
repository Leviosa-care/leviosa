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

// make test-func TEST_NAME=TestGetAllPartnersByCategories TEST_PATH=test/integration/authuser/partner/get_all_partners_by_categories_test.go

func TestGetAllPartnersByCategories(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("should successfully get all partners by multiple categories", func(t *testing.T) {
		// Clean state
		td.ClearPartnersTable(t, ctx, testPool)
		td.ClearSessionsRedis(t, ctx, redisClient)

		// Create 3 different categories
		category1 := uuid.New()
		category2 := uuid.New()
		category3 := uuid.New()

		// Create partner 1 with categories 1 and 2
		user1 := td.NewTestUser(t, "partner1@example.com", "Partner", "One")
		user1.State = domain.Active
		user1.Role = identity.PartnerStr
		userEncx1, err := domain.ProcessUserEncx(ctx, crypto, user1)
		require.NoError(t, err)
		err = td.InsertUserEncx(t, ctx, userEncx1, testPool)
		require.NoError(t, err)

		partner1 := td.NewTestPartner(t, user1.ID)
		partner1.CategoryIDs = []uuid.UUID{category1, category2}
		partnerEncx1, err := domain.ProcessPartnerEncx(ctx, crypto, partner1)
		require.NoError(t, err)
		td.InsertPartnerEncx(t, ctx, partnerEncx1, testPool)
		require.NoError(t, err)
		time.Sleep(10 * time.Millisecond)

		// Create partner 2 with categories 2 and 3
		user2 := td.NewTestUser(t, "partner2@example.com", "Partner", "Two")
		user2.State = domain.Active
		user2.Role = identity.PartnerStr
		userEncx2, err := domain.ProcessUserEncx(ctx, crypto, user2)
		require.NoError(t, err)
		err = td.InsertUserEncx(t, ctx, userEncx2, testPool)
		require.NoError(t, err)

		partner2 := td.NewTestPartner(t, user2.ID)
		partner2.CategoryIDs = []uuid.UUID{category2, category3}
		partnerEncx2, err := domain.ProcessPartnerEncx(ctx, crypto, partner2)
		require.NoError(t, err)
		td.InsertPartnerEncx(t, ctx, partnerEncx2, testPool)
		require.NoError(t, err)
		time.Sleep(10 * time.Millisecond)

		// Create partner 3 with only category 1
		user3 := td.NewTestUser(t, "partner3@example.com", "Partner", "Three")
		user3.State = domain.Active
		user3.Role = identity.PartnerStr
		userEncx3, err := domain.ProcessUserEncx(ctx, crypto, user3)
		require.NoError(t, err)
		err = td.InsertUserEncx(t, ctx, userEncx3, testPool)
		require.NoError(t, err)

		partner3 := td.NewTestPartner(t, user3.ID)
		partner3.CategoryIDs = []uuid.UUID{category1}
		partnerEncx3, err := domain.ProcessPartnerEncx(ctx, crypto, partner3)
		require.NoError(t, err)
		td.InsertPartnerEncx(t, ctx, partnerEncx3, testPool)
		require.NoError(t, err)

		// Act - request partners for categories 1 and 3 (should return partners 1, 2, and 3)
		req := td.NewGetAllPartnersByCategoriesRequest(t, ctx, testServerURL, []uuid.UUID{category1, category3})
		resp, err := client.Do(req)

		// Assert HTTP response
		assert.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Parse response body
		partners := td.ParsePartnersListResponse(t, resp)
		require.Len(t, partners, 3, "Expected 3 partners matching categories 1 or 3")

		// Verify partners are ordered by created_at DESC (newest first)
		for i := 0; i < len(partners)-1; i++ {
			assert.True(t, partners[i].CreatedAt.After(partners[i+1].CreatedAt) || partners[i].CreatedAt.Equal(partners[i+1].CreatedAt),
				"Partners should be ordered by created_at DESC")
		}

		// Verify all returned partners have at least one of the requested categories
		partnerIDs := map[uuid.UUID]bool{partner1.ID: true, partner2.ID: true, partner3.ID: true}
		for _, responsePartner := range partners {
			assert.True(t, partnerIDs[responsePartner.ID], "Returned partner should be one of the expected partners")

			// Verify partner has at least one of the requested categories
			hasCategory := false
			for _, catID := range responsePartner.CategoryIDs {
				if catID == category1 || catID == category3 {
					hasCategory = true
					break
				}
			}
			assert.True(t, hasCategory, "Partner should have at least one of the requested categories")
		}
	})

	t.Run("should return partners matching any of the provided categories", func(t *testing.T) {
		// Clean state
		td.ClearPartnersTable(t, ctx, testPool)

		// Create categories
		categoryA := uuid.New()
		categoryB := uuid.New()
		categoryC := uuid.New()

		// Partner with only category A
		user1 := td.NewTestUser(t, "partnerA@example.com", "Partner", "A")
		user1.State = domain.Active
		user1.Role = identity.PartnerStr
		userEncx1, err := domain.ProcessUserEncx(ctx, crypto, user1)
		require.NoError(t, err)
		err = td.InsertUserEncx(t, ctx, userEncx1, testPool)
		require.NoError(t, err)

		partner1 := td.NewTestPartner(t, user1.ID)
		partner1.CategoryIDs = []uuid.UUID{categoryA}
		partnerEncx1, err := domain.ProcessPartnerEncx(ctx, crypto, partner1)
		require.NoError(t, err)
		td.InsertPartnerEncx(t, ctx, partnerEncx1, testPool)
		require.NoError(t, err)

		// Partner with only category B
		user2 := td.NewTestUser(t, "partnerB@example.com", "Partner", "B")
		user2.State = domain.Active
		user2.Role = identity.PartnerStr
		userEncx2, err := domain.ProcessUserEncx(ctx, crypto, user2)
		require.NoError(t, err)
		err = td.InsertUserEncx(t, ctx, userEncx2, testPool)
		require.NoError(t, err)

		partner2 := td.NewTestPartner(t, user2.ID)
		partner2.CategoryIDs = []uuid.UUID{categoryB}
		partnerEncx2, err := domain.ProcessPartnerEncx(ctx, crypto, partner2)
		require.NoError(t, err)
		td.InsertPartnerEncx(t, ctx, partnerEncx2, testPool)
		require.NoError(t, err)

		// Partner with only category C (should NOT be returned)
		user3 := td.NewTestUser(t, "partnerC@example.com", "Partner", "C")
		user3.State = domain.Active
		user3.Role = identity.PartnerStr
		userEncx3, err := domain.ProcessUserEncx(ctx, crypto, user3)
		require.NoError(t, err)
		err = td.InsertUserEncx(t, ctx, userEncx3, testPool)
		require.NoError(t, err)

		partner3 := td.NewTestPartner(t, user3.ID)
		partner3.CategoryIDs = []uuid.UUID{categoryC}
		partnerEncx3, err := domain.ProcessPartnerEncx(ctx, crypto, partner3)
		require.NoError(t, err)
		td.InsertPartnerEncx(t, ctx, partnerEncx3, testPool)
		require.NoError(t, err)

		// Act - request partners for categories A and B only
		req := td.NewGetAllPartnersByCategoriesRequest(t, ctx, testServerURL, []uuid.UUID{categoryA, categoryB})
		resp, err := client.Do(req)

		// Assert
		assert.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		partners := td.ParsePartnersListResponse(t, resp)
		require.Len(t, partners, 2, "Expected 2 partners with categories A or B")

		// Verify correct partners returned
		returnedIDs := make(map[uuid.UUID]bool)
		for _, p := range partners {
			returnedIDs[p.ID] = true
		}
		assert.True(t, returnedIDs[partner1.ID], "Partner 1 should be returned")
		assert.True(t, returnedIDs[partner2.ID], "Partner 2 should be returned")
		assert.False(t, returnedIDs[partner3.ID], "Partner 3 should NOT be returned")
	})

	t.Run("should return empty array when no partners match any category", func(t *testing.T) {
		// Clean state
		td.ClearPartnersTable(t, ctx, testPool)

		// Create partners with different categories
		otherCategory := uuid.New()
		user := td.NewTestUser(t, "partner@example.com", "Partner", "Test")
		user.State = domain.Active
		user.Role = identity.PartnerStr
		userEncx, err := domain.ProcessUserEncx(ctx, crypto, user)
		require.NoError(t, err)
		err = td.InsertUserEncx(t, ctx, userEncx, testPool)
		require.NoError(t, err)

		partner := td.NewTestPartner(t, user.ID)
		partner.CategoryIDs = []uuid.UUID{otherCategory}
		partnerEncx, err := domain.ProcessPartnerEncx(ctx, crypto, partner)
		require.NoError(t, err)
		td.InsertPartnerEncx(t, ctx, partnerEncx, testPool)
		require.NoError(t, err)

		// Act - request partners for categories that don't exist
		nonExistentCategories := []uuid.UUID{uuid.New(), uuid.New()}
		req := td.NewGetAllPartnersByCategoriesRequest(t, ctx, testServerURL, nonExistentCategories)
		resp, err := client.Do(req)

		// Assert
		assert.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		partners := td.ParsePartnersListResponse(t, resp)
		assert.Empty(t, partners, "Expected empty partners array")
		assert.NotNil(t, partners, "Partners array should not be nil")
	})

	t.Run("should return 400 when no category IDs provided", func(t *testing.T) {
		// Clean state
		td.ClearPartnersTable(t, ctx, testPool)

		// Act - request with no category IDs
		req := td.NewGetAllPartnersByCategoriesRequest(t, ctx, testServerURL, []uuid.UUID{})
		resp, err := client.Do(req)

		// Assert
		assert.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should return 400 when category ID has invalid format", func(t *testing.T) {
		// Clean state
		td.ClearPartnersTable(t, ctx, testPool)

		// Act - request with invalid UUID
		req := td.NewGetAllPartnersByCategoriesRequestWithStrings(t, ctx, testServerURL, []string{"invalid-uuid", uuid.New().String()})
		resp, err := client.Do(req)

		// Assert
		assert.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
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
		req := td.NewGetAllPartnersByCategoriesRequest(t, ctx, testServerURL, []uuid.UUID{categoryID})
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
		req := td.NewGetAllPartnersByCategoriesRequest(t, ctx, testServerURL, []uuid.UUID{categoryID})
		resp, err := client.Do(req)

		// Assert - should return 500 due to decryption failure
		assert.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})

	t.Run("should handle single category ID correctly", func(t *testing.T) {
		// Clean state
		td.ClearPartnersTable(t, ctx, testPool)

		// Create single category and partner
		categoryID := uuid.New()
		user := td.NewTestUser(t, "single@example.com", "Single", "Partner")
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
		td.InsertPartnerEncx(t, ctx, partnerEncx, testPool)
		require.NoError(t, err)

		// Act - request with single category
		req := td.NewGetAllPartnersByCategoriesRequest(t, ctx, testServerURL, []uuid.UUID{categoryID})
		resp, err := client.Do(req)

		// Assert
		assert.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		partners := td.ParsePartnersListResponse(t, resp)
		require.Len(t, partners, 1, "Expected 1 partner")
		assert.Equal(t, partner.ID, partners[0].ID)
	})
}
