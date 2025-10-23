package categoryRepository_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/ports"
	td "github.com/Leviosa-care/catalog/test/helpers"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeleteCategory(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		name string
		// categoryID    string
		categoryID    uuid.UUID
		setup         func(*testing.T, context.Context, ports.CategoryRepository, uuid.UUID) // Setup action before running the function
		expectedErr   bool
		expectedErrIs error // For checking specific error types if expectedErr is true
	}{
		{
			name:       "Successfully delete an existing category",
			categoryID: uuid.New(),
			setup: func(t *testing.T, ctx context.Context, repo ports.CategoryRepository, categoryID uuid.UUID) {
				cat := td.NewValidCategory("CategoryToDelete")
				cat.ID = categoryID
				td.InsertCategory(t, ctx, cat, testPool)
			},
			expectedErr: false,
		},
		{
			name:       "Attempt to delete a non-existent category should return NotFound error",
			categoryID: uuid.New(), // A fresh, non-existent UUID
			setup: func(t *testing.T, ctx context.Context, repo ports.CategoryRepository, categoryID uuid.UUID) {
				// No category inserted for this ID
			},
			expectedErr:   true,
			expectedErrIs: errs.ErrRepositoryNotFound,
		},
		{
			name:       "Deleting a category with associated products should cascade delete products",
			categoryID: uuid.New(),
			setup: func(t *testing.T, ctx context.Context, repo ports.CategoryRepository, categoryID uuid.UUID) {
				// Insert the category
				cat := td.NewValidCategory("CategoryWithProducts")
				cat.ID = categoryID
				td.InsertCategory(t, ctx, cat, testPool)

				// Insert products linked to this category
				td.InsertProduct(t, ctx, testPool, td.NewValidProduct("Product1", cat.ID))
				td.InsertProduct(t, ctx, testPool, td.NewValidProduct("Product2", cat.ID))
				td.InsertProduct(t, ctx, testPool, td.NewValidProduct("Product3", cat.ID))

				// Verify products exist before deletion
				var initialProductCount int
				err := testPool.QueryRow(ctx, "SELECT COUNT(*) FROM catalog.products WHERE category_id = $1", cat.ID).Scan(&initialProductCount)
				require.NoError(t, err)
				require.Equal(t, 3, initialProductCount, "Expected 3 products before deletion")
			},
			expectedErr:   true,
			expectedErrIs: errs.ErrForeignKeyViolation,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean tables before each sub-test to ensure isolation
			td.ClearCategoriesTable(t, ctx, testPool) // This now also cleans products table

			// Perform test-specific setup
			if tt.setup != nil {
				tt.setup(t, ctx, repo, tt.categoryID) // Pass repo to setup if needed for pre-insertions
			}

			// Call the function under test
			err := repo.DeleteCategory(ctx, tt.categoryID)

			if tt.expectedErr {
				assert.Error(t, err, "Expected an error but got none")
				if tt.expectedErrIs != nil {
					assert.ErrorIs(t, err, tt.expectedErrIs, "Expected specific error type to be in error chain")
				}
			} else {
				assert.NoError(t, err, "Did not expect an error but got one")

				// For successful deletion, verify the category is gone
				var exists int
				err := testPool.QueryRow(ctx, "SELECT 1 FROM catalog.categories WHERE id = $1", tt.categoryID).Scan(&exists)
				assert.Error(t, err, "Expected category to be deleted, but it still exists")
				assert.True(t, errors.Is(err, sql.ErrNoRows), "Expected sql.ErrNoRows after deletion, but got different error")

				// For cascade deletion test, verify products are also gone
				if tt.name == "Deleting a category with associated products should cascade delete products" {
					var remainingProductCount int
					err := testPool.QueryRow(ctx, "SELECT COUNT(*) FROM catalog.products WHERE category_id = $1", tt.categoryID).Scan(&remainingProductCount)
					require.NoError(t, err)
					assert.Equal(t, 0, remainingProductCount, "Expected all associated products to be deleted via cascade")
				}
			}
		})
	}
}
