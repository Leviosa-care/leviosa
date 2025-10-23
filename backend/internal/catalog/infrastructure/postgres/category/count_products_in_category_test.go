package categoryRepository_test

import (
	"context"
	"testing"

	td "github.com/Leviosa-care/catalog/test/helpers"

	// "github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCountProductsInCategory(t *testing.T) {
	ctx := context.Background()

	existingCategory := td.NewValidCategory("Electronics For Counting")
	existingCategory.ID = uuid.MustParse("d0a0a0a0-d0a0-d0a0-d0a0-d0a0a0a0a0a0") // Fixed ID for easier reference

	tests := []struct {
		name string
		// categoryID    string
		categoryID    uuid.UUID
		setup         func(*testing.T) // Setup before each test case
		expectedCount int
		expectedErr   bool
		expectedErrIs error // Specific error to check with errors.Is
	}{
		{
			name:       "Category with multiple products should return correct count",
			categoryID: existingCategory.ID,
			setup: func(t *testing.T) {
				// Insert the category
				td.InsertCategory(t, ctx, existingCategory, testPool)
				// Insert products linked to this category
				td.InsertProduct(t, ctx, testPool, td.NewValidProduct("Laptop A", existingCategory.ID))
				td.InsertProduct(t, ctx, testPool, td.NewValidProduct("Mouse B", existingCategory.ID))
				td.InsertProduct(t, ctx, testPool, td.NewValidProduct("Keyboard C", existingCategory.ID))
			},
			expectedCount: 3,
			expectedErr:   false,
		},
		{
			name:       "Category with no products should return 0",
			categoryID: existingCategory.ID,
			setup: func(t *testing.T) {
				// Insert the category but no products
				td.InsertCategory(t, ctx, existingCategory, testPool)
			},
			expectedCount: 0,
			expectedErr:   false,
		},
		{
			name:          "Non-existent category ID should return 0",
			categoryID:    uuid.New(), // Use a fresh, non-existent UUID
			setup:         func(t *testing.T) { /* no category or products inserted for this ID */ },
			expectedCount: 0,
			expectedErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean tables before each sub-test to ensure isolation
			td.ClearCategoriesTable(t, ctx, testPool) // This now also cleans products table

			// Perform test-specific setup
			if tt.setup != nil {
				tt.setup(t)
			}

			// Call the function under test
			count, err := repo.CountProductsInCategory(ctx, tt.categoryID)

			if tt.expectedErr {
				assert.Error(t, err, "Expected an error but got none")
				if tt.expectedErrIs != nil {
					assert.ErrorIs(t, err, tt.expectedErrIs, "Expected specific error type to be in error chain")
				}
				assert.Equal(t, 0, count, "Expected count to be 0 on error")
			} else {
				assert.NoError(t, err, "Did not expect an error but got one")
				assert.Equal(t, tt.expectedCount, count, "Product count mismatch")
			}
		})
	}
}
