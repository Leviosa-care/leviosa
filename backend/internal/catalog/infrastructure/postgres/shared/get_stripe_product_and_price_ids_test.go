package sharedRepository_test

import (
	"context"
	"testing"

	td "github.com/Leviosa-care/leviosa/backend/test/helpers"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetStripeProductAndPriceIDs(t *testing.T) {
	ctx := context.Background()

	// Need a valid category for products
	testCategory := td.NewValidCategory("Test Category for Stripe IDs")
	td.InsertCategory(t, ctx, testCategory, testPool)
	validCategoryID := testCategory.ID

	productID := uuid.New()

	tests := []struct {
		name                    string
		setup                   func(t *testing.T, ctx context.Context, productID uuid.UUID, categoryID uuid.UUID)
		expectedStripeProductID string
		expectedStripePriceIDs  []string
		expectedErrIs           error
	}{
		// NOTE: done
		{
			name: "Product exists, has Stripe Product ID and multiple Prices",
			setup: func(t *testing.T, ctx context.Context, productID uuid.UUID, categoryID uuid.UUID) {
				prod := td.NewValidProduct("Prod With Stripe IDs", categoryID)
				prod.ID = productID
				prod.StripeProductID = "prod_full_example"
				td.InsertProduct(t, ctx, testPool, prod)

				price1 := td.NewValidPrice()
				price1.ProductID = productID
				price1.StripePriceID = "price_abc"
				td.InsertPrice(t, ctx, price1, testPool)

				price2 := td.NewValidPrice()
				price2.ProductID = productID
				price2.StripePriceID = "price_def"
				td.InsertPrice(t, ctx, price2, testPool)

				price3 := td.NewValidPrice()
				price3.ProductID = productID
				price3.StripePriceID = "price_ghi"
				td.InsertPrice(t, ctx, price3, testPool)
			},
			expectedStripeProductID: "prod_full_example",
			expectedStripePriceIDs:  []string{"price_abc", "price_def", "price_ghi"},
			expectedErrIs:           nil,
		},
		{
			name: "Product exists, has Stripe Product ID, no Prices",
			setup: func(t *testing.T, ctx context.Context, productID uuid.UUID, categoryID uuid.UUID) {
				prod := td.NewValidProduct("Prod No Prices", categoryID)
				prod.ID = productID
				prod.StripeProductID = "prod_no_prices"
				td.InsertProduct(t, ctx, testPool, prod)
			},
			expectedStripeProductID: "prod_no_prices",
			expectedStripePriceIDs:  []string{}, // Expect empty slice
			expectedErrIs:           nil,
		},
		{
			name: "Product exists, no Stripe Product ID (empty string), has Prices",
			setup: func(t *testing.T, ctx context.Context, productID uuid.UUID, categoryID uuid.UUID) {
				prod := td.NewValidProduct("Prod No Product ID", categoryID)
				prod.ID = productID
				prod.StripeProductID = "" // Empty string for StripeProductID
				td.InsertProduct(t, ctx, testPool, prod)
				price := td.NewValidPrice()
				price.ProductID = productID
				price.StripePriceID = "price_xyz"
				td.InsertPrice(t, ctx, price, testPool)
			},
			expectedStripeProductID: "", // Expect empty string
			expectedStripePriceIDs:  []string{"price_xyz"},
			expectedErrIs:           nil,
		},
		{
			name: "Product exists, no Stripe Product ID, no Prices",
			setup: func(t *testing.T, ctx context.Context, productID uuid.UUID, categoryID uuid.UUID) {
				prod := td.NewValidProduct("Prod No IDs At All", categoryID)
				prod.ID = productID
				prod.StripeProductID = "" // Empty string
				td.InsertProduct(t, ctx, testPool, prod)
			},
			expectedStripeProductID: "",
			expectedStripePriceIDs:  []string{},
			expectedErrIs:           nil,
		},
		{
			name: "Product not found should return NotFound error",
			setup: func(t *testing.T, ctx context.Context, productID uuid.UUID, categoryID uuid.UUID) {
				// No product inserted
			},
			expectedStripeProductID: "",
			expectedStripePriceIDs:  nil, // Or []string{} depending on your preference for nil vs empty slice on error
			expectedErrIs:           errs.ErrRepositoryNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean tables before each sub-test
			td.ClearProductsTable(t, ctx, testPool)
			td.ClearCategoriesTable(t, ctx, testPool) // Also cleans categories
			// Re-insert the test category as it might have been truncated
			td.InsertCategory(t, ctx, testCategory, testPool)

			// Perform test-specific setup
			tt.setup(t, ctx, productID, validCategoryID)

			// Call the function under test
			actualProductID, actualPriceIDs, err := repo.GetStripeProductAndPriceIDs(ctx, productID)

			if tt.expectedErrIs != nil {
				assert.Error(t, err, "Expected an error but got none")
				assert.ErrorIs(t, err, tt.expectedErrIs, "Expected specific error type")
				assert.Empty(t, actualProductID, "Expected empty product ID on error")
				assert.Nil(t, actualPriceIDs, "Expected nil price IDs on error") // Or assert.Empty if you prefer empty slice on error
			} else {
				assert.NoError(t, err, "Did not expect an error but got one")
				assert.Equal(t, tt.expectedStripeProductID, actualProductID, "Stripe Product ID mismatch")
				assert.ElementsMatch(t, tt.expectedStripePriceIDs, actualPriceIDs, "Stripe Price IDs mismatch") // Use ElementsMatch for slice comparison (order doesn't matter)
			}
		})
	}
}
