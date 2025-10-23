package productRepository_test

import (
	"context"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"
	td "github.com/Leviosa-care/leviosa/backend/test/helpers"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetAllProducts(t *testing.T) {
	ctx := context.Background()
	// Create a few categories for products to link to
	cat1 := td.NewValidCategory("Electronics")
	cat1.ID = uuid.New()
	cat1.Metadata = map[string]any{"type": "consumer"}
	td.InsertCategory(t, ctx, cat1, testPool)

	cat2 := td.NewValidCategory("Books")
	cat2.ID = uuid.New()
	cat2.Metadata = map[string]any{"genre": "fiction"}
	td.InsertCategory(t, ctx, cat2, testPool)

	tests := []struct {
		name                  string
		setup                 func(t *testing.T, ctx context.Context) // Custom setup for specific tests
		expectedProductsCount int
		verify                func(t *testing.T, products []*domain.ProductRes, err error)
	}{
		{
			name: "Successfully retrieve no products when table is empty",
			setup: func(t *testing.T, ctx context.Context) {
				// Tables are cleaned by the outer loop
			},
			expectedProductsCount: 0,
			verify: func(t *testing.T, products []*domain.ProductRes, err error) {
				require.NoError(t, err)
				assert.Empty(t, products, "Expected no products when table is empty")
			},
		},
		{
			name: "Successfully retrieve multiple products with correct order and category data",
			setup: func(t *testing.T, ctx context.Context) {
				// Products will be inserted in reverse order of created_at for testing DESC sorting
				// Product 3 (latest)
				prod3 := td.NewValidProduct("Product Gamma", cat1.ID)
				prod3.CreatedAt = time.Now().Add(3 * time.Minute).Truncate(time.Millisecond)
				prod3.Description = "Desc Gamma"
				prod3.Status = domain.Published
				td.InsertProduct(t, ctx, testPool, prod3)

				// Product 2
				prod2 := td.NewValidProduct("Product Beta", cat2.ID)
				prod2.CreatedAt = time.Now().Add(2 * time.Minute).Truncate(time.Millisecond)
				prod2.Description = "Desc Beta"
				prod2.Metadata = map[string]any{"color": "green"} // Test populated metadata
				td.InsertProduct(t, ctx, testPool, prod2)

				// Product 1 (oldest)
				prod1 := td.NewValidProduct("Product Alpha", cat1.ID)
				prod1.CreatedAt = time.Now().Add(1 * time.Minute).Truncate(time.Millisecond)
				prod1.Description = "Desc Alpha"
				prod1.Metadata = nil       // Test nil Metadata
				prod1.StripeProductID = "" // Test empty StripeProductID
				td.InsertProduct(t, ctx, testPool, prod1)
			},
			expectedProductsCount: 3,
			verify: func(t *testing.T, products []*domain.ProductRes, err error) {
				require.NoError(t, err)
				assert.Len(t, products, 3, "Expected 3 products to be returned")

				// Verify order (DESC by CreatedAt) and basic data
				assert.Equal(t, "Product Gamma", products[0].Name, "Expected latest product first")
				assert.Equal(t, "Product Beta", products[1].Name, "Expected second latest product second")
				assert.Equal(t, "Product Alpha", products[2].Name, "Expected oldest product last")

				// Verify details of a specific product and its joined category
				// --- Detailed Verification of products[0] (Product Gamma) ---
				gammaProd := products[0]
				assert.Equal(t, cat1.ID, gammaProd.Category.ID)
				assert.Equal(t, cat1.Name, gammaProd.Category.Name)
				assert.Equal(t, cat1.Description, gammaProd.Category.Description)
				assert.Equal(t, cat1.Metadata, gammaProd.Category.Metadata)
				assert.Equal(t, domain.Published, gammaProd.Status)

				// --- Detailed Verification of products[1] (Product Beta) ---
				betaProd := products[1]
				assert.Equal(t, cat2.ID, betaProd.Category.ID)
				assert.Equal(t, cat2.Name, betaProd.Category.Name)
				assert.Equal(t, cat2.Description, betaProd.Category.Description)
				assert.Equal(t, cat2.Metadata, betaProd.Category.Metadata)
				assert.Equal(t, "Desc Beta", betaProd.Description, "Product Beta Description should be empty string")
				assert.Equal(t, map[string]any{"color": "green"}, betaProd.Metadata, "Product Beta Metadata should be correctly populated")

				// --- Detailed Verification of products[2] (Product Alpha) ---
				alphaProd := products[2]
				assert.Equal(t, cat1.ID, alphaProd.Category.ID)
				assert.Equal(t, cat1.Name, alphaProd.Category.Name)
				assert.Equal(t, cat1.Description, alphaProd.Category.Description)
				assert.Equal(t, cat1.Metadata, alphaProd.Category.Metadata)
				assert.Equal(t, "Desc Alpha", alphaProd.Description, "Product alpha Description should be empty string")
				assert.Nil(t, alphaProd.Metadata, "Product Alpha Metadata should be nil")

			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Clean tables before each sub-test
			td.ClearProductsTable(t, ctx, testPool)
			td.ClearCategoriesTable(t, ctx, testPool)

			// Re-insert initial test categories as they might have been truncated
			td.InsertCategory(t, ctx, cat1, testPool)
			td.InsertCategory(t, ctx, cat2, testPool)

			// Perform test-specific setup
			tc.setup(t, ctx)

			// Call the function under test
			products, err := repo.GetAllProducts(ctx)

			// Verify using the custom verification function
			tc.verify(t, products, err)
		})
	}
}
