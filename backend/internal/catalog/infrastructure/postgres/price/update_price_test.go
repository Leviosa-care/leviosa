package priceRepository_test

import (
	"context"
	"testing"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"
	td "github.com/Leviosa-care/leviosa/backend/test/helpers"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// UpdatePrice updates specific fields of a Price in the database.
// This function takes a domain.Price struct, assuming it's already validated
// and contains the fields to update.
func TestUpdatePrice(t *testing.T) {
	ctx := context.Background()

	t.Run("should successfully update an inactive price to be active", func(t *testing.T) {
		clearTables(t, ctx, testPool)

		// Setup: Create a product and an inactive price
		parentCategory := td.NewValidCategory("Test Category")
		td.InsertCategory(t, ctx, parentCategory, testPool)
		parentProduct := td.NewValidProduct("Test Product", parentCategory.ID)
		td.InsertProduct(t, ctx, testPool, parentProduct)

		mockPrice := td.NewValidPrice()
		mockPrice.ProductID = parentProduct.ID
		mockPrice.IsActive = false
		err := repo.CreatePrice(ctx, mockPrice)
		assert.NoError(t, err)

		// Assert initial state is inactive
		assert.False(t, td.GetPriceStatus(t, ctx, mockPrice.ID.String(), testPool))

		// Action: Update the price to be active
		isActive := true
		patch := &domain.UpdatePriceRequest{
			Active: &isActive,
		}
		err = repo.UpdatePrice(ctx, mockPrice.ID.String(), patch)
		assert.NoError(t, err)

		// Assert final state is active
		assert.True(t, td.GetPriceStatus(t, ctx, mockPrice.ID.String(), testPool))
	})

	t.Run("should successfully update an active price to be inactive", func(t *testing.T) {
		clearTables(t, ctx, testPool)

		// Setup: Create a product and an active price
		parentCategory := td.NewValidCategory("Test Category")
		td.InsertCategory(t, ctx, parentCategory, testPool)
		parentProduct := td.NewValidProduct("Test Product", parentCategory.ID)
		td.InsertProduct(t, ctx, testPool, parentProduct)

		mockPrice := td.NewValidPrice()
		mockPrice.ProductID = parentProduct.ID
		mockPrice.IsActive = true
		err := repo.CreatePrice(ctx, mockPrice)
		assert.NoError(t, err)

		// Assert initial state is active
		assert.True(t, td.GetPriceStatus(t, ctx, mockPrice.ID.String(), testPool))

		// Action: Update the price to be inactive
		isActive := false
		patch := &domain.UpdatePriceRequest{
			Active: &isActive,
		}
		err = repo.UpdatePrice(ctx, mockPrice.ID.String(), patch)
		assert.NoError(t, err)

		// Assert final state is inactive
		assert.False(t, td.GetPriceStatus(t, ctx, mockPrice.ID.String(), testPool))
	})

	t.Run("should return a NotFound error if the price ID does not exist", func(t *testing.T) {
		clearTables(t, ctx, testPool)

		nonExistentPriceID := uuid.New().String()
		isActive := true
		patch := &domain.UpdatePriceRequest{
			Active: &isActive,
		}

		// Action: Attempt to update a non-existent price
		err := repo.UpdatePrice(ctx, nonExistentPriceID, patch)

		// Assert the error
		assert.Error(t, err)
		assert.ErrorIs(t, err, errs.ErrRepositoryNotFound, "expected not found error, got: %v", err)
	})

	t.Run("should return an InvalidInput error if no updatable fields are provided", func(t *testing.T) {
		clearTables(t, ctx, testPool)

		// Action: Attempt to update a price with an empty patch
		priceID := uuid.New().String()
		patch := &domain.UpdatePriceRequest{} // Empty struct

		err := repo.UpdatePrice(ctx, priceID, patch)

		// Assert the error
		assert.Error(t, err)
		assert.ErrorIs(t, err, errs.ErrInvalidInput, "expected invalid input error, got: %v", err)
	})

	t.Run("should return a database error for an invalid price ID format", func(t *testing.T) {
		clearTables(t, ctx, testPool)

		invalidID := "not-a-valid-uuid"
		isActive := true
		patch := &domain.UpdatePriceRequest{
			Active: &isActive,
		}

		// Action: Attempt to update a price with an invalid ID
		err := repo.UpdatePrice(ctx, invalidID, patch)

		// Assert the error is a database query error
		assert.Error(t, err)
		// The error should be caught and classified by `ClassifyPgError`
		assert.ErrorIs(t, err, errs.ErrInvalidInput, "expected query failed error, got: %v", err)
	})
}
