package categoryRepository_test

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

func TestUpdateCategory(t *testing.T) {
	ctx := context.Background()

	categoryID := uuid.New()

	tests := []struct {
		name    string
		setup   func(t *testing.T, ctx context.Context, categoryID uuid.UUID) // Setup before calling UpdateCategory
		request *domain.UpdateCategoryRequest
		verify  func(t *testing.T, ctx context.Context, categoryID uuid.UUID, expectedErr error) // Custom verification after update attempt
	}{
		{
			name: "Successfully update a single field (Name)",
			setup: func(t *testing.T, ctx context.Context, categoryID uuid.UUID) {
				cat := td.NewValidCategory("Original Name")
				cat.ID = categoryID
				cat.Description = "Original Description"
				cat.Status = domain.Draft
				td.InsertCategory(t, ctx, cat, testPool)
			},
			request: &domain.UpdateCategoryRequest{
				ID:   categoryID.String(),
				Name: td.StrPtr("Updated Name"),
			},
			verify: func(t *testing.T, ctx context.Context, categoryID uuid.UUID, err error) {
				require.NoError(t, err) // Ensure update itself was successful
				updatedCat, err := td.GetCategoryByID(t, ctx, categoryID, testPool)
				require.NoError(t, err)
				assert.Equal(t, "Updated Name", updatedCat.Name)
				assert.Equal(t, "Original Description", updatedCat.Description, "Description should not have changed")
				assert.Equal(t, domain.Draft, updatedCat.Status, "Status should not have changed")
			},
		},
		{
			name: "Successfully update multiple fields (Name, Description, Status, Metadata)",
			setup: func(t *testing.T, ctx context.Context, categoryID uuid.UUID) {
				cat := td.NewValidCategory("Initial Name")
				cat.ID = categoryID
				cat.Description = "Initial Description"
				cat.Status = domain.Draft
				cat.Metadata = map[string]any{"old_key": "old_value"}
				td.InsertCategory(t, ctx, cat, testPool)
			},
			request: &domain.UpdateCategoryRequest{
				ID:          categoryID.String(),
				Name:        td.StrPtr("New Full Name"),
				Description: td.StrPtr("New Full Description"),
				Status:      td.StatusStrPtr(string(domain.Draft)),
				Metadata:    map[string]any{"new_key": "new_value", "number": 123},
			},
			verify: func(t *testing.T, ctx context.Context, categoryID uuid.UUID, err error) {
				require.NoError(t, err)
				updatedCat, err := td.GetCategoryByID(t, ctx, categoryID, testPool)
				require.NoError(t, err)
				assert.Equal(t, "New Full Name", updatedCat.Name)
				assert.Equal(t, "New Full Description", updatedCat.Description)
				assert.Equal(t, domain.Draft, updatedCat.Status)
				assert.Equal(t, map[string]any{"new_key": "new_value", "number": float64(123)}, updatedCat.Metadata, "Metadata should be updated") // JSON numbers unmarshal to float64
			},
		},
		{
			name: "Update with nil fields in request should not change existing values",
			setup: func(t *testing.T, ctx context.Context, categoryID uuid.UUID) {
				cat := td.NewValidCategory("Existing Name")
				cat.ID = categoryID
				cat.Description = "Existing Description"
				cat.Status = domain.Draft
				cat.Metadata = map[string]any{"initial": "data"}
				td.InsertCategory(t, ctx, cat, testPool)
			},
			request: &domain.UpdateCategoryRequest{
				// All fields other than the ID are nil
				ID: categoryID.String(),
			},
			verify: func(t *testing.T, ctx context.Context, categoryID uuid.UUID, err error) {
				require.Error(t, err) // Expect the error
				// assert.ErrorIs(t, err, fmt.Errorf("no fields provided for update"))
				assert.ErrorIs(t, err, errs.ErrNoFieldsForUpdate)

				// Verify category state is unchanged
				currentCat, getErr := td.GetCategoryByID(t, ctx, categoryID, testPool)
				require.NoError(t, getErr)
				assert.Equal(t, "Existing Name", currentCat.Name)
				assert.Equal(t, "Existing Description", currentCat.Description)
				assert.Equal(t, domain.Draft, currentCat.Status)
				assert.Equal(t, map[string]any{"initial": "data"}, currentCat.Metadata)
			},
		},
		{
			name:  "Update a non-existent category should return NotFound error",
			setup: func(t *testing.T, ctx context.Context, categoryID uuid.UUID) {},
			request: &domain.UpdateCategoryRequest{
				ID:   categoryID.String(),
				Name: td.StrPtr("Non Existent Update"),
			},
			verify: func(t *testing.T, ctx context.Context, categoryID uuid.UUID, err error) {
				require.Error(t, err)
				assert.ErrorIs(t, err, errs.ErrRepositoryNotFound)
			},
		},
		{
			name: "Update with duplicate name should return UniqueViolation error",
			setup: func(t *testing.T, ctx context.Context, categoryID uuid.UUID) {
				// Category to be updated
				catToUpdate := td.NewValidCategory("CategoryA")
				catToUpdate.ID = categoryID
				td.InsertCategory(t, ctx, catToUpdate, testPool)

				// Category with the conflicting name
				conflictingCat := td.NewValidCategory("ConflictingName")
				conflictingCat.ID = uuid.New() // Different ID
				td.InsertCategory(t, ctx, conflictingCat, testPool)
			},
			request: &domain.UpdateCategoryRequest{
				ID:   categoryID.String(),
				Name: td.StrPtr("ConflictingName"), // Attempt to change name to an existing one
			},
			verify: func(t *testing.T, ctx context.Context, categoryID uuid.UUID, err error) {
				require.Error(t, err)
				assert.ErrorIs(t, err, errs.ErrUniqueViolation)
				// Verify original category name did not change
				originalCat, getErr := td.GetCategoryByID(t, ctx, categoryID, testPool)
				require.NoError(t, getErr)
				assert.Equal(t, "CategoryA", originalCat.Name, "Category name should not have changed after failed update")
			},
		},
		{
			name: "Update with invalid status value should return database error (ENUM violation)",
			setup: func(t *testing.T, ctx context.Context, categoryID uuid.UUID) {
				cat := td.NewValidCategory("CategoryForInvalidStatus")
				cat.ID = categoryID
				td.InsertCategory(t, ctx, cat, testPool)
			},
			request: &domain.UpdateCategoryRequest{
				ID:     categoryID.String(),
				Status: td.StatusStrPtr("invalid_status_value"), // Invalid ENUM value
			},
			verify: func(t *testing.T, ctx context.Context, categoryID uuid.UUID, err error) {
				require.Error(t, err)
				assert.ErrorIs(t, err, errs.ErrInvalidInput)
				// Verify original status did not change
				originalCat, getErr := td.GetCategoryByID(t, ctx, categoryID, testPool)
				require.NoError(t, getErr)
				assert.Equal(t, domain.Draft, originalCat.Status, "Category status should not have changed after failed update")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean tables before each sub-test
			td.ClearCategoriesTable(t, ctx, testPool)

			// Perform test-specific setup
			tt.setup(t, ctx, categoryID)

			// Call the function under test
			err := repo.UpdateCategory(ctx, tt.request)

			// Verify using the custom verification function
			tt.verify(t, ctx, categoryID, err)
		})
	}
}
