package productRepository_test

import (
	"context"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"
	td "github.com/Leviosa-care/catalog/test/helpers"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddProduct(t *testing.T) {
	ctx := context.Background()
	td.ClearProductsTable(t, ctx, testPool)

	// Need a valid category for FK constraints
	testCategory := td.NewValidCategory("Test Category for Products")
	td.InsertCategory(t, ctx, testCategory, testPool) // Insert a category once for all product tests
	validCategoryID := testCategory.ID
	tests := []struct {
		name    string
		product *domain.Product
		setup   func(t *testing.T, ctx context.Context) // Custom setup for specific tests
		verify  func(t *testing.T, ctx context.Context, newID string, addedProduct *domain.Product, err error)
	}{
		{
			name:    "Successfully add a basic product",
			product: td.NewValidProduct("Basic Product", validCategoryID),
			setup:   func(t *testing.T, ctx context.Context) {}, // No extra setup
			verify: func(t *testing.T, ctx context.Context, newID string, addedProduct *domain.Product, err error) {
				require.NoError(t, err, "AddProduct should not return an error")
				assert.Equal(t, addedProduct.ID.String(), newID, "Returned ID should match product ID")

				ID, err := uuid.Parse(newID)
				require.NoError(t, err)
				fetchedProduct, err := td.GetProductByID(t, ctx, ID, testPool)
				require.NoError(t, err)

				assert.Equal(t, addedProduct.ID, fetchedProduct.ID)
				assert.Equal(t, addedProduct.Name, fetchedProduct.Name)
				assert.Equal(t, addedProduct.Description, fetchedProduct.Description)
				assert.Equal(t, addedProduct.Duration, fetchedProduct.Duration)
				assert.WithinDuration(t, addedProduct.CreatedAt, fetchedProduct.CreatedAt, time.Second, "CreatedAt should be close")
				assert.Equal(t, addedProduct.Status, fetchedProduct.Status)
				assert.Equal(t, addedProduct.Availability, fetchedProduct.Availability)
				assert.Equal(t, addedProduct.BufferTime, fetchedProduct.BufferTime)
				assert.Equal(t, addedProduct.CancellationHours, fetchedProduct.CancellationHours)
				assert.Equal(t, addedProduct.StripeProductID, fetchedProduct.StripeProductID)
				assert.Equal(t, addedProduct.Metadata, fetchedProduct.Metadata)
				// UpdatedAt is set by DB on insert, so it should be close to CreatedAt
				assert.WithinDuration(t, addedProduct.CreatedAt, fetchedProduct.UpdatedAt, time.Second, "UpdatedAt should be set near CreatedAt by DB")
			},
		},
		{
			name: "Successfully add a product with optional fields nil/empty",
			product: &domain.Product{
				ID:                uuid.New(),
				Name:              "Product No Opts",
				Description:       "", // Empty string
				CategoryID:        validCategoryID,
				Duration:          90,
				CreatedAt:         time.Now().Truncate(time.Second),
				Status:            domain.Published,
				Availability:      domain.InPerson,
				BufferTime:        0,
				CancellationHours: 0,
				StripeProductID:   "",  // Empty string
				Metadata:          nil, // Nil map
			},
			setup: func(t *testing.T, ctx context.Context) {},
			verify: func(t *testing.T, ctx context.Context, newID string, addedProduct *domain.Product, err error) {
				require.NoError(t, err)
				ID, err := uuid.Parse(newID)
				require.NoError(t, err)

				fetchedProduct, err := td.GetProductByID(t, ctx, ID, testPool)
				require.NoError(t, err)

				assert.Equal(t, "", fetchedProduct.Description)     // Should be empty string, not nil
				assert.Equal(t, "", fetchedProduct.StripeProductID) // Should be empty string
				assert.Nil(t, fetchedProduct.Metadata)              // Should be nil

				// Basic check of a required field to ensure product exists
				assert.Equal(t, addedProduct.Name, fetchedProduct.Name)
			},
		},
		{
			name:    "Add product with duplicate name should return UniqueViolation error",
			product: td.NewValidProduct("Duplicate Name Product", validCategoryID),
			setup: func(t *testing.T, ctx context.Context) {
				// Insert a product with the same name first
				dupProduct := td.NewValidProduct("Duplicate Name Product", validCategoryID)
				_, err := repo.AddProduct(ctx, dupProduct)
				require.NoError(t, err, "Failed to add initial product for duplicate test")
			},
			verify: func(t *testing.T, ctx context.Context, newID string, addedProduct *domain.Product, err error) {
				require.Error(t, err)
				assert.Empty(t, newID, "ID should be empty on error")
				assert.ErrorIs(t, err, errs.ErrUniqueViolation)
			},
		},
		{
			name:    "Add product with non-existent CategoryID should return ForeignKeyViolation error",
			product: td.NewValidProduct("Product with Bad Category", uuid.New()), // Use a non-existent category ID
			setup:   func(t *testing.T, ctx context.Context) {},
			verify: func(t *testing.T, ctx context.Context, newID string, addedProduct *domain.Product, err error) {
				require.Error(t, err)
				assert.Empty(t, newID, "ID should be empty on error")
				assert.ErrorIs(t, err, errs.ErrForeignKeyViolation)
			},
		},
		{
			name: "Add product with invalid Status enum value should return ErrInvalidInput error",
			product: func() *domain.Product {
				p := td.NewValidProduct("Invalid Status Product", validCategoryID)
				p.Status = "invalid_status_value" // Set an invalid enum value
				return p
			}(),
			setup: func(t *testing.T, ctx context.Context) {},
			verify: func(t *testing.T, ctx context.Context, newID string, addedProduct *domain.Product, err error) {
				require.Error(t, err)
				assert.Empty(t, newID, "ID should be empty on error")
				assert.ErrorIs(t, err, errs.ErrInvalidInput)
			},
		},
		{
			name: "Add product with invalid Availability enum value should return CheckViolation error",
			product: func() *domain.Product {
				p := td.NewValidProduct("Invalid Availability Product", validCategoryID)
				p.Availability = "invalid_availability_value" // Set an invalid enum value
				return p
			}(),
			setup: func(t *testing.T, ctx context.Context) {},
			verify: func(t *testing.T, ctx context.Context, newID string, addedProduct *domain.Product, err error) {
				require.Error(t, err)
				assert.Empty(t, newID, "ID should be empty on error")
				assert.ErrorIs(t, err, errs.ErrCheckViolation)
			},
		},
		{
			name: "Add product with nil Metadata should be handled correctly (DB NULL)",
			product: func() *domain.Product {
				p := td.NewValidProduct("Product with Nil Metadata", validCategoryID)
				p.Metadata = nil // Explicitly set Metadata to nil
				return p
			}(),
			setup: func(t *testing.T, ctx context.Context) {},
			verify: func(t *testing.T, ctx context.Context, newID string, addedProduct *domain.Product, err error) {
				require.NoError(t, err)
				ID, err := uuid.Parse(newID)
				require.NoError(t, err)
				fetchedProduct, err := td.GetProductByID(t, ctx, ID, testPool)
				require.NoError(t, err)
				assert.Nil(t, fetchedProduct.Metadata, "Metadata should be nil for nil input")
			},
		},
		{
			name: "Add product with empty Metadata should be handled correctly (DB Empty JSONB)",
			product: func() *domain.Product {
				p := td.NewValidProduct("Product with Empty Metadata", validCategoryID)
				p.Metadata = make(map[string]any) // Explicitly set Metadata to an empty map
				return p
			}(),
			setup: func(t *testing.T, ctx context.Context) {},
			verify: func(t *testing.T, ctx context.Context, newID string, addedProduct *domain.Product, err error) {
				require.NoError(t, err)
				ID, err := uuid.Parse(newID)
				require.NoError(t, err)
				fetchedProduct, err := td.GetProductByID(t, ctx, ID, testPool)
				require.NoError(t, err)
				assert.NotNil(t, fetchedProduct.Metadata, "Metadata should not be nil for empty map input")
				assert.Empty(t, fetchedProduct.Metadata, "Metadata should be an empty map")
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Clean products and categories table
			td.ClearProductsTable(t, ctx, testPool)
			td.ClearCategoriesTable(t, ctx, testPool) // Keep this for robustness in case categories aren't cleared elsewhere

			// Re-insert the valid test category if it might have been truncated
			td.InsertCategory(t, ctx, testCategory, testPool)

			// Perform test-specific setup
			tc.setup(t, ctx) // No need for categoryID as it's provided by newValidProduct/tc.product

			// Call the function under test
			newID, err := repo.AddProduct(ctx, tc.product)

			// Verify using the custom verification function
			tc.verify(t, ctx, newID, tc.product, err)
		})
	}
}
