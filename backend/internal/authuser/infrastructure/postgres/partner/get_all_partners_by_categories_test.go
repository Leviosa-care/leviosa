package partnerRepository_test

import (
	"context"
	"testing"

	td "github.com/Leviosa-care/leviosa/backend/test/helpers"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestGetAllPartnersByCategories TEST_PATH=internal/authuser/infrastructure/postgres/partner/get_all_partners_by_categories_test.go

func TestGetAllPartnersByCategories(t *testing.T) {
	ctx := context.Background()

	t.Run("should successfully retrieve partners by multiple categories", func(t *testing.T) {
		// Arrange
		td.ClearPartnersTable(t, ctx, testPool)

		category1 := uuid.New()
		category2 := uuid.New()
		category3 := uuid.New()
		otherCategory := uuid.New()

		// Create test users
		userID1 := td.CreateTestUserForPartnerWithUniqueEmail(t, ctx, testPool, "multi_cat1")
		userID2 := td.CreateTestUserForPartnerWithUniqueEmail(t, ctx, testPool, "multi_cat2")
		userID3 := td.CreateTestUserForPartnerWithUniqueEmail(t, ctx, testPool, "multi_cat3")
		userID4 := td.CreateTestUserForPartnerWithUniqueEmail(t, ctx, testPool, "multi_cat4")

		// Partner with category1
		partner1 := td.NewTestPartnerEncx(t)
		partner1.UserID = userID1
		partner1.CategoryIDs = []uuid.UUID{category1}
		err := td.InsertPartnerEncx(t, ctx, partner1, testPool)
		require.NoError(t, err)

		// Partner with category2
		partner2 := td.NewTestPartnerEncx(t)
		partner2.UserID = userID2
		partner2.CategoryIDs = []uuid.UUID{category2}
		err = td.InsertPartnerEncx(t, ctx, partner2, testPool)
		require.NoError(t, err)

		// Partner with category1 and category3
		partner3 := td.NewTestPartnerEncx(t)
		partner3.UserID = userID3
		partner3.CategoryIDs = []uuid.UUID{category1, category3}
		err = td.InsertPartnerEncx(t, ctx, partner3, testPool)
		require.NoError(t, err)

		// Partner with only otherCategory (should not be returned)
		partner4 := td.NewTestPartnerEncx(t)
		partner4.UserID = userID4
		partner4.CategoryIDs = []uuid.UUID{otherCategory}
		err = td.InsertPartnerEncx(t, ctx, partner4, testPool)
		require.NoError(t, err)

		// Act - search for partners with category1 or category2
		partners, err := repo.GetAllPartnersByCategories(ctx, []uuid.UUID{category1, category2})

		// Assert
		assert.NoError(t, err)
		assert.Len(t, partners, 3, "Should return partners with any matching category")

		// Verify correct partners are returned
		userIDs := make([]uuid.UUID, len(partners))
		for i, p := range partners {
			userIDs[i] = p.UserID
		}
		assert.Contains(t, userIDs, userID1, "Should include partner with category1")
		assert.Contains(t, userIDs, userID2, "Should include partner with category2")
		assert.Contains(t, userIDs, userID3, "Should include partner with category1 and category3")
		assert.NotContains(t, userIDs, userID4, "Should not include partner with only otherCategory")
	})

	t.Run("should return empty slice when no partners match any category", func(t *testing.T) {
		// Arrange
		td.ClearPartnersTable(t, ctx, testPool)

		nonMatchingCategory1 := uuid.New()
		nonMatchingCategory2 := uuid.New()

		// Create partner with different categories
		userID := td.CreateTestUserForPartnerWithUniqueEmail(t, ctx, testPool, "no_match")
		partner := td.NewTestPartnerEncx(t)
		partner.UserID = userID
		partner.CategoryIDs = []uuid.UUID{uuid.New(), uuid.New()}
		err := td.InsertPartnerEncx(t, ctx, partner, testPool)
		require.NoError(t, err)

		// Act
		partners, err := repo.GetAllPartnersByCategories(ctx, []uuid.UUID{nonMatchingCategory1, nonMatchingCategory2})

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, partners, "Should return non-nil slice")
		assert.Empty(t, partners, "Should return empty slice when no matches")
	})

	t.Run("should return empty slice when category list is empty", func(t *testing.T) {
		// Arrange
		td.ClearPartnersTable(t, ctx, testPool)

		// Create some partners
		userID := td.CreateTestUserForPartnerWithUniqueEmail(t, ctx, testPool, "empty_search")
		partner := td.NewTestPartnerEncx(t)
		partner.UserID = userID
		partner.CategoryIDs = []uuid.UUID{uuid.New()}
		err := td.InsertPartnerEncx(t, ctx, partner, testPool)
		require.NoError(t, err)

		// Act - search with empty category list
		partners, err := repo.GetAllPartnersByCategories(ctx, []uuid.UUID{})

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, partners, "Should return non-nil slice")
		assert.Empty(t, partners, "Should return empty slice when search list is empty")
	})

	t.Run("should return empty slice when no partners exist", func(t *testing.T) {
		// Arrange
		td.ClearPartnersTable(t, ctx, testPool)

		categories := []uuid.UUID{uuid.New(), uuid.New()}

		// Act
		partners, err := repo.GetAllPartnersByCategories(ctx, categories)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, partners, "Should return non-nil slice")
		assert.Empty(t, partners, "Should return empty slice when no partners exist")
	})

	t.Run("should order results by created_at DESC", func(t *testing.T) {
		// Arrange
		td.ClearPartnersTable(t, ctx, testPool)

		category := uuid.New()

		// Create partners at different times
		userID1 := td.CreateTestUserForPartnerWithUniqueEmail(t, ctx, testPool, "order1")
		partner1 := td.NewTestPartnerEncx(t)
		partner1.UserID = userID1
		partner1.CategoryIDs = []uuid.UUID{category}
		err := td.InsertPartnerEncx(t, ctx, partner1, testPool)
		require.NoError(t, err)

		userID2 := td.CreateTestUserForPartnerWithUniqueEmail(t, ctx, testPool, "order2")
		partner2 := td.NewTestPartnerEncx(t)
		partner2.UserID = userID2
		partner2.CategoryIDs = []uuid.UUID{category}
		err = td.InsertPartnerEncx(t, ctx, partner2, testPool)
		require.NoError(t, err)

		userID3 := td.CreateTestUserForPartnerWithUniqueEmail(t, ctx, testPool, "order3")
		partner3 := td.NewTestPartnerEncx(t)
		partner3.UserID = userID3
		partner3.CategoryIDs = []uuid.UUID{category}
		err = td.InsertPartnerEncx(t, ctx, partner3, testPool)
		require.NoError(t, err)

		// Act
		partners, err := repo.GetAllPartnersByCategories(ctx, []uuid.UUID{category})

		// Assert
		assert.NoError(t, err)
		assert.Len(t, partners, 3)

		// Verify newest partner is first
		assert.Equal(t, userID3, partners[0].UserID, "Most recently created partner should be first")
		assert.Equal(t, userID2, partners[1].UserID, "Second partner should be in the middle")
		assert.Equal(t, userID1, partners[2].UserID, "First created partner should be last")
	})

	t.Run("should not return duplicates when partner matches multiple search categories", func(t *testing.T) {
		// Arrange
		td.ClearPartnersTable(t, ctx, testPool)

		category1 := uuid.New()
		category2 := uuid.New()
		category3 := uuid.New()

		// Partner with multiple categories that match search
		userID := td.CreateTestUserForPartnerWithUniqueEmail(t, ctx, testPool, "multi_match")
		partner := td.NewTestPartnerEncx(t)
		partner.UserID = userID
		partner.CategoryIDs = []uuid.UUID{category1, category2, category3}
		err := td.InsertPartnerEncx(t, ctx, partner, testPool)
		require.NoError(t, err)

		// Act - search with categories that all exist in partner
		partners, err := repo.GetAllPartnersByCategories(ctx, []uuid.UUID{category1, category2})

		// Assert
		assert.NoError(t, err)
		assert.Len(t, partners, 1, "Should return partner only once even if it matches multiple search categories")
		assert.Equal(t, userID, partners[0].UserID)
	})

	t.Run("should handle large category search lists", func(t *testing.T) {
		// Arrange
		td.ClearPartnersTable(t, ctx, testPool)

		// Create 20 different categories
		categories := make([]uuid.UUID, 20)
		for i := range categories {
			categories[i] = uuid.New()
		}

		// Create partners with different category combinations
		userID1 := td.CreateTestUserForPartnerWithUniqueEmail(t, ctx, testPool, "large1")
		partner1 := td.NewTestPartnerEncx(t)
		partner1.UserID = userID1
		partner1.CategoryIDs = []uuid.UUID{categories[0]}
		err := td.InsertPartnerEncx(t, ctx, partner1, testPool)
		require.NoError(t, err)

		userID2 := td.CreateTestUserForPartnerWithUniqueEmail(t, ctx, testPool, "large2")
		partner2 := td.NewTestPartnerEncx(t)
		partner2.UserID = userID2
		partner2.CategoryIDs = []uuid.UUID{categories[10]}
		err = td.InsertPartnerEncx(t, ctx, partner2, testPool)
		require.NoError(t, err)

		userID3 := td.CreateTestUserForPartnerWithUniqueEmail(t, ctx, testPool, "large3")
		partner3 := td.NewTestPartnerEncx(t)
		partner3.UserID = userID3
		partner3.CategoryIDs = []uuid.UUID{categories[19]}
		err = td.InsertPartnerEncx(t, ctx, partner3, testPool)
		require.NoError(t, err)

		// Partner with non-matching category
		userID4 := td.CreateTestUserForPartnerWithUniqueEmail(t, ctx, testPool, "large4")
		partner4 := td.NewTestPartnerEncx(t)
		partner4.UserID = userID4
		partner4.CategoryIDs = []uuid.UUID{uuid.New()}
		err = td.InsertPartnerEncx(t, ctx, partner4, testPool)
		require.NoError(t, err)

		// Act - search with all 20 categories
		partners, err := repo.GetAllPartnersByCategories(ctx, categories)

		// Assert
		assert.NoError(t, err)
		assert.Len(t, partners, 3, "Should return 3 partners matching any of the 20 categories")

		userIDs := make([]uuid.UUID, len(partners))
		for i, p := range partners {
			userIDs[i] = p.UserID
		}
		assert.Contains(t, userIDs, userID1)
		assert.Contains(t, userIDs, userID2)
		assert.Contains(t, userIDs, userID3)
		assert.NotContains(t, userIDs, userID4, "Should not include partner with non-matching category")
	})

	t.Run("should handle partners with empty category arrays", func(t *testing.T) {
		// Arrange
		td.ClearPartnersTable(t, ctx, testPool)

		categories := []uuid.UUID{uuid.New(), uuid.New()}

		// Partner with empty categories
		userID1 := td.CreateTestUserForPartnerWithUniqueEmail(t, ctx, testPool, "empty_cat")
		partner1 := td.NewTestPartnerEncx(t)
		partner1.UserID = userID1
		partner1.CategoryIDs = []uuid.UUID{}
		err := td.InsertPartnerEncx(t, ctx, partner1, testPool)
		require.NoError(t, err)

		// Partner with matching category
		userID2 := td.CreateTestUserForPartnerWithUniqueEmail(t, ctx, testPool, "with_cat")
		partner2 := td.NewTestPartnerEncx(t)
		partner2.UserID = userID2
		partner2.CategoryIDs = []uuid.UUID{categories[0]}
		err = td.InsertPartnerEncx(t, ctx, partner2, testPool)
		require.NoError(t, err)

		// Act
		partners, err := repo.GetAllPartnersByCategories(ctx, categories)

		// Assert
		assert.NoError(t, err)
		assert.Len(t, partners, 1, "Should only return partner with matching category")
		assert.Equal(t, userID2, partners[0].UserID)
	})

	t.Run("should return all partner response fields", func(t *testing.T) {
		// Arrange
		td.ClearPartnersTable(t, ctx, testPool)

		category := uuid.New()
		productID := uuid.New()

		userID := td.CreateTestUserForPartnerWithUniqueEmail(t, ctx, testPool, "fields")
		partner := td.NewTestPartnerEncx(t)
		partner.UserID = userID
		partner.CategoryIDs = []uuid.UUID{category}
		partner.ProductIDs = []uuid.UUID{productID}
		err := td.InsertPartnerEncx(t, ctx, partner, testPool)
		require.NoError(t, err)

		// Act
		partners, err := repo.GetAllPartnersByCategories(ctx, []uuid.UUID{category})

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
