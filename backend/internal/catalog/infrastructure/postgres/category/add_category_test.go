package categoryRepository_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/catalog/ports"
	td "github.com/Leviosa-care/catalog/test/helpers"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddCategory(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		category    *domain.Category
		wantErr     bool
		expectedErr error // Expected error type if wantErr is true
		setup       func(t *testing.T, ctx context.Context, repo ports.CategoryRepository, cat *domain.Category)
	}{
		{
			name:     "Successfully add a new category with metadata",
			category: td.NewValidCategory("Massage"),
			wantErr:  false,
			setup: func(t *testing.T, ctx context.Context, repo ports.CategoryRepository, cat *domain.Category) {
				cat.Metadata = map[string]any{
					"display_order": 1,
					"is_featured":   true,
				}
			},
		},
		{
			name:     "Successfully add a new category without metadata",
			category: td.NewValidCategory("Mental coaching"),
			wantErr:  false,
			setup: func(t *testing.T, ctx context.Context, repo ports.CategoryRepository, cat *domain.Category) {
				cat.Metadata = nil // Explicitly nil metadata
			},
		},
		{
			name:     "Successfully add a new category with empty metadata",
			category: td.NewValidCategory("Coaching"),
			wantErr:  false,
			setup: func(t *testing.T, ctx context.Context, repo ports.CategoryRepository, cat *domain.Category) {
				cat.Metadata = map[string]any{} // Empty metadata map
			},
		},
		{
			name:     "Attempt to add category with duplicate ID",
			category: td.NewValidCategory("Duplicate name category"),
			wantErr:  true,
			setup: func(t *testing.T, ctx context.Context, repo ports.CategoryRepository, cat *domain.Category) {
				// Insert a category with the same name first
				dupCat := td.NewValidCategory(cat.Name) // Create a new category object with same name
				dupCat.ID = uuid.New()                  // Ensure it has a different ID
				_, err := repo.AddCategory(ctx, dupCat)
				require.NoError(t, err, "Pre-insertion of duplicate name failed")
			},
			expectedErr: errs.ErrUniqueViolation, // Assuming ClassifyPgError returns a DatabaseError for unique violation
		},
		{
			name:        "Attempt to add category with invalid status (ENUM violation)",
			category:    td.NewValidCategory("InvalidStatusCategory"),
			wantErr:     true,
			expectedErr: errs.ErrInvalidInput,
			setup: func(t *testing.T, ctx context.Context, repo ports.CategoryRepository, cat *domain.Category) {
				// Directly manipulate the status string to be invalid
				cat.Status = "non_existent_status_value" // This will cause an ENUM error
			},
		},
		{
			name:     "Successfully add category with default status 'draft'",
			category: td.NewValidCategory("DraftCategory"),
			wantErr:  false,
			setup: func(t *testing.T, ctx context.Context, repo ports.CategoryRepository, cat *domain.Category) {
				cat.Status = domain.Draft
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			td.ClearCategoriesTable(t, ctx, testPool)

			if tt.setup != nil {
				tt.setup(t, ctx, repo, tt.category)
			}
			gotID, err := repo.AddCategory(ctx, tt.category)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.expectedErr != nil {
					assert.ErrorIs(t, err, tt.expectedErr, "Expected specific error type")
				}
				assert.Empty(t, gotID)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.category.ID.String(), gotID, "Returned ID should matth category ID")

				// Verify the category was inserted correctly by querying the database
				var (
					dbID       uuid.UUID
					dbName     string
					dbDesc     string
					dbStatus   string
					dbMetadata []byte
				)

				query := `SELECT id, name, description, status, metadata FROM catalog.categories WHERE id = $1`
				row := testPool.QueryRow(ctx, query, gotID)
				err = row.Scan(&dbID, &dbName, &dbDesc, &dbStatus, &dbMetadata)
				require.NoError(t, err, "Failed to query inserted category")

				assert.Equal(t, tt.category.ID, dbID)
				assert.Equal(t, tt.category.Name, dbName)
				assert.Equal(t, tt.category.Description, dbDesc)
				assert.Equal(t, string(tt.category.Status), dbStatus)

				// Compare metadata
				var actualMetadata map[string]any
				if len(dbMetadata) > 0 { // Check if metadata is not empty
					err = json.Unmarshal(dbMetadata, &actualMetadata)
					require.NoError(t, err, "Failed to unmarshal metadata from DB")
				} else {
					actualMetadata = nil // Treat empty JSON as nil map for comparison
				}

				// Normalize expected metadata for comparison (json.Marshal can reorder keys)
				expectedMetadataJSON, _ := json.Marshal(tt.category.Metadata)
				var expectedMetadata map[string]any
				if len(expectedMetadataJSON) > 0 {
					_ = json.Unmarshal(expectedMetadataJSON, &expectedMetadata)
				} else {
					expectedMetadata = nil
				}

				assert.Equal(t, expectedMetadata, actualMetadata, "Metadata mismatth")
			}
		})
	}
}
