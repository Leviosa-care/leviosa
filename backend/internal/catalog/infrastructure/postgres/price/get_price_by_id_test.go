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

func TestGetPrice(t *testing.T) {
	ctx := context.Background()

	t.Run("should successfully retrieve an existing price", func(t *testing.T) {
		clearTables(t, ctx, testPool)

		// Setup: Create a parent product with a valid category
		parentCategory := td.NewValidCategory("Test Category")
		td.InsertCategory(t, ctx, parentCategory, testPool)
		parentProduct := td.NewValidProduct("Test Product", parentCategory.ID)
		td.InsertProduct(t, ctx, testPool, parentProduct)

		// Setup: Create and insert a mock price
		mockPrice := td.NewValidPrice()
		mockPrice.ProductID = parentProduct.ID
		td.InsertPrice(t, ctx, mockPrice, testPool)

		// Action: Get the price by its ID
		retrievedPrice := td.GetPriceByID(t, ctx, mockPrice.ID, testPool)

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

	t.Run("should return a NotFound error if the price does not exist", func(t *testing.T) {
		clearTables(t, ctx, testPool)

		nonExistentID := uuid.New().String()

		// Action: Get a price that does not exist
		price, err := repo.GetPrice(ctx, nonExistentID)
		assert.Nil(t, price)
		assert.Error(t, err)

		// Assert the specific error type
		assert.ErrorIs(t, err, errs.ErrRepositoryNotFound, "expected not found error, got: %v", err)
	})

	t.Run("should return a ErrInvalidInput error for an invalid UUID format", func(t *testing.T) {
		clearTables(t, ctx, testPool)

		invalidID := "not-a-valid-uuid"

		// Action: Get a price with an invalid ID
		price, err := repo.GetPrice(ctx, invalidID)
		assert.Nil(t, price)
		assert.Error(t, err)

		// Assert the specific error type
		assert.ErrorIs(t, err, errs.ErrInvalidInput, "expected invalid input error, got: %v", err)
	})
}
