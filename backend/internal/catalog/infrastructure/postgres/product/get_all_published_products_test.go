package productRepository_test

import (
	"context"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"
	td "github.com/Leviosa-care/catalog/test/helpers"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetAllPublishedProducts(t *testing.T) {
	ctx := context.Background()
	// Create categories to link products to
	catElectronics := td.NewValidCategory("Electronics")
	catElectronics.ID = uuid.New()
	td.InsertCategory(t, ctx, catElectronics, testPool)

	catBooks := td.NewValidCategory("Books")
	catBooks.ID = uuid.New()
	td.InsertCategory(t, ctx, catBooks, testPool)

	tests := []struct {
		name   string
		setup  func(t *testing.T, ctx context.Context) // Custom setup for specific tests
		verify func(t *testing.T, products []*domain.ProductRes, err error)
	}{
		{
			name: "Successfully retrieve no products when product table is empty",
			setup: func(t *testing.T, ctx context.Context) {
				// Tables are cleaned by the outer loop
			},
			verify: func(t *testing.T, products []*domain.ProductRes, err error) {
				require.NoError(t, err)
				assert.Empty(t, products, "Expected no products when table is empty")
			},
		},
		{
			name: "Successfully retrieve no products when only draft/archived products exist",
			setup: func(t *testing.T, ctx context.Context) {
				// Insert a draft product
				draftProd := td.NewValidProduct("Draft Product 1", catElectronics.ID)
				draftProd.CreatedAt = time.Now().Add(-2 * time.Hour).Truncate(time.Millisecond)
				draftProd.Status = domain.Draft
				td.InsertProduct(t, ctx, testPool, draftProd)

				// Insert an archived product
				archivedProd := td.NewValidProduct("Archived Product 1", catBooks.ID)
				archivedProd.CreatedAt = time.Now().Add(-1 * time.Hour).Truncate(time.Millisecond)
				archivedProd.Status = domain.Archived
				td.InsertProduct(t, ctx, testPool, archivedProd)
			},
			verify: func(t *testing.T, products []*domain.ProductRes, err error) {
				require.NoError(t, err)
				assert.Empty(t, products, "Expected no published products when only draft/archived exist")
			},
		},
		{
			name: "Successfully retrieve only published products with correct order and data",
			setup: func(t *testing.T, ctx context.Context) {
				// Oldest published product
				prodA := td.NewValidProduct("Published Product A", catElectronics.ID)
				prodA.CreatedAt = time.Now().Add(-3 * time.Hour).Truncate(time.Millisecond)
				prodA.Status = domain.Published
				prodA.Description = "Desc A"
				prodA.Metadata = map[string]any{"tag": "new"}
				td.InsertProduct(t, ctx, testPool, prodA)

				// Draft product (should not be returned)
				draftProd := td.NewValidProduct("Draft Product X", catBooks.ID)
				draftProd.CreatedAt = time.Now().Add(-2 * time.Hour).Truncate(time.Millisecond)
				draftProd.Description = "Desc X"
				draftProd.Status = domain.Draft
				td.InsertProduct(t, ctx, testPool, draftProd)

				// Newer published product
				prodB := td.NewValidProduct("Published Product B", catBooks.ID)
				prodB.CreatedAt = time.Now().Add(-1 * time.Hour).Truncate(time.Millisecond)
				prodB.Status = domain.Published
				prodB.Description = "Desc B" // Test empty optional field
				td.InsertProduct(t, ctx, testPool, prodB)

				// Archived product (should not be returned)
				archivedProd := td.NewValidProduct("Archived Product Y", catElectronics.ID)
				archivedProd.CreatedAt = time.Now().Truncate(time.Millisecond) // Latest created but archived
				archivedProd.Status = domain.Archived
				archivedProd.Description = "Desc Y" // Test empty optional field
				td.InsertProduct(t, ctx, testPool, archivedProd)
			},
			verify: func(t *testing.T, products []*domain.ProductRes, err error) {
				require.NoError(t, err)
				assert.Len(t, products, 2, "Expected exactly 2 published products to be returned")

				// Verify order: Product B (newer) should be first, Product A (older) second
				assert.Equal(t, "Published Product B", products[0].Name, "Expected 'Published Product B' first")
				assert.Equal(t, "Published Product A", products[1].Name, "Expected 'Published Product A' second")

				// Verify details for Product B (products[0])
				prodB := products[0]
				assert.Equal(t, domain.Published, prodB.Status)
				assert.Equal(t, catBooks.ID, prodB.Category.ID)
				assert.Equal(t, catBooks.Name, prodB.Category.Name)
				assert.Equal(t, "Desc B", prodB.Description, "Expected empty description")
				// assert.WithinDuration(t, prodB.CreatedAt, prodB.UpdatedAt, time.Second, "UpdatedAt should be same as CreatedAt for Prod B")

				// Verify details for Product A (products[1])
				prodA := products[1]
				assert.Equal(t, domain.Published, prodA.Status)
				assert.Equal(t, catElectronics.ID, prodA.Category.ID)
				assert.Equal(t, catElectronics.Name, prodA.Category.Name)
				assert.Equal(t, "Desc A", prodA.Description)
				assert.Equal(t, map[string]any{"tag": "new"}, prodA.Metadata)
				// assert.WithinDuration(t, prodA.CreatedAt, prodA.UpdatedAt, time.Second, "UpdatedAt should be same as CreatedAt for Prod A")

				// Ensure no non-published products are present
				for _, p := range products {
					assert.Equal(t, domain.Published, p.Status, "All returned products must have 'published' status")
				}
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Clean tables before each sub-test
			td.ClearProductsTable(t, ctx, testPool)
			td.ClearCategoriesTable(t, ctx, testPool)

			// Re-insert initial test categories as they might have been truncated
			td.InsertCategory(t, ctx, catElectronics, testPool)
			td.InsertCategory(t, ctx, catBooks, testPool)

			// Perform test-specific setup
			tc.setup(t, ctx)

			// Call the function under test
			products, err := repo.GetAllPublishedProducts(ctx)

			// Verify using the custom verification function
			tc.verify(t, products, err)
		})
	}
}
