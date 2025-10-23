package productRepository_test

import (
	"context"
	"testing"

	td "github.com/Leviosa-care/catalog/test/helpers"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeleteProduct(t *testing.T) {
	ctx := context.Background()
	// Need a valid category for FK constraints on products
	testCategory := td.NewValidCategory("Test Category for Product Deletion")
	td.InsertCategory(t, ctx, testCategory, testPool)
	validCategoryID := testCategory.ID

	productID := uuid.New()

	tests := []struct {
		name   string
		setup  func(t *testing.T, ctx context.Context, productID uuid.UUID)                    // Setup before calling DeleteProduct
		verify func(t *testing.T, ctx context.Context, productID uuid.UUID, expectedErr error) // Custom verification after deletion attempt
	}{
		{
			name: "Successfully delete an existing product",
			setup: func(t *testing.T, ctx context.Context, productID uuid.UUID) {
				prod := td.NewValidProduct("Product To Delete", validCategoryID)
				prod.ID = productID
				td.InsertProduct(t, ctx, testPool, prod)
			},
			verify: func(t *testing.T, ctx context.Context, productID uuid.UUID, err error) {
				require.NoError(t, err)
				// Verify product is actually deleted
				_, getErr := td.GetProductByID(t, ctx, productID, testPool) // Assuming you have a GetProductByID
				assert.Error(t, getErr, "Expected product to be deleted, but it still exists")
				assert.ErrorIs(t, getErr, errs.ErrRepositoryNotFound, "Expected NotFound error after deletion")
			},
		},
		{
			name:  "Attempt to delete a non-existent product should return NotFound error",
			setup: func(t *testing.T, ctx context.Context, productID uuid.UUID) {},
			verify: func(t *testing.T, ctx context.Context, productID uuid.UUID, err error) {
				require.Error(t, err)
				assert.ErrorIs(t, err, errs.ErrRepositoryNotFound)
			},
		},
		{
			name: "Deleting a product with associated prices should cascade delete prices",
			setup: func(t *testing.T, ctx context.Context, productID uuid.UUID) {
				prod := td.NewValidProduct("Product With Prices", validCategoryID)
				prod.ID = productID
				td.InsertProduct(t, ctx, testPool, prod)

				// Insert prices linked to this product
				// TODO: first create the price then insert it
				price1 := td.NewValidPrice()
				price1.ProductID = prod.ID
				price1.StripePriceID = "sk_price1"
				price2 := td.NewValidPrice()
				price2.ProductID = prod.ID
				price2.StripePriceID = "sk_price2"
				td.InsertPrice(t, ctx, price1, testPool)
				td.InsertPrice(t, ctx, price2, testPool)

				// Verify prices exist before deletion
				var initialPriceCount int
				err := testPool.QueryRow(ctx, "SELECT COUNT(*) FROM catalog.prices WHERE product_id = $1", prod.ID).Scan(&initialPriceCount)
				require.NoError(t, err)
				require.Equal(t, 2, initialPriceCount, "Expected 2 prices before product deletion")
			},
			verify: func(t *testing.T, ctx context.Context, productID uuid.UUID, err error) {
				require.NoError(t, err)
				// Verify product is deleted
				_, getErr := td.GetProductByID(t, ctx, productID, testPool)
				assert.Error(t, getErr, "Expected product to be deleted")
				assert.ErrorIs(t, getErr, errs.ErrRepositoryNotFound, "Expected NotFound error for deleted product")

				// Verify associated prices are also deleted
				var remainingPriceCount int
				err = testPool.QueryRow(ctx, "SELECT COUNT(*) FROM catalog.prices WHERE product_id = $1", productID).Scan(&remainingPriceCount)
				require.NoError(t, err)
				assert.Equal(t, 0, remainingPriceCount, "Expected all associated prices to be deleted via cascade")
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Clean tables before each sub-test
			// Order matters: products before categories to avoid FK issues on TRUNCATE
			td.ClearProductsTable(t, ctx, testPool)
			td.ClearCategoriesTable(t, ctx, testPool) // This also cleans categories

			// Re-insert the valid test category as it might have been truncated
			td.InsertCategory(t, ctx, testCategory, testPool)

			// Perform test-specific setup
			tc.setup(t, ctx, productID)

			// Call the function under test
			err := repo.DeleteProduct(ctx, productID)

			// Verify using the custom verification function
			tc.verify(t, ctx, productID, err)
		})
	}
}
