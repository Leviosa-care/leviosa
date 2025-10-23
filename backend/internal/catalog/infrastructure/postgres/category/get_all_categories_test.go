package categoryRepository_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"
	td "github.com/Leviosa-care/catalog/test/helpers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetAllCategories(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		name          string
		setup         func(t *testing.T, ctx context.Context)
		expectedCount int
		expectedOrder []string // Expected order of category names
		expectedErr   error
	}{
		{
			name: "Successfully retrieve multiple categories, ordered by name",
			setup: func(t *testing.T, ctx context.Context) {
				// Insert categories in an arbitrary order to test sorting
				cat1 := td.NewValidCategory("Zeta Category")
				cat2 := td.NewValidCategory("Alpha Category")
				cat3 := td.NewValidCategory("Beta Category")
				cat4 := td.NewValidCategory("Gamma Category") // With imageKey and metadata
				cat4.Metadata = map[string]any{"color": "blue", "size": 10}
				cat4.Status = domain.Published

				td.InsertCategory(t, ctx, cat1, testPool)
				td.InsertCategory(t, ctx, cat2, testPool)
				td.InsertCategory(t, ctx, cat3, testPool)
				td.InsertCategory(t, ctx, cat4, testPool)
			},
			expectedCount: 4,
			expectedOrder: []string{"Alpha Category", "Beta Category", "Gamma Category", "Zeta Category"},
			expectedErr:   nil,
		},
		{
			name:          "Successfully retrieve no categories when table is empty",
			setup:         func(t *testing.T, ctx context.Context) {}, // Table is clean
			expectedCount: 0,
			expectedOrder: []string{},
			expectedErr:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean tables before each sub-test
			td.ClearCategoriesTable(t, ctx, testPool)

			// Perform test-specific setup
			tt.setup(t, ctx)

			// Call the function under test
			categories, err := repo.GetAllCategories(ctx)

			if tt.expectedErr != nil {
				assert.Error(t, err, "Expected an error but got none")
				assert.ErrorIs(t, err, tt.expectedErr, "Expected specific error type")
				assert.Nil(t, categories, "Expected categories slice to be nil on error")
			} else {
				assert.NoError(t, err, "Did not expect an error but got one")
				require.NotNil(t, categories, "Expected categories slice to be non-nil on success")
				assert.Len(t, categories, tt.expectedCount, "Mismatched number of categories")

				// Verify order of names
				actualNames := make([]string, len(categories))
				for i, cat := range categories {
					actualNames[i] = cat.Name
				}
				assert.Equal(t, tt.expectedOrder, actualNames, "Categories not ordered as expected")

				// For the "multiple categories" test, perform deeper checks on data mapping
				if tt.name == "Successfully retrieve multiple categories, ordered by name" {
					// Retrieve all categories directly from DB for comparison, ordered by name
					var dbCategories []*domain.Category
					rows, err := testPool.Query(ctx, `SELECT id, name, description, status, image_key, metadata, created_at, updated_at FROM catalog.categories ORDER BY name ASC`)
					require.NoError(t, err)
					defer rows.Close()

					for rows.Next() {
						var (
							cat          domain.Category
							metadataJSON []byte
							imageKey     sql.NullString
							statusStr    string
						)
						err := rows.Scan(&cat.ID, &cat.Name, &cat.Description, &statusStr, &imageKey, &metadataJSON, &cat.CreatedAt, &cat.UpdatedAt)
						require.NoError(t, err)

						cat.Status = domain.PublishedStatus(statusStr)
						if metadataJSON != nil {
							err := json.Unmarshal(metadataJSON, &cat.Metadata)
							require.NoError(t, err)
						} else {
							cat.Metadata = make(map[string]any) // Match logic in GetAllCategories
						}
						dbCategories = append(dbCategories, &cat)
					}
					require.NoError(t, rows.Err())

					require.Len(t, categories, len(dbCategories), "Length mismatch between retrieved and DB categories")

					for i := range categories {
						retrieved := categories[i]
						expected := dbCategories[i]

						assert.Equal(t, expected.ID, retrieved.ID, "ID mismatch at index %d", i)
						assert.Equal(t, expected.Name, retrieved.Name, "Name mismatch at index %d", i)
						assert.Equal(t, expected.Description, retrieved.Description, "Description mismatch at index %d", i)
						assert.Equal(t, expected.Status, retrieved.Status, "Status mismatch at index %d", i)

						// Metadata comparison (maps can't be directly compared with ==)
						assert.Equal(t, expected.Metadata, retrieved.Metadata, "Metadata mismatch at index %d", i)

						// Time comparison using WithinDuration due to potential precision differences
						assert.WithinDuration(t, expected.CreatedAt, retrieved.CreatedAt, time.Millisecond, "CreatedAt mismatch at index %d", i)
						assert.WithinDuration(t, expected.UpdatedAt, retrieved.UpdatedAt, time.Millisecond, "UpdatedAt mismatch at index %d", i)
					}
				}
			}
		})
	}
}
