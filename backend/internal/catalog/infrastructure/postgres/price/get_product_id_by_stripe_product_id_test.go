package priceRepository_test

import (
	"context"
	"testing"

	td "github.com/Leviosa-care/leviosa/backend/test/helpers"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/stretchr/testify/assert"
)

func TestGetProductIDByStripeProductID(t *testing.T) {
	ctx := context.Background()

	t.Run("should successfully retrieve an internal product ID by a valid Stripe product ID", func(t *testing.T) {
		clearTables(t, ctx, testPool)

		// Setup: Create a parent category
		parentCategory := td.NewValidCategory("Test Category")
		td.InsertCategory(t, ctx, parentCategory, testPool)

		// Setup: Create a parent product with a known Stripe ID
		mockProduct := td.NewValidProduct("Test Product", parentCategory.ID)
		mockProduct.StripeProductID = "prod_123456789"
		td.InsertProduct(t, ctx, testPool, mockProduct)

		// Action: Retrieve the internal product ID
		retrievedID, err := repo.GetProductIDByStripeProductID(ctx, mockProduct.StripeProductID)
		assert.NoError(t, err)
		assert.Equal(t, mockProduct.ID.String(), retrievedID)
	})

	t.Run("should return a NotFound error if the Stripe product ID does not exist", func(t *testing.T) {
		clearTables(t, ctx, testPool)

		// Setup: Create a parent category
		parentCategory := td.NewValidCategory("Test Category")
		td.InsertCategory(t, ctx, parentCategory, testPool)

		// Setup: Create a product, but we will search for a different Stripe ID
		mockProduct := td.NewValidProduct("Test Product", parentCategory.ID)
		mockProduct.StripeProductID = "prod_123456789"
		td.InsertProduct(t, ctx, testPool, mockProduct)

		nonExistentStripeID := "non_existent_stripe_id"

		// Action: Search for a product that doesn't exist
		retrievedID, err := repo.GetProductIDByStripeProductID(ctx, nonExistentStripeID)

		// Assert the error and that no ID was returned
		assert.Error(t, err)
		assert.ErrorIs(t, err, errs.ErrRepositoryNotFound, "expected not found error, got: %v", err)
		assert.Equal(t, "", retrievedID)
	})

	t.Run("should return a NotFound error if the provided Stripe product ID is an empty string", func(t *testing.T) {
		clearTables(t, ctx, testPool)

		// Setup: Insert a product to ensure the table isn't empty, but the search will fail on the empty string
		parentCategory := td.NewValidCategory("Test Category")
		td.InsertCategory(t, ctx, parentCategory, testPool)
		mockProduct := td.NewValidProduct("Test Product", parentCategory.ID)
		td.InsertProduct(t, ctx, testPool, mockProduct)

		emptyStripeID := ""

		// Action: Search using an empty string
		retrievedID, err := repo.GetProductIDByStripeProductID(ctx, emptyStripeID)

		// Assert the error and that no ID was returned
		assert.Error(t, err)
		assert.ErrorIs(t, err, errs.ErrRepositoryNotFound, "expected not found error, got: %v", err)
		assert.Equal(t, "", retrievedID)
	})
}
