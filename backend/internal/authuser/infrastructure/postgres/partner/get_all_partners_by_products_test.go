package partnerRepository_test

import (
	"context"
	"testing"

	td "github.com/Leviosa-care/leviosa/backend/test/helpers"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestGetAllPartnersByProducts TEST_PATH=internal/authuser/infrastructure/postgres/partner/get_all_partners_by_products_test.go

func TestGetAllPartnersByProducts(t *testing.T) {
	ctx := context.Background()

	t.Run("should successfully retrieve partners by multiple products", func(t *testing.T) {
		// Arrange
		td.ClearPartnersTable(t, ctx, testPool)

		product1 := uuid.New()
		product2 := uuid.New()
		product3 := uuid.New()
		otherProduct := uuid.New()

		// Create test users
		userID1 := td.CreateTestUserForPartnerWithUniqueEmail(t, ctx, testPool, "multi_prod1")
		userID2 := td.CreateTestUserForPartnerWithUniqueEmail(t, ctx, testPool, "multi_prod2")
		userID3 := td.CreateTestUserForPartnerWithUniqueEmail(t, ctx, testPool, "multi_prod3")
		userID4 := td.CreateTestUserForPartnerWithUniqueEmail(t, ctx, testPool, "multi_prod4")

		// Partner with product1
		partner1 := td.NewTestPartnerEncx(t)
		partner1.UserID = userID1
		partner1.ProductIDs = []uuid.UUID{product1}
		err := td.InsertPartnerEncx(t, ctx, partner1, testPool)
		require.NoError(t, err)

		// Partner with product2
		partner2 := td.NewTestPartnerEncx(t)
		partner2.UserID = userID2
		partner2.ProductIDs = []uuid.UUID{product2}
		err = td.InsertPartnerEncx(t, ctx, partner2, testPool)
		require.NoError(t, err)

		// Partner with product1 and product3
		partner3 := td.NewTestPartnerEncx(t)
		partner3.UserID = userID3
		partner3.ProductIDs = []uuid.UUID{product1, product3}
		err = td.InsertPartnerEncx(t, ctx, partner3, testPool)
		require.NoError(t, err)

		// Partner with only otherProduct (should not be returned)
		partner4 := td.NewTestPartnerEncx(t)
		partner4.UserID = userID4
		partner4.ProductIDs = []uuid.UUID{otherProduct}
		err = td.InsertPartnerEncx(t, ctx, partner4, testPool)
		require.NoError(t, err)

		// Act - search for partners with product1 or product2
		partners, err := repo.GetAllPartnersByProducts(ctx, []uuid.UUID{product1, product2})

		// Assert
		assert.NoError(t, err)
		assert.Len(t, partners, 3, "Should return partners with any matching product")

		// Verify correct partners are returned
		userIDs := make([]uuid.UUID, len(partners))
		for i, p := range partners {
			userIDs[i] = p.UserID
		}
		assert.Contains(t, userIDs, userID1, "Should include partner with product1")
		assert.Contains(t, userIDs, userID2, "Should include partner with product2")
		assert.Contains(t, userIDs, userID3, "Should include partner with product1 and product3")
		assert.NotContains(t, userIDs, userID4, "Should not include partner with only otherProduct")
	})

	t.Run("should return empty slice when no partners match any product", func(t *testing.T) {
		// Arrange
		td.ClearPartnersTable(t, ctx, testPool)

		nonMatchingProduct1 := uuid.New()
		nonMatchingProduct2 := uuid.New()

		// Create partner with different products
		userID := td.CreateTestUserForPartnerWithUniqueEmail(t, ctx, testPool, "no_match")
		partner := td.NewTestPartnerEncx(t)
		partner.UserID = userID
		partner.ProductIDs = []uuid.UUID{uuid.New(), uuid.New()}
		err := td.InsertPartnerEncx(t, ctx, partner, testPool)
		require.NoError(t, err)

		// Act
		partners, err := repo.GetAllPartnersByProducts(ctx, []uuid.UUID{nonMatchingProduct1, nonMatchingProduct2})

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, partners, "Should return non-nil slice")
		assert.Empty(t, partners, "Should return empty slice when no matches")
	})

	t.Run("should return empty slice when product list is empty", func(t *testing.T) {
		// Arrange
		td.ClearPartnersTable(t, ctx, testPool)

		// Create some partners
		userID := td.CreateTestUserForPartnerWithUniqueEmail(t, ctx, testPool, "empty_search")
		partner := td.NewTestPartnerEncx(t)
		partner.UserID = userID
		partner.ProductIDs = []uuid.UUID{uuid.New()}
		err := td.InsertPartnerEncx(t, ctx, partner, testPool)
		require.NoError(t, err)

		// Act - search with empty product list
		partners, err := repo.GetAllPartnersByProducts(ctx, []uuid.UUID{})

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, partners, "Should return non-nil slice")
		assert.Empty(t, partners, "Should return empty slice when search list is empty")
	})

	t.Run("should return empty slice when no partners exist", func(t *testing.T) {
		// Arrange
		td.ClearPartnersTable(t, ctx, testPool)

		products := []uuid.UUID{uuid.New(), uuid.New()}

		// Act
		partners, err := repo.GetAllPartnersByProducts(ctx, products)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, partners, "Should return non-nil slice")
		assert.Empty(t, partners, "Should return empty slice when no partners exist")
	})

	t.Run("should order results by created_at DESC", func(t *testing.T) {
		// Arrange
		td.ClearPartnersTable(t, ctx, testPool)

		product := uuid.New()

		// Create partners at different times
		userID1 := td.CreateTestUserForPartnerWithUniqueEmail(t, ctx, testPool, "order1")
		partner1 := td.NewTestPartnerEncx(t)
		partner1.UserID = userID1
		partner1.ProductIDs = []uuid.UUID{product}
		err := td.InsertPartnerEncx(t, ctx, partner1, testPool)
		require.NoError(t, err)

		userID2 := td.CreateTestUserForPartnerWithUniqueEmail(t, ctx, testPool, "order2")
		partner2 := td.NewTestPartnerEncx(t)
		partner2.UserID = userID2
		partner2.ProductIDs = []uuid.UUID{product}
		err = td.InsertPartnerEncx(t, ctx, partner2, testPool)
		require.NoError(t, err)

		userID3 := td.CreateTestUserForPartnerWithUniqueEmail(t, ctx, testPool, "order3")
		partner3 := td.NewTestPartnerEncx(t)
		partner3.UserID = userID3
		partner3.ProductIDs = []uuid.UUID{product}
		err = td.InsertPartnerEncx(t, ctx, partner3, testPool)
		require.NoError(t, err)

		// Act
		partners, err := repo.GetAllPartnersByProducts(ctx, []uuid.UUID{product})

		// Assert
		assert.NoError(t, err)
		assert.Len(t, partners, 3)

		// Verify newest partner is first
		assert.Equal(t, userID3, partners[0].UserID, "Most recently created partner should be first")
		assert.Equal(t, userID2, partners[1].UserID, "Second partner should be in the middle")
		assert.Equal(t, userID1, partners[2].UserID, "First created partner should be last")
	})

	t.Run("should not return duplicates when partner matches multiple search products", func(t *testing.T) {
		// Arrange
		td.ClearPartnersTable(t, ctx, testPool)

		product1 := uuid.New()
		product2 := uuid.New()
		product3 := uuid.New()

		// Partner with multiple products that match search
		userID := td.CreateTestUserForPartnerWithUniqueEmail(t, ctx, testPool, "multi_match")
		partner := td.NewTestPartnerEncx(t)
		partner.UserID = userID
		partner.ProductIDs = []uuid.UUID{product1, product2, product3}
		err := td.InsertPartnerEncx(t, ctx, partner, testPool)
		require.NoError(t, err)

		// Act - search with products that all exist in partner
		partners, err := repo.GetAllPartnersByProducts(ctx, []uuid.UUID{product1, product2})

		// Assert
		assert.NoError(t, err)
		assert.Len(t, partners, 1, "Should return partner only once even if it matches multiple search products")
		assert.Equal(t, userID, partners[0].UserID)
	})

	t.Run("should handle large product search lists", func(t *testing.T) {
		// Arrange
		td.ClearPartnersTable(t, ctx, testPool)

		// Create 20 different products
		products := make([]uuid.UUID, 20)
		for i := range products {
			products[i] = uuid.New()
		}

		// Create partners with different product combinations
		userID1 := td.CreateTestUserForPartnerWithUniqueEmail(t, ctx, testPool, "large1")
		partner1 := td.NewTestPartnerEncx(t)
		partner1.UserID = userID1
		partner1.ProductIDs = []uuid.UUID{products[0]}
		err := td.InsertPartnerEncx(t, ctx, partner1, testPool)
		require.NoError(t, err)

		userID2 := td.CreateTestUserForPartnerWithUniqueEmail(t, ctx, testPool, "large2")
		partner2 := td.NewTestPartnerEncx(t)
		partner2.UserID = userID2
		partner2.ProductIDs = []uuid.UUID{products[10]}
		err = td.InsertPartnerEncx(t, ctx, partner2, testPool)
		require.NoError(t, err)

		userID3 := td.CreateTestUserForPartnerWithUniqueEmail(t, ctx, testPool, "large3")
		partner3 := td.NewTestPartnerEncx(t)
		partner3.UserID = userID3
		partner3.ProductIDs = []uuid.UUID{products[19]}
		err = td.InsertPartnerEncx(t, ctx, partner3, testPool)
		require.NoError(t, err)

		// Partner with non-matching product
		userID4 := td.CreateTestUserForPartnerWithUniqueEmail(t, ctx, testPool, "large4")
		partner4 := td.NewTestPartnerEncx(t)
		partner4.UserID = userID4
		partner4.ProductIDs = []uuid.UUID{uuid.New()}
		err = td.InsertPartnerEncx(t, ctx, partner4, testPool)
		require.NoError(t, err)

		// Act - search with all 20 products
		partners, err := repo.GetAllPartnersByProducts(ctx, products)

		// Assert
		assert.NoError(t, err)
		assert.Len(t, partners, 3, "Should return 3 partners matching any of the 20 products")

		userIDs := make([]uuid.UUID, len(partners))
		for i, p := range partners {
			userIDs[i] = p.UserID
		}
		assert.Contains(t, userIDs, userID1)
		assert.Contains(t, userIDs, userID2)
		assert.Contains(t, userIDs, userID3)
		assert.NotContains(t, userIDs, userID4, "Should not include partner with non-matching product")
	})

	t.Run("should handle partners with empty product arrays", func(t *testing.T) {
		// Arrange
		td.ClearPartnersTable(t, ctx, testPool)

		products := []uuid.UUID{uuid.New(), uuid.New()}

		// Partner with empty products
		userID1 := td.CreateTestUserForPartnerWithUniqueEmail(t, ctx, testPool, "empty_prod")
		partner1 := td.NewTestPartnerEncx(t)
		partner1.UserID = userID1
		partner1.ProductIDs = []uuid.UUID{}
		err := td.InsertPartnerEncx(t, ctx, partner1, testPool)
		require.NoError(t, err)

		// Partner with matching product
		userID2 := td.CreateTestUserForPartnerWithUniqueEmail(t, ctx, testPool, "with_prod")
		partner2 := td.NewTestPartnerEncx(t)
		partner2.UserID = userID2
		partner2.ProductIDs = []uuid.UUID{products[0]}
		err = td.InsertPartnerEncx(t, ctx, partner2, testPool)
		require.NoError(t, err)

		// Act
		partners, err := repo.GetAllPartnersByProducts(ctx, products)

		// Assert
		assert.NoError(t, err)
		assert.Len(t, partners, 1, "Should only return partner with matching product")
		assert.Equal(t, userID2, partners[0].UserID)
	})

	t.Run("should return all partner response fields", func(t *testing.T) {
		// Arrange
		td.ClearPartnersTable(t, ctx, testPool)

		categoryID := uuid.New()
		product := uuid.New()

		userID := td.CreateTestUserForPartnerWithUniqueEmail(t, ctx, testPool, "fields")
		partner := td.NewTestPartnerEncx(t)
		partner.UserID = userID
		partner.CategoryIDs = []uuid.UUID{categoryID}
		partner.ProductIDs = []uuid.UUID{product}
		err := td.InsertPartnerEncx(t, ctx, partner, testPool)
		require.NoError(t, err)

		// Act
		partners, err := repo.GetAllPartnersByProducts(ctx, []uuid.UUID{product})

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
