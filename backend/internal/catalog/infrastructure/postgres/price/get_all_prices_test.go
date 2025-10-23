package priceRepository_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"
	td "github.com/Leviosa-care/leviosa/backend/test/helpers"

	"github.com/stretchr/testify/assert"
)

func TestGetAllPrices(t *testing.T) {
	ctx := context.Background()

	t.Run("should successfully retrieve all prices from the database, sorted oldest first", func(t *testing.T) {
		// --- Setup ---
		clearTables(t, ctx, testPool)

		// Create parent products to associate prices with.
		parentCategory := td.NewValidCategory("Test Category")
		td.InsertCategory(t, ctx, parentCategory, testPool)
		product1 := td.NewValidProduct("Test Product 1", parentCategory.ID)
		td.InsertProduct(t, ctx, testPool, product1)
		product2 := td.NewValidProduct("Test Product 2", parentCategory.ID)
		td.InsertProduct(t, ctx, testPool, product2)

		// Create prices for different catalog.
		price1 := td.NewValidPrice()
		price1.ProductID = product1.ID
		price1.CreatedAt = time.Now().Add(-3 * time.Hour) // Oldest

		price2 := td.NewValidPrice()
		price2.ProductID = product2.ID // Associated with a different product
		price2.StripePriceID = "price_xyz"
		price2.CreatedAt = time.Now().Add(-2 * time.Hour)

		price3 := td.NewValidPrice()
		price3.ProductID = product1.ID
		price3.StripePriceID = "price_abc"
		price3.CreatedAt = time.Now().Add(-1 * time.Hour) // Newest

		// Insert the prices into the database.
		td.InsertPrice(t, ctx, price1, testPool)
		td.InsertPrice(t, ctx, price2, testPool)
		td.InsertPrice(t, ctx, price3, testPool)

		// --- Action ---
		// Note: The `productID` parameter is unused in the provided function, so we pass a placeholder.
		retrievedPrices, err := repo.GetAllPrices(ctx)

		// --- Assertions ---
		assert.NoError(t, err)
		assert.Len(t, retrievedPrices, 3)

		// Assert that the prices are sorted by creation date, oldest first, as per the query.
		assert.Equal(t, price1.ID, retrievedPrices[0].ID)
		assert.Equal(t, price2.ID, retrievedPrices[1].ID)
		assert.Equal(t, price3.ID, retrievedPrices[2].ID)
	})

	t.Run("should return an empty slice if the database contains no prices", func(t *testing.T) {
		// --- Setup ---
		clearTables(t, ctx, testPool)

		// --- Action ---
		// No prices have been inserted into the database.
		retrievedPrices, err := repo.GetAllPrices(ctx)

		// --- Assertions ---
		assert.NoError(t, err)
		assert.Empty(t, retrievedPrices)
		// Ensure it's a non-nil, empty slice, as per the function's logic.
		assert.NotNil(t, retrievedPrices)
	})
	t.Run("should handle large number of prices efficiently", func(t *testing.T) {
		clearTables(t, ctx, testPool)

		// Create parent data
		parentCategory := td.NewValidCategory("Test Category")
		td.InsertCategory(t, ctx, parentCategory, testPool)
		product := td.NewValidProduct("Test Product", parentCategory.ID)
		td.InsertProduct(t, ctx, testPool, product)

		// Insert many prices (e.g., 1000)
		expectedCount := 1000
		for i := range expectedCount {
			price := td.NewValidPrice()
			price.ProductID = product.ID
			price.CreatedAt = time.Now().Add(time.Duration(-i) * time.Minute)
			price.StripePriceID = fmt.Sprintf("price_%d", i)
			td.InsertPrice(t, ctx, price, testPool)
		}

		prices, err := repo.GetAllPrices(ctx)

		assert.NoError(t, err)
		assert.Len(t, prices, expectedCount)
	})
	t.Run("should correctly convert interval strings to domain.Interval type", func(t *testing.T) {
		clearTables(t, ctx, testPool)

		// Setup parent data
		parentCategory := td.NewValidCategory("Test Category")
		td.InsertCategory(t, ctx, parentCategory, testPool)
		product := td.NewValidProduct("Test Product", parentCategory.ID)
		td.InsertProduct(t, ctx, testPool, product)

		// Create prices with different interval types
		monthlyPrice := td.NewValidPrice()
		monthlyPrice.ProductID = product.ID
		monthlyPrice.Interval = domain.Month // or however you represent this
		monthlyPrice.StripePriceID = "price_monthly"

		yearlyPrice := td.NewValidPrice()
		yearlyPrice.ProductID = product.ID
		yearlyPrice.Interval = domain.Year
		monthlyPrice.StripePriceID = "price_yearly"

		td.InsertPrice(t, ctx, monthlyPrice, testPool)
		td.InsertPrice(t, ctx, yearlyPrice, testPool)

		prices, err := repo.GetAllPrices(ctx)

		assert.NoError(t, err)
		assert.Len(t, prices, 2)

		// Verify the intervals were converted correctly
		intervalTypes := make(map[domain.Interval]bool)
		for _, price := range prices {
			intervalTypes[price.Interval] = true
		}
		assert.True(t, intervalTypes[domain.Month])
		assert.True(t, intervalTypes[domain.Year])
	})
}
