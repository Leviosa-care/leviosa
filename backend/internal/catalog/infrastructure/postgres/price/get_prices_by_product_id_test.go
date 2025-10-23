package priceRepository_test

import (
	"context"
	"testing"
	"time"

	td "github.com/Leviosa-care/catalog/test/helpers"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetPricesByProductID(t *testing.T) {
	ctx := context.Background()

	t.Run("should successfully retrieve all prices for a product, sorted newest first", func(t *testing.T) {
		clearTables(t, ctx, testPool)

		// Setup: Create a parent product
		parentCategory := td.NewValidCategory("Test Category")
		td.InsertCategory(t, ctx, parentCategory, testPool)
		parentProduct := td.NewValidProduct("Test Product", parentCategory.ID)
		td.InsertProduct(t, ctx, testPool, parentProduct)

		// Create three prices for the product
		price1 := td.NewValidPrice()
		price1.ProductID = parentProduct.ID
		price1.CreatedAt = time.Now().Add(-3 * time.Hour)
		price2 := td.NewValidPrice()
		price2.ProductID = parentProduct.ID
		price2.StripePriceID = "price_xyz" // Ensure unique StripePriceID
		price2.CreatedAt = time.Now().Add(-2 * time.Hour)
		price3 := td.NewValidPrice()
		price3.ProductID = parentProduct.ID
		price3.StripePriceID = "price_abc" // Ensure unique StripePriceID
		price3.CreatedAt = time.Now().Add(-1 * time.Hour)

		td.InsertPrice(t, ctx, price1, testPool)
		td.InsertPrice(t, ctx, price2, testPool)
		td.InsertPrice(t, ctx, price3, testPool)

		// Create a price for another product to ensure filtering works
		otherProduct := td.NewValidProduct("Other Product", parentCategory.ID)
		td.InsertProduct(t, ctx, testPool, otherProduct)
		otherPrice := td.NewValidPrice()
		otherPrice.ProductID = otherProduct.ID
		otherPrice.StripePriceID = "price_def"
		td.InsertPrice(t, ctx, otherPrice, testPool)

		// Action: List all prices for the main product
		retrievedPrices, err := repo.GetPricesByProductID(ctx, parentProduct.ID.String(), false)
		assert.NoError(t, err)
		assert.Len(t, retrievedPrices, 3)

		// Assert that the prices are sorted by creation date, newest first
		assert.Equal(t, price3.ID, retrievedPrices[0].ID)
		assert.Equal(t, price2.ID, retrievedPrices[1].ID)
		assert.Equal(t, price1.ID, retrievedPrices[2].ID)
	})

	t.Run("should successfully retrieve only active prices for a product", func(t *testing.T) {
		clearTables(t, ctx, testPool)

		// Setup: Create a parent product
		parentCategory := td.NewValidCategory("Test Category")
		td.InsertCategory(t, ctx, parentCategory, testPool)
		parentProduct := td.NewValidProduct("Test Product", parentCategory.ID)
		td.InsertProduct(t, ctx, testPool, parentProduct)

		// Create one active and two inactive prices for the product
		inactivePrice1 := td.NewValidPrice()
		inactivePrice1.ProductID = parentProduct.ID
		inactivePrice1.StripePriceID = "price_def"
		inactivePrice1.IsActive = false
		inactivePrice1.CreatedAt = time.Now().Add(-3 * time.Hour)

		activePrice := td.NewValidPrice()
		activePrice.ProductID = parentProduct.ID
		activePrice.ProductID = parentProduct.ID
		activePrice.StripePriceID = "price_ghi"
		activePrice.IsActive = true
		activePrice.CreatedAt = time.Now().Add(-2 * time.Hour)

		inactivePrice2 := td.NewValidPrice()
		inactivePrice2.ProductID = parentProduct.ID
		inactivePrice2.StripePriceID = "price_jkl"
		inactivePrice2.IsActive = false
		inactivePrice2.CreatedAt = time.Now().Add(-1 * time.Hour)

		err := repo.CreatePrice(ctx, inactivePrice1)
		assert.NoError(t, err)
		err = repo.CreatePrice(ctx, activePrice)
		assert.NoError(t, err)
		err = repo.CreatePrice(ctx, inactivePrice2)
		assert.NoError(t, err)

		// Action: List only active prices for the product
		retrievedPrices, err := repo.GetPricesByProductID(ctx, parentProduct.ID.String(), true)
		assert.NoError(t, err)
		assert.Len(t, retrievedPrices, 1)

		// Assert that the retrieved price is the active one
		assert.Equal(t, activePrice.ID, retrievedPrices[0].ID)
	})

	t.Run("should return an empty slice if the product has no prices", func(t *testing.T) {
		clearTables(t, ctx, testPool)

		// Setup: Create a parent product but no prices
		parentCategory := td.NewValidCategory("Test Category")
		td.InsertCategory(t, ctx, parentCategory, testPool)
		parentProduct := td.NewValidProduct("Test Product", parentCategory.ID)
		td.InsertProduct(t, ctx, testPool, parentProduct)

		// Action: List prices for a product with no prices
		retrievedPrices, err := repo.GetPricesByProductID(ctx, parentProduct.ID.String(), false)
		assert.NoError(t, err)
		assert.Empty(t, retrievedPrices)
		assert.NotNil(t, retrievedPrices) // Should return an empty slice, not nil
	})

	t.Run("should return an empty slice if the product ID does not exist", func(t *testing.T) {
		clearTables(t, ctx, testPool)

		nonExistentProductID := uuid.New().String()

		// Action: List prices for a non-existent product
		retrievedPrices, err := repo.GetPricesByProductID(ctx, nonExistentProductID, false)
		assert.NoError(t, err)
		assert.Empty(t, retrievedPrices)
		assert.NotNil(t, retrievedPrices)
	})

	t.Run("should return a QueryFailed error for an invalid product ID format", func(t *testing.T) {
		clearTables(t, ctx, testPool)

		invalidID := "not-a-valid-uuid"

		// Action: List prices with an invalid product ID
		retrievedPrices, err := repo.GetPricesByProductID(ctx, invalidID, false)
		assert.Error(t, err)
		assert.ErrorIs(t, err, errs.ErrDBQuery, "expected query failed error, got: %v", err)
		assert.Nil(t, retrievedPrices)
	})
}
