package partnerRepository_test

import (
	"context"
	"testing"

	td "github.com/Leviosa-care/leviosa/backend/test/helpers"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestGetAllPartnersByCategory TEST_PATH=internal/authuser/infrastructure/postgres/partner/get_all_partners_by_category_test.go

func TestGetAllPartnersByCategory(t *testing.T) {
	ctx := context.Background()

	t.Run("should successfully retrieve partners by single category", func(t *testing.T) {
		// Arrange
		td.ClearPartnersTable(t, ctx, testPool)

		targetCategoryID := uuid.New()
		otherCategoryID := uuid.New()

		// Create test users
		userID1 := td.CreateTestUserForPartnerWithUniqueEmail(t, ctx, testPool, "cat1_user1")
		userID2 := td.CreateTestUserForPartnerWithUniqueEmail(t, ctx, testPool, "cat1_user2")
		userID3 := td.CreateTestUserForPartnerWithUniqueEmail(t, ctx, testPool, "cat1_user3")

		// Partner with target category
		partner1 := td.NewTestPartnerEncx(t)
		partner1.UserID = userID1
		partner1.CategoryIDs = []uuid.UUID{targetCategoryID}
		err := td.InsertPartnerEncx(t, ctx, partner1, testPool)
		require.NoError(t, err)

		// Partner with target category + other categories
		partner2 := td.NewTestPartnerEncx(t)
		partner2.UserID = userID2
		partner2.CategoryIDs = []uuid.UUID{targetCategoryID, otherCategoryID}
		err = td.InsertPartnerEncx(t, ctx, partner2, testPool)
		require.NoError(t, err)

		// Partner without target category (should not be returned)
		partner3 := td.NewTestPartnerEncx(t)
		partner3.UserID = userID3
		partner3.CategoryIDs = []uuid.UUID{otherCategoryID}
		err = td.InsertPartnerEncx(t, ctx, partner3, testPool)
		require.NoError(t, err)

		// Act
		partners, err := repo.GetAllPartnersByCategory(ctx, targetCategoryID)

		// Assert
		assert.NoError(t, err)
		assert.Len(t, partners, 2, "Should return only partners with target category")

		// Verify correct partners are returned
		userIDs := make([]uuid.UUID, len(partners))
		for i, p := range partners {
			userIDs[i] = p.UserID
		}
		assert.Contains(t, userIDs, userID1, "Should include partner1")
		assert.Contains(t, userIDs, userID2, "Should include partner2")
		assert.NotContains(t, userIDs, userID3, "Should not include partner3")

		// Verify category IDs are returned (decrypted at service layer, not repository)
		for _, p := range partners {
			assert.NotEmpty(t, p.CategoryIDs, "Category IDs should be populated")
			assert.Contains(t, p.CategoryIDs, targetCategoryID, "Should contain target category")
		}
	})

	t.Run("should return empty slice when no partners have the category", func(t *testing.T) {
		// Arrange
		td.ClearPartnersTable(t, ctx, testPool)

		nonExistentCategoryID := uuid.New()

		// Create partner with different category
		userID := td.CreateTestUserForPartnerWithUniqueEmail(t, ctx, testPool, "no_match")
		partner := td.NewTestPartnerEncx(t)
		partner.UserID = userID
		partner.CategoryIDs = []uuid.UUID{uuid.New()}
		err := td.InsertPartnerEncx(t, ctx, partner, testPool)
		require.NoError(t, err)

		// Act
		partners, err := repo.GetAllPartnersByCategory(ctx, nonExistentCategoryID)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, partners, "Should return non-nil slice")
		assert.Empty(t, partners, "Should return empty slice when no matches")
	})

	t.Run("should return empty slice when no partners exist", func(t *testing.T) {
		// Arrange
		td.ClearPartnersTable(t, ctx, testPool)

		categoryID := uuid.New()

		// Act
		partners, err := repo.GetAllPartnersByCategory(ctx, categoryID)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, partners, "Should return non-nil slice")
		assert.Empty(t, partners, "Should return empty slice when no partners exist")
	})

	t.Run("should order results by created_at DESC", func(t *testing.T) {
		// Arrange
		td.ClearPartnersTable(t, ctx, testPool)

		categoryID := uuid.New()

		// Create partners at different times
		userID1 := td.CreateTestUserForPartnerWithUniqueEmail(t, ctx, testPool, "order1")
		partner1 := td.NewTestPartnerEncx(t)
		partner1.UserID = userID1
		partner1.CategoryIDs = []uuid.UUID{categoryID}
		err := td.InsertPartnerEncx(t, ctx, partner1, testPool)
		require.NoError(t, err)

		userID2 := td.CreateTestUserForPartnerWithUniqueEmail(t, ctx, testPool, "order2")
		partner2 := td.NewTestPartnerEncx(t)
		partner2.UserID = userID2
		partner2.CategoryIDs = []uuid.UUID{categoryID}
		err = td.InsertPartnerEncx(t, ctx, partner2, testPool)
		require.NoError(t, err)

		userID3 := td.CreateTestUserForPartnerWithUniqueEmail(t, ctx, testPool, "order3")
		partner3 := td.NewTestPartnerEncx(t)
		partner3.UserID = userID3
		partner3.CategoryIDs = []uuid.UUID{categoryID}
		err = td.InsertPartnerEncx(t, ctx, partner3, testPool)
		require.NoError(t, err)

		// Act
		partners, err := repo.GetAllPartnersByCategory(ctx, categoryID)

		// Assert
		assert.NoError(t, err)
		assert.Len(t, partners, 3)

		// Verify newest partner is first
		assert.Equal(t, userID3, partners[0].UserID, "Most recently created partner should be first")
		assert.Equal(t, userID2, partners[1].UserID, "Second partner should be in the middle")
		assert.Equal(t, userID1, partners[2].UserID, "First created partner should be last")
	})

	t.Run("should handle partners with empty category arrays", func(t *testing.T) {
		// Arrange
		td.ClearPartnersTable(t, ctx, testPool)

		categoryID := uuid.New()

		// Partner with empty categories
		userID1 := td.CreateTestUserForPartnerWithUniqueEmail(t, ctx, testPool, "empty_cat")
		partner1 := td.NewTestPartnerEncx(t)
		partner1.UserID = userID1
		partner1.CategoryIDs = []uuid.UUID{}
		err := td.InsertPartnerEncx(t, ctx, partner1, testPool)
		require.NoError(t, err)

		// Partner with the target category
		userID2 := td.CreateTestUserForPartnerWithUniqueEmail(t, ctx, testPool, "with_cat")
		partner2 := td.NewTestPartnerEncx(t)
		partner2.UserID = userID2
		partner2.CategoryIDs = []uuid.UUID{categoryID}
		err = td.InsertPartnerEncx(t, ctx, partner2, testPool)
		require.NoError(t, err)

		// Act
		partners, err := repo.GetAllPartnersByCategory(ctx, categoryID)

		// Assert
		assert.NoError(t, err)
		assert.Len(t, partners, 1, "Should only return partner with category")
		assert.Equal(t, userID2, partners[0].UserID, "Should return partner with target category")
	})

	t.Run("should return all partner response fields", func(t *testing.T) {
		// Arrange
		td.ClearPartnersTable(t, ctx, testPool)

		categoryID := uuid.New()
		productID := uuid.New()

		userID := td.CreateTestUserForPartnerWithUniqueEmail(t, ctx, testPool, "fields")
		partner := td.NewTestPartnerEncx(t)
		partner.UserID = userID
		partner.CategoryIDs = []uuid.UUID{categoryID}
		partner.ProductIDs = []uuid.UUID{productID}
		err := td.InsertPartnerEncx(t, ctx, partner, testPool)
		require.NoError(t, err)

		// Act
		partners, err := repo.GetAllPartnersByCategory(ctx, categoryID)

		// Assert
		require.NoError(t, err)
		require.Len(t, partners, 1)

		p := partners[0]
		assert.Equal(t, userID, p.UserID)
		assert.NotEmpty(t, p.Bio, "Bio should be populated")
		assert.NotEmpty(t, p.Experience, "Experience should be populated")
		assert.NotEmpty(t, p.CategoryIDs, "CategoryIDs should be populated")
		assert.NotEmpty(t, p.ProductIDs, "ProductIDs should be populated")
		assert.NotZero(t, p.CreatedAt, "CreatedAt should be populated")
		assert.NotZero(t, p.UpdatedAt, "UpdatedAt should be populated")
	})
}
