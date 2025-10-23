package sharedRepository_test

import (
	"context"
	"testing"

	td "github.com/Leviosa-care/catalog/test/helpers"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestDeactivatePrices(t *testing.T) {
	ctx := context.Background()
	t.Run("should successfully deactivate a single active price", func(t *testing.T) {
		clearTables(t, ctx, testPool)

		// Setup: Create a parent product with a valid category
		parentCategory := td.NewValidCategory("Test Category")
		td.InsertCategory(t, ctx, parentCategory, testPool)
		parentProduct := td.NewValidProduct("Test Product", parentCategory.ID)
		td.InsertProduct(t, ctx, testPool, parentProduct)

		// Create an active price
		mockPrice := td.NewValidPrice()
		mockPrice.ProductID = parentProduct.ID
		mockPrice.IsActive = true
		td.InsertPrice(t, ctx, mockPrice, testPool)

		// Deactivate the price
		err := repo.DeactivatePrices(ctx, []string{mockPrice.ID.String()})
		assert.NoError(t, err)

		// Verify the price is now inactive
		deactivatedPrice := td.GetPriceByID(t, ctx, mockPrice.ID, testPool)
		assert.False(t, deactivatedPrice.IsActive)
	})

	t.Run("should successfully deactivate multiple active prices", func(t *testing.T) {
		clearTables(t, ctx, testPool)

		// Setup: Create a parent product with a valid category
		parentCategory := td.NewValidCategory("Test Category")
		td.InsertCategory(t, ctx, parentCategory, testPool)
		parentProduct := td.NewValidProduct("Test Product", parentCategory.ID)
		td.InsertProduct(t, ctx, testPool, parentProduct)

		// Create multiple active prices
		price1 := td.NewValidPrice()
		price1.ProductID = parentProduct.ID
		price1.IsActive = true
		price2 := td.NewValidPrice()
		price2.ProductID = parentProduct.ID
		price2.IsActive = true
		price2.StripePriceID = "price_xyz" // Ensure unique StripePriceID

		td.InsertPrice(t, ctx, price1, testPool)
		td.InsertPrice(t, ctx, price2, testPool)

		// Deactivate both prices
		err := repo.DeactivatePrices(ctx, []string{price1.ID.String(), price2.ID.String()})
		assert.NoError(t, err)

		// Verify both prices are now inactive
		deactivatedPrice1 := td.GetPriceByID(t, ctx, price1.ID, testPool)
		assert.False(t, deactivatedPrice1.IsActive)

		deactivatedPrice2 := td.GetPriceByID(t, ctx, price2.ID, testPool)
		assert.False(t, deactivatedPrice2.IsActive)
	})

	t.Run("should return no error if the price ID slice is empty", func(t *testing.T) {
		clearTables(t, ctx, testPool)

		// Deactivate an empty slice
		err := repo.DeactivatePrices(ctx, []string{})
		assert.NoError(t, err)
	})

	t.Run("should not return an error when deactivating an already inactive price", func(t *testing.T) {
		clearTables(t, ctx, testPool)

		// Setup: Create a parent product
		parentCategory := td.NewValidCategory("Test Category")
		td.InsertCategory(t, ctx, parentCategory, testPool)
		parentProduct := td.NewValidProduct("Test Product", parentCategory.ID)
		td.InsertProduct(t, ctx, testPool, parentProduct)

		// Create an inactive price
		mockPrice := td.NewValidPrice()
		mockPrice.ProductID = parentProduct.ID
		mockPrice.IsActive = false
		td.InsertPrice(t, ctx, mockPrice, testPool)

		// Deactivate the already inactive price
		err := repo.DeactivatePrices(ctx, []string{mockPrice.ID.String()})
		assert.NoError(t, err)

		// Verify the price remains inactive
		deactivatedPrice := td.GetPriceByID(t, ctx, mockPrice.ID, testPool)
		assert.False(t, deactivatedPrice.IsActive)
	})

	t.Run("should return a NotFound error if no prices with the given IDs are found", func(t *testing.T) {
		clearTables(t, ctx, testPool)

		nonExistentID := uuid.New().String()

		err := repo.DeactivatePrices(ctx, []string{nonExistentID})
		assert.Error(t, err)
		assert.ErrorIs(t, err, errs.ErrRepositoryNotFound)
	})

	t.Run("should deactivate existing prices and not return an error if some IDs are not found", func(t *testing.T) {
		clearTables(t, ctx, testPool)

		// Setup: Create a parent product
		parentCategory := td.NewValidCategory("Test Category")
		td.InsertCategory(t, ctx, parentCategory, testPool)
		parentProduct := td.NewValidProduct("Test Product", parentCategory.ID)
		td.InsertProduct(t, ctx, testPool, parentProduct)

		// Create one active price
		priceToDeactivate := td.NewValidPrice()
		priceToDeactivate.ProductID = parentProduct.ID
		priceToDeactivate.IsActive = true
		td.InsertPrice(t, ctx, priceToDeactivate, testPool)

		// Setup: Create a non-existent ID
		nonExistentID := uuid.New().String()

		// Deactivate a mix of valid and invalid IDs
		err := repo.DeactivatePrices(ctx, []string{priceToDeactivate.ID.String(), nonExistentID})
		assert.NoError(t, err) // The function is not expected to return an error here

		// Verify the existing price was deactivated
		deactivatedPrice := td.GetPriceByID(t, ctx, priceToDeactivate.ID, testPool)
		assert.False(t, deactivatedPrice.IsActive)
	})

	t.Run("should handle a database error gracefully", func(t *testing.T) {
		clearTables(t, ctx, testPool)

		// Provide an invalid UUID format in the slice to trigger a database error
		err := repo.DeactivatePrices(ctx, []string{"not-a-valid-uuid"})
		assert.Error(t, err)
		assert.ErrorIs(t, err, errs.ErrInvalidInput, "expected a query failure error")
	})
}

// Helper function for clearing tables before each test (optional but recommended) ---
func clearTables(t *testing.T, ctx context.Context, testPool *pgxpool.Pool) {
	td.ClearCategoriesTable(t, ctx, testPool)
	td.ClearProductsTable(t, ctx, testPool)
	td.ClearPricesTable(t, ctx, testPool)
}
