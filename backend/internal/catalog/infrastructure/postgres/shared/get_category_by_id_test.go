package sharedRepository_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"
	td "github.com/Leviosa-care/leviosa/backend/test/helpers"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetCategoryByID(t *testing.T) {
	ctx := context.Background()

	categoryID := uuid.New()

	tests := []struct {
		name        string
		setup       func(t *testing.T, ctx context.Context, categoryID uuid.UUID) // setup might need categoryID
		expectedCat *domain.Category
		expectedErr error
	}{
		{
			name: "Successfully retrieve an existing category",
			setup: func(t *testing.T, ctx context.Context, categoryID uuid.UUID) {
				cat := td.NewValidCategory("Test Category One")
				cat.ID = categoryID // Use the generated ID
				cat.Description = "Description for test category one."
				cat.Status = domain.Published
				td.InsertCategory(t, ctx, cat, testPool)
			},
			expectedCat: func() *domain.Category {
				cat := td.NewValidCategory("Test Category One")
				// ID, CreatedAt, UpdatedAt will be populated by setup, so use placeholder values
				// and then copy from actual DB if needed for deep comparison later.
				// For now, these are just what we expect the *non-dynamic* fields to be.
				cat.Description = "Description for test category one."
				cat.Status = domain.Published
				return cat
			}(),
			expectedErr: nil,
		},
		{
			name:        "Return NotFoundErr when category does not exist",
			setup:       func(t *testing.T, ctx context.Context, categoryID uuid.UUID) {}, // No setup needed, ensuring table is empty for this ID
			expectedCat: nil,
			expectedErr: errs.ErrRepositoryNotFound, // This is the error from errs.NewNotFoundErr
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean tables before each sub-test
			td.ClearCategoriesTable(t, ctx, testPool)

			// Perform test-specific setup with the generated ID
			tt.setup(t, ctx, categoryID)

			// Call the function under test
			category, err := repo.GetCategoryByID(ctx, categoryID)

			if tt.expectedErr != nil {
				assert.Error(t, err, "Expected an error but got none")
				assert.ErrorIs(t, err, tt.expectedErr, "Expected specific error type")
				assert.Nil(t, category, "Expected category to be nil on error")
			} else {
				assert.NoError(t, err, "Did not expect an error but got one")
				require.NotNil(t, category, "Expected a category but got nil")

				// For deep comparison, fetch the expected category directly from DB
				// This handles dynamic fields like CreatedAt/UpdatedAt
				var expectedCategoryFromDB domain.Category
				var metadataJSONFromDB []byte

				err = testPool.QueryRow(ctx, `
					SELECT id, name, description, status, metadata, created_at, updated_at
					FROM catalog.categories WHERE id = $1`, categoryID).Scan(
					&expectedCategoryFromDB.ID,
					&expectedCategoryFromDB.Name,
					&expectedCategoryFromDB.Description,
					&expectedCategoryFromDB.Status,
					&metadataJSONFromDB,
					&expectedCategoryFromDB.CreatedAt,
					&expectedCategoryFromDB.UpdatedAt,
				)
				require.NoError(t, err, "Failed to retrieve category directly from DB for comparison")

				if metadataJSONFromDB != nil {
					err := json.Unmarshal(metadataJSONFromDB, &expectedCategoryFromDB.Metadata)
					require.NoError(t, err, "Failed to unmarshal metadata for DB comparison")
				} else {
					expectedCategoryFromDB.Metadata = make(map[string]any)
				}

				// Now compare the retrieved category with the one from DB
				assert.Equal(t, expectedCategoryFromDB.ID, category.ID, "ID mismatch")
				assert.Equal(t, expectedCategoryFromDB.Name, category.Name, "Name mismatch")
				assert.Equal(t, expectedCategoryFromDB.Description, category.Description, "Description mismatch")
				assert.Equal(t, expectedCategoryFromDB.Status, category.Status, "Status mismatch")

				assert.Equal(t, expectedCategoryFromDB.Metadata, category.Metadata, "Metadata mismatch")

				assert.WithinDuration(t, expectedCategoryFromDB.CreatedAt, category.CreatedAt, time.Millisecond, "CreatedAt mismatch")
				assert.WithinDuration(t, expectedCategoryFromDB.UpdatedAt, category.UpdatedAt, time.Millisecond, "UpdatedAt mismatch")
			}
		})
	}
}
