package categoryRepository_test

import (
	"context"
	"testing"

	td "github.com/Leviosa-care/catalog/test/helpers"

	"github.com/stretchr/testify/assert"
)

func TestCategoryExistsByName(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name           string
		categoryName   string
		setup          func(*testing.T) // Setup action before running the function
		expectedExists bool
		expectedErr    bool
		expectedErrIs  error // For checking specific error types if expectedErr is true
	}{
		{
			name:         "Existing category should return true",
			categoryName: "TestCategory1",
			setup: func(t *testing.T) {
				cat := td.NewValidCategory("TestCategory1")
				td.InsertCategory(t, ctx, cat, testPool)
			},
			expectedExists: true,
			expectedErr:    false,
		},
		{
			name:           "Non-existing category should return false",
			categoryName:   "NonExistentCategory",
			setup:          func(t *testing.T) { /* no setup, ensure table is empty or doesn't have this name */ },
			expectedExists: false,
			expectedErr:    false,
		},
		{
			name:         "Existing category with different casing should return true (if DB is case-insensitive, else false)",
			categoryName: "testcategory2", // Assuming DB is case-sensitive, this should return false if "TestCategory2" exists
			setup: func(t *testing.T) {
				cat := td.NewValidCategory("TestCategory2") // Insert with exact casing
				td.InsertCategory(t, ctx, cat, testPool)
			},
			expectedExists: false, // PostgreSQL `name = $1` is case-sensitive by default
			expectedErr:    false,
		},
		{
			name:           "Empty string name should return false",
			categoryName:   "",
			setup:          func(t *testing.T) {},
			expectedExists: false,
			expectedErr:    false, // Will return false, nil (no rows found for empty string)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean table before each sub-test to ensure isolation
			td.ClearCategoriesTable(t, ctx, testPool)

			// Perform test-specific setup
			if tt.setup != nil {
				tt.setup(t)
			}

			// Call the function under test
			exists, err := repo.CategoryExistsByName(ctx, tt.categoryName)

			if tt.expectedErr {
				assert.Error(t, err, "Expected an error but got none")
				if tt.expectedErrIs != nil {
					assert.ErrorIs(t, err, tt.expectedErrIs, "Expected specific error type")
				}
				assert.False(t, exists, "Expected exists to be false on error")
			} else {
				assert.NoError(t, err, "Did not expect an error but got one")
				assert.Equal(t, tt.expectedExists, exists, "Existence result mismatch")
			}
		})
	}
}
