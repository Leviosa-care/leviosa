package sharedRepository_test

import (
	"context"
	"testing"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"
	td "github.com/Leviosa-care/leviosa/backend/test/helpers"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetProductByID(t *testing.T) {
	ctx := context.Background()

	// Create a category for products to link to
	testCategory := td.NewValidCategory("Test Category for ProductByID")
	td.InsertCategory(t, ctx, testCategory, testPool)
	validCategoryID := testCategory.ID

	productID := uuid.New()

	tests := []struct {
		name   string
		setup  func(t *testing.T, ctx context.Context, categoryID uuid.UUID, productID uuid.UUID) *domain.Product // Setup function, returns the inserted product (if any)
		verify func(t *testing.T, returnedProduct *domain.Product, err error, expectedProduct *domain.Product)
	}{
		{
			name: "Successfully retrieve an existing product with all fields populated",
			setup: func(t *testing.T, ctx context.Context, categoryID uuid.UUID, productID uuid.UUID) *domain.Product {
				product := td.NewValidProduct("Test Product 1", categoryID)
				product.ID = productID
				td.InsertProduct(t, ctx, testPool, product)
				return product
			},
			verify: func(t *testing.T, returnedProd *domain.Product, err error, expectedProd *domain.Product) {
				require.NoError(t, err)
				require.NotNil(t, returnedProd)

				assert.Equal(t, expectedProd.ID, returnedProd.ID)
				assert.Equal(t, expectedProd.Name, returnedProd.Name)
				assert.Equal(t, expectedProd.Description, returnedProd.Description)
				assert.Equal(t, expectedProd.Duration, returnedProd.Duration)
				assert.Equal(t, expectedProd.Status, returnedProd.Status)
				assert.Equal(t, expectedProd.Availability, returnedProd.Availability)
				assert.Equal(t, expectedProd.BufferTime, returnedProd.BufferTime)
				assert.Equal(t, expectedProd.CancellationHours, returnedProd.CancellationHours)
				assert.Equal(t, expectedProd.Metadata, returnedProd.Metadata) // Deep equality for maps
			},
		},
		{
			name: "Successfully retrieve an existing product with optional fields NULL/empty",
			setup: func(t *testing.T, ctx context.Context, categoryID uuid.UUID, productID uuid.UUID) *domain.Product {
				prod := td.NewValidProduct("Test Product 2 (Optional Null)", categoryID)
				prod.ID = productID        // Use the generated ID
				prod.Metadata = nil        // Nil map
				prod.Status = domain.Draft // Can be any status
				td.InsertProduct(t, ctx, testPool, prod)
				return prod
			},
			verify: func(t *testing.T, returnedProd *domain.Product, err error, expectedProd *domain.Product) {
				require.NoError(t, err)
				require.NotNil(t, returnedProd)

				assert.Equal(t, expectedProd.ID, returnedProd.ID)
				assert.Empty(t, returnedProd.Metadata, "Metadata should be empty map") // Unmarshaled nil into make(map[string]any)
				assert.Equal(t, expectedProd.Status, returnedProd.Status)
				// ... other assertions as above for non-optional fields
			},
		},
		{
			name: "Product not found should return NotFound error",
			setup: func(t *testing.T, ctx context.Context, categoryID uuid.UUID, productID uuid.UUID) *domain.Product {
				// TODO:  that causes an issue when trying to fetch from the database
				// return nil // No product inserted for this test
				prod := td.NewValidProduct("Test Product 1", categoryID)
				td.InsertProduct(t, ctx, testPool, prod)
				return prod
			},
			verify: func(t *testing.T, returnedProd *domain.Product, err error, expectedProd *domain.Product) {
				require.Error(t, err)
				assert.ErrorIs(t, err, errs.ErrRepositoryNotFound, "Expected NotFound error")
				assert.Nil(t, returnedProd, "Expected nil product on NotFound error")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean tables before each sub-test
			td.ClearProductsTable(t, ctx, testPool)
			td.ClearCategoriesTable(t, ctx, testPool)

			// Re-insert initial test category as it might have been truncated
			td.InsertCategory(t, ctx, testCategory, testPool)

			// Perform test-specific setup and get the product that was inserted (if any)
			insertedProduct := tt.setup(t, ctx, validCategoryID, productID)
			// Call the function under test
			returnedProduct, err := repo.GetProductByID(ctx, productID)

			// Verify using the custom verification function
			tt.verify(t, returnedProduct, err, insertedProduct)
		})
	}
}
