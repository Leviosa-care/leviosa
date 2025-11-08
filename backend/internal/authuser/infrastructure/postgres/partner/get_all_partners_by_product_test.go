package partnerRepository_test

import (
	"context"
	"testing"

	td "github.com/Leviosa-care/leviosa/backend/test/helpers"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestGetAllPartnersByProduct TEST_PATH=internal/authuser/infrastructure/postgres/partner/get_all_partners_by_product_test.go

func TestGetAllPartnersByProduct(t *testing.T) {
	ctx := context.Background()

	t.Run("should successfully retrieve partners by single product", func(t *testing.T) {
		// Arrange
		td.ClearPartnersTable(t, ctx, testPool)

		targetProductID := uuid.New()
		otherProductID := uuid.New()

		// Create test users
		userID1 := td.CreateTestUserForPartnerWithUniqueEmail(t, ctx, testPool, "prod1_user1")
		userID2 := td.CreateTestUserForPartnerWithUniqueEmail(t, ctx, testPool, "prod1_user2")
		userID3 := td.CreateTestUserForPartnerWithUniqueEmail(t, ctx, testPool, "prod1_user3")

		// Partner with target product
		partner1 := td.NewTestPartnerEncx(t)
		partner1.UserID = userID1
		partner1.ProductIDs = []uuid.UUID{targetProductID}
		err := td.InsertPartnerEncx(t, ctx, partner1, testPool)
		require.NoError(t, err)

		// Partner with target product + other products
		partner2 := td.NewTestPartnerEncx(t)
		partner2.UserID = userID2
		partner2.ProductIDs = []uuid.UUID{targetProductID, otherProductID}
		err = td.InsertPartnerEncx(t, ctx, partner2, testPool)
		require.NoError(t, err)

		// Partner without target product (should not be returned)
		partner3 := td.NewTestPartnerEncx(t)
		partner3.UserID = userID3
		partner3.ProductIDs = []uuid.UUID{otherProductID}
		err = td.InsertPartnerEncx(t, ctx, partner3, testPool)
		require.NoError(t, err)

		// Act
		partners, err := repo.GetAllPartnersByProduct(ctx, targetProductID)

		// Assert
		assert.NoError(t, err)
		assert.Len(t, partners, 2, "Should return only partners with target product")

		// Verify correct partners are returned
		userIDs := make([]uuid.UUID, len(partners))
		for i, p := range partners {
			userIDs[i] = p.UserID
		}
		assert.Contains(t, userIDs, userID1, "Should include partner1")
		assert.Contains(t, userIDs, userID2, "Should include partner2")
		assert.NotContains(t, userIDs, userID3, "Should not include partner3")

		// Verify product IDs are returned (decrypted at service layer, not repository)
		for _, p := range partners {
			assert.NotEmpty(t, p.ProductIDs, "Product IDs should be populated")
			assert.Contains(t, p.ProductIDs, targetProductID, "Should contain target product")
		}
	})

	t.Run("should return empty slice when no partners have the product", func(t *testing.T) {
		// Arrange
		td.ClearPartnersTable(t, ctx, testPool)

		nonExistentProductID := uuid.New()

		// Create partner with different product
		userID := td.CreateTestUserForPartnerWithUniqueEmail(t, ctx, testPool, "no_match")
		partner := td.NewTestPartnerEncx(t)
		partner.UserID = userID
		partner.ProductIDs = []uuid.UUID{uuid.New()}
		err := td.InsertPartnerEncx(t, ctx, partner, testPool)
		require.NoError(t, err)

		// Act
		partners, err := repo.GetAllPartnersByProduct(ctx, nonExistentProductID)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, partners, "Should return non-nil slice")
		assert.Empty(t, partners, "Should return empty slice when no matches")
	})

	t.Run("should return empty slice when no partners exist", func(t *testing.T) {
		// Arrange
		td.ClearPartnersTable(t, ctx, testPool)

		productID := uuid.New()

		// Act
		partners, err := repo.GetAllPartnersByProduct(ctx, productID)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, partners, "Should return non-nil slice")
		assert.Empty(t, partners, "Should return empty slice when no partners exist")
	})

	t.Run("should order results by created_at DESC", func(t *testing.T) {
		// Arrange
		td.ClearPartnersTable(t, ctx, testPool)

		productID := uuid.New()

		// Create partners at different times
		userID1 := td.CreateTestUserForPartnerWithUniqueEmail(t, ctx, testPool, "order1")
		partner1 := td.NewTestPartnerEncx(t)
		partner1.UserID = userID1
		partner1.ProductIDs = []uuid.UUID{productID}
		err := td.InsertPartnerEncx(t, ctx, partner1, testPool)
		require.NoError(t, err)

		userID2 := td.CreateTestUserForPartnerWithUniqueEmail(t, ctx, testPool, "order2")
		partner2 := td.NewTestPartnerEncx(t)
		partner2.UserID = userID2
		partner2.ProductIDs = []uuid.UUID{productID}
		err = td.InsertPartnerEncx(t, ctx, partner2, testPool)
		require.NoError(t, err)

		userID3 := td.CreateTestUserForPartnerWithUniqueEmail(t, ctx, testPool, "order3")
		partner3 := td.NewTestPartnerEncx(t)
		partner3.UserID = userID3
		partner3.ProductIDs = []uuid.UUID{productID}
		err = td.InsertPartnerEncx(t, ctx, partner3, testPool)
		require.NoError(t, err)

		// Act
		partners, err := repo.GetAllPartnersByProduct(ctx, productID)

		// Assert
		assert.NoError(t, err)
		assert.Len(t, partners, 3)

		// Verify newest partner is first
		assert.Equal(t, userID3, partners[0].UserID, "Most recently created partner should be first")
		assert.Equal(t, userID2, partners[1].UserID, "Second partner should be in the middle")
		assert.Equal(t, userID1, partners[2].UserID, "First created partner should be last")
	})

	t.Run("should handle partners with empty product arrays", func(t *testing.T) {
		// Arrange
		td.ClearPartnersTable(t, ctx, testPool)

		productID := uuid.New()

		// Partner with empty products
		userID1 := td.CreateTestUserForPartnerWithUniqueEmail(t, ctx, testPool, "empty_prod")
		partner1 := td.NewTestPartnerEncx(t)
		partner1.UserID = userID1
		partner1.ProductIDs = []uuid.UUID{}
		err := td.InsertPartnerEncx(t, ctx, partner1, testPool)
		require.NoError(t, err)

		// Partner with the target product
		userID2 := td.CreateTestUserForPartnerWithUniqueEmail(t, ctx, testPool, "with_prod")
		partner2 := td.NewTestPartnerEncx(t)
		partner2.UserID = userID2
		partner2.ProductIDs = []uuid.UUID{productID}
		err = td.InsertPartnerEncx(t, ctx, partner2, testPool)
		require.NoError(t, err)

		// Act
		partners, err := repo.GetAllPartnersByProduct(ctx, productID)

		// Assert
		assert.NoError(t, err)
		assert.Len(t, partners, 1, "Should only return partner with product")
		assert.Equal(t, userID2, partners[0].UserID, "Should return partner with target product")
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
		partners, err := repo.GetAllPartnersByProduct(ctx, productID)

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
