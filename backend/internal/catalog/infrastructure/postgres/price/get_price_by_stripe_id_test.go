package priceRepository_test

import (
	"context"
	"testing"
	"time"

	td "github.com/Leviosa-care/catalog/test/helpers"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/stretchr/testify/assert"
)

func TestGetPriceByStripeID(t *testing.T) {
	ctx := context.Background()

	t.Run("should successfully retrieve an existing price by Stripe ID", func(t *testing.T) {
		clearTables(t, ctx, testPool)

		// Setup: Create a parent product with a valid category
		parentCategory := td.NewValidCategory("Test Category")
		td.InsertCategory(t, ctx, parentCategory, testPool)
		parentProduct := td.NewValidProduct("Test Product", parentCategory.ID)
		td.InsertProduct(t, ctx, testPool, parentProduct)

		// Setup: Create and insert a mock price
		mockPrice := td.NewValidPrice()
		mockPrice.ProductID = parentProduct.ID
		mockPrice.StripePriceID = "price_12345"
		td.InsertPrice(t, ctx, mockPrice, testPool)

		// Action: Get the price by its StripePriceID
		retrievedPrice, err := repo.GetPriceByStripeID(ctx, mockPrice.StripePriceID)
		assert.NoError(t, err)
		assert.NotNil(t, retrievedPrice)

		// Assert that the retrieved data matches the mock data
		assert.Equal(t, mockPrice.ID, retrievedPrice.ID)
		assert.Equal(t, mockPrice.ProductID, retrievedPrice.ProductID)
		assert.Equal(t, mockPrice.StripePriceID, retrievedPrice.StripePriceID)
		assert.Equal(t, mockPrice.Amount, retrievedPrice.Amount)
		assert.Equal(t, mockPrice.Currency, retrievedPrice.Currency)
		assert.Equal(t, mockPrice.Interval, retrievedPrice.Interval)
		assert.Equal(t, mockPrice.IsActive, retrievedPrice.IsActive)
		assert.True(t, retrievedPrice.CreatedAt.Truncate(time.Second).Equal(mockPrice.CreatedAt.Truncate(time.Second)))
	})

	t.Run("should return a NotFound error if the Stripe ID does not exist", func(t *testing.T) {
		clearTables(t, ctx, testPool)

		nonExistentStripeID := "non_existent_price_id"

		// Action: Get a price that does not exist
		price, err := repo.GetPriceByStripeID(ctx, nonExistentStripeID)
		assert.Nil(t, price)
		assert.Error(t, err)

		// Assert the specific error type
		assert.ErrorIs(t, err, errs.ErrRepositoryNotFound, "expected not found error, got: %v", err)
	})

	t.Run("should return a NotFound error if the database is empty", func(t *testing.T) {
		clearTables(t, ctx, testPool)

		// Action: Get a price from an empty table
		price, err := repo.GetPriceByStripeID(ctx, "any_stripe_id")
		assert.Nil(t, price)
		assert.Error(t, err)

		// Assert the specific error type
		assert.ErrorIs(t, err, errs.ErrRepositoryNotFound, "expected not found error, got: %v", err)
	})
}
