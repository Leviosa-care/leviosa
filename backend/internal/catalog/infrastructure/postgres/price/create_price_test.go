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

func TestCreatePrice(t *testing.T) {
	ctx := context.Background()

	t.Run("should successfully create a new price record", func(t *testing.T) {
		clearTables(t, ctx, testPool)

		category := td.NewValidCategory("Test category")
		td.InsertCategory(t, ctx, category, testPool)
		categoryID := category.ID

		// Setup: Create a parent product first to satisfy the foreign key constraint
		parentProduct := td.NewValidProduct("Test Product", categoryID)
		td.InsertProduct(t, ctx, testPool, parentProduct)

		mockPrice := td.NewValidPrice()
		mockPrice.ProductID = parentProduct.ID

		err := repo.CreatePrice(ctx, mockPrice)
		assert.NoError(t, err)

		// Verify the price was created and has the correct data
		createdPrice := td.GetPriceByID(t, ctx, mockPrice.ID, testPool)
		assert.NoError(t, err)

		// Assert that the retrieved data matches the mock data
		assert.Equal(t, mockPrice.ID, createdPrice.ID)
		assert.Equal(t, mockPrice.ProductID, createdPrice.ProductID)
		assert.Equal(t, mockPrice.StripePriceID, createdPrice.StripePriceID)
		assert.Equal(t, mockPrice.Amount, createdPrice.Amount)
		assert.Equal(t, mockPrice.Currency, createdPrice.Currency)
		assert.Equal(t, mockPrice.Interval, createdPrice.Interval)
		assert.Equal(t, mockPrice.IsActive, createdPrice.IsActive)
		assert.True(t, time.Since(createdPrice.CreatedAt) < time.Minute)
	})

	t.Run("should fail with a unique constraint error on StripePriceID", func(t *testing.T) {
		clearTables(t, ctx, testPool)

		category := td.NewValidCategory("Test category")
		td.InsertCategory(t, ctx, category, testPool)
		categoryID := category.ID

		// Setup: Create a parent product
		parentProduct := td.NewValidProduct("Test Product", categoryID)
		td.InsertProduct(t, ctx, testPool, parentProduct)

		mockPrice1 := td.NewValidPrice()
		mockPrice1.ProductID = parentProduct.ID

		// Create the first price
		err := repo.CreatePrice(ctx, mockPrice1)
		assert.NoError(t, err)

		// Attempt to create a second price with the same StripePriceID
		mockPrice2 := td.NewValidPrice()
		mockPrice2.ProductID = parentProduct.ID
		mockPrice2.StripePriceID = mockPrice1.StripePriceID // Duplicate ID

		err = repo.CreatePrice(ctx, mockPrice2)
		assert.Error(t, err)
		assert.ErrorIs(t, err, errs.ErrUniqueViolation)
	})

	t.Run("should fail with a foreign key constraint error if ProductID does not exist", func(t *testing.T) {
		clearTables(t, ctx, testPool)

		// Setup: Create a price with a non-existent ProductID
		nonExistentProductID := uuid.New()
		mockPrice := td.NewValidPrice()
		mockPrice.ProductID = nonExistentProductID

		err := repo.CreatePrice(ctx, mockPrice)
		assert.Error(t, err)
		assert.ErrorIs(t, err, errs.ErrForeignKeyViolation)
	})

	t.Run("should fail with a check constraint error for an invalid Interval", func(t *testing.T) {
		clearTables(t, ctx, testPool)

		category := td.NewValidCategory("Test category")
		td.InsertCategory(t, ctx, category, testPool)
		categoryID := category.ID

		// Setup: Create a parent product
		parentProduct := td.NewValidProduct("Test Product", categoryID)
		td.InsertProduct(t, ctx, testPool, parentProduct)

		// Setup: Create a price with an invalid interval type
		mockPrice := td.NewValidPrice()
		mockPrice.ProductID = parentProduct.ID
		mockPrice.Interval = "invalid_interval"

		err := repo.CreatePrice(ctx, mockPrice)
		assert.Error(t, err)
		assert.ErrorIs(t, err, errs.ErrInvalidInput)
	})
}
