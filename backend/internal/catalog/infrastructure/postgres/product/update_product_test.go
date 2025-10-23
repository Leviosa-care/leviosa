package productRepository_test

import (
	"context"
	"testing"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"
	td "github.com/Leviosa-care/catalog/test/helpers"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdateProduct(t *testing.T) {
	ctx := context.Background()
	// Need a valid category for products
	testCategory := td.NewValidCategory("Test Category for Product Update")
	td.InsertCategory(t, ctx, testCategory, testPool)
	validCategoryID := testCategory.ID

	// Another category for category ID update tests
	anotherCategory := td.NewValidCategory("Another Category")
	td.InsertCategory(t, ctx, anotherCategory, testPool)
	anotherCategoryID := anotherCategory.ID

	productID := uuid.New()

	tests := []struct {
		name    string
		setup   func(t *testing.T, ctx context.Context, productID uuid.UUID, categoryID uuid.UUID) *domain.Product // Setup function, returns the inserted product (if any)
		request *domain.UpdateProductRequest
		verify  func(t *testing.T, actualErr error, productID uuid.UUID, initialProduct *domain.Product, initialCategory *domain.Category, updatedRequest *domain.UpdateProductRequest)
	}{
		{
			name: "Successfully update a single field (Name)",
			setup: func(t *testing.T, ctx context.Context, productID uuid.UUID, categoryID uuid.UUID) *domain.Product {
				prod := td.NewValidProduct("Original Product Name", categoryID)
				prod.ID = productID
				td.InsertProduct(t, ctx, testPool, prod)
				return prod
			},
			request: &domain.UpdateProductRequest{
				Name: td.StrPtr("Updated Product Name"),
			},
			verify: func(t *testing.T, actualErr error, productID uuid.UUID, initialProduct *domain.Product, initialCategory *domain.Category, updatedRequest *domain.UpdateProductRequest) {
				require.NoError(t, actualErr)
				updatedProdRes, err := td.GetProductByID(t, ctx, productID, testPool)
				require.NoError(t, err)

				require.NotNil(t, updatedProdRes)
				assert.Equal(t, *updatedRequest.Name, updatedProdRes.Name)
				assert.Equal(t, initialProduct.Description, updatedProdRes.Description, "Description should not change")
				assert.Equal(t, initialProduct.Duration, updatedProdRes.Duration, "Duration should not change")
				// UpdatedAt is handled by DB trigger, no direct assertion on value
			},
		},
		{
			name: "Successfully update multiple fields",
			setup: func(t *testing.T, ctx context.Context, productID uuid.UUID, categoryID uuid.UUID) *domain.Product {
				prod := td.NewValidProduct("Initial Multi-Update Product", categoryID)
				prod.ID = productID
				prod.Description = "Initial desc"
				prod.Duration = 10
				prod.Status = domain.Draft
				prod.Availability = domain.Online
				prod.BufferTime = 5
				prod.CancellationHours = 1
				prod.Metadata = map[string]any{"old_key": "old_value"}
				td.InsertProduct(t, ctx, testPool, prod)
				return prod
			},
			request: &domain.UpdateProductRequest{
				Name:              td.StrPtr("Updated Multi-Update Name"),
				Description:       td.StrPtr("Updated desc"),
				Duration:          td.IntPtr(120),
				Status:            td.StrPtr(string(domain.Published)),
				Availability:      td.StrPtr(string(domain.InPerson)),
				BufferTime:        td.IntPtr(30),
				CancellationHours: td.IntPtr(48),
				Metadata:          map[string]any{"new_key": "new_value", "count": float64(5)},
			},
			verify: func(t *testing.T, actualErr error, productID uuid.UUID, initialProduct *domain.Product, initialCategory *domain.Category, updatedRequest *domain.UpdateProductRequest) {
				require.NoError(t, actualErr)
				updatedProdRes, err := td.GetProductByID(t, ctx, productID, testPool)
				require.NoError(t, err)

				assert.Equal(t, *updatedRequest.Name, updatedProdRes.Name)
				assert.Equal(t, *updatedRequest.Description, updatedProdRes.Description)
				assert.Equal(t, *updatedRequest.Duration, updatedProdRes.Duration)
				assert.Equal(t, domain.Published, updatedProdRes.Status)      // Convert string to PublishedStatus for comparison
				assert.Equal(t, domain.InPerson, updatedProdRes.Availability) // Convert string to AvailabilityType for comparison
				assert.Equal(t, *updatedRequest.BufferTime, updatedProdRes.BufferTime)
				assert.Equal(t, *updatedRequest.CancellationHours, updatedProdRes.CancellationHours)
				assert.Equal(t, updatedRequest.Metadata, updatedProdRes.Metadata) // Deep equality for maps
			},
		},
		{
			name: "Successfully update CategoryID",
			setup: func(t *testing.T, ctx context.Context, productID uuid.UUID, categoryID uuid.UUID) *domain.Product {
				prod := td.NewValidProduct("Product For Category Change", categoryID)
				prod.ID = productID
				td.InsertProduct(t, ctx, testPool, prod)
				return prod
			},
			request: &domain.UpdateProductRequest{
				CategoryID: td.StrPtr(anotherCategoryID.String()), // Change to another category
			},
			verify: func(t *testing.T, actualErr error, productID uuid.UUID, initialProduct *domain.Product, initialCategory *domain.Category, updatedRequest *domain.UpdateProductRequest) {
				require.NoError(t, actualErr)
				updatedProdRes, err := td.GetProductByID(t, ctx, productID, testPool)
				require.NoError(t, err)

				require.NotNil(t, updatedProdRes)
				assert.Equal(t, anotherCategoryID, updatedProdRes.CategoryID, "CategoryID should be updated")
				// Other fields should remain unchanged
				assert.Equal(t, initialProduct.Name, updatedProdRes.Name)
			},
		},
		{
			name: "Update with nil fields in request should not change existing values",
			setup: func(t *testing.T, ctx context.Context, productID uuid.UUID, categoryID uuid.UUID) *domain.Product {
				prod := td.NewValidProduct("Product Nil Fields", categoryID)
				prod.ID = productID
				prod.Description = "Original Desc"
				prod.Duration = 100
				prod.Status = domain.Published
				prod.Availability = domain.Hybrid
				prod.BufferTime = 20
				prod.CancellationHours = 24
				prod.Metadata = map[string]any{"original": "data"}
				td.InsertProduct(t, ctx, testPool, prod)
				return prod
			},
			request: &domain.UpdateProductRequest{}, // Empty request, all fields nil
			verify: func(t *testing.T, actualErr error, productID uuid.UUID, initialProduct *domain.Product, initialCategory *domain.Category, updatedRequest *domain.UpdateProductRequest) {
				require.Error(t, actualErr)
				assert.ErrorIs(t, actualErr, errs.ErrNoFieldsForUpdate, "Expected InvalidInput error for no fields")

				// Verify product state is unchanged (except UpdatedAt, handled by DB)
				currentProdRes, err := td.GetProductByID(t, ctx, productID, testPool)
				require.NoError(t, err) // Should still be able to fetch

				assert.Equal(t, initialProduct.Name, currentProdRes.Name)
				assert.Equal(t, initialProduct.Description, currentProdRes.Description)
				assert.Equal(t, initialProduct.Duration, currentProdRes.Duration)
				assert.Equal(t, initialProduct.Status, currentProdRes.Status)
				assert.Equal(t, initialProduct.Availability, currentProdRes.Availability)
				assert.Equal(t, initialProduct.BufferTime, currentProdRes.BufferTime)
				assert.Equal(t, initialProduct.CancellationHours, currentProdRes.CancellationHours)
				assert.Equal(t, initialProduct.Metadata, currentProdRes.Metadata)
			},
		},
		{
			name: "Update a non-existent product should return NotFound error",
			setup: func(t *testing.T, ctx context.Context, productID uuid.UUID, categoryID uuid.UUID) *domain.Product {
				return nil // No product inserted
			},
			request: &domain.UpdateProductRequest{
				Name: td.StrPtr("Non Existent Update"),
			},
			verify: func(t *testing.T, actualErr error, productID uuid.UUID, initialProduct *domain.Product, initialCategory *domain.Category, updatedRequest *domain.UpdateProductRequest) {
				require.Error(t, actualErr)
				assert.ErrorIs(t, actualErr, errs.ErrRepositoryNotFound, "Expected NotFound error")
			},
		},
		{
			name: "Update with duplicate name should return UniqueViolation error",
			setup: func(t *testing.T, ctx context.Context, productID uuid.UUID, categoryID uuid.UUID) *domain.Product {
				// Product to be updated
				prodToUpdate := td.NewValidProduct("Product A", categoryID)
				prodToUpdate.ID = productID
				td.InsertProduct(t, ctx, testPool, prodToUpdate)

				// Product with the conflicting name
				conflictingProd := td.NewValidProduct("Conflicting Name", categoryID)
				conflictingProd.ID = uuid.New() // Different ID
				td.InsertProduct(t, ctx, testPool, conflictingProd)
				return prodToUpdate
			},
			request: &domain.UpdateProductRequest{
				Name: td.StrPtr("Conflicting Name"), // Attempt to change name to an existing one
			},
			verify: func(t *testing.T, actualErr error, productID uuid.UUID, initialProduct *domain.Product, initialCategory *domain.Category, updatedRequest *domain.UpdateProductRequest) {
				require.Error(t, actualErr)
				assert.ErrorIs(t, actualErr, errs.ErrUniqueViolation, "Expected UniqueViolation error")
				// Verify original product name did not change
				originalProdRes, err := td.GetProductByID(t, ctx, productID, testPool)
				require.NoError(t, err)
				assert.Equal(t, initialProduct.Name, originalProdRes.Name, "Product name should not have changed after failed update")
			},
		},
		{
			name: "Update with non-existent CategoryID should return ForeignKeyViolation error",
			setup: func(t *testing.T, ctx context.Context, productID uuid.UUID, categoryID uuid.UUID) *domain.Product {
				prod := td.NewValidProduct("Product for FK Test", categoryID)
				prod.ID = productID
				td.InsertProduct(t, ctx, testPool, prod)
				return prod
			},
			request: &domain.UpdateProductRequest{
				CategoryID: td.StrPtr(uuid.New().String()), // Non-existent category ID
			},
			verify: func(t *testing.T, actualErr error, productID uuid.UUID, initialProduct *domain.Product, initialCategory *domain.Category, updatedRequest *domain.UpdateProductRequest) {
				require.Error(t, actualErr)
				assert.ErrorIs(t, actualErr, errs.ErrForeignKeyViolation, "Expected ForeignKeyViolation error")
				// Verify original category ID did not change
				originalProdRes, err := td.GetProductByID(t, ctx, productID, testPool)
				require.NoError(t, err)
				assert.Equal(t, initialProduct.CategoryID, originalProdRes.CategoryID, "Product category ID should not have changed")
			},
		},
		{
			name: "Update with invalid Status value should return ErrInvalidInput error",
			setup: func(t *testing.T, ctx context.Context, productID uuid.UUID, categoryID uuid.UUID) *domain.Product {
				prod := td.NewValidProduct("Prod Invalid Status", categoryID)
				prod.ID = productID
				td.InsertProduct(t, ctx, testPool, prod)
				return prod
			},
			request: &domain.UpdateProductRequest{
				Status: td.StrPtr("invalid_status_enum"), // Invalid ENUM value
			},
			verify: func(t *testing.T, actualErr error, productID uuid.UUID, initialProduct *domain.Product, initialCategory *domain.Category, updatedRequest *domain.UpdateProductRequest) {
				require.Error(t, actualErr)
				assert.ErrorIs(t, actualErr, errs.ErrInvalidInput, "Expected ErrInvalidInput error for invalid status")
				// Verify original status did not change
				originalProdRes, err := td.GetProductByID(t, ctx, productID, testPool)
				require.NoError(t, err)
				assert.Equal(t, initialProduct.Status, originalProdRes.Status, "Product status should not have changed")
			},
		},
		{
			name: "Update with invalid Availability value should return CheckViolation error",
			setup: func(t *testing.T, ctx context.Context, productID uuid.UUID, categoryID uuid.UUID) *domain.Product {
				prod := td.NewValidProduct("Prod Invalid Availability", categoryID)
				prod.ID = productID
				td.InsertProduct(t, ctx, testPool, prod)
				return prod
			},
			request: &domain.UpdateProductRequest{
				Availability: td.StrPtr("invalid_availability_enum"), // Invalid ENUM value
			},
			verify: func(t *testing.T, actualErr error, productID uuid.UUID, initialProduct *domain.Product, initialCategory *domain.Category, updatedRequest *domain.UpdateProductRequest) {
				require.Error(t, actualErr)
				assert.ErrorIs(t, actualErr, errs.ErrCheckViolation, "Expected CheckViolation error for invalid availability")
				// Verify original availability did not change
				originalProdRes, err := td.GetProductByID(t, ctx, productID, testPool)
				require.NoError(t, err)
				assert.Equal(t, initialProduct.Availability, originalProdRes.Availability, "Product availability should not have changed")
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Clean tables before each sub-test
			td.ClearProductsTable(t, ctx, testPool)
			td.ClearCategoriesTable(t, ctx, testPool)

			// Re-insert initial test categories as they might have been truncated
			td.InsertCategory(t, ctx, testCategory, testPool)
			td.InsertCategory(t, ctx, anotherCategory, testPool) // Ensure the second category is also present

			// Perform test-specific setup and get the initial product (if inserted)
			// initialProduct := tc.setup(t, ctx, productIDForSetup, validCategoryID)
			initialProduct := tc.setup(t, ctx, productID, validCategoryID)

			// Call the function under test
			actualErr := repo.UpdateProduct(ctx, productID, tc.request)

			// Verify using the custom verification function
			tc.verify(t, actualErr, productID, initialProduct, testCategory, tc.request)
		})
	}
}
